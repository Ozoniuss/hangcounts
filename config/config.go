package config

import (
	"errors"
	"fmt"
	"os"
	"strconv"
)

type PostgresConfig struct {
	User     string
	Password string
	DbName   string
}

func newPostgresConfig() (PostgresConfig, error) {
	user := os.Getenv("HANGCOUNTS_POSTGRES_USER")
	pw := os.Getenv("HANGCOUNTS_POSTGRES_PASSWORD")
	db := os.Getenv("HANGCOUNTS_POSTGRES_DB")

	if user == "" || pw == "" || db == "" {
		return PostgresConfig{}, errors.New("empty postgres config")
	}

	return PostgresConfig{
		User:     user,
		Password: pw,
		DbName:   db,
	}, nil
}

type AppConfig struct {
	Env      string
	Database PostgresConfig
}

func NewAppConfig() (AppConfig, error) {
	var cfgErr error
	pgconfig, err := newPostgresConfig()
	if err != nil {
		cfgErr = errors.Join(cfgErr, err)
	}

	env := os.Getenv("HANGCOUNTS_ENV")
	if env != "dev" && env != "prod" {
		cfgErr = errors.Join(cfgErr, fmt.Errorf("invalid env value %s: must be \"dev\" or \"prod\"", strconv.Quote(env)))
	}

	if cfgErr != nil {
		return AppConfig{}, cfgErr
	}
	return AppConfig{
		Database: pgconfig,
		Env:      env,
	}, nil
}
