#!/bin/bash

echo "🚀 Starting Orbia API..."
echo ""

cd "$(dirname "$0")"

# 设置默认环境为 dev
if [ -z "$ORBIA_ENV" ]; then
    export ORBIA_ENV="dev"
fi

echo "📋 Environment: $ORBIA_ENV"
echo ""

# 检查配置文件
CONFIG_FILE="conf/$ORBIA_ENV/config.yaml"
if [ ! -f "$CONFIG_FILE" ]; then
    echo "❌ Config file not found: $CONFIG_FILE"
    echo "💡 Available environments: dev, prod"
    exit 1
fi

echo "✅ Using config: $CONFIG_FILE"
echo ""

# 下载依赖
echo "📥 Installing dependencies..."
go mod tidy

# 运行服务
echo ""
echo "✨ Starting server..."
echo ""
go run .

