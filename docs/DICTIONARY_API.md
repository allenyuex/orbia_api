# 数据字典管理 API 文档

## 概述

数据字典管理功能提供了一个灵活的配置系统，用于管理系统中的预定义数据，如国家、地区、分类等。支持无限层级的树形结构。

### 特性

- ✅ 支持无限层级的树形结构
- ✅ 软删除机制，保证数据一致性
- ✅ 字典编码唯一性验证（只允许大小写字母）
- ✅ 状态管理（启用/禁用）
- ✅ 支持图标URL配置
- ✅ 排序功能
- ✅ 完整的CRUD操作
- ✅ 仅管理员权限访问（除公开的树形查询接口外）

## 数据库设计

### 字典表 (orbia_dictionary)

| 字段 | 类型 | 说明 |
|------|------|------|
| id | BIGINT | 主键ID |
| code | VARCHAR(100) | 字典编码（唯一，只能大小写字母）|
| name | VARCHAR(100) | 字典名称 |
| description | TEXT | 字典描述 |
| status | TINYINT | 状态：1-启用，0-禁用 |
| created_at | TIMESTAMP | 创建时间 |
| updated_at | TIMESTAMP | 更新时间 |
| deleted_at | TIMESTAMP | 软删除时间 |

### 字典项表 (orbia_dictionary_item)

| 字段 | 类型 | 说明 |
|------|------|------|
| id | BIGINT | 主键ID |
| dictionary_id | BIGINT | 字典ID |
| parent_id | BIGINT | 父级ID（0表示根节点）|
| code | VARCHAR(100) | 字典项编码 |
| name | VARCHAR(200) | 字典项名称 |
| description | TEXT | 字典项描述 |
| icon_url | VARCHAR(500) | 图标URL（如国旗图标）|
| sort_order | INT | 排序序号（升序）|
| level | INT | 层级（1开始）|
| path | VARCHAR(1000) | 路径（如: 1/2/3）|
| status | TINYINT | 状态：1-启用，0-禁用 |
| created_at | TIMESTAMP | 创建时间 |
| updated_at | TIMESTAMP | 更新时间 |
| deleted_at | TIMESTAMP | 软删除时间 |

## API 接口

### 认证说明

除了 `GetDictionaryTree` 接口外，其他所有接口都需要管理员权限。

**请求头：**
```
Authorization: Bearer {admin_token}
```

---

## 字典管理接口

### 1. 创建字典

**接口地址：** `POST /api/v1/admin/dictionary/create`

**权限要求：** 管理员

**请求参数：**
```json
{
  "code": "COUNTRY",
  "name": "国家列表",
  "description": "全球国家和地区列表"
}
```

| 参数 | 类型 | 必填 | 说明 |
|------|------|------|------|
| code | string | 是 | 字典编码，只能包含大小写字母，全局唯一 |
| name | string | 是 | 字典名称 |
| description | string | 否 | 字典描述 |

**响应示例：**
```json
{
  "base_resp": {
    "code": 0,
    "message": "success"
  },
  "dictionary": {
    "id": 1,
    "code": "COUNTRY",
    "name": "国家列表",
    "description": "全球国家和地区列表",
    "status": 1,
    "created_at": "2025-10-26 16:00:00",
    "updated_at": "2025-10-26 16:00:00"
  }
}
```

---

### 2. 更新字典

**接口地址：** `POST /api/v1/admin/dictionary/update`

**权限要求：** 管理员

**请求参数：**
```json
{
  "id": 1,
  "name": "国家和地区",
  "description": "全球国家和地区完整列表",
  "status": 1
}
```

| 参数 | 类型 | 必填 | 说明 |
|------|------|------|------|
| id | int64 | 是 | 字典ID |
| name | string | 是 | 字典名称 |
| description | string | 否 | 字典描述 |
| status | int32 | 否 | 状态：1-启用，0-禁用 |

**注意：** 字典编码（code）创建后不能修改

**响应示例：**
```json
{
  "base_resp": {
    "code": 0,
    "message": "success"
  },
  "dictionary": {
    "id": 1,
    "code": "COUNTRY",
    "name": "国家和地区",
    "description": "全球国家和地区完整列表",
    "status": 1,
    "created_at": "2025-10-26 16:00:00",
    "updated_at": "2025-10-26 16:05:00"
  }
}
```

---

### 3. 删除字典

**接口地址：** `POST /api/v1/admin/dictionary/delete`

**权限要求：** 管理员

**请求参数：**
```json
{
  "id": 1
}
```

| 参数 | 类型 | 必填 | 说明 |
|------|------|------|------|
| id | int64 | 是 | 字典ID |

**注意：** 
- 删除为软删除
- 删除字典会同时删除该字典下的所有字典项

**响应示例：**
```json
{
  "base_resp": {
    "code": 0,
    "message": "success"
  }
}
```

