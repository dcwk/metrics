package client

import (
	"bytes"
	"compress/gzip"
	"crypto/hmac"
	"crypto/rand"
	"crypto/rsa"
	"crypto/sha256"
	"crypto/x509"
	"encoding/hex"
	"encoding/pem"
	"fmt"
	"log"
	"os"

	"github.com/go-resty/resty/v2"
	"github.com/mailru/easyjson"

	"github.com/dcwk/metrics/internal/logger"
	"github.com/dcwk/metrics/internal/models"
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
		return err
	}

	body, err = encrypt(body, cryptoKey)
	if err != nil {
		return err
	}

	client := resty.New()
	_, err = client.R().
		SetHeaders(map[string]string{
			"Content-Type":     "application/json;charset=UTF-8",
			"Accept-Encoding":  "gzip",
			"Content-Encoding": "gzip",
			"HashSHA256":       hex.EncodeToString(sign),
		}).
		SetBody(string(body)).
		Post(path)
	if err != nil {
		logger.Log.Error(fmt.Sprintf("Can't send request to server: %s", err.Error()))
		return err
	}

	return nil
}

func encrypt(b []byte, cryptoKeyPath string) ([]byte, error) {
	publicKeyPEM, err := os.ReadFile(cryptoKeyPath)
	if err != nil {
		panic(err)
	}

	publicKeyBlock, _ := pem.Decode(publicKeyPEM)
	publicKey, err := x509.ParsePKIXPublicKey(publicKeyBlock.Bytes)
	if err != nil {
		panic(err)
	}

	cipherText, err := rsa.EncryptPKCS1v15(rand.Reader, publicKey.(*rsa.PublicKey), b)
	if err != nil {
		return nil, err
	}

	return cipherText, nil
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
