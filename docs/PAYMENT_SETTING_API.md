# 收款钱包设置 API 文档

## 概述

收款钱包设置功能用于管理员配置用户充值时的目标钱包地址。普通用户在充值时可以从启用的收款钱包列表中选择目标地址进行转账。

### 特性

- ✅ 支持多个区块链网络配置
- ✅ 软删除机制，保证数据一致性
- ✅ 状态管理（启用/禁用）
- ✅ 完整的CRUD操作
- ✅ 管理员权限管理收款钱包
- ✅ 用户端查询启用的收款钱包

## 数据库设计

### 收款钱包设置表 (orbia_payment_setting)

| 字段 | 类型 | 说明 |
|------|------|------|
| id | BIGINT | 主键ID |
| network | VARCHAR(100) | 区块链网络（如：TRC-20 - TRON Network (TRC-20)）|
| address | VARCHAR(500) | 钱包地址 |
| label | VARCHAR(200) | 钱包标签（如：USDT-TRC20 主钱包）|
| status | TINYINT | 状态：1-启用，0-禁用 |
| created_at | TIMESTAMP | 创建时间 |
| updated_at | TIMESTAMP | 更新时间 |
| deleted_at | TIMESTAMP | 软删除时间 |

## API 接口

### 认证说明

管理员接口需要管理员权限，用户端接口需要用户权限。

**请求头：**
```
Authorization: Bearer {token}
```

---

## 管理员接口

### 1. 获取收款钱包设置列表

**接口地址：** `POST /api/v1/admin/payment-settings/list`

**权限要求：** 管理员

**请求参数：**
```json
{
  "network": "TRC-20 - TRON Network (TRC-20)",
  "status": 1,
  "page": 1,
  "page_size": 20
}
```

| 参数 | 类型 | 必填 | 说明 |
|------|------|------|------|
| network | string | 否 | 区块链网络筛选 |
| status | int | 否 | 状态筛选：1-启用，0-禁用 |
| page | int | 否 | 页码，默认1 |
| page_size | int | 否 | 每页条数，默认20 |

**响应示例：**
```json
{
  "base_resp": {
    "code": 0,
    "message": "success"
  },
  "list": [
    {
      "id": 1,
      "network": "TRC-20 - TRON Network (TRC-20)",
      "address": "TYVNBvUExGmYqJsrjq3dKyJy8CfPKkrmPL",
      "label": "USDT-TRC20 主钱包",
      "status": 1,
      "created_at": "2024-01-01 12:00:00",
      "updated_at": "2024-01-01 12:00:00"
    }
  ],
  "total": 1,
  "page": 1,
  "page_size": 20
}
```

---

### 2. 获取收款钱包设置详情

**接口地址：** `POST /api/v1/admin/payment-settings/:id`

**权限要求：** 管理员

**请求参数：**

路径参数：
- `id`: 收款钱包设置ID

**响应示例：**
```json
{
  "base_resp": {
    "code": 0,
    "message": "success"
  },
  "setting": {
    "id": 1,
    "network": "TRC-20 - TRON Network (TRC-20)",
    "address": "TYVNBvUExGmYqJsrjq3dKyJy8CfPKkrmPL",
    "label": "USDT-TRC20 主钱包",
    "status": 1,
    "created_at": "2024-01-01 12:00:00",
    "updated_at": "2024-01-01 12:00:00"
  }
}
```

---

### 3. 创建收款钱包设置

**接口地址：** `POST /api/v1/admin/payment-settings/create`

**权限要求：** 管理员

**请求参数：**
```json
{
  "network": "TRC-20 - TRON Network (TRC-20)",
  "address": "TYVNBvUExGmYqJsrjq3dKyJy8CfPKkrmPL",
  "label": "USDT-TRC20 主钱包",
  "status": 1
}
```

| 参数 | 类型 | 必填 | 说明 |
|------|------|------|------|
| network | string | 是 | 区块链网络 |
| address | string | 是 | 钱包地址 |
| label | string | 是 | 钱包标签 |
| status | int | 否 | 状态：1-启用，0-禁用，默认1 |

