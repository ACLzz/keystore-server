package api

import (
	"fmt"
	"github.com/ACLzz/keystore-server/src/database"
	"github.com/ACLzz/keystore-server/src/errors"
	"github.com/gorilla/mux"
	log "github.com/sirupsen/logrus"
	"net/http"
)

func CreateCollection(w http.ResponseWriter, r *http.Request) {
	body := ConvBody(w, r)
	if body == nil {
		return
	}
	token := VerifyAuth(w, body)
	if token == nil {
		return
	}
	user := token.GetUser()
	_title, ok := body["title"]
	if !ok {
		SendError(w, errors.NoTitleError, 400)
		return
	}
	title := _title.(string)
	if !CheckCollectionTitle(title, w) {
		return
	}
	
	collection := database.Collection{
		Title:     title,
		User:      *user,
	}

	if collection.IsExist() {
		log.Info(fmt.Sprintf("Collection %s for user %d already exist", collection.Title, collection.User.Id))
		SendError(w, errors.CollectionExist, 422)
		return
	}
	collection.Add()
	log.Infof("Created %s collection for %d user", collection.Title, collection.User.Id)

	SendResp(w, nil, 201)
}

func FetchCollections(w http.ResponseWriter, r *http.Request) {}

func ListCollection(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	collection := vars["collection"]
	log.Info("Listing ", collection, " collection")
}

func UpdateCollection(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	collection := vars["collection"]
	log.Info("Updating ", collection, " collection")
}

func DeleteCollection(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	collection := vars["collection"]
	log.Info("Delete ", collection, " collection")
}