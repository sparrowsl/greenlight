package data

import "time"

type Movie struct {
	ID        int64
	Title     string
	Runtime   int32
	Year      int32
	Genres    []string
	Version   int32
	CreatedAt time.Time
}
