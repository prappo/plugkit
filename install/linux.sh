#!/bin/bash

# Create a temporary directory
TEMP_DIR=$(mktemp -d)
cd "$TEMP_DIR"

# Download the plugkit zip file
echo "Downloading plugkit..."
curl -L -o plugkit.zip https://github.com/prappo/plugkit/releases/download/v0.0.1/plugkit_linux_amd64_v1.zip

# Unzip the file
echo "Extracting plugkit..."
unzip plugkit.zip

# Move the binary to /usr/local/bin (requires sudo)
echo "Installing plugkit..."
sudo mv plugkit_linux_amd64_v1/plugkit /usr/local/bin/

# Make the binary executable
echo "Setting executable permissions..."
sudo chmod +x /usr/local/bin/plugkit

# Clean up
echo "Cleaning up..."
cd - > /dev/null
rm -rf "$TEMP_DIR"
rm -f plugkit.zip

echo "Installation complete! You can now use 'plugkit' command."
