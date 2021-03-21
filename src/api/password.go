package api

import (
	"github.com/gorilla/mux"
	log "github.com/sirupsen/logrus"
	"net/http"
)

func CreatePassword(w http.ResponseWriter, r *http.Request) {}

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
