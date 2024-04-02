package handlers

import (
	"encoding/json"
	"net/http"

	"github.com/dcwk/metrics/internal/logger"
	"github.com/dcwk/metrics/internal/models"
	"github.com/dcwk/metrics/internal/service"
	"github.com/mailru/easyjson"
)

func (h *Handlers) UpdateMetricByJSON(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	var metrics *models.Metrics
	metricsService := service.NewMetricsService(h.Storage)
	err := json.NewDecoder(r.Body).Decode(&metrics)
	if err != nil {
		logger.Log.Error(err.Error())
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	metricsJSON, err := easyjson.Marshal(metrics)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
	logger.Log.Info(string(metricsJSON))

	switch metrics.MType {
	default:
		return
	case models.Gauge:
		if err := h.Storage.AddGauge(metrics.ID, metrics.Value); err != nil {
			return
		}
	case models.Counter:
		if err := h.Storage.AddCounter(metrics.ID, metrics.Delta); err != nil {
			return
		}
	}

	metrics, err = metricsService.GetMetrics(metrics)
	if err != nil {
		logger.Log.Error(err.Error())
		logger.Log.Info(string(metricsJSON))
		http.Error(w, "", http.StatusNotFound)
		return
	}

	w.WriteHeader(http.StatusOK)
	if err := json.NewEncoder(w).Encode(metrics); err != nil {
		logger.Log.Error(err.Error())
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}
