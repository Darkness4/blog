package middleware

import (
	"net/http"
	"strings"
)

func CSP(policy string) func(http.Handler) http.Handler {
	cleanedPolicy := strings.ReplaceAll(policy, "\n", "")
	cleanedPolicy = strings.TrimSpace(cleanedPolicy)
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			w.Header().Set("Content-Security-Policy", cleanedPolicy)
			next.ServeHTTP(w, r)
		})
	}
}
