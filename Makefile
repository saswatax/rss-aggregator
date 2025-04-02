include ./.env
export

.PHONY: $(wildcard *)
.DEFAULT_GOAL := run-api

MIGRATION_FOLDER := ./database/migrations
SEEDER_FOLDER := ./database/seeders

run-api: build-api
	./bin/api

build-api:
	go build -o ./bin/api ./cmd/api

run-scraper: build-scraper
	./bin/scraper

build-scraper:
	go build -o ./bin/scraper ./cmd/scraper

clean:
	rm bin/*

migrate/create:
	goose create -dir ${MIGRATION_FOLDER} ${name} sql

migrate/up:
	goose -dir ${MIGRATION_FOLDER} postgres ${DB_URL} up

migrate/down:
	goose -dir ${MIGRATION_FOLDER} postgres ${DB_URL} down

migrate/reset:
	goose -dir ${MIGRATION_FOLDER} postgres ${DB_URL} reset

seed/create:
	goose create -dir ${SEEDER_FOLDER} ${name} sql

seed/up:
	goose -dir ${SEEDER_FOLDER} -no-versioning postgres ${DB_URL} up