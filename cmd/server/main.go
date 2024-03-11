package main

import (
	"fmt"
	"github.com/dcwk/metrics/internal/config"
	"github.com/dcwk/metrics/internal/storage"
	"github.com/go-chi/chi/v5"
	"net/http"
)

const (
	gauge   = "gauge"
	counter = "counter"
)

func main() {
	conf := config.NewServerConf()
	fmt.Println("Running server on", conf.ServerAddr)

	if err := http.ListenAndServe(conf.ServerAddr, Router()); err != nil {
		panic(err)
	}
}

func Router() chi.Router {
	r := chi.NewRouter()

	r.Get("/", getAllMetricsHandler)

	r.Route("/value/", func(r chi.Router) {
		r.Get("/{type}/{name}", getMetricHandler)
	})

	r.Route("/update/", func(r chi.Router) {
		r.Post("/{type}/{name}/{value}", updateMetricHandler)
	})

	return r
}

func getAllMetricsHandler(w http.ResponseWriter, r *http.Request) {
	r.Method = http.MethodGet
	r.Header.Set("Content-Type", "text/plain")

	s := storage.GetStorage()
	gauges := s.GetAllGauges()
	counters := s.GetAllCounters()
	res := ""

	for n, v := range gauges {
		res += n + " " + fmt.Sprintf("%.3f", v) + "\n\r"
	}

	for n, v := range counters {
		res += n + " " + fmt.Sprintf("%d", v) + "\n\r"
	}

	w.WriteHeader(http.StatusOK)
	if _, err := w.Write([]byte(res)); err != nil {
		http.Error(w, "", http.StatusInternalServerError)
		return
	}
}

func getMetricHandler(w http.ResponseWriter, r *http.Request) {
	r.Method = http.MethodGet
	r.Header.Set("Content-Type", "text/plain")

	t := chi.URLParam(r, "type")
	n := chi.URLParam(r, "name")
	v := ""
	s := storage.GetStorage()

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
		v = fmt.Sprintf("%v", metricValue)
	case counter:
		metricValue, err := s.GetCounter(n)
		if err != nil {
			http.Error(w, "", http.StatusNotFound)
			return
		}
		v = fmt.Sprintf("%d", metricValue)
	}

	w.WriteHeader(http.StatusOK)
	if _, err := w.Write([]byte(v)); err != nil {
		http.Error(w, "", http.StatusInternalServerError)
		return
	}
}

func updateMetricHandler(w http.ResponseWriter, r *http.Request) {
	r.Method = http.MethodPost
	r.Header.Set("Content-Type", "text/plain")

	t := chi.URLParam(r, "type")
	mn := chi.URLParam(r, "name")
	mv := chi.URLParam(r, "value")
	s := storage.GetStorage()

	switch t {
	default:
		http.Error(w, "", http.StatusBadRequest)
		return
	case gauge:
		if err := s.AddGauge(mn, mv); err != nil {
			http.Error(w, "", http.StatusBadRequest)
			return
		}
	case counter:
		if err := s.AddCounter(mn, mv); err != nil {
			http.Error(w, "", http.StatusBadRequest)
			return
		}
	}

	w.WriteHeader(http.StatusOK)
}
