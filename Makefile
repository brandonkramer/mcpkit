.PHONY: build test lint tidy check examples

build:
	go build ./...

test:
	go test -race -cover ./...

lint:
	./scripts/lint.sh

tidy:
	go mod tidy

examples:
	go build -o /dev/null ./examples/...

check:
	./scripts/check.sh
