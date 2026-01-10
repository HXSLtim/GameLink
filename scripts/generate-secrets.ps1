# GameLink 安全密钥生成工具
# 用于生成生产环境所需的各类安全密钥

Write-Host "========================================" -ForegroundColor Cyan
Write-Host "  GameLink 安全密钥生成工具" -ForegroundColor Cyan
Write-Host "========================================" -ForegroundColor Cyan
Write-Host ""

# 检查 OpenSSL 是否安装
$opensslExists = $null -ne (Get-Command "openssl" -ErrorAction SilentlyContinue)

if (-not $opensslExists) {
    Write-Host "错误: 未找到 OpenSSL 命令" -ForegroundColor Red
    Write-Host ""
    Write-Host "请安装 OpenSSL:" -ForegroundColor Yellow
    Write-Host "  Windows: 下载并安装 from https://slproweb.com/products/Win32OpenSSL.html" -ForegroundColor White
    Write-Host "  或使用 Git Bash / WSL 运行此脚本" -ForegroundColor White
    Write-Host ""
    exit 1
}

Write-Host "生成的密钥如下（请妥善保存，不要泄露）：" -ForegroundColor Green
Write-Host ""

# 生成 32 字节加密密钥 (AES-256-CBC)
$secretKey = openssl rand -base64 32
Write-Host "1. 加密密钥 (CRYPTO_SECRET_KEY) - 32字节:" -ForegroundColor Yellow
Write-Host "   $secretKey" -ForegroundColor White
Write-Host ""

# 生成 16 字节初始化向量
$iv = openssl rand -base64 16
Write-Host "2. 初始化向量 (CRYPTO_IV) - 16字节:" -ForegroundColor Yellow
Write-Host "   $iv" -ForegroundColor White
Write-Host ""

# 生成 32 字节 JWT 密钥
$jwtSecret = openssl rand -base64 32
Write-Host "3. JWT 密钥 (JWT_SECRET_KEY) - 32字节:" -ForegroundColor Yellow
Write-Host "   $jwtSecret" -ForegroundColor White
Write-Host ""

# 生成 24 字节超级管理员密码
$adminPassword = openssl rand -base64 24
Write-Host "4. 超级管理员密码 (SUPER_ADMIN_PASSWORD) - 24字节:" -ForegroundColor Yellow
Write-Host "   $adminPassword" -ForegroundColor White
Write-Host ""

Write-Host "========================================" -ForegroundColor Cyan
Write-Host "使用方法：" -ForegroundColor Cyan
Write-Host "========================================" -ForegroundColor Cyan
Write-Host ""
Write-Host "Windows PowerShell 环境变量设置:" -ForegroundColor Yellow
Write-Host "`$env:CRYPTO_SECRET_KEY='$secretKey'" -ForegroundColor White
Write-Host "`$env:CRYPTO_IV='$iv'" -ForegroundColor White
Write-Host "`$env:JWT_SECRET_KEY='$jwtSecret'" -ForegroundColor White
Write-Host "`$env:SUPER_ADMIN_PASSWORD='$adminPassword'" -ForegroundColor White
Write-Host ""

Write-Host "Linux/Mac 环境变量设置:" -ForegroundColor Yellow
Write-Host "export CRYPTO_SECRET_KEY='$secretKey'" -ForegroundColor White
Write-Host "export CRYPTO_IV='$iv'" -ForegroundColor White
Write-Host "export JWT_SECRET_KEY='$jwtSecret'" -ForegroundColor White
Write-Host "export SUPER_ADMIN_PASSWORD='$adminPassword'" -ForegroundColor White
Write-Host ""

Write-Host "Docker Compose 环境变量 (.env 文件):" -ForegroundColor Yellow
Write-Host "CRYPTO_SECRET_KEY=$secretKey" -ForegroundColor White
Write-Host "CRYPTO_IV=$iv" -ForegroundColor White
Write-Host "JWT_SECRET_KEY=$jwtSecret" -ForegroundColor White
Write-Host "SUPER_ADMIN_PASSWORD=$adminPassword" -ForegroundColor White
Write-Host ""

Write-Host "========================================" -ForegroundColor Cyan
Write-Host "重要提示：" -ForegroundColor Red
Write-Host "========================================" -ForegroundColor Red
Write-Host "1. 请将这些密钥保存到安全的位置（密码管理器）" -ForegroundColor White
Write-Host "2. 不要将密钥提交到 Git 仓库" -ForegroundColor White
Write-Host "3. 生产环境每次部署都应使用不同的密钥" -ForegroundColor White
Write-Host "4. 密钥泄露后立即重新生成并更新" -ForegroundColor White
Write-Host ""

# 可选：导出到 .env 文件
$exportToFile = Read-Host "是否导出到 .env 文件? (y/N)"
if ($exportToFile -eq 'y' -or $exportToFile -eq 'Y') {
    $envFile = ".env"
    @"
# GameLink 生产环境密钥
# 警告: 请勿将此文件提交到版本控制系统

# 加密配置 (AES-256-CBC)
CRYPTO_ENABLED=true
CRYPTO_SECRET_KEY=$secretKey
CRYPTO_IV=$iv

# JWT 配置
JWT_SECRET_KEY=$jwtSecret

# 超级管理员配置
SUPER_ADMIN_EMAIL=admin@gamelink.com
SUPER_ADMIN_PASSWORD=$adminPassword
SUPER_ADMIN_NAME=Super Admin

# 数据库配置 (请根据实际情况修改)
POSTGRES_USER=gamelink
POSTGRES_PASSWORD=your_secure_db_password_here
POSTGRES_DB=gamelink

# Redis 配置
REDIS_PASSWORD=your_secure_redis_password_here
"@ | Out-File -FilePath $envFile -Encoding UTF8

    Write-Host ""
    Write-Host "已导出到 $envFile" -ForegroundColor Green
    Write-Host "请编辑文件并填写数据库和 Redis 密码" -ForegroundColor Yellow
    Write-Host "然后将 .env 添加到 .gitignore" -ForegroundColor Yellow
}
