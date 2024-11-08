package suite

import (
	"JacuteCE/internal/app"
	"JacuteCE/internal/config"
	jacutesql "JacuteCE/internal/storage/JacuteSQL"
	"log/slog"
	"time"

	"github.com/jacute/prettylogger"
)

type Suite struct {
	App *app.App
	Cfg *config.AppConfig
	DB  *jacutesql.Storage
}

func New() *Suite {
	lots := []string{"RUB", "BTC", "ETH", "USDT", "USDC"}
	cfg := &config.AppConfig{
		IP:           "127.0.0.1",
		Port:         8080,
		ReadTimeout:  30 * time.Second,
		WriteTimeout: 30 * time.Second,
		IdleTimeout:  120 * time.Second,
		TokenLen:     32,
		Lots:         lots,
	}
	db := jacutesql.New("127.0.0.1", 7432, lots)
	db.Destroy()
	db.MakeMigrations(lots)

	application := app.New(cfg, slog.New(prettylogger.NewDiscardHandler()), db)

	return &Suite{
		App: application,
		Cfg: cfg,
		DB:  db,
	}
}
