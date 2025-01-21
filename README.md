# RMedia

## Install

```bash
wget -qO- https://raw.githubusercontent.com/renatormc/scripts/refs/heads/main/install.sh | bash -s -- rmedia

curl -sSL https://raw.githubusercontent.com/renatormc/scripts/refs/heads/main/install.sh | bash -s -- rmedia


```

```powershell
$src_folder = "$env:USERPROFILE\.renatormc\src"
New-Item -Path $src_folder -ItemType Directory -Force
Set-Location -Path $src_folder
git clone "https://github.com/renatormc/$args[0]"
Set-Location -Path $args[0]
.\install.ps1
```