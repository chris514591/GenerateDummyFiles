$url = "https://github.com/chris514591/GenerateDummyFiles/archive/master.zip"
$outputPath = "C:\Users\%USERPROFILE%\Downloads"

# Create the output folder if it doesn't exist
if (-not (Test-Path -Path $outputPath)) {
    New-Item -ItemType Directory -Path $outputPath | Out-Null
}

# Download the repository zip file
$zipFilePath = Join-Path -Path $outputPath -ChildPath "GenerateDummyFiles.zip"
Invoke-WebRequest -Uri $url -OutFile $zipFilePath

# Extract the contents of the zip file
$destinationPath = Join-Path -Path $outputPath -ChildPath "GenerateDummyFiles"
Expand-Archive -Path $zipFilePath -DestinationPath $destinationPath -Force

# Clean up the zip file
Remove-Item -Path $zipFilePath

Write-Host "Download and extraction complete."
Write-Host "Repository contents are located at: $destinationPath"
