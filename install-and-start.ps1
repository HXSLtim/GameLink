# GameLink 自动安装和启动脚本 (PowerShell)
# 用途: 安装PWA依赖并启动服务

$ErrorActionPreference = "Stop"

function Write-Header {
    param([string]$Text)
    Write-Host "`n========================================" -ForegroundColor Cyan
    Write-Host $Text -ForegroundColor Cyan
    Write-Host "========================================`n" -ForegroundColor Cyan
}

function Write-Step {
    param([string]$Text)
    Write-Host "`n📋 $Text" -ForegroundColor Yellow
}

function Write-Success {
    param([string]$Text)
    Write-Host "`n✅ $Text`n" -ForegroundColor Green
}

function Write-Error {
    param([string]$Text)
    Write-Host "`n❌ $Text`n" -ForegroundColor Red
}

# 检查是否在项目根目录
if (-not (Test-Path "admin\package.json")) {
    Write-Error "错误: 请在项目根目录中运行此脚本"
    Read-Host "按Enter键退出"
    exit 1
}

Write-Header "GameLink 自动安装和启动"

# 步骤 1: 安装PWA依赖
Write-Step "步骤 1: 安装前端 PWA 依赖"

try {
    Set-Location admin

    Write-Host "正在安装: vite-plugin-pwa, picocolors, workbox-window..." -ForegroundColor White

    $result = npm install -D vite-plugin-pwa picocolors workbox-window 2>&1

    if ($LASTEXITCODE -eq 0) {
        Write-Success "依赖安装成功!"
    } else {
        Write-Error "依赖安装失败!"
        Write-Host "错误信息:" -ForegroundColor Red
        Write-Host $result

        Write-Host "`n💡 可能的解决方案:" -ForegroundColor Yellow
        Write-Host "1. 检查网络连接" -ForegroundColor White
        Write-Host "2. 尝试使用淘宝镜像:" -ForegroundColor White
        Write-Host "   npm config set registry https://registry.npmmirror.com" -ForegroundColor Gray
        Write-Host "3. 手动删除 node_modules 后重新安装" -ForegroundColor White

        Read-Host "按Enter键退出"
        exit 1
    }

    Set-Location ..
}
catch {
    Write-Error "安装过程出错: $_"
    Read-Host "按Enter键退出"
    exit 1
}

# 步骤 2: 选择启动方式
Write-Step "步骤 2: 选择要启动的服务"

Write-Host "`n请选择:" -ForegroundColor Cyan
Write-Host "1. 仅启动前端开发服务器" -ForegroundColor White
Write-Host "2. 仅启动后端API服务 (开发模式)" -ForegroundColor White
Write-Host "3. 仅启动后端API服务 (Docker)" -ForegroundColor White
Write-Host "4. 同时启动前后端 (推荐)" -ForegroundColor White
Write-Host ""

$choice = Read-Host "请输入选择 (1-4, 默认: 4)"

if ([string]::IsNullOrEmpty($choice)) { $choice = "4" }

switch ($choice) {
    "1" {
        Write-Step "启动前端开发服务器"
        Set-Location admin
        npm run dev
    }

    "2" {
        Write-Step "启动后端API服务 (开发模式)"
        Set-Location api
        go run cmd/main.go
    }

    "3" {
        Write-Step "启动后端API服务 (Docker)"

        Write-Host "检查Docker状态..." -ForegroundColor Yellow
        docker --version

        if ($LASTEXITCODE -ne 0) {
            Write-Error "Docker未安装或未运行"
            Read-Host "按Enter键退出"
            exit 1
        }

        docker-compose up -d

        Write-Success "Docker服务已启动!"
        Write-Host "`n📊 查看日志: docker-compose logs -f" -ForegroundColor White
        Write-Host "🔍 查看状态: docker-compose ps" -ForegroundColor White
    }

    "4" {
        Write-Step "同时启动前后端服务"

        Write-Host "`n正在后台启动后端API服务..." -ForegroundColor Yellow

        # 使用Start-Process在新窗口启动后端
        $backendPath = Join-Path (Get-Location).Path "api"
        Start-Process powershell -ArgumentList "-NoExit", "-Command", "cd '$backendPath'; go run cmd/main.go"

        Write-Host "后端服务已在后台启动" -ForegroundColor Green

        Start-Sleep -Seconds 2

        Write-Host "`n正在启动前端开发服务器..." -ForegroundColor Yellow
        Set-Location admin
        npm run dev
    }

    default {
        Write-Error "无效的选择"
        Read-Host "按Enter键退出"
        exit 1
    }
}
