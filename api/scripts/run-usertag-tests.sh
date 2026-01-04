#!/bin/bash
# Script to run UserTag handler tests

set -e

echo "========================================="
echo "UserTag Handler Test Suite"
echo "========================================="
echo ""

cd "$(dirname "$0")/.."

echo "1. Running all UserTag handler tests..."
go test ./internal/handler/admin -run TestUserTagHandler -v

echo ""
echo "========================================="
echo "All UserTag tests passed!"
echo "========================================="
echo ""
echo "To run specific test categories:"
echo "  - CreateTag tests: go test ./internal/handler/admin -run TestUserTagHandler_CreateTag -v"
echo "  - Batch operations: go test ./internal/handler/admin -run TestUserTagHandler_Batch -v"
echo "  - With coverage: go test ./internal/handler/admin -run TestUserTagHandler -cover"
echo ""
