package main

import (
	"math/rand"
	"net/http"
	"net/http/httptest"
	"strconv"
	"testing"

	"github.com/dcwk/metrics/internal/logger"
	"github.com/dcwk/metrics/internal/models"
	"github.com/dcwk/metrics/internal/server"
	"github.com/dcwk/metrics/internal/storage"
	"github.com/go-resty/resty/v2"
	"github.com/mailru/easyjson"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestUpdateMetricsByParams(t *testing.T) {
	s := storage.NewStorage()
	if err := logger.Initialize("info"); err != nil {
		panic(err)
	}
	ts := httptest.NewServer(server.Router(s))
	defer ts.Close()

	var testTable = []struct {
		name        string
		url         string
		bodyString  string
		contentType string
		want        string
		status      int
	}{
		{
			"Test can save gauge with none value",
			"/update/gauge/testCounter/none",
			``,
			``,
			"\n",
			http.StatusBadRequest,
		},
		{
			"Test can save gauge without params",
			"/update/gauge/",
			``,
			``,
			"404 page not found\n",
			http.StatusNotFound,
		},
		{
			"Test can save counter1",
			"/update/counter/someMetric/527",
			``,
			``,
			"",
			http.StatusOK,
		},
		{
			"Test can save counter2",
			"/update/counter/testSetGet247/1965",
			``,
			``,
			"",
			http.StatusOK,
		},
		{
			"Test can save counter3",
			"/update/counter/testSetGet247/977",
			``,
			``,
			"",
			http.StatusOK,
		},
		{
			"Test can save counter with none value",
			"/update/counter/testCounter/none",
			``,
			``,
			"\n",
			http.StatusBadRequest,
		},
		{
			"Test can save counter without params",
			"/update/counter/",
			``,
			``,
			"404 page not found\n",
			http.StatusNotFound,
		},
		{
			"Test can post with unknown type",
			"/update/unknown/testCounter/100",
			``,
			``,
			"\n",
			http.StatusBadRequest,
		},
	}

	client := resty.New()
	for _, tt := range testTable {
		t.Run(tt.name, func(t *testing.T) {
			resp, err := client.R().
				SetHeader("Content-Type", tt.contentType).
				SetBody([]byte(tt.bodyString)).
				Post(ts.URL + tt.url)
			require.NoError(t, err)
			assert.Equal(t, tt.status, resp.StatusCode())
		})
	}
}

func TestUpdateMetricsByJson(t *testing.T) {
	s := storage.NewStorage()
	if err := logger.Initialize("info"); err != nil {
		panic(err)
	}
	ts := httptest.NewServer(server.Router(s))
	defer ts.Close()

	var testTable = []struct {
		name        string
		url         string
		bodyString  string
		contentType string
		want        string
		status      int
	}{
		{
			"Test can save gauge",
			"/update/",
			`{"id":"StackInuse","type":"gauge","value":327680}`,
			"application/json",
			`{"id":"StackInuse","type":"gauge","value":327680}`,
			http.StatusOK,
		},
	}

	client := resty.New()
	for _, tt := range testTable {
		t.Run(tt.name, func(t *testing.T) {
			resp, err := client.R().
				SetHeader("Content-Type", tt.contentType).
				SetBody([]byte(tt.bodyString)).
				Post(ts.URL + tt.url)
			require.NoError(t, err)
			assert.Equal(t, tt.status, resp.StatusCode())
			assert.Equal(t, tt.want+"\n", string(resp.Body()))
		})
	}
}

//func TestGetMetricsByParams(t *testing.T) {
//	s := storage.NewStorage()
//	ts := httptest.NewServer(server.Router(s))
//	defer ts.Close()
//	path := ""
//	id := strconv.Itoa(rand.Intn(256))
//	count := 1000
//	a := 0
//	client := resty.New()
//
//	for i := 0; i < count; i++ {
//		v := rand.Intn(1024) + 1
//		a += v
//		path = "/update/counter/testSetGet" + id + "/" + strconv.Itoa(v)
//		resp, err := client.R().
//			SetHeader("Content-Type", "text/html").
//			Post(ts.URL + path)
//		require.NoError(t, err)
//		assert.Equal(t, http.StatusOK, resp.StatusCode())
//
//		path = "/value/counter/testSetGet" + id
//		resp, err = client.R().
//			SetHeader("Content-Type", "text/html").
//			Get(ts.URL + path)
//
//		require.NoError(t, err)
//		assert.Equal(t, fmt.Sprintf("%d", a), string(resp.Body()))
//	}
//}

func TestGetCounterMetricsByJson(t *testing.T) {
	s := storage.NewStorage()
	ts := httptest.NewServer(server.Router(s))
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
	ts := httptest.NewServer(server.Router(s))
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
