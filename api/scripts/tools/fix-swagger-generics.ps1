# Fix Swagger Generic Syntax
# Usage: .\fix-swagger-generics.ps1

Write-Host "Starting to fix Swagger generic syntax..." -ForegroundColor Green

$count = 0

# Get all handler files (exclude test files)
$files = Get-ChildItem -Path "internal\handler" -Filter "*.go" -Recurse |
    Where-Object { $_.Name -notlike "*_test.go" }

Write-Host "Found $($files.Count) files to process" -ForegroundColor Yellow

foreach ($file in $files) {
    $modified = $false
    $content = Get-Content $file.FullName -Raw
    $originalContent = $content

    # Replace generic syntax in Failure annotations
    $content = $content -replace '(// @Failure\s+)(\d+)(\s+\{object\}\s+)model\.APIResponse\[any\]', '$1$2            {object}  model.ErrorResponse'
    $content = $content -replace '(// @Failure\s+)(\d+)(\s+\{object\}\s+)model\.APIResponse\[interface\{\}\]', '$1$2            {object}  model.ErrorResponse'

    # Replace generic syntax in common Success annotations
    $content = $content -replace '(// @Success\s+)(\d+)(\s+\{object\}\s+)model\.APIResponse\[any\]', '$1$2            {object}  model.SuccessResponse'

    if ($content -ne $originalContent) {
        Set-Content -Path $file.FullName -Value $content -NoNewline
        Write-Host "  [FIXED] $($file.FullName)" -ForegroundColor Cyan
        $count++
        $modified = $true
    }
}

Write-Host "`nFix completed! Modified $count files" -ForegroundColor Green
Write-Host "`nNOTE:" -ForegroundColor Yellow
Write-Host "  1. Some Success responses still need manual Response type creation"  -ForegroundColor Yellow
Write-Host "  2. Run 'swag init -g cmd/main.go -o docs/swagger' to verify" -ForegroundColor Yellow
Write-Host "`nSee SWAGGER_FIX_GUIDE.md for detailed instructions" -ForegroundColor Cyan
