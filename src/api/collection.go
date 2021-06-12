package api

import (
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
	user := GetUser(GetToken(r))
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
		log.Infof("Collection %s for user %d already exist", collection.Title, collection.User.Id)
		SendError(w, errors.CollectionExist, 422)
		return
	}
	collection.Add()
	log.Infof("Created %s collection for %d user", collection.Title, collection.User.Id)

	SendResp(w, nil, 201)
}

func FetchCollections(w http.ResponseWriter, r *http.Request) {
	user := GetUser(GetToken(r))
	if user == nil {
		return
	}
	collections := user.FetchCollections()
	var collTitles []interface{}
	for _, coll := range collections {
		collTitles = append(collTitles, coll.Title)
	}

	SendArray(w, &collTitles, 200)
}

func ListCollection(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	collection := vars["collection"]
	user := GetUser(GetToken(r))
	if user == nil {
		return
	}
	coll := GetCollection(collection, user)
	if coll == nil {
		SendError(w, errors.CollectionNotExist, 404)
		return
	}

	log.Info("Listing ", collection, " collection for ", user.Id, " user")
	passwords := make([]interface{}, 0)
	_passwords := coll.FetchPasswords()
	if _passwords == nil {
		SendError(w, errors.InternalError, 500)
		return
	}
	for _, pass := range _passwords {
		passwords = append(passwords, map[string]interface{}{"id": pass.Id, "title": pass.Title})
	}
	
	SendArray(w, &passwords, 200)
}

func UpdateCollection(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	title := vars["collection"]

	body := ConvBody(w, r)
	if body == nil {
		return
	}
	user := GetUser(GetToken(r))
	if !CheckCollectionTitle(title, w) {
		return
	}

	coll := database.Collection{
		Title: title,
		User: *user,
	}
	if !coll.IsExist() {
		log.Errorf("collections with %s title for user %d not found", title, user.Id)
		SendError(w, errors.CollectionNotExist, 404)
		return
	}
	
	if _title, ok := body["title"]; ok {
		if !CheckCollectionTitle(_title.(string), w) {
			return
		}
		coll.Title = _title.(string)
		if coll.IsExist() {
			log.Infof("Collection %s for user %d already exist", coll.Title, coll.User.Id)
			SendError(w, errors.CollectionExist, 422)
			return
		}
	}

	log.Info("Updating ", title, " collection")
	if !coll.Update() {
		SendError(w, errors.InternalError, 500)
	} else {
		SendResp(w, nil, 200)
	}
}

func DeleteCollection(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	title := vars["collection"]

	user := GetUser(GetToken(r))
	if !CheckCollectionTitle(title, w) {
		return
	}
	coll := database.Collection{
		Title: title,
		User: *user,
	}

	if !coll.IsExist() {
		log.Errorf("collections with %s title for user %d not found", title, user.Id)
		SendError(w, errors.CollectionNotExist, 404)
		return
	}
	
	log.Info("Delete ", title, " collection")
	if !coll.Delete() {
		SendError(w, errors.InternalError, 500)
	} else {
		SendResp(w, nil, 200)
	}
}