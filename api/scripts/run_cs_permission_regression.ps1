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
  [string]$PreferredGameName = "英雄联盟"
)

$ErrorActionPreference = "Stop"

$results = New-Object System.Collections.Generic.List[object]
$failed = $false

function Add-CheckResult {
  param(
    [string]$Name,
    [bool]$Passed,
    [string]$Detail
  )

  $script:results.Add([pscustomobject]@{
      Check  = $Name
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

  $headers = @{}
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
      Message = $message
      Raw     = $raw
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
  $userId = [uint64]$resp.Data.data.user.id
  if ([string]::IsNullOrWhiteSpace($token)) {
    throw "登录返回 token 为空: $Username"
  }

  return [pscustomobject]@{
    Token    = $token
    UserID   = $userId
    Username = $Username
  }
}

function Resolve-TargetPlayer {
  param(
    [string]$UserToken,
    [uint64]$PlayerUserID
  )

  $resp = Invoke-Api -Method "GET" -Path "/user/players?page=1&pageSize=100" -Token $UserToken
  if (-not $resp.Ok) {
    throw "获取陪玩师列表失败: $($resp.Message)"
  }

  $players = @($resp.Data.data.players)
  if ($players.Count -eq 0) {
    throw "陪玩师列表为空，无法继续"
  }

  $matched = $players | Where-Object { [uint64]$_.userId -eq $PlayerUserID } | Select-Object -First 1
  if ($null -eq $matched) {
    $matched = $players | Select-Object -First 1
  }

  return [pscustomobject]@{
    PlayerID = [uint64]$matched.id
    UserID   = [uint64]$matched.userId
    Nickname = [string]$matched.nickname
  }
}

function Resolve-TargetGame {
  $resp = Invoke-Api -Method "GET" -Path "/public/games?page=1&pageSize=100"
  if (-not $resp.Ok) {
    throw "获取游戏列表失败: $($resp.Message)"
  }

  $games = @($resp.Data.data.games)
  if ($games.Count -eq 0) {
    throw "游戏列表为空，无法继续"
  }

  $target = $games | Where-Object { [string]$_.name -eq $PreferredGameName } | Select-Object -First 1
  if ($null -eq $target) {
    $target = $games | Select-Object -First 1
  }

  return [pscustomobject]@{
    GameID   = [uint64]$target.id
    GameName = [string]$target.name
  }
}

function New-DisputeCase {
  param(
    [string]$UserToken,
    [string]$AdminToken,
    [uint64]$PlayerID,
    [uint64]$GameID,
    [string]$TitlePrefix
  )

  $title = "$TitlePrefix-$(Get-Date -Format yyyyMMddHHmmssfff)"
  $createOrder = Invoke-Api -Method "POST" -Path "/user/orders" -Token $UserToken -Body @{
    playerId       = $PlayerID
    gameId         = $GameID
    title          = $title
    description    = "cs permission regression"
    scheduledStart = (Get-Date).AddHours(2).ToString("o")
    durationHours  = 1
  }
  if (-not $createOrder.Ok) {
    throw "创建订单失败: $($createOrder.Message)"
  }

  $orderID = [uint64]$createOrder.Data.data.orderId

  $confirm = Invoke-Api -Method "POST" -Path "/admin/orders/$orderID/confirm" -Token $AdminToken -Body @{ note = "confirm for dispute regression" }
  if (-not $confirm.Ok) {
    throw "确认订单失败(orderId=$orderID): $($confirm.Message)"
  }

  $start = Invoke-Api -Method "POST" -Path "/admin/orders/$orderID/start" -Token $AdminToken -Body @{ note = "start for dispute regression" }
  if (-not $start.Ok) {
    throw "开始订单失败(orderId=$orderID): $($start.Message)"
  }

  $dispute = Invoke-Api -Method "POST" -Path "/user/orders/$orderID/dispute" -Token $UserToken -Body @{
    orderId      = $orderID
    type         = "service_quality"
    reason       = "客服权限回归测试"
    description  = "auto regression"
    evidenceUrls = @("https://example.com/evidence.png")
  }
  if (-not $dispute.Ok) {
    throw "发起争议失败(orderId=$orderID): $($dispute.Message)"
  }

  return [pscustomobject]@{
    OrderID   = $orderID
    DisputeID = [uint64]$dispute.Data.data.disputeId
  }
}

Write-Host "[cs-perm-regression] 登录账号..."
$admin = Login -Username $AdminUsername -Password $AdminPassword
$user = Login -Username $UserUsername -Password $UserPassword
$player = Login -Username $PlayerUsername -Password $PlayerPassword
$csLeader = Login -Username $CSLeaderUsername -Password $CSLeaderPassword
$csAgent = Login -Username $CSAgentUsername -Password $CSAgentPassword

$targetPlayer = Resolve-TargetPlayer -UserToken $user.Token -PlayerUserID $player.UserID
$targetGame = Resolve-TargetGame
Write-Host "[cs-perm-regression] 使用陪玩师: $($targetPlayer.Nickname) (playerId=$($targetPlayer.PlayerID)), 游戏: $($targetGame.GameName) (gameId=$($targetGame.GameID))"

$health = Invoke-Api -Method "GET" -Path "/healthz"
Add-CheckResult -Name "C0. healthz" -Passed $health.Ok -Detail "status=$($health.Code)"

$leaderLogs = Invoke-Api -Method "GET" -Path "/admin/operation-logs?page=1&pageSize=5" -Token $csLeader.Token
Add-CheckResult -Name "C1. csLeader can view operation logs" -Passed $leaderLogs.Ok -Detail "status=$($leaderLogs.Code)"

$agentLogs = Invoke-Api -Method "GET" -Path "/admin/operation-logs?page=1&pageSize=5" -Token $csAgent.Token
$agentLogsExpected = (-not $agentLogs.Ok) -and ($agentLogs.Code -eq 403)
Add-CheckResult -Name "C2. csAgent cannot view operation logs" -Passed $agentLogsExpected -Detail "status=$($agentLogs.Code), message=$($agentLogs.Message)"

$agentListDisputes = Invoke-Api -Method "GET" -Path "/admin/disputes?page=1&pageSize=5" -Token $csAgent.Token
Add-CheckResult -Name "C3. csAgent can list disputes" -Passed $agentListDisputes.Ok -Detail "status=$($agentListDisputes.Code)"

$disputeCaseA = New-DisputeCase -UserToken $user.Token -AdminToken $admin.Token -PlayerID $targetPlayer.PlayerID -GameID $targetGame.GameID -TitlePrefix "cs-agent-resolve"
Add-CheckResult -Name "C4. create dispute case A" -Passed $true -Detail "orderId=$($disputeCaseA.OrderID), disputeId=$($disputeCaseA.DisputeID)"

$agentAssign = Invoke-Api -Method "POST" -Path "/admin/disputes/$($disputeCaseA.DisputeID)/assign" -Token $csAgent.Token -Body @{
  assignedServiceId = $csAgent.UserID
}
$agentAssignExpected = (-not $agentAssign.Ok) -and ($agentAssign.Code -eq 403)
Add-CheckResult -Name "C5. csAgent cannot assign disputes" -Passed $agentAssignExpected -Detail "status=$($agentAssign.Code), message=$($agentAssign.Message)"

$agentResolve = Invoke-Api -Method "POST" -Path "/admin/disputes/$($disputeCaseA.DisputeID)/resolve" -Token $csAgent.Token -Body @{
  resolution    = "reject"
  resolveRemark = "csAgent regression resolve"
}
Add-CheckResult -Name "C6. csAgent can resolve disputes" -Passed $agentResolve.Ok -Detail "status=$($agentResolve.Code)"

$disputeCaseB = New-DisputeCase -UserToken $user.Token -AdminToken $admin.Token -PlayerID $targetPlayer.PlayerID -GameID $targetGame.GameID -TitlePrefix "cs-leader-assign-resolve"
Add-CheckResult -Name "C7. create dispute case B" -Passed $true -Detail "orderId=$($disputeCaseB.OrderID), disputeId=$($disputeCaseB.DisputeID)"

$leaderAssign = Invoke-Api -Method "POST" -Path "/admin/disputes/$($disputeCaseB.DisputeID)/assign" -Token $csLeader.Token -Body @{
  assignedServiceId = $csAgent.UserID
}
Add-CheckResult -Name "C8. csLeader can assign disputes" -Passed $leaderAssign.Ok -Detail "status=$($leaderAssign.Code)"

$leaderResolve = Invoke-Api -Method "POST" -Path "/admin/disputes/$($disputeCaseB.DisputeID)/resolve" -Token $csLeader.Token -Body @{
  resolution    = "reject"
  resolveRemark = "csLeader regression resolve"
}
Add-CheckResult -Name "C9. csLeader can resolve disputes" -Passed $leaderResolve.Ok -Detail "status=$($leaderResolve.Code)"

Write-Host ""
Write-Host "[cs-perm-regression] 回归结果:"
$results | Format-Table -AutoSize

if ($failed) {
  Write-Error "cs permission regression failed"
  exit 1
}

Write-Host "[cs-perm-regression] 全部检查通过"
