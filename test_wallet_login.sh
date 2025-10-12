#!/bin/bash

# 钱包登录测试脚本

echo "🧪 Testing Orbia API Wallet Login..."

# 服务器地址
BASE_URL="http://localhost:8888"

# 测试健康检查
echo "1. Testing health check..."
curl -s "$BASE_URL/health" | jq .
echo ""

# 测试钱包登录（需要真实的钱包地址和签名）
echo "2. Testing wallet login..."
echo "⚠️  Note: This requires a real wallet address and signature"
echo "   You can generate test data using MetaMask or other wallet tools"
echo ""

# 示例请求（需要替换为真实数据）
cat << 'EOF'
Example wallet login request:

curl -X POST http://localhost:8888/api/v1/auth/wallet-login \
  -H "Content-Type: application/json" \
  -d '{
    "wallet_address": "0x742d35Cc6634C0532925a3b8D4C9db96C4b4d8b6",
    "signature": "0x...",
    "message": "Welcome to Orbia!\n\nWallet: 0x742d35Cc6634C0532925a3b8D4C9db96C4b4d8b6\nTimestamp: 1234567890\n\nThis request will not trigger a blockchain transaction or cost any gas fees."
  }'

Expected response:
{
  "token": "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9...",
  "expires_in": 604800,
  "base_resp": {
    "code": 200,
    "message": "Login successful"
  }
}
EOF

echo ""
echo "3. Testing user profile (requires JWT token)..."
echo "   First login to get a token, then use it in Authorization header:"
echo ""

cat << 'EOF'
Example profile request:

curl -X POST http://localhost:8888/api/v1/user/profile \
  -H "Content-Type: application/json" \
  -H "Authorization: Bearer YOUR_JWT_TOKEN_HERE" \
  -d '{}'

Expected response:
{
  "user": {
    "id": 1,
    "wallet_address": "0x742d35Cc6634C0532925a3b8D4C9db96C4b4d8b6",
    "nickname": null,
    "avatar_url": null,
    "created_at": "2024-01-01 12:00:00",
    "updated_at": "2024-01-01 12:00:00"
  },
  "base_resp": {
    "code": 200,
    "message": "Success"
  }
}
EOF

echo ""
echo "✅ Test script completed!"
echo "💡 To run actual tests, replace the example data with real wallet signatures."