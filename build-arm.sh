#!/bin/bash

# Simple ARM64 Docker Build Script for Apple Silicon

echo "🍎 Building Docker Image for ARM64 (Apple Silicon)..."

# Build for ARM64 only
docker build --platform linux/arm64 -t leetcode-telegram-bot:arm64 .

if [ $? -eq 0 ]; then
    echo "✅ ARM64 build completed successfully!"
    echo ""
    echo "🚀 Run with:"
    echo "  docker run -d --name leetcode-bot \\"
    echo "    -e TELEGRAM_BOT_TOKEN=\"574513532:AAFN3cEsV48DfFUv90wYhITiPb-nlFQ81Pg\" \\"
    echo "    -e TELEGRAM_GROUP_ID=\"-4867864977\" \\"
    echo "    -v \$(pwd)/data:/data \\"
    echo "    leetcode-telegram-bot:arm64"
    echo ""
    echo "📊 Check logs with:"
    echo "  docker logs -f leetcode-bot"
else
    echo "❌ Build failed!"
    exit 1
fi 