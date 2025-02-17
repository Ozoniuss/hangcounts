package infrastructure

import (
	"context"
	"log/slog"
	"os"
	"testing"

	"github.com/Ozoniuss/hangcounts/config"
	"github.com/Ozoniuss/hangcounts/domain/model"
	"github.com/Ozoniuss/hangcounts/domain/storage"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"
)

type PostgresStoreTestSuite struct {
	suite.Suite
	pgStore *PostgresStore
	logger  *slog.Logger
}

func TestPostgresStoreTestSuite(t *testing.T) {
	if os.Getenv("HANGCOUNTS_RUN_INTEGRATION_TESTS") == "true" {
		suite.Run(t, new(PostgresStoreTestSuite))
	} else {
		t.Skipf("Skipping integration tests, HANGCOUNTS_RUN_INTEGRATION_TESTS is not true")
	}
}

func (suite *PostgresStoreTestSuite) SetupSuite() {
	logopts := &slog.HandlerOptions{
		Level: slog.LevelError,
	}
	handler := slog.NewTextHandler(os.Stdout, logopts)
	logger := slog.New(handler)

	suite.logger = logger
	pg, err := NewPostgresStore(context.TODO(), config.PostgresConfig{
		User:     "test",
		DbName:   "test",
		Password: "test",
		Host:     "localhost",
		Port:     5433,
	}, logger)
	if err != nil {
		suite.FailNow("could not start suite", err.Error())
	}
	suite.pgStore = pg
}

func (suite *PostgresStoreTestSuite) SetupTest() {
	_, err := suite.pgStore.conn.Exec(context.Background(), "DELETE FROM hangout_individuals WHERE 1=1; DELETE FROM hangouts WHERE 1=1; DELETE FROM individuals WHERE 1=1;")
	if err != nil {
		suite.FailNow("could not truncate tables", err.Error())
	}
}

func (suite *PostgresStoreTestSuite) TestGetIndividual_ReturnsCorrespondingError_IfIndividualIsNotFound() {
	_, err := suite.pgStore.GetIndividual(suite.T().Context(), 1)
	assert.ErrorIs(suite.T(), err, storage.ErrNotFound, "expected storage error if individual is not found")
}
func (suite *PostgresStoreTestSuite) TestStoreIndividual_ReturnsIndividualWithCorrectData() {

	individual := model.Individual{
		Id:       1,
		Name:     "name",
		Email:    "email",
		Username: "username",
	}
	err := suite.pgStore.StoreIndividual(suite.T().Context(), individual)
	suite.Require().NoError(err, "expected no error when inserting individual")

	ind, err := suite.pgStore.GetIndividual(suite.T().Context(), 1)
	suite.Require().NoError(err, "expected no error when retrieving individual")

	assert.Equal(suite.T(), individual, ind, "expected stored and retrieved individual to match")
}
