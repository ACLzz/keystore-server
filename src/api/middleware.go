package api

import (
	"github.com/ACLzz/keystore-server/src/database"
	"github.com/ACLzz/keystore-server/src/errors"
	"github.com/gorilla/mux"
	"net/http"
	"strconv"
	"strings"
)


func PasswordMiddleware(h http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if !VerifyAuth(w, r) {
			return
		}
		user := GetUser(GetToken(r))

		vars := mux.Vars(r)
		collection := vars["collection"]
		coll := database.Collection{Title: collection, User: *user}
		if !coll.IsExist() {
			SendError(w, errors.CollectionNotExist, 404)
			return
		}
		
		pid, ok := vars["pid"]
		if ok {
			_, err := strconv.Atoi(pid)
			if err != nil {
				SendError(w, errors.InvalidPasswordId, 422)
				return
			}
		}
		
		h.ServeHTTP(w, r)
	})
}

func CollectionMiddleWare(h http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if !VerifyAuth(w, r) {
			return
		}
		h.ServeHTTP(w, r)
	})
}

func TrailingMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		r.URL.Path = strings.TrimSuffix(r.URL.Path, "/")
		next.ServeHTTP(w, r)
	})
}
