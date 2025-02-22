package storage

import (
	"context"
	"errors"

	"github.com/Ozoniuss/hangcounts/domain/model"
)

type Individuals interface {
	StoreIndividual(context.Context, model.Individual) error
	GetIndividual(context.Context, model.IndividualId) error
	MarkIndividualAsDeleted(context.Context, model.IndividualId) error
	StoreHangoutOfIndividuals(context.Context, model.Hangout) error
	UpdateHangoutDetails(context.Context, model.HangoutId, model.HangoutDetails)
	UpdateHangoutParticipants(context.Context, model.HangoutId, []model.IndividualId)
}

// generic
var ErrAlreadyExists = errors.New("record already exists")
var ErrNotFound = errors.New("record is not found in database")
var ErrDeleted = errors.New("record is soft-deleted")
var ErrUnknown = errors.New("unknown database error")

// individual errors
var ErrIndividualEmailAlreadyExists = errors.New("email already exists")
var ErrIndividualUsernameAlreadyExists = errors.New("username already exists")

// hangout errors
var ErrHangoutCreatorNotFound = errors.New("hangout creator not found in database")
var ErrHangoutCreatorDeleted = errors.New("hangout creator is deleted")
var ErrHangoutParticipantNotFound = errors.New("hangout participant not found in database")
var ErrHangoutParticipantDeleted = errors.New("hangout participant is deleted")
var ErrParticipantHangoutNotFound = errors.New("hangout not found when inserting a participant")
var ErrParticipantIndividualNotFound = errors.New("individual not found when inserting a participant")
