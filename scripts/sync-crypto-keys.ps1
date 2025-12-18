# 同步加密密钥到管理后台环境变量
# 从 .env 文件读取后端的加密密钥，并写入管理后台的 .env.production 文件

param(
    [string]$EnvFile = ".env",
    [string]$AdminEnvFile = "admin/.env.production"
)

Write-Host "🔐 同步加密密钥到管理后台..." -ForegroundColor Cyan

# 检查后端 .env 文件是否存在
if (-not (Test-Path $EnvFile)) {
    Write-Host "❌ 错误: 找不到后端环境变量文件 $EnvFile" -ForegroundColor Red
    exit 1
}

# 检查管理后台 .env.production 文件是否存在
if (-not (Test-Path $AdminEnvFile)) {
    Write-Host "❌ 错误: 找不到管理后台环境变量文件 $AdminEnvFile" -ForegroundColor Red
    exit 1
}

# 读取后端环境变量
$envContent = Get-Content $EnvFile -Raw
$cryptoSecretKey = ""
$cryptoIV = ""

# 提取 CRYPTO_SECRET_KEY
if ($envContent -match "CRYPTO_SECRET_KEY=(.+)") {
    $cryptoSecretKey = $matches[1].Trim()
}

# 提取 CRYPTO_IV
if ($envContent -match "CRYPTO_IV=(.+)") {
    $cryptoIV = $matches[1].Trim()
}

if ([string]::IsNullOrWhiteSpace($cryptoSecretKey) -or [string]::IsNullOrWhiteSpace($cryptoIV)) {
    Write-Host "❌ 错误: 无法从 $EnvFile 中提取加密密钥" -ForegroundColor Red
    Write-Host "   请确保 CRYPTO_SECRET_KEY 和 CRYPTO_IV 已配置" -ForegroundColor Yellow
    exit 1
}

Write-Host "✅ 从后端读取到加密密钥" -ForegroundColor Green
Write-Host "   CRYPTO_SECRET_KEY: $($cryptoSecretKey.Substring(0, [Math]::Min(10, $cryptoSecretKey.Length)))..." -ForegroundColor Gray
Write-Host "   CRYPTO_IV: $($cryptoIV.Substring(0, [Math]::Min(10, $cryptoIV.Length)))..." -ForegroundColor Gray

# 读取管理后台环境变量文件
$adminEnvContent = Get-Content $AdminEnvFile -Raw

# 更新管理后台环境变量
$adminEnvContent = $adminEnvContent -replace "VITE_CRYPTO_SECRET_KEY=.*", "VITE_CRYPTO_SECRET_KEY=$cryptoSecretKey"
$adminEnvContent = $adminEnvContent -replace "VITE_CRYPTO_IV=.*", "VITE_CRYPTO_IV=$cryptoIV"

# 写回文件
Set-Content -Path $AdminEnvFile -Value $adminEnvContent -NoNewline

Write-Host "✅ 加密密钥已同步到管理后台环境变量文件" -ForegroundColor Green
Write-Host "   文件: $AdminEnvFile" -ForegroundColor Gray

Write-Host ""
Write-Host "📝 下一步:" -ForegroundColor Cyan
Write-Host "   1. 安装 crypto-js 依赖: cd admin && npm install crypto-js" -ForegroundColor Yellow
Write-Host "   2. 重新构建管理后台: npm run build" -ForegroundColor Yellow
Write-Host "   3. 重启 Docker 服务: docker-compose -f docker-compose.prod.yml restart admin" -ForegroundColor Yellow

exit 0
