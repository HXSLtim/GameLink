<#
.SYNOPSIS
    GameLink Docker 健康检查脚本
.DESCRIPTION
    检查所有服务的健康状态和连接性
.PARAMETER Environment
    环境类型: dev, prod-local 或 prod
.EXAMPLE
    .\scripts\docker-health-check.ps1 -Environment prod-local
#>

param(
    [Parameter(Mandatory=$false)]
    [ValidateSet("dev", "prod-local", "prod")]
    [string]$Environment = "prod-local"
)

$ErrorActionPreference = "Continue"

function Write-ColorOutput {
    param([string]$Message, [string]$Color = "White")
    Write-Host $Message -ForegroundColor $Color
}

function Test-ServiceHealth {
    param(
        [string]$ServiceName,
        [string]$Url,
        [int]$Timeout = 5
    )

    try {
        $response = Invoke-WebRequest -Uri $Url -TimeoutSec $Timeout -ErrorAction Stop
        if ($response.StatusCode -eq 200) {
            Write-ColorOutput "  ✅ $ServiceName - 健康" "Green"
            return $true
        }
        else {
            Write-ColorOutput "  ⚠️  $ServiceName - 状态码: $($response.StatusCode)" "Yellow"
            return $false
        }
    }
    catch {
        Write-ColorOutput "  ❌ $ServiceName - 无法访问" "Red"
        return $false
    }
}

function Test-DatabaseConnection {
    param([string]$ContainerName)

    try {
        $result = docker exec $ContainerName pg_isready -U gamelink 2>&1
        if ($LASTEXITCODE -eq 0) {
            Write-ColorOutput "  ✅ PostgreSQL - 连接正常" "Green"
            
            # 检查数据库大小
            $dbSize = docker exec $ContainerName psql -U gamelink -d gamelink -t -c "SELECT pg_size_pretty(pg_database_size('gamelink'));" 2>&1
            Write-ColorOutput "     数据库大小: $($dbSize.Trim())" "Cyan"
            return $true
        }
        else {
            Write-ColorOutput "  ❌ PostgreSQL - 连接失败" "Red"
            return $false
        }
    }
    catch {
        Write-ColorOutput "  ❌ PostgreSQL - 检查失败: $_" "Red"
        return $false
    }
}

function Test-RedisConnection {
    param([string]$ContainerName, [string]$Password)

    try {
        $result = docker exec $ContainerName redis-cli -a $Password PING 2>&1
        if ($result -match "PONG") {
            Write-ColorOutput "  ✅ Redis - 连接正常" "Green"
            
            # 检查 Redis 信息
            $info = docker exec $ContainerName redis-cli -a $Password INFO stats 2>&1 | Select-String "total_commands_processed"
            Write-ColorOutput "     $info" "Cyan"
            return $true
        }
        else {
            Write-ColorOutput "  ❌ Redis - 连接失败" "Red"
            return $false
        }
    }
    catch {
        Write-ColorOutput "  ❌ Redis - 检查失败: $_" "Red"
        return $false
    }
}

function Get-ContainerStats {
    param([string]$ContainerName)

    try {
        $stats = docker stats $ContainerName --no-stream --format "table {{.CPUPerc}}\t{{.MemUsage}}" 2>&1
        if ($stats) {
            $lines = $stats -split "`n"
            if ($lines.Count -gt 1) {
                Write-ColorOutput "     资源使用: $($lines[1])" "Cyan"
            }
        }
    }
    catch {
        # 忽略错误
    }
}

