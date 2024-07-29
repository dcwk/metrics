package client

import (
	"context"
	"fmt"
	"log"
	"os"
	"os/signal"
	"sync"
	"syscall"
	"time"

	"github.com/dcwk/metrics/internal/config"
	"github.com/dcwk/metrics/internal/logger"
)

func Run(ctx context.Context, conf *config.ClientConf) error {
	if err := logger.Initialize(conf.LogLevel); err != nil {
		return err
	}
	log.Printf("Sending metrics to %s\n", conf.ServerAddr)
	shutdouwnSignal := make(chan os.Signal, 1)
	signal.Notify(shutdouwnSignal, syscall.SIGINT, syscall.SIGTERM, syscall.SIGQUIT)
	ctx, cancel := context.WithCancel(ctx)

	wg := new(sync.WaitGroup)
	wg.Add(2)
	agent := NewAgent(conf.PollInterval, conf.ReportInterval)

	go func() {
		signal := <-shutdouwnSignal
		fmt.Printf("fired signal %v\n", signal)
		cancel()
	}()

	go agent.pollMetrics(ctx, wg)
	go reportMetrics(ctx, wg, conf, agent)

	wg.Wait()

	return nil
}

func reportMetrics(ctx context.Context, wg *sync.WaitGroup, conf *config.ClientConf, agent *Agent) {
	reportTicker := time.NewTicker(time.Duration(conf.ReportInterval) * time.Second)
	defer reportTicker.Stop()

	for {
		select {
		case <-ctx.Done():
			fmt.Printf("stopped reporting metrics\n")
			wg.Done()
			return
		case <-reportTicker.C:
			_ = SendBatchMetrics(agent.Metrics, conf.ServerAddr, conf.HashKey, conf.CryptoKey, &agent.PollCount)
		}
	}
}
