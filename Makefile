.PHONY: lint test build vuln check

lint:
	golangci-lint run --config .github/.golangci.yml

build:
	go build ./...

test:
	go test ./...

vuln:
	govulncheck ./...

check: lint build test vuln