**响应示例：**
```json
{
  "base_resp": {
    "code": 0,
    "message": "success"
  },
  "setting": {
    "id": 1,
    "network": "TRC-20 - TRON Network (TRC-20)",
    "address": "TYVNBvUExGmYqJsrjq3dKyJy8CfPKkrmPL",
    "label": "USDT-TRC20 主钱包",
    "status": 1,
    "created_at": "2024-01-01 12:00:00",
    "updated_at": "2024-01-01 12:00:00"
  }
}
```

---

### 4. 更新收款钱包设置

**接口地址：** `POST /api/v1/admin/payment-settings/update`

**权限要求：** 管理员

**请求参数：**
```json
{
  "id": 1,
  "network": "TRC-20 - TRON Network (TRC-20)",
  "address": "TYVNBvUExGmYqJsrjq3dKyJy8CfPKkrmPL",
  "label": "USDT-TRC20 备用钱包",
  "status": 0
}
```

| 参数 | 类型 | 必填 | 说明 |
|------|------|------|------|
| id | int64 | 是 | 收款钱包设置ID |
| network | string | 否 | 区块链网络 |
| address | string | 否 | 钱包地址 |
| label | string | 否 | 钱包标签 |
| status | int | 否 | 状态：1-启用，0-禁用 |

**响应示例：**
```json
{
  "base_resp": {
    "code": 0,
    "message": "success"
  },
  "setting": {
    "id": 1,
    "network": "TRC-20 - TRON Network (TRC-20)",
    "address": "TYVNBvUExGmYqJsrjq3dKyJy8CfPKkrmPL",
    "label": "USDT-TRC20 备用钱包",
    "status": 0,
    "created_at": "2024-01-01 12:00:00",
    "updated_at": "2024-01-01 13:00:00"
  }
}
```

---

### 5. 删除收款钱包设置

**接口地址：** `POST /api/v1/admin/payment-settings/delete`

**权限要求：** 管理员

**请求参数：**
```json
{
  "id": 1
}
```

| 参数 | 类型 | 必填 | 说明 |
|------|------|------|------|
| id | int64 | 是 | 收款钱包设置ID |

**响应示例：**
```json
{
  "base_resp": {
    "code": 0,
    "message": "success"
  }
}
```

**注意：** 删除操作为软删除，不会真正从数据库中删除记录。

---

## 用户端接口

### 6. 获取启用的收款钱包设置列表

**接口地址：** `POST /api/v1/payment-settings/active`

**权限要求：** 用户或管理员

**请求参数：**
```json
{
  "network": "TRC-20 - TRON Network (TRC-20)"
}
```

| 参数 | 类型 | 必填 | 说明 |
|------|------|------|------|
| network | string | 否 | 区块链网络筛选 |

**响应示例：**
```json
{
  "base_resp": {
    "code": 0,
    "message": "success"
  },
  "list": [
    {
      "id": 1,
      "network": "TRC-20 - TRON Network (TRC-20)",
      "address": "TYVNBvUExGmYqJsrjq3dKyJy8CfPKkrmPL",
      "label": "USDT-TRC20 主钱包",
      "status": 1,
      "created_at": "2024-01-01 12:00:00",
      "updated_at": "2024-01-01 12:00:00"
    },
    {
      "id": 2,
      "network": "ERC-20 - Ethereum Network (ERC-20)",
      "address": "0x1234567890123456789012345678901234567890",
      "label": "USDT-ERC20 主钱包",
      "status": 1,
      "created_at": "2024-01-01 12:00:00",
      "updated_at": "2024-01-01 12:00:00"
    }
  ]
}
```

**说明：** 此接口只返回 status=1（启用）的收款钱包设置。

---

## 错误码说明

| 错误码 | 说明 |
|--------|------|
| 0 | 成功 |
| 400 | 请求参数错误 |
| 401 | 未授权 |
| 403 | 无权限 |
| 404 | 资源不存在 |
| 500 | 服务器内部错误 |

---

## 使用场景示例

### 场景1: 管理员添加新的收款钱包

1. 管理员登录系统获取 admin token
2. 调用创建收款钱包设置接口：
```bash
curl -X POST 'http://localhost:8080/api/v1/admin/payment-settings/create' \
  -H 'Authorization: Bearer {admin_token}' \
  -H 'Content-Type: application/json' \
  -d '{
    "network": "TRC-20 - TRON Network (TRC-20)",
    "address": "TYVNBvUExGmYqJsrjq3dKyJy8CfPKkrmPL",
    "label": "USDT-TRC20 主钱包",
    "status": 1
  }'
```

