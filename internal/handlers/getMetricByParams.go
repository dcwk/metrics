package handlers

import (
	"encoding/json"
	"net/http"

	"github.com/dcwk/metrics/internal/logger"
	"github.com/dcwk/metrics/internal/models"
	"github.com/dcwk/metrics/internal/service"
	"github.com/go-chi/chi/v5"
)

func (h *Handlers) GetMetricByParams(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	metricsService := service.NewMetricsService(h.Storage)

	t := chi.URLParam(r, "type")
	mn := chi.URLParam(r, "name")
	metrics := &models.Metrics{
		ID:    mn,
		MType: t,
	}

	metrics, err := metricsService.GetMetrics(metrics)
	if err != nil {
		logger.Log.Error(err.Error())
		http.Error(w, "", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
	if err := json.NewEncoder(w).Encode(metrics); err != nil {
		logger.Log.Error(err.Error())
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}
