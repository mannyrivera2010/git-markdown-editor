# Makefile for the Git Markdown Editor
#
# This Makefile contains commands for building, running, testing, and cleaning the application.

.PHONY: build run test clean fmt lint

# ==============================================================================
# Variables
# ==============================================================================

APP_NAME := git-markdown-editor
CMD_DIR := .
INTERNAL_DIR := ./internal

# ==============================================================================
# Build and Run
# ==============================================================================

# Build the application
build:
	@echo "Building $(APP_NAME)..."
	@go build -o $(APP_NAME) $(CMD_DIR)

# Run the application
run: 
	@echo "Running $(APP_NAME)..."
	@./$(APP_NAME)

# ==============================================================================
# Testing and Code Quality
# ==============================================================================

# Run the tests
test:
	@echo "Running tests..."
	@go test ../../...

# Format the code
fmt:
	@echo "Formatting code..."
	@go fmt ../../...

# Lint the code
lint:
	@echo "Linting code... (requires golangci-lint)"
	@golangci-lint run ../../... || true # Allow lint to fail without stopping build

# ==============================================================================
# Cleanup
# ==============================================================================

# Clean up the build artifacts
clean:
	@echo "Cleaning up..."
	@rm -f $(APP_NAME)
	@go clean
