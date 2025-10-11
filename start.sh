#!/bin/bash

echo "🚀 Starting Orbia API..."
echo ""

cd "$(dirname "$0")"

# 检查配置文件
if [ ! -f "conf/config.yaml" ]; then
    echo "❌ Config file not found: conf/config.yaml"
    exit 1
fi

# 下载依赖
echo "📥 Installing dependencies..."
go mod tidy

# 运行服务
echo ""
echo "✨ Starting server..."
echo ""
go run .

