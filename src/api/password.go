package api

import (
	"github.com/ACLzz/keystore-server/src/database"
	"github.com/ACLzz/keystore-server/src/errors"
	"github.com/gorilla/mux"
	log "github.com/sirupsen/logrus"
	"net/http"
)

func CreatePassword(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	collection := vars["collection"]
	body := ConvBody(w, r)
	if body == nil {
		return
	}
	if !VerifyAuth(w, body) {
		return
	}
	token := GetTokenObj(body["token"].(string))
	if token == nil {
		return
	}
	user := token.GetUser()

	_login, ok := body["login"]
	if !ok {
		SendError(w, errors.NoLoginError, 400)
		return
	}
	_password, ok := body["password"]
	if !ok {
		SendError(w, errors.NoPasswordError, 400)
		return
	}
	_title, ok := body["title"]
	if !ok {
		SendError(w, errors.NoTitleError, 400)
		return
	}

	login := _login.(string)
	password := _password.(string)
	title := _title.(string)
	email := ""

	_email, ok := body["email"]
	if ok {
		email = _email.(string)
	}

	passwd := database.Password{
		Title: title,
		Login: login,
		Password: password,
		Email: email,
	}

	if !passwd.Add(collection, user.Id) {
		SendError(w, errors.InternalError, 500)
		return
	}
	log.Infof("Created %s password in collection %s for %d user", passwd.Title, passwd.Collection.Title, user.Id)
	
	SendResp(w, nil, 201)
}

func ReadPassword(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	pid := vars["pid"]
	log.Info("Getting ", pid, " password")
}

func UpdatePassword(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	pid := vars["pid"]
	log.Info("Updating ", pid, " password")
}

func DeletePassword(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	pid := vars["pid"]
	log.Info("Delete ", pid, " password")
}
