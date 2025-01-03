package mwlogger

import (
	"log/slog"
	"net/http"
	"time"

	"github.com/go-chi/chi/v5/middleware"
)

// New creates a logger middleware
func New(log *slog.Logger) func(next http.Handler) http.Handler {
	const op = "middlewares.logger.New"
	return func(next http.Handler) http.Handler {
		log = log.With(
			slog.String("op", op),
		)

		log.Info("Logger middleware enabled")

		fn := func(w http.ResponseWriter, r *http.Request) {
			rid := middleware.GetReqID(r.Context())
			entry := log.With(
				slog.String("method", r.Method),
				slog.String("path", r.URL.Path),
				slog.String("remote_addr", r.RemoteAddr),
				slog.String("user_agent", r.UserAgent()),
				slog.String("request_id", rid),
			)
			ww := middleware.NewWrapResponseWriter(w, r.ProtoMajor)

			t1 := time.Now()
			defer func() {
				entry.Info("Request completed",
					slog.Int("status", ww.Status()),
					slog.Int("bytes", ww.BytesWritten()),
					slog.String("duration", time.Since(t1).String()),
				)
			}()

			next.ServeHTTP(ww, r)
		}

		return http.HandlerFunc(fn)
	}
}
