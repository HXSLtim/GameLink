<#
.SYNOPSIS
    GameLink 本地生产环境测试启动脚本
.DESCRIPTION
    在本地启动完整的生产环境配置（PostgreSQL + Redis），用于生产环境测试
.EXAMPLE
    .\scripts\docker-prod-local-start.ps1
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
        Write-ColorOutput "❌ Docker 未安装或未运行，请先安装 Docker Desktop" "Red"
        return $false
    }
}

# 检查环境变量文件
function Test-EnvFile {
    Write-ColorOutput "`n🔍 检查环境变量配置..." "Cyan"

    if (-not (Test-Path ".env.production.local")) {
        Write-ColorOutput "⚠️  .env.production.local 文件不存在，将使用默认配置" "Yellow"
        Write-ColorOutput "💡 如需自定义配置，请创建 .env.production.local 文件" "Yellow"
        return $true
    }

    Write-ColorOutput "✅ 找到 .env.production.local 配置文件" "Green"
    return $true
}

# 检查端口占用
function Test-Ports {
    Write-ColorOutput "`n🔍 检查端口占用..." "Cyan"

    $ports = @{
        "8081" = "后端服务"
        "5433" = "PostgreSQL"
        "6380" = "Redis"
    }

    $portsInUse = @()

    foreach ($port in $ports.Keys) {
        $connection = Get-NetTCPConnection -LocalPort $port -ErrorAction SilentlyContinue
        if ($connection) {
            $portsInUse += "$port ($($ports[$port]))"
        }
    }

    if ($portsInUse.Count -gt 0) {
        Write-ColorOutput "⚠️  以下端口已被占用:" "Yellow"
        foreach ($portInfo in $portsInUse) {
            Write-ColorOutput "   - $portInfo" "Yellow"
        }
        Write-ColorOutput "`n💡 请停止占用端口的程序，或修改 docker-compose.prod.local.yml 中的端口映射" "Yellow"
        
        $continue = Read-Host "`n是否继续启动? (yes/no)"
        if ($continue -ne "yes") {
            return $false
        }
    }
    else {
        Write-ColorOutput "✅ 所有端口可用" "Green"
    }

    return $true
}

# 清理旧容器和数据
function Clear-OldData {
    param([bool]$Force = $false)

    if ($Force) {
        Write-ColorOutput "`n🧹 清理旧数据..." "Cyan"
        docker-compose --env-file .env.production.local -f docker-compose.prod.local.yml down -v
        Write-ColorOutput "✅ 旧数据已清理" "Green"
    }
    else {
        Write-ColorOutput "`n🛑 停止现有容器..." "Cyan"
        docker-compose --env-file .env.production.local -f docker-compose.prod.local.yml down
    }
}

