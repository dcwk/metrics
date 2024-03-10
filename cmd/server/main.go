package main

import (
	"fmt"
	"github.com/go-chi/chi/v5"
	"metrics/internal/storage"
	"net/http"
)

const (
	gauge   = "gauge"
	counter = "counter"
)

func Router() chi.Router {
	r := chi.NewRouter()
	r.Route("/update/", func(r chi.Router) {
		r.Post("/gauge/{name}/{value}", updateGaugeHandler)
		r.Post("/counter/{name}/{value}", updateCounterHandler)
		r.Post("/unknown/*", updateUnknownHandler)
	})
	r.Route("/value/", func(r chi.Router) {
		r.Get("/{type}/{name}", getDataHandler)
	})

	return r
}

func main() {
	err := http.ListenAndServe("localhost:8080", Router())
	if err != nil {
		panic(err)
	}
}

func updateGaugeHandler(w http.ResponseWriter, r *http.Request) {
	r.Method = http.MethodPost
	r.Header.Set("Content-Type", "text/plain")

	mn := chi.URLParam(r, "name")
	mv := chi.URLParam(r, "value")

	stor := storage.NewStorage()
	err := stor.AddGauge(mn, mv)
	if err != nil {
		http.Error(w, "", http.StatusBadRequest)
		return
	}

	w.WriteHeader(http.StatusOK)
}

func updateCounterHandler(w http.ResponseWriter, r *http.Request) {
	r.Method = http.MethodPost
	r.Header.Set("Content-Type", "text/plain")

	mn := chi.URLParam(r, "name")
	mv := chi.URLParam(r, "value")

	stor := storage.NewStorage()
	err := stor.AddCounter(mn, mv)
	if err != nil {
		http.Error(w, "", http.StatusBadRequest)
		return
	}

	w.WriteHeader(http.StatusOK)
}

func updateUnknownHandler(w http.ResponseWriter, r *http.Request) {
	r.Method = http.MethodPost
	r.Header.Set("Content-Type", "text/plain")

	http.Error(w, "", http.StatusBadRequest)
}

func getDataHandler(w http.ResponseWriter, r *http.Request) {
	r.Method = http.MethodGet
	r.Header.Set("Content-Type", "text/plain")

	t := chi.URLParam(r, "type")
	n := chi.URLParam(r, "name")
	v := ""
	s := storage.NewStorage()

	switch t {
	default:
		http.Error(w, "", http.StatusBadRequest)
		return
	case gauge:
		metricValue, err := s.GetGauge(n)
		if err != nil {
			http.Error(w, "", http.StatusNotFound)
			return
		}
		v = fmt.Sprintf("%f", metricValue)
	case counter:
		metricValue, err := s.GetCounter(n)
		if err != nil {
			http.Error(w, "", http.StatusNotFound)
			return
		}
		v = fmt.Sprintf("%d", metricValue)
	}

	w.WriteHeader(http.StatusOK)
	w.Write([]byte(v))
}
