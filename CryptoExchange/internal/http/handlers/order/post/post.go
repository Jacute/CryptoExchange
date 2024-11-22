package orderpost

import (
	mwauth "CryptoExchange/internal/http/middlewares/auth"
	"CryptoExchange/internal/lib/api/response"
	"CryptoExchange/internal/models"
	"CryptoExchange/internal/storage"
	"errors"
	"log/slog"
	"net/http"
	"strconv"

	"github.com/go-chi/chi/v5/middleware"
	"github.com/go-chi/render"
	"github.com/jacute/prettylogger"
)

type OrderSaver interface {
	SaveOrder(userID, pairID int, quantity, price float64, orderType string) (string, error)
}

type UserPayer interface {
	Pay(userID, lotID string, price float64) (string, error)
}

type PairProvider interface {
	GetPairByID(pairID int) (*models.Pair, error)
}

type Request struct {
	PairId   int     `json:"pair_id"`
	Quantity float64 `json:"quantity"`
	Price    float64 `json:"price"`
	Type     string  `json:"type"`
}

type Response struct {
	response.Response
	OrderID int `json:"order_id"`
}

func New(log *slog.Logger, orderSaver OrderSaver, userPayer UserPayer, pairProvider PairProvider) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()
		user := ctx.Value(mwauth.UserContextKey).(*models.User)
		log = log.With(
			slog.String("op", "handlers.order.post.New"),
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

		// validators
		if req.Type != "sell" && req.Type != "buy" {
			render.JSON(w, r, response.Error("invalid order type, should be 'sell' or 'buy'"))
			return
		}
		if req.Price < 0 {
			render.JSON(w, r, response.Error("invalid order price, should be positive"))
			return
		}
		if req.Quantity <= 0 {
			render.JSON(w, r, response.Error("invalid order quantity, should be positive"))
			return
		}

		pair, err := pairProvider.GetPairByID(req.PairId)
		if err != nil {
			if errors.Is(err, storage.ErrPairNotFound) {
				log.Warn("invalid pair id")
				render.JSON(w, r, response.Error("invalid pair ID"))
				return
			}
			log.Error("failed to get pair", prettylogger.Err(err))
			render.JSON(w, r, response.Error("failed to get pair"))
			return
		}

		idStr, err := orderSaver.SaveOrder(user.ID, req.PairId, req.Quantity, req.Price, req.Type)
		if err != nil {
			log.Error("error saving order", prettylogger.Err(err))
			render.JSON(w, r, &Response{
				Response: response.Error("error saving order"),
			})
			return
		}

		id, err := strconv.Atoi(idStr)
		if err != nil {
			log.Error("error converting order id to int", prettylogger.Err(err))
			render.JSON(w, r, &Response{
				Response: response.Error("error saving order"),
			})
			return
		}

		if req.Type == "buy" {
			_, err := userPayer.Pay(strconv.Itoa(user.ID), strconv.Itoa(pair.BuyLotID), req.Price)
			if err != nil {
				if errors.Is(err, storage.ErrNotEnoughMoney) {
					render.JSON(w, r, response.Error("not enough money"))
					return
				}
				log.Error("error paying user", prettylogger.Err(err))
				render.JSON(w, r, response.Error("error paying user"))
				return
			}
		}

		log.Info("order created successfully")
		render.JSON(w, r, &Response{
			Response: response.OK(),
			OrderID:  id,
		})
	}
}