# 主函数
function Start-GameLinkProdLocal {
    param(
        [switch]$Clean,
        [switch]$NoBuild
    )

    Write-ColorOutput "`n🚀 GameLink 本地生产环境启动脚本" "Yellow"
    Write-ColorOutput "====================================`n" "Yellow"

    # 检查 Docker
    if (-not (Test-Docker)) {
        exit 1
    }

    # 检查环境变量
    if (-not (Test-EnvFile)) {
        exit 1
    }

    # 检查端口
    if (-not (Test-Ports)) {
        exit 1
    }

    # 清理旧数据
    if ($Clean) {
        Write-ColorOutput "`n⚠️  警告: 将删除所有数据（包括数据库）" "Red"
        $confirm = Read-Host "确认删除? (yes/no)"
        if ($confirm -eq "yes") {
            Clear-OldData -Force $true
        }
        else {
            Write-ColorOutput "已取消清理" "Yellow"
            exit 0
        }
    }
    else {
        Clear-OldData
    }

    # 构建镜像
    if (-not $NoBuild) {
        Write-ColorOutput "`n🔨 构建 Docker 镜像..." "Cyan"
        docker-compose --env-file .env.production.local -f docker-compose.prod.local.yml build

        if ($LASTEXITCODE -ne 0) {
            Write-ColorOutput "❌ 镜像构建失败" "Red"
            exit 1
        }
        Write-ColorOutput "✅ 镜像构建成功" "Green"
    }

    # 启动服务
    Write-ColorOutput "`n▶️  启动服务..." "Cyan"
    docker-compose --env-file .env.production.local -f docker-compose.prod.local.yml up -d

    if ($LASTEXITCODE -ne 0) {
        Write-ColorOutput "❌ 服务启动失败" "Red"
        Write-ColorOutput "`n💡 查看日志: docker-compose --env-file .env.production.local -f docker-compose.prod.local.yml logs" "Yellow"
        exit 1
    }

    # 等待服务就绪
    Write-ColorOutput "`n⏳ 等待服务启动（约30秒）..." "Cyan"
    
    $maxWait = 60
    $waited = 0
    $interval = 5

    while ($waited -lt $maxWait) {
        Start-Sleep -Seconds $interval
        $waited += $interval

        # 检查后端健康状态
        try {
            $response = Invoke-WebRequest -Uri "http://localhost:8081/api/v1/health" -TimeoutSec 2 -ErrorAction SilentlyContinue
            if ($response.StatusCode -eq 200) {
                Write-ColorOutput "✅ 后端服务已就绪" "Green"
                break
            }
        }
        catch {
            Write-Host "." -NoNewline
        }
    }

    Write-Host ""

    # 检查服务状态
    Write-ColorOutput "`n📊 服务状态:" "Cyan"
    docker-compose --env-file .env.production.local -f docker-compose.prod.local.yml ps

    # 显示访问信息
    Write-ColorOutput "`n✅ GameLink 本地生产环境启动成功!" "Green"
    Write-ColorOutput "`n📍 访问地址:" "Yellow"
    Write-ColorOutput "   后端API:      http://localhost:8081" "White"
    Write-ColorOutput "   Swagger文档:  http://localhost:8081/swagger/index.html" "White"
    Write-ColorOutput "   PostgreSQL:   localhost:5433" "White"
    Write-ColorOutput "   Redis:        localhost:6380" "White"

    Write-ColorOutput "`n👤 默认管理员账号:" "Yellow"
    Write-ColorOutput "   邮箱: admin@gamelink.com" "White"
    Write-ColorOutput "   密码: admin123456" "White"

    Write-ColorOutput "`n🔑 数据库连接信息:" "Yellow"
    Write-ColorOutput "   数据库: gamelink" "White"
    Write-ColorOutput "   用户名: gamelink" "White"
    Write-ColorOutput "   密码:   gamelink123" "White"
    Write-ColorOutput "   端口:   5433" "White"

    Write-ColorOutput "`n📝 常用命令:" "Yellow"
    Write-ColorOutput "   查看日志:     docker-compose --env-file .env.production.local -f docker-compose.prod.local.yml logs -f" "White"
    Write-ColorOutput "   查看后端日志: docker-compose --env-file .env.production.local -f docker-compose.prod.local.yml logs -f backend" "White"
    Write-ColorOutput "   停止服务:     docker-compose --env-file .env.production.local -f docker-compose.prod.local.yml down" "White"
    Write-ColorOutput "   重启服务:     docker-compose --env-file .env.production.local -f docker-compose.prod.local.yml restart" "White"
    Write-ColorOutput "   进入数据库:   docker exec -it gamelink-postgres psql -U gamelink -d gamelink" "White"
    Write-ColorOutput "   进入Redis:    docker exec -it gamelink-redis redis-cli -a redis123" "White"

    Write-ColorOutput "`n💡 提示:" "Yellow"
    Write-ColorOutput "   - 使用 -Clean 参数清理所有数据重新开始" "White"
    Write-ColorOutput "   - 使用 -NoBuild 参数跳过镜像构建" "White"
    Write-ColorOutput "   - 端口已调整避免与本地服务冲突（8081, 5433, 6380）" "White"

    Write-ColorOutput "`n" "White"
}

# 解析参数并运行
$params = @{}
if ($args -contains "-Clean") { $params.Clean = $true }
if ($args -contains "-NoBuild") { $params.NoBuild = $true }

Start-GameLinkProdLocal @params
