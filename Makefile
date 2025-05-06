all: build

.PHONY: build check lint test

check: test

lint:
	gofmt -w .
	go vet ./...

test: lint
	go test -v -coverprofile=coverage.out ./...
