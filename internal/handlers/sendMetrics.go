package handlers

import (
	"bytes"
	"compress/gzip"
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/dcwk/metrics/internal/models"
)

func (h *Handlers) SendMetrics(metrics map[string]float64, addr string, pollCount *int64) error {
	for k, v := range metrics {
		metric := models.Metrics{
			ID:    k,
			MType: models.Gauge,
			Value: &v,
		}
		//json, err := easyjson.Marshal(&metric)
		//if err != nil {
		//	return err
		//}

		//logger.Log.Info(string(json))

		if err := send(&metric, addr); err != nil {
			return err
		}
	}

	metric := models.Metrics{
		ID:    "PollCount",
		MType: models.Counter,
		Delta: pollCount,
	}

	//json, err := easyjson.Marshal(&metric)
	//if err != nil {
	//	return err
	//}

	if err := send(&metric, addr); err != nil {
		return err
	}

	return nil
}

func send(metric *models.Metrics, addr string) error {
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
	metricJSON, err := json.Marshal(metric)
	if err != nil {
		return err
	}

	endpoint := fmt.Sprintf("http://%s/update/", addr)
	client := http.Client{}
	req, err := http.NewRequest("POST", endpoint, bytes.NewBuffer(metricJSON))
	if err != nil {
		return err
	}
	req.Header.Set("Content-Type", "application/json")

	resp, err := client.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusOK {
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
