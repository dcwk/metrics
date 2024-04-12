package handlers

import (
	"bytes"
	"fmt"
	"net/http"
	"runtime"
)

func SendMetrics(addr string) error {
	for k, v := range getGauges() {
		r := bytes.NewReader([]byte(""))
		resp, err := http.Post(
			fmt.Sprintf("http://%s/update/gauge/%s/%f", addr, k, v),
			"Content-Type: text/plain",
			r,
		)
		if err != nil {
			return err
		}

		if err := resp.Body.Close(); err != nil {
			return err
		}
	}

	return nil
}

func getGauges() map[string]float64 {
	gauges := map[string]float64{}
	ms := runtime.MemStats{}
	runtime.ReadMemStats(&ms)

	gauges["Alloc"] = float64(ms.Alloc)
	gauges["BuckHashSys"] = float64(ms.BuckHashSys)
	gauges["Frees"] = float64(ms.Frees)
	gauges["GCSys"] = float64(ms.GCSys)
	gauges["HeapAlloc"] = float64(ms.HeapAlloc)
	gauges["HeapIdle"] = float64(ms.HeapIdle)
	gauges["HeapInuse"] = float64(ms.HeapInuse)
	gauges["HeapObjects"] = float64(ms.HeapObjects)
	gauges["HeapReleased"] = float64(ms.HeapReleased)
	gauges["HeapSys"] = float64(ms.HeapSys)
	gauges["LastGC"] = float64(ms.LastGC)
	gauges["Lookups"] = float64(ms.Lookups)
	gauges["MCacheInuse"] = float64(ms.MCacheInuse)
	gauges["MCacheSys"] = float64(ms.MCacheSys)
	gauges["MSpanInuse"] = float64(ms.MSpanInuse)
	gauges["MSpanSys"] = float64(ms.MSpanSys)
	gauges["Mallocs"] = float64(ms.Mallocs)
	gauges["NextGC"] = float64(ms.NextGC)
	gauges["NumForcedGC"] = float64(ms.NumForcedGC)
	gauges["NumGC"] = float64(ms.NumGC)
	gauges["OtherSys"] = float64(ms.OtherSys)
	gauges["PauseTotalNs"] = float64(ms.PauseTotalNs)
	gauges["StackInuse"] = float64(ms.StackInuse)
	gauges["StackSys"] = float64(ms.StackSys)
	gauges["Sys"] = float64(ms.Sys)
	gauges["TotalAlloc"] = float64(ms.TotalAlloc)

	return gauges
}
