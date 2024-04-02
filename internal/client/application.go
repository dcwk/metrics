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
		return err
	}
	var pollCount int64

	logger.Log.Info("Sending metrics to" + conf.ServerAddr)

	go updateMemStat(conf.PollInterval, &pollCount)
	updateMetrics(conf.ServerAddr, conf.ReportInterval, &pollCount)

	return nil
}

func updateMemStat(pollInterval int64, pollCount *int64) {
	for {
		time.Sleep(time.Duration(pollInterval) * time.Second)
		ms := runtime.MemStats{}
		runtime.ReadMemStats(&ms)
		*pollCount++
	}
}

func updateMetrics(serverAddr string, reportInterval int64, pollCount *int64) error {
	h := handlers.Handlers{
		Storage: storage.NewStorage(),
	}

	for {
		time.Sleep(time.Duration(reportInterval) * time.Second)

		if err := h.SendMetrics(serverAddr, pollCount); err != nil {

			return err
		}
	}
}
