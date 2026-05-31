.PHONY: build test lint tidy check examples

build:
	go build ./...

test:
	go test -race -cover ./...

lint:
	golangci-lint run ./...

tidy:
	go mod tidy

examples:
	go build -o /dev/null ./examples/...

check:
	./scripts/check.sh
