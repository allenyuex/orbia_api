#!/bin/bash

BASE_URL="http://localhost:8888"

echo "🧪 Testing Orbia API..."
echo ""

# 颜色
GREEN='\033[0;32m'
RED='\033[0;31m'
BLUE='\033[0;34m'
NC='\033[0m'

test_api() {
    local method=$1
    local endpoint=$2
    local data=$3
    local desc=$4
    
    echo -e "${BLUE}Testing: ${desc}${NC}"
    echo "  ${method} ${endpoint}"
    
    if [ -z "$data" ]; then
        response=$(curl -s -w "\n%{http_code}" -X ${method} "${BASE_URL}${endpoint}")
    else
        response=$(curl -s -w "\n%{http_code}" -X ${method} "${BASE_URL}${endpoint}" \
            -H "Content-Type: application/json" \
            -d "${data}")
    fi
    
    http_code=$(echo "$response" | tail -n1)
    body=$(echo "$response" | sed '$d')
    
    if [ "$http_code" -ge 200 ] && [ "$http_code" -lt 300 ]; then
        echo -e "  ${GREEN}✓ ${http_code}${NC}"
        echo "  ${body}"
    else
        echo -e "  ${RED}✗ ${http_code}${NC}"
        echo "  ${body}"
    fi
    echo ""
}

# 等待服务启动
echo "⏳ Waiting for server..."
for i in {1..10}; do
    if curl -s "${BASE_URL}/health" > /dev/null 2>&1; then
        echo -e "${GREEN}✓ Server is ready${NC}"
        echo ""
        break
    fi
    if [ $i -eq 10 ]; then
        echo -e "${RED}✗ Server not responding${NC}"
        exit 1
    fi
    sleep 1
done

# 测试接口
test_api "GET" "/" "" "Welcome page"
test_api "GET" "/health" "" "Health check"
test_api "GET" "/api/v1/demo/hello?name=Orbia" "" "Hello Demo"

test_api "POST" "/api/v1/users" \
    '{"name":"张三","email":"zhangsan@example.com","phone":"13800138000"}' \
    "Create user"

test_api "GET" "/api/v1/users/1" "" "Get user by ID"
test_api "GET" "/api/v1/users" "" "List users"

echo -e "${GREEN}✅ Tests completed!${NC}"

