install:
	go mod download

dev:
	go run ./cmd/api/main.go

lint:
	golangci-lint run

format:
	./check-format.sh

test:
	go test ./...

PHONY: install lint format test