package api

import (
	"github.com/ACLzz/keystore-server/src/utils"
	"github.com/gorilla/mux"
	log "github.com/sirupsen/logrus"
	"net/http"
)

type routesMap []routeType

type routeType struct {
	Route	string
	Handler func(http.ResponseWriter, *http.Request)
	Methods []string
}


func MainRouter() *mux.Router {
	log.Info("Initializing main router")
	r := mux.NewRouter()

	AuthRouter(r)
	CollectionRouter(r)
	if utils.Config.Test {
		DevRouter(r)
	}

	return r
}

func AuthRouter(parent *mux.Router) *mux.Router {
	log.Info("Initializing auth router")
	return buildRouter(parent, "/auth/", routesMap{
		{"/login", Login, []string{"POST"}},
		{"/", ReadUser, []string{"GET"}},
		{"/", Register, []string{"POST"}},
		{"/", UpdateUser, []string{"PUT"}},
		{"/", DeleteUser, []string{"DELETE"}},
	})
}

func CollectionRouter(parent *mux.Router) *mux.Router {
	log.Info("Initializing collections router")
	r := buildRouter(parent, "/collection/", routesMap{
		{"/", FetchCollections, []string{"GET"}},
		{"/", CreateCollection, []string{"POST"}},
		{"/{collection}", ListCollection, []string{"GET"}},
		{"/{collection}", UpdateCollection, []string{"PUT"}},
		{"/{collection}", DeleteCollection, []string{"DELETE"}},
	})
	
	PasswordRouter(r)
	return r
}

func PasswordRouter(parent *mux.Router) *mux.Router {
	log.Info("Initializing passwords router")
	
	return buildRouter(parent, "/pswd/", routesMap{
		{"/", CreatePassword, []string{"POST"}},
		{"/{pid}", ReadPassword, []string{"GET"}},
		{"/{pid}", UpdatePassword, []string{"PUT"}},
		{"/{pid}", DeletePassword, []string{"DELETE"}},
	})
}

func DevRouter(parent *mux.Router) *mux.Router {
	log.Info("Initializing dev router")
	return buildRouter(parent, "/dev/", routesMap{
		{"/shutdown-server", ShutdownServer, []string{"GET"}},
	})
}

func buildRouter(parent *mux.Router, path string, routes routesMap) *mux.Router {
	/*
	parent: parent router
	path: sub-path for router
	routes: map of sub-route
	 */
	r := parent.PathPrefix(path).Subrouter()
	for _, route := range routes {
		for _, method := range route.Methods {
			r.HandleFunc(route.Route, route.Handler).Methods(method)
		}
	}
	return r
}
