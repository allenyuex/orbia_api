# Dashboard API 文档

## 概述

Dashboard API 提供了运营数据管理和展示功能，包括：
1. 优秀广告案例管理
2. 内容趋势管理
3. 平台数据统计管理
4. Dashboard 数据展示（普通用户接口）

## 接口列表

### 管理员接口（需要管理员权限）

#### 1. 优秀广告案例管理

##### 1.1 创建优秀案例

**接口地址**：`POST /api/v1/admin/dashboard/excellent-case/create`

**权限要求**：管理员

**请求参数**：
```json
{
  "video_url": "https://example.com/video.mp4",      // 必填，视频URL
  "cover_url": "https://example.com/cover.jpg",      // 必填，封面URL
  "title": "优秀广告案例标题",                        // 必填，案例标题
  "description": "案例描述信息",                      // 可选，案例描述
  "sort_order": 1                                    // 可选，排序序号（数字越小越靠前）
}
```

**响应示例**：
```json
{
  "base_resp": {
    "code": 0,
    "message": "success"
  },
  "id": 1                                            // 创建的案例ID
}
```

##### 1.2 更新优秀案例

**接口地址**：`POST /api/v1/admin/dashboard/excellent-case/update`

**权限要求**：管理员

**请求参数**：
```json
{
  "id": 1,                                           // 必填，案例ID
  "video_url": "https://example.com/video.mp4",      // 可选，视频URL
  "cover_url": "https://example.com/cover.jpg",      // 可选，封面URL
  "title": "优秀广告案例标题",                        // 可选，案例标题
  "description": "案例描述信息",                      // 可选，案例描述
  "sort_order": 1,                                   // 可选，排序序号
  "status": 1                                        // 可选，状态：1-启用，0-禁用
}
```

**响应示例**：
```json
{
  "base_resp": {
    "code": 0,
    "message": "success"
  }
}
```

##### 1.3 删除优秀案例

**接口地址**：`POST /api/v1/admin/dashboard/excellent-case/delete`

**权限要求**：管理员

**请求参数**：
```json
{
  "id": 1                                            // 必填，案例ID
}
```

**响应示例**：
```json
{
  "base_resp": {
    "code": 0,
    "message": "success"
  }
}
```

##### 1.4 获取优秀案例列表

**接口地址**：`POST /api/v1/admin/dashboard/excellent-case/list`

**权限要求**：管理员

**请求参数**：
```json
{
  "status": 1,                                       // 可选，状态筛选：1-启用，0-禁用
  "page": 1,                                         // 可选，页码，默认1
  "page_size": 10                                    // 可选，每页数量，默认10
}
```

**响应示例**：
```json
{
  "base_resp": {
    "code": 0,
    "message": "success"
  },
  "cases": [
    {
      "id": 1,
      "video_url": "https://example.com/video.mp4",
      "cover_url": "https://example.com/cover.jpg",
      "title": "优秀广告案例标题",
      "description": "案例描述信息",
      "sort_order": 1,
      "status": 1,
      "created_at": "2025-01-01T00:00:00Z",
      "updated_at": "2025-01-01T00:00:00Z"
    }
  ],
  "page_info": {
    "page": 1,
    "page_size": 10,
    "total": 100,
    "total_page": 10
  }
}
```

##### 1.5 获取优秀案例详情

**接口地址**：`POST /api/v1/admin/dashboard/excellent-case/:id`

**权限要求**：管理员

**URL参数**：
- `id`: 案例ID

**响应示例**：
```json
{
  "base_resp": {
    "code": 0,
    "message": "success"
  },
  "case_detail": {
    "id": 1,
    "video_url": "https://example.com/video.mp4",
    "cover_url": "https://example.com/cover.jpg",
    "title": "优秀广告案例标题",
    "description": "案例描述信息",
    "sort_order": 1,
    "status": 1,
    "created_at": "2025-01-01T00:00:00Z",
    "updated_at": "2025-01-01T00:00:00Z"
  }
}
```

#### 2. 内容趋势管理

##### 2.1 创建内容趋势

**接口地址**：`POST /api/v1/admin/dashboard/content-trend/create`

**权限要求**：管理员

**请求参数**：
```json
{
  "ranking": 1,                                      // 必填，排名（1,2,3,4,5...）
  "hot_keyword": "DeFi",                             // 必填，热点词
  "value_level": "high",                             // 必填，价值等级：low-低，medium-中，high-高
  "heat": 100000,                                    // 必填，热度值
  "growth_rate": 15.5                                // 必填，增长比例（百分比）
}
```

