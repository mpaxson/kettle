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
VERSION ?= $(shell git describe --tags --always --dirty)
LDFLAGS = -ldflags="-X main.version=${VERSION}"

# Phony targets are not real files
.PHONY: all build run test clean deps install docs lint

# Default target
all: build

# Build the binary
build: deps
	@echo "Building $(BINARY_NAME)..."
	$(GOBUILD) -o bin/$(BINARY_NAME) $(LDFLAGS) .

# Run the application
run: build
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
	go run ./src/internal/tools/docgen -out ./docs/cli -format markdown

lint:
	@echo "Linting..."
	golangci-lint run
	
# Install dependencies
deps:
	@echo "Installing dependencies..."
	$(GOMOD) tidy
	$(GOMOD) download
