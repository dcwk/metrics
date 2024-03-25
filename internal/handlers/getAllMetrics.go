package handlers

import (
	"fmt"
	"net/http"
)

func (h *Handlers) GetAllMetrics(w http.ResponseWriter, r *http.Request) {
	r.Method = http.MethodGet
	r.Header.Set("Content-Type", "text/plain")

	gauges := h.Storage.GetAllGauges()
	counters := h.Storage.GetAllCounters()
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
