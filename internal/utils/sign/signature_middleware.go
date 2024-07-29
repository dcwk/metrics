package sign

import (
	"bytes"
	"crypto/hmac"
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"io"
	"net/http"

	"github.com/dcwk/metrics/internal/logger"
)

func SignMiddleware(signKey string) func(next http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			var requestCopy bytes.Buffer
			var currentSign []byte

			sign := r.Header.Get("HashSHA256")
			if sign == "" {
				next.ServeHTTP(w, r)
				return
			}

			logger.Log.Info(fmt.Sprintf("Sign data: %s", sign))
			defer r.Body.Close()
			if _, err := requestCopy.ReadFrom(r.Body); err != nil {
				w.WriteHeader(http.StatusInternalServerError)
				return
			}

			h := hmac.New(sha256.New, []byte(signKey))
			currentSign = h.Sum(requestCopy.Bytes())
			if hex.EncodeToString(currentSign) != sign {
				w.WriteHeader(http.StatusBadRequest)
				return
			}

			r.Body = io.NopCloser(&requestCopy)
			next.ServeHTTP(w, r)
		})
	}
}
