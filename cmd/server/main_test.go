package main

import (
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"io"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestServer(t *testing.T) {
	type want struct {
		statusCode  int
		response    string
		contentType string
	}

	tests := []struct {
		name    string
		want    want
		request string
	}{
		{
			name: "Test can save gauge",
			want: want{
				statusCode:  http.StatusOK,
				response:    "",
				contentType: "text/plain",
			},
			request: "/update/gauge/someMetric/527",
		},
		{
			name: "Test fail unsupported value",
			want: want{
				statusCode:  http.StatusBadRequest,
				response:    "",
				contentType: "text/plain ",
			},
			request: "/update/gauge/testCounter/none",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			request := httptest.NewRequest(http.MethodPost, tt.request, nil)
			w := httptest.NewRecorder()
			GaugeHandler(w, request)

			res := w.Result()

			assert.Equal(t, tt.want.statusCode, res.StatusCode)
			defer res.Body.Close()
			_, err := io.ReadAll(res.Body)
			require.NoError(t, err)
		})
	}
}
