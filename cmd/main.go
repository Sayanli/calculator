package main

import (
	"flag"
	"log/slog"
	"os"

	"github.com/sayanli/calculator/internal/app"
	"github.com/sayanli/calculator/internal/config"
)

const (
	envDev  = "dev"
	envProd = "prod"

	configPath = "./config/config.yaml"
)

func main() {
	golimit := flag.Int("golimit", 100, "maximum number of working goroutines")
	cfg, err := config.NewConfig(configPath)
	if err != nil {
		panic(err)
	}
	log := setupLogger(cfg.Log.Env)
	application := app.NewApp(log, cfg.GRPC.Port, cfg.HTTP.Port, *golimit)
	application.Run()
}

func setupLogger(env string) *slog.Logger {
	var log *slog.Logger

	switch env {
	case envDev:
		log = slog.New(
			slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelDebug}),
		)
	case envProd:
		log = slog.New(
			slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelInfo}),
		)
	}

	return log
}
