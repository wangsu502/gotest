package service

import (
	"time"

	"git.neds.sh/matty/entain/racing/db"
	"git.neds.sh/matty/entain/racing/proto/racing"
	"golang.org/x/net/context"
)

type Racing interface {
	// ListRaces will return a collection of races.
	ListRaces(ctx context.Context, in *racing.ListRacesRequest) (*racing.ListRacesResponse, error)
}

// racingService implements the Racing interface.
type racingService struct {
	racesRepo db.RacesRepo
}

// NewRacingService instantiates and returns a new racingService.
func NewRacingService(racesRepo db.RacesRepo) Racing {
	return &racingService{racesRepo}
}

func (s *racingService) ListRaces(ctx context.Context, in *racing.ListRacesRequest) (*racing.ListRacesResponse, error) {
	races, err := s.racesRepo.List(in.Filter, in.OrderBy)
	if err != nil {
		return nil, err
	}

	setRaceStatuses(races, time.Now())

	return &racing.ListRacesResponse{Races: races}, nil
}

// setRaceStatuses derives the status for each race based on current time.
// Races with advertised_start_time in the past are "CLOSED", otherwise "OPEN".
func setRaceStatuses(races []*racing.Race, now time.Time) {
	for _, r := range races {
		if r.AdvertisedStartTime.AsTime().Before(now) {
			r.Status = "CLOSED"
		} else {
			r.Status = "OPEN"
		}
	}
}
