package aggregate

import (
	"errors"

	"github.com/Ozoniuss/hangcounts/domain/model"
	"github.com/Ozoniuss/hangcounts/domain/storage"
)

// Explicit Individual validation errors
type IndividualValidationError error

var ErrInvalidEmail = errors.New("invalid email")
var ErrEmptyName = errors.New("name cannot be empty")
var ErrEmptyUsername = errors.New("username cannot be empty")
var ErrDuplicateUser = errors.New("user already exists")

// Explicit Hangout errors
var ErrNegativeMinutes = errors.New("duration cannot be negative")

type IndividualAgg struct {
	model.Individual

	storage storage.Individuals
}

func (agg *IndividualAgg) NewIndividual(id uint64, name, email, username string) error {
	var errs error

	// may move those to their own time
	if name == "" {
		errs = errors.Join(errs, ErrEmptyName)
	}
	if username == "" {
		errs = errors.Join(errs, ErrEmptyUsername)
	}
	_email, err := model.NewEmail(email)
	if err != nil {
		errs = errors.Join(errs, err)
	}
	// eager return to avoid database call
	if errs != nil {
		return IndividualValidationError(errs)
	}

	agg.Individual = model.Individual{
		Id:       model.IndividualId(id),
		Name:     name,
		Email:    model.Email(_email),
		Username: username,
	}
	err = agg.storage.StoreIndividual(agg.Individual)
	if errors.Is(err, storage.ErrNotFound) {
		return IndividualValidationError(ErrDuplicateUser)
	}

	return nil
}
