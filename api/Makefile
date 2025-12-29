.PHONY: test test-coverage test-watch test-verbose test-race test-bench test-integration clean

# Default target
all: test

# Run all tests
test:
	@echo "Running tests..."
	@go test ./... -v

# Run tests with coverage
test-coverage:
	@echo "Running tests with coverage..."
	@go test ./... -coverprofile=coverage.out -covermode=atomic
	@go tool cover -func=coverage.out | grep total
	@echo "Coverage report generated: coverage.out"

# Run tests and generate HTML coverage report
test-coverage-html: test-coverage
	@go tool cover -html=coverage.out -o coverage.html
	@echo "HTML coverage report generated: coverage.html"
	@start coverage.html

# Run tests in watch mode
test-watch:
	@echo "Running tests in watch mode..."
	@which gow >/dev/null 2>&1 || go install github.com/mitranim/gow@latest
	@gow test ./... -v

# Run tests with verbose output
test-verbose:
	@echo "Running tests with verbose output..."
	@go test ./... -v -count=1

# Run tests with race detector
test-race:
	@echo "Running tests with race detector..."
	@go test ./... -race -count=1

# Run benchmarks
test-bench:
	@echo "Running benchmarks..."
	@go test ./... -bench=. -benchmem

# Run integration tests (if any)
test-integration:
	@echo "Running integration tests..."
	@go test ./... -tags=integration -v

# Run integration tests with test database environment
test-integration-db:
	@echo "Running integration tests with test database..."
	@TEST_DB_HOST=localhost \
	TEST_DB_PORT=5433 \
	TEST_DB_USER=gamelink \
	TEST_DB_PASSWORD=gamelink \
	TEST_DB_NAME=gamelink_test \
	go test ./internal/service/integration/... -v -count=1

# Run tests for a specific package
run-test:
	@echo "Running tests for $(PKG)..."
	@go test ./$(PKG)/... -v

# Run coverage for a specific package
cover-pkg:
	@echo "Running coverage for $(PKG)..."
	@go test ./$(PKG)/... -coverprofile=coverage.out -covermode=atomic
	@go tool cover -func=coverage.out

# Clean test artifacts
clean:
	@echo "Cleaning test artifacts..."
	@rm -f coverage.out coverage.html
	@rm -f */*.test
	@rm -f */*/*.test
	@rm -f */*/*/*.test

# Generate mocks
generate-mocks:
	@echo "Generating mocks..."
	@which mockery >/dev/null 2>&1 || go install github.com/vektra/mockery/v2@latest
	@mockery --all --output=mocks --outpkg=mocks --with-expecter

# Run linter
lint:
	@echo "Running linter..."
	@which golangci-lint >/dev/null 2>&1 || go install github.com/golangci/golangci-lint/cmd/golangci-lint@latest
	@golangci-lint run ./...

# Format code
fmt:
	@echo "Formatting code..."
	@go fmt ./...

# Vet code
vet:
	@echo "Vetting code..."
	@go vet ./...

# Run all checks (test, lint, fmt, vet)
check: fmt vet lint test
	@echo "All checks passed!"

# Install dependencies
deps:
	@echo "Installing dependencies..."
	@go mod download
	@go mod tidy

# Install test tools
test-tools:
	@echo "Installing test tools..."
	@go install github.com/stretchr/testify@latest
	@go install github.com/golang/mock/mockgen@latest
	@go install github.com/vektra/mockery/v2@latest
	@go install github.com/mitranim/gow@latest
	@go install github.com/golangci/golangci-lint/cmd/golangci-lint@latest

# Generate Swagger documentation (supports generics)
swagger:
	@echo "Generating Swagger documentation (with generics support)..."
	@which swag >/dev/null 2>&1 || go install github.com/swaggo/swag/cmd/swag@latest
	@swag init \
		-g cmd/main.go \
		--output "docs" \
		--dir "./" \
		--parseDependency \
		--parseInternal \
		--parseDepth 10
	@echo "Swagger docs generated in docs/"

# Generate Swagger documentation with vendor
swagger-vendor:
	@echo "Generating Swagger documentation (with vendor support)..."
	@which swag >/dev/null 2>&1 || go install github.com/swaggo/swag/cmd/swag@latest
	@swag init \
		-g cmd/main.go \
		--output "docs" \
		--dir "./" \
		--parseDependency \
		--parseInternal \
		--parseDepth 10 \
		--parseVendor
	@echo "Swagger docs with vendor support generated in docs/"

# Clean Swagger docs
clean-swagger:
	@echo "Cleaning Swagger documentation..."
	@rm -rf docs/swagger.json docs/swagger.yaml docs/docs.go docs/swagger/*
	@echo "Swagger docs cleaned"

# Help
help:
	@echo "Available targets:"
	@echo "  test              Run all tests"
	@echo "  test-coverage     Run tests with coverage"
	@echo "  test-coverage-html Run tests and generate HTML coverage report"
	@echo "  test-watch        Run tests in watch mode"
	@echo "  test-verbose      Run tests with verbose output"
	@echo "  test-race         Run tests with race detector"
	@echo "  test-bench        Run benchmarks"
	@echo "  test-integration  Run integration tests"
	@echo "  test-integration-db Run integration tests with test DB (auto-sets env vars)"
	@echo "  run-test PKG=pkg Run tests for specific package"
	@echo "  cover-pkg PKG=pkg Run coverage for specific package"
	@echo "  clean             Clean test artifacts"
	@echo "  generate-mocks    Generate mocks"
	@echo "  swagger           Generate Swagger documentation (supports generics)"
	@echo "  swagger-vendor    Generate Swagger documentation with vendor support"
	@echo "  clean-swagger     Clean Swagger documentation"
	@echo "  lint              Run linter"
	@echo "  fmt               Format code"
	@echo "  vet               Vet code"
	@echo "  check             Run all checks"
	@echo "  deps              Install dependencies"
	@echo "  test-tools        Install test tools"
	@echo "  help              Show this help"
