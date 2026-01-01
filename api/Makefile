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
	@go test ./... -bench=. -benchmem -benchtime=10s

# Run benchmarks with specific suite (order, auth, payment, database, http, all)
bench-suite:
	@echo "Running benchmarks for suite $(SUITE)..."
	@bash ./tests/load/run_bench.sh $(SUITE) $(BENCHTIME)

# Run order benchmarks
bench-order:
	@echo "Running order service benchmarks..."
	@go test -v -bench=. -benchtime=10s -run=^$$ ./internal/service/order/... -benchmem

# Run auth benchmarks
bench-auth:
	@echo "Running auth service benchmarks..."
	@go test -v -bench=. -benchtime=10s -run=^$$ ./internal/service/auth/... -benchmem

# Run payment benchmarks
bench-payment:
	@echo "Running payment service benchmarks..."
	@go test -v -bench=. -benchtime=10s -run=^$$ ./internal/service/payment/... -benchmem

# Run database benchmarks
bench-db:
	@echo "Running database benchmarks..."
	@go test -v -bench=. -benchtime=10s -run=^$$ ./internal/repository/benchmarks/... -benchmem

# Run HTTP benchmarks
bench-http:
	@echo "Running HTTP benchmarks..."
	@go test -v -bench=. -benchtime=10s -run=^$$ ./internal/handler/benchmarks/... -benchmem

# Run all benchmarks with profiling
bench-all:
	@echo "Running all benchmarks with profiling..."
	@mkdir -p test/results/bench
	@go test -v -bench=. -benchtime=10s -run=^$$ \
		./internal/service/order/... \
		./internal/service/auth/... \
		./internal/service/payment/... \
		./internal/repository/benchmarks/... \
		./internal/handler/benchmarks/... \
		-benchmem \
		-cpuprofile=test/results/bench/cpu.prof \
		-memprofile=test/results/bench/mem.prof \
		| tee test/results/bench/results.txt

# Compare benchmark results
bench-compare:
	@echo "Comparing benchmark results..."
	@if [ -z "$(OLD)" ] || [ -z "$(NEW)" ]; then \
		echo "Usage: make bench-compare OLD=<old_file> NEW=<new_file>"; \
		exit 1; \
	fi
	@which benchstat >/dev/null 2>&1 || go install golang.org/x/perf/cmd/benchstat@latest
	@benchstat $(OLD) $(NEW)

# Run load test with vegeta
load-test:
	@echo "Running load tests..."
	@bash ./tests/load/run_load_test.sh $(SCENARIO) $(RATE) $(DURATION)

# Load test for auth endpoints
load-test-auth:
	@echo "Running load test for auth endpoints..."
	@which vegeta >/dev/null 2>&1 || go install github.com/tsenart/vegeta@latest
	@mkdir -p test/results/load
	@echo "Make sure the API server is running on localhost:8080"
	@echo "Update auth_targets.txt with valid token before running"
	@vegeta attack -targets=tests/load/vegeta/auth_targets.txt -rate=100 -duration=30s -output=test/results/load/auth_results.bin
	@vegeta report -type=text -input=test/results/load/auth_results.bin
	@vegeta report -type=html -input=test/results/load/auth_results.bin -output=test/results/load/auth_report.html
	@echo "Report generated: test/results/load/auth_report.html"

# Load test for order endpoints
load-test-order:
	@echo "Running load test for order endpoints..."
	@which vegeta >/dev/null 2>&1 || go install github.com/tsenart/vegeta@latest
	@mkdir -p test/results/load
	@echo "Make sure the API server is running on localhost:8080"
	@echo "Update order_targets.txt with valid token before running"
	@vegeta attack -targets=tests/load/vegeta/order_targets.txt -rate=50 -duration=60s -output=test/results/load/order_results.bin
	@vegeta report -type=text -input=test/results/load/order_results.bin
	@vegeta report -type=html -input=test/results/load/order_results.bin -output=test/results/load/order_report.html
	@echo "Report generated: test/results/load/order_report.html"

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
	@rm -rf test/results/bench
	@rm -rf test/results/load

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

# Install benchmark and load testing tools
bench-tools:
	@echo "Installing benchmark and load testing tools..."
	@go install golang.org/x/perf/cmd/benchstat@latest
	@go install github.com/tsenart/vegeta@latest
	@echo "Benchmark tools installed: benchstat, vegeta"

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
	@echo "  test                    Run all tests"
	@echo "  test-coverage           Run tests with coverage"
	@echo "  test-coverage-html      Run tests and generate HTML coverage report"
	@echo "  test-watch              Run tests in watch mode"
	@echo "  test-verbose            Run tests with verbose output"
	@echo "  test-race               Run tests with race detector"
	@echo "  test-bench              Run all benchmarks (10s each)"
	@echo "  test-integration        Run integration tests"
	@echo "  test-integration-db     Run integration tests with test DB (auto-sets env vars)"
	@echo "  run-test PKG=pkg        Run tests for specific package"
	@echo "  cover-pkg PKG=pkg       Run coverage for specific package"
	@echo ""
	@echo "Benchmark targets:"
	@echo "  bench-tools             Install benchmark tools (benchstat, vegeta)"
	@echo "  bench-all               Run all benchmarks with profiling"
	@echo "  bench-order             Run order service benchmarks"
	@echo "  bench-auth              Run auth service benchmarks"
	@echo "  bench-payment           Run payment service benchmarks"
	@echo "  bench-db                Run database benchmarks"
	@echo "  bench-http              Run HTTP endpoint benchmarks"
	@echo "  bench-suite SUITE=name   Run specific benchmark suite"
	@echo "  bench-compare OLD=x NEW=y Compare benchmark results"
	@echo ""
	@echo "Load testing targets:"
	@echo "  load-test SCENARIO=x RATE=y DURATION=z  Run load tests with vegeta"
	@echo "  load-test-auth          Load test auth endpoints"
	@echo "  load-test-order         Load test order endpoints"
	@echo ""
	@echo "Other targets:"
	@echo "  clean                   Clean test and benchmark artifacts"
	@echo "  generate-mocks          Generate mocks"
	@echo "  swagger                 Generate Swagger documentation (supports generics)"
	@echo "  swagger-vendor          Generate Swagger documentation with vendor support"
	@echo "  clean-swagger           Clean Swagger documentation"
	@echo "  lint                    Run linter"
	@echo "  fmt                     Format code"
	@echo "  vet                     Vet code"
	@echo "  check                   Run all checks"
	@echo "  deps                    Install dependencies"
	@echo "  test-tools              Install test tools"
	@echo "  help                    Show this help"
