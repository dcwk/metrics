package handlers

import (
	"fmt"
	"net/http"

	"github.com/dcwk/metrics/internal/logger"
)

func (h *Handlers) Ping(w http.ResponseWriter, r *http.Request) {
	err := h.Storage.Ping(r.Context())
	if err != nil {
		logger.Log.Info(fmt.Sprintf("can't connect to db: %s", err.Error()))
		http.Error(w, err.Error(), http.StatusInternalServerError)

		return
	}

	w.WriteHeader(http.StatusOK)
}
