-- +goose Up
-- +goose StatementBegin
CREATE EXTENSION IF NOT EXISTS CITEXT;

CREATE TABLE IF NOT EXISTS movies(
  id bigserial PRIMARY KEY,
  created_at timestamp(0) WITH time zone NOT NULL DEFAULT NOW(),
  title TEXT NOT NULL, 
  year integer NOT NULL,
  runtime INTEGER NOT NULL,
  genres text[] NOT NULL,
  version INTEGER NOT NULL DEFAULT 1
);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE IF EXISTS movies;
-- +goose StatementEnd
