# KOL Order API 接口文档

## 目录
- [概述](#概述)
- [订单状态说明](#订单状态说明)
- [订单状态转换流程](#订单状态转换流程)
- [API 接口列表](#api-接口列表)
  - [用户端接口](#用户端接口)
  - [KOL 端接口](#kol-端接口)
- [数据模型](#数据模型)
- [典型业务流程](#典型业务流程)
- [错误处理](#错误处理)

---

## 概述

KOL Order API 提供了完整的 KOL 订单管理功能，包括订单创建、查询、状态管理等。

**Base URL**: `/api/v1/kol-order`

**认证方式**: 所有接口均需要在请求头中携带 JWT Token
```
Authorization: Bearer <token>
```

**请求方式**: 所有接口均使用 **POST** 方法

**Content-Type**: `application/json`

---

## 订单状态说明

| 状态 | 状态码 | 说明 | 参与者 |
|-----|--------|------|--------|
| 待支付 | `pending_payment` | 订单已创建，等待用户支付 | 用户 |
| 待确认 | `pending` | 用户已支付，等待 KOL 确认接单 | KOL |
| 已确认 | `confirmed` | KOL 已确认接单，自动创建会话 | 双方 |
| 进行中 | `in_progress` | 订单执行中 | KOL |
| 已完成 | `completed` | 订单已完成 | 双方 |
| 已取消 | `cancelled` | 订单已取消（用户或 KOL 取消） | 双方 |
| 已退款 | `refunded` | 订单已退款 | 系统/管理员 |

---

## 订单状态转换流程

### 正常流程
```
pending_payment (待支付)
    ↓ [用户支付完成]
pending (待确认)
    ↓ [KOL 确认接单]
confirmed (已确认) → [自动创建会话]
    ↓ [KOL 开始制作]
in_progress (进行中)
    ↓ [KOL 完成交付]
completed (已完成)
```

### KOL 可操作的状态转换

**从 pending（待确认）:**
- → `confirmed`: KOL 确认接单
- → `cancelled`: KOL 拒绝接单（需提供拒绝原因）

**从 confirmed（已确认）:**
- → `in_progress`: KOL 开始制作
- → `cancelled`: KOL 取消订单（需提供取消原因）

**从 in_progress（进行中）:**
- → `completed`: KOL 完成交付
- → `cancelled`: KOL 取消订单（需提供取消原因）

### 用户可操作的状态转换

用户可以在以下状态取消订单：
- `pending_payment`（待支付）
- `pending`（待确认）
- `confirmed`（已确认）
- `in_progress`（进行中）

> **注意**: 不能取消状态为 `completed`、`cancelled`、`refunded` 的订单

---

## API 接口列表

### 用户端接口

#### 1. 创建 KOL 订单

**接口**: `POST /api/v1/kol-order/create`

**描述**: 用户创建一个新的 KOL 订单

**权限**: 需要登录

**请求参数**:
```json
{
  "kol_id": 123,                          // 必填，KOL ID
  "plan_id": 456,                         // 必填，报价计划 ID
  "title": "产品推广视频制作",               // 必填，订单标题，最长 200 字符
  "requirement_description": "需要一个 60 秒的产品介绍视频...",  // 必填，合作需求描述
  "video_type": "Product Review",         // 必填，视频类型，最长 100 字符
  "video_duration": 60,                   // 必填，视频预计时长（秒）
  "target_audience": "18-35岁科技爱好者",  // 必填，目标受众，最长 500 字符
  "expected_delivery_date": "2024-12-31", // 必填，期望交付日期（格式：YYYY-MM-DD）
  "additional_requirements": "希望加入品牌 logo...",  // 可选，额外要求
  "team_id": 789                          // 可选，团队 ID（如果是团队下单）
}
```

**响应**:
```json
{
  "base_resp": {
    "code": 0,
    "msg": "success"
  },
  "order_id": "KORD_1703123456_ABC123"    // 订单 ID
}
```

**业务逻辑**:
1. 验证 KOL 是否存在且状态为 `approved`（已审核通过）
2. 验证 Plan 是否存在且属于该 KOL
3. 如果指定了 `team_id`，验证用户是否属于该团队（待实现）
4. 生成唯一订单 ID（格式：`KORD_{timestamp}_{random}`）
5. 保存 Plan 的快照信息（title、description、price、type）
6. 初始订单状态为 `pending_payment`（待支付）

**错误情况**:
- KOL 不存在: `"KOL 不存在"`
- KOL 未审核通过: `"该 KOL 尚未通过审核"`
- Plan 不存在: `"报价计划不存在"`
- Plan 不属于该 KOL: `"报价计划不属于该 KOL"`

---

#### 2. 确认订单支付

**接口**: `POST /api/v1/kol-order/payment/confirm`

**描述**: 用户支付完成后，调用此接口确认支付

**权限**: 需要登录，且必须是订单创建者

**请求参数**:
```json
{
  "order_id": "KORD_1703123456_ABC123"   // 必填，订单 ID
}
```

**响应**:
```json
{
  "base_resp": {
    "code": 0,
    "msg": "success"
  }
}
```

**业务逻辑**:
1. 验证订单是否存在
2. 验证用户是否是订单创建者
3. 验证订单状态是否为 `pending_payment`
4. 更新订单状态为 `pending`（待 KOL 确认）

**错误情况**:
- 订单不存在: `"订单不存在"`
- 无权操作: `"无权操作该订单"`
- 状态错误: `"订单状态不是待支付，无法确认支付"`

---

#### 3. 获取订单详情

**接口**: `POST /api/v1/kol-order/detail`

**描述**: 获取指定订单的详细信息

**权限**: 需要登录，且必须是订单创建者或订单相关的 KOL

**请求参数**:
```json
{
  "order_id": "KORD_1703123456_ABC123"   // 必填，订单 ID
}
```

**响应**:
```json
{
  "base_resp": {
    "code": 0,
    "msg": "success"
  },
  "order": {
    // KolOrderInfo 对象，详见"数据模型"章节
  }
}
```

**业务逻辑**:
1. 获取订单信息（包含 KOL 和用户信息）
2. 权限验证：只有订单创建者、订单相关的 KOL 可以查看
3. 返回完整订单信息

**错误情况**:
- 订单不存在: `"订单不存在"`
- 无权查看: `"无权查看该订单"`

---

#### 4. 获取用户订单列表

**接口**: `POST /api/v1/kol-order/user/list`

**描述**: 获取当前登录用户创建的所有订单列表

**权限**: 需要登录

**请求参数**:
```json
{
  "status": "pending",           // 可选，订单状态筛选
  "keyword": "产品推广",          // 可选，关键词搜索（搜索订单标题、订单ID、KOL名称）
  "kol_id": 123,                 // 可选，筛选指定 KOL 的订单
  "team_id": 789,                // 可选，筛选指定团队的订单
  "page": 1,                     // 可选，页码，默认 1
  "page_size": 10                // 可选，每页数量，默认 10
}
```

**响应**:
```json
{
  "base_resp": {
    "code": 0,
    "msg": "success"
  },
  "orders": [
    // KolOrderInfo 对象数组，详见"数据模型"章节
  ],
  "total": 100                   // 总记录数
}
```

**业务逻辑**:
1. 获取当前用户的所有订单
2. 支持按状态筛选
3. 支持按 KOL ID 筛选
4. 支持关键词模糊搜索（订单标题、订单ID、KOL名称）
5. 按创建时间倒序排列

---

#### 5. 取消订单

**接口**: `POST /api/v1/kol-order/cancel`

**描述**: 用户取消自己创建的订单

**权限**: 需要登录，且必须是订单创建者

**请求参数**:
```json
{
  "order_id": "KORD_1703123456_ABC123",  // 必填，订单 ID
  "reason": "项目暂时搁置"                // 必填，取消原因
}
```

**响应**:
```json
{
  "base_resp": {
    "code": 0,
    "msg": "success"
  }
}
```

**业务逻辑**:
1. 验证订单是否存在
2. 验证用户是否是订单创建者
3. 验证订单状态是否可以取消（不能取消已完成、已取消、已退款的订单）
4. 更新订单状态为 `cancelled`
5. 记录取消原因和取消时间

**错误情况**:
- 订单不存在: `"订单不存在"`
- 无权操作: `"无权操作该订单"`
- 状态错误: `"该订单无法取消"`

---

### KOL 端接口

#### 6. 获取 KOL 收到的订单列表

**接口**: `POST /api/v1/kol-order/kol/list`

**描述**: KOL 查看收到的所有订单

**权限**: 需要登录，且必须是 KOL

**请求参数**:
```json
{
  "status": "pending",           // 可选，订单状态筛选
  "keyword": "产品推广",          // 可选，关键词搜索（搜索订单标题、订单ID、用户名称）
  "page": 1,                     // 可选，页码，默认 1
  "page_size": 10                // 可选，每页数量，默认 10
}
```

**响应**:
```json
{
  "base_resp": {
    "code": 0,
    "msg": "success"
  },
  "orders": [
    // KolOrderInfo 对象数组，详见"数据模型"章节
  ],
  "total": 100                   // 总记录数
}
```

**业务逻辑**:
1. 验证用户是否是 KOL
2. 获取该 KOL 收到的所有订单
3. 支持按状态筛选
4. 支持关键词模糊搜索（订单标题、订单ID、用户名称）
5. 按创建时间倒序排列

**错误情况**:
- 不是 KOL: `"您还不是 KOL，无法查看订单"`

---

#### 7. 更新订单状态

**接口**: `POST /api/v1/kol-order/status/update`

**描述**: KOL 更新订单状态（确认、拒绝、进行中、完成等）

**权限**: 需要登录，且必须是该订单相关的 KOL

**请求参数**:
```json
{
  "order_id": "KORD_1703123456_ABC123",  // 必填，订单 ID
  "status": "confirmed",                 // 必填，新状态
  "reject_reason": "档期已满"            // 可选，拒绝/取消原因（当状态为 cancelled 时需要提供）
}
```

**可用状态值**:
- `confirmed`: 确认接单
- `in_progress`: 进行中
- `completed`: 已完成
- `cancelled`: 取消订单

**响应**:
```json
{
  "base_resp": {
    "code": 0,
    "msg": "success"
  }
}
```

**业务逻辑**:
1. 验证用户是否是 KOL
2. 验证订单是否属于该 KOL
3. 验证状态转换是否合法（见"订单状态转换流程"）
4. 更新订单状态和相关时间戳
5. **特殊处理**: 当状态变更为 `confirmed` 时，自动创建会话
   - 会话标题: `"KOL订单: {订单标题}"`
   - 会话类型: `kol_order`
   - 参与者: 订单创建用户 + KOL 用户

**状态转换规则**:
```
pending (待确认) → confirmed, cancelled
confirmed (已确认) → in_progress, cancelled
in_progress (进行中) → completed, cancelled
```

**错误情况**:
- 不是 KOL: `"您还不是 KOL，无法操作订单"`
- 订单不存在: `"订单不存在"`
- 无权操作: `"无权操作该订单"`
- 状态转换非法: `"当前订单状态无法变更"` 或 `"不允许从 {当前状态} 状态转换到 {目标状态} 状态"`

---

## 数据模型

### KolOrderInfo

完整的订单信息对象，包含订单详情、用户信息、KOL 信息等。

```json
{
  "order_id": "KORD_1703123456_ABC123",           // 订单ID（格式：KORD_{timestamp}_{random}）
  "user_id": 1001,                                // 下单用户 ID
  "user_nickname": "张三",                         // 用户昵称
  "team_id": 789,                                 // 团队 ID（可选）
  "team_name": "营销团队A",                        // 团队名称（可选）
  "kol_id": 123,                                  // KOL ID
  "kol_display_name": "李四",                      // KOL 显示名称
  "kol_avatar_url": "https://example.com/avatar.jpg",  // KOL 头像 URL
  "plan_id": 456,                                 // Plan ID
  "plan_title": "标准推广套餐",                     // Plan 标题（快照）
  "plan_description": "包含 1 个 60 秒视频...",    // Plan 描述（快照）
  "plan_price": 999.99,                           // Plan 价格（美元，快照）
  "plan_type": "standard",                        // Plan 类型（快照）: basic/standard/premium
  "title": "产品推广视频制作",                      // 订单标题
  "requirement_description": "需要一个 60 秒的产品介绍视频...",  // 合作需求描述
  "video_type": "Product Review",                 // 视频类型
  "video_duration": 60,                           // 视频预计时长（秒）
  "target_audience": "18-35岁科技爱好者",          // 目标受众
  "expected_delivery_date": "2024-12-31",         // 期望交付日期（YYYY-MM-DD）
  "additional_requirements": "希望加入品牌 logo...",  // 额外要求（可选）
  "status": "pending",                            // 订单状态
  "reject_reason": "档期已满",                     // 拒绝/取消原因（可选）
  "confirmed_at": "2024-01-15T10:30:00Z",        // 确认时间（可选，ISO 8601 格式）
  "completed_at": "2024-01-20T15:45:00Z",        // 完成时间（可选）
  "cancelled_at": null,                           // 取消时间（可选）
  "created_at": "2024-01-10T08:00:00Z",          // 创建时间（ISO 8601 格式）
  "updated_at": "2024-01-15T10:30:00Z"           // 更新时间（ISO 8601 格式）
}
```

**字段说明**:

| 字段 | 类型 | 说明 |
|-----|------|------|
| `order_id` | string | 订单唯一ID，格式：`KORD_{timestamp}_{random}` |
| `user_id` | int64 | 下单用户ID |
| `user_nickname` | string | 下单用户昵称（关联查询） |
| `team_id` | int64? | 下单团队ID（可选） |
| `team_name` | string? | 团队名称（可选，关联查询） |
| `kol_id` | int64 | KOL ID |
| `kol_display_name` | string | KOL 显示名称（关联查询） |
| `kol_avatar_url` | string | KOL 头像URL（关联查询） |
| `plan_id` | int64 | 报价计划ID |
| `plan_title` | string | Plan 标题（快照，创建订单时保存） |
| `plan_description` | string | Plan 描述（快照） |
| `plan_price` | double | Plan 价格（美元，快照） |
| `plan_type` | string | Plan 类型（快照）: `basic`, `standard`, `premium` |
| `title` | string | 订单标题 |
| `requirement_description` | string | 合作需求描述 |
| `video_type` | string | 视频类型 |
| `video_duration` | int32 | 视频预计时长（秒） |
| `target_audience` | string | 目标受众 |
| `expected_delivery_date` | string | 期望交付日期（格式：YYYY-MM-DD） |
| `additional_requirements` | string? | 额外要求（可选） |
| `status` | string | 订单状态（见"订单状态说明"） |
| `reject_reason` | string? | 拒绝/取消原因（可选） |
| `confirmed_at` | string? | 确认时间（可选，ISO 8601 格式） |
| `completed_at` | string? | 完成时间（可选） |
| `cancelled_at` | string? | 取消时间（可选） |
| `created_at` | string | 创建时间（ISO 8601 格式） |
| `updated_at` | string | 更新时间（ISO 8601 格式） |

---

## 典型业务流程

### 流程 1: 用户下单并完成订单

```
1. 用户浏览 KOL 主页，选择一个 Plan
   GET /api/v1/kol/detail { "kol_id": 123 }

2. 用户创建订单
   POST /api/v1/kol-order/create
   {
     "kol_id": 123,
     "plan_id": 456,
     "title": "产品推广视频",
     "requirement_description": "...",
     "video_type": "Product Review",
     "video_duration": 60,
     "target_audience": "18-35岁",
     "expected_delivery_date": "2024-12-31"
   }
   → 返回 order_id: "KORD_xxx"
   → 订单状态: pending_payment

3. 用户进行支付（对接支付接口）
   ...

4. 支付成功后，确认支付
   POST /api/v1/kol-order/payment/confirm
   { "order_id": "KORD_xxx" }
   → 订单状态: pending (等待 KOL 确认)

5. KOL 查看订单列表
   POST /api/v1/kol-order/kol/list
   { "status": "pending" }

6. KOL 确认接单
   POST /api/v1/kol-order/status/update
   { "order_id": "KORD_xxx", "status": "confirmed" }
   → 订单状态: confirmed
   → 系统自动创建会话（双方可以开始沟通）

7. KOL 开始制作
   POST /api/v1/kol-order/status/update
   { "order_id": "KORD_xxx", "status": "in_progress" }
   → 订单状态: in_progress

8. KOL 完成交付
   POST /api/v1/kol-order/status/update
   { "order_id": "KORD_xxx", "status": "completed" }
   → 订单状态: completed
```

### 流程 2: KOL 拒绝订单

```
1. KOL 查看待确认订单
   POST /api/v1/kol-order/kol/list
   { "status": "pending" }

2. KOL 拒绝订单
   POST /api/v1/kol-order/status/update
   {
     "order_id": "KORD_xxx",
     "status": "cancelled",
     "reject_reason": "档期已满，无法接单"
   }
   → 订单状态: cancelled
```

### 流程 3: 用户取消订单

```
1. 用户查看自己的订单列表
   POST /api/v1/kol-order/user/list

2. 用户取消订单
   POST /api/v1/kol-order/cancel
   {
     "order_id": "KORD_xxx",
     "reason": "项目暂时搁置"
   }
   → 订单状态: cancelled
```

---

## 错误处理

### 统一错误响应格式

```json
{
  "base_resp": {
    "code": 40001,              // 错误码（非 0）
    "msg": "订单不存在"          // 错误消息
  }
}
```

### 常见错误码

| 错误码 | 说明 |
|-------|------|
| `0` | 成功 |
| `40001` | 请求参数错误 |
| `40101` | 未登录或 Token 无效 |
| `40301` | 无权访问该资源 |
| `40401` | 资源不存在（订单、KOL、Plan 等） |
| `50001` | 服务器内部错误 |

### 常见业务错误消息

| 错误消息 | 触发条件 |
|---------|---------|
| `"KOL 不存在"` | 指定的 KOL ID 不存在 |
| `"该 KOL 尚未通过审核"` | KOL 状态不是 `approved` |
| `"报价计划不存在"` | 指定的 Plan ID 不存在 |
| `"报价计划不属于该 KOL"` | Plan 和 KOL 不匹配 |
| `"订单不存在"` | 指定的订单 ID 不存在 |
| `"无权操作该订单"` | 用户不是订单创建者或相关 KOL |
| `"无权查看该订单"` | 用户不是订单创建者或相关 KOL |
| `"订单状态不是待支付，无法确认支付"` | 订单状态不符合支付确认条件 |
| `"该订单无法取消"` | 订单状态为 `completed`/`cancelled`/`refunded` |
| `"当前订单状态无法变更"` | 订单当前状态不支持任何状态转换 |
| `"不允许从 X 状态转换到 Y 状态"` | 状态转换不符合业务规则 |
| `"您还不是 KOL，无法查看订单"` | 非 KOL 用户访问 KOL 接口 |
| `"您还不是 KOL，无法操作订单"` | 非 KOL 用户尝试更新订单状态 |

---

## 注意事项

### 1. Plan 快照机制
- 订单创建时会保存 Plan 的快照（title、description、price、type）
- 即使 Plan 后续被修改或删除，订单中的价格和描述不会改变
- 这确保了订单的历史信息完整性

### 2. 会话自动创建
- 当 KOL 确认订单（状态变为 `confirmed`）时，系统会自动创建一个会话
- 会话标题: `"KOL订单: {订单标题}"`
- 会话类型: `kol_order`
- 参与者: 订单创建用户 + KOL 用户
- 如果会话创建失败，不影响订单状态更新（仅记录错误日志）

### 3. 团队订单
- 用户可以代表团队下单（传递 `team_id`）
- 团队成员验证功能待实现

### 4. 时间格式
- 所有日期时间字段均使用 **ISO 8601** 格式（例如：`2024-01-15T10:30:00Z`）
- `expected_delivery_date` 使用日期格式（例如：`2024-12-31`）

### 5. 分页参数
- `page`: 默认 1（从 1 开始）
- `page_size`: 默认 10，建议范围 1-100

### 6. 关键词搜索
- 用户端搜索: 订单标题、订单ID、KOL名称
- KOL 端搜索: 订单标题、订单ID、用户名称
- 搜索为模糊匹配（LIKE 查询）

---

## 更新日志

| 版本 | 日期 | 说明 |
|-----|------|------|
| 1.0.0 | 2024-01-10 | 初始版本，包含所有基础功能 |

---

## 联系方式

如有问题，请联系后端开发团队或查阅相关技术文档。

