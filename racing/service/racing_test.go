package service

import (
	"context"
	"database/sql"
	"testing"
	"time"

	"git.neds.sh/matty/entain/racing/proto/racing"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/types/known/timestamppb"
)

func TestSetRaceStatuses_Open(t *testing.T) {
	now := time.Date(2026, 3, 23, 12, 0, 0, 0, time.UTC)
	races := []*racing.Race{
		{Id: 1, AdvertisedStartTime: timestamppb.New(now.Add(1 * time.Hour))},
	}

	setRaceStatuses(races, now)

	if races[0].Status != "OPEN" {
		t.Errorf("expected OPEN, got %s", races[0].Status)
	}
}

func TestSetRaceStatuses_Closed(t *testing.T) {
	now := time.Date(2026, 3, 23, 12, 0, 0, 0, time.UTC)
	races := []*racing.Race{
		{Id: 1, AdvertisedStartTime: timestamppb.New(now.Add(-1 * time.Hour))},
	}

	setRaceStatuses(races, now)

	if races[0].Status != "CLOSED" {
		t.Errorf("expected CLOSED, got %s", races[0].Status)
	}
}

func TestSetRaceStatuses_Mixed(t *testing.T) {
	now := time.Date(2026, 3, 23, 12, 0, 0, 0, time.UTC)
	races := []*racing.Race{
		{Id: 1, AdvertisedStartTime: timestamppb.New(now.Add(1 * time.Hour))},
		{Id: 2, AdvertisedStartTime: timestamppb.New(now.Add(-1 * time.Hour))},
		{Id: 3, AdvertisedStartTime: timestamppb.New(now.Add(24 * time.Hour))},
	}

	setRaceStatuses(races, now)

	expected := []string{"OPEN", "CLOSED", "OPEN"}
	for i, e := range expected {
		if races[i].Status != e {
			t.Errorf("race %d: expected %s, got %s", races[i].Id, e, races[i].Status)
		}
	}
}

// mockRacesRepo is a minimal mock of db.RacesRepo for service-layer tests.
type mockRacesRepo struct {
	races map[int64]*racing.Race
}

func (m *mockRacesRepo) Init() error { return nil }

// List is not used in service-layer tests but required by the RacesRepo interface.
func (m *mockRacesRepo) List(_ *racing.ListRacesRequestFilter, _ string) ([]*racing.Race, error) {
	return nil, nil
}

func (m *mockRacesRepo) Get(id int64) (*racing.Race, error) {
	r, ok := m.races[id]
	if !ok {
		return nil, sql.ErrNoRows
	}
	return r, nil
}

func TestGetRace_Found(t *testing.T) {
	now := time.Date(2026, 3, 23, 12, 0, 0, 0, time.UTC)
	// Build a mock repo with one race starting in the future.
	repo := &mockRacesRepo{
		races: map[int64]*racing.Race{
			1: {Id: 1, Name: "Race A", AdvertisedStartTime: timestamppb.New(now.Add(1 * time.Hour))},
		},
	}
	svc := &racingService{racesRepo: repo}

	race, err := svc.GetRace(context.Background(), &racing.GetRaceRequest{Id: 1})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if race.Name != "Race A" {
		t.Errorf("expected Race A, got %s", race.Name)
	}
	if race.Status != "OPEN" {
		t.Errorf("expected OPEN status, got %s", race.Status)
	}
}

func TestGetRace_NotFound(t *testing.T) {
	repo := &mockRacesRepo{races: map[int64]*racing.Race{}}
	svc := &racingService{racesRepo: repo}

	_, err := svc.GetRace(context.Background(), &racing.GetRaceRequest{Id: 999})
	if err == nil {
		t.Fatal("expected error for non-existent race")
	}
	st, ok := status.FromError(err)
	if !ok {
		t.Fatalf("expected gRPC status error, got %v", err)
	}
	if st.Code() != codes.NotFound {
		t.Errorf("expected NotFound, got %s", st.Code())
	}
}
