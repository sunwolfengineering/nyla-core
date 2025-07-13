VERSION=$(shell git describe --tags --candidates=1 --dirty)
BUILD_FLAGS=-ldflags="-X main.Version=$(VERSION)" -trimpath
SRC=$(shell find . -name '*.go') go.mod
.PHONY: nyla-core

nyla-core:
	go build $(BUILD_FLAGS) -o bin/nyla-core ./cmd/nyla-core

# Database operations (built-in migration system)
seed:
	rm -f nyla.db nyla.db-shm nyla.db-wal
	go run migrations/seed/main.go

