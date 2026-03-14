BINARY    := lazys3
CMD       := ./cmd/lazys3
GOFLAGS   :=
COVER_OUT := coverage.out
COVER_HTML := coverage.html

.PHONY: build build-duckdb run test test-verbose test-race cover cover-html fmt vet tidy clean pre-commit help

## build: compile the binary
build:
	go build $(GOFLAGS) -o $(BINARY) $(CMD)

## build-duckdb: build with DuckDB support (requires CGO)
build-duckdb:
	CGO_ENABLED=1 go build $(GOFLAGS) -o $(BINARY) $(CMD)

## run: run without building a binary
run:
	go run $(CMD)

## test: run all tests
test:
	go test ./...

## test-verbose: run all tests with verbose output
test-verbose:
	go test -v ./...

## test-race: run all tests with the race detector
test-race:
	go test -race ./...

## cover: run tests with coverage and print summary per package
cover:
	go test -coverprofile=$(COVER_OUT) -covermode=atomic ./...
	@echo ""
	@echo "==> Coverage summary:"
	go tool cover -func=$(COVER_OUT)

## cover-html: generate HTML coverage report and open it
cover-html: cover
	go tool cover -html=$(COVER_OUT) -o $(COVER_HTML)
	@echo "Coverage report written to $(COVER_HTML)"

## fmt: format all Go source files in place
fmt:
	gofmt -w .

## vet: run static analysis
vet:
	go vet ./...

## tidy: tidy and verify the module graph
tidy:
	go mod tidy

## clean: remove build artifacts
clean:
	rm -f $(BINARY) $(COVER_OUT) $(COVER_HTML)

## pre-commit: fmt check + vet + tests (must all pass before committing)
pre-commit:
	@echo "==> Checking formatting..."
	@out=$$(gofmt -l .); if [ -n "$$out" ]; then echo "$$out"; echo "Run 'make fmt' to fix formatting."; exit 1; fi
	@echo "==> Running vet..."
	go vet ./...
	@echo "==> Running tests..."
	go test ./...
	@echo "==> All checks passed."

## help: show this help
help:
	@grep -E '^## ' Makefile | sed 's/^## //'
