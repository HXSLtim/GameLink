# GameLink 完整服务启动脚本
# 启动前端开发服务器和后端API服务

Write-Host "========================================" -ForegroundColor Cyan
Write-Host "GameLink 服务启动" -ForegroundColor Cyan
Write-Host "========================================" -ForegroundColor Cyan
Write-Host ""

# 检查是否在项目根目录
if (-not (Test-Path "admin\package.json") -or -not (Test-Path "api\go.mod")) {
    Write-Host "❌ 错误: 请在项目根目录中运行此脚本" -ForegroundColor Red
    exit 1
}

Write-Host "📋 启动选项:" -ForegroundColor Yellow
Write-Host "1. 仅启动前端开发服务器" -ForegroundColor White
Write-Host "2. 仅启动后端API服务" -ForegroundColor White
Write-Host "3. 同时启动前后端服务 (推荐)" -ForegroundColor White
Write-Host "4. 部署后端Docker服务" -ForegroundColor White
Write-Host ""

$choice = Read-Host "请选择 (1-4, 默认: 3)"

if ([string]::IsNullOrEmpty($choice)) { $choice = "3" }

switch ($choice) {
    "1" {
        Write-Host ""
        Write-Host "🚀 启动前端开发服务器..." -ForegroundColor Green
        Write-Host ""

        Set-Location admin
        npm run dev
    }

    "2" {
        Write-Host ""
        Write-Host "🚀 启动后端API服务..." -ForegroundColor Green
        Write-Host ""

        Set-Location api
        go run cmd/main.go
    }

    "3" {
        Write-Host ""
        Write-Host "🚀 同时启动前后端服务..." -ForegroundColor Green
        Write-Host ""

        # 启动后端
        $backendJob = Start-Job -ScriptBlock {
            Set-Location $args[0]
            go run cmd/main.go
        } -ArgumentList (Get-Location).Path + "\api"

        # 启动前端
        Set-Location admin
        npm run dev

        # 清理
        Stop-Job $backendJob
        Remove-Job $backendJob
    }

    "4" {
        Write-Host ""
        Write-Host "🐳 部署后端Docker服务..." -ForegroundColor Green
        Write-Host ""

        Write-Host "📦 检查Docker状态..." -ForegroundColor Yellow
        docker --version

        if ($LASTEXITCODE -ne 0) {
            Write-Host "❌ Docker未安装或未运行" -ForegroundColor Red
            exit 1
        }

        Write-Host ""
        Write-Host "🐋 启动Docker Compose服务..." -ForegroundColor Yellow
        docker-compose up -d

        Write-Host ""
        Write-Host "✅ Docker服务已启动!" -ForegroundColor Green
        Write-Host ""
        Write-Host "📊 查看日志: docker-compose logs -f" -ForegroundColor White
        Write-Host "🛑 停止服务: docker-compose down" -ForegroundColor White
        Write-Host "🔄 重启服务: docker-compose restart" -ForegroundColor White
    }

    default {
        Write-Host "❌ 无效的选择" -ForegroundColor Red
        exit 1
    }
}
