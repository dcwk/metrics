package handlers

import (
	"bytes"
	"compress/gzip"
	"fmt"
	"net/http"
	"strings"

	"github.com/dcwk/metrics/internal/logger"
	"github.com/dcwk/metrics/internal/models"
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

		if err := send(json, addr); err != nil {
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

	if err := send(json, addr); err != nil {
		return err
	}

	return nil
}

func send(metricsJSON []byte, addr string) error {
	//body, err := compress(metricsJSON)
	//if err != nil {
	//	return err
	//}

	//client := resty.New()
	//_, err = client.R().
	//	SetHeaders(map[string]string{
	//		"Content-Type":     "application/json;charset=UTF-8",
	//		"Accept-Encoding":  "gzip",
	//		"Content-Encoding": "gzip",
	//	}).
	//	SetBody(string(body)).
	//	Post(fmt.Sprintf("http://%s/update/", addr))
	request, err := http.NewRequest(
		http.MethodPost,
		fmt.Sprintf("http://%s/update/", addr),
		strings.NewReader(string(metricsJSON)),
	)
	if err != nil {
		return err
	}
	request.Header.Set("Content-Type", "application/json;charset=UTF-8")
	//request.Header.Set("Content-Encoding", "gzip")
	client := &http.Client{}
	response, err := client.Do(request)
	if err := response.Body.Close(); err != nil {
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
