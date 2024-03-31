package handlers

import (
	"net/http"
	"strconv"
	"strings"

	"github.com/dcwk/metrics/internal/models"
	"github.com/dcwk/metrics/internal/service"
	"github.com/go-chi/chi/v5"
)

func (h *Handlers) UpdateMetricByParams(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	metricsService := service.NewMetricsService(h.Storage)

	t := chi.URLParam(r, "type")
	mn := chi.URLParam(r, "name")
	mv := chi.URLParam(r, "value")
	metrics := &models.Metrics{
		ID:    mn,
		MType: t,
	}

	if metrics.MType == gauge {
		convertedVal, err := strconv.ParseFloat(strings.TrimSpace(mv), 64)
		if err != nil {
			http.Error(w, "unsupported gauge value", http.StatusBadRequest)
			return
		}
		metrics.Value = &convertedVal
	} else if metrics.MType == counter {
		convertedVal, err := strconv.ParseInt(strings.TrimSpace(mv), 10, 64)
		if err != nil {
			http.Error(w, "unsupported counter value", http.StatusBadRequest)
			return
		}
		metrics.Delta = &convertedVal
	} else {
		http.Error(w, "unsupported metric type", http.StatusBadRequest)
		return
	}

	if err := metricsService.UpdateMetrics(metrics); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
	}

	w.WriteHeader(http.StatusOK)
}
