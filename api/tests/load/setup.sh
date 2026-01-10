#!/bin/bash
# GameLink Performance Benchmarking Setup Script
#
# This script helps you set up the benchmarking environment
#
# Usage: ./setup.sh

set -e

echo "=========================================="
echo "GameLink Benchmarking Environment Setup"
echo "=========================================="
echo ""

# Check if Docker is installed
if ! command -v docker &> /dev/null; then
    echo "⚠️  Docker is not installed. Please install Docker first."
    echo "   Visit: https://docs.docker.com/get-docker/"
    exit 1
fi

echo "✅ Docker is installed"

# Check if Go is installed
if ! command -v go &> /dev/null; then
    echo "⚠️  Go is not installed. Please install Go 1.25+ first."
    echo "   Visit: https://golang.org/dl/"
    exit 1
fi

echo "✅ Go is installed"

# Check Go version
GO_VERSION=$(go version | awk '{print $3}' | sed 's/go//')
echo "   Go version: $GO_VERSION"

echo ""
echo "=========================================="
echo "Step 1: Installing Benchmark Tools"
echo "=========================================="
echo ""

make bench-tools

echo ""
echo "✅ Benchmark tools installed:"
echo "   - benchstat: For comparing benchmark results"
echo "   - vegeta: For HTTP load testing"
echo ""

echo "=========================================="
echo "Step 2: Setting Up Benchmark Database"
echo "=========================================="
echo ""

# Check if gamelink-bench-db container already exists
if docker ps -a | grep -q gamelink-bench-db; then
    echo "⚠️  Benchmark database container already exists"
    read -p "   Recreate database? (y/N): " -n 1 -r
    echo
    if [[ $REPLY =~ ^[Yy]$ ]]; then
        echo "   Removing old container..."
        docker rm -f gamelink-bench-db
        sleep 2
    else
        echo "   Using existing container"
    fi
fi

# Create container if it doesn't exist
if ! docker ps -a | grep -q gamelink-bench-db; then
    echo "   Creating benchmark database container..."

    docker run -d \
        --name gamelink-bench-db \
        -e POSTGRES_USER=gamelink \
        -e POSTGRES_PASSWORD=gamelink \
        -e POSTGRES_DB=gamelink_bench \
        -p 5432:5432 \
        postgres:16

    echo "   Waiting for database to be ready..."
    sleep 5

    # Check if database is ready
    for i in {1..30}; do
        if docker exec gamelink-bench-db pg_isready -U gamelink &> /dev/null; then
            echo "   ✅ Database is ready!"
            break
        fi
        if [ $i -eq 30 ]; then
            echo "   ❌ Database failed to start"
            exit 1
        fi
        sleep 1
    done
else
    echo "   ✅ Database container exists and is running"
fi

echo ""
echo "=========================================="
echo "Step 3: Creating Required Directories"
echo "=========================================="
echo ""

mkdir -p test/results/bench
mkdir -p test/results/load
mkdir -p docs/baseline/archive

echo "✅ Directories created:"
echo "   - test/results/bench  (for benchmark results)"
echo "   - test/results/load   (for load test results)"
echo "   - docs/baseline/archive (for baseline history)"
echo ""

echo "=========================================="
echo "Step 4: Setting Environment Variables"
echo "=========================================="
echo ""

# Create .env file for benchmarking
cat > .env.bench << 'EOF'
# Benchmark Database Configuration
export BENCH_DB_HOST=localhost
export BENCH_DB_PORT=5432
export BENCH_DB_USER=gamelink
export BENCH_DB_PASSWORD=gamelink
export BENCH_DB_NAME=gamelink_bench

# Optional: Enable profiling
export BENCH_CPU_PROFILE=test/results/bench/cpu.prof
export BENCH_MEM_PROFILE=test/results/bench/mem.prof
EOF

echo "✅ Environment variables saved to .env.bench"
echo ""
echo "   To use these variables, run:"
echo "   source .env.bench"
echo ""

echo "=========================================="
echo "Step 5: Running Health Check"
echo "=========================================="
echo ""

# Test database connection
echo "Testing database connection..."
if docker exec gamelink-bench-db psql -U gamelink -d gamelink_bench -c "SELECT 1;" &> /dev/null; then
    echo "✅ Database connection successful"
else
    echo "❌ Database connection failed"
    exit 1
fi

echo ""
echo "=========================================="
echo "Setup Complete! 🎉"
echo "=========================================="
echo ""
echo "Next Steps:"
echo ""
echo "1. Load environment variables:"
echo "   source .env.bench"
echo ""
echo "2. Run initial benchmarks:"
echo "   make bench-all"
echo ""
echo "3. View results:"
echo "   cat test/results/bench/results.txt"
echo ""
echo "4. Run load tests (requires API server):"
echo "   # Terminal 1:"
echo "   go run cmd/main.go"
echo "   # Terminal 2:"
echo "   make load-test-auth"
echo ""
echo "5. For more information, see:"
echo "   - docs/PERFORMANCE_BENCHMARKING.md"
echo "   - docs/BASELINE_METRICS.md"
echo "   - tests/load/README.md"
echo ""
echo "=========================================="
