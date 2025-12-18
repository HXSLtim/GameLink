# GameLink 生产环境部署脚本（标准版 - 不含加密）
# 用于快速部署，不启用加密中间件

param(
    [switch]$SkipBuild,
    [switch]$Skipadmin,
    [switch]$NoPull
)

$ErrorActionPreference = "Stop"

Write-Host "🚀 GameLink 生产环境部署（标准版）" -ForegroundColor Cyan
Write-Host "===================================" -ForegroundColor Cyan
Write-Host ""

# 1. 检查环境变量
Write-Host "📋 步骤 1/6: 检查环境变量..." -ForegroundColor Yellow
if (-not (Test-Path ".env")) {
    Write-Host "❌ 错误: 找不到 .env 文件" -ForegroundColor Red
    Write-Host "   请先创建 .env 文件，参考 .env.example" -ForegroundColor Yellow
    exit 1
}

# 检查必需的环境变量
$envContent = Get-Content .env -Raw
$requiredVars = @("POSTGRES_PASSWORD", "REDIS_PASSWORD", "JWT_SECRET_KEY", "SUPER_ADMIN_PASSWORD")
$missingVars = @()

foreach ($var in $requiredVars) {
    if ($envContent -notmatch "$var=(.+)") {
        $missingVars += $var
    }
}

if ($missingVars.Count -gt 0) {
    Write-Host "❌ 错误: .env 文件中缺少以下必需变量:" -ForegroundColor Red
    foreach ($var in $missingVars) {
        Write-Host "   - $var" -ForegroundColor Yellow
    }
    exit 1
}

Write-Host "✅ 环境变量检查通过" -ForegroundColor Green

# 2. 安装管理后台依赖
if (-not $Skipadmin) {
    Write-Host ""
    Write-Host "📦 步骤 2/6: 检查管理后台依赖..." -ForegroundColor Yellow
    Push-Location admin
    
    if (-not (Test-Path "node_modules")) {
        Write-Host "   安装管理后台依赖..." -ForegroundColor Gray
        npm install
        if ($LASTEXITCODE -ne 0) {
            Write-Host "❌ 管理后台依赖安装失败" -ForegroundColor Red
            Pop-Location
            exit 1
        }
    } else {
        Write-Host "   管理后台依赖已存在" -ForegroundColor Gray
    }
    
    Pop-Location
    Write-Host "✅ 管理后台依赖检查完成" -ForegroundColor Green
}

# 3. 构建管理后台
if (-not $Skipadmin -and -not $SkipBuild) {
    Write-Host ""
    Write-Host "🔨 步骤 3/6: 构建管理后台..." -ForegroundColor Yellow
    Push-Location admin
    npm run build
    if ($LASTEXITCODE -ne 0) {
        Write-Host "❌ 管理后台构建失败" -ForegroundColor Red
        Pop-Location
        exit 1
    }
    Pop-Location
    Write-Host "✅ 管理后台构建完成" -ForegroundColor Green
}

# 4. 构建 Docker 镜像
if (-not $SkipBuild) {
    Write-Host ""
    Write-Host "🐳 步骤 4/6: 构建 Docker 镜像..." -ForegroundColor Yellow
    
    if (-not $NoPull) {
        Write-Host "   拉取基础镜像..." -ForegroundColor Gray
        docker-compose -f docker-compose.prod.yml pull postgres redis
    }
    
    docker-compose -f docker-compose.prod.yml build
    if ($LASTEXITCODE -ne 0) {
        Write-Host "❌ Docker 镜像构建失败" -ForegroundColor Red
        exit 1
    }
    Write-Host "✅ Docker 镜像构建完成" -ForegroundColor Green
}

# 5. 停止旧服务
Write-Host ""
Write-Host "🛑 步骤 5/6: 停止旧服务..." -ForegroundColor Yellow
docker-compose -f docker-compose.prod.yml down
Write-Host "✅ 旧服务已停止" -ForegroundColor Green

# 6. 启动新服务
Write-Host ""
Write-Host "▶️  步骤 6/6: 启动新服务..." -ForegroundColor Yellow
docker-compose -f docker-compose.prod.yml up -d
if ($LASTEXITCODE -ne 0) {
    Write-Host "❌ 服务启动失败" -ForegroundColor Red
    exit 1
}

Write-Host ""
Write-Host "⏳ 等待服务启动..." -ForegroundColor Gray
Start-Sleep -Seconds 20

# 检查服务状态
Write-Host ""
Write-Host "📊 服务状态:" -ForegroundColor Cyan
docker ps --filter "name=gamelink" --format "table {{.Names}}\t{{.Status}}\t{{.Ports}}"

