package main

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

type HangoutsService struct {
}

func (h *HangoutsService) NewHangout(Atendees []Individual) (Hangout, error) {
	return Hangout{}, nil
}

func (h *HangoutsService) DeleteHangout(hangoutId uint64) (Hangout, error) {
	return Hangout{}, nil
}

func main() {

}
