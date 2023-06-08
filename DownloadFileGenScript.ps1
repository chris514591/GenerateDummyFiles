$repositoryUrl = "https://github.com/chris514591/GenerateDummyFiles"
# $accessToken = "YOUR_ACCESS_TOKEN" # Set this to your access token if the repository is private
$outputPath = "C:\Users\xevic\Downloads\GenerateDummyFiles"

# Create the output folder if it doesn't exist
if (-not (Test-Path -Path $outputPath)) {
    New-Item -ItemType Directory -Path $outputPath | Out-Null
}

# Download the repository zip file
$zipUrl = "$repositoryUrl/archive/master.zip"
# $headers = @{ "Authorization" = "Bearer $accessToken" }
$zipFilePath = Join-Path -Path $outputPath -ChildPath "GenerateDummyFiles.zip"
Invoke-WebRequest -Uri $zipUrl -OutFile $zipFilePath -Headers $headers

# Extract the contents of the zip file
$destinationPath = Join-Path -Path $outputPath -ChildPath "GenerateDummyFiles"
Expand-Archive -Path $zipFilePath -DestinationPath $destinationPath -Force

# Clean up the zip file
Remove-Item -Path $zipFilePath -Force

# Remove unnecessary files
$filesToRemove = @("fileGen.go", "go.mod", "go.sum", "DownloadFileGenScript.ps1")
$filesToRemove | ForEach-Object {
    $filePath = Join-Path -Path $destinationPath -ChildPath "GenerateDummyFiles-master\$_"
    if (Test-Path -Path $filePath) {
        Remove-Item -Path $filePath -Force
    }
}
