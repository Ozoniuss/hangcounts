package main

import "github.com/Ozoniuss/hangcounts/domain"

type HangoutsService struct {
}

func (h *HangoutsService) NewHangout(Atendees []domain.Individual) (domain.Hangout, error) {
	return domain.Hangout{}, nil
}

func (h *HangoutsService) DeleteHangout(hangoutId uint64) (domain.Hangout, error) {
	return domain.Hangout{}, nil
}

func main() {

}
