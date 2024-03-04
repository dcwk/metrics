package main

import (
	"metrics/internal/storage"
	"metrics/internal/util"
	"net/http"
)

func main() {
	mux := http.NewServeMux()
	mux.HandleFunc("/update/gauge/", GaugeHandler)
	mux.HandleFunc("/update/counter/", CounterHandler)
	mux.HandleFunc("/update/unknown/", UnknownHandler)

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
	err = stor.AddGauge(mn, mv)
	if err != nil {
		http.Error(w, "", http.StatusInternalServerError)
		return
	}

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
	err = stor.AddCounter(mn, mv)
	if err != nil {
		http.Error(w, "", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
}

func UnknownHandler(w http.ResponseWriter, r *http.Request) {
	r.Method = http.MethodPost
	r.Header.Set("Content-Type", "text/plain")

	w.WriteHeader(http.StatusOK)
}
