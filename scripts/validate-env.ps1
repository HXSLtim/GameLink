#!/usr/bin/env pwsh
# GameLink Environment Variables Validation Script
# This script checks if all required environment variables are properly configured

Write-Host "GameLink Environment Variables Validation" -ForegroundColor Cyan
Write-Host "=========================================" -ForegroundColor Cyan
Write-Host ""

$envFile = ".env"
$errors = 0
$warnings = 0

# Check if .env file exists
if (-not (Test-Path $envFile)) {
    Write-Host "✗ ERROR: .env file not found!" -ForegroundColor Red
    Write-Host "  Copy .env.example to .env and configure your environment variables." -ForegroundColor Yellow
    exit 1
}

Write-Host "✓ .env file found" -ForegroundColor Green
Write-Host ""

# Load environment variables from .env file
Get-Content $envFile | ForEach-Object {
    if ($_ -match "^([^#][^=]+)=(.*)$") {
        $name = $matches[1].Trim()
        $value = $matches[2].Trim()
        Set-Item -Path "env:$name" -Value $value
    }
}

# Validation functions
function Test-RequiredVar {
    param(
        [string]$Name,
        [string]$Description,
        [int]$MinLength = 0
    )

    $value = [Environment]::GetEnvironmentVariable($Name)
    if ([string]::IsNullOrWhiteSpace($value)) {
        Write-Host "✗ ERROR: $Name is required ($Description)" -ForegroundColor Red
        return $false
    }

    if ($MinLength -gt 0 -and $value.Length -lt $MinLength) {
        Write-Host "✗ ERROR: $Name must be at least $MinLength characters (current: $($value.Length))" -ForegroundColor Red
        return $false
    }

    Write-Host "✓ $Name configured" -ForegroundColor Green
    return $true
}

function Test-OptionalVar {
    param(
        [string]$Name,
        [string]$Description
    )

    $value = [Environment]::GetEnvironmentVariable($Name)
    if ([string]::IsNullOrWhiteSpace($value)) {
        Write-Host "⚠ WARNING: $Name not set ($Description)" -ForegroundColor Yellow
        return $false
    }

    Write-Host "✓ $Name configured" -ForegroundColor Green
    return $true
}

# Check required variables
Write-Host "Required Variables:" -ForegroundColor Cyan
Write-Host "-------------------" -ForegroundColor Cyan

if (-not (Test-RequiredVar "APP_ENV" "Application environment")) { $errors++ }

$dbType = [Environment]::GetEnvironmentVariable("DB_TYPE")
if ($dbType -eq "postgres") {
    if (-not (Test-RequiredVar "DB_DSN" "Database connection string")) { $errors++ }
}

if (-not (Test-RequiredVar "CACHE_TYPE" "Cache type (redis/memory)")) { $errors++ }

$cacheType = [Environment]::GetEnvironmentVariable("CACHE_TYPE")
if ($cacheType -eq "redis") {
    if (-not (Test-RequiredVar "REDIS_ADDR" "Redis server address")) { $errors++ }
}

# Check security variables
Write-Host ""
Write-Host "Security Variables:" -ForegroundColor Cyan
Write-Host "--------------------" -ForegroundColor Cyan

$cryptoEnabled = [Environment]::GetEnvironmentVariable("CRYPTO_ENABLED")
if ($cryptoEnabled -eq "true") {
    if (-not (Test-RequiredVar "CRYPTO_SECRET_KEY" "Encryption secret key (32 bytes)" 32)) {
        $errors++
    } else {
        $key = [Environment]::GetEnvironmentVariable("CRYPTO_SECRET_KEY")
        if ($key.Length -ne 32) {
            Write-Host "✗ WARNING: CRYPTO_SECRET_KEY should be exactly 32 bytes for AES-256-CBC (current: $($key.Length))" -ForegroundColor Yellow
            $warnings++
        }
    }

    if (-not (Test-RequiredVar "CRYPTO_IV" "Encryption IV (16 bytes)" 16)) {
        $errors++
    } else {
        $iv = [Environment]::GetEnvironmentVariable("CRYPTO_IV")
        if ($iv.Length -ne 16) {
            Write-Host "✗ WARNING: CRYPTO_IV should be exactly 16 bytes for AES-CBC (current: $($iv.Length))" -ForegroundColor Yellow
            $warnings++
        }
    }
}

