# GameLink 改进项目结构快速搭建脚本 (PowerShell版本)
# 使用方法: .\scripts\setup-improvement-structure.ps1

Write-Host "🚀 GameLink 改进项目结构搭建开始..." -ForegroundColor Cyan
Write-Host ""

# 获取项目根目录
$ProjectRoot = Split-Path -Parent $PSScriptRoot
Set-Location $ProjectRoot

Write-Host "📂 项目根目录: $ProjectRoot" -ForegroundColor Blue
Write-Host ""

# 辅助函数: 创建文件
function New-FileIfNotExists {
    param([string]$Path)
    
    if (Test-Path $Path) {
        Write-Host "⚠  已存在: $Path" -ForegroundColor Yellow
    } else {
        New-Item -Path $Path -ItemType File -Force | Out-Null
        Write-Host "✓ 创建: $Path" -ForegroundColor Green
    }
}

# 辅助函数: 创建目录
function New-DirectoryIfNotExists {
    param([string]$Path)
    
    if (Test-Path $Path) {
        Write-Host "⚠  目录已存在: $Path" -ForegroundColor Yellow
    } else {
        New-Item -Path $Path -ItemType Directory -Force | Out-Null
        Write-Host "✓ 创建目录: $Path" -ForegroundColor Green
    }
}

# ============================================
# 第一部分: 后端数据模型文件
# ============================================

Write-Host "📊 第一步: 创建后端数据模型文件..." -ForegroundColor Blue

$Models = @(
    "dispute",
    "ticket",
    "notification",
    "chat",
    "favorite",
    "tag"
)

foreach ($model in $Models) {
    $file = "backend\internal\model\$model.go"
    New-FileIfNotExists -Path $file
}

Write-Host ""

# ============================================
# 第二部分: Repository 层
# ============================================

Write-Host "📚 第二步: 创建 Repository 层文件..." -ForegroundColor Blue

$Repos = @(
    "dispute",
    "ticket",
    "notification",
    "chat",
    "favorite",
    "tag"
)

foreach ($repo in $Repos) {
    $dir = "backend\internal\repository\$repo"
    New-DirectoryIfNotExists -Path $dir
    
    $files = @("repository.go", "repository_test.go")
    foreach ($file in $files) {
        $fullPath = "$dir\$file"
        New-FileIfNotExists -Path $fullPath
    }
}

Write-Host ""

# ============================================
# 第三部分: Service 层
# ============================================

Write-Host "💼 第三步: 创建 Service 层文件..." -ForegroundColor Blue

$Services = @(
    "dispute",
    "ticket",
    "notification",
    "chat",
    "favorite",
    "upload"
)

foreach ($service in $Services) {
    $dir = "backend\internal\service\$service"
    New-DirectoryIfNotExists -Path $dir
    
    $files = @("service.go", "service_test.go")
    foreach ($file in $files) {
        $fullPath = "$dir\$file"
        New-FileIfNotExists -Path $fullPath
    }
}

# 创建支付服务文件
$PaymentFiles = @(
    "backend\internal\service\payment\alipay.go",
    "backend\internal\service\payment\wechat.go"
)

foreach ($file in $PaymentFiles) {
    New-FileIfNotExists -Path $file
}

# 创建聊天Hub
New-FileIfNotExists -Path "backend\internal\service\chat\hub.go"

Write-Host ""

# ============================================
# 第四部分: Handler 层
# ============================================

Write-Host "🎯 第四步: 创建 Handler 层文件..." -ForegroundColor Blue

# User Handler
$UserHandlers = @(
    "dispute",
    "ticket",
    "notification",
    "favorite"
)

foreach ($handler in $UserHandlers) {
    $file = "backend\internal\handler\user\$handler.go"
    New-FileIfNotExists -Path $file
}

# Player Handler
$PlayerHandlers = @("online")

foreach ($handler in $PlayerHandlers) {
    $file = "backend\internal\handler\player\$handler.go"
    New-FileIfNotExists -Path $file
}

# WebSocket Handler
New-DirectoryIfNotExists -Path "backend\internal\handler\websocket"
$WebSocketFiles = @("chat.go", "notification.go")

foreach ($file in $WebSocketFiles) {
    $fullPath = "backend\internal\handler\websocket\$file"
    New-FileIfNotExists -Path $fullPath
}

# Upload Handler
New-DirectoryIfNotExists -Path "backend\internal\handler\upload"
New-FileIfNotExists -Path "backend\internal\handler\upload\upload.go"

Write-Host ""

# ============================================
# 第五部分: 调度器和中间件
# ============================================

Write-Host "⏰ 第五步: 创建调度器和中间件文件..." -ForegroundColor Blue

# 调度器
New-DirectoryIfNotExists -Path "backend\internal\scheduler"
$SchedulerFiles = @(
    "order_scheduler.go",
    "settlement_scheduler.go"
)

foreach ($file in $SchedulerFiles) {
    $fullPath = "backend\internal\scheduler\$file"
    New-FileIfNotExists -Path $fullPath
}

# Prometheus中间件
New-FileIfNotExists -Path "backend\internal\middleware\prometheus.go"

Write-Host ""

# ============================================
# 第六部分: 前端用户端页面
# ============================================

