package orderdelete

import (
	mwauth "CryptoExchange/internal/http/middlewares/auth"
	"CryptoExchange/internal/lib/api/response"
	"CryptoExchange/internal/models"
	"CryptoExchange/internal/storage"
	"errors"
	"log/slog"
	"net/http"
	"strconv"

	"github.com/go-chi/chi/middleware"
	"github.com/go-chi/render"
	"github.com/jacute/prettylogger"
)

type OrderDeleter interface {
	DeleteOrder(orderID string) error
}

type OrderProvider interface {
	GetOrderByID(orderID string) (*models.Order, error)
}

type UserMoneyAdder interface {
	AddMoney(userID, lotID string, quantity float64) (string, error)
}

type PairProvider interface {
	GetPairByID(pairID int) (*models.Pair, error)
}

type Request struct {
	OrderID int `json:"order_id"`
}

func New(log *slog.Logger, orderDeleter OrderDeleter, orderProvider OrderProvider, moneyAdder UserMoneyAdder, pairProvider PairProvider) http.HandlerFunc {
	const op = "handlers.order.delete.New"
	return func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()
		user := ctx.Value(mwauth.UserContextKey).(*models.User)
		log = log.With(
			slog.String("op", op),
			slog.String("request_id", middleware.GetReqID(ctx)),
			slog.String("username", user.Username),
			slog.Int("user_id", user.ID),
		)

		var req Request
		err := render.DecodeJSON(r.Body, &req)
		if err != nil {
			log.Warn("invalid request", prettylogger.Err(err))
			render.JSON(w, r, response.Error("invalid request"))
			return
		}

		log = log.With(
			slog.Int("order_id", req.OrderID),
			slog.Int("user_id", user.ID),
		)

		order, err := orderProvider.GetOrderByID(strconv.Itoa(req.OrderID))
		if err != nil {
			if errors.Is(err, storage.ErrOrderNotFound) {
				log.Warn("order not found")
				render.JSON(w, r, response.Error("order not found"))
				return
			}
			log.Error("cannot get order", prettylogger.Err(err))
			render.JSON(w, r, response.Error("failed to get order"))
			return
		}

		if order.UserID != user.ID {
			log.Warn("user does not own the order")
			render.JSON(w, r, response.Error("you do not own this order"))
			return
		}

		pair, err := pairProvider.GetPairByID(order.PairID)
		if err != nil {
			log.Error("cannot get pair", prettylogger.Err(err))
			render.JSON(w, r, response.Error("failed to get pair"))
			return
		}

		err = orderDeleter.DeleteOrder(strconv.Itoa(req.OrderID))
		if err != nil {
			log.Error("cannot delete order", prettylogger.Err(err))
			render.JSON(w, r, response.Error("failed to delete order"))
			return
		}

		if order.Type == "buy" {
			_, err := moneyAdder.AddMoney(strconv.Itoa(user.ID), strconv.Itoa(pair.BuyLotID), order.Price)
			if err != nil {
				log.Error("cannot return money to user", prettylogger.Err(err))
				render.JSON(w, r, response.Error("failed to return money to user"))
				return
			}
		}

		log.Info("order deleted successfully")
		render.JSON(w, r, response.OK())
	}
}
