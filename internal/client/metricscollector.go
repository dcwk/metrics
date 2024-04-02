package client

import (
	"math/rand"
	"runtime"
)

func getGauges(pollCount *int64) map[string]float64 {
	gauges := map[string]float64{}
	var ms runtime.MemStats
	runtime.ReadMemStats(&ms)

	gauges["Alloc"] = float64(ms.Alloc)
	gauges["GCCPUFraction"] = float64(ms.GCCPUFraction)
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
