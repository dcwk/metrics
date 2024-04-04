package client

import (
	"context"
	"log"
	"sync"
	"time"

	"github.com/dcwk/metrics/internal/config"
	"github.com/dcwk/metrics/internal/handlers"
	"github.com/dcwk/metrics/internal/logger"
	"github.com/dcwk/metrics/internal/storage"
)

func Run(ctx context.Context, conf *config.ClientConf) error {
	if err := logger.Initialize(conf.LogLevel); err != nil {
		return err
	}

	log.Printf("Sending metrics to %s\n", conf.ServerAddr)
	wg := new(sync.WaitGroup)
	wg.Add(2)

	agent := NewAgent(conf.PollInterval, conf.ReportInterval)
	go agent.pollMetrics(ctx, wg)
	go reportMetrics(ctx, wg, conf, agent)

	<-ctx.Done()

	wg.Wait()

	return nil
}

func reportMetrics(ctx context.Context, wg *sync.WaitGroup, conf *config.ClientConf, agent *Agent) {
	reportTicker := time.NewTicker(time.Duration(conf.ReportInterval) * time.Second)
	defer reportTicker.Stop()

	h := handlers.Handlers{
		Storage: storage.NewStorage(),
	}

	for {
		select {
		case <-ctx.Done():
			wg.Done()
			return
		case <-reportTicker.C:
			time.Sleep(time.Duration(conf.ReportInterval) * time.Second)
			if err := h.SendMetrics(agent.Metrics, conf.ServerAddr, &agent.PollCount); err != nil {
			}
		}
	}
}
