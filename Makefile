.DEFAULT_GOAL := run
.PHONY: run build

include ./.env
export

run: build
	./bin/api

build:
	go build -o bin/api ./cmd/api