package entities

import "time"

type Song struct {
	ID          int64     `db:"id"`
	Group       string    `db:"group"`
	Name        string    `db:"name"`
	Link        string    `db:"link"`
	ReleaseDate time.Time `db:"release_date"`
}
