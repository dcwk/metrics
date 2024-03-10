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
		url    string
		want   string
		status int
	}{
		{"/update/gauge/someMetric/527", "", http.StatusOK},
		{"/update/gauge/testCounter/none", "", http.StatusBadRequest},
	}

	for _, tt := range testTable {
		resp, data := testRequest(t, ts, "POST", tt.url)
		assert.Equal(t, tt.status, resp.StatusCode)
		assert.Equal(t, tt.want, data)
	}
}
