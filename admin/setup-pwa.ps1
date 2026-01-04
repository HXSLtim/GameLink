# GameLink Admin PWA 设置脚本
# 用途: 安装 PWA 相关依赖

Write-Host "========================================" -ForegroundColor Cyan
Write-Host "GameLink Admin - PWA 依赖安装" -ForegroundColor Cyan
Write-Host "========================================" -ForegroundColor Cyan
Write-Host ""

$ErrorActionPreference = "Stop"

# 检查是否在 admin 目录
if (-not (Test-Path "package.json")) {
    Write-Host "❌ 错误: 请在 admin 目录中运行此脚本" -ForegroundColor Red
    exit 1
}

Write-Host "📦 安装 PWA 依赖..." -ForegroundColor Yellow
Write-Host ""

try {
    # 安装依赖
    npm install -D vite-plugin-pwa picocolors workbox-window

    Write-Host ""
    Write-Host "✅ PWA 依赖安装成功!" -ForegroundColor Green
    Write-Host ""
    Write-Host "📝 下一步操作:" -ForegroundColor Yellow
    Write-Host "1. 按照 public/PWA_ICONS_README.md 生成 PNG 图标" -ForegroundColor White
    Write-Host "2. 运行 'npm run dev' 启动开发服务器" -ForegroundColor White
    Write-Host "3. 在浏览器中访问显示的任意 IP 地址" -ForegroundColor White
    Write-Host ""
    Write-Host "🎉 完成!" -ForegroundColor Green

} catch {
    Write-Host ""
    Write-Host "❌ 安装失败: $_" -ForegroundColor Red
    Write-Host ""
    Write-Host "💡 提示: 您可以手动运行以下命令:" -ForegroundColor Yellow
    Write-Host "npm install -D vite-plugin-pwa picocolors workbox-window" -ForegroundColor White
    exit 1
}
