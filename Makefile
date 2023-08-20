DB_URL=postgresql://root:postgres@localhost:5432/bwabackerdb?sslmode=disable

## help: print this help messages
help:
	@echo 'Usage:'
	@sed -n 's/^##//p' ${MAKEFILE_LIST} | column -t -s ':' | sed -e 's/^/ /'

confirm:
	@echo 'Are you sure? [y/N] ' && read ans && [ $${ans:-N} = y ]

## db/createdb: create db with docker using predefined database name (bwabackerdb) 
db/create:
	docker exec -it postgres12 createdb --username=root --owner=root bwabackerdb

## db/createdb: drop db 
db/drop: confirm
	docker exec -it postgres12 dropdb bwabackerdb

## db/psql: connect 
db/psql:
	docker exec -it postgres12 psql ${DB_URL}
	
## db/migrations/up: apply by 1 up database migrations
db/migration/up:
	goose -dir ./db/migrations postgres ${DB_URL} up-by-one

## db/migrations/all: apply all up database migrations
db/migration/all:
	@echo 'running all migrations...'
	goose -dir ./db/migrations postgres ${DB_URL} up

## db/migrations/down: delete by one database migrations
db/migration/down:
	goose -dir ./db/migrations postgres ${DB_URL} down

## db/migrations/reset: reset all database migrations
db/migration/reset:
	goose -dir ./db/migrations postgres ${DB_URL} reset

## db/migrations/new name=$1: create a new database migrations
db/migration/new:
	@echo 'create migrations file for ${name}'
	goose -dir ./db/migrations postgres ${DB_URL} create ${name} sql 
	goose -dir ./db/migrations postgres ${DB_URL} fix

.PHONY: help confirm db/create db/drop db/psql db/migration/up db/migration/all db/migration/down db/migration/reset db/migration/new
