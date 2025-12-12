<#
.SYNOPSIS
    GameLink Docker 清理工具
.DESCRIPTION
    清理 Docker 容器、镜像、数据卷和网络
.PARAMETER Level
    清理级别: soft (仅停止容器), medium (删除容器和镜像), hard (删除所有包括数据)
.PARAMETER Environment
    环境类型: dev, prod-local, prod 或 all
.EXAMPLE
    .\scripts\docker-clean.ps1 -Level soft -Environment prod-local
    .\scripts\docker-clean.ps1 -Level hard -Environment all
#>

param(
    [Parameter(Mandatory=$false)]
    [ValidateSet("soft", "medium", "hard")]
    [string]$Level = "soft",
    
    [Parameter(Mandatory=$false)]
    [ValidateSet("dev", "prod-local", "prod", "all")]
    [string]$Environment = "prod-local"
)

$ErrorActionPreference = "Continue"

function Write-ColorOutput {
    param([string]$Message, [string]$Color = "White")
    Write-Host $Message -ForegroundColor $Color
}

function Clean-Environment {
    param([string]$Env)

    $composeFile = switch ($Env) {
        "dev" { "docker-compose.yml" }
        "prod-local" { "docker-compose.prod.local.yml" }
        "prod" { "docker-compose.prod.yml" }
    }

    $envFile = switch ($Env) {
        "dev" { ".env.development" }
        "prod-local" { ".env.production.local" }
        "prod" { ".env" }
    }

    Write-ColorOutput "`n🧹 清理环境: $Env" "Cyan"

    if (-not (Test-Path $composeFile)) {
        Write-ColorOutput "⚠️  配置文件不存在: $composeFile" "Yellow"
        return
    }

    $baseCmd = "docker-compose"
    if (Test-Path $envFile) {
        $baseCmd += " --env-file $envFile"
    }
    $baseCmd += " -f $composeFile"

    switch ($Level) {
        "soft" {
            Write-ColorOutput "  停止容器..." "Cyan"
            Invoke-Expression "$baseCmd stop"
        }
        "medium" {
            Write-ColorOutput "  删除容器和镜像..." "Cyan"
            Invoke-Expression "$baseCmd down --rmi local"
        }
        "hard" {
            Write-ColorOutput "  删除所有（包括数据卷）..." "Cyan"
            Invoke-Expression "$baseCmd down -v --rmi local"
        }
    }

    Write-ColorOutput "  ✅ 完成" "Green"
}

Write-ColorOutput "`n🧹 GameLink Docker 清理工具" "Yellow"
Write-ColorOutput "==========================`n" "Yellow"

Write-ColorOutput "清理级别: $Level" "Cyan"
Write-ColorOutput "目标环境: $Environment`n" "Cyan"

# 显示清理级别说明
switch ($Level) {
    "soft" {
        Write-ColorOutput "📝 Soft 清理: 仅停止容器，保留所有数据" "White"
    }
    "medium" {
        Write-ColorOutput "📝 Medium 清理: 删除容器和镜像，保留数据卷" "White"
    }
    "hard" {
        Write-ColorOutput "📝 Hard 清理: 删除所有容器、镜像和数据卷" "Red"
        Write-ColorOutput "⚠️  警告: 这将删除所有数据，无法恢复！" "Red"
    }
}

# 确认操作
if ($Level -eq "hard") {
    Write-Host ""
    $confirm = Read-Host "确认执行 HARD 清理? 输入 'DELETE' 确认"
    if ($confirm -ne "DELETE") {
        Write-ColorOutput "已取消清理" "Yellow"
        exit 0
    }
}
elseif ($Level -eq "medium") {
    Write-Host ""
    $confirm = Read-Host "确认清理? (yes/no)"
    if ($confirm -ne "yes") {
        Write-ColorOutput "已取消清理" "Yellow"
        exit 0
    }
}

# 执行清理
if ($Environment -eq "all") {
    Clean-Environment -Env "dev"
    Clean-Environment -Env "prod-local"
    Clean-Environment -Env "prod"
}
else {
    Clean-Environment -Env $Environment
}

# 额外清理（仅在 hard 模式）
if ($Level -eq "hard") {
    Write-ColorOutput "`n🧹 清理未使用的资源..." "Cyan"
    
    Write-ColorOutput "  清理未使用的网络..." "Cyan"
    docker network prune -f 2>&1 | Out-Null
    
    Write-ColorOutput "  清理未使用的数据卷..." "Cyan"
    docker volume prune -f 2>&1 | Out-Null
    
    Write-ColorOutput "  清理未使用的镜像..." "Cyan"
    docker image prune -a -f 2>&1 | Out-Null
    
    Write-ColorOutput "  ✅ 完成" "Green"
}

# 显示清理后的状态
Write-ColorOutput "`n📊 当前 Docker 状态:" "Yellow"
Write-ColorOutput "`n容器:" "Cyan"
docker ps -a --filter "name=gamelink" --format "table {{.Names}}\t{{.Status}}\t{{.Size}}"

Write-ColorOutput "`n镜像:" "Cyan"
docker images --filter "reference=*gamelink*" --format "table {{.Repository}}\t{{.Tag}}\t{{.Size}}"

Write-ColorOutput "`n数据卷:" "Cyan"
docker volume ls --filter "name=gamelink" --format "table {{.Name}}\t{{.Driver}}"

Write-ColorOutput "`n✅ 清理完成!" "Green"
