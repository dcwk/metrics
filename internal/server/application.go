package server

import (
	"net/http"
	"os"
	"time"

	"github.com/dcwk/metrics/internal/config"
	"github.com/dcwk/metrics/internal/handlers"
	"github.com/dcwk/metrics/internal/logger"
	"github.com/dcwk/metrics/internal/storage"
	"github.com/dcwk/metrics/internal/utils"
	"github.com/go-chi/chi/v5"
	"go.uber.org/zap"
)

func Run(conf *config.ServerConf) {
	if err := logger.Initialize(conf.LogLevel); err != nil {
		panic(err)
	}
	s := storage.NewStorage()
	logger.Log.Info("Running server", zap.String("address", conf.ServerAddr))
	go Flush(s, conf)

	if err := http.ListenAndServe(conf.ServerAddr, Router(s)); err != nil {
		panic(err)
	}
}

func Flush(s storage.DataKeeper, conf *config.ServerConf) {
	for {
		logger.Log.Info("start flush data")
		metricsJSON, err := s.GetJsonMetrics()
		if err != nil {
			panic(err)
		}

		file, err := os.OpenFile(conf.FileStoragePath, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, 0666)
		if err != nil {
			panic(err)
		}
		if _, err := file.Write([]byte(metricsJSON)); err != nil {
			panic(err)
		}
		if err := file.Close(); err != nil {
			panic(err)
		}

		time.Sleep(time.Duration(conf.StoreInterval) * time.Second)
	}
}

func Router(s storage.DataKeeper) chi.Router {
	r := chi.NewRouter()
	r.Use(logger.RequestLogger)
	r.Use(utils.GzipMiddleware)

	h := handlers.Handlers{
		Storage: s,
	}

	r.Get("/", h.GetAllMetrics)
	r.Get("/value/{type}/{name}", h.GetMetricByParams)
	r.Post("/value/", h.GetMetricByJSON)
	r.Post("/update/{type}/{name}/{value}", h.UpdateMetricByParams)
	r.Post("/update/", h.UpdateMetricByJSON)

	return r
}
