param(
    [string]$BaseUrl = "http://localhost:8080/api/v1",
    [int]$Count = 100,
    [int]$StartIndex = 1,
    [string]$Password = "123456",
    [string]$EmailPrefix = "bench",
    [string]$EmailDomain = "example.com",
    [string]$OutFile = "scripts/websocket_benchmark/tokens.txt",
    [string]$DetailFile = "scripts/websocket_benchmark/token_pool.jsonl",
    [string]$RecipientFile = "scripts/websocket_benchmark/recipients.txt",
    [int]$ThrottleMs = 20,
    [switch]$DryRun,
    [switch]$Overwrite
)

Set-StrictMode -Version Latest
$ErrorActionPreference = "Stop"

if ($Count -le 0) { throw "Count must be > 0" }
if ([string]::IsNullOrWhiteSpace($BaseUrl)) { throw "BaseUrl is required" }

$BaseUrl = $BaseUrl.TrimEnd('/')

$root = Resolve-Path (Join-Path $PSScriptRoot "..\..")
$tokenPath = Join-Path $root $OutFile
$detailPath = Join-Path $root $DetailFile
$recipientPath = Join-Path $root $RecipientFile

$tokenDir = Split-Path -Parent $tokenPath
$detailDir = Split-Path -Parent $detailPath
$recipientDir = Split-Path -Parent $recipientPath
if (-not (Test-Path $tokenDir)) { New-Item -ItemType Directory -Path $tokenDir | Out-Null }
if (-not (Test-Path $detailDir)) { New-Item -ItemType Directory -Path $detailDir | Out-Null }
if (-not (Test-Path $recipientDir)) { New-Item -ItemType Directory -Path $recipientDir | Out-Null }

if ($Overwrite) {
    Set-Content -Path $tokenPath -Value "" -NoNewline
    Set-Content -Path $detailPath -Value "" -NoNewline
    Set-Content -Path $recipientPath -Value "" -NoNewline
}
else {
    if (-not (Test-Path $tokenPath)) { New-Item -ItemType File -Path $tokenPath | Out-Null }
    if (-not (Test-Path $detailPath)) { New-Item -ItemType File -Path $detailPath | Out-Null }
    if (-not (Test-Path $recipientPath)) { New-Item -ItemType File -Path $recipientPath | Out-Null }
}

$signupUrl = "$BaseUrl/auth/sign-up"
$loginUrl = "$BaseUrl/auth/login"

$ok = 0
$fail = 0

function ConvertTo-Base64Url {
    param([Parameter(Mandatory = $true)][string]$Text)

    $base64 = [Convert]::ToBase64String([Text.Encoding]::UTF8.GetBytes($Text))
    return $base64.TrimEnd('=').Replace('+', '-').Replace('/', '_')
}

function ConvertFrom-Base64Url {
    param([Parameter(Mandatory = $true)][string]$Text)

    $base64 = $Text.Replace('-', '+').Replace('_', '/')
    switch ($base64.Length % 4) {
        2 { $base64 += '==' }
        3 { $base64 += '=' }
        0 { }
        default { throw "invalid base64url payload" }
    }

    return [Text.Encoding]::UTF8.GetString([Convert]::FromBase64String($base64))
}

function New-FakeJwtToken {
    param([Parameter(Mandatory = $true)][string]$UserId)

    $header = @{ alg = 'none'; typ = 'JWT' } | ConvertTo-Json -Compress
    $payload = @{ UserID = $UserId; iat = [DateTimeOffset]::UtcNow.ToUnixTimeSeconds(); exp = [DateTimeOffset]::UtcNow.AddDays(1).ToUnixTimeSeconds() } | ConvertTo-Json -Compress
    return "{0}.{1}." -f (ConvertTo-Base64Url $header), (ConvertTo-Base64Url $payload)
}

function Get-JwtRecipientId {
    param([Parameter(Mandatory = $true)][string]$Token)

    $parts = $Token.Trim().Split('.')
    if ($parts.Length -lt 2) { throw "invalid jwt format" }

    $payloadJson = ConvertFrom-Base64Url $parts[1]
    $payload = $payloadJson | ConvertFrom-Json
    foreach ($name in @('UserID', 'user_id', 'sub')) {
        $prop = $payload.PSObject.Properties[$name]
        if ($null -ne $prop) {
            $value = [string]$prop.Value
            if (-not [string]::IsNullOrWhiteSpace($value)) {
                return $value.Trim()
            }
        }
    }

    throw "jwt payload user id is empty"
}

for ($i = 0; $i -lt $Count; $i++) {
    $idx = $StartIndex + $i
    $email = "{0}.{1}.{2}@{3}" -f $EmailPrefix, (Get-Date -Format "yyyyMMddHHmmss"), $idx, $EmailDomain

    try {
        if ($DryRun) {
            $userNumber = "dryrun-$idx"
            $accessToken = New-FakeJwtToken -UserId $userNumber
        }
        else {
            $signUpBody = @{ email = $email; password = $Password } | ConvertTo-Json -Compress
            $signUpResp = Invoke-RestMethod -Method Post -Uri $signupUrl -ContentType "application/json" -Body $signUpBody
            $userNumber = $signUpResp.data.user_number

            if ([string]::IsNullOrWhiteSpace($userNumber)) {
                throw "signup success but user_number is empty: $($signUpResp | ConvertTo-Json -Depth 5 -Compress)"
            }

            $session = New-Object Microsoft.PowerShell.Commands.WebRequestSession
            $loginBody = @{ number = $userNumber; password = $Password } | ConvertTo-Json -Compress
            $loginResp = Invoke-RestMethod -Method Post -Uri $loginUrl -WebSession $session -ContentType "application/json" -Body $loginBody
            $accessToken = $loginResp.data.access_token

            if ([string]::IsNullOrWhiteSpace($accessToken)) {
                throw "login success but access_token is empty: $($loginResp | ConvertTo-Json -Depth 5 -Compress)"
            }
        }

        $recipientId = Get-JwtRecipientId -Token $accessToken

        Add-Content -Path $tokenPath -Value ("Bearer {0}" -f $accessToken)
        Add-Content -Path $recipientPath -Value $recipientId

        $line = [ordered]@{
            index = $idx
            email = $email
            user_number = $userNumber
            user_id = $recipientId
            access_token = $accessToken
        } | ConvertTo-Json -Compress
        Add-Content -Path $detailPath -Value $line

        $ok++
        Write-Host ("[{0}/{1}] OK user={2}" -f ($i + 1), $Count, $userNumber)
    }
    catch {
        $fail++
        Write-Warning ("[{0}/{1}] FAIL email={2} error={3}" -f ($i + 1), $Count, $email, $_.Exception.Message)
    }

    if ($ThrottleMs -gt 0) {
        Start-Sleep -Milliseconds $ThrottleMs
    }
}

Write-Host ""
Write-Host "Done"
Write-Host ("Success: {0}" -f $ok)
Write-Host ("Failed:  {0}" -f $fail)
Write-Host ("Token file:  {0}" -f $tokenPath)
Write-Host ("Detail file: {0}" -f $detailPath)
Write-Host ("Recipient file: {0}" -f $recipientPath)

