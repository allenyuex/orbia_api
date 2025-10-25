# 管理员接口文档

本文档描述所有管理员专用接口。所有接口都需要管理员权限，使用 JWT Bearer Token 认证。

## 鉴权说明

所有管理员接口都需要：
1. 在请求头中携带 `Authorization: Bearer {token}`
2. 用户必须具有 `admin` 角色
3. 如果用户不是管理员，将返回 403 Forbidden

## 接口列表

### 1. 用户管理

#### 1.1 获取所有用户列表

**接口地址：** `POST /api/v1/admin/users`

**功能描述：** 查询所有用户的列表，支持分页和筛选

**请求参数：**
```json
{
  "keyword": "search_text",    // 可选，搜索关键字（用户名、邮箱、钱包地址）
  "role": "user",              // 可选，角色筛选：user | admin
  "status": "normal",          // 可选，状态筛选：normal | disabled | deleted
  "page": 1,                   // 可选，页码，默认 1
  "page_size": 10              // 可选，每页数量，默认 10，最大 100
}
```

**响应示例：**
```json
{
  "code": 0,
  "message": "success",
  "users": [
    {
      "id": 1,
      "wallet_address": "0x1234...",
      "email": "user@example.com",
      "nickname": "用户名",
      "avatar_url": "https://...",
      "role": "user",
      "status": "normal",
      "kol_id": 1,
      "created_at": "2024-01-01 00:00:00",
      "updated_at": "2024-01-01 00:00:00"
    }
  ],
  "page_info": {
    "page": 1,
    "page_size": 10,
    "total": 100,
    "total_pages": 10
  }
}
```

#### 1.2 设置用户状态

**接口地址：** `POST /api/v1/admin/user/status`

**功能描述：** 设置用户的状态（禁用、正常、已删除）

**请求参数：**
```json
{
  "user_id": 1,
  "status": "disabled"  // normal | disabled | deleted
}
```

**注意事项：**
- 管理员用户不能被禁用或删除
- 只有普通用户才可以修改状态

**响应示例：**
```json
{
  "code": 0,
  "message": "success"
}
```

**错误响应：**
```json
{
  "code": 500,
  "message": "cannot modify admin user status"
}
```

### 2. KOL 管理

#### 2.1 获取所有 KOL 列表

**接口地址：** `POST /api/v1/admin/kols`

**功能描述：** 查询所有 KOL 的列表，支持分页和筛选

**请求参数：**
```json
{
  "keyword": "kol_name",       // 可选，搜索关键字（显示名称、国家）
  "status": "pending",         // 可选，状态筛选：pending | approved | rejected
  "country": "US",             // 可选，国家筛选
  "tag": "DeFi",              // 可选，标签筛选
  "page": 1,                   // 可选，页码，默认 1
  "page_size": 10              // 可选，每页数量，默认 10
}
```

**响应示例：**
```json
{
  "code": 0,
  "message": "success",
  "kols": [
    {
      "id": 1,
      "user_id": 1,
      "display_name": "KOL Name",
      "avatar_url": "https://...",
      "country": "US",
      "status": "approved",
      "total_followers": 100000,
      "created_at": "2024-01-01 00:00:00",
      "updated_at": "2024-01-01 00:00:00"
    }
  ],
  "page_info": {
    "page": 1,
    "page_size": 10,
    "total": 50,
    "total_pages": 5
  }
}
```

#### 2.2 审核 KOL 申请

**接口地址：** `POST /api/v1/admin/kol/review`

**功能描述：** 审核 KOL 的申请，可以批准或拒绝

**请求参数：**
```json
{
  "kol_id": 1,
  "status": "approved",             // approved | rejected
  "reject_reason": "不符合要求"     // 可选，拒绝时必填
}
```

**响应示例：**
```json
{
  "code": 0,
  "message": "success"
}
```

### 3. 团队管理

#### 3.1 获取所有团队列表

**接口地址：** `POST /api/v1/admin/teams`

**功能描述：** 查询所有团队的列表，支持分页和模糊搜索

**请求参数：**
```json
{
  "keyword": "team_name",      // 可选，搜索关键字（团队名称）
  "page": 1,                   // 可选，页码，默认 1
  "page_size": 10              // 可选，每页数量，默认 10
}
```

**响应示例：**
```json
{
  "code": 0,
  "message": "success",
  "teams": [
    {
      "id": 1,
      "name": "Team Name",
      "icon_url": "https://...",
      "creator_id": 1,
      "creator_name": "Creator",
      "member_count": 5,
      "created_at": "2024-01-01 00:00:00"
    }
  ],
  "page_info": {
    "page": 1,
    "page_size": 10,
    "total": 30,
    "total_pages": 3
  }
}
```

#### 3.2 获取特定团队的所有成员

**接口地址：** `POST /api/v1/admin/team/:team_id/members`

**功能描述：** 查询特定团队的所有成员列表

**URL 参数：**
- `team_id`: 团队 ID

**请求参数：**
```json
{
  "page": 1,                   // 可选，页码，默认 1
  "page_size": 10              // 可选，每页数量，默认 10
}
```

**响应示例：**
```json
{
  "code": 0,
  "message": "success",
  "members": [
    {
      "user_id": 1,
      "nickname": "User Name",
      "email": "user@example.com",
      "avatar_url": "https://...",
      "role": "creator",        // creator | owner | member
      "joined_at": "2024-01-01 00:00:00"
    }
  ],
  "page_info": {
    "page": 1,
    "page_size": 10,
    "total": 5,
    "total_pages": 1
  }
}
```

