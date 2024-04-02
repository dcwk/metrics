package handlers

import (
	"fmt"

	"github.com/dcwk/metrics/internal/logger"
	"github.com/dcwk/metrics/internal/models"
	"github.com/go-resty/resty/v2"
	"github.com/mailru/easyjson"
)

func (h *Handlers) SendMetrics(metrics map[string]float64, addr string, pollCount *int64) error {
	for k, v := range metrics {
		metric := models.Metrics{
			ID:    k,
			MType: models.Gauge,
			Value: &v,
		}
		json, err := easyjson.Marshal(&metric)
		if err != nil {
			return err
		}

		logger.Log.Info(string(json))

		if err := send(string(json), addr); err != nil {
			return err
		}
	}

	metric := models.Metrics{
		ID:    "PollCount",
		MType: models.Counter,
		Delta: pollCount,
	}

	json, err := easyjson.Marshal(&metric)
	if err != nil {
		return err
	}

	if err := send(string(json), addr); err != nil {
		return err
	}

	return nil
}

func send(metricsJSON string, addr string) error {
	client := resty.New()
	_, err := client.R().
		SetHeader("Content-Type", "application/json").
		SetHeader("Accept-Encoding", "gzip").
		SetHeader("Content-Encoding", "gzip").
		SetBody(metricsJSON).
		Post(fmt.Sprintf("http://%s/update/", addr))
	if err != nil {
		return err
	}

	return nil
}
