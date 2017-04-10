API DOCUMENTATION:

Global Routes:

`GET /` -- index page, shows the title of the project
`GET /health` -- returns 200: OK when the app is running (shallow healthcheck)

Version 1 Routes:

`POST /v1/metrics` -- accepts json data as shown below and creates metric datapoints
EXAMPLE
```
$ curl -X POST -H 'Content-Type: application/json' -d
'[
	{
		"name": "c1",
		"type": "counter",
		"value": 1,
		"timestamp": "2016-12-31T21:49:38Z" (RFC3339)
	},
	{
		"name": "c2",
		"type": "counter",
		"value": 4.7,
		"timestamp": "2016-12-31T22:49:50Z" (RFC3339)
	}
]'
http://localhost:8080/v1/metrics
```

`GET /v1/metrics` -- returns metric data, accepts metric `name` and `timestamp` as optional query params
EXAMPLE
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

`GET /v1/stats` -- returns a set of statistics about a metric, accepts metric `name`, `start`, and `end` as required query params
EXAMPLE
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


DEVELOPMENT SETUP:

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

[4] clone this repository into `~/go/src/github.com/gkhasel1/deadmoon`

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
deadmoon --init
```

[5] run the application from binary
```
deadmoon
```