if (-not (Test-RequiredVar "JWT_SECRET_KEY" "JWT signing secret (min 32 chars recommended)" 16)) {
    $errors++
}

$superAdminPassword = [Environment]::GetEnvironmentVariable("SUPER_ADMIN_PASSWORD")
if (-not (Test-RequiredVar "SUPER_ADMIN_PASSWORD" "Super admin password" 8)) {
    $errors++
} else {
    # Check password strength for production
    $appEnv = [Environment]::GetEnvironmentVariable("APP_ENV")
    if ($appEnv -eq "production") {
        $hasUpper = $superAdminPassword -cmatch '[A-Z]'
        $hasLower = $superAdminPassword -cmatch '[a-z]'
        $hasDigit = $superAdminPassword -match '\d'
        $hasSpecial = $superAdminPassword -match '[^a-zA-Z0-9]'

        if (-not ($hasUpper -and $hasLower -and $hasDigit -and $hasSpecial)) {
            Write-Host "✗ WARNING: SUPER_ADMIN_PASSWORD should contain uppercase, lowercase, digit, and special character in production" -ForegroundColor Yellow
            $warnings++
        }
    }
}

if (-not (Test-RequiredVar "SUPER_ADMIN_EMAIL" "Super admin email")) {
    $errors++
}

# Optional variables
Write-Host ""
Write-Host "Optional Variables:" -ForegroundColor Cyan
Write-Host "--------------------" -ForegroundColor Cyan

Test-OptionalVar "SERVICE_PORT" "Server port (default: 8080)" | Out-Null
Test-OptionalVar "ENABLE_SWAGGER" "Enable Swagger documentation" | Out-Null

# Production-specific checks
Write-Host ""
Write-Host "Production Checks:" -ForegroundColor Cyan
Write-Host "--------------------" -ForegroundColor Cyan

if ([Environment]::GetEnvironmentVariable("APP_ENV") -eq "production") {
    if ($cryptoEnabled -ne "true") {
        Write-Host "✗ ERROR: CRYPTO_ENABLED must be true in production" -ForegroundColor Red
        $errors++
    } else {
        Write-Host "✓ CRYPTO_ENABLED is true (required for production)" -ForegroundColor Green
    }

    if ($cacheType -ne "redis") {
        Write-Host "✗ ERROR: CACHE_TYPE must be 'redis' in production (current: $cacheType)" -ForegroundColor Red
        $errors++
    } else {
        Write-Host "✓ CACHE_TYPE is redis (required for production)" -ForegroundColor Green
    }
} else {
    Write-Host "ℹ Not in production mode, skipping production checks" -ForegroundColor Gray
}

# Summary
Write-Host ""
Write-Host "=========================================" -ForegroundColor Cyan
Write-Host "Validation Summary" -ForegroundColor Cyan
Write-Host "=========================================" -ForegroundColor Cyan

if ($errors -eq 0 -and $warnings -eq 0) {
    Write-Host "✓ All checks passed! Environment is properly configured." -ForegroundColor Green
    exit 0
} elseif ($errors -eq 0) {
    Write-Host "⚠ Validation passed with $warnings warning(s). Review the warnings above." -ForegroundColor Yellow
    exit 0
} else {
    Write-Host "✗ Validation failed with $errors error(s) and $warnings warning(s)." -ForegroundColor Red
    Write-Host "  Please fix the errors before starting the application." -ForegroundColor Yellow
    exit 1
}
