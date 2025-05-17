# Save the original working directory
$originalLocation = Get-Location

# Create a temporary directory
$tempDir = Join-Path $env:TEMP "plugkit_install"
New-Item -ItemType Directory -Force -Path $tempDir | Out-Null
Set-Location $tempDir

Write-Host "Downloading plugkit..."
$url = "https://github.com/prappo/plugkit/releases/download/latest/plugkit_windows_amd64_v1.zip"
$output = "plugkit.zip"
try {
    Invoke-WebRequest -Uri $url -OutFile $output -ErrorAction Stop
} catch {
    Write-Host "Failed to download plugkit: $_" -ForegroundColor Red
    exit 1
}

Write-Host "Extracting plugkit..."
try {
    Expand-Archive -Path $output -DestinationPath $tempDir -Force -ErrorAction Stop
    Write-Host "Extracted files:"
    Get-ChildItem -Path $tempDir -Recurse | ForEach-Object { Write-Host $_.FullName }
} catch {
    Write-Host "Failed to extract plugkit: $_" -ForegroundColor Red
    exit 1
}

# Create a directory in user's local app data
$installDir = Join-Path $env:LOCALAPPDATA "plugkit"
try {
    if (-not (Test-Path $installDir)) {
        New-Item -ItemType Directory -Force -Path $installDir -ErrorAction Stop | Out-Null
    }
} catch {
    Write-Host "Failed to create installation directory: $_" -ForegroundColor Red
    exit 1
}

Write-Host "Installing plugkit..."
try {
    # Look for the executable in the extracted files
    $exeFile = Get-ChildItem -Path $tempDir -Filter "plugkit.exe" -Recurse -ErrorAction Stop
    if (-not $exeFile) {
        throw "plugkit.exe not found in extracted files"
    }
    Copy-Item -Path $exeFile.FullName -Destination $installDir -Force -ErrorAction Stop
} catch {
    Write-Host "Failed to install plugkit: $_" -ForegroundColor Red
    exit 1
}

# Add to user's PATH if not already present
$currentPath = [Environment]::GetEnvironmentVariable("Path", "User")
if (-not $currentPath.Contains($installDir)) {
    Write-Host "Adding plugkit to user PATH..."
    try {
        [Environment]::SetEnvironmentVariable("Path", $currentPath + ";$installDir", "User")
        $env:Path = [System.Environment]::GetEnvironmentVariable("Path","User")
    } catch {
        Write-Host "Failed to update user PATH: $_" -ForegroundColor Red
        Write-Host "You may need to manually add $installDir to your PATH" -ForegroundColor Yellow
    }
}

# Clean up
Write-Host "Cleaning up..."
try {
    Set-Location $originalLocation
    Remove-Item -Path $tempDir -Recurse -Force -ErrorAction Stop
    Remove-Item -Path $output -Force -ErrorAction SilentlyContinue
} catch {
    Write-Host "Warning: Failed to clean up temporary files: $_" -ForegroundColor Yellow
}

Write-Host "Installation complete! You can now use 'plugkit' command."
Write-Host "Note: You may need to restart your terminal for the PATH changes to take effect." 