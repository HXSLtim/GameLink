<#
.SYNOPSIS
    GameLink 生产环境 Docker 启动脚本
.DESCRIPTION
    自动检查环境变量、构建镜像并启动 GameLink 生产环境
.EXAMPLE
    .\scripts\docker-prod-start.ps1
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

# 检查 .env 文件
function Test-EnvFile {
    Write-ColorOutput "`n🔍 检查环境变量配置..." "Cyan"

    if (-not (Test-Path ".env")) {
        Write-ColorOutput "❌ .env 文件不存在" "Red"
        Write-ColorOutput "💡 请复制 .env.example 为 .env 并配置相应的值" "Yellow"
        Write-ColorOutput "   Copy-Item .env.example .env" "White"
        return $false
    }

    # 检查必需的环境变量
    $requiredVars = @(
        "JWT_SECRET_KEY",
        "CRYPTO_SECRET_KEY",
        "CRYPTO_IV",
        "SUPER_ADMIN_PASSWORD"
    )

    $envContent = Get-Content ".env" -Raw
    $missingVars = @()

    foreach ($var in $requiredVars) {
        if ($envContent -notmatch "$var=.+") {
            $missingVars += $var
        }
    }

    if ($missingVars.Count -gt 0) {
        Write-ColorOutput "❌ 以下必需的环境变量未设置:" "Red"
        foreach ($var in $missingVars) {
            Write-ColorOutput "   - $var" "Yellow"
        }
        Write-ColorOutput "`n💡 请编辑 .env 文件并设置这些变量" "Yellow"
        return $false
    }

    Write-ColorOutput "✅ 环境变量配置完整" "Green"
    return $true
}

# 检查 Docker
function Test-Docker {
    Write-ColorOutput "`n🔍 检查 Docker 环境..." "Cyan"

    try {
        $dockerVersion = docker --version
        Write-ColorOutput "✅ Docker 已安装: $dockerVersion" "Green"

        $composeVersion = docker-compose --version
        Write-ColorOutput "✅ Docker Compose 已安装: $composeVersion" "Green"

        docker info | Out-Null
        Write-ColorOutput "✅ Docker 服务正在运行" "Green"

        return $true
    }
    catch {
        Write-ColorOutput "❌ Docker 未安装或未运行" "Red"
        return $false
    }
}

# 主函数
function Start-GameLinkProd {
    Write-ColorOutput "`n🚀 GameLink 生产环境启动脚本" "Yellow"
    Write-ColorOutput "================================`n" "Yellow"

    # 检查 Docker
    if (-not (Test-Docker)) {
        exit 1
    }

    # 检查环境变量
    if (-not (Test-EnvFile)) {
        exit 1
    }

    # 安全确认
    Write-ColorOutput "`n⚠️  警告: 即将启动生产环境" "Yellow"
    $confirm = Read-Host "是否继续? (yes/no)"
    if ($confirm -ne "yes") {
        Write-ColorOutput "已取消启动" "Yellow"
        exit 0
    }

    # 停止现有容器
    Write-ColorOutput "`n🛑 停止现有容器..." "Cyan"
    docker-compose -f docker-compose.prod.yml down

    # 构建镜像
    Write-ColorOutput "`n🔨 构建 Docker 镜像..." "Cyan"
    docker-compose -f docker-compose.prod.yml build

    if ($LASTEXITCODE -ne 0) {
        Write-ColorOutput "❌ 镜像构建失败" "Red"
        exit 1
    }

    # 启动服务
    Write-ColorOutput "`n▶️  启动服务..." "Cyan"
    docker-compose -f docker-compose.prod.yml up -d

    if ($LASTEXITCODE -ne 0) {
        Write-ColorOutput "❌ 服务启动失败" "Red"
        exit 1
    }

    # 等待服务就绪
    Write-ColorOutput "`n⏳ 等待服务启动..." "Cyan"
    Start-Sleep -Seconds 15

    # 检查服务状态
    Write-ColorOutput "`n📊 服务状态:" "Cyan"
    docker-compose -f docker-compose.prod.yml ps

    # 显示访问信息
    Write-ColorOutput "`n✅ GameLink 生产环境启动成功!" "Green"
    Write-ColorOutput "`n📍 访问地址:" "Yellow"
    Write-ColorOutput "   前端应用:     http://localhost" "White"
    Write-ColorOutput "   后端API:      http://localhost:8080" "White"

    Write-ColorOutput "`n📝 常用命令:" "Yellow"
    Write-ColorOutput "   查看日志:     docker-compose -f docker-compose.prod.yml logs -f" "White"
    Write-ColorOutput "   停止服务:     docker-compose -f docker-compose.prod.yml down" "White"
    Write-ColorOutput "   重启服务:     docker-compose -f docker-compose.prod.yml restart" "White"

    Write-ColorOutput "`n⚠️  安全提醒:" "Yellow"
    Write-ColorOutput "   1. 确保修改了所有默认密码" "White"
    Write-ColorOutput "   2. 生产环境建议启用 HTTPS" "White"
    Write-ColorOutput "   3. 定期备份数据库" "White"
    Write-ColorOutput "   4. 监控系统日志" "White"

    Write-ColorOutput "`n" "White"
}

# 运行主函数
Start-GameLinkProd
