# Fix Swagger Pointer Types in Generic Annotations
# This script removes pointer symbols (*) from APIResponse generic type parameters
# to make them compatible with Swag CLI
#
# 使用方法：
# 1. 以管理员身份打开 PowerShell
# 2. 切换到 backend 目录: cd "C:\Users\a2778\Desktop\code\GameLink\backend"
# 3. 运行脚本: .\scripts\fix-swagger-pointers.ps1
#
# 脚本将：
# - 创建 handler 目录的备份
# - 查找所有 APIResponse[* 注解并替换为 APIResponse[
# - 生成修复报告
#
# 修复完成后，重新运行: .\scripts\generate-swagger.ps1

Write-Host "🔧 开始修复 Swagger 指针类型注解..." -ForegroundColor Yellow
Write-Host "================================================" -ForegroundColor Yellow

# 配置
$projectRoot = "C:\Users\a2778\Desktop\code\GameLink\backend"
$handlerDir = "$projectRoot\internal\handler"
$backupDir = "$projectRoot\internal\handler.backup-$(Get-Date -Format 'yyyyMMdd-HHmmss')"

# 初始化计数器
$totalFiles = 0
$modifiedFiles = 0
$totalReplacements = 0

# Create backup
Write-Host "📁 创建备份目录: $backupDir" -ForegroundColor Cyan
Copy-Item -Path $handlerDir -Destination $backupDir -Recurse -Force
Write-Host "✅ 备份完成" -ForegroundColor Green
Write-Host ""

# Get all Go files in handler directory
$goFiles = Get-ChildItem -Path $handlerDir -Filter "*.go" -Recurse

foreach ($file in $goFiles) {
    $content = Get-Content -Path $file.FullName -Raw
    $originalContent = $content

    # Count occurrences of APIResponse[*
    $occurrences = [regex]::Matches($content, 'APIResponse\[\*').Count

    if ($occurrences -gt 0) {
        Write-Host "📄 处理文件: $($file.FullName.Replace($handlerDir, ''))" -ForegroundColor White
        Write-Host "   发现 $occurrences 处指针类型注解" -ForegroundColor Gray

        # Replace APIResponse[* with APIResponse[
        $content = $content -replace 'APIResponse\[\*', 'APIResponse['

        # Write back to file
        Set-Content -Path $file.FullName -Value $content -NoNewline

        Write-Host "   ✅ 已修复" -ForegroundColor Green

        $modifiedFiles++
        $totalReplacements += $occurrences
    }

    $totalFiles++
}

Write-Host ""
Write-Host "================================================" -ForegroundColor Yellow
Write-Host "📊 修复完成摘要:" -ForegroundColor Yellow
Write-Host "================================================" -ForegroundColor Yellow
Write-Host "总文件数: $totalFiles" -ForegroundColor White
Write-Host "修改文件数: $modifiedFiles" -ForegroundColor Green
Write-Host "总替换次数: $totalReplacements" -ForegroundColor Green
Write-Host ""

if ($modifiedFiles -gt 0) {
    Write-Host "✅ 修复成功完成！" -ForegroundColor Green
    Write-Host "📁 备份位置: $backupDir" -ForegroundColor Cyan
    Write-Host ""
    Write-Host "下一步: 重新运行 Swagger 生成脚本" -ForegroundColor Yellow
    Write-Host "> .\scripts\generate-swagger.ps1" -ForegroundColor Cyan
} else {
    Write-Host "ℹ️  没有找到需要修复的指针类型注解" -ForegroundColor Yellow
}

Write-Host ""
