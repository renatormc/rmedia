$dest = "$env:USERPROFILE\.local\bin\rmedia.exe"
go build -o $dest
Write-Host "Installed at $dest"