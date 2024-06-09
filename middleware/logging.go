package middleware

import (
	"log/slog"
	"net/http"
	"runtime/debug"
	"time"
)

func LoggingHandler(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()
		attrs := []any{
			slog.Group(
				"request",
				slog.String("method", r.Method),
				slog.String("url", r.URL.String()),
				slog.String("host", r.Host),
				slog.String("proto", r.Proto),
				slog.String("user_agent", r.UserAgent()),
			),
		}

		defer func() {
			if err := recover(); err != nil {
				w.WriteHeader(http.StatusInternalServerError)
				attrs = append(
					attrs,
					slog.Group("response",
						slog.Int("status_code", http.StatusInternalServerError),
						slog.Duration("duration", time.Since(start)),
					),
					slog.Group("panic",
						slog.Any("message", err),
						slog.String("stack", string(debug.Stack())),
					),
				)
				slog.ErrorContext(r.Context(), "server recovered from panic", attrs...)
			}
		}()

		flusher := w.(http.Flusher)
		hijacker := w.(http.Hijacker)
		mw := &metricsHttpWriter{
			ResponseWriter: w,
			Flusher:        flusher,
			Hijacker:       hijacker,
			statusCode:     http.StatusOK,
		}

		next.ServeHTTP(mw, r)

		duration := time.Since(start)
		attrs = append(
			attrs,
			slog.Group(
				"response",
				slog.Int("status_code", mw.statusCode),
				slog.Duration("duration", duration),
			),
		)

		if mw.statusCode >= 500 {
			slog.ErrorContext(r.Context(), "request resulted with server error", attrs...)
		} else if mw.statusCode >= 400 {
			slog.WarnContext(r.Context(), "request resulted with client error", attrs...)
		} else {
			slog.InfoContext(r.Context(), "request finished successfully", attrs...)
		}
	})
}

type metricsHttpWriter struct {
	http.ResponseWriter
	http.Flusher
	http.Hijacker
	statusCode int
}

func (m *metricsHttpWriter) WriteHeader(statusCode int) {
	m.statusCode = statusCode
	m.ResponseWriter.WriteHeader(statusCode)
}
