BINARY_NAME=microsql-ago
BUILD_DIR=build
DOCKER_IMAGE=microsql-ago
VERSION=1.0.0

.PHONY: all build clean test coverage lint fmt docker-build docker-push

all: clean build test

build:
	@echo "Building..."
	@go build -o $(BUILD_DIR)/$(BINARY_NAME) ./cmd/api

clean:
	@echo "Cleaning..."
	@rm -rf $(BUILD_DIR)

test:
	@echo "Running tests..."
	@go test -v ./...

test-integration:
	@echo "Running integration tests..."
	@go test -v ./test/integration/...

test-security:
	@echo "Running security tests..."
	@go test -v ./test/security/...

test-benchmark:
	@echo "Running benchmarks..."
	@go test -v -bench=. -benchmem ./test/benchmarks/...

coverage:
	@echo "Generating coverage report..."
	@go test -coverprofile=coverage.out ./...
	@go tool cover -html=coverage.out

lint:
	@echo "Running linter..."
	@golangci-lint run

fmt:
	@echo "Formatting code..."
	@go fmt ./...

run:
	@echo "Running application..."
	@go run ./cmd/api

dev:
	@echo "Running in development mode..."
	@air

docker-build:
	@echo "Building Docker image..."
	@docker build -t $(DOCKER_IMAGE):$(VERSION) .

docker-push:
	@echo "Pushing Docker image..."
	@docker push $(DOCKER_IMAGE):$(VERSION)

generate-docs:
	@echo "Generating API documentation..."
	@swag init -g cmd/api/main.go

install-tools:
	@echo "Installing development tools..."
	@go install github.com/cosmtrek/air@latest
	@go install github.com/golangci/golangci-lint/cmd/golangci-lint@latest
	@go install github.com/swaggo/swag/cmd/swag@latest

migrate:
	@echo "Running database migrations..."
	@go run ./cmd/migrate

seed:
	@echo "Seeding database..."
	@go run ./cmd/seed