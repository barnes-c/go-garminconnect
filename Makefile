.PHONY: lint test build check

lint:
	golangci-lint run

build:
	go build ./...

test:
	go test ./...

check: lint build test
	govulncheck ./...
