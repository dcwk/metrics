package main

import (
	"flag"
	"fmt"
	"github.com/go-chi/chi/v5"
	"metrics/internal/storage"
	"net/http"
	"os"
)

const (
	gauge   = "gauge"
	counter = "counter"
)

var flagRunAddr string

func main() {
	parseFlags()
	fmt.Println("Running server on", flagRunAddr)

	if err := http.ListenAndServe(flagRunAddr, Router()); err != nil {
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

func parseFlags() {
	flag.StringVar(&flagRunAddr, "a", ":8080", "address and port to run server")
	flag.Parse()

	if envAddress := os.Getenv("ADDRESS"); envAddress != "" {
		flagRunAddr = envAddress
	}
}

func getAllMetricsHandler(w http.ResponseWriter, r *http.Request) {
	r.Method = http.MethodGet
	r.Header.Set("Content-Type", "text/plain")

	stor := storage.NewStorage()
	gauges := stor.GetAllGauges()
	counters := stor.GetAllCounters()
	res := ""

	for n, v := range gauges {
		res += n + " " + fmt.Sprintf("%.3f", v) + "\n\r"
	}

	for n, v := range counters {
		res += n + " " + fmt.Sprintf("%d", v) + "\n\r"
	}

	w.WriteHeader(http.StatusOK)
	w.Write([]byte(res))
}

func getMetricHandler(w http.ResponseWriter, r *http.Request) {
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
		v = fmt.Sprintf("%.3f", metricValue)
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

func updateMetricHandler(w http.ResponseWriter, r *http.Request) {
	r.Method = http.MethodPost
	r.Header.Set("Content-Type", "text/plain")

	t := chi.URLParam(r, "type")
	mn := chi.URLParam(r, "name")
	mv := chi.URLParam(r, "value")
	stor := storage.NewStorage()

	switch t {
	default:
		http.Error(w, "", http.StatusBadRequest)
		return
	case gauge:
		if err := stor.AddGauge(mn, mv); err != nil {
			http.Error(w, "", http.StatusBadRequest)
			return
		}
	case counter:
		if err := stor.AddCounter(mn, mv); err != nil {
			http.Error(w, "", http.StatusBadRequest)
			return
		}
	}

	w.WriteHeader(http.StatusOK)
}
