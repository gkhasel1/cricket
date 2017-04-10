package main

import (
	"fmt"
	"net/http"
	"encoding/json"
)

func IndexHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Fprint(w, "DEADMOON\na metrics server written in go, backed by elasticseach\n")
}

func HealthHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Fprint(w, "OK")
}

func PostMetricsHandler(w http.ResponseWriter, r *http.Request) {
	var metrics *Metrics

	decoder := json.NewDecoder(r.Body)
	err := decoder.Decode(&metrics)
	if err != nil {
		http.Error(w, err.Error(), 500)
		return
	}

	defer r.Body.Close()

	for _, metric := range *metrics {
		err := metric.Save()
		if err != nil {
			http.Error(w, err.Error(), 500)
			return
		}
	}

	fmt.Fprint(w, "Success")
}

func GetMetricsHandler(w http.ResponseWriter, r *http.Request) {
	name := r.URL.Query().Get("name")
	time :=  r.URL.Query().Get("timestamp")

	metrics, err := GetMetrics(name, time)
	if err != nil {
		http.Error(w, err.Error(), 500)
		return
	}

	if len(metrics) < 1 {
		http.NotFound(w, r)
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
	name := r.URL.Query().Get("name")
	start :=  r.URL.Query().Get("start")
	end :=  r.URL.Query().Get("end")

	if name == "" || start == "" || end == "" {
		fmt.Fprint(w, "fu")
		http.Error(w, "Invalid Query: name, start, and end are required.", 500)
		return
	}

	stat, err := GetStats(name, start, end)
	if err != nil {
		http.Error(w, err.Error(), 500)
		return
	}

	if stat.Count < 1 {
		http.NotFound(w, r)
		return
	}

	jsonStat, err := json.Marshal(stat)
	if err != nil {
		http.Error(w, err.Error(), 500)
		return
	}

	fmt.Fprint(w, string(jsonStat))
}
