# Makefile for forter

.PHONY: build clean test lint install uninstall fmt vet run help

# Variables
BINARY_NAME=forter
MAIN_PATH=cmd/forter/main.go
INSTALL_PATH=$(GOPATH)/bin
VERSION?=dev
COMMIT=$(shell git rev-parse --short HEAD 2>/dev/null || echo "unknown")
BUILD_DATE=$(shell date -u +%Y-%m-%dT%H:%M:%SZ)
LDFLAGS=-ldflags "-X main.Version=$(VERSION) -X main.Commit=$(COMMIT) -X main.BuildDate=$(BUILD_DATE)"

# Colors for output
BLUE=\033[34m
GREEN=\033[32m
YELLOW=\033[33m
RED=\033[31m
NC=\033[0m

## help: Show this help message
help:
	@echo "Usage: make [target]"
	@echo ""
	@echo "Targets:"
	@awk '/^[a-zA-Z_-]+:/ { \
		helpMessage = match(lastLine, /^## (.*)/); \
		if (helpMessage) { \
			target = $$1; sub(/:$$/, "", target); \
			printf "  $(BLUE)%-15s$(NC) %s\n", target, substr(lastLine, RSTART + 3, RLENGTH); \
		} \
	} { lastLine = $$0 }' $(MAKEFILE_LIST)

## build: Build the binary
build:
	@echo "$(BLUE)Building $(BINARY_NAME)...$(NC)"
	@mkdir -p dist
	go build $(LDFLAGS) -o dist/$(BINARY_NAME) $(MAIN_PATH)
	@echo "$(GREEN)Build complete: dist/$(BINARY_NAME)$(NC)"

## build-all: Build for multiple platforms
build-all: build-darwin build-linux

build-darwin:
	@echo "$(BLUE)Building for macOS...$(NC)"
	@mkdir -p dist
	GOOS=darwin GOARCH=amd64 go build $(LDFLAGS) -o dist/$(BINARY_NAME)-darwin-amd64 $(MAIN_PATH)
	GOOS=darwin GOARCH=arm64 go build $(LDFLAGS) -o dist/$(BINARY_NAME)-darwin-arm64 $(MAIN_PATH)
	@echo "$(GREEN)macOS builds complete$(NC)"

build-linux:
	@echo "$(BLUE)Building for Linux...$(NC)"
	@mkdir -p dist
	GOOS=linux GOARCH=amd64 go build $(LDFLAGS) -o dist/$(BINARY_NAME)-linux-amd64 $(MAIN_PATH)
	GOOS=linux GOARCH=arm64 go build $(LDFLAGS) -o dist/$(BINARY_NAME)-linux-arm64 $(MAIN_PATH)
	@echo "$(GREEN)Linux builds complete$(NC)"

## clean: Clean build artifacts
clean:
	@echo "$(YELLOW)Cleaning build artifacts...$(NC)"
	@rm -rf dist/
	@go clean
	@echo "$(GREEN)Clean complete$(NC)"

## test: Run all tests
test:
	@echo "$(BLUE)Running tests...$(NC)"
	@go test -v ./...

## test-coverage: Run tests with coverage report
test-coverage:
	@echo "$(BLUE)Running tests with coverage...$(NC)"
	@go test -coverprofile=coverage.out ./...
	@go tool cover -html=coverage.out -o coverage.html
	@echo "$(GREEN)Coverage report generated: coverage.html$(NC)"

## lint: Run linter
lint:
	@echo "$(BLUE)Running linter...$(NC)"
	@if command -v golangci-lint >/dev/null 2>&1; then \
		golangci-lint run; \
	else \
		echo "$(YELLOW)golangci-lint not installed, using go vet$(NC)"; \
		go vet ./...; \
	fi

## fmt: Format code
fmt:
	@echo "$(BLUE)Formatting code...$(NC)"
	@go fmt ./...

## vet: Run go vet
vet:
	@echo "$(BLUE)Running go vet...$(NC)"
	@go vet ./...

## install: Install binary to GOPATH/bin
install: build
	@echo "$(BLUE)Installing to $(INSTALL_PATH)...$(NC)"
	@cp dist/$(BINARY_NAME) $(INSTALL_PATH)/$(BINARY_NAME)
	@echo "$(GREEN)Installed to $(INSTALL_PATH)/$(BINARY_NAME)$(NC)"

## uninstall: Remove binary from GOPATH/bin
uninstall:
	@echo "$(YELLOW)Removing from $(INSTALL_PATH)...$(NC)"
	@rm -f $(INSTALL_PATH)/$(BINARY_NAME)
	@echo "$(GREEN)Uninstalled$(NC)"

## run: Run the application (specify ARGS)
run:
	@go run $(MAIN_PATH) $(ARGS)

## run-dry: Run with dry-run mode
run-dry:
	@go run $(MAIN_PATH) --dry-run $(ARGS)

## deps: Download and verify dependencies
deps:
	@echo "$(BLUE)Downloading dependencies...$(NC)"
	@go mod download
	@go mod verify

## update-deps: Update all dependencies
update-deps:
	@echo "$(BLUE)Updating dependencies...$(NC)"
	@go get -u ./...
	@go mod tidy

## tidy: Tidy go modules
tidy:
	@echo "$(BLUE)Tidying modules...$(NC)"
	@go mod tidy

## check: Run all checks (fmt, vet, lint, test)
check: fmt vet lint test
	@echo "$(GREEN)All checks passed!$(NC)"

.DEFAULT_GOAL := build
