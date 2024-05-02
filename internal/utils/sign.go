package utils

import (
	"fmt"
	"net/http"

	"github.com/dcwk/metrics/internal/logger"
)

func SignMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		sign := r.Header.Get("HashSHA256")
		if sign == "" {
			next.ServeHTTP(w, r)
			return
		}

		logger.Log.Info(fmt.Sprintf("Sign data: %s", sign))
	})
}
