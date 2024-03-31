package handlers

import (
	"fmt"
	"net/http"

	"github.com/go-chi/chi/v5"
)

func (h *Handlers) GetMetricByParams(w http.ResponseWriter, r *http.Request) {
	r.Method = http.MethodGet
	r.Header.Set("Content-Type", "text/plain")

	t := chi.URLParam(r, "type")
	n := chi.URLParam(r, "name")
	v := ""

	switch t {
	default:
		http.Error(w, "", http.StatusBadRequest)
		return
	case gauge:
		metricValue, err := h.Storage.GetGauge(n)
		if err != nil {
			http.Error(w, "", http.StatusNotFound)
			return
		}
		v = fmt.Sprintf("%v", metricValue)
	case counter:
		metricValue, err := h.Storage.GetCounter(n)
		if err != nil {
			http.Error(w, "", http.StatusNotFound)
			return
		}
		v = fmt.Sprintf("%d", metricValue)
	}

	w.WriteHeader(http.StatusOK)
	if _, err := w.Write([]byte(v)); err != nil {
		http.Error(w, "", http.StatusInternalServerError)
		return
	}
}
