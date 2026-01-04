@echo off
REM GameLink Order Handler Test Runner (Windows)
REM This script provides convenient commands to run order handler tests

SETLOCAL EnableDelayedExpansion

SET GREEN=[92m
SET RED=[91m
SET YELLOW=[93m
SET NC=[0m

REM Change to api directory
CD /D "%~dp0..\.."

echo.
echo %GREEN%=== GameLink Order Handler Test Suite ===%NC%
echo.

IF "%1"=="help" GOTO :help
IF "%1"=="-h" GOTO :help
IF "%1"=="--help" GOTO :help

SET "OPTION=%~1"
IF "%OPTION%"=="" SET "OPTION=all"

IF "%OPTION%"=="all" GOTO :all
IF "%OPTION%"=="unit" GOTO :unit
IF "%OPTION%"=="comprehensive" GOTO :comprehensive
IF "%OPTION%"=="coverage" GOTO :coverage
IF "%OPTION%"=="race" GOTO :race
IF "%OPTION%"=="verbose" GOTO :verbose
IF "%OPTION%"=="state" GOTO :state
IF "%OPTION%"=="refund" GOTO :refund
IF "%OPTION%"=="create" GOTO :create
IF "%OPTION%"=="quick" GOTO :quick

echo %RED%Unknown option: %OPTION%%NC%
echo.
GOTO :help

:all
echo.
echo %YELLOW%=== Running All Order Handler Tests ===%NC%
echo.
go test ./internal/handler/admin -run "TestOrderHandler" -v
IF %ERRORLEVEL% EQU 0 (
    echo.
    echo %GREEN%✓ All Order Handler Tests PASSED%NC%
) ELSE (
    echo.
    echo %RED%✗ All Order Handler Tests FAILED%NC%
)
GOTO :end

:unit
echo.
echo %YELLOW%=== Unit Tests Only (Skip Database) ===%NC%
echo.
go test ./internal/handler/admin -run "TestOrderHandler_Unit" -v -short
IF %ERRORLEVEL% EQU 0 (
    echo.
    echo %GREEN%✓ Unit Tests PASSED%NC%
) ELSE (
    echo.
    echo %RED%✗ Unit Tests FAILED%NC%
)
GOTO :end

:comprehensive
echo.
echo %YELLOW%=== Comprehensive Tests ===%NC%
echo.
go test ./internal/handler/admin -run "TestOrderHandler_Comprehensive" -v
IF %ERRORLEVEL% EQU 0 (
    echo.
    echo %GREEN%✓ Comprehensive Tests PASSED%NC%
) ELSE (
    echo.
    echo %RED%✗ Comprehensive Tests FAILED%NC%
)
GOTO :end

:coverage
echo.
echo %YELLOW%=== Generating Coverage Report ===%NC%
echo.
echo Running tests with coverage...
go test ./internal/handler/admin -run "TestOrderHandler" -coverprofile=coverage.out -covermode=atomic

IF %ERRORLEVEL% EQU 0 (
    echo.
    echo Generating HTML coverage report...
    go tool cover -html=coverage.out -o coverage.html

    echo Coverage summary:
    go tool cover -func=coverage.out | findstr /C:"total:"

    echo.
    echo %GREEN%✓ Coverage report generated: api\coverage.html%NC%
) ELSE (
    echo.
    echo %RED%✗ Coverage generation FAILED%NC%
)
GOTO :end

:race
echo.
echo %YELLOW%=== Race Detector Tests ===%NC%
echo.
go test ./internal/handler/admin -run "TestOrderHandler" -v -race
IF %ERRORLEVEL% EQU 0 (
    echo.
    echo %GREEN%✓ Race Detector Tests PASSED%NC%
) ELSE (
    echo.
    echo %RED%✗ Race Detector Tests FAILED%NC%
)
GOTO :end

:verbose
echo.
echo %YELLOW%=== Verbose Tests ===%NC%
echo.
go test ./internal/handler/admin -run "TestOrderHandler" -v -count=1
IF %ERRORLEVEL% EQU 0 (
    echo.
    echo %GREEN%✓ Verbose Tests PASSED%NC%
) ELSE (
    echo.
    echo %RED%✗ Verbose Tests FAILED%NC%
)
GOTO :end

:state
echo.
echo %YELLOW%=== State Machine Tests ===%NC%
echo.
go test ./internal/handler/admin -run "TestOrderHandler.*StateMachine" -v
IF %ERRORLEVEL% EQU 0 (
    echo.
    echo %GREEN%✓ State Machine Tests PASSED%NC%
) ELSE (
    echo.
    echo %RED%✗ State Machine Tests FAILED%NC%
)
GOTO :end

:refund
echo.
echo %YELLOW%=== Refund Tests ===%NC%
echo.
go test ./internal/handler/admin -run "TestOrderHandler.*RefundOrder" -v
IF %ERRORLEVEL% EQU 0 (
    echo.
    echo %GREEN%✓ Refund Tests PASSED%NC%
) ELSE (
    echo.
    echo %RED%✗ Refund Tests FAILED%NC%
)
GOTO :end

:create
echo.
echo %YELLOW%=== Order Creation Tests ===%NC%
echo.
go test ./internal/handler/admin -run "TestOrderHandler.*CreateOrder" -v
IF %ERRORLEVEL% EQU 0 (
    echo.
    echo %GREEN%✓ Order Creation Tests PASSED%NC%
) ELSE (
    echo.
    echo %RED%✗ Order Creation Tests FAILED%NC%
)
GOTO :end

:quick
echo.
echo %YELLOW%=== Quick Test Run (Short Mode) ===%NC%
echo.
go test ./internal/handler/admin -run "TestOrderHandler" -short
IF %ERRORLEVEL% EQU 0 (
    echo.
    echo %GREEN%✓ Quick Tests PASSED%NC%
) ELSE (
    echo.
    echo %RED%✗ Quick Tests FAILED%NC%
)
GOTO :end

:help
echo Usage: test-order-handler.bat [option]
echo.
echo Options:
echo   all           Run all order handler tests (default)
echo   unit          Run unit tests only (without database)
echo   comprehensive Run comprehensive test scenarios
echo   coverage      Generate coverage report
echo   race          Run with race detector
echo   verbose       Run with verbose output
echo   state         Run state machine tests
echo   refund        Run refund tests
echo   create        Run order creation tests
echo   quick         Quick test run (short mode)
echo   help, -h      Show this help message
echo.
echo Examples:
echo   test-order-handler.bat all
echo   test-order-handler.bat coverage
echo   test-order-handler.bat unit
echo.
GOTO :end

:end
echo.
echo %GREEN%Order handler test execution finished!%NC%
echo.
echo To view detailed coverage:
echo   go tool cover -html=api\coverage.out
echo.
echo To run specific tests:
echo   go test ./internal/handler/admin -run TestName -v
echo.

ENDLOCAL
