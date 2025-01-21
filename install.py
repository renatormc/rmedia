from pathlib import Path
import os

dest = Path.home() / ".local/bin"
try:
    dest.mkdir(parents=True)
except FileExistsError:
    pass
filename = "rmedia.exe" if os.name == "nt" else "rmedia"
path = dest / filename
os.system(f"go build -o \"{path}\"")
if not os.name == "nt":
    os.system(f"chmod +x \"{path}\"")
print(f"Generated file: {path}")