**响应示例**：
```json
{
  "base_resp": {
    "code": 0,
    "message": "success"
  },
  "id": 1                                            // 创建的趋势ID
}
```

##### 2.2 更新内容趋势

**接口地址**：`POST /api/v1/admin/dashboard/content-trend/update`

**权限要求**：管理员

**请求参数**：
```json
{
  "id": 1,                                           // 必填，趋势ID
  "ranking": 1,                                      // 可选，排名
  "hot_keyword": "DeFi",                             // 可选，热点词
  "value_level": "high",                             // 可选，价值等级
  "heat": 100000,                                    // 可选，热度值
  "growth_rate": 15.5,                               // 可选，增长比例
  "status": 1                                        // 可选，状态：1-启用，0-禁用
}
```

**响应示例**：
```json
{
  "base_resp": {
    "code": 0,
    "message": "success"
  }
}
```

##### 2.3 删除内容趋势

**接口地址**：`POST /api/v1/admin/dashboard/content-trend/delete`

**权限要求**：管理员

**请求参数**：
```json
{
  "id": 1                                            // 必填，趋势ID
}
```

**响应示例**：
```json
{
  "base_resp": {
    "code": 0,
    "message": "success"
  }
}
```

##### 2.4 获取内容趋势列表

**接口地址**：`POST /api/v1/admin/dashboard/content-trend/list`

**权限要求**：管理员

**请求参数**：
```json
{
  "status": 1,                                       // 可选，状态筛选：1-启用，0-禁用
  "page": 1,                                         // 可选，页码，默认1
  "page_size": 10                                    // 可选，每页数量，默认10
}
```

**响应示例**：
```json
{
  "base_resp": {
    "code": 0,
    "message": "success"
  },
  "trends": [
    {
      "id": 1,
      "ranking": 1,
      "hot_keyword": "DeFi",
      "value_level": "high",
      "heat": 100000,
      "growth_rate": 15.5,
      "status": 1,
      "created_at": "2025-01-01T00:00:00Z",
      "updated_at": "2025-01-01T00:00:00Z"
    }
  ],
  "page_info": {
    "page": 1,
    "page_size": 10,
    "total": 100,
    "total_page": 10
  }
}
```

##### 2.5 获取内容趋势详情

**接口地址**：`POST /api/v1/admin/dashboard/content-trend/:id`

**权限要求**：管理员

**URL参数**：
- `id`: 趋势ID

**响应示例**：
```json
{
  "base_resp": {
    "code": 0,
    "message": "success"
  },
  "trend_detail": {
    "id": 1,
    "ranking": 1,
    "hot_keyword": "DeFi",
    "value_level": "high",
    "heat": 100000,
    "growth_rate": 15.5,
    "status": 1,
    "created_at": "2025-01-01T00:00:00Z",
    "updated_at": "2025-01-01T00:00:00Z"
  }
}
```

#### 3. 平台数据统计管理

##### 3.1 更新平台数据

**接口地址**：`POST /api/v1/admin/dashboard/platform-stats/update`

**权限要求**：管理员

**说明**：此接口用于更新平台数据统计。如果数据不存在，会自动创建。所有字段都是可选的，只更新提供的字段。

**请求参数**：
```json
{
  "active_kols": 1000,                               // 可选，活跃的KOLs数量
  "total_coverage": 5000000,                         // 可选，总覆盖用户数
  "total_ad_impressions": 10000000,                  // 可选，累计广告曝光次数
  "total_transaction_amount": 1000000.00,            // 可选，平台总交易额（美元）
  "average_roi": 25.5,                               // 可选，平均ROI（百分比）
  "average_cpm": 15.8,                               // 可选，平均CPM
  "web3_brand_count": 500                            // 可选，合作Web3品牌数
}
```

**响应示例**：
```json
{
  "base_resp": {
    "code": 0,
    "message": "success"
  }
}
```

##### 3.2 获取平台数据

**接口地址**：`POST /api/v1/admin/dashboard/platform-stats`

**权限要求**：管理员

**请求参数**：
```json
{}
```

**响应示例**：
```json
{
  "base_resp": {
    "code": 0,
    "message": "success"
  },
  "stats": {
    "id": 1,
    "active_kols": 1000,
    "total_coverage": 5000000,
    "total_ad_impressions": 10000000,
    "total_transaction_amount": 1000000.00,
    "average_roi": 25.5,
    "average_cpm": 15.8,
    "web3_brand_count": 500,
    "created_at": "2025-01-01T00:00:00Z",
    "updated_at": "2025-01-01T00:00:00Z"
  }
}
```

