# Campaign API 接口文档

## 概述

Campaign（广告活动）功能允许用户创建和管理 TikTok 广告活动。用户可以配置推广目标、受众定位、预算、时间等参数，并上传创意素材。

**Base URL**: `/`

**认证方式**: JWT Token (放在请求头的 `Authorization` 字段)

**请求方法**: 所有接口均使用 `POST` 方法

**Content-Type**: `application/json`

---

## 认证说明

所有 Campaign 接口都需要 JWT 认证。请在请求头中添加：

```
Authorization: Bearer {your_jwt_token}
```

接口会自动从 JWT Token 中提取 `user_id` 和 `team_id`。

---

## 接口列表

### 普通用户接口

1. [创建 Campaign](#1-创建-campaign)
2. [更新 Campaign](#2-更新-campaign)
3. [更新 Campaign 状态](#3-更新-campaign-状态)
4. [获取 Campaign 列表](#4-获取-campaign-列表)
5. [获取 Campaign 详情](#5-获取-campaign-详情)

### 管理员接口

6. [管理员获取所有 Campaign](#6-管理员获取所有-campaign)
7. [管理员更新 Campaign 状态](#7-管理员更新-campaign-状态)

---

## 接口详情

### 1. 创建 Campaign

创建一个新的广告活动。

**URL**: `/campaign/create`

**Method**: `POST`

**权限**: 普通用户、管理员

**请求参数**:

| 参数名 | 类型 | 必填 | 说明 |
|--------|------|------|------|
| campaign_name | string | 是 | 活动名称 |
| promotion_objective | string | 是 | 推广目标：`awareness`(品牌认知)、`consideration`(受众意向)、`conversion`(行为转化) |
| optimization_goal | string | 是 | 优化目标，根据 promotion_objective 不同而不同：<br>- awareness: `reach`<br>- consideration: `website`、`app`<br>- conversion: `app_promotion`、`lead_generation` |
| location | array[int64] | 否 | 地区列表（数据字典ID数组） |
| age | int64 | 否 | 年龄段（数据字典ID） |
| gender | int64 | 否 | 性别（数据字典ID） |
| languages | array[int64] | 否 | 语言列表（数据字典ID数组） |
| spending_power | int64 | 否 | 消费能力（数据字典ID） |
| operating_system | int64 | 否 | 操作系统（数据字典ID） |
| os_versions | array[int64] | 否 | 系统版本列表（数据字典ID数组） |
| device_models | array[int64] | 否 | 设备品牌列表（数据字典ID数组） |
| connection_type | int64 | 否 | 网络情况（数据字典ID） |
| device_price_type | int32 | 是 | 设备价格类型：`0`-any, `1`-specific range |
| device_price_min | float64 | 否 | 设备价格最小值（当 device_price_type=1 时需要） |
| device_price_max | float64 | 否 | 设备价格最大值（当 device_price_type=1 时需要） |
| planned_start_time | string | 是 | 计划开始时间（RFC3339格式：`2024-01-01T00:00:00Z`） |
| planned_end_time | string | 是 | 计划结束时间（RFC3339格式） |
| time_zone | int64 | 否 | 时区（数据字典ID） |
| dayparting_type | int32 | 是 | 分时段类型：`0`-全天，`1`-特定时段 |
| dayparting_schedule | string | 否 | 特定时段配置（JSON格式，当 dayparting_type=1 时需要） |
| frequency_cap_type | int32 | 是 | 频次上限类型：`0`-每七天不超过三次，`1`-每天不超过一次，`2`-自定义 |
| frequency_cap_times | int32 | 否 | 自定义频次-次数（当 frequency_cap_type=2 时需要） |
| frequency_cap_days | int32 | 否 | 自定义频次-天数（当 frequency_cap_type=2 时需要） |
| budget_type | int32 | 是 | 预算类型：`0`-每日预算，`1`-总预算 |
| budget_amount | float64 | 是 | 预算金额 |
| website | string | 否 | 网站链接 |
| ios_download_url | string | 否 | iOS下载链接 |
| android_download_url | string | 否 | Android下载链接 |
| attachment_urls | array[string] | 否 | 附件URL列表（用户上传的创意素材） |

**请求示例**:

```json
{
  "campaign_name": "春季促销活动",
  "promotion_objective": "awareness",
  "optimization_goal": "reach",
  "location": [1, 2, 3],
  "age": 10,
  "gender": 20,
  "languages": [30, 31],
  "spending_power": 40,
  "operating_system": 50,
  "os_versions": [51, 52],
  "device_models": [60, 61, 62],
  "connection_type": 70,
  "device_price_type": 1,
  "device_price_min": 500.0,
  "device_price_max": 2000.0,
  "planned_start_time": "2024-03-01T00:00:00Z",
  "planned_end_time": "2024-03-31T23:59:59Z",
  "time_zone": 80,
  "dayparting_type": 0,
  "frequency_cap_type": 1,
  "budget_type": 0,
  "budget_amount": 1000.0,
  "website": "https://example.com",
  "attachment_urls": [
    "https://s3.example.com/file1.jpg",
    "https://s3.example.com/video1.mp4"
  ]
}
```

**响应参数**:

| 参数名 | 类型 | 说明 |
|--------|------|------|
| campaign | object | Campaign详细信息（参见 [Campaign 对象](#campaign-对象)） |
| base_resp | object | 基础响应 |
| base_resp.code | int32 | 状态码，0表示成功 |
| base_resp.message | string | 响应消息 |

**响应示例**:

```json
{
  "campaign": {
    "id": 1,
    "campaign_id": "CAMPAIGN_1709251200000_12345",
    "user_id": 100,
    "team_id": 10,
    "campaign_name": "春季促销活动",
    "promotion_objective": "awareness",
    "optimization_goal": "reach",
    "location": [1, 2, 3],
    "age": 10,
    "gender": 20,
    "languages": [30, 31],
    "spending_power": 40,
    "operating_system": 50,
    "os_versions": [51, 52],
    "device_models": [60, 61, 62],
    "connection_type": 70,
    "device_price_type": 1,
    "device_price_min": 500.0,
    "device_price_max": 2000.0,
    "planned_start_time": "2024-03-01T00:00:00Z",
    "planned_end_time": "2024-03-31T23:59:59Z",
    "time_zone": 80,
    "dayparting_type": 0,
    "frequency_cap_type": 1,
    "budget_type": 0,
    "budget_amount": 1000.0,
    "website": "https://example.com",
    "status": "pending",
    "attachments": [
      {
        "id": 1,
        "file_url": "https://s3.example.com/file1.jpg",
        "file_name": "file1.jpg",
        "file_type": "image/jpeg",
        "file_size": 102400,
        "created_at": "2024-03-01T10:00:00Z"
      },
      {
        "id": 2,
        "file_url": "https://s3.example.com/video1.mp4",
        "file_name": "video1.mp4",
        "file_type": "video/mp4",
        "file_size": 5242880,
        "created_at": "2024-03-01T10:00:01Z"
      }
    ],
    "created_at": "2024-03-01T10:00:00Z",
    "updated_at": "2024-03-01T10:00:00Z"
  },
  "base_resp": {
    "code": 0,
    "message": "Campaign created successfully"
  }
}
```

---

### 2. 更新 Campaign

更新 Campaign 信息。**注意：只能在 pending 状态下更新。**

**URL**: `/campaign/update`

**Method**: `POST`

**权限**: 普通用户、管理员（只能更新自己创建的）

**请求参数**:

| 参数名 | 类型 | 必填 | 说明 |
|--------|------|------|------|
| campaign_id | string | 是 | Campaign ID |
| campaign_name | string | 否 | 活动名称 |
| promotion_objective | string | 否 | 推广目标 |
| optimization_goal | string | 否 | 优化目标 |
| ... | ... | 否 | 其他字段同创建接口，均为可选 |

**请求示例**:

```json
{
  "campaign_id": "CAMPAIGN_1709251200000_12345",
  "campaign_name": "春季促销活动-更新",
  "budget_amount": 2000.0,
  "attachment_urls": [
    "https://s3.example.com/new_file.jpg"
  ]
}
```

**响应参数**: 同创建接口

**错误情况**:
- Campaign 不存在
- 无权限修改（不是创建者）
- Campaign 状态不是 pending（已启动的活动无法修改）

---

### 3. 更新 Campaign 状态

更新 Campaign 的运行状态。**普通用户只能在 active 和 paused 之间切换。**

**URL**: `/campaign/status`

**Method**: `POST`

**权限**: 普通用户、管理员

**请求参数**:

| 参数名 | 类型 | 必填 | 说明 |
|--------|------|------|------|
| campaign_id | string | 是 | Campaign ID |
| status | string | 是 | 目标状态：`active`(已启动)、`paused`(暂停) |

**请求示例**:

```json
{
  "campaign_id": "CAMPAIGN_1709251200000_12345",
  "status": "active"
}
```

**响应参数**:

```json
{
  "base_resp": {
    "code": 0,
    "message": "Campaign status updated successfully"
  }
}
```

**状态转换规则**:
- `pending` → 无法通过此接口启动（需要管理员审核）
- `active` ⇄ `paused`（普通用户可操作）
- `ended` → 终态，无法再次启动

---

### 4. 获取 Campaign 列表

获取当前团队的 Campaign 列表，支持搜索和筛选。

**URL**: `/campaign/list`

**Method**: `POST`

**权限**: 普通用户、管理员

**请求参数**:

| 参数名 | 类型 | 必填 | 说明 |
|--------|------|------|------|
| keyword | string | 否 | 搜索关键字（匹配活动名称） |
| status | string | 否 | 状态筛选：`pending`、`active`、`paused`、`ended` |
| promotion_objective | string | 否 | 推广目标筛选 |
| page | int32 | 否 | 页码，默认 1 |
| page_size | int32 | 否 | 每页数量，默认 10 |

**请求示例**:

```json
{
  "keyword": "促销",
  "status": "active",
  "page": 1,
  "page_size": 20
}
```

**响应参数**:

| 参数名 | 类型 | 说明 |
|--------|------|------|
| campaigns | array | Campaign列表 |
| page_info | object | 分页信息 |
| page_info.page | int32 | 当前页码 |
| page_info.page_size | int32 | 每页数量 |
| page_info.total | int64 | 总记录数 |
| page_info.total_pages | int32 | 总页数 |
| base_resp | object | 基础响应 |

**响应示例**:

```json
{
  "campaigns": [
    {
      "id": 1,
      "campaign_id": "CAMPAIGN_1709251200000_12345",
      "campaign_name": "春季促销活动",
      "status": "active",
      "budget_amount": 1000.0,
      "created_at": "2024-03-01T10:00:00Z",
      "updated_at": "2024-03-01T12:00:00Z"
      // ... 其他字段
    }
  ],
  "page_info": {
    "page": 1,
    "page_size": 20,
    "total": 100,
    "total_pages": 5
  },
  "base_resp": {
    "code": 0,
    "message": "Success"
  }
}
```

---

### 5. 获取 Campaign 详情

获取指定 Campaign 的完整信息。

**URL**: `/campaign/detail`

**Method**: `POST`

**权限**: 普通用户、管理员（只能查看自己团队的）

**请求参数**:

| 参数名 | 类型 | 必填 | 说明 |
|--------|------|------|------|
| campaign_id | string | 是 | Campaign ID |

**请求示例**:

```json
{
  "campaign_id": "CAMPAIGN_1709251200000_12345"
}
```

**响应参数**:

```json
{
  "campaign": {
    // 完整的 Campaign 对象，包含所有字段和附件列表
  },
  "base_resp": {
    "code": 0,
    "message": "Success"
  }
}
```

---

### 6. 管理员获取所有 Campaign

管理员查看系统中所有的 Campaign，支持多维度筛选。

**URL**: `/admin/campaign/list`

**Method**: `POST`

**权限**: 仅管理员

**请求参数**:

| 参数名 | 类型 | 必填 | 说明 |
|--------|------|------|------|
| keyword | string | 否 | 搜索关键字 |
| status | string | 否 | 状态筛选 |
| promotion_objective | string | 否 | 推广目标筛选 |
| user_id | int64 | 否 | 按用户筛选 |
| team_id | int64 | 否 | 按团队筛选 |
| page | int32 | 否 | 页码，默认 1 |
| page_size | int32 | 否 | 每页数量，默认 10 |

**请求示例**:

```json
{
  "status": "active",
  "user_id": 100,
  "page": 1,
  "page_size": 50
}
```

**响应参数**: 同普通用户列表接口

---

### 7. 管理员更新 Campaign 状态

管理员可以将 Campaign 设置为任何状态，包括启动和结束。

**URL**: `/admin/campaign/status`

**Method**: `POST`

**权限**: 仅管理员

**请求参数**:

| 参数名 | 类型 | 必填 | 说明 |
|--------|------|------|------|
| campaign_id | string | 是 | Campaign ID |
| status | string | 是 | 目标状态：`active`、`paused`、`ended` |

**请求示例**:

```json
{
  "campaign_id": "CAMPAIGN_1709251200000_12345",
  "status": "ended"
}
```

**响应参数**:

```json
{
  "base_resp": {
    "code": 0,
    "message": "Campaign status updated successfully"
  }
}
```

---

## 数据结构

### Campaign 对象

| 字段名 | 类型 | 说明 |
|--------|------|------|
| id | int64 | 自增ID |
| campaign_id | string | 业务唯一ID |
| user_id | int64 | 创建用户ID |
| team_id | int64 | 所属团队ID |
| campaign_name | string | 活动名称 |
| promotion_objective | string | 推广目标 |
| optimization_goal | string | 优化目标 |
| location | array[int64] | 地区列表（数据字典ID） |
| age | int64 | 年龄段（数据字典ID） |
| gender | int64 | 性别（数据字典ID） |
| languages | array[int64] | 语言列表（数据字典ID） |
| spending_power | int64 | 消费能力（数据字典ID） |
| operating_system | int64 | 操作系统（数据字典ID） |
| os_versions | array[int64] | 系统版本列表（数据字典ID） |
| device_models | array[int64] | 设备品牌列表（数据字典ID） |
| connection_type | int64 | 网络情况（数据字典ID） |
| device_price_type | int32 | 设备价格类型 |
| device_price_min | float64 | 设备价格最小值 |
| device_price_max | float64 | 设备价格最大值 |
| planned_start_time | string | 计划开始时间 |
| planned_end_time | string | 计划结束时间 |
| time_zone | int64 | 时区（数据字典ID） |
| dayparting_type | int32 | 分时段类型 |
| dayparting_schedule | string | 特定时段配置 |
| frequency_cap_type | int32 | 频次上限类型 |
| frequency_cap_times | int32 | 自定义频次-次数 |
| frequency_cap_days | int32 | 自定义频次-天数 |
| budget_type | int32 | 预算类型 |
| budget_amount | float64 | 预算金额 |
| website | string | 网站链接 |
| ios_download_url | string | iOS下载链接 |
| android_download_url | string | Android下载链接 |
| status | string | 状态 |
| attachments | array | 附件列表 |
| created_at | string | 创建时间 |
| updated_at | string | 更新时间 |

### Attachment 对象

| 字段名 | 类型 | 说明 |
|--------|------|------|
| id | int64 | 附件ID |
| file_url | string | 文件URL |
| file_name | string | 文件名 |
| file_type | string | 文件类型（MIME类型） |
| file_size | int64 | 文件大小（字节） |
| created_at | string | 创建时间 |

---

## 枚举值说明

### promotion_objective（推广目标）

| 值 | 说明 | 可选的 optimization_goal |
|----|------|--------------------------|
| awareness | 品牌认知 | reach |
| consideration | 受众意向 | website, app |
| conversion | 行为转化 | app_promotion, lead_generation |

### optimization_goal（优化目标）

| 值 | 说明 |
|----|------|
| reach | 触达最多人群 |
| website | 网站推广 |
| app | 应用推广 |
| app_promotion | 应用促销 |
| lead_generation | 线索生成 |

### status（状态）

| 值 | 说明 | 可操作 |
|----|------|--------|
| pending | 待启动 | 可修改、等待管理员审核启动 |
| active | 已启动 | 可暂停 |
| paused | 暂停 | 可重新启动 |
| ended | 已结束 | 终态，不可操作 |

### device_price_type（设备价格类型）

| 值 | 说明 |
|----|------|
| 0 | 不限（any） |
| 1 | 指定范围（specific range），需提供 min 和 max |

### budget_type（预算类型）

| 值 | 说明 |
|----|------|
| 0 | 每日预算 |
| 1 | 总预算 |

### dayparting_type（分时段类型）

| 值 | 说明 |
|----|------|
| 0 | 全天 |
| 1 | 特定时段（需提供 dayparting_schedule） |

### frequency_cap_type（频次上限类型）

| 值 | 说明 |
|----|------|
| 0 | 每七天不超过三次 |
| 1 | 每天不超过一次 |
| 2 | 自定义（需提供 frequency_cap_times 和 frequency_cap_days） |

---

## 数据字典说明

Campaign 中的很多字段引用了数据字典（Dictionary）的值。需要先通过数据字典接口获取可选项：

### 需要引用数据字典的字段

| 字段 | 字典类型建议 | 说明 |
|------|--------------|------|
| location | `campaign_location` | 地区列表 |
| age | `campaign_age` | 年龄段 |
| gender | `campaign_gender` | 性别（all/male/female） |
| languages | `campaign_language` | 语言列表 |
| spending_power | `campaign_spending_power` | 消费能力 |
| operating_system | `campaign_os` | 操作系统（Android/iOS） |
| os_versions | `campaign_os_version` | 系统版本 |
| device_models | `campaign_device_model` | 设备品牌 |
| connection_type | `campaign_connection` | 网络类型（WiFi/4G/5G等） |
| time_zone | `campaign_timezone` | 时区 |

**使用方式**:
1. 调用数据字典接口获取可选项及其ID
2. 在创建/更新 Campaign 时传入对应的字典项ID

---

## 错误码说明

| 错误码 | 说明 |
|--------|------|
| 0 | 成功 |
| 400 | 参数错误 |
| 401 | 未认证 |
| 403 | 无权限 |
| 404 | Campaign不存在 |
| 500 | 服务器错误 |

**常见错误消息**:
- `User not authenticated`: 未登录或token失效
- `Team not found`: 团队不存在或用户不属于该团队
- `campaign not found`: Campaign不存在
- `permission denied`: 无权限操作此Campaign
- `only pending campaigns can be updated`: 只能更新待启动状态的Campaign
- `invalid status transition`: 无效的状态转换
- `invalid promotion_objective`: 无效的推广目标
- `invalid optimization_goal for promotion_objective xxx`: 优化目标与推广目标不匹配

---

## 业务流程

### 创建 Campaign 的完整流程

1. **准备数据字典**
   - 调用数据字典接口获取各类配置项
   
2. **上传素材文件**
   - 调用文件上传接口上传创意素材
   - 获取文件URL列表

3. **创建 Campaign**
   - 调用 `/campaign/create` 接口
   - 传入完整的配置和素材URL
   - Campaign 状态为 `pending`

4. **等待审核（可选）**
   - 管理员审核通过后将状态改为 `active`
   
5. **管理运行状态**
   - 使用 `/campaign/status` 接口暂停/重启
   - 查看 `/campaign/list` 接口监控所有活动

### 更新 Campaign 的注意事项

⚠️ **重要限制**:
- 只有 `pending` 状态的 Campaign 可以更新
- 一旦 Campaign 启动（status = `active`），无法再修改配置
- 如需修改已启动的活动，建议创建新的 Campaign

### 状态转换图

```
pending ──[admin启动]──> active ──[用户暂停]──> paused
                            │                      │
                            │                      │
                        [admin结束]          [用户重启]
                            │                      │
                            ↓                      ↓
                          ended                 active
```

---

## 附件管理

### 支持的文件类型

| 类型 | 扩展名 | MIME类型 |
|------|--------|----------|
| 图片 | .jpg, .jpeg, .png, .gif | image/jpeg, image/png, image/gif |
| 视频 | .mp4, .mov, .avi | video/mp4, video/quicktime, video/x-msvideo |
| 文档 | .pdf, .doc, .docx | application/pdf, application/msword |

### 上传流程

1. 调用上传接口（`/upload/file`）上传文件
2. 获取返回的文件URL
3. 在创建/更新 Campaign 时将URL放入 `attachment_urls` 数组

---

## 开发建议

### 前端实现建议

1. **表单验证**
   - 验证 promotion_objective 和 optimization_goal 的组合
   - 根据类型显示/隐藏相关字段（如设备价格范围、自定义频次等）
   - 时间选择器使用 RFC3339 格式

2. **数据字典缓存**
   - 在应用启动时获取所有数据字典并缓存
   - 避免每次创建Campaign都请求数据字典

3. **文件上传**
   - 支持拖拽上传
   - 显示上传进度
   - 预览已上传的文件

4. **状态管理**
   - 根据状态显示不同的操作按钮
   - 已启动的Campaign禁用编辑功能

5. **权限控制**
   - 根据用户角色显示/隐藏管理员功能
   - 显示操作权限提示

---

## 测试用例

### 测试数据字典ID（示例）

请使用实际的数据字典ID替换以下值：

```javascript
const testData = {
  location: [1, 2, 3],        // 美国、中国、日本
  age: 10,                    // 18-24岁
  gender: 20,                 // 全部
  languages: [30, 31],        // 英语、中文
  spending_power: 40,         // 中等
  operating_system: 50,       // iOS
  os_versions: [51],          // iOS 15+
  device_models: [60, 61],    // iPhone、Samsung
  connection_type: 70,        // WiFi
  time_zone: 80              // UTC+8
};
```

### cURL 示例

```bash
# 创建 Campaign
curl -X POST http://localhost:8080/campaign/create \
  -H "Authorization: Bearer YOUR_JWT_TOKEN" \
  -H "Content-Type: application/json" \
  -d '{
    "campaign_name": "测试活动",
    "promotion_objective": "awareness",
    "optimization_goal": "reach",
    "device_price_type": 0,
    "planned_start_time": "2024-03-01T00:00:00Z",
    "planned_end_time": "2024-03-31T23:59:59Z",
    "dayparting_type": 0,
    "frequency_cap_type": 1,
    "budget_type": 0,
    "budget_amount": 1000.0
  }'

# 获取列表
curl -X POST http://localhost:8080/campaign/list \
  -H "Authorization: Bearer YOUR_JWT_TOKEN" \
  -H "Content-Type: application/json" \
  -d '{
    "page": 1,
    "page_size": 20
  }'
```

---

## 常见问题 FAQ

### Q1: 为什么已启动的 Campaign 无法修改？

A: 这是为了保证广告投放的一致性和数据准确性。一旦Campaign启动，其配置就会锁定。如需修改，建议创建新的Campaign。

### Q2: 如何理解数据字典ID？

A: 数据字典是配置项的统一管理方式。例如"年龄段"可能有多个选项（18-24、25-34等），每个选项都有唯一的ID。创建Campaign时传入ID，系统会自动关联。

### Q3: 附件URL必须是S3链接吗？

A: 不是必须的，但建议使用系统的上传接口，这样可以确保文件的持久化和访问权限。

### Q4: 普通用户能否看到其他团队的Campaign？

A: 不能。普通用户只能看到自己所在团队的Campaign。只有管理员可以查看所有Campaign。

### Q5: Campaign 的状态转换规则是什么？

A: 
- 普通用户只能在 `active` 和 `paused` 之间切换
- 管理员可以设置为任何状态
- `ended` 是终态，无法再次启动

---

## 更新日志

### v1.0.0 (2024-03-01)
- 初始版本发布
- 支持完整的Campaign创建和管理功能
- 支持附件上传
- 支持数据字典集成
- 支持权限控制

---

## 联系方式

如有问题，请联系：
- 技术支持：tech@example.com
- API文档：https://docs.example.com

---

**文档生成时间**: 2024-03-01
**API版本**: v1.0.0

