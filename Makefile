.PHONY: build run test clean docker-build docker-run docker-stop lint fmt help

# Variables
BINARY_NAME=download-manager
VERSION?=dev
COMMIT?=$(shell git rev-parse --short HEAD 2>/dev/null || echo "unknown")
BUILD_DIR=bin
GO_VERSION=1.21

# Build flags
LDFLAGS=-ldflags "-s -w -X main.version=$(VERSION) -X main.commit=$(COMMIT)"

help: ## Show this help message
	@echo 'Usage: make [target]'
	@echo ''
	@echo 'Available targets:'
	@awk 'BEGIN {FS = ":.*?## "} /^[a-zA-Z_-]+:.*?## / {printf "  %-15s %s\n", $$1, $$2}' $(MAKEFILE_LIST)

build: ## Build the application for current platform
	@echo "Building $(BINARY_NAME)..."
	@mkdir -p $(BUILD_DIR)
	@go build $(LDFLAGS) -trimpath -o $(BUILD_DIR)/$(BINARY_NAME) ./cmd/main.go
	@echo "Build complete: $(BUILD_DIR)/$(BINARY_NAME)"

build-all: ## Build for all platforms
	@echo "Building for all platforms..."
	@mkdir -p $(BUILD_DIR)
	@GOOS=linux GOARCH=amd64 go build $(LDFLAGS) -trimpath -o $(BUILD_DIR)/$(BINARY_NAME)-linux-amd64 ./cmd/main.go
	@GOOS=linux GOARCH=arm64 go build $(LDFLAGS) -trimpath -o $(BUILD_DIR)/$(BINARY_NAME)-linux-arm64 ./cmd/main.go
	@GOOS=darwin GOARCH=amd64 go build $(LDFLAGS) -trimpath -o $(BUILD_DIR)/$(BINARY_NAME)-darwin-amd64 ./cmd/main.go
	@GOOS=darwin GOARCH=arm64 go build $(LDFLAGS) -trimpath -o $(BUILD_DIR)/$(BINARY_NAME)-darwin-arm64 ./cmd/main.go
	@GOOS=windows GOARCH=amd64 go build $(LDFLAGS) -trimpath -o $(BUILD_DIR)/$(BINARY_NAME)-windows-amd64.exe ./cmd/main.go
	@echo "Build complete: All binaries in $(BUILD_DIR)/"

run: ## Run the application
	@go run ./cmd/main.go

test: ## Run tests
	@echo "Running tests..."
	@go test -v -race -coverprofile=coverage.out ./...
	@go tool cover -html=coverage.out -o coverage.html
	@echo "Coverage report: coverage.html"

lint: ## Run linters
	@echo "Running linters..."
	@golangci-lint run ./...
	@go vet ./...

fmt: ## Format code
	@echo "Formatting code..."
	@gofmt -s -w .
	@go mod tidy

clean: ## Clean build artifacts
	@echo "Cleaning..."
	@rm -rf $(BUILD_DIR)
	@rm -f coverage.out coverage.html
	@go clean -cache -modcache -testcache

docker-build: ## Build Docker image
	@echo "Building Docker image..."
	@docker build \
		--build-arg VERSION=$(VERSION) \
		--build-arg COMMIT=$(COMMIT) \
		-t $(BINARY_NAME):$(VERSION) \
		-t $(BINARY_NAME):latest \
		.

docker-run: ## Run Docker container
	@echo "Running Docker container..."
	@docker-compose up -d

docker-stop: ## Stop Docker container
	@echo "Stopping Docker container..."
	@docker-compose down

docker-logs: ## View Docker logs
	@docker-compose logs -f download-manager

install: build ## Build and install to system
	@echo "Installing $(BINARY_NAME)..."
	@sudo cp $(BUILD_DIR)/$(BINARY_NAME) /usr/local/bin/
	@echo "Installed to /usr/local/bin/$(BINARY_NAME)"

deps: ## Download dependencies
	@echo "Downloading dependencies..."
	@go mod download
	@go mod verify

update-deps: ## Update dependencies
	@echo "Updating dependencies..."
	@go get -u ./...
	@go mod tidy

check: fmt lint test ## Run all checks (format, lint, test)

ci: check build-all ## Run CI pipeline locally

.DEFAULT_GOAL := help

