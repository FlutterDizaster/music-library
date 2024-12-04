package middleware

import (
	"log/slog"
	"net/http"
	"time"
)

// Logger is a middleware that logs incoming requests
// and logging some of the details to slog.
func Logger(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()

		mw := memoryWriter{
			ResponseWriter: w,
			statusCode:     http.StatusOK,
			size:           0,
		}
		next.ServeHTTP(&mw, r)

		slog.Info(
			"Incoming request",
			slog.Int("status", mw.statusCode),
			slog.String("method", r.Method),
			slog.String("path", r.URL.Path),
			slog.Duration("duration", time.Since(start)),
			slog.String("client ip", getClientIP(r)),
		)
	})
}

// getClientIP helper for getting client IP.
func getClientIP(r *http.Request) string {
	if ip := r.Header.Get("X-Forwarded-For"); ip != "" {
		return ip
	}
	if ip := r.Header.Get("X-Real-Ip"); ip != "" {
		return ip
	}
	return r.RemoteAddr
}
