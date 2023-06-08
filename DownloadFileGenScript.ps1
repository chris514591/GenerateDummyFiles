$repositoryUrl = "https://github.com/chris514591/GenerateDummyFiles"
$accessToken = "YOUR_ACCESS_TOKEN" # Set this to your access token if the repository is private
$outputPath = "C:\Users\xevic\Download\GenerateDummyFiles"
$errorLogPath = Join-Path -Path $PSScriptRoot -ChildPath "errors.log"

# Function to log errors to a file
function Write-ErrorLog {
    param([string]$ErrorMessage)
    Write-Host $ErrorMessage -ForegroundColor Red
    $ErrorMessage | Out-File -FilePath $errorLogPath -Append
}

# Create the output folder if it doesn't exist
if (-not (Test-Path -Path $outputPath)) {
    try {
        New-Item -ItemType Directory -Path $outputPath -ErrorAction Stop | Out-Null
    }
    catch {
        Write-ErrorLog "Failed to create output folder: $($_.Exception.Message)"
        exit 1
    }
}

# Download the repository zip file
$zipUrl = "$repositoryUrl/archive/master.zip"
$headers = @{ "Authorization" = "Bearer $accessToken" }
$zipFilePath = Join-Path -Path $outputPath -ChildPath "GenerateDummyFiles.zip"
try {
    Invoke-WebRequest -Uri $zipUrl -OutFile $zipFilePath -Headers $headers -ErrorAction Stop
}
catch {
    Write-ErrorLog "Failed to download repository zip file: $($_.Exception.Message)"
    exit 1
}

# Extract the contents of the zip file
$destinationPath = Join-Path -Path $outputPath -ChildPath "GenerateDummyFiles"
try {
    Expand-Archive -Path $zipFilePath -DestinationPath $destinationPath -Force -ErrorAction Stop
}
catch {
    Write-ErrorLog "Failed to extract zip file: $($_.Exception.Message)"
    exit 1
}

# Clean up the zip file
try {
    Remove-Item -Path $zipFilePath -Force -ErrorAction Stop
}
catch {
    Write-ErrorLog "Failed to delete zip file: $($_.Exception.Message)"
    exit 1
}

# Remove unnecessary files
$filesToRemove = @("fileGen.go", "go.mod", "go.sum", "DownloadFileGenScript.ps1")
$filesToRemove | ForEach-Object {
    $filePath = Join-Path -Path $destinationPath -ChildPath "GenerateDummyFiles-master\$_"
    if (Test-Path -Path $filePath) {
        try {
            Remove-Item -Path $filePath -Force -ErrorAction Stop
        }
        catch {
            Write-ErrorLog "Failed to delete file: $($_.Exception.Message)"
        }
    }
}
