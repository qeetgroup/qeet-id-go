# Tool versions. Keep them in step with .github/workflows/lint.yml and security.yml.
GOLANGCI_LINT_VERSION ?= v2.15.0
GOVULNCHECK_VERSION   ?= v1.9.0

.PHONY: build test vet fmt lint vuln tidy ci

build:
	go build ./...

test:
	go test -race -count=1 ./...

vet:
	go vet ./...

fmt:
	@test -z "$$(gofmt -l .)" || (echo "needs gofmt:"; gofmt -l .; exit 1)

lint:
	golangci-lint run ./...

vuln:
	go run golang.org/x/vuln/cmd/govulncheck@$(GOVULNCHECK_VERSION) ./...

tidy:
	go mod tidy

# ci runs what the GitHub workflows run, in the same order.
ci: fmt vet build test lint vuln
