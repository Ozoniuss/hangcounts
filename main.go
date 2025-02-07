package main

import (
	"fmt"
	"log/slog"
	"os"

	"github.com/Ozoniuss/hangcounts/config"
	"github.com/Ozoniuss/hangcounts/domain"
)

type HangoutsService struct {
}

func (h *HangoutsService) NewHangout(Atendees []domain.Individual) (domain.Hangout, error) {
	return domain.Hangout{}, nil
}

func (h *HangoutsService) DeleteHangout(hangoutId uint64) (domain.Hangout, error) {
	return domain.Hangout{}, nil
}

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
	handler := slog.NewJSONHandler(os.Stdout, nil)
	logger := slog.New(handler)

	logger.Info("read application mode", slog.String("env", config.Env))
	logger.Info("starting app")

	fmt.Println(config)
	return nil
}

func main() {
	if err := run(); err != nil {
		slog.Error("could not start app", slog.Any("error", err))
		os.Exit(1)
	}
}
