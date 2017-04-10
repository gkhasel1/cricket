package main

import (
	"log"
	"net/http"
	"context"
	"os"

	elastic "gopkg.in/olivere/elastic.v5"
)

const (
	ELASTICSEARCH_URL = "http://127.0.0.1:9200"
	ELASTICSEARCH_INDEX = "metrics"
	ELASTICSEARCH_TYPE = "metrics"

	PORT = ":8080"
)

func main() {
    if (len(os.Args) > 1 && os.Args[1] == "--init") {
    	Init()
    }

	router := NewRouter()

	log.Fatal(http.ListenAndServe(PORT, router))
}

func Init() {
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
}
