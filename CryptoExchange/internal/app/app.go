package app

import (
	"JacuteCE/internal/config"
	"JacuteCE/internal/http/handlers/user"
	mwlogger "JacuteCE/internal/http/middlewares/logger"
	mwrecoverer "JacuteCE/internal/http/middlewares/recoverer"
	jacutesql "JacuteCE/internal/storage/JacuteSQL"
	"fmt"
	"log/slog"
	"net/http"

	"github.com/go-chi/chi"
	"github.com/go-chi/chi/v5/middleware"
)

type App struct {
	cfg        *config.AppConfig
	log        *slog.Logger
	db         *jacutesql.Storage
	httpServer *http.Server
}

func New(cfg *config.AppConfig, log *slog.Logger, db *jacutesql.Storage) *App {
	return &App{
		cfg: cfg,
		log: log,
		db:  db,
	}
}

func (a *App) setupRouter() http.Handler {
	router := chi.NewRouter()
	loggerMiddleware := mwlogger.New(a.log)
	recovererMiddleware := mwrecoverer.New(a.log)

	router.Use(recovererMiddleware)
	router.Use(middleware.RequestID)
	// router.Use(middleware.URLFormat)
	router.Use(loggerMiddleware)

	router.Post("/user", user.New(a.log, a.db, a.cfg.TokenLen))

	return router
}

func (a *App) MustRun() {
	if err := a.Run(); err != nil {
		panic(err)
	}
}

func (a *App) Run() error {
	const op = "app.Run"

	a.httpServer = &http.Server{
		Addr:         fmt.Sprintf("%s:%d", a.cfg.IP, a.cfg.Port),
		Handler:      a.setupRouter(),
		ReadTimeout:  a.cfg.ReadTimeout,
		WriteTimeout: a.cfg.WriteTimeout,
		IdleTimeout:  a.cfg.IdleTimeout,
	}

	a.log.Info("http server listening", slog.String("ip", a.cfg.IP), slog.Int("port", a.cfg.Port))

	if err := a.httpServer.ListenAndServe(); err != nil {
		return fmt.Errorf("%s: %w", op, err)
	}

	return nil
}

func (a *App) Stop() {
	err := a.httpServer.Close()
	if err != nil {
		panic(err)
	}
}