### 场景2: 用户充值时获取收款地址

1. 用户登录系统获取 token（user 或 admin 都可以）
2. 调用获取启用的收款钱包列表接口：
```bash
curl -X POST 'http://localhost:8080/api/v1/payment-settings/active' \
  -H 'Authorization: Bearer {token}' \
  -H 'Content-Type: application/json' \
  -d '{
    "network": "TRC-20 - TRON Network (TRC-20)"
  }'
```
3. 从返回的列表中选择合适的钱包地址进行转账

### 场景3: 管理员禁用某个收款钱包

1. 管理员调用更新接口将 status 设置为 0：
```bash
curl -X POST 'http://localhost:8080/api/v1/admin/payment-settings/update' \
  -H 'Authorization: Bearer {admin_token}' \
  -H 'Content-Type: application/json' \
  -d '{
    "id": 1,
    "status": 0
  }'
```

---

## 常见问题

### Q1: 可以为同一个网络配置多个收款地址吗？
**A:** 可以。系统支持为同一个区块链网络配置多个收款地址，用户可以从中选择。

### Q2: 删除收款钱包设置后数据会丢失吗？
**A:** 不会。系统采用软删除机制，删除后的数据仍然保留在数据库中，只是标记为已删除状态。

### Q3: 用户端能看到禁用的收款钱包吗？
**A:** 不能。用户端接口只返回 status=1（启用）的收款钱包设置。

### Q4: 支持哪些区块链网络？
**A:** 系统不限制区块链网络类型，管理员可以配置任意网络。常见的包括：
- TRC-20 - TRON Network (TRC-20)
- ERC-20 - Ethereum Network (ERC-20)
- BEP-20 - BSC Network (BEP-20)
- Polygon Network
- 等等

### Q5: 钱包地址的格式有验证吗？
**A:** 目前系统不对钱包地址格式进行验证，管理员需要确保填写的地址正确。建议在添加钱包地址时仔细核对。

---

## 最佳实践

1. **地址验证**: 在添加收款钱包时，请仔细核对钱包地址，确保准确无误
2. **标签命名**: 建议使用清晰的标签命名，如 "USDT-TRC20 主钱包"、"USDT-ERC20 备用钱包" 等
3. **状态管理**: 不再使用的收款钱包建议禁用而不是删除，保留历史记录
4. **网络匹配**: 确保钱包地址与区块链网络匹配，避免用户转账错误
5. **备用方案**: 建议为每个网络至少配置两个收款钱包，作为备用

---

## 技术架构

### 目录结构
```
orbia_api/
├── biz/
│   ├── dal/
│   │   ├── model/
│   │   │   └── orbia_payment_setting.gen.go  # 数据模型
│   │   └── mysql/
│   │       └── payment_setting.go             # 数据访问层
│   ├── service/
│   │   └── payment_setting/
│   │       └── payment_setting_service.go     # 业务逻辑层
│   ├── handler/
│   │   └── payment_setting/
│   │       └── payment_setting_service.go     # 处理层
│   ├── router/
│   │   └── payment_setting/
│   │       ├── payment_setting.go             # 路由定义
│   │       └── middleware.go                  # 中间件
│   └── model/
│       └── payment_setting/
│           └── payment_setting.go             # API模型
├── idl/
│   └── payment_setting.thrift                 # IDL定义
└── sql/
    └── init.sql                               # 数据库初始化脚本
```

### 技术栈
- **Web框架**: CloudWeGo Hertz
- **ORM**: GORM
- **IDL**: Thrift
- **数据库**: MySQL 8.0
- **认证**: JWT

---

## 更新日志

### v1.0.0 (2024-01-01)
- ✅ 初始版本发布
- ✅ 支持收款钱包设置的CRUD操作
- ✅ 支持管理员和用户端分离的接口
- ✅ 支持按网络和状态筛选
- ✅ 软删除机制

---

## 联系方式

如有问题或建议，请联系开发团队。

