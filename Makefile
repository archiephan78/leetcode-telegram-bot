.PHONY: build run clean test docker-build docker-run help

# Variables
BINARY_NAME=leetcode-telegram-bot
DOCKER_IMAGE=leetcode-telegram-bot:latest

# Build the application
build:
	@echo "Building $(BINARY_NAME)..."
	@go build -o $(BINARY_NAME) .

# Run the application
run: build
	@echo "Running $(BINARY_NAME)..."
	@./$(BINARY_NAME)

# Clean build artifacts
clean:
	@echo "Cleaning..."
	@go clean
	@rm -f $(BINARY_NAME)
	@rm -f *.db

# Run tests manually
test:
	@echo "Please run tests manually with: go test -v ./..."

# Run tests with race detector
test-race:
	@echo "Running tests with race detector..."
	@go test -race -v ./...

# Download dependencies
deps:
	@echo "Downloading dependencies..."
	@go mod download
	@go mod tidy

# Build Docker image
docker-build:
	@echo "Building Docker image..."
	@docker build -t $(DOCKER_IMAGE) .

# Run with Docker Compose
docker-up:
	@echo "Starting with Docker Compose..."
	@docker-compose up -d

# Stop Docker Compose
docker-down:
	@echo "Stopping Docker Compose..."
	@docker-compose down

# View Docker logs
docker-logs:
	@echo "Viewing Docker logs..."
	@docker-compose logs -f

# Format code
fmt:
	@echo "Formatting code..."
	@go fmt ./...

# Lint code
lint:
	@echo "Linting code..."
	@golangci-lint run

# Install development dependencies
dev-deps:
	@echo "Installing development dependencies..."
	@go install github.com/golangci/golangci-lint/cmd/golangci-lint@latest

# Show help
help:
	@echo "Available commands:"
	@echo "  build       - Build the application"
	@echo "  run         - Build and run the application"
	@echo "  clean       - Clean build artifacts"
	@echo "  test        - Run tests"
	@echo "  test-race   - Run tests with race detector"
	@echo "  deps        - Download and tidy dependencies"
	@echo "  docker-build - Build Docker image"
	@echo "  docker-up   - Start with Docker Compose"
	@echo "  docker-down - Stop Docker Compose"
	@echo "  docker-logs - View Docker logs"
	@echo "  fmt         - Format code"
	@echo "  lint        - Lint code"
	@echo "  dev-deps    - Install development dependencies"
	@echo "  help        - Show this help message" 