package main

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/dcwk/metrics/internal/logger"
	"github.com/dcwk/metrics/internal/server"
	"github.com/dcwk/metrics/internal/storage"
	"github.com/go-resty/resty/v2"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestUpdateMetrics(t *testing.T) {
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
			"/update",
			`{"id":"StackInuse","type":"gauge","value":327680}`,
			"application/json",
			"",
			http.StatusOK,
		},
		//{
		//	"Test can save gauge with none value",
		//	"/update/gauge/testCounter/none",
		//	"\n",
		//	http.StatusBadRequest,
		//},
		//{
		//	"Test can save gauge without params",
		//	"/update/gauge/",
		//	"404 page not found\n",
		//	http.StatusNotFound,
		//},
		//{
		//	"Test can save counter1",
		//	"/update/counter/someMetric/527",
		//	"",
		//	http.StatusOK,
		//},
		//{
		//	"Test can save counter2",
		//	"/update/counter/testSetGet247/1965",
		//	"",
		//	http.StatusOK,
		//},
		//{
		//	"Test can save counter3",
		//	"/update/counter/testSetGet247/977",
		//	"",
		//	http.StatusOK,
		//},
		//{
		//	"Test can save counter with none value",
		//	"/update/counter/testCounter/none",
		//	"\n",
		//	http.StatusBadRequest,
		//},
		//{
		//	"Test can save counter without params",
		//	"/update/counter/",
		//	"404 page not found\n",
		//	http.StatusNotFound,
		//},
		//{
		//	"Test can post with unknown type",
		//	"/update/unknown/testCounter/100",
		//	"\n",
		//	http.StatusBadRequest,
		//},
	}

	client := resty.New()
	for _, tt := range testTable {
		t.Run(tt.name, func(t *testing.T) {
			resp, err := client.R().
				SetHeader("Content-Type", tt.contentType).
				SetBody([]byte(tt.bodyString)).
				Post(ts.URL + tt.url)
			require.NoError(t, err)

			assert.Equal(t, tt.status, resp.StatusCode)
			assert.Equal(t, tt.want, string(resp.Body()))
		})
	}
}

//func TestGetMetrics(t *testing.T) {
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
//		r, _ := testRequest(t, ts, "POST", path)
//		if err := r.Body.Close(); err != nil {
//			fmt.Println(err.Error())
//		}
//
//		path = "/value/counter/testSetGet" + id
//		resp, err := client.R().
//			SetHeader("Content-Type", "Content-Type").
//			Get(path)
//		require.NoError(t, err)
//		assert.Equal(t, fmt.Sprintf("%d", a), string(resp.Body()))
//	}
//}
