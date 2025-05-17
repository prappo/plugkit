# PlugKit ( WIP )

## Installation

### Mac
```bash
/bin/bash -c "$(curl -fsSL dub.sh/plugkit/mac)"
```

### Linux
```bash
/bin/bash -c "$(curl -fsSL dub.sh/plugkit/linux)"
```

### Windows
```bash
Set-ExecutionPolicy Bypass -Scope Process -Force; [System.Net.ServicePointManager]::SecurityProtocol = [System.Net.ServicePointManager]::SecurityProtocol -bor 3072; iex ((New-Object System.Net.WebClient).DownloadString('https://dub.sh/plugkit/windows'))
```