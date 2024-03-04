package main

import (
	"metrics/internal/storage"
	"metrics/internal/util"
	"net/http"
)

func main() {
	mux := http.NewServeMux()
	mux.HandleFunc("/update/gauge", GaugeHandler)
	mux.HandleFunc("/update/counter", CounterHandler)

	err := http.ListenAndServe("localhost:8080", mux)
	if err != nil {
		panic(err)
	}
}

func GaugeHandler(w http.ResponseWriter, r *http.Request) {
	r.Method = http.MethodPost
	r.Header.Set("Content-Type", "text/plain")

	mn, mv, err := util.ParamsFromURL(r.URL.Path)
	if err != nil {
		w.WriteHeader(http.StatusNotFound)

		return
	}

	stor := storage.NewStorage()
	stor.AddMetric(mn, mv)

	w.WriteHeader(http.StatusOK)
}

func CounterHandler(w http.ResponseWriter, r *http.Request) {
	r.Method = http.MethodPost
	r.Header.Set("Content-Type", "text/plain")

	mn, mv, err := util.ParamsFromURL(r.URL.Path)
	if err != nil {
		w.WriteHeader(http.StatusNotFound)

		return
	}

	stor := storage.NewStorage()
	stor.AddMetric(mn, mv)

	w.WriteHeader(http.StatusOK)
}
