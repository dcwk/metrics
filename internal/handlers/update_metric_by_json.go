package handlers

import (
	"encoding/json"
	"net/http"

	"github.com/mailru/easyjson"

	"github.com/dcwk/metrics/internal/logger"
	"github.com/dcwk/metrics/internal/models"
	"github.com/dcwk/metrics/internal/service"
)

// UpdateMetricByJSON - обновляет одну метрику переданную в body в формате JSON
func (h *Handlers) UpdateMetricByJSON(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json;charset=UTF-8")
	w.Header().Set("Content-Encoding", "gzip")

	var metrics *models.Metrics
	err := json.NewDecoder(r.Body).Decode(&metrics)
	if err != nil {
		logger.Log.Error(err.Error())
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	metricsService := service.NewMetricsService(h.Storage)
	if err := metricsService.UpdateMetrics(metrics); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}

	metricsJSON, err := easyjson.Marshal(metrics)
	logger.Log.Info(string(metricsJSON))
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}

	w.WriteHeader(http.StatusOK)
	if err := json.NewEncoder(w).Encode(metrics); err != nil {
		logger.Log.Error(err.Error())
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}
