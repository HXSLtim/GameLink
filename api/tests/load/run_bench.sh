#!/bin/bash
# Comprehensive Benchmark Script for GameLink API
#
# Usage:
#   ./run_bench.sh [suite] [benchtime]
#
# Examples:
#   ./run_bench.sh order 10s       # Run order benchmarks for 10 seconds
#   ./run_bench.sh auth 5s         # Run auth benchmarks for 5 seconds
#   ./run_bench.sh all 10s         # Run all benchmarks

set -e

# Configuration
SUITE=${1:-all}
BENCHTIME=${2:-10s}
BENCHMARK_DIR="./test/results/bench"
TIMESTAMP=$(date +%Y%m%d_%H%M%S)
CPU_PROFILE="${BENCHMARK_DIR}/cpu_${SUITE}_${TIMESTAMP}.prof"
MEM_PROFILE="${BENCHMARK_DIR}/mem_${SUITE}_${TIMESTAMP}.prof"
BENCH_OUTPUT="${BENCHMARK_DIR}/bench_${SUITE}_${TIMESTAMP}.txt"

# Create benchmark directory
mkdir -p "$BENCHMARK_DIR"

echo "======================================"
echo "GameLink API Benchmarking"
echo "======================================"
echo "Suite:       $SUITE"
echo "Bench Time:  $BENCHTIME"
echo "Output:      $BENCH_OUTPUT"
echo "======================================"

# Set environment for benchmarking
export BENCH_DB_HOST=${BENCH_DB_HOST:-localhost}
export BENCH_DB_PORT=${BENCH_DB_PORT:-5432}
export BENCH_DB_USER=${BENCH_DB_USER:-gamelink}
export BENCH_DB_PASSWORD=${BENCH_DB_PASSWORD:-gamelink}
export BENCH_DB_NAME=${BENCH_DB_NAME:-gamelink_bench}

# Run benchmarks based on suite
case $SUITE in
  order)
    echo ""
    echo "Running Order Service Benchmarks..."
    go test -v -bench=. -benchtime="$BENCHTIME" -run=^$ ./internal/service/order/... \
        -benchmem \
        | tee "$BENCH_OUTPUT"
    ;;

  auth)
    echo ""
    echo "Running Auth Service Benchmarks..."
    go test -v -bench=. -benchtime="$BENCHTIME" -run=^$ ./internal/service/auth/... \
        -benchmem \
        | tee "$BENCH_OUTPUT"
    ;;

  payment)
    echo ""
    echo "Running Payment Service Benchmarks..."
    go test -v -bench=. -benchtime="$BENCHTIME" -run=^$ ./internal/service/payment/... \
        -benchmem \
        | tee "$BENCH_OUTPUT"
    ;;

  database)
    echo ""
    echo "Running Database Benchmarks..."
    go test -v -bench=. -benchtime="$BENCHTIME" -run=^$ ./internal/repository/benchmarks/... \
        -benchmem \
        | tee "$BENCH_OUTPUT"
    ;;

  http)
    echo ""
    echo "Running HTTP Benchmarks..."
    go test -v -bench=. -benchtime="$BENCHTIME" -run=^$ ./internal/handler/benchmarks/... \
        -benchmem \
        | tee "$BENCH_OUTPUT"
    ;;

  all)
    echo ""
    echo "Running All Benchmarks..."
    go test -v -bench=. -benchtime="$BENCHTIME" -run=^$ \
        ./internal/service/order/... \
        ./internal/service/auth/... \
        ./internal/service/payment/... \
        ./internal/repository/benchmarks/... \
        ./internal/handler/benchmarks/... \
        -benchmem \
        | tee "$BENCH_OUTPUT"
    ;;

  *)
    echo "Error: Unknown suite '$SUITE'"
    echo "Available suites: order, auth, payment, database, http, all"
    exit 1
    ;;
esac

echo ""
echo "======================================"
echo "Benchmarking Complete"
echo "======================================"
echo "Results saved to: $BENCH_OUTPUT"
echo ""
echo "To generate a flame graph:"
echo "  go tool pprof -http=:8080 $CPU_PROFILE"
echo "======================================"

# Compare with previous results if available
LATEST_RESULT=$(ls -t "$BENCHMARK_DIR"/bench_${SUITE}_*.txt 2>/dev/null | head -2 | tail -1)
if [ -n "$LATEST_RESULT" ] && [ "$LATEST_RESULT" != "$BENCH_OUTPUT" ]; then
    echo ""
    echo "Comparing with previous run:"
    echo "  Previous: $LATEST_RESULT"
    echo "  Current:  $BENCH_OUTPUT"
    echo ""
    echo "To compare:"
    echo "  benchstat $LATEST_RESULT $BENCH_OUTPUT"
fi
