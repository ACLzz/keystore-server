package api

import (
	"fmt"
	"github.com/ACLzz/keystore-server/src/database"
	"github.com/ACLzz/keystore-server/src/errors"
	log "github.com/sirupsen/logrus"
	"net/http"
)

func Register(w http.ResponseWriter, r *http.Request) {
	body := ConvBody(w, r)
	if body == nil {
		return
	}
	
	if isBodyValid := CheckAuthFields(body, w); !isBodyValid {
		return
	}

	user := database.User{Username: body["login"].(string), Password: body["password"].(string)}
	if err := user.Register(); err != nil {
		log.Info(fmt.Sprintf("User %s already exist", body["login"]))
		SendError(w, errors.UserExists, 422)
		return
	}
	log.Info(fmt.Sprintf("Registered %s user", body["login"]))

	// Response that all ok
	SendResp(w, nil, 201)
}

func Login(w http.ResponseWriter, r *http.Request) {
	body := ConvBody(w, r)
	if body == nil {
		return
	}

	if isBodyValid := CheckAuthFields(body, w); !isBodyValid {
		return
	}

	user := database.User{Username: body["login"].(string), Password: body["password"].(string)}
	if !user.CheckPassword() {
		SendError(w, errors.InvalidCredentials, 401)
		return
	}
	conn := database.GetConn()
	DB, _ := conn.DB()
	defer DB.Close()
	conn.Where("username = ? and password = ?", user.Username, user.HashPassword()).First(&user)

	log.Infof("User %s logged in", user.Username)
	token := user.GenToken()
	SendResp(w, &map[string]interface{}{"token": token}, 202)
}

func ReadUser(w http.ResponseWriter, r *http.Request) {
	if !VerifyAuth(w, r) {
		return
	}

	user := GetUser(GetToken(r))
	if user == nil {
		SendError(w, errors.UserNotExists, 404)
		return
	}
	log.Info("Getting info for ", user.Id, " user")

	SendResp(w, &map[string]interface{}{
		"username": user.Username,
		"registered": user.RegistrationDate.Format("2006-01-02T15:04:05Z07:00"),
	}, 200)
}

func UpdateUser(w http.ResponseWriter, r *http.Request) {
	body := ConvBody(w, r)
	if body == nil {
		return
	}
	if !VerifyAuth(w, r) {
		return
	}

	if isBodyValid := CheckAuthFieldsLimits(body, w); !isBodyValid {
		return
	}

	user := GetUser(GetToken(r))
	if user == nil {
		SendError(w, errors.UserNotExists, 404)
		return
	}

	log.Info("Updating info for ", user.Id, " user")
	if username, ok := body["login"]; ok {
		user.Username = username.(string)
	}
	if password, ok := body["password"]; ok {
		user.Password = password.(string)
		user.Password = user.HashPassword()
	}
	
	if user.Update() {
		SendResp(w, nil, 200)
	} else {
		SendError(w, errors.InternalError, 500)
	}
}

func DeleteUser(w http.ResponseWriter, r *http.Request) {
	if !VerifyAuth(w, r) {
		return
	}
	user := GetUser(GetToken(r))
	if user == nil {
		SendError(w, errors.UserNotExists, 404)
		return
	}

	log.Info("Delete ", user.Id, " user")
	if user.Delete() {
		SendResp(w, nil, 200)
	} else {
		SendError(w, errors.InternalError, 500)
	}
}