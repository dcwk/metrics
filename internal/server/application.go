package server

import (
	"fmt"
	"net/http"

	"github.com/dcwk/metrics/internal/config"
	"github.com/dcwk/metrics/internal/handlers"
	"github.com/dcwk/metrics/internal/storage"
	"github.com/go-chi/chi/v5"
)

func Run(conf *config.ServerConf) {
	s := storage.NewStorage()
	fmt.Println("Running server on", conf.ServerAddr)

	if err := http.ListenAndServe(conf.ServerAddr, Router(s)); err != nil {
		panic(err)
	}
}

func Router(s storage.Storage) chi.Router {
	r := chi.NewRouter()

	r.Get("/", func(w http.ResponseWriter, r *http.Request) {
		handlers.GetAllMetrics(w, r, s)
	})

	r.Route("/value/", func(r chi.Router) {
		r.Get("/{type}/{name}", func(w http.ResponseWriter, r *http.Request) {
			handlers.GetMetric(w, r, s)
		})
	})

	r.Route("/update/", func(r chi.Router) {
		r.Post("/{type}/{name}/{value}", func(w http.ResponseWriter, r *http.Request) {
			handlers.UpdateMetric(w, r, s)
		})
	})

	return r
}
