package db

import (
	"math/rand"
	"time"

	"syreclabs.com/go/faker"
)

var sportTypes = []string{"football", "tennis", "basketball", "cricket", "baseball"}

func (e *eventsRepo) seed() error {
	statement, err := e.db.Prepare(`CREATE TABLE IF NOT EXISTS events (id INTEGER PRIMARY KEY, name TEXT, sport TEXT, visible INTEGER, advertised_start_time DATETIME)`)
	if err == nil {
		_, err = statement.Exec()
	}

	// Seed 100 dummy events for testing/demo purposes.
	for i := 1; i <= 100; i++ {
		statement, err = e.db.Prepare(`INSERT OR IGNORE INTO events(id, name, sport, visible, advertised_start_time) VALUES (?,?,?,?,?)`)
		if err == nil {
			_, err = statement.Exec(
				i,
				faker.Team().Name(),
				sportTypes[rand.Intn(len(sportTypes))],
				faker.Number().Between(0, 1),
				faker.Time().Between(time.Now().AddDate(0, 0, -1), time.Now().AddDate(0, 0, 2)).Format(time.RFC3339),
			)
		}
	}

	return err
}
