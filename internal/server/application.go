package server

import (
	"net/http"

	"github.com/dcwk/metrics/internal/config"
	"github.com/dcwk/metrics/internal/handlers"
	"github.com/dcwk/metrics/internal/logger"
	"github.com/dcwk/metrics/internal/storage"
	"github.com/go-chi/chi/v5"
	"go.uber.org/zap"
)

func Run(conf *config.ServerConf) {
	if err := logger.Initialize(conf.LogLevel); err != nil {
		panic(err)
	}
	s := storage.NewStorage()
	logger.Log.Info("Running server", zap.String("address", conf.ServerAddr))

	if err := http.ListenAndServe(conf.ServerAddr, Router(s)); err != nil {
		panic(err)
	}
}

func Router(s storage.DataKeeper) chi.Router {
	r := chi.NewRouter()
	h := handlers.Handlers{
		Storage: s,
	}

	r.Get("/", h.GetAllMetrics)
	r.Get("/value/{type}/{name}", h.GetMetric)
	r.Post("/update/{type}/{name}/{value}", h.UpdateMetric)

	return r
}
