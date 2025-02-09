package storage

import (
	"errors"

	"github.com/Ozoniuss/hangcounts/domain/model"
)

type Individuals interface {
	StoreIndividual(model.Individual) error
	GetIndividual(model.IndividualId) error
}

// this can be enhanced with db-specific stuff
var ErrAlreadyExists = errors.New("")
var ErrNotFound = errors.New("")
