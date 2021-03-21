package api

import (
	"github.com/gorilla/mux"
	log "github.com/sirupsen/logrus"
	"net/http"
)

func CreateCollection(w http.ResponseWriter, r *http.Request) {}

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