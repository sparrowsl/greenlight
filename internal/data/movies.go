package data

import "time"

type Movie struct {
	ID        int64     `json:"id"`
	Title     string    `json:"title"`
	Runtime   int32     `json:"runtime"`
	Year      int32     `json:"year"`
	Genres    []string  `json:"genres"`
	Version   int32     `json:"version"`
	CreatedAt time.Time `json:"created_at"`
}
