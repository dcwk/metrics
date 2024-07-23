package server

import (
	"bufio"
	"database/sql"
	"net/http"
	"os"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	_ "github.com/jackc/pgx/v5/stdlib"
	"github.com/mailru/easyjson"
	"go.uber.org/zap"

	"github.com/dcwk/metrics/internal/config"
	"github.com/dcwk/metrics/internal/handlers"
	"github.com/dcwk/metrics/internal/logger"
	"github.com/dcwk/metrics/internal/models"
	"github.com/dcwk/metrics/internal/storage"
	"github.com/dcwk/metrics/internal/utils"
)

func Run(conf *config.ServerConf) {
	if err := logger.Initialize(conf.LogLevel); err != nil {
		panic(err)
	}
	var memStorage *storage.MemStorage
	var dbStorage *storage.DatabaseStorage

	if conf.DatabaseDSN == "" {
		memStorage = storage.NewStorage()
	} else {
		//sleepDuration := 1
		var err error
		db, err := sql.Open("pgx", conf.DatabaseDSN)
		if err != nil {
			panic(err)
		}
		defer db.Close()

		//for i := 0; i < 3; i++ {
		//	err = db.Ping()
		//	if err != nil {
		//		logger.Log.Info(fmt.Sprintf("Can't connect to database sleep %v", sleepDuration))
		//		time.Sleep(time.Duration(sleepDuration) * time.Second)
		//		sleepDuration += 2
		//		continue
		//	}
		//
		//	break
		//}
		//if err != nil {
		//	panic(err)
		//}

		dbStorage, err = storage.NewDBStorage(db)
		if err != nil {
			panic(err)
		}
	}

	if conf.Restore && memStorage != nil {
		restore(memStorage, conf)
	}

	if memStorage != nil {
		go flush(memStorage, conf)
	}

	logger.Log.Info("Running server", zap.String("address", conf.ServerAddr))
	if memStorage != nil {
		if err := http.ListenAndServe(conf.ServerAddr, Router(memStorage, conf)); err != nil {
			panic(err)
		}
	} else {
		if err := http.ListenAndServe(conf.ServerAddr, Router(dbStorage, conf)); err != nil {
			panic(err)
		}
	}
}

func Router(storage storage.DataKeeper, conf *config.ServerConf) chi.Router {
	r := chi.NewRouter()

	r.Use(logger.RequestLogger)
	if conf.IsActivePprof {
		r.Mount("/debug", middleware.Profiler())
	}

	h := handlers.Handlers{
		Storage: storage,
	}

	r.Route("/", func(r chi.Router) {
		r.Use(utils.GzipMiddleware)
		r.Use(utils.SignMiddleware(conf.HashKey))

		r.Get("/", h.GetAllMetrics)
		r.Get("/ping", h.Ping)
		r.Get("/value/{type}/{name}", h.GetMetricByParams)
		r.Post("/value/", h.GetMetricByJSON)
		r.With(utils.DecodeBodyMiddleware(conf.CryptoKey)).Post("/update/{type}/{name}/{value}", h.UpdateMetricByParams)
		r.With(utils.DecodeBodyMiddleware(conf.CryptoKey)).Post("/update/", h.UpdateMetricByJSON)
		r.Post("/updates/", h.UpdateBatchMetricByJSON)
	})

	return r
}

func restore(storage storage.MemoryKeeper, conf *config.ServerConf) {
	if conf.FileStoragePath == "" {
		return
	}

	logger.Log.Info("start restore data from file" + conf.FileStoragePath)
	file, err := os.OpenFile(conf.FileStoragePath, os.O_RDONLY|os.O_CREATE, 0666)
	if err != nil {
		return
	}
	scanner := bufio.NewScanner(file)
	if !scanner.Scan() {
		return
	}
	data := scanner.Bytes()
	metricsList := models.MetricsList{}
	if err := easyjson.Unmarshal(data, &metricsList); err != nil {
		panic(err)
	}

	storage.SaveMetricsList(&metricsList)
}

func flush(storage storage.MemoryKeeper, conf *config.ServerConf) {
	if conf.FileStoragePath == "" {
		return
	}

	for {
		logger.Log.Info("start flush data to file " + conf.FileStoragePath)
		metricsJSON, err := storage.GetJSONMetrics()
		if err != nil {
			return
		}

		file, err := os.OpenFile(conf.FileStoragePath, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, 0666)
		if err != nil {
			return
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
