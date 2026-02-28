#!/usr/bin/env pwsh

<#
.SYNOPSIS
    Simple build script for Vox - local development only
.DESCRIPTION
    Builds Vox for the current platform. Cross-platform builds are handled by GitHub Actions.
.PARAMETER Clean
    Clean build artifacts before building
.PARAMETER Test
    Run tests before building
.PARAMETER Run
    Run the application after building
.PARAMETER Version
    Version string (default: git tag or 'dev')
.EXAMPLE
    .\build.ps1
    Build for current platform
.EXAMPLE
    .\build.ps1 -Test -Run
    Test, build, and run
.EXAMPLE
    .\build.ps1 -Clean
    Clean and build
#>

param(
    [switch]$Clean,
    [switch]$Test,
    [switch]$Run,
    [string]$Version = '',
    [string]$OutputDir = 'dist'
)

# Project configuration
$ProjectName = 'vox'
$CmdPath = './cmd/vox'

# Colors
$ColorReset = "`e[0m"
$ColorGreen = "`e[32m"
$ColorCyan = "`e[36m"
$ColorRed = "`e[31m"

function Write-ColorOutput {
    param([string]$Message, [string]$Color = $ColorReset)
    Write-Host "${Color}${Message}${ColorReset}"
}

function Write-Header {
    param([string]$Message)
    Write-ColorOutput "`n=== $Message ===" $ColorCyan
}

function Write-Success {
    param([string]$Message)
    Write-ColorOutput "✓ $Message" $ColorGreen
}

function Write-Error {
    param([string]$Message)
    Write-ColorOutput "✗ $Message" $ColorRed
}

# Determine version
if ([string]::IsNullOrEmpty($Version)) {
    try {
        $Version = git describe --tags --always --dirty 2>$null
        if ([string]::IsNullOrEmpty($Version)) {
            $Version = 'dev'
        }
    }
    catch {
        $Version = 'dev'
    }
}

$BuildTime = (Get-Date).ToUniversalTime().ToString('yyyy-MM-dd_HH:mm:ss')
$LdFlags = "-s -w -X main.Version=$Version -X main.BuildTime=$BuildTime"

Write-Header "Vox Build System"
Write-Host "Version: $Version"
Write-Host "Output: $OutputDir"
Write-Host ""

# Clean if requested
if ($Clean) {
    Write-Header "Cleaning build artifacts"
    if (Test-Path $OutputDir) {
        Remove-Item -Recurse -Force $OutputDir
    }
    Remove-Item -Force "$ProjectName.exe" -ErrorAction SilentlyContinue
    Remove-Item -Force $ProjectName -ErrorAction SilentlyContinue
    Write-Success "Clean complete"
}

# Run tests if requested
if ($Test) {
    Write-Header "Running tests"
    go test -v ./...
    if ($LASTEXITCODE -ne 0) {
        Write-Error "Tests failed"
        exit 1
    }
    Write-Success "Tests passed"
}

# Create output directory
New-Item -ItemType Directory -Force -Path $OutputDir | Out-Null

# Build for current platform
Write-Header "Building for current platform"

$OutputName = $ProjectName
if ($IsWindows) {
    $OutputName += '.exe'
}
$OutputPath = Join-Path $OutputDir $OutputName

go build -ldflags $LdFlags -o $OutputPath $CmdPath

if ($LASTEXITCODE -eq 0) {
    Write-Success "Build complete: $OutputPath"
}
else {
    Write-Error "Build failed"
    exit 1
}

# Run if requested
if ($Run) {
    Write-Header "Running application"
    & $OutputPath
}

Write-Host ""
Write-ColorOutput "Note: Cross-platform builds are handled by GitHub Actions." $ColorCyan
