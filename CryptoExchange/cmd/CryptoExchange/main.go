package main

import (
	"CryptoExchange/internal/app"
	"CryptoExchange/internal/config"
	"CryptoExchange/internal/logger"
	jacutesql "CryptoExchange/internal/storage/JacuteSQL"
	"log/slog"
	"os"
	"os/signal"
	"syscall"
)

func main() {
	cfg := config.MustLoad()

	logger := logger.New(cfg.Env)
	db := jacutesql.New(cfg.DatabaseConfig.IP, cfg.DatabaseConfig.Port, cfg.Lots)
	db.MakeMigrations(cfg.Lots)

	logger.Log.Info(
		"starting app",
		slog.Any("config", cfg),
	)
	application := app.New(&cfg.AppConfig, logger.Log, db)
	go application.Run()

	stop := make(chan os.Signal, 1)
	signal.Notify(stop, syscall.SIGTERM, syscall.SIGINT)
	sign := <-stop

	application.Stop()

	logger.Log.Info("app stopped", slog.String("signal", sign.String()))
}
