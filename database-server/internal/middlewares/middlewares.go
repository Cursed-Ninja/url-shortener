package middlewares

import (
	"net/http"
	"url-shortner-database/internal/utils"
)

func LoggingMiddleware(h http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		requestId := r.Header.Get("X-request-id")

		if requestId == "" {
			requestId = utils.GenerateRequestId()
		}

		r.Header.Set("X-request-id", requestId)
		h.ServeHTTP(w, r)
	})
}
