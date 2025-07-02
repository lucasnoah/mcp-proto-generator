.PHONY: build test clean install generate lint

# Build the CLI binary
build:
	@echo "Building mcp-proto-gen..."
	@go build -o bin/mcp-proto-gen ./cmd/mcp-proto-gen

# Install the CLI
install:
	@echo "Installing mcp-proto-gen..."
	@go install ./cmd/mcp-proto-gen

# Run all tests
test: test-unit test-integration

# Run unit tests only
test-unit:
	@echo "Running unit tests..."
	@go test -v -short ./pkg/...

# Run integration tests
test-integration:
	@echo "Running integration tests..."
	@go test -v ./test/...

# Test with coverage
test-coverage:
	@echo "Running tests with coverage..."
	@go test -v -coverprofile=coverage.out -covermode=atomic ./...
	@go tool cover -html=coverage.out -o coverage.html
	@echo "Coverage report: coverage.html"
	@go tool cover -func=coverage.out | grep total | awk '{print "Total coverage: " $$3}'

# Run benchmarks
test-benchmark:
	@echo "Running benchmarks..."
	@go test -bench=. -benchmem ./test/...

# Lint code
lint:
	@echo "Running linter..."
	@golangci-lint run

# Clean build artifacts
clean:
	@echo "Cleaning..."
	@rm -rf bin/
	@rm -f coverage.out coverage.html

# Generate example MCP server (for testing)
generate-example:
	@echo "Generating example MCP server..."
	@./bin/mcp-proto-gen generate --config examples/simple/mcp-gen.yaml

# Build and test end-to-end
e2e: build
	@echo "Running end-to-end test..."
	@cd examples/simple && ../../bin/mcp-proto-gen generate
	@cd examples/simple/mcp-server && go build .
	@echo "✓ E2E test passed"

# Development setup
dev:
	@echo "Setting up development environment..."
	@go mod download
	@go install github.com/golangci/golangci-lint/cmd/golangci-lint@latest

# Release build with cross-compilation
release:
	@echo "Building release binaries..."
	@mkdir -p bin/
	@GOOS=linux GOARCH=amd64 go build -o bin/mcp-proto-gen-linux-amd64 ./cmd/mcp-proto-gen
	@GOOS=darwin GOARCH=amd64 go build -o bin/mcp-proto-gen-darwin-amd64 ./cmd/mcp-proto-gen
	@GOOS=darwin GOARCH=arm64 go build -o bin/mcp-proto-gen-darwin-arm64 ./cmd/mcp-proto-gen
	@GOOS=windows GOARCH=amd64 go build -o bin/mcp-proto-gen-windows-amd64.exe ./cmd/mcp-proto-gen
	@echo "Release binaries in bin/"

# Quick help
help:
	@echo "Available targets:"
	@echo "  build           Build the CLI binary"
	@echo "  install         Install the CLI globally"
	@echo "  test            Run tests"
	@echo "  test-coverage   Run tests with coverage"
	@echo "  lint            Run linter"
	@echo "  clean           Clean build artifacts"
	@echo "  generate-example Generate example MCP server"
	@echo "  e2e             End-to-end test"
	@echo "  dev             Setup development environment"
	@echo "  release         Build release binaries"