# GameLink 后端服务启动脚本

Write-Host "========================================" -ForegroundColor Cyan
Write-Host "GameLink Backend API" -ForegroundColor Cyan
Write-Host "========================================" -ForegroundColor Cyan
Write-Host ""

# 切换到api目录
if (-not (Test-Path "api\cmd\main.go")) {
    Write-Host "❌ 错误: 请在项目根目录中运行此脚本" -ForegroundColor Red
    exit 1
}

Write-Host "📦 启动方式:" -ForegroundColor Yellow
Write-Host "1. 开发模式 (直接运行)" -ForegroundColor White
Write-Host "2. Docker Compose" -ForegroundColor White
Write-Host "3. Docker 生产部署" -ForegroundColor White
Write-Host ""

$mode = Read-Host "请选择 (1-3, 默认: 1)"

if ([string]::IsNullOrEmpty($mode)) { $mode = "1" }

switch ($mode) {
    "1" {
        Write-Host ""
        Write-Host "🚀 启动开发服务器..." -ForegroundColor Green
        Write-Host ""

        Set-Location api

        # 检查是否需要安装依赖
        if (-not (Test-Path "go.sum")) {
            Write-Host "📦 安装Go依赖..." -ForegroundColor Yellow
            go mod download
        }

        # 启动服务
        go run cmd/main.go
    }

    "2" {
        Write-Host ""
        Write-Host "🐳 启动Docker Compose..." -ForegroundColor Green
        Write-Host ""

        docker-compose up -d

        Write-Host ""
        Write-Host "✅ Docker服务已启动!" -ForegroundColor Green
        Write-Host ""
        Write-Host "📊 查看日志: docker-compose logs -f api" -ForegroundColor White
        Write-Host "🔍 查看状态: docker-compose ps" -ForegroundColor White
    }

    "3" {
        Write-Host ""
        Write-Host "🐳 Docker 生产部署..." -ForegroundColor Green
        Write-Host ""

        # 使用生产部署脚本
        if (Test-Path "scripts\deploy-production-encrypted.ps1") {
            & .\scripts\deploy-production-encrypted.ps1
        } else {
            Write-Host "❌ 生产部署脚本不存在" -ForegroundColor Red
            exit 1
        }
    }

    default {
        Write-Host "❌ 无效的选择" -ForegroundColor Red
        exit 1
    }
}
