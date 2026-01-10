@echo off
REM GameLink 自动安装和启动脚本

echo ========================================
echo GameLink 自动安装和启动
echo ========================================
echo.

REM 检查是否在项目根目录
if not exist "admin\package.json" (
    echo ❌ 错误: 请在项目根目录中运行此脚本
    pause
    exit /b 1
)

echo 📦 步骤 1: 安装前端 PWA 依赖...
echo.

cd admin

echo 安装 vite-plugin-pwa, picocolors, workbox-window...
call npm install -D vite-plugin-pwa picocolors workbox-window

if errorlevel 1 (
    echo.
    echo ❌ 依赖安装失败!
    echo.
    echo 可能的解决方案:
    echo 1. 检查网络连接
    echo 2. 尝试使用淘宝镜像: npm config set registry https://registry.npmmirror.com
    echo 3. 手动删除 node_modules 后重新安装
    echo.
    pause
    exit /b 1
)

echo.
echo ✅ 依赖安装成功!
echo.

cd ..

echo.
echo ========================================
echo 📋 启动选项
echo ========================================
echo.
echo 1. 启动前端开发服务器
echo 2. 启动后端API服务 (开发模式)
echo 3. 启动后端API服务 (Docker)
echo 4. 同时启动前后端 (推荐)
echo.

set /p choice="请选择 (1-4, 默认: 4)"

if "%choice%"=="" set choice=4

if "%choice%"=="1" goto frontend
if "%choice%"=="2" goto backend_dev
if "%choice%"=="3" goto backend_docker
if "%choice%"=="4" goto both

:frontend
echo.
echo 🚀 启动前端开发服务器...
echo.
cd admin
call npm run dev
goto end

:backend_dev
echo.
echo 🚀 启动后端API服务 (开发模式)...
echo.
cd api
go run cmd/main.go
goto end

:backend_docker
echo.
echo 🐳 启动后端API服务 (Docker)...
echo.
docker-compose up -d
echo.
echo ✅ Docker服务已启动!
echo 📊 查看日志: docker-compose logs -f
goto end

:both
echo.
echo 🚀 同时启动前后端服务...
echo.
echo 正在后台启动后端API服务...
start "GameLink Backend" cmd /k "cd api && go run cmd/main.go"

timeout /t 3 /nobreak >nul

echo 正在启动前端开发服务器...
cd admin
call npm run dev

:end
pause
