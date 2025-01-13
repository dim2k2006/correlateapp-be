install:
	go mod download

lint:
	golangci-lint run

make format:
	./check-format.sh