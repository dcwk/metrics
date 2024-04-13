package handlers

import (
	"net/http"
)

func (h *Handlers) Ping(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json;charset=UTF-8")
	w.Header().Set("Content-Encoding", "gzip")

	_, err := h.DB.Conn(r.Context())
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
}
