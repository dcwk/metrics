package handlers

import (
	"encoding/json"
	"net/http"

	"github.com/dcwk/metrics/internal/logger"
	"github.com/dcwk/metrics/internal/models"
	"github.com/dcwk/metrics/internal/service"
)

func (h *Handlers) GetMetricByJSON(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	var metrics *models.Metrics
	metricsService := service.NewMetricsService(h.Storage)
	err := json.NewDecoder(r.Body).Decode(&metrics)
	if err != nil {
		logger.Log.Error(err.Error())
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	metrics, err = metricsService.GetMetrics(metrics)
	if err != nil {
		logger.Log.Error(err.Error())
		http.Error(w, err.Error(), http.StatusNotFound)
		return
	}

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(metrics)
}
