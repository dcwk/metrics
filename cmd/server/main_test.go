package main

import (
	"fmt"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"io"
	"math"
	"math/rand"
	"net/http"
	"net/http/httptest"
	"strconv"
	"testing"
)

func testRequest(t *testing.T, ts *httptest.Server, method, path string) (*http.Response, string) {
	req, err := http.NewRequest(method, ts.URL+path, nil)
	require.NoError(t, err)

	resp, err := ts.Client().Do(req)
	require.NoError(t, err)
	defer resp.Body.Close()

	respBody, err := io.ReadAll(resp.Body)
	require.NoError(t, err)

	return resp, string(respBody)
}

func TestUpdateMetrics(t *testing.T) {
	ts := httptest.NewServer(Router())
	defer ts.Close()

	var testTable = []struct {
		name   string
		url    string
		want   string
		status int
	}{
		{
			"Test can save gauge",
			"/update/gauge/someMetric/527",
			"",
			http.StatusOK,
		},
		{
			"Test can save gauge with none value",
			"/update/gauge/testCounter/none",
			"\n",
			http.StatusBadRequest,
		},
		{
			"Test can save gauge without params",
			"/update/gauge/",
			"404 page not found\n",
			http.StatusNotFound,
		},
		{
			"Test can save counter1",
			"/update/counter/someMetric/527",
			"",
			http.StatusOK,
		},
		{
			"Test can save counter2",
			"/update/counter/testSetGet247/1965",
			"",
			http.StatusOK,
		},
		{
			"Test can save counter3",
			"/update/counter/testSetGet247/977",
			"",
			http.StatusOK,
		},
		{
			"Test can save counter with none value",
			"/update/counter/testCounter/none",
			"\n",
			http.StatusBadRequest,
		},
		{
			"Test can save counter without params",
			"/update/counter/",
			"404 page not found\n",
			http.StatusNotFound,
		},
		{
			"Test can post with unknown type",
			"/update/unknown/testCounter/100",
			"\n",
			http.StatusBadRequest,
		},
	}

	for _, tt := range testTable {
		t.Run(tt.name, func(t *testing.T) {
			resp, data := testRequest(t, ts, "POST", tt.url)
			defer resp.Body.Close()

			assert.Equal(t, tt.status, resp.StatusCode)
			assert.Equal(t, tt.want, data)
		})
	}
}

func TestGetMetrics(t *testing.T) {
	ts := httptest.NewServer(Router())
	defer ts.Close()
	path := ""
	id := strconv.Itoa(rand.Intn(256))
	count := 1000

	for i := 0; i < count; i++ {
		v := rand.Intn(1024)
		path = "/update/counter/testSetGet" + id + "/" + strconv.Itoa(v)
		testRequest(t, ts, "POST", path)

		path = "/value/counter/testSetGet" + id
		_, resp := testRequest(t, ts, "GET", path)

		assert.Equal(t, fmt.Sprintf("%d", v), resp)

		path = "/update/gauge/testSetGet" + id + "/" + strconv.Itoa(v)
		testRequest(t, ts, "POST", path)

		path = "/value/gauge/testSetGet" + id
		_, resp1 := testRequest(t, ts, "GET", path)
		val, err := strconv.ParseFloat(resp1, 64)
		if err != nil {
			fmt.Errorf(err.Error())
		}
		assert.Equal(t, fmt.Sprintf("%d", v), fmt.Sprintf("%.0f", math.Round(val)))
	}
}
