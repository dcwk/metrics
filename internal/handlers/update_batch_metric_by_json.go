package handlers

import (
	"bytes"
	"encoding/json"
	"net/http"

	"github.com/mailru/easyjson"

	"github.com/dcwk/metrics/internal/logger"
	"github.com/dcwk/metrics/internal/models"
)

// UpdateBatchMetricByJSON - обновляет все метрики переданные в body в формате JSON
func (h *Handlers) UpdateBatchMetricByJSON(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json;charset=UTF-8")
	w.Header().Set("Content-Encoding", "gzip")

	var metricsData []models.Metrics
	var buf bytes.Buffer
	_, err := buf.ReadFrom(r.Body)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	if err = json.Unmarshal(buf.Bytes(), &metricsData); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	metricsList := &models.MetricsList{
		List: metricsData,
	}

	if err := h.Storage.AddMetricsAtBatchMode(metricsList); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}

	metricsJSON, err := easyjson.Marshal(metricsList)
	logger.Log.Info(string(metricsJSON))
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}

	w.WriteHeader(http.StatusOK)
	if err := json.NewEncoder(w).Encode(metricsList); err != nil {
		logger.Log.Error(err.Error())
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}
