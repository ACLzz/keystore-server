package api

import (
	"github.com/ACLzz/keystore-server/src/database"
	"github.com/ACLzz/keystore-server/src/errors"
	"github.com/gorilla/mux"
	"net/http"
)


func PasswordMiddleware(h http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		body := ConvBody(w, r)
		if body == nil {
			return
		}

		if VerifyAuth(w, body) == nil {
			return
		}
		token := GetTokenObj(body["token"].(string))

		vars := mux.Vars(r)
		collection := vars["collection"]
		coll := database.Collection{Title: collection, User: token.User}
		if !coll.IsExist() {
			SendError(w, errors.CollectionNotExist, 404)
			return
		}
		
		h.ServeHTTP(w, r)
	})
}

func CollectionMiddleWare(h http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if VerifyAuth(w, ConvBody(w, r)) == nil {
			return
		}
		h.ServeHTTP(w, r)
	})
}
