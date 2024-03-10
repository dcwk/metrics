package main

import (
	"bytes"
	"flag"
	"fmt"
	"log"
	"net/http"
	"os"
	"runtime"
	"strconv"
	"time"
)

var (
	flagServerAddr string
	reportInterval int64
	pollInterval   int64
)

func main() {
	parseFlags()
	fmt.Println("Sending metrics to", flagServerAddr)
	var pollCount int64 = 0

	for {
		time.Sleep(time.Duration(pollInterval) * time.Second)
		gauges := getGauges()

		if pollCount%reportInterval == 0 {
			for k, v := range gauges {
				r := bytes.NewReader([]byte(""))
				resp, err := http.Post(
					fmt.Sprintf("http://%s/update/gauge/%s/%f", flagServerAddr, k, v),
					"Content-Type: text/plain",
					r,
				)
				if err != nil {
					log.Fatalln(err)
				}

				err = resp.Body.Close()
				if err != nil {
					log.Fatalln(err)
				}
			}
		}

		pollCount++
	}
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

func parseFlags() {
	flag.StringVar(&flagServerAddr, "a", ":8080", "metrics server address")
	flag.Int64Var(&reportInterval, "r", 10, "sending frequency interval")
	flag.Int64Var(&pollInterval, "p", 2, "metrics reading frequency")

	flag.Parse()

	if envAddress := os.Getenv("ADDRESS"); envAddress != "" {
		flagServerAddr = envAddress
	}

	if envReportInterval := os.Getenv("REPORT_INTERVAL"); envReportInterval != "" {
		interval, err := strconv.ParseInt(envReportInterval, 10, 64)
		if err != nil {
			log.Fatalln(err)
		}
		reportInterval = interval
	}

	if envPollInterval := os.Getenv("POLL_INTERVAL"); envPollInterval != "" {
		interval, err := strconv.ParseInt(envPollInterval, 10, 64)
		if err != nil {
			log.Fatalln(err)
		}
		pollInterval = interval
	}
}