---

### 4. 获取字典列表

**接口地址：** `POST /api/v1/admin/dictionary/list`

**权限要求：** 管理员

**请求参数：**
```json
{
  "keyword": "国家",
  "status": 1,
  "page": 1,
  "page_size": 10
}
```

| 参数 | 类型 | 必填 | 说明 |
|------|------|------|------|
| keyword | string | 否 | 搜索关键字（匹配编码或名称）|
| status | int32 | 否 | 状态筛选：1-启用，0-禁用 |
| page | int32 | 否 | 页码，默认1 |
| page_size | int32 | 否 | 每页数量，默认10 |

**响应示例：**
```json
{
  "base_resp": {
    "code": 0,
    "message": "success"
  },
  "dictionaries": [
    {
      "id": 1,
      "code": "COUNTRY",
      "name": "国家列表",
      "description": "全球国家和地区列表",
      "status": 1,
      "created_at": "2025-10-26 16:00:00",
      "updated_at": "2025-10-26 16:00:00"
    }
  ],
  "page_info": {
    "page": 1,
    "page_size": 10,
    "total": 1,
    "total_pages": 1
  }
}
```

---

### 5. 获取字典详情

**接口地址：** `POST /api/v1/admin/dictionary/:id`

**权限要求：** 管理员

**路径参数：**
- `id`: 字典ID

**请求参数：**
```json
{
  "id": 1
}
```

**响应示例：**
```json
{
  "base_resp": {
    "code": 0,
    "message": "success"
  },
  "dictionary": {
    "id": 1,
    "code": "COUNTRY",
    "name": "国家列表",
    "description": "全球国家和地区列表",
    "status": 1,
    "created_at": "2025-10-26 16:00:00",
    "updated_at": "2025-10-26 16:00:00"
  }
}
```

---

### 6. 根据编码获取字典

**接口地址：** `POST /api/v1/admin/dictionary/by-code`

**权限要求：** 管理员

**请求参数：**
```json
{
  "code": "COUNTRY"
}
```

| 参数 | 类型 | 必填 | 说明 |
|------|------|------|------|
| code | string | 是 | 字典编码 |

**响应示例：**
```json
{
  "base_resp": {
    "code": 0,
    "message": "success"
  },
  "dictionary": {
    "id": 1,
    "code": "COUNTRY",
    "name": "国家列表",
    "description": "全球国家和地区列表",
    "status": 1,
    "created_at": "2025-10-26 16:00:00",
    "updated_at": "2025-10-26 16:00:00"
  }
}
```

---

## 字典项管理接口

### 7. 创建字典项

**接口地址：** `POST /api/v1/admin/dictionary/item/create`

**权限要求：** 管理员

**请求参数：**
```json
{
  "dictionary_id": 1,
  "parent_id": 0,
  "code": "CN",
  "name": "中国",
  "description": "中华人民共和国",
  "icon_url": "https://example.com/flags/cn.png",
  "sort_order": 1
}
```

| 参数 | 类型 | 必填 | 说明 |
|------|------|------|------|
| dictionary_id | int64 | 是 | 字典ID |
| parent_id | int64 | 是 | 父级ID，0表示根节点 |
| code | string | 是 | 字典项编码 |
| name | string | 是 | 字典项名称 |
| description | string | 否 | 字典项描述 |
| icon_url | string | 否 | 图标URL（如国旗图标）|
| sort_order | int32 | 否 | 排序序号，默认0 |

**注意：** 
- 同一字典下同一父级下的编码不能重复
- 层级和路径会自动计算

**响应示例：**
```json
{
  "base_resp": {
    "code": 0,
    "message": "success"
  },
  "item": {
    "id": 1,
    "dictionary_id": 1,
    "parent_id": 0,
    "code": "CN",
    "name": "中国",
    "description": "中华人民共和国",
    "icon_url": "https://example.com/flags/cn.png",
    "sort_order": 1,
    "level": 1,
    "path": "1",
    "status": 1,
    "created_at": "2025-10-26 16:00:00",
    "updated_at": "2025-10-26 16:00:00"
  }
}
```

**创建子级示例：**
```json
{
  "dictionary_id": 1,
  "parent_id": 1,
  "code": "SH",
  "name": "上海",
  "description": "上海市",
  "icon_url": "",
  "sort_order": 1
}
```

---

### 8. 更新字典项

**接口地址：** `POST /api/v1/admin/dictionary/item/update`

**权限要求：** 管理员

**请求参数：**
```json
{
  "id": 1,
  "name": "中国大陆",
  "description": "中华人民共和国（大陆地区）",
  "icon_url": "https://example.com/flags/cn_new.png",
  "sort_order": 1,
  "status": 1
}
```

