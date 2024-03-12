package client

import (
	"fmt"
	"github.com/dcwk/metrics/internal/config"
	"github.com/dcwk/metrics/internal/handlers"
	"time"
)

func Run(conf *config.ClientConf) error {
	var pollCount int64
	fmt.Println("Sending metrics to", conf.ServerAddr)

	for {
		time.Sleep(time.Duration(conf.PollInterval) * time.Second)
		if pollCount%conf.ReportInterval != 0 {
			continue
		}

		if err := handlers.SendMetrics(conf.ServerAddr); err != nil {
			return err
		}

		pollCount++
	}
}
