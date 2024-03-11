package server

import (
	"fmt"
	"github.com/dcwk/metrics/internal/config"
	"github.com/dcwk/metrics/internal/handlers"
	"github.com/go-chi/chi/v5"
	"net/http"
)

func Run(conf config.ServerConf) {

	fmt.Println("Running server on", conf.ServerAddr)
	if err := http.ListenAndServe(conf.ServerAddr, Router()); err != nil {
		panic(err)
	}
}

func Router() chi.Router {
	r := chi.NewRouter()

	r.Get("/", handlers.GetAllMetrics)

	r.Route("/value/", func(r chi.Router) {
		r.Get("/{type}/{name}", handlers.GetMetric)
	})

	r.Route("/update/", func(r chi.Router) {
		r.Post("/{type}/{name}/{value}", handlers.UpdateMetric)
	})

	return r
}
