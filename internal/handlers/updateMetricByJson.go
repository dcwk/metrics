package handlers

import (
	"encoding/json"
	"net/http"

	"github.com/dcwk/metrics/internal/models"
	"github.com/dcwk/metrics/internal/service"
)

func (h *Handlers) UpdateMetricByJson(w http.ResponseWriter, r *http.Request) {
	r.Method = http.MethodPost
	r.Header.Set("Content-Type", "text/plain")
	var metrics *models.Metrics
	metricsService := service.NewMetricsService(h.Storage)
	err := json.NewDecoder(r.Body).Decode(&metrics)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	if err := metricsService.UpdateMetrics(metrics); err != nil {
		http.Error(w, "", http.StatusBadRequest)
		return
	}

	w.WriteHeader(http.StatusOK)
}
