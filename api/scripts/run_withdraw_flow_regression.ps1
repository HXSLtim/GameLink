param(
  [string]$BaseUrl = "http://127.0.0.1:8080/api/v1",
  [string]$AdminUsername = "admin@gamelink.com",
  [string]$AdminPassword = "Admin123456",
  [string]$UserUsername = "demo.user@gamelink.com",
  [string]$UserPassword = "User@123456",
  [string]$PlayerUsername = "pro.player@gamelink.com",
  [string]$PlayerPassword = "Player@123456",
  [string]$PreferredGameName = "英雄联盟",
  [int64]$WithdrawAmountCents = 10000,
  [int]$MaxTopupOrders = 8
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

function Get-EarningsSummary {
  param(
    [string]$PlayerToken
  )

  $resp = Invoke-Api -Method "GET" -Path "/player/earnings/summary" -Token $PlayerToken
  if (-not $resp.Ok) {
    throw "获取收益概览失败: $($resp.Message)"
  }

  $data = $resp.Data.data
  return [pscustomobject]@{
    TotalEarnings    = [int64]$data.totalEarnings
    AvailableBalance = [int64]$data.availableBalance
    WithdrawTotal    = [int64]$data.withdrawTotal
    PendingBalance   = [int64]$data.pendingBalance
  }
}

function Create-PaidCompletedOrder {
  param(
    [string]$UserToken,
    [string]$AdminToken,
    [string]$PlayerToken,
    [uint64]$PlayerID,
    [uint64]$GameID,
    [string]$TitlePrefix
  )

  $title = "$TitlePrefix-$(Get-Date -Format yyyyMMddHHmmssfff)"
  $create = Invoke-Api -Method "POST" -Path "/user/orders" -Token $UserToken -Body @{
    playerId       = $PlayerID
    gameId         = $GameID
    title          = $title
    description    = "withdraw flow regression topup"
    scheduledStart = (Get-Date).AddHours(2).ToString("o")
    durationHours  = 1
  }
  if (-not $create.Ok) {
    throw "创建订单失败: $($create.Message)"
  }

  $orderID = [uint64]$create.Data.data.orderId

  $pay = Invoke-Api -Method "POST" -Path "/user/payments" -Token $UserToken -Body @{
    orderId = $orderID
    method  = "wechat"
  }
  if (-not $pay.Ok) {
    throw "支付订单失败(orderId=$orderID): $($pay.Message)"
  }

  $start = Invoke-Api -Method "POST" -Path "/admin/orders/$orderID/start" -Token $AdminToken -Body @{
    note = "withdraw regression start"
  }
  if (-not $start.Ok) {
    # 兼容部分环境：先 confirm 再 start
    $confirm = Invoke-Api -Method "POST" -Path "/admin/orders/$orderID/confirm" -Token $AdminToken -Body @{
      note = "withdraw regression confirm"
    }
    if (-not $confirm.Ok) {
      throw "开始前确认订单失败(orderId=$orderID): $($confirm.Message)"
    }

    $start = Invoke-Api -Method "POST" -Path "/admin/orders/$orderID/start" -Token $AdminToken -Body @{
      note = "withdraw regression start"
    }
    if (-not $start.Ok) {
      throw "开始订单失败(orderId=$orderID): $($start.Message)"
    }
  }

  $complete = Invoke-Api -Method "PUT" -Path "/player/orders/$orderID/complete" -Token $PlayerToken -Body @{}
  if (-not $complete.Ok) {
    throw "陪玩师完成订单失败(orderId=$orderID): $($complete.Message)"
  }

  return $orderID
}

function Get-WithdrawsByStatus {
  param(
    [string]$AdminToken,
    [uint64]$PlayerID,
    [string]$Status
  )

  $resp = Invoke-Api -Method "GET" -Path "/admin/withdraws?page=1&pageSize=100&playerId=$PlayerID&status=$Status" -Token $AdminToken
  if (-not $resp.Ok) {
    throw "获取提现列表失败(status=$Status): $($resp.Message)"
  }

  return @($resp.Data.data.withdraws)
}

function Clear-PendingWithdraws {
  param(
    [string]$AdminToken,
    [uint64]$PlayerID
  )

  $processed = 0

  $pending = Get-WithdrawsByStatus -AdminToken $AdminToken -PlayerID $PlayerID -Status "pending"
  foreach ($item in $pending) {
    $id = [uint64]$item.id
    $approve = Invoke-Api -Method "POST" -Path "/admin/withdraws/$id/approve" -Token $AdminToken -Body @{
      remark = "auto-approve for withdraw regression precheck"
    }
    if (-not $approve.Ok) {
      throw "预处理批准提现失败(withdrawId=$id): $($approve.Message)"
    }

    $complete = Invoke-Api -Method "POST" -Path "/admin/withdraws/$id/complete" -Token $AdminToken -Body @{}
    if (-not $complete.Ok) {
      throw "预处理完成提现失败(withdrawId=$id): $($complete.Message)"
    }
    $processed++
  }

  $approved = Get-WithdrawsByStatus -AdminToken $AdminToken -PlayerID $PlayerID -Status "approved"
  foreach ($item in $approved) {
    $id = [uint64]$item.id
    $complete = Invoke-Api -Method "POST" -Path "/admin/withdraws/$id/complete" -Token $AdminToken -Body @{}
    if (-not $complete.Ok) {
      throw "预处理完成已批准提现失败(withdrawId=$id): $($complete.Message)"
    }
    $processed++
  }

  return $processed
}

Write-Host "[withdraw-regression] 登录账号..."

$createdOrders = New-Object System.Collections.Generic.List[uint64]

try {
  $admin = Login -Username $AdminUsername -Password $AdminPassword
  $user = Login -Username $UserUsername -Password $UserPassword
  $player = Login -Username $PlayerUsername -Password $PlayerPassword

  $targetPlayer = Resolve-TargetPlayer -UserToken $user.Token -PlayerUserID $player.UserID
  $targetGame = Resolve-TargetGame
  Write-Host "[withdraw-regression] 使用陪玩师: $($targetPlayer.Nickname) (playerId=$($targetPlayer.PlayerID)), 游戏: $($targetGame.GameName) (gameId=$($targetGame.GameID))"

  $health = Invoke-Api -Method "GET" -Path "/healthz"
  Add-CheckResult -Name "W0. healthz" -Passed $health.Ok -Detail "status=$($health.Code)"

  $cleared = Clear-PendingWithdraws -AdminToken $admin.Token -PlayerID $targetPlayer.PlayerID
  Add-CheckResult -Name "W1. clear pending/approved withdraws" -Passed $true -Detail "processed=$cleared"

  $summaryBefore = Get-EarningsSummary -PlayerToken $player.Token
  $topupCount = 0
  while ($summaryBefore.AvailableBalance -lt $WithdrawAmountCents -and $topupCount -lt $MaxTopupOrders) {
    $orderID = Create-PaidCompletedOrder -UserToken $user.Token -AdminToken $admin.Token -PlayerToken $player.Token -PlayerID $targetPlayer.PlayerID -GameID $targetGame.GameID -TitlePrefix "withdraw-topup"
    $createdOrders.Add($orderID)
    $topupCount++
    $summaryBefore = Get-EarningsSummary -PlayerToken $player.Token
  }

  $balanceReady = $summaryBefore.AvailableBalance -ge $WithdrawAmountCents
  Add-CheckResult -Name "W2. ensure available balance" -Passed $balanceReady -Detail "available=$($summaryBefore.AvailableBalance), target=$WithdrawAmountCents, topups=$topupCount, orders=$($createdOrders -join ',')"
  if (-not $balanceReady) {
    throw "可提现余额不足，无法继续。available=$($summaryBefore.AvailableBalance), target=$WithdrawAmountCents"
  }

  $withdrawResp = Invoke-Api -Method "POST" -Path "/player/earnings/withdraw" -Token $player.Token -Body @{
    amountCents = $WithdrawAmountCents
    method      = "wechat"
    accountInfo = "wechat:withdraw-regression-$(Get-Date -Format yyyyMMddHHmmss)"
  }
  Add-CheckResult -Name "W3. player request withdraw" -Passed $withdrawResp.Ok -Detail "status=$($withdrawResp.Code), message=$($withdrawResp.Message)"
  if (-not $withdrawResp.Ok) {
    throw "申请提现失败: $($withdrawResp.Message)"
  }

  $withdrawID = [uint64]$withdrawResp.Data.data.withdrawId

  $approve = Invoke-Api -Method "POST" -Path "/admin/withdraws/$withdrawID/approve" -Token $admin.Token -Body @{
    remark = "approved by withdraw regression script"
  }
  Add-CheckResult -Name "W4. admin approve withdraw" -Passed $approve.Ok -Detail "withdrawId=$withdrawID, status=$($approve.Code), message=$($approve.Message)"
  if (-not $approve.Ok) {
    throw "批准提现失败(withdrawId=$withdrawID): $($approve.Message)"
  }

  $complete = Invoke-Api -Method "POST" -Path "/admin/withdraws/$withdrawID/complete" -Token $admin.Token -Body @{}
  Add-CheckResult -Name "W5. admin complete withdraw" -Passed $complete.Ok -Detail "withdrawId=$withdrawID, status=$($complete.Code), message=$($complete.Message)"
  if (-not $complete.Ok) {
    throw "完成提现失败(withdrawId=$withdrawID): $($complete.Message)"
  }

  $history = Invoke-Api -Method "GET" -Path "/player/earnings/withdraw-history?page=1&pageSize=20" -Token $player.Token
  $historyRecord = $null
  if ($history.Ok) {
    $records = @($history.Data.data.records)
    $historyRecord = $records | Where-Object { [uint64]$_.id -eq $withdrawID } | Select-Object -First 1
  }
  $historyOK = $history.Ok -and $null -ne $historyRecord -and [string]$historyRecord.status -eq "completed"
  Add-CheckResult -Name "W6. player sees completed withdraw history" -Passed $historyOK -Detail "withdrawId=$withdrawID, found=$($null -ne $historyRecord), status=$([string]$historyRecord.status)"

  $summaryAfter = Get-EarningsSummary -PlayerToken $player.Token
  $withdrawDelta = $summaryAfter.WithdrawTotal - $summaryBefore.WithdrawTotal
  $summaryOK = $withdrawDelta -ge $WithdrawAmountCents
  Add-CheckResult -Name "W7. summary withdraw total updated" -Passed $summaryOK -Detail "before=$($summaryBefore.WithdrawTotal), after=$($summaryAfter.WithdrawTotal), delta=$withdrawDelta"
} catch {
  Add-CheckResult -Name "WX. fatal error" -Passed $false -Detail $_.Exception.Message
}

Write-Host ""
Write-Host "[withdraw-regression] 回归结果:"
$results | Format-Table -AutoSize

if ($failed) {
  Write-Error "withdraw flow regression failed"
  exit 1
}

Write-Host "[withdraw-regression] 全部检查通过"
