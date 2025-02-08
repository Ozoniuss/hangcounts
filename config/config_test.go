package config

import (
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
)

func Test_NewPostgresConfig_ValidConfig_ReturnsNoError(t *testing.T) {
	t.Setenv("HANGCOUNTS_POSTGRES_USER", "val")
	t.Setenv("HANGCOUNTS_POSTGRES_PASSWORD", "val")
	t.Setenv("HANGCOUNTS_POSTGRES_DB", "val")
	t.Setenv("HANGCOUNTS_POSTGRES_HOST", "val")
	t.Setenv("HANGCOUNTS_POSTGRES_PORT", "5432")
	t.Setenv("HANGCOUNTS_POSTGRES_SHOW_CONFIG", "true")

	c, err := newPostgresConfig()
	assert.NoError(t, err, "config should be valid")

	assert.Equal(t, c, PostgresConfig{
		User:       "val",
		Password:   "val",
		DbName:     "val",
		Host:       "val",
		Port:       5432,
		ShowConfig: true,
	})
}

func Test_NewPostgresConfig_MissingValue_ReturnsError(t *testing.T) {
	t.Setenv("HANGCOUNTS_POSTGRES_USER", "val")
	t.Setenv("HANGCOUNTS_POSTGRES_PASSWORD", "val")
	t.Setenv("HANGCOUNTS_POSTGRES_DB", "val")
	t.Setenv("HANGCOUNTS_POSTGRES_HOST", "val")
	t.Setenv("HANGCOUNTS_POSTGRES_PORT", "5432")
	t.Setenv("HANGCOUNTS_POSTGRES_SHOW_CONFIG", "true")

	tc := []struct {
		name       string
		envToUnset string
	}{
		{name: "user not set", envToUnset: "HANGCOUNTS_POSTGRES_USER"},
		{name: "password not set", envToUnset: "HANGCOUNTS_POSTGRES_PASSWORD"},
		{name: "db not set", envToUnset: "HANGCOUNTS_POSTGRES_DB"},
		{name: "host not set", envToUnset: "HANGCOUNTS_POSTGRES_HOST"},
		{name: "port not set", envToUnset: "HANGCOUNTS_POSTGRES_PORT"},
	}

	for _, tt := range tc {
		t.Run(tt.name, func(t *testing.T) {
			original := os.Getenv(tt.envToUnset)
			t.Setenv(tt.envToUnset, "")

			_, err := newPostgresConfig()
			assert.Error(t, err, "config should not be valid")

			t.Setenv(tt.envToUnset, original)
		})
	}
}
