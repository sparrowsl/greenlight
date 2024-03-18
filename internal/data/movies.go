package data

import (
	"context"
	"database/sql"
	"errors"
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

func (m *MovieModel) GetAll(title string, genres []string, filters Filters) ([]Movie, error) {
	statement := `SELECT id, title, year, runtime, created_at, genres, version 
                FROM movies
                WHERE (LOWER(title) = LOWER($1) OR $1 = '')
                AND (genres @> $2 OR $2 = '{}')
                ORDER BY id`

	ctx, cancel := context.WithTimeout(context.Background(), time.Second*3)
	defer cancel()

	rows, err := m.DB.QueryContext(ctx, statement, title, pq.Array(genres))
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	movies := []Movie{}

	for rows.Next() {
		var movie Movie

		err := rows.Scan(&movie.ID, &movie.Title, &movie.Year, &movie.Runtime, &movie.CreatedAt, pq.Array(&movie.Genres), &movie.Version)
		if err != nil {
			return nil, err
		}

		movies = append(movies, movie)
	}

	if err = rows.Err(); err != nil {
		return nil, err
	}

	return movies, nil
}

func (m *MovieModel) Get(id int64) (*Movie, error) {
	if id < 1 {
		return nil, ErrRecordNotFound
	}

	statement := `SELECT id, title, year, runtime, created_at, genres, version 
                FROM movies
                WHERE id = $1`

	var movie Movie

	ctx, cancel := context.WithTimeout(context.Background(), time.Second*3)
	defer cancel()

	row := m.DB.QueryRowContext(ctx, statement, id)
	err := row.Scan(&movie.ID, &movie.Title, &movie.Year, &movie.Runtime, &movie.CreatedAt, pq.Array(&movie.Genres), &movie.Version)
	if err != nil {
		switch {
		case errors.Is(err, sql.ErrNoRows):
			return nil, ErrRecordNotFound
		default:
			return nil, err
		}
	}

	return &movie, nil
}

func (m *MovieModel) Update(movie *Movie) error {
	statement := `UPDATE movies
                SET title = $1, year = $2, runtime = $3, genres = $4, version = version + 1 
                WHERE id = $5 AND version = $6
                RETURNING version`

	row := m.DB.QueryRow(statement, movie.Title, movie.Year, movie.Runtime, pq.Array(movie.Genres), movie.ID, movie.Version)
	if err := row.Scan(&movie.Version); err != nil {
		switch {
		case errors.Is(err, sql.ErrNoRows):
			return ErrEditConflict
		default:
			return err
		}
	}

	return nil
}

func (m *MovieModel) Delete(id int64) error {
	if id < 1 {
		return ErrRecordNotFound
	}

	statement := `DELETE FROM movies
                WHERE id = $1`

	ctx, cancel := context.WithTimeout(context.Background(), time.Second*3)
	defer cancel()

	result, err := m.DB.ExecContext(ctx, statement, id)
	if err != nil {
		return err
	}

	rows, err := result.RowsAffected()
	if err != nil {
		return err
	}

	if rows == 0 {
		return ErrRecordNotFound
	}

	return nil
}
