param(
  [switch]$Fix,
  [string]$Container = "gamelink-postgres",
  [string]$DBUser = "gamelink",
  [string]$DBName = "gamelink"
)

$ErrorActionPreference = "Stop"

function Resolve-PostgresContainer {
  param([string]$PreferredName)

  $names = docker ps --format "{{.Names}}"
  if ($names -contains $PreferredName) {
    return $PreferredName
  }

  $candidates = docker ps --format "{{.Names}} {{.Image}}" |
    Where-Object { $_ -match "postgres" } |
    ForEach-Object { ($_ -split "\s+")[0] }

  if ($candidates.Count -eq 0) {
    throw "no running postgres container found"
  }

  if ($candidates.Count -gt 1) {
    Write-Host "[integrity] preferred container '$PreferredName' not found, using first postgres container '$($candidates[0])'"
  }

  return $candidates[0]
}

$resolvedContainer = Resolve-PostgresContainer -PreferredName $Container

if ($Fix) {
  Write-Host "[integrity] applying fix script..."
  Get-Content "$PSScriptRoot/data_integrity_fix.sql" |
    docker exec -i $resolvedContainer psql -U $DBUser -d $DBName
}

Write-Host "[integrity] running check script..."
Get-Content "$PSScriptRoot/data_integrity_check.sql" |
  docker exec -i $resolvedContainer psql -U $DBUser -d $DBName
