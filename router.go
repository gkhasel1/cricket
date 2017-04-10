package main

import (
	"net/http"

	mux "github.com/gorilla/mux"
)

func NewRouter() *mux.Router {
	/**
	 * NewRouter: Creates a Mux Router for our API. There
	 *		are global and v1 routes in routes.go, based on our
	 *		versioning scheme we assign routes to the correct router
	 *      or subrouter. Later we could add a v2 router.
	 *
	 *		This router uses HTTP, some metrics implementations like
	 *		statsd use UDP packets. This may actually be better for
	 *		a service like this and is worth exploring.
	 **/
	router := mux.NewRouter().StrictSlash(true)
	v1 := router.PathPrefix("/v1").Subrouter()

	for _, route := range routes {
		var handler http.Handler
		handler = route.HandlerFunc

		// Use our route handlers with our convenient logger
		handler = Logger(handler, route.Name)

		var assignedRouter *mux.Router
		switch route.Version {
		case 0: // global version
			assignedRouter = router
		case 1: // version 1
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
