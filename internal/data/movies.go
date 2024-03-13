package data

import (
	"database/sql"
	"time"

	"github.com/lib/pq"
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

type MovieModel struct {
	DB *sql.DB
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

func (m *MovieModel) Insert(movie *Movie) error {
	statement := `INSERT INTO movies (title, year, runtime, genres)
                VALUES ($1, $2, $3, $4)
                RETURNING id, created_at, version`

	row := m.DB.QueryRow(statement, movie.Title, movie.Year, movie.Runtime, pq.Array(movie.Genres))

	return row.Scan(&movie.ID, &movie.CreatedAt, &movie.Version)
}

func (m *MovieModel) Get(id int64) (*Movie, error) {
	return &Movie{}, nil
}

func (m *MovieModel) Update(movie *Movie) error {
	return nil
}

func (m *MovieModel) Delete(id int64) error {
	return nil
}
