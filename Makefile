install:
	go mod download

lint:
	golangci-lint run