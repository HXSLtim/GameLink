#!/bin/bash
# Load Testing Script for GameLink API using Vegeta
#
# Requirements:
#   - Install vegeta: go install github.com/tsenart/vegeta@latest
#   - Start the API server: go run cmd/main.go
#
# Usage:
#   ./run_load_test.sh [scenario] [rate] [duration]
#
# Examples:
#   ./run_load_test.sh auth 100 30s    # Test auth endpoints at 100 req/s for 30s
#   ./run_load_test.sh order 50 60s    # Test order endpoints at 50 req/s for 60s
#   ./run_load_test.sh all 100 30s     # Test all endpoints

set -e

# Configuration
SCENARIO=${1:-auth}
RATE=${2:-100}
DURATION=${3:-30s}
RESULTS_DIR="./test/results/load"
TIMESTAMP=$(date +%Y%m%d_%H%M%S)
TARGETS_FILE="./tests/load/vegeta/${SCENARIO}_targets.txt"
RESULTS_FILE="${RESULTS_DIR}/results_${SCENARIO}_${TIMESTAMP}.bin"
REPORT_TEXT="${RESULTS_DIR}/report_${SCENARIO}_${TIMESTAMP}.txt"
REPORT_HTML="${RESULTS_DIR}/report_${SCENARIO}_${TIMESTAMP}.html"

# Create results directory
mkdir -p "$RESULTS_DIR"

echo "======================================"
echo "GameLink API Load Testing"
echo "======================================"
echo "Scenario:   $SCENARIO"
echo "Rate:        $RATE requests/second"
echo "Duration:    $DURATION"
echo "Results:     $RESULTS_FILE"
echo "======================================"

# Check if vegeta is installed
if ! command -v vegeta &> /dev/null; then
    echo "Error: vegeta is not installed"
    echo "Install it with: go install github.com/tsenart/vegeta@latest"
    exit 1
fi

# Check if targets file exists
if [ ! -f "$TARGETS_FILE" ]; then
    echo "Error: Targets file not found: $TARGETS_FILE"
    exit 1
fi

# Run load test
echo ""
echo "Starting load test..."
vegeta attack \
    -targets="$TARGETS_FILE" \
    -rate="$RATE" \
    -duration="$DURATION" \
    -output="$RESULTS_FILE" \
    -workers=50 \
    -max-workers=100

# Generate text report
echo ""
echo "Generating text report..."
vegeta report \
    -type=text \
    -input="$RESULTS_FILE" \
    -output="$REPORT_TEXT"

# Generate HTML report
echo "Generating HTML report..."
vegeta report \
    -type=html \
    -input="$RESULTS_FILE" \
    -output="$REPORT_HTML"

# Display results
echo ""
echo "======================================"
echo "Load Test Results"
echo "======================================"
cat "$REPORT_HTML"

echo ""
echo "Reports saved:"
echo "  - Text: $REPORT_TEXT"
echo "  - HTML: $REPORT_HTML"
echo "  - Binary: $RESULTS_FILE"
echo ""
echo "You can analyze the binary file later with:"
echo "  vegeta report -type=text -input=$RESULTS_FILE"
echo "  vegeta report -type=hist -input=$RESULTS_FILE"
echo "======================================"

# Open HTML report (macOS only)
if [[ "$OSTYPE" == "darwin"* ]]; then
    open "$REPORT_HTML"
fi
