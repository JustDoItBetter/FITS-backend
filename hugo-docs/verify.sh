#!/bin/bash
# Verify Hugo documentation setup

echo "🔍 Verifying FITS Documentation Setup"
echo "======================================"
echo ""

# Check config
if [ -f "config.toml" ]; then
    echo "✅ config.toml exists"
else
    echo "❌ config.toml missing!"
    exit 1
fi

# Check theme
if [ -d "themes/geekdoc" ]; then
    echo "✅ Geekdoc theme installed"
else
    echo "⚠️  Theme not installed. Run: ./setup-arch.sh"
fi

# Check content
CONTENT_COUNT=$(find content -name "*.md" -type f | wc -l)
echo "✅ Found $CONTENT_COUNT markdown files"

# Check Hugo
if command -v hugo &> /dev/null; then
    echo "✅ Hugo is installed: $(hugo version | head -n1)"
else
    echo "❌ Hugo not installed! Run: sudo pacman -S hugo"
    exit 1
fi

echo ""
echo "🎉 Setup verification complete!"
echo ""
echo "To start the server:"
echo "  ./RUN.sh"
echo ""
