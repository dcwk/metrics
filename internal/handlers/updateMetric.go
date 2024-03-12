package handlers

import (
	"github.com/dcwk/metrics/internal/storage"
	"github.com/go-chi/chi/v5"
	"net/http"
)

func UpdateMetric(w http.ResponseWriter, r *http.Request, s storage.Storage) {
	r.Method = http.MethodPost
	r.Header.Set("Content-Type", "text/plain")

	t := chi.URLParam(r, "type")
	mn := chi.URLParam(r, "name")
	mv := chi.URLParam(r, "value")

	switch t {
	default:
		http.Error(w, "", http.StatusBadRequest)
		return
	case gauge:
		if err := s.AddGauge(mn, mv); err != nil {
			http.Error(w, "", http.StatusBadRequest)
			return
		}
	case counter:
		if err := s.AddCounter(mn, mv); err != nil {
			http.Error(w, "", http.StatusBadRequest)
			return
		}
	}

	w.WriteHeader(http.StatusOK)
}
