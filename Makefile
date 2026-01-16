.PHONY: build test clean run install dev lint docker-build docker-run sdk-build sdk-test all

BINARY=sqm
VERSION=0.1.0
BUILD_DIR=./build
LDFLAGS=-ldflags "-X main.version=$(VERSION)"

# Default target
all: build

# Build the binary
build:
	@echo "Building $(BINARY)..."
	@mkdir -p $(BUILD_DIR)
	go build $(LDFLAGS) -o $(BUILD_DIR)/$(BINARY) ./cmd/sqm

# Build for multiple platforms
build-all: build-linux build-darwin build-windows

build-linux:
	@echo "Building for Linux..."
	GOOS=linux GOARCH=amd64 go build $(LDFLAGS) -o $(BUILD_DIR)/$(BINARY)-linux-amd64 ./cmd/sqm
	GOOS=linux GOARCH=arm64 go build $(LDFLAGS) -o $(BUILD_DIR)/$(BINARY)-linux-arm64 ./cmd/sqm

build-darwin:
	@echo "Building for macOS..."
	GOOS=darwin GOARCH=amd64 go build $(LDFLAGS) -o $(BUILD_DIR)/$(BINARY)-darwin-amd64 ./cmd/sqm
	GOOS=darwin GOARCH=arm64 go build $(LDFLAGS) -o $(BUILD_DIR)/$(BINARY)-darwin-arm64 ./cmd/sqm

build-windows:
	@echo "Building for Windows..."
	GOOS=windows GOARCH=amd64 go build $(LDFLAGS) -o $(BUILD_DIR)/$(BINARY)-windows-amd64.exe ./cmd/sqm

# Run tests
test:
	@echo "Running tests..."
	go test -v ./...

# Run tests with coverage
test-coverage:
	@echo "Running tests with coverage..."
	go test -v -coverprofile=coverage.out ./...
	go tool cover -html=coverage.out -o coverage.html

# Clean build artifacts
clean:
	@echo "Cleaning..."
	rm -rf $(BUILD_DIR)
	rm -f coverage.out coverage.html

# Run the binary
run: build
	$(BUILD_DIR)/$(BINARY)

# Install globally
install: build
	@echo "Installing $(BINARY)..."
	cp $(BUILD_DIR)/$(BINARY) /usr/local/bin/

# Uninstall
uninstall:
	@echo "Uninstalling $(BINARY)..."
	rm -f /usr/local/bin/$(BINARY)

# Development mode - run without building
dev:
	go run ./cmd/sqm

# Lint the code
lint:
	@echo "Linting..."
	golangci-lint run

# Format code
fmt:
	@echo "Formatting..."
	go fmt ./...

# Tidy dependencies
tidy:
	@echo "Tidying dependencies..."
	go mod tidy

# Download dependencies
deps:
	@echo "Downloading dependencies..."
	go mod download

# Docker build
docker-build:
	@echo "Building Docker image..."
	docker build -t squaremind:$(VERSION) .
	docker tag squaremind:$(VERSION) squaremind:latest

# Docker run
docker-run:
	docker run -it --rm squaremind:$(VERSION)

# SDK build (TypeScript)
sdk-build:
	@echo "Building TypeScript SDK..."
	cd sdk/squaremind-sdk && npm install && npm run build

# SDK test
sdk-test:
	@echo "Testing TypeScript SDK..."
	cd sdk/squaremind-sdk && npm test

# SDK publish (dry run)
sdk-publish-dry:
	cd sdk/squaremind-sdk && npm publish --dry-run

# Generate API documentation
docs:
	@echo "Generating documentation..."
	godoc -http=:6060

# Quick start demo
demo: build
	@echo "Running demo..."
	$(BUILD_DIR)/$(BINARY) init DemoSwarm
	$(BUILD_DIR)/$(BINARY) spawn Coder -c code.write,code.review
	$(BUILD_DIR)/$(BINARY) spawn Reviewer -c code.review,security
	$(BUILD_DIR)/$(BINARY) status

# Help
help:
	@echo "Squaremind Makefile"
	@echo ""
	@echo "Usage:"
	@echo "  make build         Build the CLI binary"
	@echo "  make test          Run tests"
	@echo "  make clean         Clean build artifacts"
	@echo "  make install       Install binary to /usr/local/bin"
	@echo "  make dev           Run in development mode"
	@echo "  make lint          Lint the code"
	@echo "  make docker-build  Build Docker image"
	@echo "  make sdk-build     Build TypeScript SDK"
	@echo "  make demo          Run a quick demo"
	@echo "  make help          Show this help"
