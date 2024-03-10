package main

import (
	"github.com/go-chi/chi/v5"
	"metrics/internal/storage"
	"net/http"
)

func Router() chi.Router {
	r := chi.NewRouter()
	r.Route("/update/", func(r chi.Router) {
		r.Post("/update/gauge/{name}/{value}", GaugeHandler)
		r.Post("/update/counter/{name}/{value}", CounterHandler)
		r.Post("/update/", UnknownHandler)
	})

	return r
}

func main() {
	err := http.ListenAndServe("localhost:8080", Router())
	if err != nil {
		panic(err)
	}
}

func GaugeHandler(w http.ResponseWriter, r *http.Request) {
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

func CounterHandler(w http.ResponseWriter, r *http.Request) {
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

func UnknownHandler(w http.ResponseWriter, r *http.Request) {
	r.Method = http.MethodPost
	r.Header.Set("Content-Type", "text/plain")

	http.Error(w, "", http.StatusBadRequest)
}
