@echo off
REM GameLink 后端API服务启动脚本

echo ========================================
echo GameLink Backend API
echo ========================================
echo.

REM 检查是否在项目根目录
if not exist "api\cmd\main.go" (
    echo ❌ 错误: 请在项目根目录中运行此脚本
    pause
    exit /b 1
)

echo 启动方式:
echo 1. 开发模式 (直接运行)
echo 2. Docker Compose
echo.

set /p mode="请选择 (1-2, 默认: 1)"

if "%mode%"=="" set mode=1

if "%mode%"=="1" (
    echo.
    echo 🚀 启动后端开发服务器...
    echo.

    cd api

    REM 检查Go环境
    go version >nul 2>&1
    if errorlevel 1 (
        echo ❌ Go未安装或未在PATH中
        pause
        exit /b 1
    )

    REM 启动服务
    go run cmd/main.go
) else if "%mode%"=="2" (
    echo.
    echo 🐳 启动Docker Compose...
    echo.

    docker-compose up -d

    echo.
    echo ✅ Docker服务已启动!
    echo.
    echo 📊 查看日志: docker-compose logs -f api
    echo 🔍 查看状态: docker-compose ps
    echo.

    pause
) else (
    echo ❌ 无效的选择
    pause
    exit /b 1
)
