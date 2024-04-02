package client

import (
	"math/rand"
	"runtime"
	"time"

	"github.com/dcwk/metrics/internal/logger"
)

type Agent struct {
	Metrics        map[string]float64
	PollCount      int64
	PollInterval   time.Duration
	ReportInterval time.Duration
}

func NewAgent(pollInterval int64, reportInterval int64) *Agent {
	return &Agent{
		Metrics:        make(map[string]float64),
		PollInterval:   time.Duration(pollInterval),
		ReportInterval: time.Duration(reportInterval),
	}
}

func (a *Agent) Update() {
	for {
		logger.Log.Info("start read metrics for agent update")
		var ms runtime.MemStats
		runtime.ReadMemStats(&ms)
		a.PollCount++

		a.Metrics["Alloc"] = float64(ms.Alloc)
		a.Metrics["GCCPUFraction"] = float64(ms.GCCPUFraction)
		a.Metrics["BuckHashSys"] = float64(ms.BuckHashSys)
		a.Metrics["Frees"] = float64(ms.Frees)
		a.Metrics["GCSys"] = float64(ms.GCSys)
		a.Metrics["HeapAlloc"] = float64(ms.HeapAlloc)
		a.Metrics["HeapIdle"] = float64(ms.HeapIdle)
		a.Metrics["HeapInuse"] = float64(ms.HeapInuse)
		a.Metrics["HeapObjects"] = float64(ms.HeapObjects)
		a.Metrics["HeapReleased"] = float64(ms.HeapReleased)
		a.Metrics["HeapSys"] = float64(ms.HeapSys)
		a.Metrics["LastGC"] = float64(ms.LastGC)
		a.Metrics["Lookups"] = float64(ms.Lookups)
		a.Metrics["MCacheInuse"] = float64(ms.MCacheInuse)
		a.Metrics["MCacheSys"] = float64(ms.MCacheSys)
		a.Metrics["MSpanInuse"] = float64(ms.MSpanInuse)
		a.Metrics["MSpanSys"] = float64(ms.MSpanSys)
		a.Metrics["Mallocs"] = float64(ms.Mallocs)
		a.Metrics["NextGC"] = float64(ms.NextGC)
		a.Metrics["NumForcedGC"] = float64(ms.NumForcedGC)
		a.Metrics["NumGC"] = float64(ms.NumGC)
		a.Metrics["OtherSys"] = float64(ms.OtherSys)
		a.Metrics["PauseTotalNs"] = float64(ms.PauseTotalNs)
		a.Metrics["StackInuse"] = float64(ms.StackInuse)
		a.Metrics["StackSys"] = float64(ms.StackSys)
		a.Metrics["Sys"] = float64(ms.Sys)
		a.Metrics["TotalAlloc"] = float64(ms.TotalAlloc)
		a.Metrics["RandomValue"] = float64(rand.Intn(1024) + 1)

		time.Sleep(time.Duration(a.PollInterval) * time.Second)
	}
}
