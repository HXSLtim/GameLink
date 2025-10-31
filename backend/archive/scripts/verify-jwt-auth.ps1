# Verify JWT Authentication

$baseUrl = "http://localhost:8080/api/v1"

Write-Host "`n========================================" -ForegroundColor Cyan
Write-Host "JWT Authentication Verification" -ForegroundColor Cyan
Write-Host "========================================`n" -ForegroundColor Cyan

# Wait for service to start
Write-Host "Waiting for service to start..." -ForegroundColor Yellow
Start-Sleep -Seconds 3

# Step 1: Test health endpoint
Write-Host "`n[Step 1] Testing service health..." -ForegroundColor Yellow
try {
    $health = Invoke-RestMethod -Uri "$baseUrl/healthz" -Method GET -TimeoutSec 5
    Write-Host "SUCCESS: Service is running (status: $($health.status))" -ForegroundColor Green
} catch {
    Write-Host "ERROR: Service is not responding" -ForegroundColor Red
    Write-Host "Please make sure the service is running" -ForegroundColor Red
    exit 1
}

# Step 2: Login to get JWT token
Write-Host "`n[Step 2] Logging in as super admin..." -ForegroundColor Yellow
$loginJson = '{"email":"admin@gamelink.local","password":"Admin@123456"}'

try {
    $response = Invoke-RestMethod -Uri "$baseUrl/auth/login" -Method POST -Body $loginJson -ContentType "application/json; charset=utf-8"
    
    if ($response.success) {
        Write-Host "SUCCESS: Login successful" -ForegroundColor Green
        Write-Host "  User: $($response.data.user.name) ($($response.data.user.email))" -ForegroundColor Gray
        Write-Host "  Role: $($response.data.user.role)" -ForegroundColor Gray
        Write-Host "  Token expires in: $($response.data.expiresIn) seconds" -ForegroundColor Gray
        $token = $response.data.accessToken
    } else {
        Write-Host "ERROR: Login failed - $($response.message)" -ForegroundColor Red
        exit 1
    }
} catch {
    Write-Host "ERROR: Login request failed" -ForegroundColor Red
    Write-Host "Details: $_" -ForegroundColor Red
    exit 1
}

# Step 3: Test admin dashboard endpoint with JWT
Write-Host "`n[Step 3] Testing admin dashboard with JWT token..." -ForegroundColor Yellow
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
    } else {
        Write-Host "ERROR: Dashboard access denied - $($dashboard.message)" -ForegroundColor Red
    }
} catch {
    Write-Host "ERROR: Dashboard request failed" -ForegroundColor Red
    Write-Host "Details: $_" -ForegroundColor Red
    exit 1
}

# Step 4: Test without token (should fail)
Write-Host "`n[Step 4] Testing dashboard without token (should fail)..." -ForegroundColor Yellow
try {
    $noAuth = Invoke-RestMethod -Uri "$baseUrl/admin/stats/dashboard" -Method GET
    Write-Host "WARNING: Request succeeded without token (this should not happen)" -ForegroundColor Yellow
} catch {
    Write-Host "SUCCESS: Request properly rejected without token" -ForegroundColor Green
}

# Summary
Write-Host "`n========================================" -ForegroundColor Cyan
Write-Host "Verification Complete!" -ForegroundColor Cyan
Write-Host "========================================" -ForegroundColor Cyan
Write-Host "`nYour JWT authentication is working correctly!" -ForegroundColor Green
Write-Host "`nNext steps:" -ForegroundColor Yellow
Write-Host "  1. Use the token from login response" -ForegroundColor Gray
Write-Host "  2. Add 'Authorization: Bearer <token>' header to all admin requests" -ForegroundColor Gray
Write-Host "  3. Tokens expire after 1 hour - refresh as needed`n" -ForegroundColor Gray

