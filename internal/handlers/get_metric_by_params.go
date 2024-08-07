package handlers

import (
	"fmt"
	"net/http"

	"github.com/go-chi/chi/v5"

	"github.com/dcwk/metrics/internal/models"
)

// GetMetricByParams - получение списка метрик с фильтром в query string
func (h *Handlers) GetMetricByParams(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json;charset=UTF-8")

	t := chi.URLParam(r, "type")
	n := chi.URLParam(r, "name")
	v := ""

	switch t {
	default:
		http.Error(w, "", http.StatusBadRequest)
		return
	case models.Gauge:
		metricValue, err := h.Storage.GetGauge(n)
		if err != nil {
			http.Error(w, "", http.StatusNotFound)
			return
		}
		v = fmt.Sprintf("%v", metricValue)
	case models.Counter:
		metricValue, err := h.Storage.GetCounter(n)
		if err != nil {
			http.Error(w, "", http.StatusNotFound)
			return
		}
		v = fmt.Sprintf("%d", metricValue)
	}

	w.WriteHeader(http.StatusOK)
	if _, err := w.Write([]byte(v)); err != nil {
		http.Error(w, "", http.StatusInternalServerError)
		return
	}
}
