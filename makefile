# Version information
VERSION := $(shell git describe --tags --always --dirty 2>/dev/null || echo "dev")
COMMIT := $(shell git rev-parse --short HEAD 2>/dev/null || echo "unknown")
BUILD_TIME := $(shell date -u '+%Y-%m-%d_%H:%M:%S')
GO_VERSION := $(shell go version | cut -d' ' -f3)

# Build settings
BINARY_NAME := thresh
BUILD_DIR := bin
DIST_DIR := dist
GO_FILES := $(shell find . -type f -name '*.go' -not -path "./vendor/*")
BUILD_FLAGS := -ldflags="-s -w -X main.version=$(VERSION) -X main.commit=$(COMMIT) -X main.buildTime=$(BUILD_TIME)"
COVER_PROFILE := coverage.out
COVER_HTML := coverage.html

# Cross-compilation settings
PLATFORMS := linux/amd64 darwin/amd64 windows/amd64

# Build and run the entire application
.PHONY: all
all: verify-prereqs deps security-check build test lint

# Verify prerequisites
.PHONY: verify-prereqs
verify-prereqs:
	@echo "Verifying prerequisites..."
	@which go >/dev/null || (echo "Error: Go is not installed" && exit 1)
	@which docker >/dev/null || (echo "Error: Docker is not installed" && exit 1)
	@which golangci-lint >/dev/null || (echo "Error: golangci-lint is not installed" && exit 1)
	@go version | grep -q "go1.21" || (echo "Error: Project requires Go 1.21" && exit 1)

# Security check
.PHONY: security-check
security-check:
	@echo "Running security checks..."
	@go install golang.org/x/vuln/cmd/govulncheck@latest
	@govulncheck ./...

# Build all services and CLI
.PHONY: build
build: $(BUILD_DIR)
	@echo "Building $(BINARY_NAME) version $(VERSION)"
	@go build $(BUILD_FLAGS) -o $(BUILD_DIR)/$(BINARY_NAME) cmd/cli/main.go
	@docker-compose build

# Cross-platform builds
.PHONY: release
release: $(DIST_DIR)
	@echo "Building for multiple platforms..."
	@for platform in $(PLATFORMS); do \
		GOOS=$${platform%/*} GOARCH=$${platform#*/} \
		go build $(BUILD_FLAGS) -o $(DIST_DIR)/$(BINARY_NAME)-$${platform%/*}-$${platform#*/}$(if $(findstring windows,$${platform%/*}),.exe,) \
		cmd/cli/main.go || exit 1; \
		echo "Built: $(DIST_DIR)/$(BINARY_NAME)-$${platform%/*}-$${platform#*/}"; \
	done
	@cd $(DIST_DIR) && for file in *; do sha256sum "$$file" >> checksums.txt; done

# Install Dependencies
.PHONY: deps
deps:
	@echo "Installing dependencies..."
	@go mod download
	@go mod verify

# Run tests with coverage
.PHONY: test
test:
	@echo "Running tests with coverage..."
	@go test -race -coverprofile=$(COVER_PROFILE) -covermode=atomic ./...
	@go tool cover -html=$(COVER_PROFILE) -o $(COVER_HTML)
	@echo "Coverage report generated at $(COVER_HTML)"
	@go tool cover -func=$(COVER_PROFILE)

# Lint
.PHONY: lint
lint:
	@echo "Running linter..."
	@golangci-lint run

# Run the application
.PHONY: run
run:
	@docker-compose up

# Stop all services
.PHONY: stop
stop:
	@docker-compose down

# Clean up
.PHONY: clean
clean:
	@echo "Cleaning up..."
	@rm -rf $(BUILD_DIR) $(DIST_DIR) $(COVER_PROFILE) $(COVER_HTML)
	@docker-compose down -v
	@docker system prune -f

# Pull required Ollama models
.PHONY: models
models:
	@docker-compose run --rm ollama pull nous-hermes2:10.7b-ctx
	@docker-compose run --rm ollama pull codellama:70b

# Development setup
.PHONY: dev-setup
dev-setup:
	@cd frontend && npm install

# Run frontend in development mode
.PHONY: dev-frontend
dev-frontend:
	@cd frontend && npm start

# Run backend in development mode
.PHONY: dev-backend
dev-backend:
	@go run cmd/web/main.go

# Run CLI
.PHONY: run-cli
run-cli:
	@go run cmd/cli/main.go -provider deepseek -prompt "Hello world"

# Installation verification
.PHONY: verify-install
verify-install: $(BUILD_DIR)/$(BINARY_NAME)
	@echo "Verifying installation..."
	@./$(BUILD_DIR)/$(BINARY_NAME) --version
	@echo "Checking Docker services..."
	@docker-compose config --quiet
	@echo "Installation verified successfully"

# Create necessary directories
$(BUILD_DIR):
	@mkdir -p $(BUILD_DIR)

$(DIST_DIR):
	@mkdir -p $(DIST_DIR)

# Help command
.PHONY: help
help:
	@echo "ThreshAI Build System $(VERSION)"
	@echo "Available commands:"
	@echo "  make all            - Run full verification, build and test suite"
	@echo "  make build          - Build all services and CLI binary"
	@echo "  make release        - Build cross-platform binaries"
	@echo "  make run           - Run the application"
	@echo "  make run-cli       - Run the CLI with deepseek provider"
	@echo "  make stop          - Stop all services"
	@echo "  make clean         - Clean up build artifacts and docker resources"
	@echo "  make models        - Pull required Ollama models"
	@echo "  make deps          - Install and verify Go dependencies"
	@echo "  make dev-setup     - Install development dependencies"
	@echo "  make dev-frontend  - Run frontend in development mode"
	@echo "  make dev-backend   - Run backend in development mode"
	@echo "  make test          - Run tests with coverage report"
	@echo "  make lint          - Run golangci-lint"
	@echo "  make security-check - Run vulnerability scanning"
	@echo "  make verify-install - Verify installation and dependencies"

.DEFAULT_GOAL := help