package main

import (
	"log/slog"
	"os"
	"sso/internal/config"
)

const (
	envLocal = "local"
	envProd  = "prod"
)

func main() {
	cfg := config.MustLoadConfig()

	log := setUpLogger(cfg.Env)

	log.Info("Starting app", slog.String("env", cfg.Env))

	// TODO: init app

	// TODO: run gRPC-server
}

func setUpLogger(env string) *slog.Logger {
	var log *slog.Logger

	switch env {
	case envLocal:
		log = slog.New(
			slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelDebug}),
		)
	case envProd:
		log = slog.New(
			slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelInfo}),
		)
	}

	return log
}
