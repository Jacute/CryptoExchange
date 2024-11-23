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
	SaveOrder(order *models.Order) (string, error)
}

type UserPayer interface {
	Pay(userID, lotID string, price float64) (string, error)
	GetOrderForOperation(price float64, pairID, orderType string) (*models.Order, error)
	Buy(buyerOrder *models.Order, sellerOrder *models.Order) error
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

		myOrder := &models.Order{
			UserID:   user.ID,
			PairID:   req.PairId,
			Quantity: req.Quantity,
			Price:    req.Price,
			Type:     req.Type,
			Closed:   "0",
		}

		idStr, err := orderSaver.SaveOrder(myOrder)
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
		myOrder.ID = id

		if req.Type == "buy" {
			_, err := userPayer.Pay(strconv.Itoa(user.ID), strconv.Itoa(pair.SellLotID), req.Quantity*req.Price)
			if err != nil {
				if errors.Is(err, storage.ErrNotEnoughMoney) {
					render.JSON(w, r, response.Error("not enough money"))
					return
				}
				log.Error("error paying user", prettylogger.Err(err))
				render.JSON(w, r, response.Error("error paying user"))
				return
			}

			sellerOrder, err := userPayer.GetOrderForOperation(req.Price, strconv.Itoa(pair.ID), "sell")
			if err != nil {
				if errors.Is(err, storage.ErrOrderNotFound) {
					log.Info("order created successfully, but seller not found yet")
					render.JSON(w, r, &Response{
						Response: response.OK(),
						OrderID:  id,
					})
					return
				}
				log.Error("error checking for buying", prettylogger.Err(err))
				render.JSON(w, r, response.Error("error checking for buying"))
				return
			}
			err = userPayer.Buy(myOrder, sellerOrder)
			if err != nil {
				log.Error("error buying order", prettylogger.Err(err))
				render.JSON(w, r, response.Error("error buying order"))
				return
			}
		} else {
			buyerOrder, err := userPayer.GetOrderForOperation(req.Price, strconv.Itoa(pair.ID), "buy")
			if err != nil {
				if errors.Is(err, storage.ErrOrderNotFound) {
					log.Info("order created successfully, but buyer not found yet")
					render.JSON(w, r, &Response{
						Response: response.OK(),
						OrderID:  id,
					})
					return
				}
				log.Error("error checking for buying", prettylogger.Err(err))
				render.JSON(w, r, response.Error("error checking for buying"))
				return
			}
			err = userPayer.Buy(buyerOrder, myOrder)
			if err != nil {
				log.Error("error buying order", prettylogger.Err(err))
				render.JSON(w, r, response.Error("error buying order"))
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
