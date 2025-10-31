# Test New Super Admin Credentials

$baseUrl = "http://localhost:8080/api/v1"

Write-Host "`n==========================================" -ForegroundColor Cyan
Write-Host "Testing New Super Admin Account" -ForegroundColor Cyan
Write-Host "==========================================" -ForegroundColor Cyan

# New credentials
$email = "superAdmin@GameLink.com"
$password = "admin123"

Write-Host "`nCredentials:" -ForegroundColor Yellow
Write-Host "  Email: $email" -ForegroundColor Gray
Write-Host "  Password: $password`n" -ForegroundColor Gray

# Check if service is running
Write-Host "[1/4] Checking service status..." -ForegroundColor Yellow
try {
    $health = Invoke-WebRequest -Uri "$baseUrl/healthz" -Method GET -TimeoutSec 5 -UseBasicParsing
    Write-Host "Service is running" -ForegroundColor Green
} catch {
    Write-Host "ERROR: Service is not running!" -ForegroundColor Red
    Write-Host "Please start the service first." -ForegroundColor Yellow
    exit 1
}

# Test login
Write-Host "`n[2/4] Testing login..." -ForegroundColor Yellow

$loginBody = @{
    email = $email
    password = $password
} | ConvertTo-Json -Compress

try {
    $loginResponse = Invoke-RestMethod -Uri "$baseUrl/auth/login" `
        -Method POST `
        -Body $loginBody `
        -ContentType "application/json" `
        -ErrorAction Stop
    
    if ($loginResponse.success) {
        Write-Host "Login successful!" -ForegroundColor Green
        Write-Host "  User ID: $($loginResponse.data.user.id)" -ForegroundColor Cyan
        Write-Host "  Name: $($loginResponse.data.user.name)" -ForegroundColor Cyan
        Write-Host "  Role: $($loginResponse.data.user.role)" -ForegroundColor Cyan
        $token = $loginResponse.data.accessToken
    } else {
        Write-Host "Login failed: $($loginResponse.message)" -ForegroundColor Red
        exit 1
    }
} catch {
    Write-Host "Login request failed!" -ForegroundColor Red
    if ($_.Exception.Response) {
        $statusCode = $_.Exception.Response.StatusCode.value__
        Write-Host "Status code: $statusCode" -ForegroundColor Red
    }
    Write-Host "Error: $_" -ForegroundColor Red
    exit 1
}

# Test dashboard access
Write-Host "`n[3/4] Testing dashboard access..." -ForegroundColor Yellow

$headers = @{
    "Authorization" = "Bearer $token"
}

try {
    $dashboard = Invoke-RestMethod -Uri "$baseUrl/admin/stats/dashboard" `
        -Method GET `
        -Headers $headers `
        -ErrorAction Stop
    
    if ($dashboard.success) {
        Write-Host "Dashboard access successful!" -ForegroundColor Green
        Write-Host "`nDashboard Stats:" -ForegroundColor Cyan
        Write-Host "  Total Users: $($dashboard.data.totalUsers)" -ForegroundColor Gray
        Write-Host "  Total Players: $($dashboard.data.totalPlayers)" -ForegroundColor Gray
        Write-Host "  Total Orders: $($dashboard.data.totalOrders)" -ForegroundColor Gray
        Write-Host "  Total Revenue: $($dashboard.data.totalRevenue)" -ForegroundColor Gray
    } else {
        Write-Host "Dashboard access failed!" -ForegroundColor Red
        Write-Host "Message: $($dashboard.message)" -ForegroundColor Red
        Write-Host "Code: $($dashboard.code)" -ForegroundColor Red
        exit 1
    }
} catch {
    Write-Host "Dashboard request failed!" -ForegroundColor Red
    if ($_.Exception.Response) {
        $statusCode = $_.Exception.Response.StatusCode.value__
        Write-Host "Status code: $statusCode" -ForegroundColor Red
    }
    exit 1
}

# Test other endpoints
Write-Host "`n[4/4] Testing other admin endpoints..." -ForegroundColor Yellow

$testEndpoints = @(
    @{url="$baseUrl/admin/users?page=1&pageSize=10"; name="Users List"}
    @{url="$baseUrl/admin/games?page=1&pageSize=10"; name="Games List"}
    @{url="$baseUrl/admin/stats/overview"; name="Stats Overview"}
)

$passCount = 0
foreach ($endpoint in $testEndpoints) {
    try {
        $result = Invoke-RestMethod -Uri $endpoint.url -Method GET -Headers $headers -ErrorAction Stop
        if ($result.success) {
            Write-Host "  $($endpoint.name): PASS" -ForegroundColor Green
            $passCount++
        } else {
            Write-Host "  $($endpoint.name): FAIL ($($result.message))" -ForegroundColor Red
        }
    } catch {
        Write-Host "  $($endpoint.name): ERROR" -ForegroundColor Red
    }
}

# Summary
Write-Host "`n==========================================" -ForegroundColor Cyan
if ($passCount -eq $testEndpoints.Count) {
    Write-Host "ALL TESTS PASSED!" -ForegroundColor Green
} else {
    Write-Host "SOME TESTS FAILED ($passCount/$($testEndpoints.Count) passed)" -ForegroundColor Yellow
}
Write-Host "==========================================" -ForegroundColor Cyan

Write-Host "`nYour new admin account is working correctly!" -ForegroundColor Green
Write-Host "`nYou can now use these credentials in your frontend:" -ForegroundColor Yellow
Write-Host "  Email: $email" -ForegroundColor Cyan
Write-Host "  Password: $password`n" -ForegroundColor Cyan


