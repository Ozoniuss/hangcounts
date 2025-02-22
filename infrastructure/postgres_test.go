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
		Level: slog.LevelDebug,
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
	_, err := suite.pgStore.GetIndividual(suite.T().Context(), "username")
	assert.ErrorIs(suite.T(), err, storage.ErrNotFound, "expected storage error if individual is not found")
}

func (suite *PostgresStoreTestSuite) TestGetIndividual_ReturnsError_IfIndividualExists() {
	// same as TestDeleteIndividual_ReturnsNoError_IfIndividualExists
}

func (suite *PostgresStoreTestSuite) TestStoreIndividual_ReturnsIndividualWithCorrectData() {

	individual := model.Individual{
		Username: "username",
		Name:     "name",
		Email:    "email",
	}
	err := suite.pgStore.StoreIndividual(suite.T().Context(), individual)
	suite.Require().NoError(err, "expected no error when inserting individual")

	ind, err := suite.pgStore.GetIndividual(suite.T().Context(), "username")
	suite.Require().NoError(err, "expected no error when retrieving individual")

	assert.Equal(suite.T(), individual, ind, "expected stored and retrieved individual to match")
}

func (suite *PostgresStoreTestSuite) TestStoreIndividual_ReturnsError_IfUsernameConstraintIsViolated() {

	individual := model.Individual{
		Username: "username",
		Name:     "name",
		Email:    "email",
	}
	err := suite.pgStore.StoreIndividual(suite.T().Context(), individual)
	suite.Require().NoError(err, "expected no error when inserting individual")

	individual.Email = "other"
	err = suite.pgStore.StoreIndividual(suite.T().Context(), individual)
	assert.ErrorIs(suite.T(), err, storage.ErrUsernameAlreadyExists, "expected error if username is taken")
}

func (suite *PostgresStoreTestSuite) TestStoreIndividual_ReturnsError_IfEmailConstraintIsViolated() {

	individual := model.Individual{
		Username: "username",
		Name:     "name",
		Email:    "email",
	}
	err := suite.pgStore.StoreIndividual(suite.T().Context(), individual)
	suite.Require().NoError(err, "expected no error when inserting individual")

	individual.Username = "other"
	err = suite.pgStore.StoreIndividual(suite.T().Context(), individual)
	assert.ErrorIs(suite.T(), err, storage.ErrEmailAlreadyExists, "expected error if email is taken")
}

func (suite *PostgresStoreTestSuite) TestDeleteIndividual_ReturnsNoError_IfIndividualExists() {
	individual := model.Individual{
		Username: "username",
		Name:     "name",
		Email:    "email",
	}
	err := suite.pgStore.StoreIndividual(suite.T().Context(), individual)
	suite.Require().NoError(err, "expected no error when inserting individual")

	err = suite.pgStore.MarkIndividualAsDeleted(suite.T().Context(), "username")
	suite.Require().NoError(err, "expected no error when deleting individual")

	_, err = suite.pgStore.GetIndividual(suite.T().Context(), "username")
	suite.Require().ErrorIs(err, storage.ErrDeleted, "expected deleted error if individual is soft-deleted")
}

func (suite *PostgresStoreTestSuite) TestDeleteIndividual_ReturnsError_IfAlreadyDeleted() {

	individual := model.Individual{
		Username: "username",
		Name:     "name",
		Email:    "email",
	}
	err := suite.pgStore.StoreIndividual(suite.T().Context(), individual)
	suite.Require().NoError(err, "expected no error when inserting individual")

	err = suite.pgStore.MarkIndividualAsDeleted(suite.T().Context(), "username")
	suite.Require().NoError(err, "expected no error when deleting individual")

	err = suite.pgStore.MarkIndividualAsDeleted(suite.T().Context(), "username")
	suite.Require().ErrorIs(err, storage.ErrDeleted, "expected deleted error if individual is soft-deleted")
}

func (suite *PostgresStoreTestSuite) TestDeleteIndividual_ReturnsError_IfIndividualDoesntExist() {
	err := suite.pgStore.MarkIndividualAsDeleted(suite.T().Context(), "username")
	suite.Require().ErrorIs(err, storage.ErrNotFound, "expected error when individual is not in the database")
}
