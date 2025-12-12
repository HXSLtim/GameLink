<#
.SYNOPSIS
    GameLink Docker 数据恢复脚本
.DESCRIPTION
    从备份文件恢复 PostgreSQL 数据库和 Redis 数据
.PARAMETER BackupFile
    备份文件路径（.zip 文件）
.PARAMETER Environment
    环境类型: prod-local 或 prod
.EXAMPLE
    .\scripts\docker-restore.ps1 -BackupFile ".\backups\20250113_120000.zip" -Environment prod-local
#>

param(
    [Parameter(Mandatory=$true)]
    [string]$BackupFile,
    
    [Parameter(Mandatory=$false)]
    [ValidateSet("prod-local", "prod")]
    [string]$Environment = "prod-local"
)

$ErrorActionPreference = "Stop"

function Write-ColorOutput {
    param([string]$Message, [string]$Color = "White")
    Write-Host $Message -ForegroundColor $Color
}

function Restore-GameLink {
    Write-ColorOutput "`n♻️  GameLink 数据恢复工具" "Yellow"
    Write-ColorOutput "========================`n" "Yellow"

    # 检查备份文件
    if (-not (Test-Path $BackupFile)) {
        Write-ColorOutput "❌ 备份文件不存在: $BackupFile" "Red"
        exit 1
    }

    Write-ColorOutput "📦 备份文件: $BackupFile" "Cyan"

    # 确认恢复操作
    Write-ColorOutput "`n⚠️  警告: 此操作将覆盖现有数据!" "Red"
    $confirm = Read-Host "确认恢复? (yes/no)"
    if ($confirm -ne "yes") {
        Write-ColorOutput "已取消恢复" "Yellow"
        exit 0
    }

    # 解压备份文件
    Write-ColorOutput "`n📂 解压备份文件..." "Cyan"
    $tempDir = Join-Path $env:TEMP "gamelink_restore_$(Get-Date -Format 'yyyyMMddHHmmss')"
    Expand-Archive -Path $BackupFile -DestinationPath $tempDir -Force
    Write-ColorOutput "✅ 解压完成" "Green"

    # 确定 docker-compose 文件
    $composeFile = if ($Environment -eq "prod-local") {
        "docker-compose.prod.local.yml"
    } else {
        "docker-compose.prod.yml"
    }

    $envFile = if ($Environment -eq "prod-local") {
        ".env.production.local"
    } else {
        ".env"
    }

    # 检查容器是否运行
    Write-ColorOutput "`n🔍 检查容器状态..." "Cyan"
    $containers = docker-compose --env-file $envFile -f $composeFile ps -q
    if (-not $containers) {
        Write-ColorOutput "❌ 容器未运行，请先启动服务" "Red"
        Remove-Item -Path $tempDir -Recurse -Force
        exit 1
    }

    # 恢复 PostgreSQL
    Write-ColorOutput "`n📥 恢复 PostgreSQL 数据库..." "Cyan"
    $pgBackupFile = Get-ChildItem -Path $tempDir -Filter "postgres_backup.sql" -Recurse | Select-Object -First 1
    
    if ($pgBackupFile) {
        try {
            # 删除现有数据库并重新创建
            docker exec gamelink-postgres psql -U gamelink -d postgres -c "DROP DATABASE IF EXISTS gamelink;"
            docker exec gamelink-postgres psql -U gamelink -d postgres -c "CREATE DATABASE gamelink;"
            
            # 恢复数据
            Get-Content $pgBackupFile.FullName | docker exec -i gamelink-postgres psql -U gamelink -d gamelink
            Write-ColorOutput "✅ PostgreSQL 恢复完成" "Green"
        }
        catch {
            Write-ColorOutput "❌ PostgreSQL 恢复失败: $_" "Red"
        }
    }
    else {
        Write-ColorOutput "⚠️  未找到 PostgreSQL 备份文件" "Yellow"
    }

    # 恢复 Redis
    Write-ColorOutput "`n📥 恢复 Redis 数据..." "Cyan"
    $redisBackupFile = Get-ChildItem -Path $tempDir -Filter "redis_dump.rdb" -Recurse | Select-Object -First 1
    
    if ($redisBackupFile) {
        try {
            # 停止 Redis 写入
            docker exec gamelink-redis redis-cli -a redis123 SHUTDOWN NOSAVE 2>&1 | Out-Null
            Start-Sleep -Seconds 2
            
            # 复制备份文件
            docker cp $redisBackupFile.FullName gamelink-redis:/data/dump.rdb
            
            # 重启 Redis 容器
            docker-compose --env-file $envFile -f $composeFile restart redis
            Start-Sleep -Seconds 5
            
            Write-ColorOutput "✅ Redis 恢复完成" "Green"
        }
        catch {
            Write-ColorOutput "❌ Redis 恢复失败: $_" "Red"
            Write-ColorOutput "💡 尝试重启 Redis: docker-compose --env-file $envFile -f $composeFile restart redis" "Yellow"
        }
    }
    else {
        Write-ColorOutput "⚠️  未找到 Redis 备份文件" "Yellow"
    }

    # 清理临时文件
    Remove-Item -Path $tempDir -Recurse -Force

    Write-ColorOutput "`n✅ 恢复完成!" "Green"
    Write-ColorOutput "`n💡 建议重启后端服务以确保连接正常:" "Yellow"
    Write-ColorOutput "   docker-compose --env-file $envFile -f $composeFile restart backend" "White"
}

Restore-GameLink
