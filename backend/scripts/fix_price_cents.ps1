# 批量修复 PriceCents 字段引用的脚本
# Order 模型已重构: PriceCents -> TotalPriceCents

Write-Host "开始修复 PriceCents 字段引用..." -ForegroundColor Green

$replacements = @(
    @{
        Pattern = '\.PriceCents(?!\s*=)'
        Replacement = '.TotalPriceCents'
        Description = "访问PriceCents字段"
    },
    @{
        Pattern = 'PriceCents:\s*(\d+)'
        Replacement = 'TotalPriceCents: $1'
        Description = "结构体初始化PriceCents"
    },
    @{
        Pattern = 'json:"priceCents"'
        Replacement = 'json:"totalPriceCents"'
        Description = "JSON标签"
    },
    @{
        Pattern = 'o\.GetPriceCents\(\)'
        Replacement = 'o.TotalPriceCents'
        Description = "GetPriceCents方法调用"
    }
)

$files = Get-ChildItem -Path "internal" -Recurse -Include "*.go" -Exclude "*_test.go"

$totalFixed = 0

foreach ($file in $files) {
    $content = Get-Content $file.FullName -Raw -Encoding UTF8
    $originalContent = $content
    $fileFixed = 0
    
    foreach ($rep in $replacements) {
        $matches = [regex]::Matches($content, $rep.Pattern)
        if ($matches.Count -gt 0) {
            $content = $content -replace $rep.Pattern, $rep.Replacement
            $fileFixed += $matches.Count
            Write-Host "  [$($file.Name)] 修复 $($matches.Count) 处: $($rep.Description)" -ForegroundColor Yellow
        }
    }
    
    if ($content -ne $originalContent) {
        Set-Content -Path $file.FullName -Value $content -Encoding UTF8 -NoNewline
        $totalFixed += $fileFixed
    }
}

Write-Host "`n✅ 修复完成！共修复 $totalFixed 处" -ForegroundColor Green
Write-Host "⚠️  请手动检查测试文件和特殊情况" -ForegroundColor Yellow

