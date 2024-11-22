package pair

import (
	"CryptoExchange/internal/lib/api/response"
	"CryptoExchange/internal/models"
	"log/slog"
	"net/http"

	"github.com/go-chi/chi/middleware"
	"github.com/go-chi/render"
	"github.com/jacute/prettylogger"
)

type PairProvider interface {
	GetPairs() ([]*models.Pair, error)
}

func New(log *slog.Logger, pairProvider PairProvider) http.HandlerFunc {
	const op = "handlers.pair.New"
	return func(w http.ResponseWriter, r *http.Request) {
		log = log.With(
			slog.String("op", op),
			slog.String("request_id", middleware.GetReqID(r.Context())),
		)

		pairs, err := pairProvider.GetPairs()
		if err != nil {
			log.Error("failed to get pairs", prettylogger.Err(err))
			render.JSON(w, r, response.Error("failed to get pairs"))
			return
		}

		log.Info("pairs got successfully")
		render.JSON(w, r, pairs)
	}
}
