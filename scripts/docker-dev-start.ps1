<#
.SYNOPSIS
    GameLink 开发环境 Docker 快速启动脚本
.DESCRIPTION
    自动检查环境、构建镜像并启动 GameLink 开发环境
.EXAMPLE
    .\scripts\docker-dev-start.ps1
#>

# 设置错误处理
$ErrorActionPreference = "Stop"

# 颜色输出函数
function Write-ColorOutput {
    param(
        [string]$Message,
        [string]$Color = "White"
    )
    Write-Host $Message -ForegroundColor $Color
}

# 检查 Docker 是否安装
function Test-Docker {
    Write-ColorOutput "`n🔍 检查 Docker 环境..." "Cyan"

    try {
        $dockerVersion = docker --version
        Write-ColorOutput "✅ Docker 已安装: $dockerVersion" "Green"

        $composeVersion = docker-compose --version
        Write-ColorOutput "✅ Docker Compose 已安装: $composeVersion" "Green"

        # 检查 Docker 服务是否运行
        docker info | Out-Null
        Write-ColorOutput "✅ Docker 服务正在运行" "Green"

        return $true
    }
    catch {
        Write-ColorOutput "❌ Docker 未安装或未运行，请先安装 Docker Desktop" "Red"
        return $false
    }
}

# 主函数
function Start-GameLinkDev {
    Write-ColorOutput "`n🚀 GameLink 开发环境启动脚本" "Yellow"
    Write-ColorOutput "================================`n" "Yellow"

    # 检查 Docker
    if (-not (Test-Docker)) {
        exit 1
    }

    # 停止现有容器
    Write-ColorOutput "`n🛑 停止现有容器..." "Cyan"
    docker-compose down

    # 构建镜像
    Write-ColorOutput "`n🔨 构建 Docker 镜像..." "Cyan"
    docker-compose build

    if ($LASTEXITCODE -ne 0) {
        Write-ColorOutput "❌ 镜像构建失败" "Red"
        exit 1
    }

    # 启动服务
    Write-ColorOutput "`n▶️  启动服务..." "Cyan"
    docker-compose up -d

    if ($LASTEXITCODE -ne 0) {
        Write-ColorOutput "❌ 服务启动失败" "Red"
        exit 1
    }

    # 等待服务就绪
    Write-ColorOutput "`n⏳ 等待服务启动..." "Cyan"
    Start-Sleep -Seconds 10

    # 检查服务状态
    Write-ColorOutput "`n📊 服务状态:" "Cyan"
    docker-compose ps

    # 显示访问信息
    Write-ColorOutput "`n✅ GameLink 开发环境启动成功!" "Green"
    Write-ColorOutput "`n📍 访问地址:" "Yellow"
    Write-ColorOutput "   前端应用:     http://localhost" "White"
    Write-ColorOutput "   后端API:      http://localhost:8080" "White"
    Write-ColorOutput "   Swagger文档:  http://localhost:8080/swagger/index.html" "White"

    Write-ColorOutput "`n👤 默认管理员账号:" "Yellow"
    Write-ColorOutput "   邮箱: admin@gamelink.com" "White"
    Write-ColorOutput "   密码: 123456" "White"

    Write-ColorOutput "`n📝 常用命令:" "Yellow"
    Write-ColorOutput "   查看日志:     docker-compose logs -f" "White"
    Write-ColorOutput "   停止服务:     docker-compose down" "White"
    Write-ColorOutput "   重启服务:     docker-compose restart" "White"

    Write-ColorOutput "`n" "White"
}

# 运行主函数
Start-GameLinkDev
