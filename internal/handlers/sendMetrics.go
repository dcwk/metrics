package handlers

import (
	"bytes"
	"compress/gzip"
	"fmt"
	"log"

	"github.com/dcwk/metrics/internal/models"
	"github.com/go-resty/resty/v2"
	"github.com/mailru/easyjson"
)

func (h *Handlers) SendMetrics(metrics map[string]float64, addr string, pollCount *int64) error {
	path := fmt.Sprintf("http://%s/update/", addr)

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

		if err := send(json, path); err != nil {
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

	if err := send(json, path); err != nil {
		return err
	}

	return nil
}

func send(metricsJSON []byte, path string) error {
	body, err := compress(metricsJSON)
	if err != nil {
		return err
	}

	client := resty.New()
	_, err = client.R().
		SetHeaders(map[string]string{
			"Content-Type":     "application/json;charset=UTF-8",
			"Accept-Encoding":  "gzip",
			"Content-Encoding": "gzip",
		}).
		SetBody(string(body)).
		Post(path)
	if err != nil {
		return err
	}

	return nil
}

func compress(b []byte) ([]byte, error) {
	buf := bytes.NewBuffer(nil)
	gz := gzip.NewWriter(buf)

	_, err := gz.Write(b)
	if err != nil {
		return nil, err
	}
	if err := gz.Close(); err != nil {
		return nil, err
	}

	return buf.Bytes(), nil
}
