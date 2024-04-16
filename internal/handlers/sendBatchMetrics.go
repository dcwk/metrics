package handlers

import (
	"encoding/json"
	"fmt"
	"log"

	"github.com/dcwk/metrics/internal/models"
)

func (h *Handlers) SendBatchMetrics(metrics map[string]float64, addr string, pollCount *int64) error {
	path := fmt.Sprintf("http://%s/updates/", addr)
	metricsList := models.MetricsList{}

	for k, v := range metrics {
		metric := models.Metrics{
			ID:    k,
			MType: models.Gauge,
			Value: &v,
		}

		metricsList.List = append(metricsList.List, metric)
	}

	metric := models.Metrics{
		ID:    "PollCount",
		MType: models.Counter,
		Delta: pollCount,
	}

	metricsList.List = append(metricsList.List, metric)

	jsonData, err := json.Marshal(metricsList.List)
	if err != nil {
		return err
	}

	if err := send(jsonData, path); err != nil {
		return err
	}
	log.Printf("reported metrics in JSON %s\n", string(jsonData))

	return nil
}