### 普通用户接口（无需认证）

#### 4. 获取Dashboard数据

**接口地址**：`POST /api/v1/dashboard/data`

**权限要求**：无

**说明**：此接口返回所有启用的优秀广告案例、内容趋势和平台数据统计，供前端展示。

**请求参数**：
```json
{}
```

**响应示例**：
```json
{
  "base_resp": {
    "code": 0,
    "message": "success"
  },
  "excellent_cases": [
    {
      "id": 1,
      "video_url": "https://example.com/video.mp4",
      "cover_url": "https://example.com/cover.jpg",
      "title": "优秀广告案例标题",
      "description": "案例描述信息",
      "sort_order": 1,
      "status": 1,
      "created_at": "2025-01-01T00:00:00Z",
      "updated_at": "2025-01-01T00:00:00Z"
    }
  ],
  "content_trends": [
    {
      "id": 1,
      "ranking": 1,
      "hot_keyword": "DeFi",
      "value_level": "high",
      "heat": 100000,
      "growth_rate": 15.5,
      "status": 1,
      "created_at": "2025-01-01T00:00:00Z",
      "updated_at": "2025-01-01T00:00:00Z"
    }
  ],
  "platform_stats": {
    "id": 1,
    "active_kols": 1000,
    "total_coverage": 5000000,
    "total_ad_impressions": 10000000,
    "total_transaction_amount": 1000000.00,
    "average_roi": 25.5,
    "average_cpm": 15.8,
    "web3_brand_count": 500,
    "created_at": "2025-01-01T00:00:00Z",
    "updated_at": "2025-01-01T00:00:00Z"
  }
}
```

## 数据模型说明

### 1. 优秀广告案例（ExcellentCaseItem）

| 字段 | 类型 | 说明 |
|------|------|------|
| id | int64 | 案例ID |
| video_url | string | 视频URL |
| cover_url | string | 封面URL |
| title | string | 案例标题 |
| description | string | 案例描述 |
| sort_order | int32 | 排序序号（升序，数字越小越靠前） |
| status | int32 | 状态：1-启用，0-禁用 |
| created_at | string | 创建时间（RFC3339格式） |
| updated_at | string | 更新时间（RFC3339格式） |

### 2. 内容趋势（ContentTrendItem）

| 字段 | 类型 | 说明 |
|------|------|------|
| id | int64 | 趋势ID |
| ranking | int32 | 排名（1,2,3,4,5...） |
| hot_keyword | string | 热点词 |
| value_level | string | 价值等级：low-低，medium-中，high-高 |
| heat | int64 | 热度值 |
| growth_rate | double | 增长比例（百分比） |
| status | int32 | 状态：1-启用，0-禁用 |
| created_at | string | 创建时间（RFC3339格式） |
| updated_at | string | 更新时间（RFC3339格式） |

### 3. 平台数据统计（PlatformStatsData）

| 字段 | 类型 | 说明 |
|------|------|------|
| id | int64 | 统计ID |
| active_kols | int64 | 活跃的KOLs数量 |
| total_coverage | int64 | 总覆盖用户数 |
| total_ad_impressions | int64 | 累计广告曝光次数 |
| total_transaction_amount | double | 平台总交易额（美元） |
| average_roi | double | 平均ROI（百分比） |
| average_cpm | double | 平均CPM |
| web3_brand_count | int64 | 合作Web3品牌数 |
| created_at | string | 创建时间（RFC3339格式） |
| updated_at | string | 更新时间（RFC3339格式） |

## 错误码说明

| 错误码 | 说明 |
|--------|------|
| 0 | 成功 |
| 400 | 请求参数错误 |
| 401 | 未认证 |
| 403 | 无权限 |
| 404 | 资源不存在 |
| 500 | 服务器内部错误 |

## 注意事项

1. 所有管理员接口都需要在请求头中携带有效的管理员Token
2. 平台数据统计表只有一行数据，更新时会自动创建或更新这一行数据
3. 内容趋势的排名字段具有唯一性约束，同一个排名不能重复
4. 优秀案例和内容趋势的排序按照 sort_order 或 ranking 字段升序排列
5. Dashboard数据接口（/api/v1/dashboard/data）只返回启用状态（status=1）的数据
6. 所有接口都使用POST方法，参数都通过JSON格式传递

