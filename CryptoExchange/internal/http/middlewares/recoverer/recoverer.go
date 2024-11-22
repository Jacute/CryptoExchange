package mwrecoverer

import (
	"log/slog"
	"net/http"
	"runtime"

	"github.com/jacute/prettylogger"
)

// New returns a recoverer middleware handler, need for recover from panic
func New(log *slog.Logger) func(next http.Handler) http.Handler {
	const op = "middlewares.recoverer.New"
	return func(next http.Handler) http.Handler {
		log = log.With(
			slog.String("op", op),
		)

		log.Info("Recoverer middleware enabled")

		fn := func(w http.ResponseWriter, r *http.Request) {
			defer func() {
				if err := recover(); err != nil {
					buf := make([]byte, 2048)
					n := runtime.Stack(buf, false)
					buf = buf[:n]

					log.Error(
						"Recovering from panic",
						prettylogger.Err(err.(error)),
						slog.String("stacktrace", string(buf)),
					)
				}
			}()
			next.ServeHTTP(w, r)
		}

		return http.HandlerFunc(fn)
	}
}