Write-Host "👥 第六步: 创建前端用户端页面..." -ForegroundColor Blue

$UserPages = @(
    "Home",
    "GameList",
    "PlayerList",
    "PlayerDetail",
    "OrderCreate",
    "MyOrders",
    "Profile"
)

foreach ($page in $UserPages) {
    $dir = "frontend\src\pages\UserPortal\$page"
    New-DirectoryIfNotExists -Path $dir
    
    $files = @("index.tsx", "$page.module.less")
    foreach ($file in $files) {
        $fullPath = "$dir\$file"
        New-FileIfNotExists -Path $fullPath
    }
}

Write-Host ""

# ============================================
# 第七部分: 前端陪玩师端页面
# ============================================

Write-Host "🎮 第七步: 创建前端陪玩师端页面..." -ForegroundColor Blue

$PlayerPages = @(
    "Dashboard",
    "Orders",
    "Earnings",
    "Services",
    "Profile",
    "Reviews",
    "Schedule"
)

foreach ($page in $PlayerPages) {
    $dir = "frontend\src\pages\PlayerPortal\$page"
    New-DirectoryIfNotExists -Path $dir
    
    $files = @("index.tsx", "$page.module.less")
    foreach ($file in $files) {
        $fullPath = "$dir\$file"
        New-FileIfNotExists -Path $fullPath
    }
}

Write-Host ""

# ============================================
# 第八部分: 前端通用组件
# ============================================

Write-Host "🧩 第八步: 创建前端通用组件..." -ForegroundColor Blue

$Components = @(
    "GameCard",
    "PlayerCard",
    "OrderStatusBadge",
    "ChatWindow",
    "DisputeModal",
    "TicketModal",
    "NotificationBell",
    "FavoriteButton"
)

foreach ($component in $Components) {
    $dir = "frontend\src\components\$component"
    New-DirectoryIfNotExists -Path $dir
    
    $files = @("index.ts", "$component.tsx", "$component.module.less")
    foreach ($file in $files) {
        $fullPath = "$dir\$file"
        New-FileIfNotExists -Path $fullPath
    }
}

Write-Host ""

# ============================================
# 第九部分: 前端服务层
# ============================================

Write-Host "🔧 第九步: 创建前端服务层文件..." -ForegroundColor Blue

$ApiFiles = @(
    "dispute",
    "ticket",
    "notification",
    "favorite",
    "chat",
    "earnings"
)

foreach ($api in $ApiFiles) {
    $file = "frontend\src\services\api\$api.ts"
    New-FileIfNotExists -Path $file
}

# WebSocket 服务
New-DirectoryIfNotExists -Path "frontend\src\services\websocket"
New-FileIfNotExists -Path "frontend\src\services\websocket\chat.ts"

Write-Host ""

# ============================================
# 第十部分: 前端类型定义
# ============================================

Write-Host "📝 第十步: 创建前端类型定义文件..." -ForegroundColor Blue

$TypeFiles = @(
    "dispute",
    "ticket",
    "notification",
    "favorite",
    "chat",
    "player"
)

foreach ($type in $TypeFiles) {
    $file = "frontend\src\types\$type.ts"
    New-FileIfNotExists -Path $file
}

Write-Host ""

# ============================================
# 完成
# ============================================

Write-Host ""
Write-Host "✅ 项目结构搭建完成!" -ForegroundColor Green
Write-Host ""
Write-Host "📊 统计信息:" -ForegroundColor Cyan
Write-Host "  - 后端模型文件: 6个"
Write-Host "  - Repository层: 6个目录, 12个文件"
Write-Host "  - Service层: 6个目录, 12+个文件"
Write-Host "  - Handler层: 10+个文件"
Write-Host "  - 前端用户端页面: 7个目录, 14个文件"
Write-Host "  - 前端陪玩师端页面: 7个目录, 14个文件"
Write-Host "  - 前端组件: 8个目录, 24个文件"
Write-Host "  - 前端服务层: 7个文件"
Write-Host "  - 前端类型定义: 6个文件"
Write-Host ""
Write-Host "📖 下一步:" -ForegroundColor Yellow
Write-Host "  1. 查看详细开发计划: Get-Content GAMELINK_IMPROVEMENT_PLAN.md"
Write-Host "  2. 查看快速摘要: Get-Content IMPROVEMENT_SUMMARY.md"
Write-Host "  3. 开始实现数据模型: cd backend\internal\model"
Write-Host "  4. 运行数据库迁移: cd backend; go run cmd\server\main.go migrate up"
Write-Host ""
Write-Host "🎯 第一周任务:" -ForegroundColor Yellow
Write-Host "  - Day 1-2: 实现6个新数据模型"
Write-Host "  - Day 3-4: 实现Repository层"
Write-Host "  - Day 5-7: 实现Service层"
Write-Host ""
Write-Host "💡 提示: 所有文件已创建为空文件,请根据GAMELINK_IMPROVEMENT_PLAN.md中的代码模板填充内容" -ForegroundColor Cyan
Write-Host ""
Write-Host "按任意键退出..."
$null = $Host.UI.RawUI.ReadKey("NoEcho,IncludeKeyDown")

