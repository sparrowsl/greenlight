include .env

## help: print this help message
help:
	@echo 'Usage:'
	@sed -n 's/^##//p' ${MAKEFILE_LIST} | column -t -s ':' | sed -e 's/^/ /'

# Confirmation target
confirm:
	@echo -n 'Are you sure? [y/N]' && read ans && [ $${ans:-N} = y ]

## run/api: run the ./cmd/api/ application
run/api:
	go run ./cmd/api

## db/psql: connect to database using psql
db/psql:
	psql ${DATABASE_URL}

## db/migrations/new name=$1: creates a new database migration
db/migrations/new:
	@echo 'Creating migration files for ${name}'
	goose -dir=migrations create ${name} sql

## db/migrations/up: apply all up database migrations
db/migrations/up: confirm
	@echo 'Running up migrations...'
	goose -dir=migrations postgres ${DATABASE_URL} up

## db/migrations/down: apply all down database migrations
db/migrations/down: confirm
	@echo 'Running down migrations...'
	goose -dir=migrations postgres ${DATABASE_URL} down
