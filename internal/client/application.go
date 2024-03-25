package client

import (
	"fmt"
	"time"

	"github.com/dcwk/metrics/internal/config"
	"github.com/dcwk/metrics/internal/handlers"
	"github.com/dcwk/metrics/internal/storage"
)

func Run(conf *config.ClientConf) error {
	var pollCount int64
	h := handlers.Handlers{
		Storage: storage.NewStorage(),
	}
	fmt.Println("Sending metrics to", conf.ServerAddr)

	for {
		time.Sleep(time.Duration(conf.PollInterval) * time.Second)
		if pollCount%conf.ReportInterval != 0 {
			continue
		}

		if err := h.SendMetrics(conf.ServerAddr); err != nil {
			return err
		}

		pollCount++
	}
}
