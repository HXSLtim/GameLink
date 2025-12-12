<#
.SYNOPSIS
    GameLink Docker 数据备份脚本
.DESCRIPTION
    备份 PostgreSQL 数据库和 Redis 数据
.PARAMETER Environment
    环境类型: prod-local 或 prod
.PARAMETER BackupDir
    备份目录路径
.EXAMPLE
    .\scripts\docker-backup.ps1 -Environment prod-local
#>

param(
    [Parameter(Mandatory=$false)]
    [ValidateSet("prod-local", "prod")]
    [string]$Environment = "prod-local",
    
    [Parameter(Mandatory=$false)]
    [string]$BackupDir = ".\backups"
)

$ErrorActionPreference = "Stop"

function Write-ColorOutput {
    param([string]$Message, [string]$Color = "White")
    Write-Host $Message -ForegroundColor $Color
}

function Backup-GameLink {
    Write-ColorOutput "`n💾 GameLink 数据备份工具" "Yellow"
    Write-ColorOutput "========================`n" "Yellow"

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
    Write-ColorOutput "🔍 检查容器状态..." "Cyan"
    $containers = docker-compose --env-file $envFile -f $composeFile ps -q
    if (-not $containers) {
        Write-ColorOutput "❌ 容器未运行，请先启动服务" "Red"
        exit 1
    }

    # 创建备份目录
    $timestamp = Get-Date -Format "yyyyMMdd_HHmmss"
    $backupPath = Join-Path $BackupDir $timestamp
    New-Item -ItemType Directory -Path $backupPath -Force | Out-Null
    Write-ColorOutput "✅ 备份目录: $backupPath" "Green"

    # 备份 PostgreSQL
    Write-ColorOutput "`n📦 备份 PostgreSQL 数据库..." "Cyan"
    $pgBackupFile = Join-Path $backupPath "postgres_backup.sql"
    
    try {
        docker exec gamelink-postgres pg_dump -U gamelink gamelink > $pgBackupFile
        $fileSize = (Get-Item $pgBackupFile).Length / 1MB
        Write-ColorOutput "✅ PostgreSQL 备份完成 (${fileSize:N2} MB)" "Green"
    }
    catch {
        Write-ColorOutput "❌ PostgreSQL 备份失败: $_" "Red"
    }

    # 备份 Redis
    Write-ColorOutput "`n📦 备份 Redis 数据..." "Cyan"
    try {
        # 触发 Redis 保存
        docker exec gamelink-redis redis-cli -a redis123 SAVE 2>&1 | Out-Null
        
        # 复制 RDB 文件
        $redisBackupFile = Join-Path $backupPath "redis_dump.rdb"
        docker cp gamelink-redis:/data/dump.rdb $redisBackupFile
        
        $fileSize = (Get-Item $redisBackupFile).Length / 1KB
        Write-ColorOutput "✅ Redis 备份完成 (${fileSize:N2} KB)" "Green"
    }
    catch {
        Write-ColorOutput "❌ Redis 备份失败: $_" "Red"
    }

    # 备份配置文件
    Write-ColorOutput "`n📦 备份配置文件..." "Cyan"
    try {
        Copy-Item $envFile (Join-Path $backupPath "env_backup.txt")
        Copy-Item "backend/configs/config.production.yaml" (Join-Path $backupPath "config.yaml") -ErrorAction SilentlyContinue
        Write-ColorOutput "✅ 配置文件备份完成" "Green"
    }
    catch {
        Write-ColorOutput "⚠️  配置文件备份失败: $_" "Yellow"
    }

    # 创建备份信息文件
    $backupInfo = @{
        Timestamp = $timestamp
        Environment = $Environment
        PostgreSQLBackup = $pgBackupFile
        RedisBackup = $redisBackupFile
    } | ConvertTo-Json

    $backupInfo | Out-File (Join-Path $backupPath "backup_info.json")

    # 压缩备份
    Write-ColorOutput "`n🗜️  压缩备份文件..." "Cyan"
    $zipFile = "$backupPath.zip"
    Compress-Archive -Path $backupPath -DestinationPath $zipFile -Force
    
    $zipSize = (Get-Item $zipFile).Length / 1MB
    Write-ColorOutput "✅ 备份已压缩: $zipFile (${zipSize:N2} MB)" "Green"

    # 清理临时目录
    Remove-Item -Path $backupPath -Recurse -Force

    Write-ColorOutput "`n✅ 备份完成!" "Green"
    Write-ColorOutput "`n📍 备份文件: $zipFile" "Yellow"
    Write-ColorOutput "`n💡 恢复命令:" "Yellow"
    Write-ColorOutput "   .\scripts\docker-restore.ps1 -BackupFile `"$zipFile`" -Environment $Environment" "White"
}

Backup-GameLink
