package handlers

import (
	"fmt"
	"net/http"

	"github.com/dcwk/metrics/internal/logger"
)

func (h *Handlers) Ping(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json;charset=UTF-8")
	w.Header().Set("Content-Encoding", "gzip")

	err := h.DB.PingContext(r.Context())
	if err != nil {
		logger.Log.Info(fmt.Sprintf("can't connect to db: %s", err.Error()))
		http.Error(w, err.Error(), http.StatusInternalServerError)

		return
	}

	w.WriteHeader(http.StatusOK)
}
