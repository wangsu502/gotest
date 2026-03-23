package service

import (
	"testing"
	"time"

	"git.neds.sh/matty/entain/racing/proto/racing"
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
