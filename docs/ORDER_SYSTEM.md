# KOL 订单系统

## 概述

KOL 订单系统允许用户针对 KOL 的报价计划（Pricing Plans）进行下单，KOL 可以管理收到的订单状态。

## 核心功能

### 1. 订单 ID 生成器

使用 Snowflake 算法生成唯一的订单 ID，格式为 `ORD{snowflake_id}`。

**特点：**
- 分布式唯一 ID
- 高性能（单机可达 400万+/秒）
- 时间有序
- 可复用于其他业务场景

**实现位置：** `biz/utils/id_generator.go`

### 2. 数据库设计

**订单表（orbia_kol_order）：**

| 字段 | 类型 | 说明 |
|------|------|------|
| id | BIGINT | 自增ID（内部使用） |
| order_id | VARCHAR(64) | 订单ID（业务唯一ID） |
| user_id | BIGINT | 下单用户ID |
| team_id | BIGINT | 下单团队ID（可选） |
| kol_id | BIGINT | KOL ID |
| plan_id | BIGINT | 报价Plan ID |
| plan_title | VARCHAR(200) | Plan标题（快照） |
| plan_description | TEXT | Plan描述（快照） |
| plan_price | DECIMAL(10,2) | Plan价格（快照） |
| plan_type | VARCHAR(20) | Plan类型（快照） |
| description | TEXT | 订单描述 |
| status | ENUM | 订单状态 |
| reject_reason | TEXT | 拒绝/取消原因 |
| confirmed_at | TIMESTAMP | 确认时间 |
| completed_at | TIMESTAMP | 完成时间 |
| cancelled_at | TIMESTAMP | 取消时间 |
| created_at | TIMESTAMP | 创建时间 |
| updated_at | TIMESTAMP | 更新时间 |
| deleted_at | TIMESTAMP | 软删除时间 |

### 3. 订单状态流转

```
pending（待确认）
    ↓ [KOL确认]
confirmed（已确认）
    ↓ [KOL开始执行]
in_progress（进行中）
    ↓ [KOL完成]
completed（已完成）
    ↓ [用户申请退款]
refunded（已退款）

可随时取消：
pending/confirmed/in_progress → cancelled（已取消）
```

**状态说明：**
- `pending`: 用户刚创建订单，等待 KOL 确认
- `confirmed`: KOL 确认接单
- `in_progress`: KOL 正在执行订单
- `completed`: 订单已完成
- `cancelled`: 订单已取消（用户或 KOL 均可取消）
- `refunded`: 订单已退款

### 4. API 接口

#### 4.1 创建订单

**接口：** `POST /api/v1/order/create`

**权限：** 需要登录

**请求参数：**
```json
{
  "kol_id": 1,
  "plan_id": 1,
  "description": "推广我的项目，需要发布3个视频",
  "team_id": 1  // 可选，如果是团队下单
}
```

**响应示例：**
```json
{
  "code": 0,
  "message": "success",
  "data": {
    "order_id": "ORD7234567890123456789"
  }
}
```

**业务逻辑：**
1. 验证 KOL 是否存在且已审核通过
2. 验证 Plan 是否存在且属于该 KOL
3. 如果指定团队ID，验证用户是否属于该团队
4. 生成唯一订单ID
5. 保存 Plan 快照（价格、标题、描述等）
6. 创建订单，初始状态为 `pending`

#### 4.2 获取订单详情

**接口：** `POST /api/v1/order/detail`

**权限：** 需要登录，只能查看自己的订单或收到的订单

**请求参数：**
```json
{
  "order_id": "ORD7234567890123456789"
}
```

**响应示例：**
```json
{
  "code": 0,
  "message": "success",
  "data": {
    "order": {
      "id": 1,
      "order_id": "ORD7234567890123456789",
      "user_id": 1,
      "kol_id": 2,
      "kol_display_name": "张三",
      "kol_avatar_url": "https://...",
      "plan_id": 1,
      "plan_title": "基础推广套餐",
      "plan_description": "包含3个视频发布",
      "plan_price": 1000.00,
      "plan_type": "basic",
      "description": "推广我的项目，需要发布3个视频",
      "status": "pending",
      "created_at": "2024-01-01T00:00:00Z",
      "updated_at": "2024-01-01T00:00:00Z"
    }
  }
}
```

#### 4.3 获取订单列表

**接口：** `POST /api/v1/order/list`

**权限：** 需要登录

**请求参数：**
```json
{
  "status": "pending",  // 可选，筛选状态
  "kol_id": 1,         // 可选，筛选指定KOL的订单
  "team_id": 1,        // 可选，筛选团队订单
  "page": 1,
  "page_size": 10
}
```

