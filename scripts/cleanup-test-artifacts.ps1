Param(
    [switch]$WhatIf
)

# 清理常见测试生成文件与目录（不影响 docs/ 归档内容）
$repoRoot = Resolve-Path (Join-Path $PSScriptRoot "..")

function Remove-IfExists {
    param(
        [string]$Path
    )
    $fullPath = Join-Path $repoRoot $Path
    if (Test-Path $fullPath) {
        Write-Host "Removing: $Path" -ForegroundColor Yellow
        if ($WhatIf) {
            Write-Host "WhatIf: would remove $fullPath" -ForegroundColor DarkYellow
        } else {
            Remove-Item -Recurse -Force -ErrorAction SilentlyContinue $fullPath
        }
    }
}

# Backend 测试产物
$backendArtifacts = @(
    "backend/coverage.out",
    "backend/cover.out",
    "backend/coverage.html",
    "backend/coverage",
    "backend/junit.xml",
    "backend/tmp",
    "backend/*.tmp"
)

# Frontend 测试产物
$frontendArtifacts = @(
    "frontend/coverage",
    "frontend/junit.xml",
    "frontend/.vitest",
    "frontend/.eslintcache",
    "frontend/tmp",
    "frontend/*.tmp"
)

# 根目录可能的产物（兼容性处理）
$rootArtifacts = @(
    "coverage.out",
    "cover.out",
    "junit.xml",
    "tmp",
    "*.tmp"
)

# 执行删除（排除 docs/）
($backendArtifacts + $frontendArtifacts + $rootArtifacts) | ForEach-Object {
    if (-not $_.StartsWith("docs/")) {
        Remove-IfExists $_
    }
}

Write-Host "Test artifacts cleanup completed." -ForegroundColor Green