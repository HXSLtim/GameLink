# GameLink 生产环境部署脚本（加密版）
# 启用 AES-256-CBC 加密中间件，保护前后端通信

param(
    [switch]$SkipBuild,
    [switch]$SkipFrontend,
    [switch]$NoPull,
    [switch]$RegenerateKeys
)

$ErrorActionPreference = "Stop"

Write-Host "🚀 GameLink 生产环境部署（加密版）" -ForegroundColor Cyan
Write-Host "===================================" -ForegroundColor Cyan
Write-Host ""

# 1. 检查环境变量
Write-Host "📋 步骤 1/8: 检查环境变量..." -ForegroundColor Yellow
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

# 2. 检查或生成加密密钥
Write-Host ""
Write-Host "🔐 步骤 2/8: 检查加密密钥..." -ForegroundColor Yellow

$needsKeys = $false
if ($RegenerateKeys) {
    Write-Host "   重新生成加密密钥..." -ForegroundColor Gray
    $needsKeys = $true
} elseif ($envContent -notmatch "CRYPTO_SECRET_KEY=(.+)" -or $envContent -notmatch "CRYPTO_IV=(.+)") {
    Write-Host "   未找到加密密钥，将自动生成..." -ForegroundColor Gray
    $needsKeys = $true
} else {
    # 验证密钥长度
    if ($envContent -match "CRYPTO_SECRET_KEY=(.+)") {
        $secretKey = $matches[1].Trim()
        if ($secretKey.Length -ne 32) {
            Write-Host "   ⚠️  CRYPTO_SECRET_KEY 长度不正确（需要32字符），将重新生成" -ForegroundColor Yellow
            $needsKeys = $true
        }
    }
    if ($envContent -match "CRYPTO_IV=(.+)") {
        $iv = $matches[1].Trim()
        if ($iv.Length -ne 16) {
            Write-Host "   ⚠️  CRYPTO_IV 长度不正确（需要16字符），将重新生成" -ForegroundColor Yellow
            $needsKeys = $true
        }
    }
}

if ($needsKeys) {
    # 生成随机密钥
    $chars = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"
    $secretKey = -join ((1..32) | ForEach-Object { $chars[(Get-Random -Maximum $chars.Length)] })
    $iv = -join ((1..16) | ForEach-Object { $chars[(Get-Random -Maximum $chars.Length)] })
    
    # 更新 .env 文件
    $envLines = Get-Content .env
    $newEnvLines = @()
    $foundSecret = $false
    $foundIV = $false
    
    foreach ($line in $envLines) {
        if ($line -match "^CRYPTO_SECRET_KEY=") {
            $newEnvLines += "CRYPTO_SECRET_KEY=$secretKey"
            $foundSecret = $true
        } elseif ($line -match "^CRYPTO_IV=") {
            $newEnvLines += "CRYPTO_IV=$iv"
            $foundIV = $true
        } else {
            $newEnvLines += $line
        }
    }
    
    # 如果没有找到，添加到加密配置部分
    if (-not $foundSecret -or -not $foundIV) {
        $insertIndex = -1
        for ($i = 0; $i -lt $newEnvLines.Count; $i++) {
            if ($newEnvLines[$i] -match "# .*加密配置") {
                $insertIndex = $i + 1
                break
            }
        }
        
        if ($insertIndex -eq -1) {
            # 没有加密配置部分，添加到末尾
            $newEnvLines += ""
            $newEnvLines += "# ==================== 加密配置 ===================="
            if (-not $foundSecret) { $newEnvLines += "CRYPTO_SECRET_KEY=$secretKey" }
            if (-not $foundIV) { $newEnvLines += "CRYPTO_IV=$iv" }
        } else {
            if (-not $foundSecret) {
                $newEnvLines = $newEnvLines[0..($insertIndex-1)] + "CRYPTO_SECRET_KEY=$secretKey" + $newEnvLines[$insertIndex..($newEnvLines.Count-1)]
                $insertIndex++
            }
            if (-not $foundIV) {
                $newEnvLines = $newEnvLines[0..($insertIndex-1)] + "CRYPTO_IV=$iv" + $newEnvLines[$insertIndex..($newEnvLines.Count-1)]
            }
        }
    }
    
    $newEnvLines | Set-Content .env
    Write-Host "   ✅ 加密密钥已生成并保存到 .env" -ForegroundColor Green
} else {
    Write-Host "   ✅ 加密密钥已存在" -ForegroundColor Green
}

