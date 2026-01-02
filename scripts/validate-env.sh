#!/bin/bash
# GameLink Environment Variables Validation Script
# This script checks if all required environment variables are properly configured

set -e

echo "GameLink Environment Variables Validation"
echo "========================================="
echo ""

ENV_FILE=".env"
ERRORS=0
WARNINGS=0

# Colors
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
CYAN='\033[0;36m'
GRAY='\033[0;37m'
NC='\033[0m' # No Color

# Check if .env file exists
if [ ! -f "$ENV_FILE" ]; then
    echo -e "${RED}✗ ERROR: .env file not found!${NC}"
    echo "  Copy .env.example to .env and configure your environment variables."
    exit 1
fi

echo -e "${GREEN}✓ .env file found${NC}"
echo ""

# Load environment variables from .env file
set -a
source "$ENV_FILE"
set +a

# Validation functions
test_required_var() {
    local name="$1"
    local description="$2"
    local min_length="${3:-0}"

    local value="${!name}"

    if [ -z "$value" ]; then
        echo -e "${RED}✗ ERROR: $name is required ($description)${NC}"
        return 1
    fi

    if [ "$min_length" -gt 0 ] && [ ${#value} -lt "$min_length" ]; then
        echo -e "${RED}✗ ERROR: $name must be at least $min_length characters (current: ${#value})${NC}"
        return 1
    fi

    echo -e "${GREEN}✓ $name configured${NC}"
    return 0
}

test_optional_var() {
    local name="$1"
    local description="$2"

    local value="${!name}"

    if [ -z "$value" ]; then
        echo -e "${YELLOW}⚠ WARNING: $name not set ($description)${NC}"
        return 1
    fi

    echo -e "${GREEN}✓ $name configured${NC}"
    return 0
}

# Check required variables
echo -e "${CYAN}Required Variables:${NC}"
echo -e "${CYAN}-------------------${NC}"

if ! test_required_var "APP_ENV" "Application environment"; then
    ((ERRORS++))
fi

if [ "$DB_TYPE" = "postgres" ]; then
    if ! test_required_var "DB_DSN" "Database connection string"; then
        ((ERRORS++))
    fi
fi

if ! test_required_var "CACHE_TYPE" "Cache type (redis/memory)"; then
    ((ERRORS++))
fi

if [ "$CACHE_TYPE" = "redis" ]; then
    if ! test_required_var "REDIS_ADDR" "Redis server address"; then
        ((ERRORS++))
    fi
fi

# Check security variables
echo ""
echo -e "${CYAN}Security Variables:${NC}"
echo -e "${CYAN}--------------------${NC}"

if [ "$CRYPTO_ENABLED" = "true" ]; then
    if ! test_required_var "CRYPTO_SECRET_KEY" "Encryption secret key (32 bytes)" 32; then
        ((ERRORS++))
    else
        if [ ${#CRYPTO_SECRET_KEY} -ne 32 ]; then
            echo -e "${YELLOW}✗ WARNING: CRYPTO_SECRET_KEY should be exactly 32 bytes for AES-256-CBC (current: ${#CRYPTO_SECRET_KEY})${NC}"
            ((WARNINGS++))
        fi
    fi

    if ! test_required_var "CRYPTO_IV" "Encryption IV (16 bytes)" 16; then
        ((ERRORS++))
    else
        if [ ${#CRYPTO_IV} -ne 16 ]; then
            echo -e "${YELLOW}✗ WARNING: CRYPTO_IV should be exactly 16 bytes for AES-CBC (current: ${#CRYPTO_IV})${NC}"
            ((WARNINGS++))
        fi
    fi
fi

if ! test_required_var "JWT_SECRET_KEY" "JWT signing secret (min 32 chars recommended)" 16; then
    ((ERRORS++))
fi

if ! test_required_var "SUPER_ADMIN_PASSWORD" "Super admin password" 8; then
    ((ERRORS++))
else
    # Check password strength for production
    if [ "$APP_ENV" = "production" ]; then
        if [[ ! "$SUPER_ADMIN_PASSWORD" =~ [A-Z] ]] || \
           [[ ! "$SUPER_ADMIN_PASSWORD" =~ [a-z] ]] || \
           [[ ! "$SUPER_ADMIN_PASSWORD" =~ [0-9] ]] || \
           [[ ! "$SUPER_ADMIN_PASSWORD" =~ [^a-zA-Z0-9] ]]; then
            echo -e "${YELLOW}✗ WARNING: SUPER_ADMIN_PASSWORD should contain uppercase, lowercase, digit, and special character in production${NC}"
            ((WARNINGS++))
        fi
    fi
fi

if ! test_required_var "SUPER_ADMIN_EMAIL" "Super admin email"; then
    ((ERRORS++))
fi

# Optional variables
echo ""
echo -e "${CYAN}Optional Variables:${NC}"
echo -e "${CYAN}--------------------${NC}"

test_optional_var "SERVICE_PORT" "Server port (default: 8080)" || true
test_optional_var "ENABLE_SWAGGER" "Enable Swagger documentation" || true

# Production-specific checks
echo ""
echo -e "${CYAN}Production Checks:${NC}"
echo -e "${CYAN}--------------------${NC}"

if [ "$APP_ENV" = "production" ]; then
    if [ "$CRYPTO_ENABLED" != "true" ]; then
        echo -e "${RED}✗ ERROR: CRYPTO_ENABLED must be true in production${NC}"
        ((ERRORS++))
    else
        echo -e "${GREEN}✓ CRYPTO_ENABLED is true (required for production)${NC}"
    fi

    if [ "$CACHE_TYPE" != "redis" ]; then
        echo -e "${RED}✗ ERROR: CACHE_TYPE must be 'redis' in production (current: $CACHE_TYPE)${NC}"
        ((ERRORS++))
    else
        echo -e "${GREEN}✓ CACHE_TYPE is redis (required for production)${NC}"
    fi
else
    echo -e "${GRAY}ℹ Not in production mode, skipping production checks${NC}"
fi

# Summary
echo ""
echo "========================================="
echo -e "${CYAN}Validation Summary${NC}"
echo "========================================="

if [ $ERRORS -eq 0 ] && [ $WARNINGS -eq 0 ]; then
    echo -e "${GREEN}✓ All checks passed! Environment is properly configured.${NC}"
    exit 0
elif [ $ERRORS -eq 0 ]; then
    echo -e "${YELLOW}⚠ Validation passed with $WARNINGS warning(s). Review the warnings above.${NC}"
    exit 0
else
    echo -e "${RED}✗ Validation failed with $ERRORS error(s) and $WARNINGS warning(s).${NC}"
    echo "  Please fix the errors before starting the application."
    exit 1
fi
