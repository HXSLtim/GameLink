param(
  [string]$BaseUrl = "http://127.0.0.1:8080/api/v1",
  [string]$AdminUsername = "admin@gamelink.com",
  [string]$AdminPassword = "Admin123456",
  [string]$UserUsername = "demo.user@gamelink.com",
  [string]$UserPassword = "User@123456",
  [string]$PlayerUsername = "pro.player@gamelink.com",
  [string]$PlayerPassword = "Player@123456",
  [string]$CSLeaderUsername = "cs.leader@gamelink.com",
  [string]$CSLeaderPassword = "CsLeader@123",
  [string]$CSAgentUsername = "cs.agent@gamelink.com",
  [string]$CSAgentPassword = "CsAgent@123",
  [string]$PostgresContainer = "gamelink-postgres",
  [string]$DBUser = "gamelink",
  [string]$DBName = "gamelink",
  [switch]$EnableCompatChecks
)

$ErrorActionPreference = "Stop"

$results = New-Object System.Collections.Generic.List[object]
$failed = $false
$scriptRoot = Split-Path -Parent $MyInvocation.MyCommand.Path

function Add-CheckResult {
  param(
    [string]$Stage,
    [string]$Check,
    [bool]$Passed,
    [string]$Detail
  )

  $script:results.Add([pscustomobject]@{
      Stage  = $Stage
      Check  = $Check
      Result = $(if ($Passed) { "PASS" } else { "FAIL" })
      Detail = $Detail
    })

  if (-not $Passed) {
    $script:failed = $true
  }
}

function Invoke-Api {
  param(
    [string]$Method,
    [string]$Path,
    [string]$Token = "",
    $Body = $null
  )

  $headers = @{ Accept = "application/json" }
  if ($Token) {
    $headers["Authorization"] = "Bearer $Token"
  }

  $params = @{
    Method      = $Method
    Uri         = "$BaseUrl$Path"
    Headers     = $headers
    ContentType = "application/json"
  }

  if ($null -ne $Body) {
    $params["Body"] = ($Body | ConvertTo-Json -Depth 10 -Compress)
  }

  try {
    $resp = Invoke-RestMethod @params
    return [pscustomobject]@{
      Ok      = $true
      Code    = 200
      Data    = $resp
      Message = ""
    }
  } catch {
    $raw = $_.ErrorDetails.Message
    $parsed = $null
    if ($raw) {
      try {
        $parsed = $raw | ConvertFrom-Json
      } catch {
        $parsed = $null
      }
    }

    $code = 0
    $message = $_.Exception.Message
    if ($parsed) {
      if ($parsed.code) {
        $code = [int]$parsed.code
      }
      if ($parsed.message) {
        $message = [string]$parsed.message
      }
    }

    return [pscustomobject]@{
      Ok      = $false
      Code    = $code
      Data    = $parsed
      Raw     = $raw
      Message = $message
    }
  }
}

function Login {
  param(
    [string]$Username,
    [string]$Password
  )

  $resp = Invoke-Api -Method "POST" -Path "/auth/login" -Body @{
    username = $Username
    password = $Password
  }

  if (-not $resp.Ok) {
    throw "登录失败: $Username -> $($resp.Message)"
  }

  $token = [string]$resp.Data.data.token
  $userID = [uint64]$resp.Data.data.user.id
  if ([string]::IsNullOrWhiteSpace($token)) {
    throw "登录返回 token 为空: $Username"
  }

  return [pscustomobject]@{
    Username = $Username
    Token    = $token
    UserID   = $userID
  }
}

