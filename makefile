# Build and run the entire application
.PHONY: all
all: build run

# Build all services
.PHONY: build
build:
	docker-compose build

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

# Run tests
.PHONY: test
test:
	go test ./...

# Help command
.PHONY: help
help:
	@echo "Available commands:"
	@echo "  make all          - Build and run the application"
	@echo "  make build        - Build all services"
	@echo "  make run          - Run the application"
	@echo "  make stop         - Stop all services"
	@echo "  make clean        - Clean up docker resources"
	@echo "  make models       - Pull required Ollama models"
	@echo "  make dev-setup    - Install development dependencies"
	@echo "  make dev-frontend - Run frontend in development mode"
	@echo "  make dev-backend  - Run backend in development mode"
	@echo "  make test         - Run tests"