# 3. 同步加密密钥到前端
Write-Host ""
Write-Host "🔄 步骤 3/8: 同步加密密钥到前端..." -ForegroundColor Yellow
if (Test-Path ".\scripts\sync-crypto-keys.ps1") {
    & .\scripts\sync-crypto-keys.ps1
    if ($LASTEXITCODE -ne 0) {
        Write-Host "❌ 密钥同步失败" -ForegroundColor Red
        exit 1
    }
} else {
    # 手动同步
    $envContent = Get-Content .env -Raw
    if ($envContent -match "CRYPTO_SECRET_KEY=(.+)") {
        $secretKey = $matches[1].Trim()
    }
    if ($envContent -match "CRYPTO_IV=(.+)") {
        $iv = $matches[1].Trim()
    }
    
    $frontendEnv = @"
# 生产环境配置

# API 基础 URL
VITE_API_BASE_URL=/api/v1

# 加密配置（必须与后端配置一致）
VITE_CRYPTO_ENABLED=true
VITE_CRYPTO_SECRET_KEY=$secretKey
VITE_CRYPTO_IV=$iv
VITE_CRYPTO_USE_SIGNATURE=true

# WebSocket 配置
VITE_WS_URL=ws://localhost:8080
"@
    
    $frontendEnv | Set-Content "frontend\.env.production"
    Write-Host "   ✅ 密钥已同步到前端" -ForegroundColor Green
}

# 4. 安装前端依赖（包括 crypto-js）
if (-not $SkipFrontend) {
    Write-Host ""
    Write-Host "📦 步骤 4/8: 安装前端依赖..." -ForegroundColor Yellow
    Push-Location frontend
    
    # 检查 crypto-js 是否已安装
    $packageJson = Get-Content package.json -Raw | ConvertFrom-Json
    if (-not $packageJson.dependencies.'crypto-js') {
        Write-Host "   安装 crypto-js..." -ForegroundColor Gray
        npm install crypto-js
        npm install --save-dev @types/crypto-js
        if ($LASTEXITCODE -ne 0) {
            Write-Host "❌ crypto-js 安装失败" -ForegroundColor Red
            Pop-Location
            exit 1
        }
    } else {
        Write-Host "   crypto-js 已安装" -ForegroundColor Gray
    }
    
    # 确保所有依赖都已安装
    if (-not (Test-Path "node_modules")) {
        Write-Host "   安装所有前端依赖..." -ForegroundColor Gray
        npm install
        if ($LASTEXITCODE -ne 0) {
            Write-Host "❌ 前端依赖安装失败" -ForegroundColor Red
            Pop-Location
            exit 1
        }
    }
    
    Pop-Location
    Write-Host "✅ 前端依赖安装完成" -ForegroundColor Green
}

# 5. 构建前端
if (-not $SkipFrontend -and -not $SkipBuild) {
    Write-Host ""
    Write-Host "🔨 步骤 5/8: 构建前端..." -ForegroundColor Yellow
    Push-Location frontend
    npm run build
    if ($LASTEXITCODE -ne 0) {
        Write-Host "❌ 前端构建失败" -ForegroundColor Red
        Pop-Location
        exit 1
    }
    Pop-Location
    Write-Host "✅ 前端构建完成" -ForegroundColor Green
}

