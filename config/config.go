package config

import (
	"errors"
	"fmt"
	"os"
	"strconv"
)

type PostgresConfig struct {
	User       string
	Password   string
	DbName     string
	Host       string
	Port       int
	ShowConfig bool
}

func newPostgresConfig() (PostgresConfig, error) {
	user := os.Getenv("HANGCOUNTS_POSTGRES_USER")
	pw := os.Getenv("HANGCOUNTS_POSTGRES_PASSWORD")
	db := os.Getenv("HANGCOUNTS_POSTGRES_DB")
	host := os.Getenv("HANGCOUNTS_POSTGRES_HOST")
	portstr := os.Getenv("HANGCOUNTS_POSTGRES_PORT")
	showConfigStr := os.Getenv("HANGCOUNTS_POSTGRES_SHOW_CONFIG")

	// let the condition below catch this failure
	port, _ := strconv.Atoi(portstr)

	if user == "" || pw == "" || db == "" || host == "" || port == 0 {
		return PostgresConfig{}, errors.New("empty postgres config")
	}

	showConfig := false
	if showConfigStr == "true" {
		showConfig = true
	}

	return PostgresConfig{
		User:       user,
		Password:   pw,
		DbName:     db,
		Host:       host,
		Port:       port,
		ShowConfig: showConfig,
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
