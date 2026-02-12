param(
  [switch]$Fix
)

$ErrorActionPreference = "Stop"

$container = "gamelink-postgres"
$dbUser = "gamelink"
$dbName = "gamelink"

if ($Fix) {
  Write-Host "[integrity] applying fix script..."
  Get-Content "$PSScriptRoot/data_integrity_fix.sql" |
    docker exec -i $container psql -U $dbUser -d $dbName
}

Write-Host "[integrity] running check script..."
Get-Content "$PSScriptRoot/data_integrity_check.sql" |
  docker exec -i $container psql -U $dbUser -d $dbName
