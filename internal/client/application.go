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

	logger.Log.Info("Sending metrics to" + conf.ServerAddr)
	agent := NewAgent(conf.PollInterval, conf.ReportInterval)
	for {
		go agent.Update()
		if err := reportMetrics(conf, agent); err != nil {
			return err
		}
	}
}

func reportMetrics(conf *config.ClientConf, agent *Agent) error {
	h := handlers.Handlers{
		Storage: storage.NewStorage(),
	}

	for {
		time.Sleep(time.Duration(conf.ReportInterval) * time.Second)
		if err := h.SendMetrics(agent.Metrics, conf.ServerAddr, &agent.PollCount); err != nil {

			return err
		}
	}
}
