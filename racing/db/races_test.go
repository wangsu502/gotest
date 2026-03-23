package db

import (
	"database/sql"
	"testing"

	_ "github.com/mattn/go-sqlite3"

	"git.neds.sh/matty/entain/racing/proto/racing"
)

func setupTestDB(t *testing.T) *racesRepo {
	t.Helper()

	// Use in-memory SQLite so each test gets an isolated DB that is
	// automatically released when the test ends. Does not touch racing.db.
	db, err := sql.Open("sqlite3", ":memory:")
	if err != nil {
		t.Fatalf("failed to open test db: %v", err)
	}

	t.Cleanup(func() { db.Close() })

	repo := &racesRepo{db: db}

	_, err = db.Exec(`CREATE TABLE races (
		id INTEGER PRIMARY KEY,
		meeting_id INTEGER,
		name TEXT,
		number INTEGER,
		visible INTEGER,
		advertised_start_time DATETIME
	)`)
	if err != nil {
		t.Fatalf("failed to create table: %v", err)
	}

	// Insert test data.
	testRaces := []struct {
		id        int
		meetingID int
		name      string
		number    int
		visible   int
		startTime string
	}{
		{1, 1, "Race A", 1, 1, "2026-03-23T10:00:00Z"},
		{2, 1, "Race B", 2, 0, "2026-03-23T11:00:00Z"},
		{3, 2, "Race C", 3, 1, "2026-03-23T12:00:00Z"},
		{4, 3, "Race D", 4, 0, "2026-03-23T13:00:00Z"},
	}

	for _, r := range testRaces {
		_, err := db.Exec(
			`INSERT INTO races(id, meeting_id, name, number, visible, advertised_start_time) VALUES (?,?,?,?,?,?)`,
			r.id, r.meetingID, r.name, r.number, r.visible, r.startTime,
		)
		if err != nil {
			t.Fatalf("failed to insert test race %d: %v", r.id, err)
		}
	}

	return repo
}

func boolPtr(b bool) *bool {
	return &b
}

func TestList_NoFilter(t *testing.T) {
	repo := setupTestDB(t)

	races, err := repo.List(nil, "")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(races) != 4 {
		t.Errorf("expected 4 races, got %d", len(races))
	}
}

func TestList_VisibleOnly(t *testing.T) {
	repo := setupTestDB(t)

	races, err := repo.List(&racing.ListRacesRequestFilter{
		VisibleOnly: boolPtr(true),
	}, "")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(races) != 2 {
		t.Errorf("expected 2 visible races, got %d", len(races))
	}
	for _, r := range races {
		if !r.Visible {
			t.Errorf("race %d should be visible", r.Id)
		}
	}
}

func TestList_VisibleOnlyFalse(t *testing.T) {
	repo := setupTestDB(t)

	races, err := repo.List(&racing.ListRacesRequestFilter{
		VisibleOnly: boolPtr(false),
	}, "")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(races) != 4 {
		t.Errorf("expected 4 races (all), got %d", len(races))
	}
}

func TestList_MeetingIdsFilter(t *testing.T) {
	repo := setupTestDB(t)

	races, err := repo.List(&racing.ListRacesRequestFilter{
		MeetingIds: []int64{1},
	}, "")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(races) != 2 {
		t.Errorf("expected 2 races for meeting 1, got %d", len(races))
	}
}

func TestList_VisibleOnlyWithMeetingIds(t *testing.T) {
	repo := setupTestDB(t)

	races, err := repo.List(&racing.ListRacesRequestFilter{
		MeetingIds:  []int64{1},
		VisibleOnly: boolPtr(true),
	}, "")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(races) != 1 {
		t.Errorf("expected 1 visible race for meeting 1, got %d", len(races))
	}
	if len(races) > 0 && races[0].Name != "Race A" {
		t.Errorf("expected Race A, got %s", races[0].Name)
	}
}

func TestList_DefaultOrder(t *testing.T) {
	repo := setupTestDB(t)

	// Empty orderBy defaults to advertised_start_time ASC.
	races, err := repo.List(nil, "")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(races) != 4 {
		t.Fatalf("expected 4 races, got %d", len(races))
	}
	// Test data times: Race A 10:00, Race B 11:00, Race C 12:00, Race D 13:00.
	expected := []string{"Race A", "Race B", "Race C", "Race D"}
	for i, name := range expected {
		if races[i].Name != name {
			t.Errorf("position %d: expected %s, got %s", i, name, races[i].Name)
		}
	}
}

func TestList_OrderByStartTimeDesc(t *testing.T) {
	repo := setupTestDB(t)

	races, err := repo.List(nil, "advertised_start_time DESC")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(races) != 4 {
		t.Fatalf("expected 4 races, got %d", len(races))
	}
	expected := []string{"Race D", "Race C", "Race B", "Race A"}
	for i, name := range expected {
		if races[i].Name != name {
			t.Errorf("position %d: expected %s, got %s", i, name, races[i].Name)
		}
	}
}

func TestList_OrderByName(t *testing.T) {
	repo := setupTestDB(t)

	races, err := repo.List(nil, "name ASC")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(races) != 4 {
		t.Fatalf("expected 4 races, got %d", len(races))
	}
	expected := []string{"Race A", "Race B", "Race C", "Race D"}
	for i, name := range expected {
		if races[i].Name != name {
			t.Errorf("position %d: expected %s, got %s", i, name, races[i].Name)
		}
	}
}

func TestList_InvalidOrderField(t *testing.T) {
	repo := setupTestDB(t)

	_, err := repo.List(nil, "invalid_field ASC")
	if err == nil {
		t.Fatal("expected error for invalid order field, got nil")
	}
}

func TestList_InvalidOrderDirection(t *testing.T) {
	repo := setupTestDB(t)

	_, err := repo.List(nil, "name INVALID")
	if err == nil {
		t.Fatal("expected error for invalid order direction, got nil")
	}
}

func TestGet_Found(t *testing.T) {
	repo := setupTestDB(t)

	race, err := repo.Get(1)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if race.Id != 1 {
		t.Errorf("expected id 1, got %d", race.Id)
	}
	if race.Name != "Race A" {
		t.Errorf("expected Race A, got %s", race.Name)
	}
}

func TestGet_NotFound(t *testing.T) {
	repo := setupTestDB(t)

	_, err := repo.Get(999)
	if err == nil {
		t.Fatal("expected error for non-existent race, got nil")
	}
}
