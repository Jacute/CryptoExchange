package balance

import (
	mwauth "CryptoExchange/internal/http/middlewares/auth"
	"CryptoExchange/internal/lib/api/response"
	"CryptoExchange/internal/models"
	"log/slog"
	"net/http"
	"strconv"

	"github.com/go-chi/chi/middleware"
	"github.com/go-chi/render"
	"github.com/jacute/prettylogger"
)

type BalanceProvider interface {
	GetUserLots(userID string) ([]*models.UserLot, error)
}

func New(log *slog.Logger, balanceProvider BalanceProvider) http.HandlerFunc {
	const op = "handlers.lot.New"
	return func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()
		user := ctx.Value(mwauth.UserContextKey).(*models.User)
		log = log.With(
			slog.String("op", op),
			slog.String("request_id", middleware.GetReqID(ctx)),
			slog.String("username", user.Username),
			slog.Int("user_id", user.ID),
		)

		lots, err := balanceProvider.GetUserLots(strconv.Itoa(user.ID))
		if err != nil {
			log.Error("failed to get lots", prettylogger.Err(err))
			render.JSON(w, r, response.Error("failed to get lots"))
			return
		}

		log.Info("user lots got successfully")
		render.JSON(w, r, lots)
	}
}
