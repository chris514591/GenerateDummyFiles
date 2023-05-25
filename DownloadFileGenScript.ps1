$url = "https://github.com/chris514591/GenerateDummyFiles/archive/master.zip"
$outputPath = "C:\Users\xevic\Downloads"
# $accessToken = "YOUR_ACCESS_TOKEN"
# Use accessToken above when repository is private

# Create the output folder if it doesn't exist
if (-not (Test-Path -Path $outputPath)) {
    New-Item -ItemType Directory -Path $outputPath | Out-Null
}

# Download the repository zip file with authentication headers
$zipFilePath = Join-Path -Path $outputPath -ChildPath "GenerateDummyFiles.zip"
Invoke-WebRequest -Uri $url -OutFile $zipFilePath # Remove comment for private repo -Headers @{ "Authorization" = "Bearer $accessToken" }

# Extract the contents of the zip file
$destinationPath = Join-Path -Path $outputPath -ChildPath "GenerateDummyFiles"
Expand-Archive -Path $zipFilePath -DestinationPath $destinationPath -Force

# Clean up the zip file
Remove-Item -Path $zipFilePath

Write-Host "Download and extraction complete."
Write-Host "Repository contents are located at: $destinationPath"
