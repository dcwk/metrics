package realip

import (
	"fmt"
	"net"
	"net/http"

	"github.com/dcwk/metrics/internal/logger"
)

// CheckXRealIpMiddleware checks the ip in the X-Real-IP http header
// whether it is part of a trusted network
func CheckXRealIpMiddleware(trustedNetwork string) func(next http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			if trustedNetwork == "" {
				next.ServeHTTP(w, r)
				return
			}

			xRealIP := r.Header.Get("X-Real-IP")
			logger.Log.Info(fmt.Sprintf("client made a request with the x-real-ip header %s", xRealIP))
			clientIP := net.ParseIP(xRealIP)
			if clientIP == nil {
				w.WriteHeader(http.StatusForbidden)
				return
			}

			_, network, err := net.ParseCIDR(trustedNetwork)
			if err != nil {
				logger.Log.Fatal(fmt.Sprintf("trusted network variable couldn't be parsed %s", trustedNetwork))
				w.WriteHeader(http.StatusInternalServerError)
				return
			}

			if !network.Contains(clientIP) {
				w.WriteHeader(http.StatusForbidden)
				return
			}

			next.ServeHTTP(w, r)
		})
	}
}
