package client

import (
	"fmt"
	"log"

	"github.com/dcwk/metrics/internal/logger"
	"github.com/dcwk/metrics/internal/models"
	"github.com/mailru/easyjson"
)

func SendMetricsInPool(
	metrics map[string]float64,
	addr string,
	hashKey string,
	rateLimit int64,
	pollCount *int64,
) error {
	path := fmt.Sprintf("http://%s/update/", addr)
	workerPool := NewWorkerPool(rateLimit)
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
		log.Printf("reported metric JSON %s with value %f\n", k, v)

		workerPool.Produce(func() {
			if err := send(json, path, hashKey); err != nil {
				logger.Log.Error(err.Error())
			}
		})
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

	workerPool.Produce(func() {
		if err := send(json, path, hashKey); err != nil {
			logger.Log.Error(err.Error())
		}
	})

	return nil
}
