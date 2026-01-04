@echo off
REM GameLink 前端开发服务器启动脚本

echo ========================================
echo GameLink Frontend Dev Server
echo ========================================
echo.

REM 切换到admin目录
cd /d "%~dp0admin"

echo 检查依赖...
if not exist "node_modules" (
    echo 安装npm依赖...
    call npm install
)

echo.
echo 🚀 启动前端开发服务器...
echo.

REM 启动开发服务器
call npm run dev

pause
