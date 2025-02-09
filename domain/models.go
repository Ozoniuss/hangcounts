package domain

import (
	"errors"
	"net/mail"
	"time"
)

type IndividualId uint64
type HangoutId uint64

type Email string
type Minutes int

var ErrInvalidEmail = errors.New("invalid email")
var ErrNegativeMinutes = errors.New("duration cannot be negative")

func newEmail(address string) (Email, error) {
	parsed, err := mail.ParseAddress(address)
	if err != nil {
		return Email(""), ErrInvalidEmail
	}

	return Email(parsed.Address), nil
}

func newMinute(d int) (Minutes, error) {
	if d < 0 {
		return 0, ErrNegativeMinutes
	}
	return Minutes(d), nil
}

type Individual struct {
	Id       IndividualId
	Name     string
	Email    Email
	Username string
}

type Hangout struct {
	Id          HangoutId
	Location    string
	Description *string
	Duration    Minutes
	Date        time.Time

	// Note that the use of Ids is deliberate. Hangouts are merely a collection
	// of individuals and if an Individual deletes his account, the other
	// participants may still be interested in the hangout.
	//
	// Given hangouts would still live if their creator removes their account
	// or participants remove their account, hangouts would work correctly with
	// eventually consistent Individuals. (Reading Individual data is obviously
	// functional under eventual consistency)
	CreatedBy   IndividualId
	Individuals []IndividualId
}
