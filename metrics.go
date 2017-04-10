package main

import (
	"time"
	"context"
	"encoding/json"

	elastic "gopkg.in/olivere/elastic.v5"
	uuid "github.com/google/uuid"
)

/**
 * STAT: the response object to a call to GetStats
 **/
type Stat struct {
	Count   int64		`json:"count"`
	Sum		float64		`json:"sum"`
	Min 	float64		`json:"min"`
	Max 	float64		`json:"max"`
	Average float64		`json:"average"`
}


/**
 * METRIC: the metrics object passed into PostMetrics and returned
 *		by GetMetrics. For convenience we also define the slice,
 *		Metrics []Metric below (this is basically just an array struct).
 **/
type Metric struct {
	Name        string		`json:"name"`
	Type        string		`json:"type"`
	Value       float64		`json:"value"`
	Timestamp	time.Time 	`json:"timestamp"`
}
type Metrics []Metric


func (metric *Metric) Save() error {
	/**
	 * Save: This is method on a Metric struct. Once a metric is constructed,
	 *		we save it to elasticsearch using the process below. The ID of a
	 *		metric is a uuidV4.
	 **/
	ctx := context.Background()

	// Get an ES client
	client, err := elastic.NewSimpleClient(elastic.SetURL(ELASTICSEARCH_URL))
	if err != nil {
		return err
	}

	// Write our metric to the metric index, refresh afterwards.
	_, err = client.Index().
		Index(ELASTICSEARCH_INDEX).
		Type(ELASTICSEARCH_TYPE).
		Id(uuid.New().String()).
		BodyJson(metric).
		Refresh("true").
		Do(ctx)
	if err != nil {
		return err
	}

	return nil
}

func GetMetrics(name string, timeStr string) (Metrics, error) {
	/**
	 * GetMetrics: Given an optional name and timestamp string, we get metrics
	 *		based on these criteria. Because we are using Elasticsearch v5
	 *		we must use Boolean Queries instead of more traditional filters.
	 **/
	metrics := make(Metrics, 0)

	ctx := context.Background()

	// Get an ES client
	client, err := elastic.NewSimpleClient(elastic.SetURL(ELASTICSEARCH_URL))
	if err != nil {
		return metrics , err
	}

	/* Start building our query. We will combine name and timestamp
	 * to compose the query. If neither is specified its the same as a `*` query.
	 */
	baseQuery := client.Search().Index(ELASTICSEARCH_INDEX)
	boolQuery := elastic.NewBoolQuery()

	// Add name query if we have it
	if name != "" {
		nameQuery := elastic.NewTermQuery("name", name)
		boolQuery = boolQuery.Must(nameQuery)
	}

	// Add time query if we have it
	if timeStr != "" {
		// This is just golang's strange timestring formatting mechanism
		timestamp, err := time.Parse("2006-01-02T15:04:05Z", timeStr)
		if err != nil {
			return metrics , err
		}
		timestampQuery := elastic.NewTermQuery("timestamp", timestamp)
		boolQuery = boolQuery.Must(timestampQuery)
	}

	// Run our completed query
	searchResult, err := baseQuery.
		Query(boolQuery).
		Sort("timestamp", true).
		Pretty(true).
		From(0).Size(9999). // For safety, should use paginatation instead
		Do(ctx)
	if err != nil {
		return metrics , err
	}

	// Convert our hits from json to structs and return them to the handler
	if searchResult.Hits.TotalHits > 0 {
		for _, hit := range searchResult.Hits.Hits {
			var metric Metric
			err := json.Unmarshal(*hit.Source, &metric)
			if err != nil {
				return metrics , err
			}
			metrics = append(metrics, metric)
		}
	}

	return metrics, nil
}

func GetStats(name string, start string, end string) (Stat, error) {
	/**
	 * GetStats: Given a required name and start/end timestamp strings, we get
	 *		statistics on a metrics set bounded by those paramenters. Because we
	 *		are using Elasticsearch v5 we must use Boolean Queries instead of
	 *		more traditional filters. We don't need to validate the existance
	 * 		of these params here, it is done in the handler that calls this method.
	 **/
	var stat Stat

	ctx := context.Background()

	// Get an ES client
	client, err := elastic.NewSimpleClient(elastic.SetURL(ELASTICSEARCH_URL))
	if err != nil {
		return stat , err
	}

	// Parse our start time string into a time struct
	startTime, err := time.Parse("2006-01-02T15:04:05Z", start)
	if err != nil {
		return stat , err
	}

	// Parse our end time string into a time struct
	endTime, err := time.Parse("2006-01-02T15:04:05Z", end)
	if err != nil {
		return stat , err
	}

	// Create our name query, our timerange query, and combine them into a bool query
	nameQuery := elastic.NewTermQuery("name", name)
	rangeQuery := elastic.NewRangeQuery("timestamp").Gte(startTime).Lte(endTime)
	boolQuery := elastic.NewBoolQuery().Must(nameQuery).Must(rangeQuery)

	// Run our completed query
	searchResult, err :=  client.Search().
		Index(ELASTICSEARCH_INDEX).
		Query(boolQuery).
		Sort("timestamp", true).
		Pretty(true).
		From(0).Size(9999). // For safety, should use paginatation instead
		Do(ctx)
	if err != nil {
		return stat , err
	}

	sum := 0.0

	// Convert our hits from json to structs and calculate the desired stats.
	// All stats are float64 for better precision.
	if searchResult.Hits.TotalHits > 0 {
		for _, hit := range searchResult.Hits.Hits {
			var metric Metric
			err := json.Unmarshal(*hit.Source, &metric)
			if err != nil {
				return stat , err
			}

			if metric.Value > stat.Max {
				stat.Max = metric.Value
			}
			if metric.Value < stat.Min {
				stat.Min = metric.Value
			}

			sum = sum + metric.Value
		}

		stat.Sum = sum
		stat.Count = searchResult.Hits.TotalHits
		stat.Average = sum / float64(searchResult.Hits.TotalHits)
	}

	return stat, nil
}