function Check-GameLinkHealth {
    Write-ColorOutput "`n🏥 GameLink 健康检查" "Yellow"
    Write-ColorOutput "===================`n" "Yellow"

    # 确定配置
    $composeFile = switch ($Environment) {
        "dev" { "docker-compose.yml" }
        "prod-local" { "docker-compose.prod.local.yml" }
        "prod" { "docker-compose.prod.yml" }
    }

    $envFile = switch ($Environment) {
        "dev" { ".env.development" }
        "prod-local" { ".env.production.local" }
        "prod" { ".env" }
    }

    $backendPort = switch ($Environment) {
        "dev" { "8080" }
        "prod-local" { "8081" }
        "prod" { "8080" }
    }

    Write-ColorOutput "📋 环境: $Environment" "Cyan"
    Write-ColorOutput "📋 配置文件: $composeFile`n" "Cyan"

    # 检查容器状态
    Write-ColorOutput "🔍 检查容器状态..." "Cyan"
    $containers = docker-compose -f $composeFile ps 2>&1
    
    if ($LASTEXITCODE -ne 0) {
        Write-ColorOutput "❌ 无法获取容器状态，服务可能未启动" "Red"
        exit 1
    }

    Write-Host $containers
    Write-Host ""

    # 检查后端服务
    Write-ColorOutput "🔍 检查后端服务..." "Cyan"
    $backendHealthy = Test-ServiceHealth -ServiceName "后端API" -Url "http://localhost:$backendPort/api/v1/health"
    Get-ContainerStats -ContainerName "gamelink-backend"
    Write-Host ""

    # 检查 Swagger 文档
    Write-ColorOutput "🔍 检查 Swagger 文档..." "Cyan"
    Test-ServiceHealth -ServiceName "Swagger文档" -Url "http://localhost:$backendPort/swagger/index.html" | Out-Null
    Write-Host ""

    # 检查数据库（生产环境）
    if ($Environment -ne "dev") {
        Write-ColorOutput "🔍 检查 PostgreSQL..." "Cyan"
        $dbHealthy = Test-DatabaseConnection -ContainerName "gamelink-postgres"
        Get-ContainerStats -ContainerName "gamelink-postgres"
        Write-Host ""

        Write-ColorOutput "🔍 检查 Redis..." "Cyan"
        $redisHealthy = Test-RedisConnection -ContainerName "gamelink-redis" -Password "redis123"
        Get-ContainerStats -ContainerName "gamelink-redis"
        Write-Host ""
    }

    # 检查管理后台服务（如果存在）
    $adminContainer = docker ps --filter "name=gamelink-admin" --format "{{.Names}}" 2>&1
    if ($adminContainer -eq "gamelink-admin") {
        Write-ColorOutput "🔍 检查管理后台服务..." "Cyan"
        Test-ServiceHealth -ServiceName "管理后台应用" -Url "http://localhost/health" | Out-Null
        Get-ContainerStats -ContainerName "gamelink-admin"
        Write-Host ""
    }

    # 检查网络连接
    Write-ColorOutput "🔍 检查网络配置..." "Cyan"
    $network = docker network inspect gamelink-network 2>&1
    if ($LASTEXITCODE -eq 0) {
        $networkJson = $network | ConvertFrom-Json
        $containerCount = $networkJson[0].Containers.Count
        Write-ColorOutput "  ✅ 网络正常 - $containerCount 个容器已连接" "Green"
    }
    else {
        Write-ColorOutput "  ❌ 网络检查失败" "Red"
    }
    Write-Host ""

    # 检查数据卷
    Write-ColorOutput "🔍 检查数据卷..." "Cyan"
    $volumes = docker volume ls --filter "name=gamelink" --format "{{.Name}}" 2>&1
    if ($volumes) {
        foreach ($vol in $volumes) {
            $volInfo = docker volume inspect $vol 2>&1 | ConvertFrom-Json
            Write-ColorOutput "  ✅ $vol" "Green"
        }
    }
    else {
        Write-ColorOutput "  ⚠️  未找到数据卷" "Yellow"
    }
    Write-Host ""

    # 总结
    Write-ColorOutput "📊 健康检查总结" "Yellow"
    Write-ColorOutput "===============" "Yellow"
    
    $allHealthy = $backendHealthy
    if ($Environment -ne "dev") {
        $allHealthy = $allHealthy -and $dbHealthy -and $redisHealthy
    }

    if ($allHealthy) {
        Write-ColorOutput "✅ 所有核心服务运行正常" "Green"
    }
    else {
        Write-ColorOutput "⚠️  部分服务存在问题，请查看上方详情" "Yellow"
    }

    Write-ColorOutput "`n💡 常用命令:" "Cyan"
    Write-ColorOutput "   查看日志: docker-compose -f $composeFile logs -f" "White"
    Write-ColorOutput "   重启服务: docker-compose -f $composeFile restart" "White"
    Write-ColorOutput "   查看资源: docker stats" "White"
}

Check-GameLinkHealth
