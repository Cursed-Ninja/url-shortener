package middlewares

import (
	"main-server/internal/utils"
	"net/http"
)

func LoggingMiddleware(h http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		r.Header.Set("X-request-id", utils.GenerateRequestId())
		h.ServeHTTP(w, r)
	})
}
