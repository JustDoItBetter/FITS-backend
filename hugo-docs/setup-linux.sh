#!/bin/bash
# FITS Backend Documentation Setup Script for Linux

set -e

echo "🐧 FITS Backend Documentation Setup (Linux)"
echo "=========================================="
echo ""

# Check if Hugo is installed
if ! command -v hugo &> /dev/null; then
    echo "📦 Hugo not found. Installing Hugo..."
    
    # Download Hugo Extended
    HUGO_VERSION="0.120.0"
    wget -q https://github.com/gohugoio/hugo/releases/download/v${HUGO_VERSION}/hugo_extended_${HUGO_VERSION}_linux-amd64.tar.gz
    
    # Extract
    tar -xzf hugo_extended_${HUGO_VERSION}_linux-amd64.tar.gz
    
    # Move to /usr/local/bin (requires sudo)
    echo "Moving Hugo to /usr/local/bin (requires sudo)..."
    sudo mv hugo /usr/local/bin/
    
    # Cleanup
    rm hugo_extended_${HUGO_VERSION}_linux-amd64.tar.gz LICENSE README.md 2>/dev/null || true
    
    echo "✅ Hugo installed successfully!"
else
    echo "✅ Hugo is already installed: $(hugo version | head -n1)"
fi

echo ""

# Check if theme exists
if [ ! -d "themes/geekdoc" ]; then
    echo "🎨 Installing Geekdoc theme..."
    git clone https://github.com/thegeeklab/hugo-geekdoc.git themes/geekdoc
    echo "✅ Theme installed!"
else
    echo "✅ Theme already installed!"
fi

echo ""
echo "🎉 Setup complete!"
echo ""
echo "To start the documentation server:"
echo "  cd hugo-docs"
echo "  hugo server -D"
echo ""
echo "Then open: http://localhost:1313"
echo ""
