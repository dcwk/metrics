package main

import (
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"io"
	"net/http"
	"net/http/httptest"
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

func TestServer(t *testing.T) {
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
			"Test can save counter",
			"/update/counter/someMetric/527",
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
