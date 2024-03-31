package handlers

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/dcwk/metrics/internal/logger"
	"github.com/dcwk/metrics/internal/models"
)

func (h *Handlers) UpdateMetric(w http.ResponseWriter, r *http.Request) {
	r.Method = http.MethodPost
	r.Header.Set("Content-Type", "text/plain")
	var metrics models.Metrics
	err := json.NewDecoder(r.Body).Decode(&metrics)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	logger.Log.Info(fmt.Sprintf("%v", metrics))

	switch metrics.MType {
	default:
		http.Error(w, "", http.StatusBadRequest)
		return
	case gauge:
		if err := h.Storage.AddGauge(metrics.ID, metrics.Value); err != nil {
			http.Error(w, "", http.StatusBadRequest)
			return
		}
	case counter:
		if err := h.Storage.AddCounter(metrics.ID, metrics.Delta); err != nil {
			http.Error(w, "", http.StatusBadRequest)
			return
		}
	}

	w.WriteHeader(http.StatusOK)
}
