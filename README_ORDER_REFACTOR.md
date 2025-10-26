# 订单系统重构说明

## 重构概述

本次重构将原来的统一订单系统（Order）分离为两个独立的订单类型：
1. **KOL 订单系统** (KolOrder) - 用于用户向 KOL 下单
2. **广告订单系统** (AdOrder) - 用于用户创建广告订单

## 主要变更

### 1. IDL 定义

#### 新增文件
- `idl/kol_order.thrift` - KOL 订单接口定义
- `idl/ad_order.thrift` - 广告订单接口定义

#### 旧文件（已废弃）
- `idl/order.thrift` - 原统一订单接口（已被新的 thrift 文件替代）

### 2. 数据库表

#### 新增/更新的表
- `orbia_kol_order` - KOL 订单表（已更新，增加了更多字段）
  - 新增字段：title, requirement_description, video_type, video_duration, target_audience, expected_delivery_date, additional_requirements
  - 订单 ID 格式：`KORD_{timestamp}_{random}`
  
- `orbia_ad_order` - 广告订单表（新增）
  - 订单 ID 格式：`ADORD_{timestamp}_{random}`

### 3. API 接口路由

#### KOL 订单接口（需要用户认证）
- `POST /api/v1/kol-order/create` - 用户创建 KOL 订单
- `POST /api/v1/kol-order/detail` - 获取 KOL 订单详情
- `POST /api/v1/kol-order/user/list` - 用户查询自己的 KOL 订单列表（支持模糊搜索）
- `POST /api/v1/kol-order/cancel` - 用户取消 KOL 订单
- `POST /api/v1/kol-order/kol/list` - KOL 查询收到的订单列表（支持模糊搜索）
- `POST /api/v1/kol-order/status/update` - KOL 更新订单状态

#### 广告订单接口（需要用户认证）
- `POST /api/v1/ad-order/create` - 用户创建广告订单
- `POST /api/v1/ad-order/detail` - 获取广告订单详情
- `POST /api/v1/ad-order/user/list` - 用户查询自己的广告订单列表（支持模糊搜索）
- `POST /api/v1/ad-order/cancel` - 用户取消广告订单
- `POST /api/v1/ad-order/admin/list` - 管理员查询所有广告订单列表（需要管理员权限）
- `POST /api/v1/ad-order/status/update` - 管理员更新广告订单状态（需要管理员权限）

#### 旧接口（已废弃）
- `/api/v1/order/*` - 原统一订单接口

### 4. 代码结构

#### 新增目录和文件

##### DAL 层
- `biz/dal/mysql/ad_order.go` - 广告订单数据库操作
- `biz/dal/mysql/order.go` - KOL 订单数据库操作（已更新）

##### Service 层
- `biz/service/kol_order/kol_order_service.go` - KOL 订单业务逻辑
- `biz/service/ad_order/ad_order_service.go` - 广告订单业务逻辑

##### Handler 层
- `biz/handler/kol_order/kol_order_service.go` - KOL 订单处理器
- `biz/handler/ad_order/ad_order_service.go` - 广告订单处理器

##### Router 层
- `biz/router/kol_order/` - KOL 订单路由配置
- `biz/router/ad_order/` - 广告订单路由配置

##### Model 层
- `biz/model/kol_order/` - KOL 订单模型
- `biz/model/ad_order/` - 广告订单模型

#### 旧目录（待清理）
- `biz/model/order/` - 原订单模型目录
- `biz/handler/order/` - 原订单处理器目录
- `biz/service/order/` - 原订单服务目录
- `biz/router/order/` - 原订单路由目录

### 5. 订单 ID 生成

新增了两个订单 ID 生成函数：
- `utils.GenerateKolOrderID()` - 生成 KOL 订单 ID（格式：`KORD_{timestamp}_{random}`）
- `utils.GenerateAdOrderID()` - 生成广告订单 ID（格式：`ADORD_{timestamp}_{random}`）

旧函数（已标记为废弃）：
- `utils.GenerateOrderID()` - 原订单 ID 生成函数

### 6. 功能特性

#### KOL 订单
- 用户下单时需提供：订单标题、合作需求描述、视频类型、视频预计时长、目标受众、期望交付日期、额外要求、选择的方案
- 支持用户查询自己的订单列表（模糊搜索：标题、订单ID、KOL名称）
- 支持 KOL 查询收到的订单列表（模糊搜索：标题、订单ID、用户名称）
- KOL 可以更新订单状态（确认、进行中、完成、取消）
- 用户可以取消订单

#### 广告订单
- 用户下单时需提供：标题、描述、预算、广告类型、目标受众、开始日期、结束日期
- 支持用户查询自己的订单列表（模糊搜索：标题、订单ID、描述）
- 管理员可以查询所有广告订单列表
- 管理员可以更新订单状态（批准、进行中、完成、取消）
- 用户可以取消订单

## 迁移说明

由于服务尚未上线，无需考虑数据迁移问题。旧的 order 表和代码可以直接删除。

## 清理计划

以下旧代码可以在确认新系统运行正常后删除：
1. `idl/order.thrift`
2. `biz/model/order/` 目录
3. `biz/handler/order/` 目录
4. `biz/service/order/` 目录
5. `biz/router/order/` 目录

## 数据库初始化

运行以下命令初始化数据库：
```bash
bash script/init_db.sh
```

## 测试建议

1. 测试 KOL 订单创建流程
2. 测试用户查询自己的 KOL 订单列表（包括模糊搜索）
3. 测试 KOL 查询收到的订单列表（包括模糊搜索）
4. 测试 KOL 更新订单状态
5. 测试用户取消 KOL 订单
6. 测试广告订单创建流程
7. 测试用户查询自己的广告订单列表（包括模糊搜索）
8. 测试管理员查询所有广告订单列表
9. 测试管理员更新广告订单状态
10. 测试用户取消广告订单

## 注意事项

1. 所有订单接口都需要用户认证
2. 广告订单的管理员接口需要管理员权限
3. 订单 ID 由程序自动生成，格式为字符串
4. 两种订单系统完全解耦，互不影响
5. 代码符合 Golang 和 Hertz 的最佳实践

---

重构完成日期：2025-10-26

