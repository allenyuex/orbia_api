#!/bin/bash

# 角色权限系统测试脚本

set -e

# 颜色定义
GREEN='\033[0;32m'
RED='\033[0;31m'
YELLOW='\033[1;33m'
BLUE='\033[0;34m'
NC='\033[0m' # No Color

# 配置
API_URL=${API_URL:-"http://localhost:8080"}
TEST_EMAIL="test_user_$(date +%s)@example.com"
ADMIN_EMAIL="admin_$(date +%s)@example.com"
PASSWORD="test123456"

echo -e "${BLUE}========================================${NC}"
echo -e "${BLUE}   用户角色权限系统测试${NC}"
echo -e "${BLUE}========================================${NC}"
echo ""

# 测试 1: 注册普通用户
echo -e "${YELLOW}[测试 1] 注册普通用户${NC}"
REGISTER_RESP=$(curl -s -X POST "$API_URL/api/v1/auth/email-login" \
  -H "Content-Type: application/json" \
  -d "{\"email\":\"$TEST_EMAIL\",\"password\":\"$PASSWORD\"}")

USER_TOKEN=$(echo $REGISTER_RESP | jq -r '.token // empty')

if [ -z "$USER_TOKEN" ]; then
  echo -e "${RED}✗ 用户注册失败${NC}"
  echo $REGISTER_RESP | jq .
  exit 1
fi

echo -e "${GREEN}✓ 普通用户注册成功${NC}"
echo "  Token: ${USER_TOKEN:0:20}..."
echo ""

# 测试 2: 获取用户信息，验证角色为 user
echo -e "${YELLOW}[测试 2] 验证普通用户角色${NC}"
PROFILE_RESP=$(curl -s -X POST "$API_URL/api/v1/user/profile" \
  -H "Authorization: Bearer $USER_TOKEN" \
  -H "Content-Type: application/json")

USER_ROLE=$(echo $PROFILE_RESP | jq -r '.user.role // empty')

if [ "$USER_ROLE" != "user" ]; then
  echo -e "${RED}✗ 用户角色不正确，期望: user, 实际: $USER_ROLE${NC}"
  echo $PROFILE_RESP | jq .
  exit 1
fi

echo -e "${GREEN}✓ 用户角色正确: $USER_ROLE${NC}"
echo ""

# 测试 3: 普通用户尝试访问管理员 API（应该失败）
echo -e "${YELLOW}[测试 3] 普通用户访问管理员 API（应该被拒绝）${NC}"
REVIEW_RESP=$(curl -s -X POST "$API_URL/api/v1/kol/review" \
  -H "Authorization: Bearer $USER_TOKEN" \
  -H "Content-Type: application/json" \
  -d '{"kol_id":1,"status":"approved"}')

REVIEW_CODE=$(echo $REVIEW_RESP | jq -r '.base_resp.code // .code // empty')

if [ "$REVIEW_CODE" == "403" ]; then
  echo -e "${GREEN}✓ 权限验证成功：普通用户被正确拒绝${NC}"
else
  echo -e "${RED}✗ 权限验证失败：普通用户不应该能访问管理员 API${NC}"
  echo "  响应码: $REVIEW_CODE"
  echo $REVIEW_RESP | jq .
  exit 1
fi
echo ""

# 测试 4: 注册管理员用户
echo -e "${YELLOW}[测试 4] 注册管理员用户${NC}"
ADMIN_REGISTER_RESP=$(curl -s -X POST "$API_URL/api/v1/auth/email-login" \
  -H "Content-Type: application/json" \
  -d "{\"email\":\"$ADMIN_EMAIL\",\"password\":\"$PASSWORD\"}")

ADMIN_TOKEN=$(echo $ADMIN_REGISTER_RESP | jq -r '.token // empty')

if [ -z "$ADMIN_TOKEN" ]; then
  echo -e "${RED}✗ 管理员注册失败${NC}"
  echo $ADMIN_REGISTER_RESP | jq .
  exit 1
fi

echo -e "${GREEN}✓ 管理员账号创建成功${NC}"
echo "  Token: ${ADMIN_TOKEN:0:20}..."
echo ""

