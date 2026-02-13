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
    description    = "partial refund regression"
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
    note = "confirm for partial refund regression"
  }
  if (-not $confirm.Ok) {
    throw "确认订单失败(orderId=$orderID): $($confirm.Message)"
  }

  return $orderID
}

function Get-OrderPrimaryPayment {
  param(
    [string]$AdminToken,
    [uint64]$OrderID
  )

  $paymentsResp = Invoke-Api -Method "GET" -Path "/admin/orders/$OrderID/payments" -Token $AdminToken
  if (-not $paymentsResp.Ok) {
    throw "获取订单支付记录失败(orderId=$OrderID): $($paymentsResp.Message)"
  }

  $payments = @($paymentsResp.Data.data)
  if ($payments.Count -eq 0) {
    throw "订单支付记录为空(orderId=$OrderID)"
  }

  $target = $payments |
    Where-Object { [string]$_.status -eq "paid" -or [string]$_.status -eq "refunded" } |
    Select-Object -First 1
  if ($null -eq $target) {
    $target = $payments | Select-Object -First 1
  }

  $amount = [int64]$target.amountCents
  $refunded = [int64]$target.refundedAmountCents
  if ($refunded -lt 0) {
    $refunded = 0
  }
  if ($refunded -gt $amount) {
    $refunded = $amount
  }
  $refundable = $amount - $refunded

  return [pscustomobject]@{
    PaymentID   = [uint64]$target.id
    AmountCents = $amount
    Refunded    = $refunded
    Refundable  = $refundable
  }
}

function Sum-RefundAmounts {
  param(
    [array]$Items
  )

  $sum = [int64]0
  foreach ($item in $Items) {
    $value = $null
    if ($item.PSObject.Properties.Name -contains "amount_cents") {
      $value = $item.amount_cents
    } elseif ($item.PSObject.Properties.Name -contains "amountCents") {
      $value = $item.amountCents
    }

    if ($null -eq $value) {
      continue
    }
    $sum += [int64]$value
  }
  return $sum
}

Write-Host "[partial-refund-regression] 登录账号..."

