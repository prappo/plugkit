# PlugKit
Plugkit is a CLI tool designed for rapid WordPress plugin development. It streamlines the process using a modern tech stack, helping you build plugins faster and more efficiently.

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

## Usage

Plugkit provides several commands to help you develop WordPress plugins:

### Create a new plugin
```bash
plugkit create my-plugin
```
This command will create a new WordPress plugin with the given name. It will prompt you for additional configuration details like plugin description, version, author information, etc.

### Run in development mode
```bash
plugkit serve my-plugin
# or
plugkit serve
```
This command runs the plugin in development mode, starting the development server.

### Check version
```bash
plugkit version
```
Shows the current version of Plugkit.

### Build plugin
```bash
plugkit build
```
Builds the plugin for production.

After creating a new plugin, you'll need to:
1. Navigate to the plugin directory: `cd my-plugin`
2. Install dependencies: `npm install`
3. Install PHP dependencies: `composer install`
4. Start development: `npm run dev`

