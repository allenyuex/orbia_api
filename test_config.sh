#!/bin/bash

echo "=========================================="
echo "🧪 测试配置文件加载"
echo "=========================================="
echo ""

cd "$(dirname "$0")"

# 测试 1: 默认环境 (dev)
echo "📋 测试 1: 默认环境 (无 ORBIA_ENV 环境变量)"
echo "------------------------------------------"
unset ORBIA_ENV
echo "Expected: 使用 dev 环境"
if [ -f "conf/dev/config.yaml" ]; then
    echo "✅ conf/dev/config.yaml 存在"
else
    echo "❌ conf/dev/config.yaml 不存在"
fi
echo ""

# 测试 2: dev 环境
echo "📋 测试 2: 明确指定 dev 环境"
echo "------------------------------------------"
export ORBIA_ENV=dev
echo "ORBIA_ENV=$ORBIA_ENV"
if [ -f "conf/$ORBIA_ENV/config.yaml" ]; then
    echo "✅ conf/$ORBIA_ENV/config.yaml 存在"
else
    echo "❌ conf/$ORBIA_ENV/config.yaml 不存在"
fi
echo ""

# 测试 3: prod 环境
echo "📋 测试 3: 明确指定 prod 环境"
echo "------------------------------------------"
export ORBIA_ENV=prod
echo "ORBIA_ENV=$ORBIA_ENV"
if [ -f "conf/$ORBIA_ENV/config.yaml" ]; then
    echo "✅ conf/$ORBIA_ENV/config.yaml 存在"
else
    echo "❌ conf/$ORBIA_ENV/config.yaml 不存在"
fi
echo ""

# 测试 4: 无效环境
echo "📋 测试 4: 无效环境 (staging)"
echo "------------------------------------------"
export ORBIA_ENV=staging
echo "ORBIA_ENV=$ORBIA_ENV"
if [ -f "conf/$ORBIA_ENV/config.yaml" ]; then
    echo "✅ conf/$ORBIA_ENV/config.yaml 存在"
else
    echo "❌ conf/$ORBIA_ENV/config.yaml 不存在 (预期行为)"
fi
echo ""

# 显示配置文件结构
echo "=========================================="
echo "📁 配置文件结构"
echo "=========================================="
tree conf/ 2>/dev/null || find conf/ -type f

echo ""
echo "=========================================="
echo "✅ 配置测试完成"
echo "=========================================="
echo ""
echo "💡 使用方法："
echo "   默认环境: ./start.sh"
echo "   指定环境: ORBIA_ENV=prod ./start.sh"

