package service

import (
	"context"
	"fmt"
	"testing"

	"git.neds.sh/matty/entain/sports/proto/sports"
)

// mockEventsRepo is a minimal mock of db.EventsRepo for service-layer tests.
type mockEventsRepo struct {
	events []*sports.Event
	err    error
}

func (m *mockEventsRepo) Init() error { return nil }

// List returns pre-configured events/error, filter arg is ignored in mock.
func (m *mockEventsRepo) List(_ *sports.ListEventsRequestFilter) ([]*sports.Event, error) {
	return m.events, m.err
}

// Normal case: repo returns 2 events, verify count and first event name.
func TestListEvents_Success(t *testing.T) {
	repo := &mockEventsRepo{
		events: []*sports.Event{
			{Id: 1, Name: "Grand Final", Sport: "football"},
			{Id: 2, Name: "Semifinal", Sport: "tennis"},
		},
	}
	svc := &sportsService{eventsRepo: repo}

	resp, err := svc.ListEvents(context.Background(), &sports.ListEventsRequest{})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(resp.Events) != 2 {
		t.Errorf("expected 2 events, got %d", len(resp.Events))
	}
	if resp.Events[0].Name != "Grand Final" {
		t.Errorf("expected Grand Final, got %s", resp.Events[0].Name)
	}
}

// Empty result from repo should return an empty slice, not an error.
func TestListEvents_Empty(t *testing.T) {
	repo := &mockEventsRepo{events: []*sports.Event{}}
	svc := &sportsService{eventsRepo: repo}

	resp, err := svc.ListEvents(context.Background(), &sports.ListEventsRequest{})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(resp.Events) != 0 {
		t.Errorf("expected 0 events, got %d", len(resp.Events))
	}
}

// Repo error should propagate through the service layer.
func TestListEvents_RepoError(t *testing.T) {
	repo := &mockEventsRepo{err: fmt.Errorf("db connection lost")}
	svc := &sportsService{eventsRepo: repo}

	_, err := svc.ListEvents(context.Background(), &sports.ListEventsRequest{})
	if err == nil {
		t.Fatal("expected error from repo, got nil")
	}
}
