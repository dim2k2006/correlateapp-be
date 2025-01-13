install:
	go mod download

lint:
	golangci-lint run

format:
	./check-format.sh

test:
	go test ./...

PHONY: test