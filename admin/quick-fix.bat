@echo off
REM GameLink 快速修复安装问题

echo ========================================
echo GameLink 安装快速修复
echo ========================================
echo.

cd admin

echo 检测到安装问题，正在尝试修复...
echo.

REM 方案1: 使用淘宝镜像
echo 方案1: 配置淘宝镜像并安装...
echo.

call npm config set registry https://registry.npmmirror.com
call npm cache clean --force
call npm install -D vite-plugin-pwa picocolors workbox-window

if %errorlevel% equ 0 (
    echo.
    echo ========================================
    echo ✅ 安装成功!
    echo ========================================
    echo.
    echo 下一步运行: npm run dev
    echo.
    pause
    exit /b 0
)

echo.
echo 方案1失败，尝试方案2...
echo.

REM 方案2: 使用legacy模式
echo 方案2: 使用legacy-peer-deps模式...
echo.

call npm install -D vite-plugin-pwa picocolors workbox-window --legacy-peer-deps

if %errorlevel% equ 0 (
    echo.
    echo ========================================
    echo ✅ 安装成功!
    echo ========================================
    echo.
    echo 下一步运行: npm run dev
    echo.
    pause
    exit /b 0
)

echo.
echo 方案2也失败了，尝试方案3...
echo.

REM 方案3: 分步安装
echo 方案3: 分步安装每个包...
echo.

call npm install -D vite-plugin-pwa
if %errorlevel% neq 0 (
    echo vite-plugin-pwa 安装失败，继续尝试下一个...
)

call npm install -D picocolors
if %errorlevel% neq 0 (
    echo picocolors 安装失败，继续尝试下一个...
)

call npm install -D workbox-window
if %errorlevel% neq 0 (
    echo workbox-window 安装失败
)

echo.
echo ========================================
echo 安装尝试完成
echo ========================================
echo.

if %errorlevel% equ 0 (
    echo ✅ 部分或全部依赖可能已安装
    echo 尝试运行: npm run dev
) else (
    echo ❌ 所有方案都失败了
    echo.
    echo 建议手动执行以下命令:
    echo.
    echo cd admin
    echo npm config set registry https://registry.npmmirror.com
    echo npm install -D vite-plugin-pwa picocolors workbox-window
    echo.
    echo 或查看详细故障排除指南:
    echo INSTALL_TROUBLESHOOTING.md
)

echo.
pause
