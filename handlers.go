package main

import (
	"fmt"
	"net/http"
	"encoding/json"
)

func IndexHandler(w http.ResponseWriter, r *http.Request) {
	/**
	 * Index: Returns the name of the project, just for fun really.
	 **/
	fmt.Fprint(w, "CRICKET\na metrics server written in go, backed by elasticseach\n")
}

func HealthHandler(w http.ResponseWriter, r *http.Request) {
	/**
	 * Health: this method would be called by a loadbalancing mechanism
	 *		to ensure that the application is running correctly.
	 **/
	fmt.Fprint(w, "OK")
}

func PostMetricsHandler(w http.ResponseWriter, r *http.Request) {
	/**
	 * PostMetrics: This handler takes a metric json payload and
	 *		writes it to Elasticsearch in the metrics index. The user
	 *      posting the metric can send a list with multiple metrics
	 *      for convenient batching.
	 **/
	var metrics *Metrics

	decoder := json.NewDecoder(r.Body)
	err := decoder.Decode(&metrics)
	if err != nil {
		http.Error(w, err.Error(), 500)
		return
	}

	// Close the response body after this function is done
	defer r.Body.Close()

	for _, metric := range *metrics {
		/* We could potentially use a goroutine in the save call and
		 * improve performance. This would either require us to ignore errors (bad™)
		 * or use a coroutine communication mechanism (complicated™).
		 */
		err := metric.Save()
		if err != nil {
			http.Error(w, err.Error(), 500)
			return
		}
	}

	fmt.Fprint(w, "Success")
}

func GetMetricsHandler(w http.ResponseWriter, r *http.Request) {
	/**
	 * GetMetrics: Takes name and timestamp as optional params. If
	 *		neither is provided, we return all metrics ever seen (up to
	 *		10000 items just to be safe). If only name is provided, all metrics
	 * 		with that name are returned. If only timestamp is provided all metrics
	 * 		at that time are returned. If both, we return metrics by the specified
	 * 		name sent at a the given time.
	 **/
	name := r.URL.Query().Get("name")
	time :=  r.URL.Query().Get("timestamp")

	metrics, err := GetMetrics(name, time)
	if err != nil {
		http.Error(w, err.Error(), 500)
		return
	}

	if len(metrics) < 1 {
		http.NotFound(w, r) // 404 response code
		return
	}

	jsonMetrics, err := json.Marshal(metrics)
	if err != nil {
		http.Error(w, err.Error(), 500)
		return
	}

	fmt.Fprint(w, string(jsonMetrics))
}

func GetStatsHandler(w http.ResponseWriter, r *http.Request) {
	/**
	 * GetStats: Given a name, start, and end timestamp (all required)
	 *		this handler gets the Count, Sum, Min, Max, Average of the given
	 *		metrics within the time window specified. The time window is
	 *		inclusive.
	 **/
	name := r.URL.Query().Get("name")
	start :=  r.URL.Query().Get("start")
	end :=  r.URL.Query().Get("end")

	if name == "" || start == "" || end == "" {
		http.Error(w, "Invalid Query: name, start, and end are required.", 500)
		return
	}

	stat, err := GetStats(name, start, end)
	if err != nil {
		http.Error(w, err.Error(), 500)
		return
	}

	if stat.Count < 1 {
		http.NotFound(w, r) // 404 response code
		return
	}

	jsonStat, err := json.Marshal(stat)
	if err != nil {
		http.Error(w, err.Error(), 500)
		return
	}

	fmt.Fprint(w, string(jsonStat))
}
