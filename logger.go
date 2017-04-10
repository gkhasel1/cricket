package main

import (
	"log"
	"net/http"
	"time"
)

func Logger(inner http.Handler, name string) http.Handler {
	/**
	 * Logger: This function wraps a handler and provides some nice
	 *		event logging for the server. Example output can be seen
	 *      in the README. We could make this logger better by adding
	 *      http response status codes or using a middleware framework
	 *      that supports good logging such as Negroni.
	 **/
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()

		inner.ServeHTTP(w, r)

		log.Printf(
			"%s\t%s\t%s\t%s",
			r.Method,
			r.RequestURI,
			name,
			time.Since(start),
		)
	})
}
