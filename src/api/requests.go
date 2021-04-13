package api

import (
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/ACLzz/keystore-server/src/database"
	"github.com/ACLzz/keystore-server/src/errors"
	"github.com/ACLzz/keystore-server/src/utils"
	log "github.com/sirupsen/logrus"
	"gorm.io/gorm"
	"io/ioutil"
	"net/http"
	"time"
)

func SendError(w http.ResponseWriter, err error, errCode int) {
	msg := fmt.Sprintf("{\"error\":\"%s\"}", err.Error())
	http.Error(w, msg, errCode)
}

func SendResp(w  http.ResponseWriter, resp *map[string]interface{}, statusCode int) {
	if resp != nil {
		jresp, err := json.Marshal(resp)
		if err != nil {
			log.Error("Error in json marshalling: ", err)
			SendError(w, errors.InternalError, 500)
			return
		}
		w.WriteHeader(statusCode)
		if _, err := w.Write(jresp); err != nil {
			log.WithFields(log.Fields{"err": err}).Error("Error in send response.")
		}
	} else {
		w.WriteHeader(statusCode)
	}
}

func SendArray(w  http.ResponseWriter, resp *[]interface{}, statusCode int) {
	if resp != nil {
		jresp, err := json.Marshal(resp)
		if err != nil {
			log.Error("Error in json marshalling: ", err)
			SendError(w, errors.InternalError, 500)
			return
		}
		w.WriteHeader(statusCode)
		if _, err := w.Write(jresp); err != nil {
			log.WithFields(log.Fields{"err": err}).Error("Error in send response.")
		}
	} else {
		w.WriteHeader(statusCode)
	}
}


func ConvBody(w http.ResponseWriter, r *http.Request) map[string]interface{} {
	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		log.Errorf("Error in getting body %v\n", err)
		SendError(w, errors.InvalidBodyError, 400)
		return nil
	}

	if len(string(body[:])) <= 2 {
		SendError(w, errors.EmptyBodyError, 400)
		return nil
	}

	jbody := make(map[string]interface{})
	err = json.Unmarshal(body, &jbody)
	if err != nil {
		log.Error("Error in body unmarshalling: ", err)
		SendError(w, errors.InvalidBodyError, 400)
		return nil
	}

	r.Body = ioutil.NopCloser(bytes.NewBuffer(body))
	return jbody
}

func VerifyAuth(w http.ResponseWriter, body map[string]interface{}) bool {
	token, ok := body["token"]
	if !ok || token == "" {
		SendError(w, errors.NoToken, 400)
		return false
	}

	// Check token for non-digit and non-letters symbols
	for _, run := range token.(string) {
		if (run <= 90 && run >= 65) || (run <= 57 && run >= 48) || (run <= 122 && run >= 97) {
			continue
		}
		SendError(w, errors.TokenDeniedSymbolsError, 422)
		log.Warn("Strange token was sent: ", token.(string))
		return false
	}

	conn := database.GetConn()
	DB, _ := conn.DB()
	defer DB.Close()

	dbToken := database.Token{}
	tx := conn.Where("token = ?", token).First(&dbToken)
	if tx.Error == gorm.ErrRecordNotFound {
		SendError(w, errors.InvalidToken, 422)
		log.Warn("Invalid token request with token ", token)
		return false
	} else if tx.Error != nil {
		SendError(w, errors.InternalError, 500)
		log.Errorf("In ValidateAuth: %v", tx.Error)
		return false
	}

	if dbToken.ExpireDate.Before(time.Now()) {
		SendError(w, errors.ExpiredToken, 403)
		return false
	}

	return true
}

func GetUser(token string) *database.User {
	conn := database.GetConn()
	DB, _ := conn.DB()
	defer DB.Close()
	
	var dbToken database.Token
	if tx := conn.Where("token = ?", token).First(&dbToken); tx.Error == gorm.ErrRecordNotFound {
		log.Warn("Can't find token ", token)
		return nil
	}

	return dbToken.GetUser()
}

func GetCollection(title string, user *database.User) *database.Collection {
	conn := database.GetConn()
	DB, _ := conn.DB()
	defer DB.Close()
	
	var collection database.Collection
	if tx := conn.Where("title = ? and user_refer = ?", title, user.Id).First(&collection); tx.Error == gorm.ErrRecordNotFound {
		return nil
	}
	return &collection
}

func GetPassword(collection string, user *database.User, pid int) (*database.Password, error) {
	coll := GetCollection(collection, user)
	conn := database.GetConn()
	DB, _ := conn.DB()
	defer DB.Close()
	var password database.Password

	if tx := conn.Where("id = ? and collection_refer = ?", pid, coll.Id).First(&password); tx.Error != nil {
		return nil, tx.Error
	}
	return &password, nil
}

func CheckAuthFields(body map[string]interface{}, w http.ResponseWriter) bool {
	if _, ok := body["login"]; !ok {
		SendError(w, errors.NoLoginError, 400)
		return false
	} else if _, ok := body["password"]; !ok {
		SendError(w, errors.NoPasswordError, 400)
		return false
	}

	return CheckAuthFieldsLimits(body, w)
}

func CheckAuthFieldsLimits(body map[string]interface{}, w http.ResponseWriter) bool {
	if login, ok := body["login"]; ok {
		if err := CheckLimits(body["login"].(string), utils.LoginMinLengthLimit, utils.LoginMaxLengthLimit,
			errors.LoginMinLengthError, errors.LoginMaxLengthError); err != nil {
			SendError(w, err, 422)
			return false
		}
		
		if !CheckNonAscii(login.(string)) {
			SendError(w, errors.LoginLocaleError, 422)
			return false
		}
	}

	if password, ok := body["password"]; ok {
		if err := CheckLimits(body["password"].(string), utils.PasswordMinLengthLimit, utils.PasswordMaxLengthLimit,
			errors.PasswordMinLengthError, errors.PasswordMaxLengthError); err != nil {
			SendError(w, err, 422)
			return false
		}
		
		if !CheckNonAscii(password.(string)) {
			SendError(w, errors.PasswordLocaleError, 422)
			return false
		}
	}

	return true
}

func CheckLimits(field string, lowLimit int, highLimit int, lowLimitError error, highLimitError error) error {
	if len(field) > highLimit {
		return highLimitError
	}
	if len(field) < lowLimit {
		return lowLimitError
	}
	return nil
}

func CheckNonAscii(field string) bool {
	for _, run := range field {
		if (run <= 125 && run >= 65) || (run <= 57 && run >= 48) {
			continue
		}
		return false
	}
	return true
}

func CheckCollectionTitle(title string, w http.ResponseWriter) bool {
	if err := CheckLimits(title, utils.CollectionTitleMinLengthLimit, utils.CollectionTitleMaxLengthLimit,
		errors.CollectionTitleMinLengthError, errors.CollectionTitleMaxLengthError); err != nil {
		SendError(w, err, 422)
		return false
	}

	if !CheckNonAscii(title) {
		SendError(w, errors.CollectionLocaleError, 422)
		return false
	}
	return true
}
