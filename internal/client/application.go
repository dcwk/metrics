package client

import (
	"context"
	"log"
	"sync"
	"time"

	"github.com/dcwk/metrics/internal/config"
	"github.com/dcwk/metrics/internal/logger"
)

func Run(ctx context.Context, conf *config.ClientConf) error {
	if err := logger.Initialize("info"); err != nil {
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
	workerPool := NewWorkerPool(conf.RateLimit)

	for {
		select {
		case <-ctx.Done():
			wg.Done()
			return
		case <-reportTicker.C:
			_ = SendMetricsInPool(agent.Metrics, conf.ServerAddr, conf.HashKey, workerPool, &agent.PollCount)
		}
	}
}
