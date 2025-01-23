# Build and run the entire application
.PHONY: all
all: deps build test lint

# Build all services and CLI
.PHONY: build
build:
	docker-compose build
	go build -o bin/thresh cmd/cli/main.go

# Run the application
.PHONY: run
run:
	docker-compose up

# Stop all services
.PHONY: stop
stop:
	docker-compose down

# Clean up docker resources
.PHONY: clean
clean:
	docker-compose down -v
	docker system prune -f

# Pull required Ollama models
.PHONY: models
models:
	docker-compose run --rm ollama pull nous-hermes2:10.7b-ctx
	docker-compose run --rm ollama pull codellama:70b

# Development setup
.PHONY: dev-setup
dev-setup:
	cd frontend && npm install

# Run frontend in development mode
.PHONY: dev-frontend
dev-frontend:
	cd frontend && npm start

# Run backend in development mode
.PHONY: dev-backend
dev-backend:
	go run cmd/web/main.go
# Install Dependencies
.PHONY: deps
deps:
	go mod download

# Run tests
.PHONY: test
test:
	go test -cover ./...

# Run CLI
.PHONY: run-cli
run-cli:
	go run cmd/cli/main.go -provider deepseek -prompt "Hello world"

# Lint
.PHONY: lint
lint:
	golangci-lint run

# Help command
# Help command
.PHONY: help
help:
	@echo "Available commands:"
	@echo "  make all          - Run deps, build, test, and lint"
	@echo "  make build        - Build all services and CLI binary"
	@echo "  make run          - Run the application"
	@echo "  make run-cli      - Run the CLI with deepseek provider"
	@echo "  make stop         - Stop all services"
	@echo "  make clean        - Clean up docker resources"
	@echo "  make models       - Pull required Ollama models"
	@echo "  make deps         - Install Go dependencies"
	@echo "  make dev-setup    - Install development dependencies"
	@echo "  make dev-frontend - Run frontend in development mode"
	@echo "  make dev-backend  - Run backend in development mode"
	@echo "  make test         - Run tests with coverage"
	@echo "  make lint         - Run golangci-lint"