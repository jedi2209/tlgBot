# Go parameters
GOCMD=go
GOBUILD=$(GOCMD) build
GOCLEAN=$(GOCMD) clean
GOTEST=$(GOCMD) test
GOGET=$(GOCMD) get
GOMOD=$(GOCMD) mod

# Main package info
MAIN_PATH=./cmd/telegram-bot
BINARY_NAME=telegram-bot
BINARY_LINUX=$(BINARY_NAME)_linux

.PHONY: all build clean test coverage help

all: test build

## Build the binary file
build:
	$(GOBUILD) -o $(BINARY_NAME) -v $(MAIN_PATH)

## Build the binary file for Linux
build-linux:
	CGO_ENABLED=0 GOOS=linux GOARCH=amd64 $(GOBUILD) -o $(BINARY_LINUX) -v $(MAIN_PATH)

## Clean build files
clean:
	$(GOCLEAN)
	rm -f $(BINARY_NAME)
	rm -f $(BINARY_LINUX)
	rm -f coverage.out coverage.html

## Run tests
test:
	$(GOTEST) -v ./...

## Run tests with coverage
coverage:
	$(GOTEST) -v -coverprofile=coverage.out ./...
	$(GOCMD) tool cover -html=coverage.out

## Run tests with race detection
test-race:
	$(GOTEST) -v -race ./...

## Run tests with coverage and generate report
test-coverage:
	$(GOTEST) -v -coverprofile=coverage.out -covermode=atomic ./...
	$(GOCMD) tool cover -func=coverage.out
	$(GOCMD) tool cover -html=coverage.out -o coverage.html

## Run benchmarks
bench:
	$(GOTEST) -bench=. -benchmem ./...

## Run tests with verbose output and show coverage
test-verbose:
	$(GOTEST) -v -cover ./...

## Clean test cache and coverage files
test-clean:
	$(GOCMD) clean -testcache
	rm -f coverage.out coverage.html

## Download dependencies
deps:
	$(GOMOD) download
	$(GOMOD) tidy

## Run the application
run:
	$(GOCMD) run $(MAIN_PATH)

## Format code
fmt:
	$(GOCMD) fmt ./...

## Lint code
lint:
	golangci-lint run

## Help
help:
	@echo "Available commands:"
	@echo "  build         - Build the binary file"
	@echo "  build-linux   - Build the binary file for Linux"
	@echo "  clean         - Clean build files and test artifacts"
	@echo "  test          - Run tests"
	@echo "  test-race     - Run tests with race detection"
	@echo "  test-coverage - Run tests with coverage and generate report"
	@echo "  test-verbose  - Run tests with verbose output and show coverage"
	@echo "  test-clean    - Clean test cache and coverage files"
	@echo "  coverage      - Run tests with coverage (legacy)"
	@echo "  bench         - Run benchmarks"
	@echo "  deps          - Download dependencies"
	@echo "  run           - Run the application"
	@echo "  fmt           - Format code"
	@echo "  lint          - Lint code"
	@echo "  help          - Show this help" 