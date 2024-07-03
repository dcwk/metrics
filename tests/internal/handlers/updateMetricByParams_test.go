package handlers

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/go-resty/resty/v2"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/dcwk/metrics/internal/config"
	"github.com/dcwk/metrics/internal/logger"
	"github.com/dcwk/metrics/internal/server"
	"github.com/dcwk/metrics/internal/storage"
)

func TestUpdateMetricsByParams(t *testing.T) {
	s := storage.NewStorage()
	if err := logger.Initialize("info"); err != nil {
		panic(err)
	}
	conf := &config.ServerConf{}
	ts := httptest.NewServer(server.Router(s, conf))
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
