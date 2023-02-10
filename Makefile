ifneq (,$(wildcard ./.env))
    include .env
    export
endif

build:
	go build -o server

run:
	go run .

test:
	go test ./...

test-integration:
	go test ./... --tags=integration

compose-up:
	docker-compose up --build -d mongo mongo-express && docker-compose logs -f

compose-down:
	docker-compose down --remove-orphans

.PHONY: test test-integration compose-up compose-down