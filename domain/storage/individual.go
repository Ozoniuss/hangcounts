package storage

import (
	"context"
	"errors"

	"github.com/Ozoniuss/hangcounts/domain/model"
)

type Individuals interface {
	StoreIndividual(context.Context, model.Individual) error
	GetIndividual(context.Context, model.IndividualId) error
}

// this can be enhanced with db-specific stuff
var ErrAlreadyExists = errors.New("record already exists")
var ErrNotFound = errors.New("record is not found in database")
var ErrDeleted = errors.New("record is soft-deleted")
var ErrUnknown = errors.New("unknown database error")
