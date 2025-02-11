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
var ErrAlreadyExists = errors.New("")
var ErrNotFound = errors.New("")
