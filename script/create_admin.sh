#!/bin/bash

# 创建管理员账号脚本

set -e

# 颜色定义
GREEN='\033[0;32m'
RED='\033[0;31m'
YELLOW='\033[1;33m'
BLUE='\033[0;34m'
NC='\033[0m' # No Color

echo -e "${BLUE}========================================${NC}"
echo -e "${BLUE}   创建管理员账号${NC}"
echo -e "${BLUE}========================================${NC}"
echo ""

# 数据库配置
DB_HOST=${DB_HOST:-"127.0.0.1"}
DB_PORT=${DB_PORT:-"3306"}
DB_USER=${DB_USER:-"root"}
DB_PASSWORD=${DB_PASSWORD:-"root123"}
DB_NAME=${DB_NAME:-"orbia"}

# 检查参数
if [ $# -eq 0 ]; then
    echo -e "${YELLOW}使用方法:${NC}"
    echo "  $0 <email>           # 将指定邮箱的用户提升为管理员"
    echo "  $0 <email> <id>      # 将指定邮箱或ID的用户提升为管理员"
    echo ""
    echo -e "${YELLOW}示例:${NC}"
    echo "  $0 admin@orbia.com"
    echo "  $0 user@example.com 1"
    echo ""
    exit 1
fi

EMAIL="$1"
USER_ID="${2:-}"

# 检查 mysql 命令是否可用
if ! command -v mysql &> /dev/null; then
    echo -e "${RED}错误: mysql 命令未找到${NC}"
    exit 1
fi

# 测试数据库连接
echo -e "${BLUE}测试数据库连接...${NC}"
if ! mysql -h"$DB_HOST" -P"$DB_PORT" -u"$DB_USER" -p"$DB_PASSWORD" -e "SELECT 1;" &> /dev/null; then
    echo -e "${RED}错误: 无法连接到数据库${NC}"
    exit 1
fi
echo -e "${GREEN}✓ 数据库连接成功${NC}"
echo ""

# 构建 SQL 查询
if [ -n "$USER_ID" ]; then
    WHERE_CLAUSE="email = '$EMAIL' OR id = $USER_ID"
    IDENTIFIER="email='$EMAIL' 或 id=$USER_ID"
else
    WHERE_CLAUSE="email = '$EMAIL'"
    IDENTIFIER="email='$EMAIL'"
fi

# 检查用户是否存在
echo -e "${BLUE}检查用户是否存在...${NC}"
USER_EXISTS=$(mysql -h"$DB_HOST" -P"$DB_PORT" -u"$DB_USER" -p"$DB_PASSWORD" -sN \
    -e "USE $DB_NAME; SELECT COUNT(*) FROM orbia_user WHERE $WHERE_CLAUSE;" 2>/dev/null)

if [ "$USER_EXISTS" == "0" ]; then
    echo -e "${RED}错误: 用户不存在 ($IDENTIFIER)${NC}"
    echo ""
    echo -e "${YELLOW}提示: 请先通过以下方式创建用户：${NC}"
    echo "  1. 通过 API 注册: curl -X POST http://localhost:8080/api/v1/auth/email-login"
    echo "  2. 或直接在数据库中插入用户"
    exit 1
fi

# 获取用户当前信息
echo -e "${GREEN}✓ 用户存在${NC}"
echo ""
echo -e "${BLUE}当前用户信息:${NC}"
mysql -h"$DB_HOST" -P"$DB_PORT" -u"$DB_USER" -p"$DB_PASSWORD" \
    -e "USE $DB_NAME; SELECT id, email, nickname, role, created_at FROM orbia_user WHERE $WHERE_CLAUSE;" 2>/dev/null
echo ""

# 检查用户是否已经是管理员
CURRENT_ROLE=$(mysql -h"$DB_HOST" -P"$DB_PORT" -u"$DB_USER" -p"$DB_PASSWORD" -sN \
    -e "USE $DB_NAME; SELECT role FROM orbia_user WHERE $WHERE_CLAUSE LIMIT 1;" 2>/dev/null)

if [ "$CURRENT_ROLE" == "admin" ]; then
    echo -e "${YELLOW}⚠ 用户已经是管理员${NC}"
    exit 0
fi

# 提升用户为管理员
echo -e "${BLUE}提升用户为管理员...${NC}"
mysql -h"$DB_HOST" -P"$DB_PORT" -u"$DB_USER" -p"$DB_PASSWORD" \
    -e "USE $DB_NAME; UPDATE orbia_user SET role = 'admin' WHERE $WHERE_CLAUSE;" 2>/dev/null

if [ $? -eq 0 ]; then
    echo -e "${GREEN}✓ 用户已成功提升为管理员${NC}"
    echo ""
    echo -e "${BLUE}更新后的用户信息:${NC}"
    mysql -h"$DB_HOST" -P"$DB_PORT" -u"$DB_USER" -p"$DB_PASSWORD" \
        -e "USE $DB_NAME; SELECT id, email, nickname, role, created_at FROM orbia_user WHERE $WHERE_CLAUSE;" 2>/dev/null
    echo ""
    echo -e "${GREEN}完成！${NC}"
    echo ""
    echo -e "${YELLOW}注意: 用户需要重新登录以获取新的权限${NC}"
else
    echo -e "${RED}错误: 更新用户角色失败${NC}"
    exit 1
fi


