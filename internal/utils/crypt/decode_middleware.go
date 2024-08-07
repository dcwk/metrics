package crypt

import (
	"bytes"
	"fmt"
	"io"
	"net/http"
	"strings"

	"github.com/dcwk/metrics/internal/logger"
)

func DecodeBodyMiddleware(privateKeyPath string) func(next http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			var requestCopy bytes.Buffer

			defer r.Body.Close()
			if _, err := io.ReadAll(r.Body); err != nil {
				w.WriteHeader(http.StatusInternalServerError)
				return
			}

			decodedBody, err := DecryptWithRSA(requestCopy.Bytes(), privateKeyPath)
			if err != nil {
				logger.Log.Fatal(fmt.Sprintf("couldn't decrypt request body : %s", requestCopy.Bytes()))
				return
			}

			r.Body = io.NopCloser(strings.NewReader(string(decodedBody)))
			next.ServeHTTP(w, r)
		})
	}
}
