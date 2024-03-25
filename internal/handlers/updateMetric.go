package handlers

import (
	"net/http"

	"github.com/go-chi/chi/v5"
)

func (h *Handlers) UpdateMetric(w http.ResponseWriter, r *http.Request) {
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
		if err := h.Storage.AddGauge(mn, mv); err != nil {
			http.Error(w, "", http.StatusBadRequest)
			return
		}
	case counter:
		if err := h.Storage.AddCounter(mn, mv); err != nil {
			http.Error(w, "", http.StatusBadRequest)
			return
		}
	}

	w.WriteHeader(http.StatusOK)
}
