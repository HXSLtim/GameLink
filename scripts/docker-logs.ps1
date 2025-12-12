<#
.SYNOPSIS
    GameLink Docker 日志查看工具
.DESCRIPTION
    方便地查看和过滤 Docker 容器日志
.PARAMETER Environment
    环境类型: dev, prod-local 或 prod
.PARAMETER Service
    服务名称: backend, frontend, postgres, redis 或 all
.PARAMETER Follow
    是否持续跟踪日志
.PARAMETER Lines
    显示最近的行数
.EXAMPLE
    .\scripts\docker-logs.ps1 -Environment prod-local -Service backend -Follow
    .\scripts\docker-logs.ps1 -Service postgres -Lines 100
#>

param(
    [Parameter(Mandatory=$false)]
    [ValidateSet("dev", "prod-local", "prod")]
    [string]$Environment = "prod-local",
    
    [Parameter(Mandatory=$false)]
    [ValidateSet("backend", "frontend", "postgres", "redis", "all")]
    [string]$Service = "all",
    
    [Parameter(Mandatory=$false)]
    [switch]$Follow,
    
    [Parameter(Mandatory=$false)]
    [int]$Lines = 50
)

function Write-ColorOutput {
    param([string]$Message, [string]$Color = "White")
    Write-Host $Message -ForegroundColor $Color
}

# 确定配置文件
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

Write-ColorOutput "`n📋 GameLink 日志查看器" "Yellow"
Write-ColorOutput "=====================`n" "Yellow"
Write-ColorOutput "环境: $Environment" "Cyan"
Write-ColorOutput "服务: $Service" "Cyan"
Write-ColorOutput "行数: $Lines`n" "Cyan"

# 构建命令
$cmd = "docker-compose"
if (Test-Path $envFile) {
    $cmd += " --env-file $envFile"
}
$cmd += " -f $composeFile logs --tail=$Lines"

if ($Follow) {
    $cmd += " -f"
    Write-ColorOutput "💡 按 Ctrl+C 停止跟踪日志`n" "Yellow"
}

if ($Service -ne "all") {
    $cmd += " $Service"
}

Write-ColorOutput "执行命令: $cmd`n" "Gray"
Write-ColorOutput "==================== 日志开始 ====================`n" "Cyan"

# 执行命令
Invoke-Expression $cmd
