#!/usr/bin/env pwsh

<#
.SYNOPSIS
    Creates and pushes a new release tag based on current date
.DESCRIPTION
    Generates version tag in format YY.M.D or YY.M.D.B where B is build number
    if tag already exists. Automatically increments build number until unique tag is found.
.EXAMPLE
    .\release.ps1
#>

param(
    [switch]$DryRun,
    [switch]$Help
)

if ($Help) {
    Get-Help $MyInvocation.MyCommand.Path -Detailed
    exit 0
}

# Get current date
$now = Get-Date
$year = $now.ToString("yy")
$month = $now.Month
$day = $now.Day

# Base version without leading zeros
$baseVersion = "$year.$month.$day"

Write-Host "Determining release version..." -ForegroundColor Cyan
Write-Host "Base version: $baseVersion" -ForegroundColor Gray

# Fetch all tags from remote
Write-Host "Fetching tags from remote..." -ForegroundColor Gray
git fetch --tags 2>$null

# Get all existing tags (local and remote)
$allTags = git tag -l | ForEach-Object { $_.Trim() }

# Function to check if tag exists
function Test-TagExists {
    param([string]$tag)
    return $allTags -contains $tag
}

# Determine final version
$version = $baseVersion
$buildNumber = 0

if (Test-TagExists $version) {
    Write-Host "Tag '$version' already exists, checking build numbers..." -ForegroundColor Yellow

    # Find next available build number
    $buildNumber = 1
    while ($true) {
        $version = "$baseVersion.$buildNumber"
        if (-not (Test-TagExists $version)) {
            break
        }
        $buildNumber++
    }

    Write-Host "Found available version: $version (build #$buildNumber)" -ForegroundColor Green
}
else {
    Write-Host "Version '$version' is available" -ForegroundColor Green
}

# Final version tag with 'v' prefix
$tag = "v$version"

Write-Host ""
Write-Host "==================================" -ForegroundColor Cyan
Write-Host "Release version: $tag" -ForegroundColor Green
Write-Host "==================================" -ForegroundColor Cyan
Write-Host ""

if ($DryRun) {
    Write-Host "DRY RUN: Would create and push tag '$tag'" -ForegroundColor Yellow
    exit 0
}

# Confirm with user
$confirmation = Read-Host "Create and push tag '$tag'? (y/N)"
if ($confirmation -ne 'y' -and $confirmation -ne 'Y') {
    Write-Host "Release cancelled" -ForegroundColor Yellow
    exit 0
}

# Check if there are uncommitted changes
$status = git status --porcelain
if ($status) {
    Write-Host ""
    Write-Host "WARNING: You have uncommitted changes:" -ForegroundColor Yellow
    Write-Host $status
    Write-Host ""
    $proceed = Read-Host "Continue anyway? (y/N)"
    if ($proceed -ne 'y' -and $proceed -ne 'Y') {
        Write-Host "Release cancelled" -ForegroundColor Yellow
        exit 0
    }
}

# Create tag
Write-Host ""
Write-Host "Creating tag '$tag'..." -ForegroundColor Cyan
git tag -a $tag -m "Release $tag"

if ($LASTEXITCODE -ne 0) {
    Write-Host "ERROR: Failed to create tag" -ForegroundColor Red
    exit 1
}

Write-Host "Tag created successfully" -ForegroundColor Green

# Push tag
Write-Host "Pushing tag to remote..." -ForegroundColor Cyan
git push origin $tag

if ($LASTEXITCODE -ne 0) {
    Write-Host "ERROR: Failed to push tag" -ForegroundColor Red
    Write-Host "You can manually push it later with: git push origin $tag" -ForegroundColor Yellow
    exit 1
}

Write-Host ""
Write-Host "==================================" -ForegroundColor Green
Write-Host "Release $tag created successfully!" -ForegroundColor Green
Write-Host "==================================" -ForegroundColor Green
Write-Host ""
Write-Host "GitHub Actions will now build and publish the release." -ForegroundColor Cyan
Write-Host "Check progress at: https://github.com/d-mozulyov/vox/actions" -ForegroundColor Cyan
Write-Host ""
