#!/bin/bash

# Create a temporary directory
TEMP_DIR=$(mktemp -d)
cd "$TEMP_DIR"

# Create bin directory in user's home if it doesn't exist
mkdir -p ~/bin

# Download the plugkit zip file
echo "Downloading plugkit..."
curl -L -o plugkit.zip https://github.com/prappo/plugkit/releases/download/latest/plugkit_darwin_amd64_v1.zip

# Unzip the file
echo "Extracting plugkit..."
unzip plugkit.zip

# Move the binary to user's bin directory
echo "Installing plugkit..."
mv plugkit_darwin_amd64_v1/plugkit ~/bin/

# Make the binary executable
chmod +x ~/bin/plugkit

# Remove quarantine attribute
echo "Removing quarantine attribute..."
xattr -d com.apple.quarantine ~/bin/plugkit 2>/dev/null || true

# Clean up
echo "Cleaning up..."
cd - > /dev/null
rm -rf "$TEMP_DIR"
rm -f plugkit.zip

# Add to PATH if not already present
if [[ ":$PATH:" != *":$HOME/bin:"* ]]; then
    # Add to both shell configs
    echo 'export PATH="$HOME/bin:$PATH"' >> ~/.zshrc
    echo 'export PATH="$HOME/bin:$PATH"' >> ~/.bash_profile
    
    # Add to current session
    export PATH="$HOME/bin:$PATH"
fi

# Verify installation
echo "Verifying installation..."
if command -v plugkit &> /dev/null; then
    echo "✅ plugkit is now available in your terminal!"
    echo "You can verify by running: plugkit --version"
else
    echo "⚠️  Installation complete, but you may need to restart your terminal."
    echo "Please run: source ~/.zshrc (or source ~/.bash_profile)"
fi