Write-Host ""
Write-Host "🔍 健康检查..." -ForegroundColor Yellow

# 检查后端 - 使用多种方式确保准确检测
$backendHealthy = $false
try {
    $response = Invoke-WebRequest -Uri "http://localhost:8080/api/v1/healthz" -TimeoutSec 10 -UseBasicParsing -ErrorAction Stop
    if ($response.StatusCode -eq 200) {
        $backendHealthy = $true
    }
} catch {
    try {
        $curlResult = curl.exe -s -o NUL -w "%{http_code}" "http://localhost:8080/api/v1/healthz" 2>$null
        if ($curlResult -eq "200") {
            $backendHealthy = $true
        }
    } catch { }
}

if ($backendHealthy) {
    Write-Host "   ✅ 后端服务健康" -ForegroundColor Green
} else {
    $backendContainer = docker ps --filter "name=gamelink-backend" --format "{{.Status}}" 2>$null
    if ($backendContainer -match "Up") {
        Write-Host "   ⚠️  后端容器运行中，但健康检查未通过（可能仍在初始化）" -ForegroundColor Yellow
        Write-Host "      容器状态: $backendContainer" -ForegroundColor Gray
    } else {
        Write-Host "   ❌ 后端服务未运行" -ForegroundColor Red
    }
}

# 检查管理后台
$adminHealthy = $false
try {
    $response = Invoke-WebRequest -Uri "http://localhost" -TimeoutSec 10 -UseBasicParsing -ErrorAction Stop
    if ($response.StatusCode -eq 200) {
        $adminHealthy = $true
    }
} catch {
    try {
        $curlResult = curl.exe -s -o NUL -w "%{http_code}" "http://localhost" 2>$null
        if ($curlResult -eq "200") {
            $adminHealthy = $true
        }
    } catch { }
}

if ($adminHealthy) {
    Write-Host "   ✅ 管理后台服务健康" -ForegroundColor Green
} else {
    $adminContainer = docker ps --filter "name=gamelink-admin" --format "{{.Status}}" 2>$null
    if ($adminContainer -match "Up") {
        Write-Host "   ⚠️  管理后台容器运行中，但健康检查未通过（可能仍在初始化）" -ForegroundColor Yellow
        Write-Host "      容器状态: $adminContainer" -ForegroundColor Gray
    } else {
        Write-Host "   ❌ 管理后台服务未运行" -ForegroundColor Red
    }
}

Write-Host ""
Write-Host "🎉 部署完成！" -ForegroundColor Green
Write-Host ""
Write-Host "📝 访问信息:" -ForegroundColor Cyan
Write-Host "   管理后台: http://localhost" -ForegroundColor Gray
Write-Host "   后端 API: http://localhost:8080/api/v1" -ForegroundColor Gray
Write-Host "   健康检查: http://localhost:8080/api/v1/healthz" -ForegroundColor Gray
Write-Host ""
Write-Host "🔑 管理员账号:" -ForegroundColor Cyan
$adminEmail = if ($envContent -match "SUPER_ADMIN_EMAIL=(.+)") { $matches[1].Trim() } else { "admin@gamelink.com" }
Write-Host "   邮箱: $adminEmail" -ForegroundColor Gray
Write-Host "   密码: 查看 .env 文件中的 SUPER_ADMIN_PASSWORD" -ForegroundColor Gray
Write-Host ""
Write-Host "🔐 加密状态:" -ForegroundColor Cyan
Write-Host "   后端: 未启用（标准部署）" -ForegroundColor Yellow
Write-Host "   管理后台: 未启用（标准部署）" -ForegroundColor Yellow
Write-Host ""
Write-Host "💡 提示:" -ForegroundColor Cyan
Write-Host "   - 如需启用加密，请使用: .\scripts\deploy-production-encrypted.ps1" -ForegroundColor Gray
Write-Host "   - 查看日志: docker logs gamelink-backend --tail=50" -ForegroundColor Gray
Write-Host "   - 查看所有服务: docker-compose -f docker-compose.prod.yml ps" -ForegroundColor Gray
Write-Host ""
Write-Host "📚 相关文档:" -ForegroundColor Cyan
Write-Host "   - 部署指南: DOCKER_DEPLOYMENT.md" -ForegroundColor Gray
Write-Host "   - 快速参考: DOCKER_QUICK_REFERENCE.md" -ForegroundColor Gray
Write-Host "   - 生产部署: PRODUCTION_DEPLOYMENT_GUIDE.md" -ForegroundColor Gray

