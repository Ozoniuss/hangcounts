package infrastructure

import (
	"context"
	"fmt"
	"log/slog"
	"os"
	"strconv"
	"testing"
	"time"

	"github.com/Ozoniuss/hangcounts/config"
	"github.com/Ozoniuss/hangcounts/domain/model"
	"github.com/Ozoniuss/hangcounts/domain/storage"
	"github.com/google/uuid"
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
	assert.ErrorIs(suite.T(), err, storage.ErrIndividualUsernameAlreadyExists, "expected error if username is taken")
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
	assert.ErrorIs(suite.T(), err, storage.ErrIndividualEmailAlreadyExists, "expected error if email is taken")
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

func (suite *PostgresStoreTestSuite) TestStoreHangout_ReturnsError_WhenHangoutCreatorDoesntExist() {

	currentDate := time.Now()
	hangout := model.Hangout{
		PublicId: model.HangoutId(uuid.New()),
		HangoutDetails: model.HangoutDetails{
			Description: func() *string { s := "description"; return &s }(),
			Location:    "location",
			Duration:    10,
			Date:        currentDate,
		},
		CreatedBy: model.IndividualId("emil"),
	}
	err := suite.pgStore.StoreHangoutOfIndividuals(suite.T().Context(), hangout)
	suite.Require().ErrorIs(err, storage.ErrHangoutCreatorNotFound, "expected error when creator doesn't exist")
}

func (suite *PostgresStoreTestSuite) TestStoreHangout_ReturnsError_WhenHangoutCreatorIsDeleted() {

	individual := model.Individual{
		Username: "creator",
		Name:     "name",
		Email:    "email",
	}

	err := suite.pgStore.StoreIndividual(suite.T().Context(), individual)
	suite.Require().NoError(err, "expected no error when storing hangout creator")

	err = suite.pgStore.MarkIndividualAsDeleted(suite.T().Context(), "creator")
	suite.Require().NoError(err, "expected no error when soft deleting hangout creator")

	currentDate := time.Now()
	hangout := model.Hangout{
		PublicId: model.HangoutId(uuid.New()),
		HangoutDetails: model.HangoutDetails{
			Description: func() *string { s := "description"; return &s }(),
			Location:    "location",
			Duration:    10,
			Date:        currentDate,
		},
		CreatedBy: model.IndividualId("emil"),
	}
	err = suite.pgStore.StoreHangoutOfIndividuals(suite.T().Context(), hangout)
	suite.Require().Error(err, "expected error when creator was soft deleted")
}

func (suite *PostgresStoreTestSuite) TestStoreHangout_ReturnsNoError_WhenCreatorIsValidAndThereAreNoParticipants() {

	individual := model.Individual{
		Username: "creator",
		Name:     "name",
		Email:    "email",
	}
	err := suite.pgStore.StoreIndividual(suite.T().Context(), individual)
	suite.Require().NoError(err, "expected no error when storing hangout creator")

	currentDate := time.Now()
	hangout := model.Hangout{
		PublicId: model.HangoutId(uuid.New()),
		HangoutDetails: model.HangoutDetails{
			Description: func() *string { s := "description"; return &s }(),
			Location:    "location",
			Duration:    10,
			Date:        currentDate,
		},
		CreatedBy: model.IndividualId(individual.Username),
	}
	err = suite.pgStore.StoreHangoutOfIndividuals(suite.T().Context(), hangout)
	suite.Require().NoError(err, "expected no error when creating a valid hangout")
}

func (suite *PostgresStoreTestSuite) TestStoreHangout_ReturnsNoError_WhenCreatorAndParticipantsExist() {
	individualIds := make([]model.IndividualId, 0)
	individual := model.Individual{
		Username: "creator",
		Name:     "name",
		Email:    "email",
	}
	err := suite.pgStore.StoreIndividual(suite.T().Context(), individual)
	suite.Require().NoError(err, "expected no error when storing hangout creator")
	individualIds = append(individualIds, "creator")

	participants := make([]model.Individual, 0, 10)
	for i := range 10 {
		participants = append(participants, model.Individual{
			Name:     strconv.Itoa(i),
			Username: model.IndividualId(strconv.Itoa(i)),
			Email:    model.Email(strconv.Itoa(i)),
		})
	}
	for i, p := range participants {
		err := suite.pgStore.StoreIndividual(suite.T().Context(), p)
		suite.Require().NoError(err, fmt.Sprintf("expected no error when storing participant %d", i))
		individualIds = append(individualIds, p.Username)
	}

	currentDate := time.Now()
	hangout := model.Hangout{
		PublicId: model.HangoutId(uuid.New()),
		HangoutDetails: model.HangoutDetails{
			Description: func() *string { s := "description"; return &s }(),
			Location:    "location",
			Duration:    10,
			Date:        currentDate,
		},
		CreatedBy:   model.IndividualId(individual.Username),
		Individuals: individualIds,
	}
	err = suite.pgStore.StoreHangoutOfIndividuals(suite.T().Context(), hangout)
	suite.Require().NoError(err, "expected no error when creating a valid hangout")
}

func (suite *PostgresStoreTestSuite) TestStoreHangout_ReturnsError_WhenAParticipantDoesntExist() {
	individualIds := make([]model.IndividualId, 0)
	creator := model.Individual{
		Username: "creator",
		Name:     "name",
		Email:    "email",
	}
	err := suite.pgStore.StoreIndividual(suite.T().Context(), creator)
	suite.Require().NoError(err, "expected no error when storing hangout creator")
	individualIds = append(individualIds, "creator")

	participants := make([]model.Individual, 0, 10)
	for i := range 10 {
		participants = append(participants, model.Individual{
			Name:     strconv.Itoa(i),
			Username: model.IndividualId(strconv.Itoa(i)),
			Email:    model.Email(strconv.Itoa(i)),
		})
	}

	for i, p := range participants {
		// do not store first participant
		if i != 0 {
			err := suite.pgStore.StoreIndividual(suite.T().Context(), p)
			suite.Require().NoError(err, fmt.Sprintf("expected no error when storing participant %d", i))
		}
		// add the first participant as well
		individualIds = append(individualIds, p.Username)
	}

	currentDate := time.Now()
	hangout := model.Hangout{
		PublicId: model.HangoutId(uuid.New()),
		HangoutDetails: model.HangoutDetails{
			Description: func() *string { s := "description"; return &s }(),
			Location:    "location",
			Duration:    10,
			Date:        currentDate,
		},
		CreatedBy:   model.IndividualId(creator.Username),
		Individuals: individualIds,
	}
	err = suite.pgStore.StoreHangoutOfIndividuals(suite.T().Context(), hangout)
	suite.Require().ErrorIs(err, storage.ErrHangoutParticipantNotFound, "expected error when inserting a hangout where a participant is not in the database")
}

func (suite *PostgresStoreTestSuite) TestStoreHangout_ReturnsError_WhenAParticipantIsDeleted() {
	individualIds := make([]model.IndividualId, 0)
	individual := model.Individual{
		Username: "creator",
		Name:     "name",
		Email:    "email",
	}
	err := suite.pgStore.StoreIndividual(suite.T().Context(), individual)
	suite.Require().NoError(err, "expected no error when storing hangout creator")
	individualIds = append(individualIds, "creator")

	participants := make([]model.Individual, 0, 10)
	for i := range 10 {
		participants = append(participants, model.Individual{
			Name:     strconv.Itoa(i),
			Username: model.IndividualId(strconv.Itoa(i)),
			Email:    model.Email(strconv.Itoa(i)),
		})
	}
	for i, p := range participants {
		err := suite.pgStore.StoreIndividual(suite.T().Context(), p)
		suite.Require().NoError(err, fmt.Sprintf("expected no error when storing participant %d", i))
		individualIds = append(individualIds, p.Username)
	}
	// Fourth participant exists in the database but is soft deleted
	suite.pgStore.MarkIndividualAsDeleted(suite.T().Context(), participants[3].Username)

	currentDate := time.Now()
	hangout := model.Hangout{
		PublicId: model.HangoutId(uuid.New()),
		HangoutDetails: model.HangoutDetails{
			Description: func() *string { s := "description"; return &s }(),
			Location:    "location",
			Duration:    10,
			Date:        currentDate,
		},
		CreatedBy:   model.IndividualId(individual.Username),
		Individuals: individualIds,
	}
	err = suite.pgStore.StoreHangoutOfIndividuals(suite.T().Context(), hangout)
	suite.Require().ErrorIs(err, storage.ErrHangoutParticipantDeleted, "expected error when inserting a hangout where a participant is soft-deleted")
}

func (suite *PostgresStoreTestSuite) TestStoreHangout_ReturnsError_WhenAHangoutWithTheSamePublicIdExists() {
	individualIds := make([]model.IndividualId, 0)
	individual := model.Individual{
		Username: "creator",
		Name:     "name",
		Email:    "email",
	}
	err := suite.pgStore.StoreIndividual(suite.T().Context(), individual)
	suite.Require().NoError(err, "expected no error when storing hangout creator")
	individualIds = append(individualIds, "creator")

	hangoutUuid := uuid.New()
	currentDate := time.Now()
	hangout := model.Hangout{
		PublicId: model.HangoutId(hangoutUuid),
		HangoutDetails: model.HangoutDetails{
			Description: func() *string { s := "description"; return &s }(),
			Location:    "location",
			Duration:    10,
			Date:        currentDate,
		},
		CreatedBy:   model.IndividualId(individual.Username),
		Individuals: individualIds,
	}
	err = suite.pgStore.StoreHangoutOfIndividuals(suite.T().Context(), hangout)
	suite.Require().NoError(err, "expected no error when inserting a valid hangout")
	err = suite.pgStore.StoreHangoutOfIndividuals(suite.T().Context(), hangout)
	suite.Require().ErrorIs(err, storage.ErrAlreadyExists, "expected an error when inserting a hangout with the same public id")
}
