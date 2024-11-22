package user

import (
	"CryptoExchange/internal/lib/api/response"
	"CryptoExchange/internal/lib/utils"
	"CryptoExchange/internal/storage"
	"errors"
	"log/slog"
	"net/http"

	"github.com/go-chi/chi/v5/middleware"
	"github.com/go-chi/render"
	"github.com/jacute/prettylogger"
)

type Request struct {
	Username string `json:"username"`
}

type Response struct {
	response.Response
	ID    string `json:"id"`
	Token string `json:"token"`
}

type UserProvider interface {
	SaveUser(username string, token string) (string, error)
	AddLots(userID string, quantity string) error
}

func New(log *slog.Logger, userProvider UserProvider, tokenLen int) http.HandlerFunc {
	const op = "handlers.user.New"
	return func(w http.ResponseWriter, r *http.Request) {
		log = log.With(
			slog.String("op", op),
			slog.String("request_id", middleware.GetReqID(r.Context())),
		)

		var req Request
		err := render.DecodeJSON(r.Body, &req)
		if err != nil {
			log.Warn("invalid request", prettylogger.Err(err))
			render.JSON(w, r, response.Error("invalid request"))
			return
		}
		log = log.With(slog.String("username", req.Username))

		token := utils.GenerateToken(tokenLen)

		id, err := userProvider.SaveUser(req.Username, token)
		if err != nil {
			if errors.Is(err, storage.ErrMaliciousParameter) {
				log.Warn("malicious parameter", prettylogger.Err(err))
				render.JSON(w, r, response.Error("malicious parameter"))
				return
			}
			if errors.Is(err, storage.ErrUserExists) {
				log.Warn("user already exists", prettylogger.Err(err))
				render.JSON(w, r, response.Error("user already exists"))
				return
			}
			log.Error("failed to save user", prettylogger.Err(err))
			render.JSON(w, r, response.Error("error saving user"))
			return
		}

		err = userProvider.AddLots(id, "1000")
		if err != nil {
			log.Error("failed to save lot", prettylogger.Err(err))
			render.JSON(w, r, response.Error("error creating lots for user"))
			return
		}

		log.Info("user created")
		render.JSON(w, r, Response{
			Response: response.OK(),
			ID:       id,
			Token:    token,
		})
	}
}
