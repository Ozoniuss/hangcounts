package model

import (
	"errors"
	"time"

	"github.com/google/uuid"
)

type HangoutId uuid.UUID
type Minutes int

func NewMinute(d int) (Minutes, error) {
	if d < 0 {
		return 0, errors.New("negative miuntes")
	}
	return Minutes(d), nil
}

type HangoutDetails struct {
	Location    string
	Description *string
	Duration    Minutes
	Date        time.Time
}

type Hangout struct {
	PublicId HangoutId
	HangoutDetails

	// Note that the use of Ids is deliberate. Hangouts are merely a collection
	// of individuals and if an Individual deletes his account, the other
	// participants may still be interested in the hangout.
	//
	// Given hangouts would still live if their creator removes their account
	// or participants remove their account, hangouts would work correctly with
	// eventually consistent Individuals. (Reading Individual data is obviously
	// functional under eventual consistency)
	CreatedBy IndividualId

	// The creator must be part of the individuals. This should be enforced by
	// the aggregate.
	Individuals []IndividualId
}
