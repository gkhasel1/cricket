package main

import (
	"log"
	"net/http"
	"context"
	"os"

	elastic "gopkg.in/olivere/elastic.v5"
)

/**
 * These could be passed in as Environment Variables at runtime
 * to improve configurability.
 */
const (
	ELASTICSEARCH_URL = "http://127.0.0.1:9200"
	ELASTICSEARCH_INDEX = "metrics"
	ELASTICSEARCH_TYPE = "metrics"

	PORT = ":8080"
)

func main() {
	/**
	 * MAIN: Provisions the DB with an index if `--init` is passed.
	 *       Runs the webserver on port specified above.
	 **/
    if (len(os.Args) > 1 && os.Args[1] == "--init") {
    	Init()
    }

	router := NewRouter()

	log.Fatal(http.ListenAndServe(PORT, router))
}

func Init() {
	/**
	 * INIT: Creates the `metrics` Elasticsearch index and exits.
	 *	     This should be run to make this web service is fully operational.
	 **/
	ctx := context.Background()

	client, err := elastic.NewSimpleClient(elastic.SetURL(ELASTICSEARCH_URL))
	if err != nil {
		panic(err)
	}

	exists, err := client.IndexExists(ELASTICSEARCH_INDEX).Do(context.Background())
	if err != nil {
		panic(err)
	}

	if !exists {
		_, err = client.CreateIndex(ELASTICSEARCH_INDEX).Do(ctx)
		if err != nil {
			panic(err)
		}
	}

	return
}