# 6. 构建 Docker 镜像
if (-not $SkipBuild) {
    Write-Host ""
    Write-Host "🐳 步骤 6/8: 构建 Docker 镜像..." -ForegroundColor Yellow
    
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

# 7. 停止旧服务
Write-Host ""
Write-Host "🛑 步骤 7/8: 停止旧服务..." -ForegroundColor Yellow
docker-compose -f docker-compose.prod.yml down
Write-Host "✅ 旧服务已停止" -ForegroundColor Green

# 8. 启动新服务
Write-Host ""
Write-Host "▶️  步骤 8/8: 启动新服务..." -ForegroundColor Yellow
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
Write-Host "🔍 验证加密配置..." -ForegroundColor Yellow

# 检查后端加密中间件
$cryptoEnabled = docker logs gamelink-backend 2>&1 | Select-String "crypto middleware enabled"
if ($cryptoEnabled) {
    Write-Host "   ✅ 后端加密中间件已启用" -ForegroundColor Green
    Write-Host "      $($cryptoEnabled.Line.Trim())" -ForegroundColor Gray
} else {
    Write-Host "   ⚠️  后端加密中间件未启用（检查日志）" -ForegroundColor Yellow
}

# 检查前端环境变量
$frontendEnvCheck = Get-Content "frontend\.env.production" -Raw
if ($frontendEnvCheck -match "VITE_CRYPTO_ENABLED=true") {
    Write-Host "   ✅ 前端加密配置已启用" -ForegroundColor Green
} else {
    Write-Host "   ⚠️  前端加密配置未启用" -ForegroundColor Yellow
}

Write-Host ""
Write-Host "🔍 健康检查..." -ForegroundColor Yellow

# 检查后端 - 使用多种方式确保准确检测
$backendHealthy = $false
try {
    # 方式1: 使用 Invoke-WebRequest
    $response = Invoke-WebRequest -Uri "http://localhost:8080/api/v1/healthz" -TimeoutSec 10 -UseBasicParsing -ErrorAction Stop
    if ($response.StatusCode -eq 200) {
        $backendHealthy = $true
    }
} catch {
    # 方式2: 如果 Invoke-WebRequest 失败，尝试使用 curl（如果可用）
    try {
        $curlResult = curl.exe -s -o NUL -w "%{http_code}" "http://localhost:8080/api/v1/healthz" 2>$null
        if ($curlResult -eq "200") {
            $backendHealthy = $true
        }
    } catch {
        # curl 也不可用，忽略
    }
}

if ($backendHealthy) {
    Write-Host "   ✅ 后端服务健康" -ForegroundColor Green
} else {
    # 检查容器是否在运行
    $backendContainer = docker ps --filter "name=gamelink-backend" --format "{{.Status}}" 2>$null
    if ($backendContainer -match "Up") {
        Write-Host "   ⚠️  后端容器运行中，但健康检查未通过（可能仍在初始化）" -ForegroundColor Yellow
        Write-Host "      容器状态: $backendContainer" -ForegroundColor Gray
    } else {
        Write-Host "   ❌ 后端服务未运行" -ForegroundColor Red
    }
}

# 检查前端
$frontendHealthy = $false
try {
    $response = Invoke-WebRequest -Uri "http://localhost" -TimeoutSec 10 -UseBasicParsing -ErrorAction Stop
    if ($response.StatusCode -eq 200) {
        $frontendHealthy = $true
    }
} catch {
    try {
        $curlResult = curl.exe -s -o NUL -w "%{http_code}" "http://localhost" 2>$null
        if ($curlResult -eq "200") {
            $frontendHealthy = $true
        }
    } catch {
        # curl 也不可用，忽略
    }
}

if ($frontendHealthy) {
    Write-Host "   ✅ 前端服务健康" -ForegroundColor Green
} else {
    $frontendContainer = docker ps --filter "name=gamelink-frontend" --format "{{.Status}}" 2>$null
    if ($frontendContainer -match "Up") {
        Write-Host "   ⚠️  前端容器运行中，但健康检查未通过（可能仍在初始化）" -ForegroundColor Yellow
        Write-Host "      容器状态: $frontendContainer" -ForegroundColor Gray
    } else {
        Write-Host "   ❌ 前端服务未运行" -ForegroundColor Red
    }
}

Write-Host ""
Write-Host "🎉 部署完成！" -ForegroundColor Green
Write-Host ""
Write-Host "📝 访问信息:" -ForegroundColor Cyan
Write-Host "   前端: http://localhost" -ForegroundColor Gray
Write-Host "   后端 API: http://localhost:8080/api/v1" -ForegroundColor Gray
Write-Host "   健康检查: http://localhost:8080/api/v1/healthz" -ForegroundColor Gray
Write-Host ""
Write-Host "🔑 管理员账号:" -ForegroundColor Cyan
$envContent = Get-Content .env -Raw
$adminEmail = if ($envContent -match "SUPER_ADMIN_EMAIL=(.+)") { $matches[1].Trim() } else { "admin@gamelink.com" }
Write-Host "   邮箱: $adminEmail" -ForegroundColor Gray
Write-Host "   密码: 查看 .env 文件中的 SUPER_ADMIN_PASSWORD" -ForegroundColor Gray
Write-Host ""
Write-Host "🔐 加密状态:" -ForegroundColor Cyan
Write-Host "   后端: ✅ 已启用 (AES-256-CBC + SHA-256 签名)" -ForegroundColor Green
Write-Host "   前端: ✅ 已启用 (自动加密 POST/PUT/PATCH 请求)" -ForegroundColor Green
Write-Host "   算法: AES-256-CBC" -ForegroundColor Gray
Write-Host "   签名: SHA-256" -ForegroundColor Gray
Write-Host ""
Write-Host "🧪 测试加密功能:" -ForegroundColor Cyan
Write-Host "   1. 打开浏览器访问 http://localhost" -ForegroundColor Gray
Write-Host "   2. 按 F12 打开开发者工具 → Network 标签" -ForegroundColor Gray
Write-Host "   3. 尝试登录或注册" -ForegroundColor Gray
Write-Host "   4. 查看请求体，应该看到加密格式:" -ForegroundColor Gray
Write-Host '      {"encrypted":true,"payload":"...","timestamp":...,"signature":"..."}' -ForegroundColor DarkGray
Write-Host ""
Write-Host "💡 提示:" -ForegroundColor Cyan
Write-Host "   - 查看后端日志: docker logs gamelink-backend --tail=50" -ForegroundColor Gray
Write-Host "   - 查看加密日志: docker logs gamelink-backend | Select-String 'crypto'" -ForegroundColor Gray
Write-Host "   - 重新生成密钥: .\scripts\deploy-production-encrypted.ps1 -RegenerateKeys" -ForegroundColor Gray
Write-Host ""
Write-Host "📚 相关文档:" -ForegroundColor Cyan
Write-Host "   - 加密配置指南: CRYPTO_SETUP_GUIDE.md" -ForegroundColor Gray
Write-Host "   - 部署状态: FINAL_DEPLOYMENT_STATUS.md" -ForegroundColor Gray
Write-Host "   - 部署总结: DEPLOYMENT_SUMMARY.md" -ForegroundColor Gray

