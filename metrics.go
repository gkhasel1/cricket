package main

import (
	"time"
	"context"
	"encoding/json"

	elastic "gopkg.in/olivere/elastic.v5"
	uuid "github.com/google/uuid"
)

type Stat struct {
	Count   int64		`json:"count"`
	Sum		float64		`json:"sum"`
	Min 	float64		`json:"min"`
	Max 	float64		`json:"max"`
	Average float64		`json:"average"`
}

type Metric struct {
	Name        string		`json:"name"`
	Type        string		`json:"type"`
	Value       float64		`json:"value"`
	Timestamp	time.Time 	`json:"timestamp"`
}

type Metrics []Metric


func (metric *Metric) Save() error {
	ctx := context.Background()

	client, err := elastic.NewSimpleClient(elastic.SetURL(ELASTICSEARCH_URL))
	if err != nil {
		return err
	}

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
	metrics := make(Metrics, 0)

	ctx := context.Background()

	client, err := elastic.NewSimpleClient(elastic.SetURL(ELASTICSEARCH_URL))
	if err != nil {
		return metrics , err
	}

	baseQuery := client.Search().Index(ELASTICSEARCH_INDEX)
	boolQuery := elastic.NewBoolQuery()

	if name != "" {
		nameQuery := elastic.NewTermQuery("name", name)
		boolQuery = boolQuery.Must(nameQuery)
	}

	if timeStr != "" {
		timestamp, err := time.Parse("2006-01-02T15:04:05Z", timeStr)
		if err != nil {
			return metrics , err
		}
		timestampQuery := elastic.NewTermQuery("timestamp", timestamp)
		boolQuery = boolQuery.Must(timestampQuery)
	}

	searchResult, err := baseQuery.
		Query(boolQuery).
		Sort("timestamp", true).
		Pretty(true).
		From(0).Size(9999).
		Do(ctx)
	if err != nil {
		return metrics , err
	}

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
	var stat Stat

	ctx := context.Background()

	client, err := elastic.NewSimpleClient(elastic.SetURL(ELASTICSEARCH_URL))
	if err != nil {
		return stat , err
	}

	startTime, err := time.Parse("2006-01-02T15:04:05Z", start)
	if err != nil {
		return stat , err
	}

	endTime, err := time.Parse("2006-01-02T15:04:05Z", end)
	if err != nil {
		return stat , err
	}

	nameQuery := elastic.NewTermQuery("name", name)
	rangeQuery := elastic.NewRangeQuery("timestamp").Gte(startTime).Lte(endTime)
	boolQuery := elastic.NewBoolQuery().Must(nameQuery).Must(rangeQuery)

	searchResult, err :=  client.Search().
		Index(ELASTICSEARCH_INDEX).
		Query(boolQuery).
		Sort("timestamp", true).
		Pretty(true).
		From(0).Size(9999).
		Do(ctx)
	if err != nil {
		return stat , err
	}

	sum := 0.0
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
