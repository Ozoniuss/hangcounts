package main

import (
	"context"
	"fmt"
	"log/slog"
	"os"

	"github.com/Ozoniuss/hangcounts/config"
	"github.com/Ozoniuss/hangcounts/infrastructure"
)

func run() error {
	config, err := config.NewAppConfig()
	if err != nil {
		return fmt.Errorf("could not read config: %w", err)
	}

	logopts := &slog.HandlerOptions{
		Level: slog.LevelInfo,
	}
	if config.Env == "dev" {
		logopts.Level = slog.LevelDebug
	}
	handler := slog.NewJSONHandler(os.Stdout, logopts)
	logger := slog.New(handler)

	logger.Info("read application mode", slog.String("env", config.Env))
	if config.Database.ShowConfig {
		logger.Debug("database config", slog.Any("config", config.Database))
	}

	ctx := context.Background()

	pgStore, err := infrastructure.NewPostgresStore(ctx, config.Database, logger)
	if err != nil {
		return fmt.Errorf("could not create a postgres store: %w", err)
	}
	logger.Info("connected to postgres database", slog.String("host", config.Database.Host), slog.Int("port", config.Database.Port))

	logger.Info("starting app")

	fmt.Println(config, pgStore)
	return nil
}

func main() {
	if err := run(); err != nil {
		slog.Error("could not start app", slog.Any("error", err))
		os.Exit(1)
	}
}
