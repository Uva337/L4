package api

import (
	"net/http"
	"time"

	"calendar/pkg/logger"
)

// оборачивает HTTP-обработчик и замеряет время выполнения
func LoggingMiddleware(next http.HandlerFunc, l *logger.AsyncLogger) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()

		next.ServeHTTP(w, r)

		duration := time.Since(start).String()
		l.Info(r.Method, r.URL.Path, duration)
	}
}
