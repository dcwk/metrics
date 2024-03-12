package handlers

import (
	"fmt"
	"github.com/dcwk/metrics/internal/storage"
	"github.com/go-chi/chi/v5"
	"net/http"
)

func GetMetric(w http.ResponseWriter, r *http.Request, s storage.Storage) {
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
		metricValue, err := s.GetGauge(n)
		if err != nil {
			http.Error(w, "", http.StatusNotFound)
			return
		}
		v = fmt.Sprintf("%v", metricValue)
	case counter:
		metricValue, err := s.GetCounter(n)
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
