# Orbia API

Orbia API 是一个基于 Hertz 框架的 Go 语言后端服务。

## 目录

- [快速开始](#快速开始)
- [配置管理](#配置管理)
- [开发指南](#开发指南)
- [部署](#部署)
- [API 文档](#api-文档)

## 快速开始

### 前置要求

- Go 1.21+
- MySQL 8.0+
- Redis 6.0+

### 安装依赖

```bash
go mod tidy
```

### 初始化数据库

```bash
./script/init_db.sh
```

### 启动服务

```bash
# 使用开发环境配置（默认）
./start.sh

# 或者直接运行
go run .
```

服务将在 `http://localhost:8888` 启动。

## 配置管理

### 配置文件结构

```
conf/
├── dev/          # 开发环境配置
│   └── config.yaml
└── prod/         # 生产环境配置
    └── config.yaml
```

### 环境切换

通过环境变量 `ORBIA_ENV` 来切换配置环境：

```bash
# 开发环境（默认）
ORBIA_ENV=dev ./start.sh

# 生产环境
ORBIA_ENV=prod ./start.sh
```

### 生产环境配置

生产环境支持通过环境变量设置敏感信息。详细说明请参考：

- [配置文档](docs/CONFIG.md)
- [环境变量示例](env.example)

## 开发指南

### 项目结构

```
.
├── biz/                    # 业务逻辑层
│   ├── consts/            # 常量定义
│   ├── dal/               # 数据访问层
│   │   ├── model/         # 数据库模型（GORM 生成）
│   │   └── mysql/         # MySQL 仓储层
│   ├── handler/           # HTTP 处理器
│   ├── model/             # API 模型
│   ├── router/            # 路由定义
│   ├── service/           # 业务逻辑服务层
│   │   └── */rpc.go      # RPC 通用逻辑
│   ├── mw/                # 中间件
│   └── utils/             # 工具函数
├── conf/                   # 配置文件
│   ├── dev/               # 开发环境
│   └── prod/              # 生产环境
├── docs/                   # 文档
├── idl/                    # Thrift IDL 定义
├── script/                # 脚本工具
└── sql/                    # SQL 脚本
```

### 开发规范

1. **API 定义**: 所有接口使用 POST 方法，入参出参使用 JSON 格式
2. **IDL 优先**: 使用 hz 框架生成基础脚手架代码
3. **分层架构**:
   - `router/handler`: 路由和基本请求处理
   - `service`: 业务逻辑实现
   - `service/rpc`: 通用逻辑抽象
   - `dal/mysql`: 数据库操作

### 添加新功能

1. 定义或更新 IDL 文件（`idl/*.thrift`）
2. 运行 hz 命令生成代码：
   ```bash
   hz update -idl idl/your_feature.thrift -module orbia_api
   ```
3. 实现 service 层业务逻辑
4. 实现 dal 层数据访问

### 数据库相关

- 数据库名：`orbia`
- 表命名：`orbia_xxxx`
- 模型位置：`biz/dal/model/`
- 仓储位置：`biz/dal/mysql/`

更新数据库结构：

```bash
# 修改 sql/init.sql 后运行
./script/init_db.sh
```

## 部署

### Docker 部署

```bash
# 构建镜像
docker build -t orbia-api .

# 运行容器
docker run -d \
  -p 8888:8888 \
  -e ORBIA_ENV=prod \
  -e MYSQL_HOST=your_db_host \
  -e MYSQL_PASSWORD=your_password \
  -e JWT_SECRET=your_jwt_secret \
  --name orbia-api \
  orbia-api
```

### 直接部署

```bash
# 构建
./build.sh

# 设置环境变量
export ORBIA_ENV=prod
export MYSQL_HOST=your_db_host
export MYSQL_PASSWORD=your_password
export JWT_SECRET=your_jwt_secret

# 运行
./orbia_api
```

## API 文档

详细的 API 文档请参考 `docs/` 目录：

- [Campaign API](docs/CAMPAIGN_API.md)
- [Conversation API](docs/CONVERSATION_API.md)
- [Dashboard API](docs/DASHBOARD_API.md)
- [KOL Order API](docs/KOL_ORDER_API.md)
- [配置文档](docs/CONFIG.md)

## 常用脚本

```bash
# 初始化数据库
./script/init_db.sh

# 初始化字典数据
./script/init_dictionary.sh

# 创建管理员账号
./script/create_admin.sh

# 测试配置
./test_config.sh
```

## 技术栈

- **框架**: CloudWeGo Hertz
- **数据库**: MySQL 8.0, Redis
- **ORM**: GORM
- **对象存储**: Cloudflare R2
- **邮件**: SMTP
- **认证**: JWT

## 开发工具

- **hz**: Hertz 代码生成工具
- **thriftgo**: Thrift IDL 编译器

## License

Copyright © 2024 Orbia