| 参数 | 类型 | 必填 | 说明 |
|------|------|------|------|
| id | int64 | 是 | 字典项ID |
| name | string | 否 | 字典项名称 |
| description | string | 否 | 字典项描述 |
| icon_url | string | 否 | 图标URL |
| sort_order | int32 | 否 | 排序序号 |
| status | int32 | 否 | 状态：1-启用，0-禁用 |

**注意：** 
- 字典项编码、父级ID、字典ID创建后不能修改
- 要移动字典项到其他父级，需要重新创建

**响应示例：**
```json
{
  "base_resp": {
    "code": 0,
    "message": "success"
  },
  "item": {
    "id": 1,
    "dictionary_id": 1,
    "parent_id": 0,
    "code": "CN",
    "name": "中国大陆",
    "description": "中华人民共和国（大陆地区）",
    "icon_url": "https://example.com/flags/cn_new.png",
    "sort_order": 1,
    "level": 1,
    "path": "1",
    "status": 1,
    "created_at": "2025-10-26 16:00:00",
    "updated_at": "2025-10-26 16:10:00"
  }
}
```

---

### 9. 删除字典项

**接口地址：** `POST /api/v1/admin/dictionary/item/delete`

**权限要求：** 管理员

**请求参数：**
```json
{
  "id": 1
}
```

| 参数 | 类型 | 必填 | 说明 |
|------|------|------|------|
| id | int64 | 是 | 字典项ID |

**注意：** 
- 删除为软删除
- 删除字典项会递归删除其所有子节点

**响应示例：**
```json
{
  "base_resp": {
    "code": 0,
    "message": "success"
  }
}
```

---

### 10. 获取字典项列表

**接口地址：** `POST /api/v1/admin/dictionary/item/list`

**权限要求：** 管理员

**请求参数：**
```json
{
  "dictionary_id": 1,
  "parent_id": 0,
  "status": 1,
  "page": 1,
  "page_size": 100
}
```

| 参数 | 类型 | 必填 | 说明 |
|------|------|------|------|
| dictionary_id | int64 | 是 | 字典ID |
| parent_id | int64 | 否 | 父级ID筛选（不传则获取所有）|
| status | int32 | 否 | 状态筛选：1-启用，0-禁用 |
| page | int32 | 否 | 页码，默认1 |
| page_size | int32 | 否 | 每页数量，默认100 |

**响应示例：**
```json
{
  "base_resp": {
    "code": 0,
    "message": "success"
  },
  "items": [
    {
      "id": 1,
      "dictionary_id": 1,
      "parent_id": 0,
      "code": "CN",
      "name": "中国",
      "description": "中华人民共和国",
      "icon_url": "https://example.com/flags/cn.png",
      "sort_order": 1,
      "level": 1,
      "path": "1",
      "status": 1,
      "created_at": "2025-10-26 16:00:00",
      "updated_at": "2025-10-26 16:00:00"
    }
  ],
  "page_info": {
    "page": 1,
    "page_size": 100,
    "total": 1,
    "total_pages": 1
  }
}
```

---

### 11. 获取字典树形结构（公开接口）

**接口地址：** `POST /api/v1/dictionary/tree`

**权限要求：** 无（所有用户可访问）

**请求参数：**
```json
{
  "dictionary_code": "COUNTRY",
  "only_enabled": 1
}
```

| 参数 | 类型 | 必填 | 说明 |
|------|------|------|------|
| dictionary_code | string | 是 | 字典编码 |
| only_enabled | int32 | 否 | 是否只返回启用的：1-是（默认），0-否 |

**响应示例：**
```json
{
  "base_resp": {
    "code": 0,
    "message": "success"
  },
  "tree": [
    {
      "id": 1,
      "code": "CN",
      "name": "中国",
      "description": "中华人民共和国",
      "icon_url": "https://example.com/flags/cn.png",
      "sort_order": 1,
      "level": 1,
      "status": 1,
      "children": [
        {
          "id": 2,
          "code": "SH",
          "name": "上海",
          "description": "上海市",
          "icon_url": "",
          "sort_order": 1,
          "level": 2,
          "status": 1,
          "children": [
            {
              "id": 3,
              "code": "PD",
              "name": "浦东新区",
              "description": "",
              "icon_url": "",
              "sort_order": 1,
              "level": 3,
              "status": 1,
              "children": []
            }
          ]
        }
      ]
    }
  ]
}
```

---

## 错误码说明

| 错误码 | 说明 |
|--------|------|
| 0 | 成功 |
| 400 | 参数错误（如：字典编码格式错误、编码已存在等）|
| 401 | 未授权（token无效或已过期）|
| 403 | 权限不足（需要管理员权限）|
| 404 | 资源不存在（字典或字典项不存在）|
| 500 | 服务器内部错误 |

**常见错误示例：**

```json
{
  "base_resp": {
    "code": 400,
    "message": "字典编码只能包含大小写字母"
  }
}
```

