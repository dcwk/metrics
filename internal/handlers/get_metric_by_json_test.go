package handlers_test

import (
	"math/rand"
	"net/http"
	"net/http/httptest"
	"strconv"
	"testing"

	"github.com/go-resty/resty/v2"
	"github.com/mailru/easyjson"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/dcwk/metrics/internal/config"
	"github.com/dcwk/metrics/internal/models"
	"github.com/dcwk/metrics/internal/server"
	"github.com/dcwk/metrics/internal/storage"
)

func TestGetCounterMetricsByJson(t *testing.T) {
	s := storage.NewStorage()
	conf := &config.ServerConf{}

	ts := httptest.NewServer(server.Router(s, conf))
	defer ts.Close()
	path := ""
	id := strconv.Itoa(rand.Intn(256))
	count := 1000
	var a int64
	client := resty.New()

	for i := 0; i < count; i++ {
		v := rand.Intn(1024) + 1
		a += int64(v)
		metricID := "testSetGet" + id
		metricVal := int64(v)
		metricsPost := &models.Metrics{
			ID:    metricID,
			MType: "counter",
			Delta: &metricVal,
		}
		path = "/update/"
		body, err := easyjson.Marshal(metricsPost)
		assert.NoError(t, err)
		resp, err := client.R().
			SetHeader("Content-Type", "application/json").
			SetBody(body).
			Post(ts.URL + path)
		require.NoError(t, err)
		status := resp.StatusCode()
		assert.Equal(t, http.StatusOK, status)

		metricsGet := &models.Metrics{
			ID:    metricID,
			MType: "counter",
		}
		path = "/value/"
		body, err = easyjson.Marshal(metricsGet)
		require.NoError(t, err)
		resp, err = client.R().
			SetHeader("Content-Type", "application/json").
			SetBody(body).
			Post(ts.URL + path)

		require.NoError(t, err)
		data := string(resp.Body())
		expMetrics := &models.Metrics{
			ID:    metricID,
			MType: "counter",
			Delta: &a,
		}
		expJSON, err := easyjson.Marshal(expMetrics)
		assert.NoError(t, err)
		assert.Equal(t, string(expJSON)+"\n", data)
	}
}

func TestGetGaugeMetricsByJson(t *testing.T) {
	s := storage.NewStorage()
	conf := &config.ServerConf{}

	ts := httptest.NewServer(server.Router(s, conf))
	defer ts.Close()
	path := ""
	id := strconv.Itoa(rand.Intn(256))
	count := 1000
	client := resty.New()

	for i := 0; i < count; i++ {
		v := rand.Intn(1024) + 1
		metricID := "testSetGet" + id
		metricVal := float64(v)
		metricsPost := &models.Metrics{
			ID:    metricID,
			MType: "gauge",
			Value: &metricVal,
		}
		path = "/update/"
		body, err := easyjson.Marshal(metricsPost)
		assert.NoError(t, err)
		resp, err := client.R().
			SetHeader("Content-Type", "application/json").
			SetBody(body).
			Post(ts.URL + path)
		require.NoError(t, err)
		status := resp.StatusCode()
		assert.Equal(t, http.StatusOK, status)

		metricsGet := &models.Metrics{
			ID:    metricID,
			MType: "gauge",
		}
		path = "/value/"
		body, err = easyjson.Marshal(metricsGet)
		require.NoError(t, err)
		resp, err = client.R().
			SetHeader("Content-Type", "application/json").
			SetBody(body).
			Post(ts.URL + path)

		require.NoError(t, err)
		data := string(resp.Body())
		expMetrics := &models.Metrics{
			ID:    metricID,
			MType: "gauge",
			Value: &metricVal,
		}
		expJSON, err := easyjson.Marshal(expMetrics)
		assert.NoError(t, err)
		assert.Equal(t, string(expJSON)+"\n", data)
	}
}
