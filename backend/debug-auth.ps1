# JWT 认证调试脚本

Write-Host "========================================" -ForegroundColor Cyan
Write-Host "GameLink JWT 认证调试工具" -ForegroundColor Cyan
Write-Host "========================================`n" -ForegroundColor Cyan

$baseUrl = "http://localhost:8080/api/v1"

# 1. 测试健康检查
Write-Host "[1] 测试服务健康状态..." -ForegroundColor Yellow
try {
    $health = Invoke-RestMethod -Uri "$baseUrl/health" -Method GET -ErrorAction Stop
    Write-Host "✓ 服务运行正常`n" -ForegroundColor Green
} catch {
    Write-Host "✗ 服务未运行或无法访问`n" -ForegroundColor Red
    Write-Host "错误: $($_.Exception.Message)" -ForegroundColor Red
    exit 1
}

# 2. 测试超级管理员登录
Write-Host "[2] 测试超级管理员登录..." -ForegroundColor Yellow
$loginBody = @{
    email = "admin@gamelink.local"
    password = "Admin@123456"
} | ConvertTo-Json

try {
    $loginResponse = Invoke-RestMethod -Uri "$baseUrl/auth/login" -Method POST -Body $loginBody -ContentType "application/json" -ErrorAction Stop
    
    if ($loginResponse.success) {
        Write-Host "✓ 登录成功" -ForegroundColor Green
        Write-Host "  用户ID: $($loginResponse.data.user.id)" -ForegroundColor Gray
        Write-Host "  用户名: $($loginResponse.data.user.name)" -ForegroundColor Gray
        Write-Host "  角色: $($loginResponse.data.user.role)" -ForegroundColor Gray
        Write-Host "  Token (前20字符): $($loginResponse.data.accessToken.Substring(0, 20))...`n" -ForegroundColor Gray
        
        $token = $loginResponse.data.accessToken
    } else {
        Write-Host "✗ 登录失败: $($loginResponse.message)`n" -ForegroundColor Red
        exit 1
    }
} catch {
    Write-Host "✗ 登录请求失败`n" -ForegroundColor Red
    Write-Host "错误: $($_.Exception.Message)" -ForegroundColor Red
    exit 1
}

# 3. 测试 Token 访问受保护接口
Write-Host "[3] 测试访问管理员仪表盘..." -ForegroundColor Yellow
$headers = @{
    "Authorization" = "Bearer $token"
}

try {
    $dashboardResponse = Invoke-RestMethod -Uri "$baseUrl/admin/stats/dashboard" -Method GET -Headers $headers -ErrorAction Stop
    
    if ($dashboardResponse.success) {
        Write-Host "✓ 成功访问仪表盘接口" -ForegroundColor Green
        Write-Host "  返回数据:`n$($dashboardResponse.data | ConvertTo-Json -Depth 2)`n" -ForegroundColor Gray
    } else {
        Write-Host "✗ 访问失败: $($dashboardResponse.message)`n" -ForegroundColor Red
    }
} catch {
    $statusCode = $_.Exception.Response.StatusCode.value__
    Write-Host "✗ 请求失败 (HTTP $statusCode)" -ForegroundColor Red
    
    try {
        $errorStream = $_.Exception.Response.GetResponseStream()
        $reader = New-Object System.IO.StreamReader($errorStream)
        $errorBody = $reader.ReadToEnd() | ConvertFrom-Json
        Write-Host "  错误消息: $($errorBody.message)" -ForegroundColor Red
        Write-Host "  错误代码: $($errorBody.code)`n" -ForegroundColor Red
    } catch {
        Write-Host "  无法解析错误响应`n" -ForegroundColor Red
    }
}

# 4. 检查环境变量
Write-Host "[4] 检查认证配置..." -ForegroundColor Yellow
$appEnv = $env:APP_ENV
$adminAuthMode = $env:ADMIN_AUTH_MODE
$jwtSecret = $env:JWT_SECRET

