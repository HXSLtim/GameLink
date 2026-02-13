param(
  [string]$BaseUrl = "http://127.0.0.1:8080/api/v1",
  [string]$AdminUsername = "admin@gamelink.com",
  [string]$AdminPassword = "Admin123456",
  [string]$UserUsername = "demo.user@gamelink.com",
  [string]$UserPassword = "User@123456",
  [string]$PlayerUsername = "pro.player@gamelink.com",
  [string]$PlayerPassword = "Player@123456"
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

Write-Host "[cs-flow-regression] 登录账号..."

try {
  $user = Login -Username $UserUsername -Password $UserPassword

  $health = Invoke-Api -Method "GET" -Path "/healthz"
  Add-CheckResult -Name "CS0. healthz" -Passed $health.Ok -Detail "status=$($health.Code)"

  $session = Invoke-Api -Method "GET" -Path "/user/customer-service/session" -Token $user.Token
  $sessionOK = $session.Ok -and $null -ne $session.Data.data -and [uint64]$session.Data.data.groupId -gt 0 -and [uint64]$session.Data.data.agent.userId -gt 0
  Add-CheckResult -Name "CS1. get customer service session" -Passed $sessionOK -Detail "status=$($session.Code), groupId=$([string]$session.Data.data.groupId), agentId=$([string]$session.Data.data.agent.userId)"
  if (-not $sessionOK) {
    throw "获取客服会话失败: $($session.Message)"
  }

  $groupId = [uint64]$session.Data.data.groupId

  $before = Invoke-Api -Method "GET" -Path "/user/customer-service/messages?page=1&pageSize=50" -Token $user.Token
  $beforeCount = 0
  if ($before.Ok) {
    $beforeCount = @($before.Data.data.messages).Count
  }
  Add-CheckResult -Name "CS2. list customer service messages" -Passed $before.Ok -Detail "status=$($before.Code), groupId=$groupId, count=$beforeCount"
  if (-not $before.Ok) {
    throw "获取客服消息失败: $($before.Message)"
  }

  $content = "customer-service-regression-$(Get-Date -Format yyyyMMddHHmmssfff)"
  $send = Invoke-Api -Method "POST" -Path "/user/customer-service/messages" -Token $user.Token -Body @{
    content = $content
  }
  $sendOK = $send.Ok -and $null -ne $send.Data.data -and [uint64]$send.Data.data.groupId -eq $groupId
  Add-CheckResult -Name "CS3. send customer service message" -Passed $sendOK -Detail "status=$($send.Code), groupId=$([string]$send.Data.data.groupId), messageId=$([string]$send.Data.data.id)"
  if (-not $sendOK) {
    throw "发送客服消息失败: $($send.Message)"
  }

  $after = Invoke-Api -Method "GET" -Path "/user/customer-service/messages?page=1&pageSize=50" -Token $user.Token
  $matched = $false
  $afterCount = 0
  if ($after.Ok) {
    $items = @($after.Data.data.messages)
    $afterCount = $items.Count
    $matched = $items | Where-Object { [string]$_.content -eq $content } | Select-Object -First 1
    $matched = $null -ne $matched
  }
  Add-CheckResult -Name "CS4. message appears in customer service history" -Passed ($after.Ok -and $matched) -Detail "status=$($after.Code), before=$beforeCount, after=$afterCount, matched=$matched"
} catch {
  Add-CheckResult -Name "CSX. fatal error" -Passed $false -Detail $_.Exception.Message
}

Write-Host ""
Write-Host "[cs-flow-regression] 回归结果:"
$results | Format-Table -AutoSize

if ($failed) {
  Write-Error "customer service flow regression failed"
  exit 1
}

Write-Host "[cs-flow-regression] 全部检查通过"
