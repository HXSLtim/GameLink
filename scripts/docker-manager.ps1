<#
.SYNOPSIS
    GameLink Docker 管理工具
.DESCRIPTION
    统一的 Docker 管理命令入口，提供所有常用操作
.PARAMETER Command
    要执行的命令
.EXAMPLE
    .\scripts\docker-manager.ps1 help
    .\scripts\docker-manager.ps1 start
    .\scripts\docker-manager.ps1 health
#>

param(
    [Parameter(Mandatory=$false, Position=0)]
    [string]$Command = "help",
    
    [Parameter(Mandatory=$false, ValueFromRemainingArguments=$true)]
    [string[]]$Args
)

$ErrorActionPreference = "Continue"

function Write-ColorOutput {
    param([string]$Message, [string]$Color = "White")
    Write-Host $Message -ForegroundColor $Color
}

function Show-Help {
    Write-ColorOutput "`n🎮 GameLink Docker 管理工具" "Yellow"
    Write-ColorOutput "==========================`n" "Yellow"
    
    Write-ColorOutput "使用方式: .\scripts\docker-manager.ps1 <command>`n" "Cyan"
    
    Write-ColorOutput "📦 环境启动:" "Yellow"
    Write-ColorOutput "  start              启动本地生产环境" "White"
    Write-ColorOutput "  start-dev          启动开发环境" "White"
    Write-ColorOutput "  start-clean        清理数据并启动" "White"
    Write-ColorOutput "  stop               停止本地生产环境" "White"
    Write-ColorOutput "  restart            重启本地生产环境" "White"
    
    Write-ColorOutput "`n🔍 监控和日志:" "Yellow"
    Write-ColorOutput "  health             健康检查" "White"
    Write-ColorOutput "  logs               查看所有日志" "White"
    Write-ColorOutput "  logs-backend       查看后端日志（持续）" "White"
    Write-ColorOutput "  logs-db            查看数据库日志" "White"
    Write-ColorOutput "  logs-redis         查看 Redis 日志" "White"
    Write-ColorOutput "  ps                 查看容器状态" "White"
    Write-ColorOutput "  stats              查看资源使用" "White"
    
    Write-ColorOutput "`n💾 数据管理:" "Yellow"
    Write-ColorOutput "  backup             备份数据" "White"
    Write-ColorOutput "  restore <file>     恢复数据" "White"
    Write-ColorOutput "  db-shell           进入数据库 Shell" "White"
    Write-ColorOutput "  redis-shell        进入 Redis Shell" "White"
    
    Write-ColorOutput "`n🧹 清理操作:" "Yellow"
    Write-ColorOutput "  clean              软清理（停止容器）" "White"
    Write-ColorOutput "  clean-medium       中等清理（删除容器和镜像）" "White"
    Write-ColorOutput "  clean-hard         完全清理（删除所有数据）" "White"
    Write-ColorOutput "  prune              清理未使用的资源" "White"
    
    Write-ColorOutput "`n🔧 构建和更新:" "Yellow"
    Write-ColorOutput "  build              构建镜像" "White"
    Write-ColorOutput "  rebuild            强制重新构建" "White"
    Write-ColorOutput "  update-backend     更新后端服务" "White"
    Write-ColorOutput "  update-admin    更新管理后台服务" "White"
    
    Write-ColorOutput "`n🚀 快捷操作:" "Yellow"
    Write-ColorOutput "  quick-start        快速启动" "White"
    Write-ColorOutput "  quick-check        快速检查（健康+状态+资源）" "White"
    Write-ColorOutput "  quick-restart      快速重启并检查" "White"
    
    Write-ColorOutput "`n📚 文档:" "Yellow"
    Write-ColorOutput "  docs               打开文档目录" "White"
    
    Write-ColorOutput "`n💡 示例:" "Cyan"
    Write-ColorOutput "  .\scripts\docker-manager.ps1 start" "Gray"
    Write-ColorOutput "  .\scripts\docker-manager.ps1 health" "Gray"
    Write-ColorOutput "  .\scripts\docker-manager.ps1 logs-backend" "Gray"
    Write-ColorOutput "  .\scripts\docker-manager.ps1 backup" "Gray"
    Write-ColorOutput "`n"
}

