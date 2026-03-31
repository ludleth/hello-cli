$ErrorActionPreference = "Stop"

$Repo = "ludleth/hello-cli"
$ApiUrl = "https://api.github.com/repos/$Repo/releases/latest"
$Release = Invoke-RestMethod -Uri $ApiUrl

$Asset = $Release.assets | Where-Object { $_.name -match "Windows_x86_64.zip$" }
if (-not $Asset) {
    Write-Error "Could not find a Windows x86_64 release."
    exit 1
}

$DownloadUrl = $Asset.browser_download_url
$TempZip = Join-Path $env:TEMP "hello-cli.zip"
$ExtractDir = Join-Path $env:TEMP "hello-cli-extract"

Write-Host "Downloading $DownloadUrl..."
Invoke-WebRequest -Uri $DownloadUrl -OutFile $TempZip

Write-Host "Extracting..."
if (Test-Path $ExtractDir) { Remove-Item -Recurse -Force $ExtractDir }
New-Item -ItemType Directory -Force -Path $ExtractDir | Out-Null
Expand-Archive -Path $TempZip -DestinationPath $ExtractDir -Force

$InstallDir = Join-Path $env:LOCALAPPDATA "hello-cli"
if (-not (Test-Path $InstallDir)) {
    New-Item -ItemType Directory -Force -Path $InstallDir | Out-Null
}

$ExePath = Join-Path $ExtractDir "hello-cli.exe"
Move-Item -Path $ExePath -Destination (Join-Path $InstallDir "hello-cli.exe") -Force

$UserPath = [Environment]::GetEnvironmentVariable("PATH", "User")
if ($UserPath -notmatch [regex]::Escape($InstallDir)) {
    Write-Host "Adding $InstallDir to PATH..."
    [Environment]::SetEnvironmentVariable("PATH", "$UserPath;$InstallDir", "User")
    $env:PATH = "$env:PATH;$InstallDir"
}

Remove-Item -Path $TempZip -Force
Remove-Item -Recurse -Force $ExtractDir

Write-Host "Installation complete! You may need to restart your terminal to use 'hello-cli'."
