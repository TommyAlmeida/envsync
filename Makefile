.PHONY: build test clean install lint

build:
	go build -o bin/envsync .

deps:
	go mod download
	go mod tidy

test:
	go test ./...

test-coverage:
	go test -coverprofile=coverage.out ./...
	go tool cover -html=coverage.out -o coverage.html

lint:
	golangci-lint run

clean:
	rm -rf bin/
	rm -f coverage.out coverage.html