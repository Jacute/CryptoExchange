package mwauth

import (
	"CryptoExchange/internal/lib/api/response"
	"CryptoExchange/internal/models"
	"CryptoExchange/internal/storage"
	"context"
	"errors"
	"log/slog"
	"net/http"

	"github.com/go-chi/render"
)

type ContextKey string

const UserContextKey ContextKey = "username"

type UserProvider interface {
	GetUserByToken(token string) (*models.User, error)
}

// New creates auth checker middleware
func New(log *slog.Logger, userProvider UserProvider) func(next http.Handler) http.Handler {
	const op = "middlewares.auth.New"
	return func(next http.Handler) http.Handler {
		log = log.With(
			slog.String("op", op),
		)

		fn := func(w http.ResponseWriter, r *http.Request) {
			token := r.Header.Get("X-USER-TOKEN")
			if token == "" {
				render.JSON(w, r, response.Error("header X-USER-TOKEN is missing"))
				return
			}
			user, err := userProvider.GetUserByToken(token)
			if err != nil {
				if errors.Is(err, storage.ErrUserNotFound) {
					render.JSON(w, r, response.Error("invalid token"))
					return
				}
				render.JSON(w, r, response.Error("internal error"))
				return
			}
			log.Info("user is authenticated", slog.String("token", token))

			ctx := context.WithValue(r.Context(), UserContextKey, user)
			r = r.WithContext(ctx)

			next.ServeHTTP(w, r)
		}
		return http.HandlerFunc(fn)
	}
}