```json
{
  "base_resp": {
    "code": 400,
    "message": "字典编码已存在"
  }
}
```

```json
{
  "base_resp": {
    "code": 404,
    "message": "字典不存在"
  }
}
```

---

## 使用场景示例

### 场景1：创建国家-省份-城市三级结构

**步骤1：创建字典**
```bash
POST /api/v1/admin/dictionary/create
{
  "code": "REGION",
  "name": "地区列表",
  "description": "国家-省份-城市三级结构"
}
```

**步骤2：创建国家（一级）**
```bash
POST /api/v1/admin/dictionary/item/create
{
  "dictionary_id": 1,
  "parent_id": 0,
  "code": "CN",
  "name": "中国",
  "icon_url": "https://example.com/flags/cn.png",
  "sort_order": 1
}
```

**步骤3：创建省份（二级）**
```bash
POST /api/v1/admin/dictionary/item/create
{
  "dictionary_id": 1,
  "parent_id": 1,  // 中国的ID
  "code": "SH",
  "name": "上海",
  "sort_order": 1
}
```

**步骤4：创建城区（三级）**
```bash
POST /api/v1/admin/dictionary/item/create
{
  "dictionary_id": 1,
  "parent_id": 2,  // 上海的ID
  "code": "PD",
  "name": "浦东新区",
  "sort_order": 1
}
```

**步骤5：前端获取树形结构**
```bash
POST /api/v1/dictionary/tree
{
  "dictionary_code": "REGION",
  "only_enabled": 1
}
```

---

### 场景2：管理业务分类

**创建业务分类字典**
```bash
POST /api/v1/admin/dictionary/create
{
  "code": "BUSINESS_CATEGORY",
  "name": "业务分类",
  "description": "系统业务分类"
}

// 添加一级分类
POST /api/v1/admin/dictionary/item/create
{
  "dictionary_id": 2,
  "parent_id": 0,
  "code": "TECH",
  "name": "科技",
  "icon_url": "https://example.com/icons/tech.png",
  "sort_order": 1
}

// 添加二级分类
POST /api/v1/admin/dictionary/item/create
{
  "dictionary_id": 2,
  "parent_id": 3,  // 科技的ID
  "code": "AI",
  "name": "人工智能",
  "icon_url": "https://example.com/icons/ai.png",
  "sort_order": 1
}
```

---

## 最佳实践

### 1. 字典编码命名规范
- 使用大写字母表示类别：`COUNTRY`, `REGION`, `CATEGORY`
- 使用小驼峰表示复合词：`businessType`, `userStatus`
- 只使用字母，不使用数字和特殊字符

### 2. 字典项编码建议
- 国家/地区：使用ISO标准代码（如：CN, US, JP）
- 省份/城市：使用拼音首字母或简称（如：SH, BJ, GD）
- 业务分类：使用英文简写（如：TECH, FIN, EDU）

### 3. 性能优化
- 字典树查询接口已优化，可直接在前端缓存
- 建议前端定时刷新字典数据（如每天更新一次）
- 字典数据变更频率低，适合使用Redis缓存

### 4. 权限设计
- 所有字典管理接口需要管理员权限
- 树形查询接口公开，供前端选择器使用
- 建议在后台管理系统中集中管理字典

### 5. 数据维护
- 禁用而非删除常用字典项，保证历史数据完整性
- 定期清理长期未使用的字典和字典项
- 重要字典变更需要做好数据迁移规划

---

## 技术实现说明

### 数据库设计特点
- **软删除**：使用 `deleted_at` 字段，保证数据一致性
- **树形结构**：通过 `parent_id` 和 `path` 字段实现
- **层级计算**：自动计算并维护 `level` 和 `path`
- **排序支持**：通过 `sort_order` 字段控制显示顺序
- **唯一性约束**：字典编码全局唯一，字典项编码在同一父级下唯一

### 核心功能
- **递归删除**：删除父节点自动删除所有子节点
- **层级验证**：创建字典项时自动验证父节点有效性
- **状态过滤**：支持只获取启用状态的数据
- **树形构建**：高效的树形数据结构构建算法

---

## 注意事项

1. **字典编码不可修改**：字典的 `code` 字段创建后不能修改，如需修改请重新创建
2. **字典项关联不可修改**：字典项的 `dictionary_id`、`parent_id`、`code` 创建后不能修改
3. **删除影响范围**：删除字典会删除所有字典项，删除字典项会删除所有子节点
4. **性能考虑**：建议字典项总数控制在10000以内，单个字典的层级不超过5层
5. **并发控制**：高并发场景下可能出现编码重复，建议在业务层加锁

---

## 联系方式

如有问题或建议，请联系技术支持团队。

**文档版本：** v1.0
**最后更新：** 2025-10-26

