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

func TestUpdateMetricsByJson(t *testing.T) {
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
