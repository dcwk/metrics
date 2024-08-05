package client

import (
	"bytes"
	"compress/gzip"
	"context"
	"crypto/hmac"
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"log"

	"github.com/go-resty/resty/v2"
	"github.com/mailru/easyjson"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"

	"github.com/dcwk/metrics/internal/grpchandler"
	"github.com/dcwk/metrics/internal/logger"
	"github.com/dcwk/metrics/internal/models"
	"github.com/dcwk/metrics/internal/utils/crypt"
)

func SendMetrics(metrics map[string]float64, addr string, hashKey string, cryptoKey string, pollCount *int64) error {
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

		if err := send(json, path, hashKey, cryptoKey); err != nil {
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

	if err := send(json, path, hashKey, cryptoKey); err != nil {
		return err
	}

	return nil
}

func send(metricsJSON []byte, path string, hashKey string, cryptoKey string) error {
	var sign []byte
	if hashKey != "" {
		h := hmac.New(sha256.New, []byte(hashKey))
		sign = h.Sum(metricsJSON)
	}

	body, err := compress(metricsJSON)
	if err != nil {
		logger.Log.Fatal(fmt.Sprintf("Failed to compress metrics: %s", err))
	}
	body, err = crypt.EncryptWithRSA(body, cryptoKey)
	if err != nil {
		logger.Log.Fatal(fmt.Sprintf("Failed to encrypt metrics: %s", err))
	}
	client := resty.New()
	_, err = client.R().
		SetHeaders(map[string]string{
			"Content-Type":     "application/json;charset=UTF-8",
			"Accept-Encoding":  "gzip",
			"Content-Encoding": "gzip",
			"HashSHA256":       hex.EncodeToString(sign),
			"X-Real-IP":        "127.0.0.1",
		}).
		SetBody(string(body)).
		Post(path)
	if err != nil {
		logger.Log.Error(fmt.Sprintf("Can't send request to server: %s", err.Error()))
		return err
	}

	return nil
}

func sendMetricsByGRPC(metricsJSON []byte, grpcServerAddr string) error {
	conn, err := grpc.NewClient(grpcServerAddr, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		log.Fatal(err)
	}
	defer conn.Close()

	c := grpchandler.NewMetricsServiceClient(conn)

	respUpdate, err := c.UpdateBatchMetricByJSON(
		context.Background(),
		&grpchandler.UpdateBatchMetricByJSONRequest{
			Metrics: string(metricsJSON),
		},
	)
	if err != nil {
		return err
	}
	logger.Log.Info(fmt.Sprintf("success update metrics data: %s", respUpdate))

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