# 测试 5: 将用户提升为管理员
echo -e "${YELLOW}[测试 5] 提升用户为管理员${NC}"
mysql -h127.0.0.1 -P3306 -uroot -proot123 -e \
  "USE orbia; UPDATE orbia_user SET role = 'admin' WHERE email = '$ADMIN_EMAIL';" 2>/dev/null

if [ $? -eq 0 ]; then
  echo -e "${GREEN}✓ 用户角色已提升为管理员${NC}"
else
  echo -e "${RED}✗ 提升用户角色失败${NC}"
  exit 1
fi
echo ""

# 测试 6: 验证管理员角色
echo -e "${YELLOW}[测试 6] 验证管理员角色${NC}"
# 需要重新登录获取新 token（包含新角色信息）
ADMIN_LOGIN_RESP=$(curl -s -X POST "$API_URL/api/v1/auth/email-login" \
  -H "Content-Type: application/json" \
  -d "{\"email\":\"$ADMIN_EMAIL\",\"password\":\"$PASSWORD\"}")

ADMIN_TOKEN=$(echo $ADMIN_LOGIN_RESP | jq -r '.token // empty')

ADMIN_PROFILE_RESP=$(curl -s -X POST "$API_URL/api/v1/user/profile" \
  -H "Authorization: Bearer $ADMIN_TOKEN" \
  -H "Content-Type: application/json")

ADMIN_ROLE=$(echo $ADMIN_PROFILE_RESP | jq -r '.user.role // empty')

if [ "$ADMIN_ROLE" != "admin" ]; then
  echo -e "${RED}✗ 管理员角色不正确，期望: admin, 实际: $ADMIN_ROLE${NC}"
  echo $ADMIN_PROFILE_RESP | jq .
  exit 1
fi

echo -e "${GREEN}✓ 管理员角色正确: $ADMIN_ROLE${NC}"
echo ""

# 测试 7: 管理员访问管理员 API（应该成功）
echo -e "${YELLOW}[测试 7] 管理员访问管理员 API（应该成功）${NC}"
ADMIN_REVIEW_RESP=$(curl -s -X POST "$API_URL/api/v1/kol/review" \
  -H "Authorization: Bearer $ADMIN_TOKEN" \
  -H "Content-Type: application/json" \
  -d '{"kol_id":1,"status":"approved"}')

ADMIN_REVIEW_CODE=$(echo $ADMIN_REVIEW_RESP | jq -r '.base_resp.code // .code // empty')

if [ "$ADMIN_REVIEW_CODE" == "403" ]; then
  echo -e "${RED}✗ 管理员应该能访问管理员 API${NC}"
  echo "  响应码: $ADMIN_REVIEW_CODE"
  echo $ADMIN_REVIEW_RESP | jq .
  exit 1
else
  echo -e "${GREEN}✓ 管理员成功访问管理员 API${NC}"
  echo "  响应码: $ADMIN_REVIEW_CODE"
fi
echo ""

# 清理测试数据
echo -e "${YELLOW}[清理] 删除测试账号${NC}"
mysql -h127.0.0.1 -P3306 -uroot -proot123 -e \
  "USE orbia; DELETE FROM orbia_user WHERE email IN ('$TEST_EMAIL', '$ADMIN_EMAIL');" 2>/dev/null

if [ $? -eq 0 ]; then
  echo -e "${GREEN}✓ 测试数据已清理${NC}"
else
  echo -e "${YELLOW}⚠ 清理测试数据失败（可能需要手动清理）${NC}"
fi
echo ""

# 测试总结
echo -e "${BLUE}========================================${NC}"
echo -e "${GREEN}✓ 所有测试通过！${NC}"
echo -e "${BLUE}========================================${NC}"
echo ""
echo "测试摘要："
echo "  - 普通用户注册和认证: ✓"
echo "  - 角色信息正确返回: ✓"
echo "  - 普通用户访问权限限制: ✓"
echo "  - 管理员角色提升: ✓"
echo "  - 管理员访问权限: ✓"
echo ""
echo -e "${GREEN}角色权限系统运行正常！${NC}"


