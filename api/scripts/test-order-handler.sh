#!/bin/bash

# GameLink Order Handler Test Runner
# This script provides convenient commands to run order handler tests

set -e

# Colors for output
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
NC='\033[0m' # No Color

# Script directory
SCRIPT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"
PROJECT_ROOT="$(cd "$SCRIPT_DIR/.." && pwd)"
cd "$PROJECT_ROOT"

echo -e "${GREEN}=== GameLink Order Handler Test Suite ===${NC}\n"

# Function to print section headers
print_header() {
    echo -e "\n${YELLOW}=== $1 ===${NC}\n"
}

# Function to run tests
run_tests() {
    local test_name="$1"
    local test_pattern="$2"
    local extra_args="$3"

    print_header "Running: $test_name"
    echo "Command: go test ./internal/handler/admin -run '$test_pattern' $extra_args"

    if go test ./internal/handler/admin -run "$test_pattern" $extra_args; then
        echo -e "\n${GREEN}✓ $test_name PASSED${NC}"
    else
        echo -e "\n${RED}✗ $test_name FAILED${NC}"
        return 1
    fi
}

# Check if PostgreSQL test database is available
check_test_db() {
    print_header "Checking Test Database"

    if docker ps | grep -q "gamelink-test-db"; then
        echo -e "${GREEN}✓ Test database container is running${NC}"
        return 0
    elif docker ps | grep -q "postgres"; then
        echo -e "${YELLOW}⚠ PostgreSQL container found (may not be test DB)${NC}"
        return 0
    else
        echo -e "${RED}✗ Test database not found${NC}"
        echo -e "\n${YELLOW}To start test database:${NC}"
        echo "  docker-compose -f docker-compose.test.yml up -d"
        echo -e "\n${YELLOW}To skip database tests:${NC}"
        echo "  Run with: ./scripts/test-order-handler.sh unit"
        return 1
    fi
}

# Main menu
if [ "$1" == "--help" ] || [ "$1" == "-h" ]; then
    echo "Usage: ./scripts/test-order-handler.sh [option]"
    echo ""
    echo "Options:"
    echo "  all           Run all order handler tests (default)"
    echo "  unit          Run unit tests only (without database)"
    echo "  coverage      Generate coverage report"
    echo "  race          Run with race detector"
    echo "  verbose       Run with verbose output"
    echo "  state         Run state machine tests"
    echo "  refund        Run refund tests"
    echo "  create        Run order creation tests"
    echo "  quick         Quick test run (short mode)"
    echo "  help, -h      Show this help message"
    echo ""
    echo "Examples:"
    echo "  ./scripts/test-order-handler.sh all"
    echo "  ./scripts/test-order-handler.sh coverage"
    echo "  ./scripts/test-order-handler.sh unit"
    exit 0
fi

# Parse command line arguments
OPTION="${1:-all}"

case "$OPTION" in
    all)
        check_test_db || true
        run_tests "All Order Handler Tests" "TestOrderHandler" "-v"
        ;;

    unit)
        print_header "Unit Tests Only (Skip Database)"
        run_tests "Unit Tests" "TestOrderHandler_Unit" "-v -short"
        ;;

    comprehensive)
        check_test_db || true
        run_tests "Comprehensive Tests" "TestOrderHandler_Comprehensive" "-v"
        ;;

    coverage)
        check_test_db || true
        print_header "Generating Coverage Report"

        echo "Running tests with coverage..."
        go test ./internal/handler/admin -run "TestOrderHandler" -coverprofile=coverage.out -covermode=atomic

        echo "Generating HTML coverage report..."
        go tool cover -html=coverage.out -o coverage.html

        echo "Coverage summary:"
        go tool cover -func=coverage.out | tail -1

        echo -e "\n${GREEN}✓ Coverage report generated: coverage.html${NC}"
        ;;

    race)
        check_test_db || true
        run_tests "Race Detector Tests" "TestOrderHandler" "-v -race"
        ;;

    verbose)
        check_test_db || true
        run_tests "Verbose Tests" "TestOrderHandler" "-v -count=1"
        ;;

    state)
        check_test_db || true
        run_tests "State Machine Tests" "TestOrderHandler.*StateMachine" "-v"
        ;;

    refund)
        check_test_db || true
        run_tests "Refund Tests" "TestOrderHandler.*RefundOrder" "-v"
        ;;

    create)
        check_test_db || true
        run_tests "Order Creation Tests" "TestOrderHandler.*CreateOrder" "-v"
        ;;

    quick)
        print_header "Quick Test Run (Short Mode)"
        run_tests "Quick Tests" "TestOrderHandler" "-short"
        ;;

    *)
        echo -e "${RED}Unknown option: $OPTION${NC}"
        echo "Run './scripts/test-order-handler.sh help' for usage"
        exit 1
        ;;
esac

# Print summary
print_header "Test Suite Complete"

echo -e "${GREEN}Order handler test execution finished!${NC}"
echo ""
echo "To view detailed coverage:"
echo "  go tool cover -html=api/coverage.out"
echo ""
echo "To run specific tests:"
echo "  go test ./internal/handler/admin -run TestName -v"
