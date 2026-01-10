# 批量增加积分API测试脚本 (PowerShell版本)

# 配置
$API_URL = "http://localhost:8080/api/v1/admin/users/batch/points"
# 请替换为你的实际token
$TOKEN = "YOUR_ADMIN_TOKEN_HERE"

$headers = @{
    "Authorization" = "Bearer $TOKEN"
    "Content-Type" = "application/json"
}

Write-Host "========================================" -ForegroundColor Cyan
Write-Host "批量增加积分API测试脚本" -ForegroundColor Cyan
Write-Host "========================================" -ForegroundColor Cyan
Write-Host ""

# 测试1: 指定用户列表模式
Write-Host "测试1: 为指定用户增加积分" -ForegroundColor Green
Write-Host "目标: 用户ID 1,2,3"
Write-Host "积分: 500分 (5元)"
$body1 = @{
    target = "users"
    userIds = @(1, 2, 3)
    cents = 500
    reason = "测试-指定用户奖励"
    type = "admin"
} | ConvertTo-Json

try {
    $response = Invoke-RestMethod -Uri $API_URL -Method Post -Headers $headers -Body $body1
    $response | ConvertTo-Json -Depth 10 | Write-Host
} catch {
    Write-Host "错误: $_" -ForegroundColor Red
}

Write-Host ""
Write-Host "========================================" -ForegroundColor Cyan
Write-Host ""

# 测试2: 按角色筛选模式 - 单个角色
Write-Host "测试2: 为所有player角色增加积分" -ForegroundColor Green
Write-Host "角色: player"
Write-Host "积分: 1000分 (10元)"
$body2 = @{
    target = "role"
    roles = @("player")
    cents = 1000
    reason = "测试-陪玩师月度奖励"
    type = "admin"
} | ConvertTo-Json

try {
    $response = Invoke-RestMethod -Uri $API_URL -Method Post -Headers $headers -Body $body2
    $response | ConvertTo-Json -Depth 10 | Write-Host
} catch {
    Write-Host "错误: $_" -ForegroundColor Red
}

Write-Host ""
Write-Host "========================================" -ForegroundColor Cyan
Write-Host ""

# 测试3: 按角色筛选模式 - 多个角色
Write-Host "测试3: 为user和player角色增加积分" -ForegroundColor Green
Write-Host "角色: user, player"
Write-Host "积分: 200分 (2元)"
$body3 = @{
    target = "role"
    roles = @("user", "player")
    cents = 200
    reason = "测试-系统升级补偿"
    type = "compensation"
} | ConvertTo-Json

try {
    $response = Invoke-RestMethod -Uri $API_URL -Method Post -Headers $headers -Body $body3
    $response | ConvertTo-Json -Depth 10 | Write-Host
} catch {
    Write-Host "错误: $_" -ForegroundColor Red
}

Write-Host ""
Write-Host "========================================" -ForegroundColor Cyan
Write-Host ""

# 测试4: 全体用户模式
Write-Host "测试4: 为全体用户增加积分" -ForegroundColor Green
Write-Host "积分: 100分 (1元)"
$body4 = @{
    target = "all"
    cents = 100
    reason = "测试-平台活动奖励"
    type = "activity"
} | ConvertTo-Json

try {
    $response = Invoke-RestMethod -Uri $API_URL -Method Post -Headers $headers -Body $body4
    $response | ConvertTo-Json -Depth 10 | Write-Host
} catch {
    Write-Host "错误: $_" -ForegroundColor Red
}

Write-Host ""
Write-Host "========================================" -ForegroundColor Cyan
Write-Host ""

# 测试5: 错误场景 - target=users但未提供userIds
Write-Host "测试5: 错误场景 - 缺少userIds" -ForegroundColor Yellow
$body5 = @{
    target = "users"
    cents = 100
    reason = "测试-错误场景"
    type = "admin"
} | ConvertTo-Json

try {
    $response = Invoke-RestMethod -Uri $API_URL -Method Post -Headers $headers -Body $body5
    $response | ConvertTo-Json -Depth 10 | Write-Host
} catch {
    Write-Host "预期错误: $_" -ForegroundColor Yellow
}

Write-Host ""
Write-Host "========================================" -ForegroundColor Cyan
Write-Host ""

# 测试6: 错误场景 - target=role但未提供roles
Write-Host "测试6: 错误场景 - 缺少roles" -ForegroundColor Yellow
$body6 = @{
    target = "role"
    cents = 100
    reason = "测试-错误场景"
    type = "admin"
} | ConvertTo-Json

try {
    $response = Invoke-RestMethod -Uri $API_URL -Method Post -Headers $headers -Body $body6
    $response | ConvertTo-Json -Depth 10 | Write-Host
} catch {
    Write-Host "预期错误: $_" -ForegroundColor Yellow
}

Write-Host ""
Write-Host "========================================" -ForegroundColor Cyan
Write-Host ""

Write-Host "测试完成！" -ForegroundColor Green
Write-Host ""
Write-Host "注意事项：" -ForegroundColor Yellow
Write-Host "1. 请确保已替换脚本中的TOKEN为有效的管理员token"
Write-Host "2. 请确保数据库中存在测试用户"
Write-Host "3. 检查用户钱包余额确认积分是否正确增加"
Write-Host "4. 查看操作日志确认记录是否正确"
Write-Host ""
Write-Host "查询用户积分的SQL示例："
Write-Host "SELECT u.id, u.name, u.role, w.balance_cents FROM users u LEFT JOIN wallets w ON u.id = w.user_id;" -ForegroundColor Cyan
