package api

import (
	"encoding/json"
	"fmt"
	"github.com/ACLzz/keystore-server/src/database"
	"github.com/ACLzz/keystore-server/src/errors"
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


func ConvBody(w http.ResponseWriter, r *http.Request) map[string]interface{} {
	defer r.Body.Close()
	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		log.Errorf("Error in getting body, Register: %v\n", err)
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
	return jbody
}

func VerifyAuth(w http.ResponseWriter, body map[string]interface{}) *database.Token {
	token, ok := body["token"]
	if !ok || token == "" {
		SendError(w, errors.NoToken, 400)
		return nil
	}

	// Check token for non-digit and non-letters symbols
	for _, run := range token.(string) {
		if (run <= 90 && run >= 65) || (run <= 57 && run >= 48) || (run <= 122 && run >= 97) {
			continue
		}
		SendError(w, errors.TokenDeniedSymbolsError, 422)
		log.Warn("Strange token was sent: ", token.(string))
		return nil
	}

	conn := database.GetConn()
	DB, _ := conn.DB()
	defer DB.Close()

	dbToken := database.Token{}
	tx := conn.First(&dbToken, "token = ?", token)
	if tx.Error == gorm.ErrRecordNotFound {
		SendError(w, errors.InvalidToken, 422)
		log.Warn("Invalid token request with token ", token)
		return nil
	} else if tx.Error != nil {
		SendError(w, errors.InternalError, 500)
		log.Errorf("In ValidateAuth: %v", tx.Error)
		return nil
	}

	if dbToken.ExpireDate.Before(time.Now()) {
		SendError(w, errors.ExpiredToken, 403)
		return nil
	}

	return &dbToken
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
		if len(login.(string)) > errors.LoginMaxLengthLimit {
			SendError(w, errors.LoginMaxLengthError, 422)
			return false
		}
		if len(login.(string)) < errors.LoginMinLengthLimit {
			SendError(w, errors.LoginMinLengthError, 422)
			return false
		}

		// Check if login contains non-ascii chars
		for _, run := range login.(string) {
			if (run <= 125 && run >= 65) || (run <= 57 && run >= 48) {
				continue
			}
			SendError(w, errors.LoginLocaleError, 422)
			return false
		}
	}

	if password, ok := body["password"]; ok {
		if len(password.(string)) > errors.PasswordMaxLengthLimit {
			SendError(w, errors.PasswordMaxLengthError, 422)
			return false
		}
		if len(password.(string)) < errors.PasswordMinLengthLimit {
			SendError(w, errors.PasswordMinLengthError, 422)
			return false
		}

		// Check if password contains non-ascii chars
		for _, run := range password.(string) {
			if (run <= 125 && run >= 65) || (run <= 57 && run >= 48) {
				continue
			}
			SendError(w, errors.PasswordLocaleError, 422)
			return false
		}
	}

	return true
}
