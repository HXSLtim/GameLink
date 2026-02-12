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
      Ok         = $true
      StatusCode = 200
      Data       = $resp
      Message    = ""
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
      Ok         = $false
      StatusCode = $code
      Data       = $parsed
      Raw        = $raw
      Message    = $message
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
    throw "陪玩师列表为空，无法继续联调"
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
    throw "游戏列表为空，无法继续联调"
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

function Create-Order {
  param(
    [string]$UserToken,
    [uint64]$PlayerID,
    [uint64]$GameID,
    [string]$TitlePrefix
  )

  $scheduledStart = (Get-Date).AddHours(2).ToString("o")
  $title = "$TitlePrefix-$(Get-Date -Format yyyyMMddHHmmssfff)"

  $resp = Invoke-Api -Method "POST" -Path "/user/orders" -Token $UserToken -Body @{
    playerId       = $PlayerID
    gameId         = $GameID
    title          = $title
    description    = "regression check"
    scheduledStart = $scheduledStart
    durationHours  = 1
  }

  if (-not $resp.Ok) {
    throw "创建订单失败: $($resp.Message)"
  }

  return [uint64]$resp.Data.data.orderId
}

function Assert-ApiErrorContains {
  param(
    [string]$Name,
    $Resp,
    [int]$ExpectedCode,
    [string]$ExpectedText
  )

  $actualText = ""
  if ($Resp.Data -and $Resp.Data.message) {
    $actualText = [string]$Resp.Data.message
  } else {
    $actualText = [string]$Resp.Message
  }

  $ok = (-not $Resp.Ok) -and ($Resp.StatusCode -eq $ExpectedCode) -and $actualText.Contains($ExpectedText)
  Add-CheckResult -Name $Name -Passed $ok -Detail "status=$($Resp.StatusCode), message=$actualText"
}

Write-Host "[flow-regression] 登录账号..."
$admin = Login -Username $AdminUsername -Password $AdminPassword
$user = Login -Username $UserUsername -Password $UserPassword
$player = Login -Username $PlayerUsername -Password $PlayerPassword

$targetPlayer = Resolve-TargetPlayer -UserToken $user.Token -PlayerUserID $player.UserID
$targetGame = Resolve-TargetGame

Write-Host "[flow-regression] 使用陪玩师: $($targetPlayer.Nickname) (playerId=$($targetPlayer.PlayerID)), 游戏: $($targetGame.GameName) (gameId=$($targetGame.GameID))"

# Case A: 未支付完成拦截 + 评价前置拦截
Write-Host "[flow-regression] Case A: 未支付订单完成拦截"
$orderA = Create-Order -UserToken $user.Token -PlayerID $targetPlayer.PlayerID -GameID $targetGame.GameID -TitlePrefix "guard-unpaid"

$confirmA = Invoke-Api -Method "POST" -Path "/admin/orders/$orderA/confirm" -Token $admin.Token -Body @{ note = "regression unpaid confirm" }
Add-CheckResult -Name "A1. admin confirm unpaid order" -Passed $confirmA.Ok -Detail "orderId=$orderA"

$startA = Invoke-Api -Method "POST" -Path "/admin/orders/$orderA/start" -Token $admin.Token -Body @{ note = "regression unpaid start" }
Add-CheckResult -Name "A2. admin start unpaid order" -Passed $startA.Ok -Detail "orderId=$orderA"

$playerCompleteA = Invoke-Api -Method "PUT" -Path "/player/orders/$orderA/complete" -Token $player.Token -Body @{}
Assert-ApiErrorContains -Name "A3. player complete blocked before payment" -Resp $playerCompleteA -ExpectedCode 400 -ExpectedText "order must be paid before completion"

$userCompleteA = Invoke-Api -Method "PUT" -Path "/user/orders/$orderA/complete" -Token $user.Token -Body @{}
Assert-ApiErrorContains -Name "A4. user complete blocked before payment" -Resp $userCompleteA -ExpectedCode 400 -ExpectedText "order must be paid before completion"

$reviewA = Invoke-Api -Method "POST" -Path "/user/reviews" -Token $user.Token -Body @{
  orderId = $orderA
  rating  = 5
  comment = "should fail before completion"
}
Assert-ApiErrorContains -Name "A5. review blocked before completed" -Resp $reviewA -ExpectedCode 400 -ExpectedText "订单未完成"

# Case B: 支付后完成 + 评价成功
# 使用 wechat mock，避免钱包余额波动导致回归脚本不稳定
Write-Host "[flow-regression] Case B: 支付后完成与评价"
$orderB = Create-Order -UserToken $user.Token -PlayerID $targetPlayer.PlayerID -GameID $targetGame.GameID -TitlePrefix "guard-paid"

$payB = Invoke-Api -Method "POST" -Path "/user/payments" -Token $user.Token -Body @{
  orderId = $orderB
  method  = "wechat"
}
Add-CheckResult -Name "B1. create payment (wechat mock)" -Passed $payB.Ok -Detail "orderId=$orderB, status=$($payB.StatusCode), message=$($payB.Message)"

$startB = Invoke-Api -Method "POST" -Path "/admin/orders/$orderB/start" -Token $admin.Token -Body @{ note = "regression paid start" }
Add-CheckResult -Name "B2. admin start paid order" -Passed $startB.Ok -Detail "orderId=$orderB, status=$($startB.StatusCode), message=$($startB.Message)"

$playerCompleteB = Invoke-Api -Method "PUT" -Path "/player/orders/$orderB/complete" -Token $player.Token -Body @{}
Add-CheckResult -Name "B3. player complete succeeds after payment" -Passed $playerCompleteB.Ok -Detail "orderId=$orderB, status=$($playerCompleteB.StatusCode), message=$($playerCompleteB.Message)"

$reviewB = Invoke-Api -Method "POST" -Path "/user/reviews" -Token $user.Token -Body @{
  orderId = $orderB
  rating  = 5
  comment = "regression success"
}
Add-CheckResult -Name "B4. review succeeds after completed" -Passed $reviewB.Ok -Detail "orderId=$orderB, status=$($reviewB.StatusCode), message=$($reviewB.Message)"

Write-Host ""
Write-Host "[flow-regression] 回归结果:"
$results | Format-Table -AutoSize

if ($failed) {
  Write-Error "flow guard regression failed"
  exit 1
}

Write-Host "[flow-regression] 全部检查通过"
