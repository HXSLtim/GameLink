# 使用正确的超级管理员账户测试

$baseUrl = "http://localhost:8080/api/v1"

Write-Host "`n========================================" -ForegroundColor Cyan
Write-Host "Testing with Correct Super Admin" -ForegroundColor Cyan
Write-Host "========================================`n" -ForegroundColor Cyan

# 正确的超级管理员凭证
$correctEmail = "superAdmin@GameLink.com"
$correctPassword = "admin123"

Write-Host "Using credentials:" -ForegroundColor Yellow
Write-Host "  Email: $correctEmail" -ForegroundColor Gray
Write-Host "  Password: $correctPassword`n" -ForegroundColor Gray

# Step 1: Test login
Write-Host "[Step 1] Testing login..." -ForegroundColor Yellow

$loginJson = @"
{
    "email": "$correctEmail",
    "password": "$correctPassword"
}
"@

try {
    $response = Invoke-RestMethod -Uri "$baseUrl/auth/login" -Method POST -Body ([System.Text.Encoding]::UTF8.GetBytes($loginJson)) -ContentType "application/json; charset=utf-8"
    
    if ($response.success) {
        Write-Host "SUCCESS: Login successful!" -ForegroundColor Green
        Write-Host "  User ID: $($response.data.user.id)" -ForegroundColor Cyan
        Write-Host "  Name: $($response.data.user.name)" -ForegroundColor Cyan
        Write-Host "  Email: $($response.data.user.email)" -ForegroundColor Cyan
        Write-Host "  Role: $($response.data.user.role)" -ForegroundColor Cyan
        
        $token = $response.data.accessToken
        Write-Host "`n  Token: $($token.Substring(0, [Math]::Min(50, $token.Length)))...`n" -ForegroundColor Gray
    } else {
        Write-Host "FAILED: $($response.message)" -ForegroundColor Red
        exit 1
    }
} catch {
    Write-Host "ERROR: Login request failed" -ForegroundColor Red
    Write-Host $_.Exception.Message -ForegroundColor Red
    exit 1
}

# Step 2: Test dashboard access
Write-Host "[Step 2] Testing dashboard access..." -ForegroundColor Yellow

$headers = @{
    "Authorization" = "Bearer $token"
}

try {
    $dashboard = Invoke-RestMethod -Uri "$baseUrl/admin/stats/dashboard" -Method GET -Headers $headers
    
    if ($dashboard.success) {
        Write-Host "SUCCESS: Dashboard access granted!" -ForegroundColor Green
        Write-Host "`nDashboard Data:" -ForegroundColor Cyan
        Write-Host "  Total Users: $($dashboard.data.totalUsers)" -ForegroundColor Gray
        Write-Host "  Total Players: $($dashboard.data.totalPlayers)" -ForegroundColor Gray
        Write-Host "  Total Orders: $($dashboard.data.totalOrders)" -ForegroundColor Gray
        Write-Host "  Total Revenue: $($dashboard.data.totalRevenue)" -ForegroundColor Gray
        Write-Host "  Active Orders: $($dashboard.data.activeOrders)" -ForegroundColor Gray
    } else {
        Write-Host "FAILED: $($dashboard.message)" -ForegroundColor Red
        Write-Host "Code: $($dashboard.code)" -ForegroundColor Red
        exit 1
    }
} catch {
    Write-Host "ERROR: Dashboard request failed" -ForegroundColor Red
    Write-Host $_.Exception.Message -ForegroundColor Red
    exit 1
}

# Step 3: Test other admin endpoints
Write-Host "`n[Step 3] Testing other admin endpoints..." -ForegroundColor Yellow

$endpoints = @(
    @{method="GET"; path="/admin/users"; name="Users List"}
    @{method="GET"; path="/admin/games"; name="Games List"}
    @{method="GET"; path="/admin/orders"; name="Orders List"}
)

foreach ($endpoint in $endpoints) {
    try {
        $result = Invoke-RestMethod -Uri "$baseUrl$($endpoint.path)" -Method $endpoint.method -Headers $headers
        if ($result.success) {
            Write-Host "  $($endpoint.name): OK" -ForegroundColor Green
        } else {
            Write-Host "  $($endpoint.name): FAILED ($($result.message))" -ForegroundColor Red
        }
    } catch {
        Write-Host "  $($endpoint.name): ERROR" -ForegroundColor Red
    }
}

Write-Host "`n========================================" -ForegroundColor Cyan
Write-Host "All Tests Passed!" -ForegroundColor Green
Write-Host "========================================`n" -ForegroundColor Cyan

Write-Host "Your admin credentials are working correctly." -ForegroundColor Green
Write-Host "You can now use the frontend with these credentials:`n" -ForegroundColor Yellow
Write-Host "  Email: $correctEmail" -ForegroundColor Cyan
Write-Host "  Password: $correctPassword`n" -ForegroundColor Cyan

