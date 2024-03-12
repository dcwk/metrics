package handlers

import (
	"fmt"
	"github.com/dcwk/metrics/internal/storage"
	"net/http"
)

const (
	gauge   = "gauge"
	counter = "counter"
)

func GetAllMetrics(w http.ResponseWriter, r *http.Request, s storage.Storage) {
	r.Method = http.MethodGet
	r.Header.Set("Content-Type", "text/plain")

	gauges := s.GetAllGauges()
	counters := s.GetAllCounters()
	res := ""

	for n, v := range gauges {
		res += n + " " + fmt.Sprintf("%.3f", v) + "\n\r"
	}

	for n, v := range counters {
		res += n + " " + fmt.Sprintf("%d", v) + "\n\r"
	}

	w.WriteHeader(http.StatusOK)
	if _, err := w.Write([]byte(res)); err != nil {
		http.Error(w, "", http.StatusInternalServerError)
		return
	}
}