**响应示例：**
```json
{
  "code": 0,
  "message": "success",
  "data": {
    "orders": [...],
    "total": 100
  }
}
```

#### 4.4 取消订单

**接口：** `POST /api/v1/order/cancel`

**权限：** 需要登录，只有下单用户可以取消

**请求参数：**
```json
{
  "order_id": "ORD7234567890123456789",
  "reason": "不需要了"
}
```

**业务逻辑：**
- 只有 `pending` 状态的订单可以取消
- 记录取消原因和取消时间

#### 4.5 获取 KOL 收到的订单列表

**接口：** `POST /api/v1/order/kol/list`

**权限：** 需要登录且是 KOL

**请求参数：**
```json
{
  "status": "pending",  // 可选
  "page": 1,
  "page_size": 10
}
```

#### 4.6 更新订单状态（KOL 端）

**接口：** `POST /api/v1/order/status/update`

**权限：** 需要登录且是订单的 KOL

**请求参数：**
```json
{
  "order_id": "ORD7234567890123456789",
  "status": "confirmed",
  "reject_reason": "时间冲突"  // 取消时必填
}
```

**业务逻辑：**
- 验证状态转换是否合法
- 记录相应的时间戳
- 如果是取消，必须提供原因

## 使用示例

### 场景1：用户下单流程

```bash
# 1. 用户登录
TOKEN="..."

# 2. 查看 KOL 的报价计划
curl -X POST "http://localhost:8888/api/v1/kol/plans" \
  -H "Authorization: Bearer ${TOKEN}" \
  -H "Content-Type: application/json" \
  -d '{"kol_id": 1}'

# 3. 创建订单
curl -X POST "http://localhost:8888/api/v1/order/create" \
  -H "Authorization: Bearer ${TOKEN}" \
  -H "Content-Type: application/json" \
  -d '{
    "kol_id": 1,
    "plan_id": 1,
    "description": "推广我的项目"
  }'

# 4. 查看订单列表
curl -X POST "http://localhost:8888/api/v1/order/list" \
  -H "Authorization: Bearer ${TOKEN}" \
  -H "Content-Type: application/json" \
  -d '{"page": 1, "page_size": 10}'
```

### 场景2：KOL 管理订单

```bash
# 1. KOL 登录
TOKEN="..."

# 2. 查看收到的订单
curl -X POST "http://localhost:8888/api/v1/order/kol/list" \
  -H "Authorization: Bearer ${TOKEN}" \
  -H "Content-Type: application/json" \
  -d '{"status": "pending"}'

# 3. 确认接单
curl -X POST "http://localhost:8888/api/v1/order/status/update" \
  -H "Authorization: Bearer ${TOKEN}" \
  -H "Content-Type: application/json" \
  -d '{
    "order_id": "ORD7234567890123456789",
    "status": "confirmed"
  }'

# 4. 开始执行
curl -X POST "http://localhost:8888/api/v1/order/status/update" \
  -H "Authorization: Bearer ${TOKEN}" \
  -H "Content-Type: application/json" \
  -d '{
    "order_id": "ORD7234567890123456789",
    "status": "in_progress"
  }'

# 5. 完成订单
curl -X POST "http://localhost:8888/api/v1/order/status/update" \
  -H "Authorization: Bearer ${TOKEN}" \
  -H "Content-Type: application/json" \
  -d '{
    "order_id": "ORD7234567890123456789",
    "status": "completed"
  }'
```

## 测试

运行测试脚本：

```bash
bash test_order.sh
```

## 技术架构

### 分层架构

```
Controller (Handler)
    ↓
Service (Business Logic)
    ↓
Repository (Data Access)
    ↓
Database
```

### 文件结构

```
biz/
├── handler/order/order/order_service.go  # Controller层
├── service/order/order_service.go         # Service层
├── dal/mysql/order.go                     # Repository层
├── model/order/order/order.go             # API模型
└── utils/id_generator.go                  # ID生成器
```

## 最佳实践

1. **订单快照**：保存 Plan 的快照数据，避免 Plan 修改后影响历史订单
2. **状态机**：严格的状态转换控制，避免非法状态变更
3. **权限控制**：每个操作都验证用户权限
4. **唯一ID**：使用 Snowflake 算法生成分布式唯一ID
5. **软删除**：使用 deleted_at 字段实现软删除
6. **时间戳**：记录关键状态变更的时间点

## 未来优化

1. **支付集成**：集成支付系统
2. **消息通知**：订单状态变更时通知用户和 KOL
3. **评价系统**：订单完成后的评价功能
4. **争议处理**：订单纠纷处理流程
5. **订单统计**：KOL 和用户的订单统计分析

