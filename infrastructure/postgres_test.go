package infrastructure

import (
	"context"
	"fmt"
	"log/slog"
	"os"
	"testing"

	"github.com/Ozoniuss/hangcounts/config"
	"github.com/stretchr/testify/suite"
)

type PostgresStoreTestSuite struct {
	suite.Suite
	pgStore *PostgresStore
	logger  *slog.Logger
}

// In order for 'go test' to run this suite, we need to create
// a normal test function and pass our suite to suite.Run
func TestPostgresStoreTestSuite(t *testing.T) {
	fmt.Println("a")
	if os.Getenv("HANGCOUNTS_RUN_INTEGRATION_TESTS") == "true" {
		suite.Run(t, new(PostgresStoreTestSuite))
	} else {
		t.Skipf("Skipping integration tests, HANGCOUNTS_RUN_INTEGRATION_TESTS is not true")
	}
}

func (s *PostgresStoreTestSuite) SetupSuite() {
	logopts := &slog.HandlerOptions{
		Level: slog.LevelError,
	}
	handler := slog.NewTextHandler(os.Stdout, logopts)
	logger := slog.New(handler)

	s.logger = logger
	pg, err := NewPostgresStore(context.TODO(), config.PostgresConfig{
		User:     "test",
		DbName:   "test",
		Password: "test",
		Host:     "localhost",
		Port:     5433,
	}, logger)
	if err != nil {
		s.FailNowf("could not start suite", "reason: %s", err.Error())
	}
	s.pgStore = pg
}

func (s *PostgresStoreTestSuite) TestDummy() {
	s.Assert().Equal(1, 1)
}
