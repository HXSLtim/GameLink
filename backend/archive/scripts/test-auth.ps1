# JWT Authentication Test Script

$baseUrl = "http://localhost:8080/api/v1"

Write-Host "`n=== Testing Authentication ===" -ForegroundColor Cyan

# Test 1: Login as super admin
Write-Host "`n[1] Testing login..." -ForegroundColor Yellow
$loginBody = (@{
    email = "admin@gamelink.local"
    password = "Admin@123456"
} | ConvertTo-Json)

$loginResponse = Invoke-RestMethod -Uri "$baseUrl/auth/login" -Method POST -Body ([System.Text.Encoding]::UTF8.GetBytes($loginBody)) -ContentType "application/json; charset=utf-8"

if ($loginResponse.success) {
    Write-Host "Login successful!" -ForegroundColor Green
    Write-Host "User ID: $($loginResponse.data.user.id)" -ForegroundColor Gray
    Write-Host "Role: $($loginResponse.data.user.role)" -ForegroundColor Gray
    $token = $loginResponse.data.accessToken
} else {
    Write-Host "Login failed: $($loginResponse.message)" -ForegroundColor Red
    exit 1
}

# Test 2: Access dashboard with new token
Write-Host "`n[2] Testing dashboard access with new token..." -ForegroundColor Yellow
$headers = @{
    "Authorization" = "Bearer $token"
}

try {
    $dashboardResponse = Invoke-RestMethod -Uri "$baseUrl/admin/stats/dashboard" -Method GET -Headers $headers
    Write-Host "Dashboard access successful!" -ForegroundColor Green
    Write-Host "Data received:" -ForegroundColor Gray
    Write-Host ($dashboardResponse.data | ConvertTo-Json -Depth 2)
} catch {
    Write-Host "Dashboard access failed!" -ForegroundColor Red
    Write-Host "Error: $_" -ForegroundColor Red
}

# Test 3: Test your token
Write-Host "`n[3] Testing your token..." -ForegroundColor Yellow
$yourToken = "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJ1c2VyX2lkIjo3LCJyb2xlIjoiYWRtaW4iLCJpc3MiOiJnYW1lbGluayIsImV4cCI6MTc2MTk3MjUyMiwibmJmIjoxNzYxODg2MTIyLCJpYXQiOjE3NjE4ODYxMjJ9.KxYF5xV8SVg7BFlPH4JxX8Xzsbsodg5YTzG4qycsOp8"

$yourHeaders = @{
    "Authorization" = "Bearer $yourToken"
}

try {
    $yourResponse = Invoke-RestMethod -Uri "$baseUrl/admin/stats/dashboard" -Method GET -Headers $yourHeaders
    Write-Host "Your token works!" -ForegroundColor Green
} catch {
    Write-Host "Your token failed!" -ForegroundColor Red
    Write-Host "Error: $_" -ForegroundColor Red
}

# Test 4: Check environment
Write-Host "`n[4] Environment check..." -ForegroundColor Yellow
Write-Host "APP_ENV: $(if ($env:APP_ENV) { $env:APP_ENV } else { '(not set)' })" -ForegroundColor Gray
Write-Host "ADMIN_AUTH_MODE: $(if ($env:ADMIN_AUTH_MODE) { $env:ADMIN_AUTH_MODE } else { '(not set - using AdminAuth)' })" -ForegroundColor Gray

if (-not $env:ADMIN_AUTH_MODE -or $env:ADMIN_AUTH_MODE -notin @("jwt", "JWT")) {
    Write-Host "`nWARNING: JWT auth is not enabled in dev mode!" -ForegroundColor Yellow
    Write-Host "Try: `$env:ADMIN_AUTH_MODE='jwt'" -ForegroundColor Yellow
}

Write-Host "`n=== Test Complete ===" -ForegroundColor Cyan

