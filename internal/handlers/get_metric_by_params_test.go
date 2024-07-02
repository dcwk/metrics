package handlers_test

import (
	"fmt"
	"math/rand"
	"net/http"
	"net/http/httptest"
	"strconv"
	"testing"

	"github.com/go-resty/resty/v2"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/dcwk/metrics/internal/config"
	"github.com/dcwk/metrics/internal/server"
	"github.com/dcwk/metrics/internal/storage"
)

func TestGetMetricsByParams(t *testing.T) {
	s := storage.NewStorage()
	conf := &config.ServerConf{}

	ts := httptest.NewServer(server.Router(s, conf))
	defer ts.Close()
	path := ""
	id := strconv.Itoa(rand.Intn(256))
	count := 1000
	a := 0
	client := resty.New()

	for i := 0; i < count; i++ {
		v := rand.Intn(1024) + 1
		a += v
		path = "/update/counter/testSetGet" + id + "/" + strconv.Itoa(v)
		resp, err := client.R().
			SetHeader("Content-Type", "application/json").
			Post(ts.URL + path)
		require.NoError(t, err)
		assert.Equal(t, http.StatusOK, resp.StatusCode())

		path = "/value/counter/testSetGet" + id
		resp, err = client.R().
			SetHeader("Content-Type", "application/json").
			Get(ts.URL + path)

		require.NoError(t, err)
		assert.Equal(t, fmt.Sprintf("%d", a), string(resp.Body()))
	}
}