### 4. 订单管理

#### 4.1 获取所有订单列表

**接口地址：** `POST /api/v1/admin/orders`

**功能描述：** 查询所有订单的列表，支持分页和模糊搜索

**请求参数：**
```json
{
  "keyword": "search_text",    // 可选，搜索关键字（订单ID、用户名、邮箱、钱包地址）
  "status": "pending",         // 可选，状态筛选
  "page": 1,                   // 可选，页码，默认 1
  "page_size": 10              // 可选，每页数量，默认 10
}
```

**订单状态说明：**
- `pending`: 待确认
- `confirmed`: 已确认
- `in_progress`: 进行中
- `completed`: 已完成
- `cancelled`: 已取消
- `refunded`: 已退款

**响应示例：**
```json
{
  "code": 0,
  "message": "success",
  "orders": [
    {
      "order_id": "ORD123456789",
      "user_id": 1,
      "user_name": "User Name",
      "user_email": "user@example.com",
      "kol_id": 1,
      "kol_name": "KOL Name",
      "plan_title": "Basic Plan",
      "plan_price": 100.00,
      "status": "completed",
      "created_at": "2024-01-01 00:00:00",
      "completed_at": "2024-01-02 00:00:00"
    }
  ],
  "page_info": {
    "page": 1,
    "page_size": 10,
    "total": 200,
    "total_pages": 20
  }
}
```

### 5. 钱包管理

#### 5.1 获取特定用户的钱包信息

**接口地址：** `POST /api/v1/admin/user/:user_id/wallet`

**功能描述：** 查询特定用户的钱包详细信息

**URL 参数：**
- `user_id`: 用户 ID

**请求参数：** 无需请求体

**响应示例：**
```json
{
  "code": 0,
  "message": "success",
  "wallet": {
    "user_id": 1,
    "user_name": "User Name",
    "user_email": "user@example.com",
    "balance": 1000.00,
    "frozen_balance": 100.00,
    "total_recharge": 5000.00,
    "total_consume": 4000.00,
    "created_at": "2024-01-01 00:00:00",
    "updated_at": "2024-01-02 00:00:00"
  }
}
```

## 通用响应格式

### 成功响应
```json
{
  "code": 0,
  "message": "success",
  "data": { ... }
}
```

### 错误响应
```json
{
  "code": 400/403/500,
  "message": "错误信息描述"
}
```

## 错误码说明

| 错误码 | 说明 |
|--------|------|
| 0 | 成功 |
| 400 | 请求参数错误 |
| 401 | 未认证或 Token 无效 |
| 403 | 权限不足（非管理员） |
| 404 | 资源不存在 |
| 500 | 服务器内部错误 |

## 分页说明

所有列表接口都支持分页，分页参数：
- `page`: 页码，从 1 开始，默认 1
- `page_size`: 每页数量，默认 10，最大 100

分页响应包含 `page_info` 字段：
```json
{
  "page": 1,           // 当前页码
  "page_size": 10,     // 每页数量
  "total": 100,        // 总记录数
  "total_pages": 10    // 总页数
}
```

## 搜索说明

支持模糊搜索的字段会在接口描述中标注。搜索关键字通过 `keyword` 参数传递，系统会在多个字段中进行模糊匹配。

## 使用示例

### 请求示例（使用 curl）

```bash
# 获取所有用户列表
curl -X POST "http://localhost:8888/api/v1/admin/users" \
  -H "Content-Type: application/json" \
  -H "Authorization: Bearer YOUR_ADMIN_TOKEN" \
  -d '{
    "keyword": "test",
    "page": 1,
    "page_size": 20
  }'

# 设置用户状态
curl -X POST "http://localhost:8888/api/v1/admin/user/status" \
  -H "Content-Type: application/json" \
  -H "Authorization: Bearer YOUR_ADMIN_TOKEN" \
  -d '{
    "user_id": 1,
    "status": "disabled"
  }'

# 审核 KOL
curl -X POST "http://localhost:8888/api/v1/admin/kol/review" \
  -H "Content-Type: application/json" \
  -H "Authorization: Bearer YOUR_ADMIN_TOKEN" \
  -d '{
    "kol_id": 1,
    "status": "approved"
  }'
```

## 注意事项

1. **权限控制**：所有接口都严格检查管理员权限，非管理员用户无法访问
2. **保护机制**：管理员用户自身不能被禁用或删除
3. **数据完整性**：修改操作会记录操作时间和相关信息
4. **性能优化**：大数据量查询建议使用合适的分页参数
5. **搜索优化**：关键字搜索支持部分匹配，建议输入至少 2 个字符

## 数据字典

### 用户状态 (User Status)
- `normal`: 正常
- `disabled`: 禁用
- `deleted`: 已删除

### 用户角色 (User Role)
- `user`: 普通用户
- `admin`: 管理员

### KOL 状态 (KOL Status)
- `pending`: 待审核
- `approved`: 已通过
- `rejected`: 已拒绝

### 团队成员角色 (Team Member Role)
- `creator`: 创建者
- `owner`: 拥有者
- `member`: 成员

### 订单状态 (Order Status)
- `pending`: 待确认
- `confirmed`: 已确认
- `in_progress`: 进行中
- `completed`: 已完成
- `cancelled`: 已取消
- `refunded`: 已退款

