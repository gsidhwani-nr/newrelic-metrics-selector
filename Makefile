# Define variables
BINARY_NAME := nrms
BUILD_DIR := build
CMD_DIR := cmd/newrelic/metric/selector

.PHONY: all clean build lint test build-linux build-mac

# Default target
all: build

# Clean build directory
clean:
	@echo "Cleaning build directory..."
	@rm -rf $(BUILD_DIR)/*
	@rm -rf bin/*

# Build the nrms binary for the current platform
build:
	@echo "Building $(BINARY_NAME) for current platform..."
	@mkdir -p bin
	@cd $(CMD_DIR) && GO111MODULE=on go build -o ../../../../bin/$(BINARY_NAME)

# Build the nrms binary for Linux
build-linux:
	@echo "Building $(BINARY_NAME) for Linux..."
	@mkdir -p bin
	@cd $(CMD_DIR) && GO111MODULE=on GOOS=linux GOARCH=amd64 go build -o ../../../../bin/$(BINARY_NAME)-linux

# Build the nrms binary for macOS
build-mac:
	@echo "Building $(BINARY_NAME) for macOS..."
	@mkdir -p bin
	@cd $(CMD_DIR) && GO111MODULE=on GOOS=darwin GOARCH=amd64 go build -o ../../../../bin/$(BINARY_NAME)-mac

# Lint the code
lint:
	@echo "Linting the code..."
	@golangci-lint run

# Run tests
test:
	@echo "Running tests..."
	@go test ./...

# Install dependencies
deps:
	@echo "Installing dependencies..."
	@go mod tidy

# Check commit message format
lint-commit:
	@echo "Checking commit message format..."
	@git log -1 --pretty=%B | grep -P '^(chore|docs|feat|fix|refactor|test|tests?)\s?(\([^\)]+\))?!?: .+$$' || (echo "Commit message format is incorrect" && exit 1)

# Default target for tools
tools: deps lint test lint-commit build

# Help message
help:
	@echo "Usage: make [target]"
	@echo ""
	@echo "Targets:"
	@echo "  all          - Build the project (default)"
	@echo "  clean        - Clean the build directory"
	@echo "  build        - Build the nrms binary for the current platform"
	@echo "  build-linux  - Build the nrms binary for Linux"
	@echo "  build-mac    - Build the nrms binary for macOS"
	@echo "  lint         - Lint the code"
	@echo "  test         - Run tests"
	@echo "  deps         - Install dependencies"
	@echo "  lint-commit  - Check commit message format"
	@echo "  tools        - Run all tools (deps, lint, test, lint-commit, build)"
	@echo "  help         - Show this help message"
