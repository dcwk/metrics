package client

import (
	"encoding/json"
	"fmt"
	"log"

	"github.com/dcwk/metrics/internal/logger"
	"github.com/dcwk/metrics/internal/models"
)

func SendBatchMetrics(
	metrics map[string]float64,
	addr string,
	grpcAddr string,
	hashKey string,
	cryptoKey string,
	pollCount *int64,
) error {
	path := fmt.Sprintf("http://%s/updates/", addr)
	metricsList := models.MetricsList{}

	for k, v := range metrics {
		value := v
		metric := models.Metrics{
			ID:    k,
			MType: models.Gauge,
			Value: &value,
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

	if err := send(jsonData, path, hashKey, cryptoKey); err != nil {
		return err
	}
	logger.Log.Info("Report metrics by grpc")
	if err := sendMetricsByGRPC(jsonData, grpcAddr); err != nil {
		return err
	}
	log.Printf("reported metrics in JSON %s\n", string(jsonData))

	return nil
}
