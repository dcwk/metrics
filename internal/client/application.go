package client

import (
	"runtime"
	"time"

	"github.com/dcwk/metrics/internal/config"
	"github.com/dcwk/metrics/internal/handlers"
	"github.com/dcwk/metrics/internal/logger"
	"github.com/dcwk/metrics/internal/storage"
)

func Run(conf *config.ClientConf) error {
	if err := logger.Initialize(conf.LogLevel); err != nil {
		panic(err)
	}
	var pollCount int64
	h := handlers.Handlers{
		Storage: storage.NewStorage(),
	}
	logger.Log.Info("Sending metrics to" + conf.ServerAddr)

	for {
		time.Sleep(time.Duration(conf.PollInterval) * time.Second)
		if pollCount%(conf.ReportInterval/conf.PollInterval) != 0 {
			pollCount++
			ms := runtime.MemStats{}
			runtime.ReadMemStats(&ms)
			continue
		}

		if err := h.SendMetrics(conf.ServerAddr, pollCount); err != nil {
			return err
		}

		pollCount++
	}
}
