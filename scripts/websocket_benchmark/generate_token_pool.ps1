param(
    [string]$BaseUrl = "http://localhost:8080/api/v1",
    [int]$Count = 100,
    [int]$StartIndex = 1,
    [string]$Password = "123456",
    [string]$EmailPrefix = "bench",
    [string]$EmailDomain = "example.com",
    [string]$OutFile = "scripts/benchmark/tokens.txt",
    [string]$DetailFile = "scripts/benchmark/token_pool.jsonl",
    [int]$ThrottleMs = 20,
    [switch]$DryRun
)

Set-StrictMode -Version Latest
$ErrorActionPreference = "Stop"

if ($Count -le 0) { throw "Count must be > 0" }
if ([string]::IsNullOrWhiteSpace($BaseUrl)) { throw "BaseUrl is required" }

$root = Resolve-Path (Join-Path $PSScriptRoot "..\..")
$tokenPath = Join-Path $root $OutFile
$detailPath = Join-Path $root $DetailFile

$tokenDir = Split-Path -Parent $tokenPath
$detailDir = Split-Path -Parent $detailPath
if (-not (Test-Path $tokenDir)) { New-Item -ItemType Directory -Path $tokenDir | Out-Null }
if (-not (Test-Path $detailDir)) { New-Item -ItemType Directory -Path $detailDir | Out-Null }

Set-Content -Path $tokenPath -Value "" -NoNewline
Set-Content -Path $detailPath -Value "" -NoNewline

$signupUrl = "$BaseUrl/auth/sign-up"
$loginUrl = "$BaseUrl/auth/login"

$ok = 0
$fail = 0

for ($i = 0; $i -lt $Count; $i++) {
    $idx = $StartIndex + $i
    $email = "{0}.{1}.{2}@{3}" -f $EmailPrefix, (Get-Date -Format "yyyyMMddHHmmss"), $idx, $EmailDomain

    try {
        if ($DryRun) {
            $userNumber = "dryrun-$idx"
            $accessToken = [Convert]::ToBase64String([Text.Encoding]::UTF8.GetBytes("token-$idx"))
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

        Add-Content -Path $tokenPath -Value ("Bearer {0}" -f $accessToken)

        $line = [ordered]@{
            index = $idx
            email = $email
            user_number = $userNumber
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

