.PHONY: build test vet fmt lint tidy

build:
	go build ./...

test:
	go test -race -count=1 ./...

vet:
	go vet ./...

fmt:
	@test -z "$$(gofmt -l .)" || (echo "needs gofmt:"; gofmt -l .; exit 1)

lint:
	golangci-lint run

tidy:
	go mod tidy
