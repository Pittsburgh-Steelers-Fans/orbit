.PHONY: build test lint migrate docker-up generate tidy

BIN := bin
MODULE := github.com/Pittsburgh-Steelers-Fans/orbit

build:
	go build -o $(BIN)/orbit ./cmd/server
	go build -o $(BIN)/orbitctl ./cmd/orbitctl

test:
	go test ./... -race -count=1

lint:
	go vet ./...
	gofmt -l .

migrate:
	@echo "applying migrations from db/migrations"
	@ls db/migrations

generate:
	sqlc generate

docker-up:
	docker compose up --build

tidy:
	go mod tidy
