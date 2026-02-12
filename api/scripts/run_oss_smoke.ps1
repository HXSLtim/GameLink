param(
  [string]$BaseUrl = "http://127.0.0.1:8080/api/v1",
  [string]$Username = "demo.user@gamelink.com",
  [string]$Password = "User@123456",
  [string]$ExpectedHost = "",
  [string]$UploadPath = "/chat/upload/image"
)

$ErrorActionPreference = "Stop"

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

  return Invoke-RestMethod @params
}

function Resolve-ExpectedHost {
  param([string]$InputHost)

  $trimmedInput = $InputHost.Trim()
  if ($trimmedInput -ne "") {
    return $trimmedInput
  }

  $endpoint = [string]$env:OSS_ENDPOINT
  if ([string]::IsNullOrWhiteSpace($endpoint)) {
    return ""
  }

  try {
    if ($endpoint.StartsWith("http://") -or $endpoint.StartsWith("https://")) {
      return ([uri]$endpoint).Host
    }
    return ([uri]("https://" + $endpoint)).Host
  } catch {
    return ""
  }
}

Write-Host "[oss-smoke] checking backend health..."
$health = Invoke-Api -Method "GET" -Path "/healthz"
if (-not $health.success) {
  throw "health check failed"
}

Write-Host "[oss-smoke] logging in..."
$login = Invoke-Api -Method "POST" -Path "/auth/login" -Body @{
  username = $Username
  password = $Password
}

$token = [string]$login.data.token
if ([string]::IsNullOrWhiteSpace($token)) {
  throw "login token is empty"
}

$expectedHostResolved = Resolve-ExpectedHost -InputHost $ExpectedHost

Write-Host "[oss-smoke] uploading tiny image..."
$tmpDir = Join-Path $PSScriptRoot "..\tmp"
New-Item -ItemType Directory -Force -Path $tmpDir | Out-Null
$imgPath = Join-Path $tmpDir "oss-smoke.png"
$pngBase64 = "iVBORw0KGgoAAAANSUhEUgAAAAEAAAABCAQAAAC1HAwCAAAAC0lEQVR42mP8/x8AAwMCAO7ZzXcAAAAASUVORK5CYII="
[IO.File]::WriteAllBytes($imgPath, [Convert]::FromBase64String($pngBase64))

Add-Type -AssemblyName System.Net.Http
$client = New-Object System.Net.Http.HttpClient
$client.DefaultRequestHeaders.Authorization = New-Object System.Net.Http.Headers.AuthenticationHeaderValue("Bearer", $token)
$form = New-Object System.Net.Http.MultipartFormDataContent
$fileBytes = [IO.File]::ReadAllBytes($imgPath)
$fileContent = New-Object System.Net.Http.ByteArrayContent (, $fileBytes)
$fileContent.Headers.ContentType = New-Object System.Net.Http.Headers.MediaTypeHeaderValue("image/png")
$form.Add($fileContent, "image", "oss-smoke.png")

$response = $client.PostAsync("$BaseUrl$UploadPath", $form).Result
$raw = $response.Content.ReadAsStringAsync().Result
if (-not $response.IsSuccessStatusCode) {
  throw "upload failed: status=$($response.StatusCode), body=$raw"
}

$payload = $raw | ConvertFrom-Json
$uploadURL = [string]$payload.data.url

if ([string]::IsNullOrWhiteSpace($uploadURL)) {
  throw "upload url is empty"
}

if ($uploadURL.StartsWith("/uploads/")) {
  throw "upload still fallback to local path: $uploadURL"
}

try {
  $uri = [uri]$uploadURL
  if ([string]::IsNullOrWhiteSpace($uri.Host)) {
    throw "url host empty"
  }

  if ($expectedHostResolved -ne "" -and $uri.Host -ne $expectedHostResolved) {
    throw "host mismatch: actual=$($uri.Host), expected=$expectedHostResolved"
  }
} catch {
  throw "invalid upload url: $uploadURL; err=$($_.Exception.Message)"
}

if ($uploadURL -notmatch "myqcloud\.com") {
  Write-Host "[oss-smoke] warning: upload url host is not myqcloud.com"
}

Write-Host "[oss-smoke] PASS"
Write-Host "upload_url=$uploadURL"
if ($expectedHostResolved -ne "") {
  Write-Host "expected_host=$expectedHostResolved"
}
