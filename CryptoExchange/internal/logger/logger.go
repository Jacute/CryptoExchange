package logger

import (
	"log/slog"
	"os"

	"github.com/jacute/prettylogger"
)

const (
	envLocal = "local"
	envProd  = "prod"
	logPath  = "JacuteCE.log"
)

type Logger struct {
	Log    *slog.Logger
	Writer *os.File
}

func New(env string) *Logger {
	log := &Logger{}

	switch env {
	case envLocal:
		log.Log = slog.New(
			prettylogger.NewColoredHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelDebug}),
		)
	case envProd:
		file, err := os.OpenFile(logPath, os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0644)
		if err != nil {
			panic("Can't open log file: " + err.Error())
		}
		log.Log = slog.New(
			prettylogger.NewJsonHandler(file, &slog.HandlerOptions{Level: slog.LevelInfo}),
		)
	}
	return log
}
