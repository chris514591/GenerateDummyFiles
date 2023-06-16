$repositoryUrl = "https://github.com/chris514591/GenerateDummyFiles"
$accessToken = "YOUR_ACCESS_TOKEN" # Set this to your access token if the repository is private
$outputPath = "C:\Users\xevic\Downloads\GenerateDummyFiles"

# Create the output folder if it doesn't exist
if (-not (Test-Path -Path $outputPath)) {
    try {
        New-Item -ItemType Directory -Path $outputPath -ErrorAction Stop | Out-Null
    }
    catch {
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
    exit 1
}

# Extract the contents of the zip file
$destinationPath = Join-Path -Path $outputPath -ChildPath "GenerateDummyFiles"
try {
    Expand-Archive -Path $zipFilePath -DestinationPath $destinationPath -Force -ErrorAction Stop
}
catch {
    exit 1
}

# Clean up the zip file
try {
    Remove-Item -Path $zipFilePath -Force -ErrorAction Stop
}
catch {
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
        }
    }
}
