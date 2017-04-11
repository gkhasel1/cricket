# CRICKET

A Metrics Server written in Go, running on Elasticsearch.

## API DOCUMENTATION

### Global Routes

`GET /` -- index page, shows the title of the project

`GET /health` -- returns 200: OK when the app is running (shallow healthcheck)

### Version 1 Routes

`POST /v1/metrics` -- accepts json data as shown below and creates metric datapoints

Example:
```
$ curl -X POST -H 'Content-Type: application/json' -d '[
	{
		"name": "c1",
		"type": "counter",
		"value": 1,
		"timestamp": "2016-12-31T21:49:38Z"
	},
	{
		"name": "c2",
		"type": "counter",
		"value": 4.7,
		"timestamp": "2016-12-31T22:49:50Z"
	}
]' http://localhost:8080/v1/metrics
```
Timestamps passed to the metrics server should follow `RFC3339` standards.

`GET /v1/metrics` -- returns metric data, accepts metric `name` and `timestamp` as optional query params.
When neither is specified, all results are returned (up to 1000 records just for safety).

Example:
```
curl -X GET http://localhost:8080/v1/metrics\?timestamp\=1990-12-31T22%3A49%3A50Z\&name\=c2 (url encoded timestamp)
[
  	{
	    "name": "c2",
	    "type": "counter",
	    "value": 1,
	    "timestamp": "1990-12-31T22:49:50Z"
	},
	{
	    "name": "c2",
	    "type": "counter",
	    "value": 1,
	    "timestamp": "1990-12-31T22:49:50Z"
	}
]
```

`GET /v1/stats` -- returns a set of statistics about a metric, accepts metric `name`, `start`, and `end`
as required query params

Example:
```
curl -X GET http://localhost:8080/v1/stats\?start\=1980-12-31T22%3A49%3A50Z\&end\=2016-12-31T22%3A49%3A50Z\&name\=c1
{
	"count": 3,
	"sum": 3,
	"min": 0,
	"max": 1,
	"average": 1
}
```
In the example above, there were 3 c1 records written with a timestamp in 1990. Each record has 1 as the value.

## DESIGN DECISIONS

I am neither a Go programmer nor an Elaticsearch expert, but I thought it would be fun and challenging
to learn something new while doing this project.

Go seemed like a suitable language for this project given its performance characteristics and ease of parallelization
via go routines. Ultimately, this project doesn't take advantage of routines, but its an easy optimization to
make in the future.

Elasticsearch was chosen as the backend for the timeseries data being collected. Its indexing strategy makes
it very fast for these types of work loads, additionally it is a common backend for many production metrics setups.

Locally you can develop with go directly. I have provided a Dockerfile and docker-compose.yml to
illustrate a non local setup. In production we would potentially use nginx as a web proxy in front of the
metrics server itself.

## FUTURE WORK

Performance optimization via parrallel goroutines, better elasticsearch queries, and more api's would be some
obvious steps to improving this project.

Given more time we could also implement better validation of inputs to ensure data quality.

Additionally, we could use different libraries for database access, routing etc. or even switch to a
more robust web framework.

Lastly, some visualizations and dashboards on top of this data could be useful.

## DEVELOPMENT SETUP

If you already have the go toolchain setup you can skip this!
Otherwise this may serve as a very abridged getting started guide.

[1] install `go` (https://golang.org/dl/)
    install `docker` and `docker-compose` (https://docs.docker.com/engine/installation/)

[2] create the following directory tree in your home directory (`~`)
```
go/
├── bin/
├── pkg/
└── src/
    └── github.com/
    	└── gkhasel1/
```

[3] Setup your environment variables for the go toolchain
```
export GOPATH=$HOME/go
export PATH=$PATH:$GOPATH/bin
```
You might want to put these in your `.profile` etc.

[4] clone this repository into `~/go/src/github.com/gkhasel1/cricket`

[5] install `godep` using `go get github.com/tools/godep` then get dependencies
```
godep get
```

[6] install the application
```
go install
```

[7] setup the db
```
docker-compose up -d elasticsearch
cricket --init
```

[8] run the application from binary
```
cricket
```

for convenience, this applicaiton has a request log
that streams to stdout and looks like

```
2017/04/10 05:36:32 POST	/v1/metrics	            PostMetrics	32.575µs
2017/04/10 05:36:54 POST	/v1/metrics	            PostMetrics	105.681µs
2017/04/10 05:37:12 GET	    /v1/metrics?name=c1		GetMetrics	15.085372ms
```

## Testing

You can run the test suite by executing `go test` from the root directory.