Write-Host "  APP_ENV: $(if ($appEnv) { $appEnv } else { '(未设置，默认为开发模式)' })" -ForegroundColor Gray
Write-Host "  ADMIN_AUTH_MODE: $(if ($adminAuthMode) { $adminAuthMode } else { '(未设置，使用 AdminAuth)' })" -ForegroundColor Gray
Write-Host "  JWT_SECRET: $(if ($jwtSecret) { '已设置 (' + $jwtSecret.Length + ' 字符)' } else { '(未设置，使用默认值)' })" -ForegroundColor Gray

if (-not $adminAuthMode -or $adminAuthMode -notin @("jwt", "JWT")) {
    Write-Host "`n⚠️  警告: 开发环境未启用 JWT 认证模式" -ForegroundColor Yellow
    Write-Host "  建议设置环境变量: `$env:ADMIN_AUTH_MODE='jwt'" -ForegroundColor Yellow
}

# 5. 测试你提供的 Token
Write-Host "`n[5] 测试你的 Token..." -ForegroundColor Yellow
$yourToken = "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJ1c2VyX2lkIjo3LCJyb2xlIjoiYWRtaW4iLCJpc3MiOiJnYW1lbGluayIsImV4cCI6MTc2MTk3MjUyMiwibmJmIjoxNzYxODg2MTIyLCJpYXQiOjE3NjE4ODYxMjJ9.KxYF5xV8SVg7BFlPH4JxX8Xzsbsodg5YTzG4qycsOp8"

# 解码 JWT payload (Base64)
$payloadPart = $yourToken.Split('.')[1]
# 补齐 Base64 padding
$padding = 4 - ($payloadPart.Length % 4)
if ($padding -ne 4) {
    $payloadPart += "=" * $padding
}
$payloadBytes = [System.Convert]::FromBase64String($payloadPart)
$payloadJson = [System.Text.Encoding]::UTF8.GetString($payloadBytes) | ConvertFrom-Json

Write-Host "  Token Payload:" -ForegroundColor Gray
Write-Host "    用户ID: $($payloadJson.user_id)" -ForegroundColor Gray
Write-Host "    角色: $($payloadJson.role)" -ForegroundColor Gray
Write-Host "    颁发时间: $([DateTimeOffset]::FromUnixTimeSeconds($payloadJson.iat).LocalDateTime)" -ForegroundColor Gray
Write-Host "    过期时间: $([DateTimeOffset]::FromUnixTimeSeconds($payloadJson.exp).LocalDateTime)" -ForegroundColor Gray

$now = [DateTimeOffset]::UtcNow.ToUnixTimeSeconds()
if ($payloadJson.exp -lt $now) {
    Write-Host "  ✗ Token 已过期！" -ForegroundColor Red
} else {
    Write-Host "  ✓ Token 未过期" -ForegroundColor Green
}

# 使用你的 Token 测试
$yourHeaders = @{
    "Authorization" = "Bearer $yourToken"
}

try {
    $testResponse = Invoke-RestMethod -Uri "$baseUrl/admin/stats/dashboard" -Method GET -Headers $yourHeaders -ErrorAction Stop
    Write-Host "  ✓ 你的 Token 可以正常使用`n" -ForegroundColor Green
} catch {
    $statusCode = $_.Exception.Response.StatusCode.value__
    Write-Host "  ✗ 你的 Token 无法使用 (HTTP $statusCode)" -ForegroundColor Red
    
    try {
        $errorStream = $_.Exception.Response.GetResponseStream()
        $reader = New-Object System.IO.StreamReader($errorStream)
        $errorBody = $reader.ReadToEnd() | ConvertFrom-Json
        Write-Host "    原因: $($errorBody.message)`n" -ForegroundColor Red
    } catch {
        Write-Host "    无法解析错误响应`n" -ForegroundColor Red
    }
}

Write-Host "========================================" -ForegroundColor Cyan
Write-Host "调试完成" -ForegroundColor Cyan
Write-Host "========================================" -ForegroundColor Cyan

