package server

import (
	"bufio"
	"context"
	"net/http"
	"os"
	"time"

	"github.com/dcwk/metrics/internal/config"
	"github.com/dcwk/metrics/internal/handlers"
	"github.com/dcwk/metrics/internal/logger"
	"github.com/dcwk/metrics/internal/models"
	"github.com/dcwk/metrics/internal/storage"
	"github.com/dcwk/metrics/internal/utils"
	"github.com/go-chi/chi/v5"
	"github.com/mailru/easyjson"
	"go.uber.org/zap"
)

func Run(ctx context.Context, conf *config.ServerConf) {
	if err := logger.Initialize(conf.LogLevel); err != nil {
		panic(err)
	}
	stor := storage.NewStorage()
	if conf.Restore {
		restore(stor, conf)
	}

	go func() {
		logger.Log.Info("Running server", zap.String("address", conf.ServerAddr))
		if err := http.ListenAndServe(conf.ServerAddr, Router(stor)); err != nil {
			panic(err)
		}
	}()

	flush(ctx, stor, conf)
}

func Router(storage storage.DataKeeper) chi.Router {
	r := chi.NewRouter()
	r.Use(logger.RequestLogger)
	r.Use(utils.GzipMiddleware)

	h := handlers.Handlers{
		Storage: storage,
	}

	r.Get("/", h.GetAllMetrics)
	r.Get("/value/{type}/{name}", h.GetMetricByParams)
	r.Post("/value/", h.GetMetricByJSON)
	r.Post("/update/{type}/{name}/{value}", h.UpdateMetricByParams)
	r.Post("/update/", h.UpdateMetricByJSON)

	return r
}

func restore(storage storage.DataKeeper, conf *config.ServerConf) {
	if conf.FileStoragePath == "" {
		return
	}

	logger.Log.Info("start restore data from file" + conf.FileStoragePath)
	file, err := os.OpenFile(conf.FileStoragePath, os.O_RDONLY|os.O_CREATE, 0666)
	if err != nil {
		panic(err)
	}
	scanner := bufio.NewScanner(file)
	if !scanner.Scan() {
		return
	}
	data := scanner.Bytes()

	logger.Log.Info("file data: " + string(data))

	metricsList := models.MetricsList{}
	if err := easyjson.Unmarshal(data, &metricsList); err != nil {
		panic(err)
	}

	storage.SaveMetricsList(&metricsList)
}

func flush(ctx context.Context, storage storage.DataKeeper, conf *config.ServerConf) {
	if conf.FileStoragePath == "" {
		return
	}
	duration := time.Duration(conf.StoreInterval) * time.Second
	if conf.StoreInterval == 0 {
		duration = time.Duration(100) * time.Millisecond
	}

	timer := time.NewTicker(duration)
	for {
		select {
		case <-ctx.Done():
			return
		case <-timer.C:
			logger.Log.Info("start flush data to file " + conf.FileStoragePath)
			metricsJSON, err := storage.GetJSONMetrics()
			if err != nil {
				return
			}

			file, err := os.OpenFile(conf.FileStoragePath, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, 0666)
			if err != nil {
				return
			}

			n, err := file.Write([]byte(metricsJSON))
			if err != nil {
				panic(err)
			}
			if n < len(metricsJSON) {
				return
			}

			if err := file.Close(); err != nil {
				panic(err)
			}
		}
	}
}
