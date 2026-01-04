# GameLink 安装诊断和修复脚本

Write-Host "========================================" -ForegroundColor Cyan
Write-Host "GameLink 安装诊断和修复" -ForegroundColor Cyan
Write-Host "========================================" -ForegroundColor Cyan
Write-Host ""

$ErrorActionPreference = "Continue"

# 检查Node.js和npm版本
Write-Host "📋 检查环境..." -ForegroundColor Yellow
Write-Host ""

try {
    $nodeVersion = node --version
    $npmVersion = npm --version

    Write-Host "✅ Node.js: $nodeVersion" -ForegroundColor Green
    Write-Host "✅ npm: $npmVersion" -ForegroundColor Green

    if ([version]$nodeVersion -lt [version]"18.0.0") {
        Write-Host ""
        Write-Host "⚠️  警告: Node.js版本过低，建议使用18+版本" -ForegroundColor Yellow
        Write-Host "   当前版本: $nodeVersion" -ForegroundColor Yellow
    }
}
catch {
    Write-Host "❌ 未找到Node.js或npm" -ForegroundColor Red
    Write-Host ""
    Write-Host "请先安装Node.js: https://nodejs.org/" -ForegroundColor Yellow
    pause
    exit 1
}

Write-Host ""
Write-Host "📋 检查网络连接..." -ForegroundColor Yellow
Write-Host ""

# 测试npm registry连接
try {
    $registry = npm config get registry
    Write-Host "当前registry: $registry" -ForegroundColor White

    $ping = Test-Connection -ComputerName "registry.npmjs.org" -Count 1 -Quiet
    if ($ping) {
        Write-Host "✅ 可以访问npm registry" -ForegroundColor Green
    } else {
        Write-Host "⚠️  无法访问npm registry，可能需要配置代理或使用镜像" -ForegroundColor Yellow
    }
}
catch {
    Write-Host "⚠️  网络检测失败" -ForegroundColor Yellow
}

Write-Host ""
Write-Host "📋 诊断结果和建议:" -ForegroundColor Yellow
Write-Host ""

# 检查package.json是否存在
if (Test-Path "package.json") {
    Write-Host "✅ 找到package.json" -ForegroundColor Green
} else {
    Write-Host "❌ 未找到package.json，请在admin目录中运行" -ForegroundColor Red
    pause
    exit 1
}

# 检查node_modules
if (Test-Path "node_modules") {
    Write-Host "⚠️  node_modules已存在" -ForegroundColor Yellow
    Write-Host ""
    $clean = Read-Host "是否删除node_modules并重新安装? (y/N)"

    if ($clean -eq "y" -or $clean -eq "Y") {
        Write-Host ""
        Write-Host "🗑️  删除node_modules..." -ForegroundColor Yellow
        Remove-Item -Recurse -Force node_modules
        Remove-Item package-lock.json -ErrorAction SilentlyContinue
        Write-Host "✅ 已删除" -ForegroundColor Green
    }
}

Write-Host ""
Write-Host "========================================" -ForegroundColor Cyan
Write-Host "修复选项" -ForegroundColor Cyan
Write-Host "========================================" -ForegroundColor Cyan
Write-Host ""

Write-Host "请选择修复方案:" -ForegroundColor Yellow
Write-Host "1. 使用淘宝镜像安装 (推荐，国内更快)" -ForegroundColor White
Write-Host "2. 清除npm缓存后重试" -ForegroundColor White
Write-Host "3. 使用legacy-peer-dependencies模式" -ForegroundColor White
Write-Host "4. 跳过PWA依赖，仅安装基础依赖" -ForegroundColor White
Write-Host "5. 查看详细错误并手动修复" -ForegroundColor White
Write-Host ""

$choice = Read-Host "请选择 (1-5, 默认: 1)"

if ([string]::IsNullOrEmpty($choice)) { $choice = "1" }

Write-Host ""

switch ($choice) {
    "1" {
        Write-Host "🔧 配置淘宝镜像..." -ForegroundColor Yellow
        npm config set registry https://registry.npmmirror.com

        Write-Host ""
        Write-Host "📦 安装依赖..." -ForegroundColor Yellow
        npm install -D vite-plugin-pwa picocolors workbox-window

        if ($LASTEXITCODE -eq 0) {
            Write-Host ""
            Write-Host "✅ 安装成功!" -ForegroundColor Green
            Write-Host ""
            Write-Host "💡 提示: 淘宝镜像已永久配置" -ForegroundColor White
            Write-Host "   如需恢复官方源:" -ForegroundColor White
            Write-Host "   npm config set registry https://registry.npmjs.org" -ForegroundColor Gray
        } else {
            Write-Host ""
            Write-Host "❌ 安装仍然失败" -ForegroundColor Red
        }
    }

    "2" {
        Write-Host "🗑️  清除npm缓存..." -ForegroundColor Yellow
        npm cache clean --force

        Write-Host ""
        Write-Host "📦 重新安装依赖..." -ForegroundColor Yellow
        npm install -D vite-plugin-pwa picocolors workbox-window
    }

    "3" {
        Write-Host "📦 使用legacy-peer-dependencies模式安装..." -ForegroundColor Yellow
        npm install -D vite-plugin-pwa picocolors workbox-window --legacy-peer-deps
    }

    "4" {
        Write-Host "📦 仅安装基础依赖 (跳过PWA)..." -ForegroundColor Yellow
        npm install

        Write-Host ""
        Write-Host "⚠️  PWA功能将不可用" -ForegroundColor Yellow
        Write-Host "   如需PWA，稍后可以单独安装:" -ForegroundColor Yellow
        Write-Host "   npm install -D vite-plugin-pwa picocolors workbox-window" -ForegroundColor Gray
    }

    "5" {
        Write-Host "📋 显示详细错误信息..." -ForegroundColor Yellow
        Write-Host ""
        Write-Host "运行以下命令查看详细错误:" -ForegroundColor White
        Write-Host "npm install -D vite-plugin-pwa picocolors workbox-window --verbose" -ForegroundColor Cyan
        Write-Host ""
        Write-Host "常见错误解决方案:" -ForegroundColor Yellow
        Write-Host "1. ECONNREFUSED: 网络问题，使用淘宝镜像 (选项1)" -ForegroundColor White
        Write-Host "2. ETIMEDOUT: 连接超时，检查网络或使用代理" -ForegroundColor White
        Write-Host "3. E404: 包不存在，检查包名是否正确" -ForegroundColor White
        Write-Host "4. peer dependency: 使用选项3或--legacy-peer-deps" -ForegroundColor White
    }

    default {
        Write-Host "❌ 无效的选择" -ForegroundColor Red
    }
}

Write-Host ""
Write-Host "========================================" -ForegroundColor Cyan
Write-Host "完成" -ForegroundColor Cyan
Write-Host "========================================" -ForegroundColor Cyan
Write-Host ""

Write-Host "💡 下一步操作:" -ForegroundColor Yellow
Write-Host "1. 如果安装成功，运行: npm run dev" -ForegroundColor White
Write-Host "2. 如果仍然失败，尝试其他选项" -ForegroundColor White
Write-Host "3. 查看完整文档: SERVICE_STARTUP_GUIDE.md" -ForegroundColor White
Write-Host ""

pause
