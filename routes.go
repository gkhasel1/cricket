package main

import "net/http"

type Route struct {
	Name        string
	Version     int
	Method      string
	Pattern     string
	HandlerFunc http.HandlerFunc
}

type Routes []Route

var routes = Routes{
	// v1 routes
	Route{
		"PostMetrics",
		1,
		"POST",
		"/metrics",
		PostMetricsHandler,
	},
	Route{
		"GetMetrics",
		1,
		"GET",
		"/metrics",
		GetMetricsHandler,
	},
	Route{
		"GetStats",
		1,
		"GET",
		"/stats",
		GetStatsHandler,
	},
	// global routes
	Route{
		"IndexHandler",
		0,
		"GET",
		"/",
		IndexHandler,
	},
	Route{
		"Health",
		0,
		"GET",
		"/health",
		HealthHandler,
	},
}
