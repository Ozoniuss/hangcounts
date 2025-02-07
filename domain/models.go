package domain

import "time"

type Individual struct {
	Id       uint64
	Name     string
	Email    string
	Username string
}

type Hangout struct {
	Id          uint64
	Location    string
	Description string
	Duration    time.Duration
	Date        time.Time
	CreatedBy   Individual

	Individuals []Individual
}
