package api

import (
	"github.com/ACLzz/keystore-server/src/database"
	"github.com/ACLzz/keystore-server/src/errors"
	"github.com/gorilla/mux"
	log "github.com/sirupsen/logrus"
	"gorm.io/gorm"
	"net/http"
)

func CreatePassword(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	collection := vars["collection"]
	body := ConvBody(w, r)
	if body == nil {
		return
	}
	user := GetUser(body["token"].(string))
	if user == nil {
		return
	}

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
	body := ConvBody(w, r)
	if body == nil {
		return
	}
	user := GetUser(body["token"].(string))
	if user == nil {
		return
	}
	collection := GetCollection(vars["collection"], user)

	log.Info("Getting ", pid, " password from ", collection.Title, " collection for ", user.Id, " user")

	conn := database.GetConn()
	DB, _ := conn.DB()
	defer DB.Close()
	var password database.Password

	if tx := conn.Where("id = ? and collection_refer = ?", pid, collection.Id).First(&password); tx.Error != nil {
		if tx.Error == gorm.ErrRecordNotFound {
			SendError(w, errors.PasswordNotExist, 404)
		} else {
			log.Error(tx.Error)
			SendError(w, errors.InternalError, 500)
		}
		return
	}

	SendResp(w, &map[string]interface{}{
		"title": password.Title, "login": password.Login, "password": password.Password, "email": password.Email,
	}, 200)
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
