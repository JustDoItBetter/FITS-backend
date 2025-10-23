#!/bin/bash
# Quick run script for Hugo documentation

cd "$(dirname "$0")"

if [ ! -d "themes/geekdoc" ]; then
    echo "⚠️  Theme not installed!"
    echo "Run: ./setup-linux.sh"
    exit 1
fi

echo "🚀 Starting Hugo documentation server..."
echo "📖 Documentation will be available at: http://localhost:1313"
echo "Press Ctrl+C to stop"
echo ""

hugo server -D
