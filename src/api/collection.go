package api

import (
	"github.com/ACLzz/keystore-server/src/database"
	"github.com/ACLzz/keystore-server/src/errors"
	"github.com/gorilla/mux"
	log "github.com/sirupsen/logrus"
	"gorm.io/gorm"
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
		log.Infof("Collection %s for user %d already exist", collection.Title, collection.User.Id)
		SendError(w, errors.CollectionExist, 422)
		return
	}
	collection.Add()
	log.Infof("Created %s collection for %d user", collection.Title, collection.User.Id)

	SendResp(w, nil, 201)
}

func FetchCollections(w http.ResponseWriter, r *http.Request) {
	body := ConvBody(w, r)
	if body == nil {
		return
	}
	token := VerifyAuth(w, body)
	if token == nil {
		return
	}
	collections := token.FetchCollections()
	var collTitles []interface{}
	for _, coll := range collections {
		collTitles = append(collTitles, coll.Title)
	}

	SendArray(w, &collTitles, 200)
}

func ListCollection(w http.ResponseWriter, r *http.Request) {
	// TODO
	vars := mux.Vars(r)
	collection := vars["collection"]
	log.Info("Listing ", collection, " collection")
}

func UpdateCollection(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	title := vars["collection"]

	body := ConvBody(w, r)
	if body == nil {
		return
	}
	token := VerifyAuth(w, body)
	if token == nil {
		return
	}
	user := token.GetUser()
	if !CheckCollectionTitle(title, w) {
		return
	}

	conn := database.GetConn()
	DB, _ := conn.DB()
	defer DB.Close()
	var coll database.Collection
	if tx := conn.Where("title = ? and user_refer = ?", title, user.Id).First(&coll); tx.Error != nil {
		if tx.Error == gorm.ErrRecordNotFound {
			log.Errorf("collections with %s title for user %d not found", title, user.Id)
			SendError(w, errors.CollectionNotExist, 404)
		} else {
			log.Error(tx.Error)
			SendError(w, errors.InternalError, 500)
		}
		return
	}
	coll.User = *user
	
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
	collection := vars["collection"]
	log.Info("Delete ", collection, " collection")
}