try {
  $admin = Login -Username $AdminUsername -Password $AdminPassword
  $user = Login -Username $UserUsername -Password $UserPassword
  $player = Login -Username $PlayerUsername -Password $PlayerPassword

  $targetPlayer = Resolve-TargetPlayer -UserToken $user.Token -PlayerUserID $player.UserID
  $targetGame = Resolve-TargetGame
  Write-Host "[partial-refund-regression] 使用陪玩师: $($targetPlayer.Nickname) (playerId=$($targetPlayer.PlayerID)), 游戏: $($targetGame.GameName) (gameId=$($targetGame.GameID))"

  $health = Invoke-Api -Method "GET" -Path "/healthz"
  Add-CheckResult -Name "R0. healthz" -Passed $health.Ok -Detail "status=$($health.Code)"

  $orderID = Create-PaidConfirmedOrder -UserToken $user.Token -AdminToken $admin.Token -PlayerID $targetPlayer.PlayerID -GameID $targetGame.GameID -TitlePrefix "partial-refund"
  Add-CheckResult -Name "R1. create paid confirmed order" -Passed $true -Detail "orderId=$orderID"

  $payment = Get-OrderPrimaryPayment -AdminToken $admin.Token -OrderID $orderID
  $refundableBefore = $payment.Refundable
  $canPartial = $refundableBefore -ge 2
  Add-CheckResult -Name "R2. refundable amount is enough for partial flow" -Passed $canPartial -Detail "paymentId=$($payment.PaymentID), amount=$($payment.AmountCents), refunded=$($payment.Refunded), refundable=$refundableBefore"
  if (-not $canPartial) {
    throw "可退款金额不足以验证部分退款: refundable=$refundableBefore"
  }

  $firstRefund = [int64][Math]::Floor($refundableBefore / 2)
  if ($firstRefund -le 0) {
    $firstRefund = 1
  }
  if ($firstRefund -ge $refundableBefore) {
    $firstRefund = $refundableBefore - 1
  }
  $secondRefund = $refundableBefore - $firstRefund

  $first = Invoke-Api -Method "POST" -Path "/admin/orders/$orderID/refund" -Token $admin.Token -Body @{
    reason       = "partial refund regression #1"
    amount_cents = $firstRefund
    note         = "first partial refund"
  }
  Add-CheckResult -Name "R3. first partial refund succeeds" -Passed $first.Ok -Detail "orderId=$orderID, amount=$firstRefund, status=$($first.Code), message=$($first.Message)"
  if (-not $first.Ok) {
    throw "第一次部分退款失败(orderId=$orderID): $($first.Message)"
  }

  $orderAfterFirst = Invoke-Api -Method "GET" -Path "/admin/orders/$orderID" -Token $admin.Token
  $statusAfterFirst = ""
  if ($orderAfterFirst.Ok) {
    $statusAfterFirst = [string]$orderAfterFirst.Data.data.status
  }
  $firstStatusOK = $orderAfterFirst.Ok -and ($statusAfterFirst -ne "refunded")
  Add-CheckResult -Name "R4. order not fully refunded after first refund" -Passed $firstStatusOK -Detail "status=$statusAfterFirst"

  $refundsAfterFirstResp = Invoke-Api -Method "GET" -Path "/admin/orders/$orderID/refunds" -Token $admin.Token
  $refundsAfterFirst = @()
  if ($refundsAfterFirstResp.Ok) {
    $refundsAfterFirst = @($refundsAfterFirstResp.Data.data)
  }
  $sumAfterFirst = Sum-RefundAmounts -Items $refundsAfterFirst
  $firstRefundListed = $refundsAfterFirstResp.Ok -and ($sumAfterFirst -ge $firstRefund)
  Add-CheckResult -Name "R5. refund records include first partial refund" -Passed $firstRefundListed -Detail "records=$($refundsAfterFirst.Count), sum=$sumAfterFirst, expectedAtLeast=$firstRefund"

  $second = Invoke-Api -Method "POST" -Path "/admin/orders/$orderID/refund" -Token $admin.Token -Body @{
    reason       = "partial refund regression #2"
    amount_cents = $secondRefund
    note         = "second refund to full"
  }
  Add-CheckResult -Name "R6. second refund succeeds" -Passed $second.Ok -Detail "orderId=$orderID, amount=$secondRefund, status=$($second.Code), message=$($second.Message)"
  if (-not $second.Ok) {
    throw "第二次退款失败(orderId=$orderID): $($second.Message)"
  }

  $refundsAfterSecondResp = Invoke-Api -Method "GET" -Path "/admin/orders/$orderID/refunds" -Token $admin.Token
  $refundsAfterSecond = @()
  if ($refundsAfterSecondResp.Ok) {
    $refundsAfterSecond = @($refundsAfterSecondResp.Data.data)
  }
  $sumAfterSecond = Sum-RefundAmounts -Items $refundsAfterSecond
  $secondRefundListed = $refundsAfterSecondResp.Ok -and ($sumAfterSecond -ge $refundableBefore)
  Add-CheckResult -Name "R7. refund records reach full refunded amount" -Passed $secondRefundListed -Detail "records=$($refundsAfterSecond.Count), sum=$sumAfterSecond, expectedAtLeast=$refundableBefore"

  $orderAfterSecond = Invoke-Api -Method "GET" -Path "/admin/orders/$orderID" -Token $admin.Token
  $statusAfterSecond = ""
  if ($orderAfterSecond.Ok) {
    $statusAfterSecond = [string]$orderAfterSecond.Data.data.status
  }
  $secondStatusOK = $orderAfterSecond.Ok -and ($statusAfterSecond -eq "refunded")
  Add-CheckResult -Name "R8. order status becomes refunded after full amount refunded" -Passed $secondStatusOK -Detail "status=$statusAfterSecond"

  $overRefund = Invoke-Api -Method "POST" -Path "/admin/orders/$orderID/refund" -Token $admin.Token -Body @{
    reason       = "over refund regression"
    amount_cents = 1
    note         = "should be blocked"
  }
  $overRefundBlocked = (-not $overRefund.Ok) -and ($overRefund.Code -eq 400)
  Add-CheckResult -Name "R9. over-refund is blocked" -Passed $overRefundBlocked -Detail "status=$($overRefund.Code), message=$($overRefund.Message)"
} catch {
  Add-CheckResult -Name "RX. fatal error" -Passed $false -Detail $_.Exception.Message
}

Write-Host ""
Write-Host "[partial-refund-regression] 回归结果:"
$results | Format-Table -AutoSize

if ($failed) {
  Write-Error "partial refund regression failed"
  exit 1
}

Write-Host "[partial-refund-regression] 全部检查通过"
