.PHONY: lint test build vuln check
.DEFAULT_GOAL := check

lint:
	golangci-lint run --config .github/.golangci.yml
	yamllint -c .github/.yamllint .

build:
	go build ./...

test:
	go test ./...

vuln:
	govulncheck ./...

check: lint build test vuln
