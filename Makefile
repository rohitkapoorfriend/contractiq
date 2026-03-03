.PHONY: build run test lint clean migrate-up migrate-down docker-up docker-down

APP_NAME := contractiq
BUILD_DIR := bin

# Build
build:
	@echo "Building $(APP_NAME)..."
	@go build -o $(BUILD_DIR)/$(APP_NAME) ./cmd/api

run: build
	@./$(BUILD_DIR)/$(APP_NAME)

# Development
dev:
	@go run ./cmd/api

# Testing
test:
	@go test -v -race -count=1 ./...

test-cover:
	@go test -v -race -coverprofile=coverage.out ./...
	@go tool cover -html=coverage.out -o coverage.html
	@echo "Coverage report: coverage.html"

# Linting
lint:
	@golangci-lint run ./...

vet:
	@go vet ./...

# Database migrations
migrate-up:
	@migrate -path migrations -database "postgres://$(DB_USER):$(DB_PASSWORD)@$(DB_HOST):$(DB_PORT)/$(DB_NAME)?sslmode=$(DB_SSL_MODE)" up

migrate-down:
	@migrate -path migrations -database "postgres://$(DB_USER):$(DB_PASSWORD)@$(DB_HOST):$(DB_PORT)/$(DB_NAME)?sslmode=$(DB_SSL_MODE)" down 1

migrate-create:
	@migrate create -ext sql -dir migrations -seq $(name)

# Docker
docker-up:
	@docker compose up -d --build

docker-down:
	@docker compose down -v

# Cleanup
clean:
	@rm -rf $(BUILD_DIR) coverage.out coverage.html
