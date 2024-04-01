package handlers

import (
	"fmt"
	"math/rand"
	"runtime"

	"github.com/dcwk/metrics/internal/models"
	"github.com/go-resty/resty/v2"
	"github.com/mailru/easyjson"
)

func (h *Handlers) SendMetrics(addr string, pollCount int64) error {
	for k, v := range getGauges() {
		metric := models.Metrics{
			ID:    k,
			MType: gauge,
			Value: &v,
		}
		json, err := easyjson.Marshal(&metric)
		if err != nil {
			return err
		}

		//logger.Log.Info(string(json))

		if err := send(string(json), addr); err != nil {
			return err
		}
	}

	metric := models.Metrics{
		ID:    "PollCount",
		MType: counter,
		Delta: &pollCount,
	}

	json, err := easyjson.Marshal(&metric)
	if err != nil {
		return err
	}

	if err := send(string(json), addr); err != nil {
		return err
	}

	return nil
}

func send(metricsJSON string, addr string) error {
	client := resty.New()
	_, err := client.R().
		SetHeader("Content-Type", "application/json").
		SetBody(metricsJSON).
		Post(fmt.Sprintf("http://%s/update/", addr))
	if err != nil {
		return err
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
	gauges["RandomValue"] = float64(rand.Intn(1024) + 1)

	return gauges
}