# 执行命令
switch ($Command.ToLower()) {
    # 环境启动
    "start" {
        Write-ColorOutput "🚀 启动本地生产环境..." "Green"
        & "$PSScriptRoot\docker-prod-local-start.ps1"
    }
    "start-dev" {
        Write-ColorOutput "🚀 启动开发环境..." "Green"
        & "$PSScriptRoot\docker-dev-start.ps1"
    }
    "start-clean" {
        Write-ColorOutput "🚀 清理并启动..." "Green"
        & "$PSScriptRoot\docker-prod-local-start.ps1" -Clean
    }
    "stop" {
        Write-ColorOutput "🛑 停止服务..." "Yellow"
        docker-compose --env-file .env.production.local -f docker-compose.prod.local.yml down
    }
    "restart" {
        Write-ColorOutput "🔄 重启服务..." "Yellow"
        docker-compose --env-file .env.production.local -f docker-compose.prod.local.yml restart
    }
    
    # 监控和日志
    "health" {
        & "$PSScriptRoot\docker-health-check.ps1" -Environment prod-local
    }
    "logs" {
        & "$PSScriptRoot\docker-logs.ps1" -Environment prod-local
    }
    "logs-backend" {
        & "$PSScriptRoot\docker-logs.ps1" -Environment prod-local -Service backend -Follow
    }
    "logs-db" {
        & "$PSScriptRoot\docker-logs.ps1" -Environment prod-local -Service postgres
    }
    "logs-redis" {
        & "$PSScriptRoot\docker-logs.ps1" -Environment prod-local -Service redis
    }
    "ps" {
        Write-ColorOutput "`n📊 容器状态:" "Cyan"
        docker ps --filter "name=gamelink" --format "table {{.Names}}\t{{.Status}}\t{{.Ports}}"
    }
    "stats" {
        Write-ColorOutput "`n📊 资源使用:" "Cyan"
        docker stats --no-stream --filter "name=gamelink"
    }
    
    # 数据管理
    "backup" {
        & "$PSScriptRoot\docker-backup.ps1" -Environment prod-local
    }
    "restore" {
        if ($Args.Count -eq 0) {
            Write-ColorOutput "❌ 请指定备份文件" "Red"
            Write-ColorOutput "示例: .\scripts\docker-manager.ps1 restore .\backups\20250113_120000.zip" "Yellow"
        }
        else {
            & "$PSScriptRoot\docker-restore.ps1" -BackupFile $Args[0] -Environment prod-local
        }
    }
    "db-shell" {
        Write-ColorOutput "🔌 进入 PostgreSQL..." "Cyan"
        docker exec -it gamelink-postgres psql -U gamelink -d gamelink
    }
    "redis-shell" {
        Write-ColorOutput "🔌 进入 Redis..." "Cyan"
        docker exec -it gamelink-redis redis-cli -a redis123
    }
    
    # 清理操作
    "clean" {
        & "$PSScriptRoot\docker-clean.ps1" -Level soft -Environment prod-local
    }
    "clean-medium" {
        & "$PSScriptRoot\docker-clean.ps1" -Level medium -Environment prod-local
    }
    "clean-hard" {
        & "$PSScriptRoot\docker-clean.ps1" -Level hard -Environment prod-local
    }
    "prune" {
        Write-ColorOutput "🧹 清理未使用的资源..." "Yellow"
        docker system prune -f
        docker volume prune -f
        docker network prune -f
    }
    
    # 构建和更新
    "build" {
        Write-ColorOutput "🔨 构建镜像..." "Green"
        docker-compose --env-file .env.production.local -f docker-compose.prod.local.yml build
    }
    "rebuild" {
        Write-ColorOutput "🔨 强制重新构建..." "Green"
        docker-compose --env-file .env.production.local -f docker-compose.prod.local.yml build --no-cache
    }
    "update-backend" {
        Write-ColorOutput "🔄 更新后端服务..." "Green"
        docker-compose --env-file .env.production.local -f docker-compose.prod.local.yml build backend
        docker-compose --env-file .env.production.local -f docker-compose.prod.local.yml up -d --no-deps backend
    }
    "update-admin" {
        Write-ColorOutput "🔄 更新管理后台服务..." "Green"
        docker-compose --env-file .env.production.local -f docker-compose.prod.local.yml build admin
        docker-compose --env-file .env.production.local -f docker-compose.prod.local.yml up -d --no-deps admin
    }
    
    # 快捷操作
    "quick-start" {
        Write-ColorOutput "🚀 快速启动..." "Green"
        & "$PSScriptRoot\docker-prod-local-start.ps1"
    }
    "quick-check" {
        Write-ColorOutput "🔍 快速检查..." "Cyan"
        & "$PSScriptRoot\docker-health-check.ps1" -Environment prod-local
        Write-Host ""
        docker ps --filter "name=gamelink" --format "table {{.Names}}\t{{.Status}}\t{{.Ports}}"
        Write-Host ""
        docker stats --no-stream --filter "name=gamelink"
    }
    "quick-restart" {
        Write-ColorOutput "🔄 快速重启..." "Yellow"
        docker-compose --env-file .env.production.local -f docker-compose.prod.local.yml restart
        Start-Sleep -Seconds 5
        & "$PSScriptRoot\docker-health-check.ps1" -Environment prod-local
    }
    
    # 文档
    "docs" {
        Write-ColorOutput "`n📚 Docker 文档:" "Yellow"
        Write-ColorOutput "  DOCKER_DEPLOYMENT.md       - 完整部署指南" "White"
        Write-ColorOutput "  DOCKER_QUICK_REFERENCE.md  - 快速参考手册" "White"
        Write-ColorOutput "  README.docker.md           - 工具集说明" "White"
        Write-ColorOutput "  DOCKER_SETUP_COMPLETE.md   - 完成总结" "White"
        Write-Host ""
    }
    
    # 帮助
    { $_ -in "help", "-h", "--help", "?" } {
        Show-Help
    }
    
    default {
        Write-ColorOutput "❌ 未知命令: $Command" "Red"
        Write-ColorOutput "使用 'help' 查看所有可用命令`n" "Yellow"
        Show-Help
    }
}
