package main

import (
	"net/http"

	mux "github.com/gorilla/mux"
)

func NewRouter() *mux.Router {

	router := mux.NewRouter().StrictSlash(true)
	v1 := router.PathPrefix("/v1").Subrouter()

	for _, route := range routes {
		var handler http.Handler

		handler = route.HandlerFunc
		handler = Logger(handler, route.Name)

		var assignedRouter *mux.Router
		switch route.Version {
		case 0:
			assignedRouter = router
		case 1:
			assignedRouter = v1
		}

		assignedRouter.
			Methods(route.Method).
			Path(route.Pattern).
			Name(route.Name).
			Handler(handler)
	}

	return router
}
