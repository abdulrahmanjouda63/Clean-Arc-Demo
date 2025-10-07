# Makefile for Go project with testing and Swagger

.PHONY: test test-verbose test-coverage test-unit test-integration clean deps swagger swagger-install swagger-serve

# Default target
all: deps test

# Install dependencies
deps:
	go mod download
	go mod tidy

# Run all tests
test:
	go test ./...

# Run tests with verbose output
test-verbose:
	go test -v ./...

# Run tests with coverage
test-coverage:
	go test -coverprofile=coverage.out ./...
	go tool cover -html=coverage.out -o coverage.html

# Run only unit tests (excluding integration tests)
test-unit:
	go test -short ./...

# Run integration tests
test-integration:
	go test -run Integration ./...

# Clean test artifacts
clean:
	rm -f coverage.out coverage.html
	go clean -testcache

# Run tests for specific package
test-handlers:
	go test ./handlers/...

test-services:
	go test ./services/...

test-repositories:
	go test ./repositories/...

test-utils:
	go test ./utils/...

# Run tests with race detection
test-race:
	go test -race ./...

# Benchmark tests
benchmark:
	go test -bench=. ./...

# Format code
fmt:
	go fmt ./...

# Lint code
lint:
	golangci-lint run

# Build the application
build:
	go build -o bin/app main.go

PORT ?=
DEFAULT_PORT := :8081
EPORT := $(if $(PORT),$(PORT),$(DEFAULT_PORT))

# Run the application (optional: PORT=":8082")
run:
ifdef PORT
	@echo "Running on custom port $(PORT)"
	@powershell -Command "$env:SERVER_PORT='$(PORT)'; go run main.go"
else
	go run main.go
endif

# Run in release mode (GIN_MODE=release)
run-release:
ifdef PORT
	@echo "Running in release mode on custom port $(PORT)"
	@powershell -Command "$env:GIN_MODE='release'; $env:SERVER_PORT='$(PORT)'; go run main.go"
else
	@powershell -Command "$env:GIN_MODE='release'; go run main.go"
endif

# Stop process listening on PORT (default :8081)
stop:
	@echo "Stopping process on $(EPORT) if running..."
	@powershell -NoProfile -Command "$$p='$(EPORT)'.TrimStart(':'); $$conns = @(Get-NetTCPConnection -LocalPort $$p -State Listen -ErrorAction SilentlyContinue); if ($$conns.Count -gt 0) { $$pids = $$conns \| Select-Object -ExpandProperty OwningProcess -Unique; foreach ($$pid in $$pids) { Write-Host 'Killing PID' $$pid; Stop-Process -Id $$pid -Force -ErrorAction SilentlyContinue } } else { Write-Host 'No process found on' (':'+$$p) }"

# Stop then run on desired port
run-free: stop run

# Swagger commands
swagger-install:
	@echo "Installing Swagger tools..."
	go install github.com/swaggo/swag/cmd/swag@latest

swagger:
	@echo "Generating Swagger documentation..."
	swag init -g main.go -o ./docs
	@echo "✅ Swagger documentation generated!"
	@echo "🌐 Access Swagger UI at: http://localhost:8080/swagger/index.html"

swagger-serve:
	@echo "Starting Swagger UI server..."
	@echo "🌐 Swagger UI will be available at: http://localhost:8080/swagger/index.html"
	@echo "🚀 Opening browser to Swagger UI..."
	@powershell -Command "Start-Sleep 3; Start-Process 'http://localhost:8080/swagger/index.html'" &
	go run main.go

swagger-open:
	@echo "🚀 Opening Swagger UI in browser..."
	@powershell -Command "Start-Process 'http://localhost:8080/swagger/index.html'"

# Help
help:
	@echo "Available targets:"
	@echo "  test          - Run all tests"
	@echo "  test-verbose  - Run tests with verbose output"
	@echo "  test-coverage - Run tests with coverage report"
	@echo "  test-unit     - Run only unit tests"
	@echo "  test-integration - Run integration tests"
	@echo "  test-handlers - Run handler tests"
	@echo "  test-services - Run service tests"
	@echo "  test-repositories - Run repository tests"
	@echo "  test-utils    - Run utility tests"
	@echo "  test-race     - Run tests with race detection"
	@echo "  benchmark     - Run benchmark tests"
	@echo "  clean         - Clean test artifacts"
	@echo "  deps          - Install dependencies"
	@echo "  fmt           - Format code"
	@echo "  lint          - Lint code"
	@echo "  build         - Build the application"
	@echo "  run           - Run the application"
	@echo "      VARS: PORT=':8082' to override port (else :8081)"
	@echo "  run-release   - Run the application in release mode"
	@echo "      VARS: PORT=':8082' to override port (else :8081)"
	@echo "  stop          - Kill process listening on PORT (default :8081)"
	@echo "  run-free      - Stop then run on desired PORT"
	@echo "  swagger-install - Install Swagger tools"
	@echo "  swagger       - Generate Swagger documentation"
	@echo "  swagger-serve - Start server with Swagger UI (auto-opens browser)"
	@echo "  swagger-open - Open Swagger UI in browser"
