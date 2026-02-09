@echo off
REM GameLink 健康检查脚本 (Windows 版本)
setlocal enabledelayedexpansion

echo =========================================
echo GameLink 服务健康检查
echo 检查时间: %date% %time%
echo =========================================
echo.

set TOTAL=0
set PASSED=0
set FAILED=0

REM 检查 Docker 容器
echo === Docker 容器检查 ===
docker ps --format "table {{.Names}}\t{{.Status}}" 2>nul
if %errorlevel% equ 0 (
    echo [√] Docker 运行正常
    set /a PASSED+=1
) else (
    echo [X] Docker 未运行
    set /a FAILED+=1
)
set /a TOTAL+=1
echo.

REM 检查端口监听
echo === 服务端口检查 ===
netstat -ano | findstr "LISTENING" | findstr ":5433" >nul
if %errorlevel% equ 0 (
    echo [√] PostgreSQL (端口 5433) 正在监听
    set /a PASSED+=1
) else (
    echo [X] PostgreSQL (端口 5433) 未监听
    set /a FAILED+=1
)
set /a TOTAL+=1

netstat -ano | findstr "LISTENING" | findstr ":6380" >nul
if %errorlevel% equ 0 (
    echo [√] Redis (端口 6380) 正在监听
    set /a PASSED+=1
) else (
    echo [X] Redis (端口 6380) 未监听
    set /a FAILED+=1
)
set /a TOTAL+=1

netstat -ano | findstr "LISTENING" | findstr ":8080" >nul
if %errorlevel% equ 0 (
    echo [√] 后端 API (端口 8080) 正在监听
    set /a PASSED+=1
) else (
    echo [X] 后端 API (端口 8080) 未监听
    set /a FAILED+=1
)
set /a TOTAL+=1

netstat -ano | findstr "LISTENING" | findstr ":5173" >nul
if %errorlevel% equ 0 (
    echo [√] Admin 前端 (端口 5173) 正在监听
    set /a PASSED+=1
) else (
    echo [X] Admin 前端 (端口 5173) 未监听
    set /a FAILED+=1
)
set /a TOTAL+=1

netstat -ano | findstr "LISTENING" | findstr ":3000" >nul
if %errorlevel% equ 0 (
    echo [√] App 前端 (端口 3000) 正在监听
    set /a PASSED+=1
) else (
    echo [X] App 前端 (端口 3000) 未监听
    set /a FAILED+=1
)
set /a TOTAL+=1
echo.

REM 检查磁盘空间
echo === 系统资源检查 ===
for /f "tokens=3" %%a in ('dir /-c ^| find "bytes free"') do set FREE=%%a
echo 磁盘空间检查完成
set /a PASSED+=1
set /a TOTAL+=1
echo.

REM 容器资源使用
echo === 容器资源使用 ===
docker stats --no-stream --format "table {{.Container}}\t{{.CPUPerc}}\t{{.MemUsage}}" 2>nul
if %errorlevel% equ 0 (
    echo [√] 容器资源使用获取成功
    set /a PASSED+=1
) else (
    echo [X] 无法获取容器资源使用
    set /a FAILED+=1
)
set /a TOTAL+=1
echo.

REM 汇总
echo =========================================
echo === 检查汇总 ===
echo 总检查项: %TOTAL%
echo 通过: %PASSED%
echo 失败: %FAILED%
echo =========================================

if %FAILED% gtr 0 (
    echo [X] 健康检查失败
    exit /b 1
) else (
    echo [√] 健康检查通过
    exit /b 0
)
