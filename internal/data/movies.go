package data

import (
	"time"

	"github.com/sparrowsl/greenlight/internal/validator"
)

type Movie struct {
	ID        int64     `json:"id"`
	Title     string    `json:"title"`
	Runtime   Runtime   `json:"runtime,omitempty"`
	Year      int32     `json:"year,omitempty"`
	Genres    []string  `json:"genres,omitempty"`
	Version   int32     `json:"version"`
	CreatedAt time.Time `json:"created_at"`
}

func ValidateMovie(val *validator.Validator, movie *Movie) {
	val.Check(movie.Title != "", "title", "must be provided")
	val.Check(len(movie.Title) <= 500, "title", "must not be more than 500 bytes long")

	val.Check(movie.Year != 0, "year", "must be provided")
	val.Check(movie.Year >= 1888, "year", "must be greater than 1888")
	val.Check(movie.Year <= int32(time.Now().Year()), "year", "must not be in the future")

	val.Check(movie.Runtime != 0, "runtime", "must be provided")
	val.Check(movie.Runtime > 0, "runtime", "must be a positive integer")

	val.Check(movie.Genres != nil, "genres", "must be provided")
	val.Check(len(movie.Genres) >= 1, "genres", "must contain at least 1 genre")
	val.Check(len(movie.Genres) <= 5, "genres", "must not contain more than 5 genres")
	val.Check(validator.Unique(movie.Genres), "genres", "must not contain duplicate values")
}
