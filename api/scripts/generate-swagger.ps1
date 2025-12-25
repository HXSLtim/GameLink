# GameLink Swagger 文档生成脚本
# 用于生成支持泛型的 Swagger 文档

Write-Host "🚀 开始生成 GameLink Swagger 文档..." -ForegroundColor Green

# 切换到 cmd 目录
$cmdDir = "C:\Users\a2778\Desktop\code\GameLink\backend\cmd"
$outputDir = "C:\Users\a2778\Desktop\code\GameLink\backend\docs"

Write-Host "📁 工作目录: $cmdDir" -ForegroundColor Yellow
Write-Host "📄 输出目录: $outputDir" -ForegroundColor Yellow

# 设置环境变量确保 UTF-8
$env:GOOS = "windows"
$env:GOARCH = "amd64"

# 切换到 cmd 目录
Set-Location -Path $cmdDir

# 运行 swag init
try {
    & swag init `
        --output $outputDir `
        --generalInfo "main.go" `
        --dir "." `
        --parseDependency `
        --parseInternal `
        --parseDepth 10

    if ($LASTEXITCODE -eq 0) {
        Write-Host "✅ Swagger 文档生成成功!" -ForegroundColor Green
        Write-Host "📊 文件位置:" -ForegroundColor Yellow
        Write-Host "  - docs/docs.go" -ForegroundColor Cyan
        Write-Host "  - docs/swagger.json" -ForegroundColor Cyan
        Write-Host "  - docs/swagger.yaml" -ForegroundColor Cyan
    } else {
        Write-Host "❌ Swagger 文档生成失败!" -ForegroundColor Red
        exit 1
    }
} catch {
    Write-Host "❌ 执行错误: $_" -ForegroundColor Red
    exit 1
}
