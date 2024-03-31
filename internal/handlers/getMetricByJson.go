package handlers

import (
	"encoding/json"
	"net/http"

	"github.com/dcwk/metrics/internal/models"
	"github.com/dcwk/metrics/internal/service"
	"github.com/mailru/easyjson"
)

func (h *Handlers) GetMetricByJson(w http.ResponseWriter, r *http.Request) {
	r.Method = http.MethodPost
	r.Header.Set("Content-Type", "application/json")
	var metrics *models.Metrics
	metricsService := service.NewMetricsService(h.Storage)
	err := json.NewDecoder(r.Body).Decode(&metrics)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	metrics, err = metricsService.GetMetrics(metrics)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	resp, err := easyjson.Marshal(metrics)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
	if _, err := w.Write(resp); err != nil {
		http.Error(w, "", http.StatusInternalServerError)
		return
	}
}