function Invoke-SmokeChecks {
  Write-Host "[full-acceptance] 阶段1: 前端关键链路烟测"

  try {
    $admin = Login -Username $AdminUsername -Password $AdminPassword
    $user = Login -Username $UserUsername -Password $UserPassword
    $player = Login -Username $PlayerUsername -Password $PlayerPassword
  } catch {
    Add-CheckResult -Stage "Smoke" -Check "登录初始化" -Passed $false -Detail $_.Exception.Message
    return
  }

  $playersResp = Invoke-Api -Method "GET" -Path "/public/players?page=1&pageSize=20"
  if (-not $playersResp.Ok) {
    Add-CheckResult -Stage "Smoke" -Check "获取公开陪玩师列表" -Passed $false -Detail $playersResp.Message
    return
  }

  $players = @($playersResp.Data.data.players)
  if ($players.Count -eq 0) {
    Add-CheckResult -Stage "Smoke" -Check "公开陪玩师列表非空" -Passed $false -Detail "players=0"
    return
  }

  $playerID = [uint64]$players[0].id
  $checks = @(
    @{ Name = "public/games/categories"; Method = "GET"; Path = "/public/games/categories"; Token = "" },
    @{ Name = "public/players"; Method = "GET"; Path = "/public/players?page=1&pageSize=10"; Token = "" },
    @{ Name = "public/players/:id/services"; Method = "GET"; Path = "/public/players/$playerID/services"; Token = "" },
    @{ Name = "user/settings"; Method = "GET"; Path = "/user/settings"; Token = $user.Token },
    @{ Name = "user/notification-settings"; Method = "GET"; Path = "/user/notification-settings"; Token = $user.Token },
    @{ Name = "player/orders (compat)"; Method = "GET"; Path = "/player/orders?page=1&page_size=10"; Token = $player.Token },
    @{ Name = "player/certification/identity"; Method = "GET"; Path = "/player/certification/identity"; Token = $player.Token },
    @{ Name = "users/me (legacy)"; Method = "GET"; Path = "/users/me"; Token = $admin.Token },
    @{ Name = "users/chat/groups (legacy)"; Method = "GET"; Path = "/users/chat/groups?page=1&page_size=10"; Token = $admin.Token },
    @{ Name = "user/chat/groups"; Method = "GET"; Path = "/user/chat/groups?page=1&page_size=10"; Token = $admin.Token },
    @{ Name = "admin/stats/dashboard"; Method = "GET"; Path = "/admin/stats/dashboard"; Token = $admin.Token }
  )

  if ($EnableCompatChecks) {
    $checks += @(
      @{ Name = "admin/stats (compat)"; Method = "GET"; Path = "/admin/stats"; Token = $admin.Token },
      @{ Name = "admin/dashboard/stats (compat)"; Method = "GET"; Path = "/admin/dashboard/stats"; Token = $admin.Token },
      @{ Name = "admin/user-behavior (compat)"; Method = "GET"; Path = "/admin/user-behavior"; Token = $admin.Token },
      @{ Name = "admin/user-distribution (compat)"; Method = "GET"; Path = "/admin/user-distribution"; Token = $admin.Token }
    )
  }

  foreach ($item in $checks) {
    $resp = Invoke-Api -Method $item.Method -Path $item.Path -Token $item.Token
    Add-CheckResult -Stage "Smoke" -Check $item.Name -Passed $resp.Ok -Detail $(if ($resp.Ok) { "code=200" } else { "code=$($resp.Code), message=$($resp.Message)" })
  }
}

function Invoke-SubScript {
  param(
    [string]$Stage,
    [string]$ScriptName,
    [hashtable]$Arguments
  )

  $scriptPath = Join-Path $scriptRoot $ScriptName
  if (-not (Test-Path $scriptPath)) {
    Add-CheckResult -Stage $Stage -Check $ScriptName -Passed $false -Detail "脚本不存在: $scriptPath"
    return
  }

  $argList = @("-ExecutionPolicy", "Bypass", "-File", $scriptPath)
  foreach ($key in $Arguments.Keys) {
    $argList += "-$key"
    $argList += [string]$Arguments[$key]
  }

  Write-Host "[full-acceptance] 阶段: $Stage ($ScriptName)"
  $output = & powershell @argList 2>&1
  if ($output) {
    $output | ForEach-Object { Write-Host $_ }
  }
  $exitCode = $LASTEXITCODE

  Add-CheckResult -Stage $Stage -Check $ScriptName -Passed ($exitCode -eq 0) -Detail "exitCode=$exitCode"
}

function Check-DataIntegritySummary {
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
      Write-Host "[full-acceptance] preferred container '$PreferredName' not found, using first postgres container '$($candidates[0])'"
    }

    return $candidates[0]
  }

  $sql = @"
