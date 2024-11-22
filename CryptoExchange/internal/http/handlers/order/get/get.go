package orderget

import (
	"CryptoExchange/internal/lib/api/response"
	"CryptoExchange/internal/models"
	"log/slog"
	"net/http"

	"github.com/go-chi/chi/middleware"
	"github.com/go-chi/render"
	"github.com/jacute/prettylogger"
)

type OrderProvider interface {
	GetOrders() ([]*models.Order, error)
}

func New(log *slog.Logger, orderProvider OrderProvider) http.HandlerFunc {
	const op = "handlers.order.get.New"
	return func(w http.ResponseWriter, r *http.Request) {
		log = log.With(
			slog.String("op", op),
			slog.String("request_id", middleware.GetReqID(r.Context())),
		)

		orders, err := orderProvider.GetOrders()
		if err != nil {
			log.Error("failed to get orders", prettylogger.Err(err))
			render.JSON(w, r, response.Error("failed to get orders"))
			return
		}
		log.Info("orders got successfully")
		render.JSON(w, r, orders)
	}
}
