# Define variables
BINARY_NAME := nrms
VERSION := 1.0.1
BUILD_DIR := build
BIN_DIR := bin
CMD_DIR := cmd/newrelic/metric/selector

.PHONY: all clean build lint test build-linux build-mac package-linux package-mac

# Default target
all: build

# Clean build directory
clean:
	@echo "Cleaning build directory..."
	@rm -rf $(BUILD_DIR)/*
	@rm -rf $(BIN_DIR)/*

# Build the nrms binary for the current platform
build:
	@echo "Building $(BINARY_NAME) for current platform..."
	@mkdir -p $(BIN_DIR)
	GO111MODULE=on go build -o $(BIN_DIR)/$(BINARY_NAME) ./$(CMD_DIR)

# Build the nrms binary for Linux
build-linux:
	@echo "Building $(BINARY_NAME) for Linux..."
	@mkdir -p $(BIN_DIR)
	GO111MODULE=on GOOS=linux GOARCH=amd64 go build -o $(BIN_DIR)/$(BINARY_NAME) ./$(CMD_DIR)

# Build the nrms binary for macOS
build-mac:
	@echo "Building $(BINARY_NAME) for macOS..."
	@mkdir -p $(BIN_DIR)
	GO111MODULE=on GOOS=darwin GOARCH=amd64 go build -o $(BIN_DIR)/$(BINARY_NAME) ./$(CMD_DIR)

# Package the nrms binary for Linux
package-linux: build-linux
	@echo "Packaging $(BINARY_NAME) for Linux..."
	@mkdir -p $(BUILD_DIR)
	@tar -czvf $(BUILD_DIR)/$(BINARY_NAME)-linux-$(VERSION).tar.gz -C $(BIN_DIR) $(BINARY_NAME)

# Package the nrms binary for macOS
package-mac: build-mac
	@echo "Packaging $(BINARY_NAME) for macOS..."
	@mkdir -p $(BUILD_DIR)
	@tar -czvf $(BUILD_DIR)/$(BINARY_NAME)-mac-$(VERSION).tar.gz -C $(BIN_DIR) $(BINARY_NAME)

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
	@echo "  package-linux - Package the nrms binary for Linux"
	@echo "  package-mac  - Package the nrms binary for macOS"
	@echo "  lint         - Lint the code"
	@echo "  test         - Run tests"
	@echo "  deps         - Install dependencies"
	@echo "  lint-commit  - Check commit message format"
	@echo "  tools        - Run all tools (deps, lint, test, lint-commit, build)"
	@echo "  help         - Show this help message"
