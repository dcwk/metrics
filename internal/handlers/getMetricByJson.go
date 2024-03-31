package handlers

import (
	"encoding/json"
	"net/http"

	"github.com/dcwk/metrics/internal/logger"
	"github.com/dcwk/metrics/internal/models"
	"github.com/dcwk/metrics/internal/service"
	"github.com/mailru/easyjson"
)

func (h *Handlers) GetMetricByJSON(w http.ResponseWriter, r *http.Request) {
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
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	resp, err := easyjson.Marshal(metrics)
	if err != nil {
		logger.Log.Error(err.Error())
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	if _, err := w.Write(resp); err != nil {
		logger.Log.Error(err.Error())
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
}
