package middlewares

import (
	"log/slog"
	"net/http"
)

func LoggingMiddleware(next http.Handler, log *slog.Logger) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		log.Info(
			"Incoming request",
			slog.String("method", r.Method),
			slog.String("url", r.URL.Path),
		)
		next.ServeHTTP(w, r)
	})
}
