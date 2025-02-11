package model

import (
	"fmt"
	"net/mail"
)

type IndividualId uint64

type Email string

func NewEmail(address string) (Email, error) {
	parsed, err := mail.ParseAddress(address)
	if err != nil {
		return Email(""), fmt.Errorf("could not parse email: %w", err)
	}

	return Email(parsed.Address), nil
}

type Individual struct {
	Id       IndividualId
	Name     string
	Email    Email
	Username string
}
