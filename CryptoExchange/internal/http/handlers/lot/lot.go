package lot

import (
	"CryptoExchange/internal/lib/api/response"
	"CryptoExchange/internal/models"
	"log/slog"
	"net/http"

	"github.com/go-chi/chi/middleware"
	"github.com/go-chi/render"
	"github.com/jacute/prettylogger"
)

type LotProvider interface {
	GetLots() ([]*models.Lot, error)
}

func New(log *slog.Logger, lotProvider LotProvider) http.HandlerFunc {
	const op = "handlers.lot.New"
	return func(w http.ResponseWriter, r *http.Request) {
		log = log.With(
			slog.String("op", op),
			slog.String("request_id", middleware.GetReqID(r.Context())),
		)

		lots, err := lotProvider.GetLots()
		if err != nil {
			log.Error("failed to get lots", prettylogger.Err(err))
			render.JSON(w, r, response.Error("failed to get lots"))
			return
		}

		log.Info("lots got successfully")
		render.JSON(w, r, lots)
	}
}
