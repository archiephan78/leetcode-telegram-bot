#!/bin/bash

# Multi-Architecture Docker Build Script
# Supports both AMD64 (Intel) and ARM64 (Apple Silicon)

echo "🐳 Building Multi-Architecture Docker Image..."

# Create buildx builder if not exists
if ! docker buildx ls | grep -q multiarch; then
    echo "📦 Creating multi-arch builder..."
    docker buildx create --name multiarch --use
fi

# Use the multiarch builder
docker buildx use multiarch

# Build for multiple platforms
echo "🔨 Building for linux/amd64 and linux/arm64..."
docker buildx build \
    --platform linux/amd64,linux/arm64 \
    -t leetcode-telegram-bot:latest \
    --load \
    .

if [ $? -eq 0 ]; then
    echo "✅ Multi-architecture build completed successfully!"
    echo ""
    echo "📋 Available commands:"
    echo "  docker run -d --name leetcode-bot \\"
    echo "    -e TELEGRAM_BOT_TOKEN=\"574513532:AAFN3cEsV48DfFUv90wYhITiPb-nlFQ81Pg\" \\"
    echo "    -e TELEGRAM_GROUP_ID=\"-4867864977\" \\"
    echo "    -v \$(pwd)/data:/data \\"
    echo "    leetcode-telegram-bot:latest"
    echo ""
    echo "  Or use Docker Compose:"
    echo "  docker-compose up -d"
else
    echo "❌ Build failed!"
    exit 1
fi 