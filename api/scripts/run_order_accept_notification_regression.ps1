param(
  [string]$BaseUrl = "http://127.0.0.1:8080/api/v1",
  [string]$AdminUsername = "admin@gamelink.com",
  [string]$AdminPassword = "Admin123456",
  [string]$UserUsername = "demo.user@gamelink.com",
  [string]$UserPassword = "User@123456",
  [string]$PlayerUsername = "pro.player@gamelink.com",
  [string]$PlayerPassword = "Player@123456",
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

function Create-PaidConfirmedOrder {
  param(
    [string]$UserToken,
    [string]$AdminToken,
    [uint64]$PlayerID,
    [uint64]$GameID,
    [string]$TitlePrefix
  )

  $title = "$TitlePrefix-$(Get-Date -Format yyyyMMddHHmmssfff)"
  $create = Invoke-Api -Method "POST" -Path "/user/orders" -Token $UserToken -Body @{
    playerId       = $PlayerID
    gameId         = $GameID
    title          = $title
    description    = "order accept notification regression"
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

  $confirm = Invoke-Api -Method "POST" -Path "/admin/orders/$orderID/confirm" -Token $AdminToken -Body @{
    note = "confirm for accept notification regression"
  }
  if (-not $confirm.Ok) {
    throw "确认订单失败(orderId=$orderID): $($confirm.Message)"
  }

  return $orderID
}

function Get-UnreadCount {
  param(
    [string]$UserToken
  )

  $resp = Invoke-Api -Method "GET" -Path "/user/notifications/unread-count" -Token $UserToken
  if (-not $resp.Ok) {
    throw "获取未读通知数失败: $($resp.Message)"
  }

  return [int64]$resp.Data.data.unread
}

function Find-OrderAcceptNotification {
  param(
    [array]$Items,
    [uint64]$OrderID
  )

  foreach ($item in $Items) {
    if ([string]$item.referenceType -ne "order") {
      continue
    }

    if ($null -eq $item.referenceId) {
      continue
    }

    $isSameOrder = $false
    try {
      $isSameOrder = ([uint64]$item.referenceId -eq $OrderID)
    } catch {
      $isSameOrder = $false
    }
    if (-not $isSameOrder) {
      continue
    }

    $title = [string]$item.title
    $message = [string]$item.message
    $isAcceptMessage = $title.Contains("接单") -or $message.Contains("已接受您的订单")
    if ($isAcceptMessage) {
      return $item
    }
  }

  return $null
}

Write-Host "[order-accept-notify-regression] 登录账号..."

try {
  $admin = Login -Username $AdminUsername -Password $AdminPassword
  $user = Login -Username $UserUsername -Password $UserPassword
  $player = Login -Username $PlayerUsername -Password $PlayerPassword

  $targetPlayer = Resolve-TargetPlayer -UserToken $user.Token -PlayerUserID $player.UserID
  $targetGame = Resolve-TargetGame
  Write-Host "[order-accept-notify-regression] 使用陪玩师: $($targetPlayer.Nickname) (playerId=$($targetPlayer.PlayerID)), 游戏: $($targetGame.GameName) (gameId=$($targetGame.GameID))"

  $health = Invoke-Api -Method "GET" -Path "/healthz"
  Add-CheckResult -Name "N0. healthz" -Passed $health.Ok -Detail "status=$($health.Code)"

  $clearBeforeCreate = Invoke-Api -Method "POST" -Path "/user/notifications/read-all" -Token $user.Token -Body @{}
  Add-CheckResult -Name "N1. clear notifications before create order" -Passed $clearBeforeCreate.Ok -Detail "status=$($clearBeforeCreate.Code), message=$($clearBeforeCreate.Message)"
  if (-not $clearBeforeCreate.Ok) {
    throw "清空通知失败: $($clearBeforeCreate.Message)"
  }

  $orderID = Create-PaidConfirmedOrder -UserToken $user.Token -AdminToken $admin.Token -PlayerID $targetPlayer.PlayerID -GameID $targetGame.GameID -TitlePrefix "notify-accept"
  Add-CheckResult -Name "N2. create paid confirmed order" -Passed $true -Detail "orderId=$orderID"

  $clearBeforeAccept = Invoke-Api -Method "POST" -Path "/user/notifications/read-all" -Token $user.Token -Body @{}
  Add-CheckResult -Name "N3. clear notifications before accept order" -Passed $clearBeforeAccept.Ok -Detail "status=$($clearBeforeAccept.Code), message=$($clearBeforeAccept.Message)"
  if (-not $clearBeforeAccept.Ok) {
    throw "接单前清空通知失败: $($clearBeforeAccept.Message)"
  }

  $unreadBefore = Get-UnreadCount -UserToken $user.Token
  Add-CheckResult -Name "N4. unread count before accept" -Passed $true -Detail "unread=$unreadBefore"

  $accept = Invoke-Api -Method "POST" -Path "/player/orders/$orderID/accept" -Token $player.Token -Body @{}
  Add-CheckResult -Name "N5. player accepts order" -Passed $accept.Ok -Detail "orderId=$orderID, status=$($accept.Code), message=$($accept.Message)"
  if (-not $accept.Ok) {
    throw "接单失败(orderId=$orderID): $($accept.Message)"
  }

  $unreadAfter = Get-UnreadCount -UserToken $user.Token
  $delta = $unreadAfter - $unreadBefore
  Add-CheckResult -Name "N6. unread count increases after accept" -Passed ($delta -gt 0) -Detail "before=$unreadBefore, after=$unreadAfter, delta=$delta"

  $list = Invoke-Api -Method "GET" -Path "/user/notifications?page=1&pageSize=50" -Token $user.Token
  $matched = $null
  if ($list.Ok) {
    $items = @($list.Data.data.items)
    $matched = Find-OrderAcceptNotification -Items $items -OrderID $orderID
  }
  $matchedOK = $list.Ok -and ($null -ne $matched)
  Add-CheckResult -Name "N7. user notification contains accepted-order event" -Passed $matchedOK -Detail "orderId=$orderID, found=$($null -ne $matched), title=$([string]$matched.title), message=$([string]$matched.message)"
} catch {
  Add-CheckResult -Name "NX. fatal error" -Passed $false -Detail $_.Exception.Message
}

Write-Host ""
Write-Host "[order-accept-notify-regression] 回归结果:"
$results | Format-Table -AutoSize

if ($failed) {
  Write-Error "order accept notification regression failed"
  exit 1
}

Write-Host "[order-accept-notify-regression] 全部检查通过"