SELECT
  (SELECT COUNT(*) FROM orders o LEFT JOIN users u ON u.id = o.user_id WHERE u.id IS NULL)
  + (SELECT COUNT(*) FROM orders o LEFT JOIN players p ON p.id = o.player_id WHERE o.player_id IS NOT NULL AND p.id IS NULL)
  + (SELECT COUNT(*) FROM payments p LEFT JOIN orders o ON o.id = p.order_id WHERE o.id IS NULL)
  + (SELECT COUNT(*) FROM payments p LEFT JOIN users u ON u.id = p.user_id WHERE u.id IS NULL)
  + (SELECT COUNT(*) FROM reviews r LEFT JOIN orders o ON o.id = r.order_id WHERE o.id IS NULL)
  + (SELECT COUNT(*) FROM reviews r LEFT JOIN users u ON u.id = r.user_id WHERE u.id IS NULL)
  + (SELECT COUNT(*) FROM reviews r LEFT JOIN players p ON p.id = r.player_id WHERE p.id IS NULL)
  + (SELECT COUNT(*) FROM player_schedules s LEFT JOIN players p ON p.id = s.player_id WHERE p.id IS NULL)
  + (SELECT COUNT(*) FROM reviews r JOIN orders o ON o.id = r.order_id WHERE o.status <> 'completed')
  + (SELECT COUNT(*) FROM (SELECT order_id FROM payments WHERE status = 'pending' GROUP BY order_id HAVING COUNT(*) > 1) x)
  + (SELECT COUNT(*) FROM orders o WHERE o.status = 'completed' AND NOT EXISTS (SELECT 1 FROM payments p WHERE p.order_id = o.id AND p.status IN ('paid','refunded')))
  + (SELECT COUNT(*) FROM payments p JOIN orders o ON o.id = p.order_id WHERE p.status = 'paid' AND o.status = 'canceled')
  AS total_violations;
"@

  try {
    $resolvedContainer = Resolve-PostgresContainer -PreferredName $PostgresContainer
    $totalViolations = docker exec -i $resolvedContainer psql -U $DBUser -d $DBName -At -c $sql
    $trimmed = [string]$totalViolations
    $trimmed = $trimmed.Trim()
    $count = [int]$trimmed
    Add-CheckResult -Stage "Integrity" -Check "core/business 违规项归零" -Passed ($count -eq 0) -Detail "violations=$count"
  } catch {
    Add-CheckResult -Stage "Integrity" -Check "core/business 违规项归零" -Passed $false -Detail $_.Exception.Message
  }
}

Write-Host "[full-acceptance] 开始执行全流程验收..."

Invoke-SmokeChecks

$commonArgs = @{
  BaseUrl        = $BaseUrl
  AdminUsername  = $AdminUsername
  AdminPassword  = $AdminPassword
  UserUsername   = $UserUsername
  UserPassword   = $UserPassword
  PlayerUsername = $PlayerUsername
  PlayerPassword = $PlayerPassword
}

Invoke-SubScript -Stage "OrderFlow" -ScriptName "run_flow_guard_regression.ps1" -Arguments $commonArgs
Invoke-SubScript -Stage "NotificationFlow" -ScriptName "run_order_accept_notification_regression.ps1" -Arguments $commonArgs

$csArgs = @{
  BaseUrl          = $BaseUrl
  AdminUsername    = $AdminUsername
  AdminPassword    = $AdminPassword
  UserUsername     = $UserUsername
  UserPassword     = $UserPassword
  PlayerUsername   = $PlayerUsername
  PlayerPassword   = $PlayerPassword
  CSLeaderUsername = $CSLeaderUsername
  CSLeaderPassword = $CSLeaderPassword
  CSAgentUsername  = $CSAgentUsername
  CSAgentPassword  = $CSAgentPassword
}
Invoke-SubScript -Stage "DisputeFlow" -ScriptName "run_cs_permission_regression.ps1" -Arguments $csArgs

Invoke-SubScript -Stage "RefundFlow" -ScriptName "run_partial_refund_regression.ps1" -Arguments $commonArgs
Invoke-SubScript -Stage "WithdrawFlow" -ScriptName "run_withdraw_flow_regression.ps1" -Arguments $commonArgs
Invoke-SubScript -Stage "Integrity" -ScriptName "run_data_integrity.ps1" -Arguments @{
  Container = $PostgresContainer
  DBUser    = $DBUser
  DBName    = $DBName
}
Check-DataIntegritySummary

Write-Host ""
Write-Host "[full-acceptance] 汇总:"
$results | Format-Table -AutoSize

if ($failed) {
  Write-Error "full service flow acceptance failed"
  exit 1
}

Write-Host "[full-acceptance] 全部通过"
