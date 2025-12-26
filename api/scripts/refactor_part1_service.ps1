# Part 1: Service层文件重命名脚本 (PowerShell)
# 执行前请确保代码已提交

$ErrorActionPreference = "Stop"

Write-Host "🚀 Part 1: Service层文件重命名开始..." -ForegroundColor Green

$scriptPath = Split-Path -Parent $MyInvocation.MyCommand.Path
Set-Location "$scriptPath\.."

# 备份当前状态
Write-Host "📦 创建备份分支..." -ForegroundColor Yellow
git checkout -b refactor/part1-service-backup
git checkout -b refactor/part1-service

Write-Host "📝 重命名Service文件..." -ForegroundColor Cyan

# Auth
Set-Location internal\service\auth
if (Test-Path "auth_service.go") {
    git mv auth_service.go auth.go
}
if (Test-Path "auth_service_test.go") {
    git mv auth_service_test.go auth_test.go
}
Set-Location ..\..\..

# Order
Set-Location internal\service\order
if (Test-Path "order_service.go") {
    git mv order_service.go order.go
}
if (Test-Path "order_service_test.go") {
    git mv order_service_test.go order_test.go
}
Set-Location ..\..\..

# Player
Set-Location internal\service\player
if (Test-Path "player_service.go") {
    git mv player_service.go player.go
}
if (Test-Path "player_service_test.go") {
    git mv player_service_test.go player_test.go
}
Set-Location ..\..\..

# Payment
Set-Location internal\service\payment
if (Test-Path "payment_service.go") {
    git mv payment_service.go payment.go
}
if (Test-Path "payment_service_test.go") {
    git mv payment_service_test.go payment_test.go
}
Set-Location ..\..\..

# Review
Set-Location internal\service\review
if (Test-Path "review_service.go") {
    git mv review_service.go review.go
}
if (Test-Path "review_service_test.go") {
    git mv review_service_test.go review_test.go
}
Set-Location ..\..\..

# Earnings
Set-Location internal\service\earnings
if (Test-Path "earnings_service.go") {
    git mv earnings_service.go earnings.go
}
if (Test-Path "earnings_service_test.go") {
    git mv earnings_service_test.go earnings_test.go
}
Set-Location ..\..\..

# Gift
Set-Location internal\service\gift
if (Test-Path "gift_service.go") {
    git mv gift_service.go gift.go
}
if (Test-Path "gift_service_test.go") {
    git mv gift_service_test.go gift_test.go
}
Set-Location ..\..\..

# ServiceItem → Item
Write-Host "📦 重命名 serviceitem → item ..." -ForegroundColor Yellow
if (Test-Path "internal\service\serviceitem") {
    git mv internal\service\serviceitem internal\service\item
    Set-Location internal\service\item
    if (Test-Path "service_item.go") {
        git mv service_item.go item.go
    }
    if (Test-Path "service_item_test.go") {
        git mv service_item_test.go item_test.go
    }
    Set-Location ..\..\..
}

# Commission
Set-Location internal\service\commission
if (Test-Path "commission_service.go") {
    git mv commission_service.go commission.go
}
if (Test-Path "commission_service_test.go") {
    git mv commission_service_test.go commission_test.go
}
Set-Location ..\..\..

# Ranking
if (Test-Path "internal\service\ranking") {
    Set-Location internal\service\ranking
    if (Test-Path "ranking_service.go") {
        git mv ranking_service.go ranking.go
    }
    if (Test-Path "ranking_service_test.go") {
        git mv ranking_service_test.go ranking_test.go
    }
    Set-Location ..\..\..
}

Write-Host "🔄 更新cmd/main.go中的导入路径..." -ForegroundColor Yellow
# 更新 serviceitem → item
(Get-Content cmd\main.go) -replace 'serviceitemservice', 'itemservice' | Set-Content cmd\main.go
(Get-Content cmd\main.go) -replace 'service/serviceitem', 'service/item' | Set-Content cmd\main.go

Write-Host "✅ 编译测试..." -ForegroundColor Green
go build .\...

Write-Host "✅ 运行测试..." -ForegroundColor Green
go test .\internal\service\... -v

Write-Host "✅ Part 1 完成！" -ForegroundColor Green
Write-Host "📝 请检查修改，确认无误后提交：" -ForegroundColor Cyan
Write-Host "   git add ."
Write-Host "   git commit -m 'refactor(service): remove redundant _service suffix'"
Write-Host "   git push origin refactor/part1-service"


