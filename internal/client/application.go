package client

import (
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
	if err := updateMetrics(conf, &pollCount); err != nil {
		return err
	}

	return nil
}

func updateMetrics(conf *config.ClientConf, pollCount *int64) error {
	h := handlers.Handlers{
		Storage:    storage.NewStorage(),
		ClientConf: conf,
	}

	for {
		time.Sleep(time.Duration(conf.ReportInterval) * time.Second)
		metrics := getGauges(pollCount)
		if err := h.SendMetrics(metrics, conf.ServerAddr, pollCount); err != nil {

			return err
		}
	}
}
