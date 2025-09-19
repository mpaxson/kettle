# Go parameters
GOCMD=go
GOBUILD=$(GOCMD) build
GOCLEAN=$(GOCMD) clean
GOTEST=$(GOCMD) test
GOMOD=$(GOCMD) mod

# Binary name - default is current directory name
BINARY_NAME := $(shell basename "$(CURDIR)")

# Setup the -ldflags option for go build here, conditionally
# adding the version information.
VERSION := $(shell git describe --tags --always)
COMMIT  := $(shell git rev-parse --short HEAD)
DATE    := $(shell date -u +%Y-%m-%dT%H:%M:%SZ)
VERSION_MODULE := github.com/mpaxson/kettle/src/internal/version
LDFLAGS := -ldflags "-s -w \
    -X main.version=$(VERSION) \
    -X $(VERSION_MODULE).Version=$(VERSION) \
    -X $(VERSION_MODULE).Commit=$(COMMIT) \
    -X $(VERSION_MODULE).BuildDate=$(DATE)"

# Phony targets are not real files
.PHONY: all build run test clean deps install docs lint

# Default target
all: build

# Build the binary
build: deps
	@echo "Building $(BINARY_NAME)..."
	$(GOBUILD) -trimpath -o bin/$(BINARY_NAME) $(LDFLAGS) .
# Run the application
run: build
	@echo "Version: $(VERSION), Commit: $(COMMIT), Date: $(DATE)"
	@echo "Running $(BINARY_NAME)..."
	./bin/$(BINARY_NAME) 

# Run the application
install: build
	@echo "Running $(BINARY_NAME)..."
	./bin/$(BINARY_NAME) install
# Run tests
test:
	@echo "Running tests..."
	$(GOTEST) -v ./...

# Clean up build artifacts
clean:
	@echo "Cleaning..."
	$(GOCLEAN)
	rm -f $(BINARY_NAME)

docs:
	@echo "Generating documentation..."
	rm -r ./docs/cli || true
	go run ./src/internal/tools/docgen -out ./docs/cli -format markdown

format: deps
	@echo "Formatting code..."
	$(GOCMD) fmt ./...

lint: format
	@echo "Linting..."

	golangci-lint run


dev-tools:
	go install github.com/evilmartians/lefthook@latest
	go install github.com/golangci/golangci-lint/cmd/golangci-lint@latest
	
# Install dependencies
deps:
	@echo "Installing dependencies..."
	$(GOMOD) tidy
	$(GOMOD) download
