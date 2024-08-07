package handlers

import (
	"encoding/json"
	"net/http"

	"github.com/mailru/easyjson"

	"github.com/dcwk/metrics/internal/logger"
	"github.com/dcwk/metrics/internal/models"
)

// GetMetricByJSON - получение конкретных метрик с фильтром в формате JSON
func (h *Handlers) GetMetricByJSON(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json;charset=UTF-8")
	w.Header().Set("Content-Encoding", "gzip")

	var metrics *models.Metrics
	err := json.NewDecoder(r.Body).Decode(&metrics)
	if err != nil {
		logger.Log.Error(err.Error())
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	switch metrics.MType {
	default:
		return
	case models.Gauge:
		metricValue, err := h.Storage.GetGauge(metrics.ID)
		if err != nil {
			http.Error(w, err.Error(), http.StatusNotFound)
			return
		}

		metrics.Value = &metricValue
	case models.Counter:
		metricValue, err := h.Storage.GetCounter(metrics.ID)
		if err != nil {
			http.Error(w, err.Error(), http.StatusNotFound)
			return
		}
		metrics.Delta = &metricValue
	}

	metricsJSON2, err := easyjson.Marshal(metrics)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
	logger.Log.Info(string(metricsJSON2))

	w.WriteHeader(http.StatusOK)
	if err := json.NewEncoder(w).Encode(metrics); err != nil {
		logger.Log.Error(err.Error())
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}
