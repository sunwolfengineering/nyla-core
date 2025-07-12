VERSION=$(shell git describe --tags --candidates=1 --dirty)
BUILD_FLAGS=-ldflags="-X main.Version=$(VERSION)" -trimpath
SRC=$(shell find . -name '*.go') go.mod
.PHONY: nyla-api nyla-ui

nyla-api:
	go build $(BUILD_FLAGS) -o bin/nyla-api ./cmd/api

nyla-ui:
	go build $(BUILD_FLAGS) -o bin/nyla-ui ./cmd/ui

migrate:
	goose up

migrate-reset:
	goose reset

migrate-status:
	goose status

seed:
	goose reset
	goose up
	go run migrations/seed/main.go
	sqlite3 nyla.db ".mode tabs" ".import data/dump.data events"

