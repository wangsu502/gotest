package db

import (
	"database/sql"
	"testing"

	_ "github.com/mattn/go-sqlite3"

	"git.neds.sh/matty/entain/sports/proto/sports"
)

func setupTestDB(t *testing.T) *eventsRepo {
	t.Helper()

	// In-memory SQLite so each test gets an isolated DB,
	// automatically cleaned up when the test ends.
	db, err := sql.Open("sqlite3", ":memory:")
	if err != nil {
		t.Fatalf("failed to open test db: %v", err)
	}

	t.Cleanup(func() { db.Close() })

	repo := &eventsRepo{db: db}

	_, err = db.Exec(`CREATE TABLE events (
		id INTEGER PRIMARY KEY,
		name TEXT,
		sport TEXT,
		visible INTEGER,
		advertised_start_time DATETIME
	)`)
	if err != nil {
		t.Fatalf("failed to create table: %v", err)
	}

	// Test data: 2 visible (id 1,3) and 2 hidden (id 2,4), ordered by start time.
	testEvents := []struct {
		id        int
		name      string
		sport     string
		visible   int
		startTime string
	}{
		{1, "Grand Final", "football", 1, "2026-03-23T10:00:00Z"},
		{2, "Semifinal", "tennis", 0, "2026-03-23T11:00:00Z"},
		{3, "Quarter Cup", "basketball", 1, "2026-03-23T12:00:00Z"},
		{4, "League Match", "cricket", 0, "2026-03-23T13:00:00Z"},
	}

	for _, e := range testEvents {
		_, err := db.Exec(
			`INSERT INTO events(id, name, sport, visible, advertised_start_time) VALUES (?,?,?,?,?)`,
			e.id, e.name, e.sport, e.visible, e.startTime,
		)
		if err != nil {
			t.Fatalf("failed to insert test event %d: %v", e.id, err)
		}
	}

	return repo
}

func boolPtr(b bool) *bool {
	return &b
}

// No filter should return all 4 events.
func TestList_NoFilter(t *testing.T) {
	repo := setupTestDB(t)

	events, err := repo.List(nil)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(events) != 4 {
		t.Errorf("expected 4 events, got %d", len(events))
	}
}

// visible_only=true should only return events with visible=1.
func TestList_VisibleOnly(t *testing.T) {
	repo := setupTestDB(t)

	events, err := repo.List(&sports.ListEventsRequestFilter{
		VisibleOnly: boolPtr(true),
	})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(events) != 2 {
		t.Errorf("expected 2 visible events, got %d", len(events))
	}
	for _, e := range events {
		if !e.Visible {
			t.Errorf("event %d should be visible", e.Id)
		}
	}
}

// visible_only=false should not filter anything.
func TestList_VisibleOnlyFalse(t *testing.T) {
	repo := setupTestDB(t)

	events, err := repo.List(&sports.ListEventsRequestFilter{
		VisibleOnly: boolPtr(false),
	})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(events) != 4 {
		t.Errorf("expected 4 events (all), got %d", len(events))
	}
}

// Results should be sorted by advertised_start_time ASC by default.
func TestList_DefaultOrder(t *testing.T) {
	repo := setupTestDB(t)

	events, err := repo.List(nil)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(events) != 4 {
		t.Fatalf("expected 4 events, got %d", len(events))
	}
	// Test data: Grand Final 10:00, Semifinal 11:00, Quarter Cup 12:00, League Match 13:00.
	expected := []string{"Grand Final", "Semifinal", "Quarter Cup", "League Match"}
	for i, name := range expected {
		if events[i].Name != name {
			t.Errorf("position %d: expected %s, got %s", i, name, events[i].Name)
		}
	}
}

// An empty filter (no fields set) should behave the same as nil.
func TestList_EmptyFilter(t *testing.T) {
	repo := setupTestDB(t)

	events, err := repo.List(&sports.ListEventsRequestFilter{})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(events) != 4 {
		t.Errorf("expected 4 events, got %d", len(events))
	}
}
