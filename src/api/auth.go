package api

import (
	"fmt"
	"github.com/ACLzz/keystore-server/src/database"
	"github.com/ACLzz/keystore-server/src/errors"
	"github.com/gorilla/mux"
	log "github.com/sirupsen/logrus"
	"net/http"
)

func Register(w http.ResponseWriter, r *http.Request) {
	body := ConvBody(w, r)
	if body == nil {
		return
	}
	
	if isBodyValid := checkAuthFields(body, w); !isBodyValid {
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
	SendResp(w, map[string]interface{}{"ok": true}, 201)
}

func Login(w http.ResponseWriter, r *http.Request) {}

func ReadUser(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	uid := vars["uid"]
	log.Info("Getting info for ", uid, " user")
}

func UpdateUser(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	uid := vars["uid"]
	log.Info("Updating info for ", uid, " user")
}

func DeleteUser(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	uid := vars["uid"]
	log.Info("Delete ", uid, " user")
}

func checkAuthFields(body map[string]interface{}, w http.ResponseWriter) bool {
	// Check if fields has been sent
	if _, ok := body["login"]; !ok {
		SendError(w, errors.NoLoginError, 400)
		return false
	} else if _, ok := body["password"]; !ok {
		SendError(w, errors.NoPasswordError, 400)
		return false
	}

	// Check fields on limits
	// TODO min length
	if len(body["login"].(string)) > errors.LoginMaxLengthLimit {
		SendError(w, errors.LoginMaxLengthError, 422)
		return false
	}
	if len(body["password"].(string)) > errors.PasswordMaxLengthLimit {
		SendError(w, errors.PasswordMaxLengthError, 422)
		return false
	}

	if len(body["login"].(string)) < errors.LoginMinLengthLimit {
		SendError(w, errors.LoginMinLengthError, 422)
		return false
	}
	if len(body["password"].(string)) < errors.PasswordMinLengthLimit {
		SendError(w, errors.PasswordMinLengthError, 422)
		return false
	}

	// Check if login contains non-ascii chars
	for _, run := range body["login"].(string) {
		if (run <= 125 && run >= 65) || (run <= 57 && run >= 48) {
			continue
		}
		SendError(w, errors.LoginLocaleError, 422)
		return false
	}

	// Check if password contains non-ascii chars
	for _, run := range body["password"].(string) {
		if (run <= 125 && run >= 65) || (run <= 57 && run >= 48) {
			continue
		}
		SendError(w, errors.PasswordLocaleError, 422)
		return false
	}
	
	return true
}