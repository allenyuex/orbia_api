# Orbia API

基于 Hertz 框架的高性能 API 服务

## 📁 项目目录结构

```
orbia_api/
├── biz/                          # 业务逻辑层
│   ├── consts/                  # 常量定义（错误码、通用常量）
│   │   ├── errors.go           # 错误码和错误消息
│   │   └── common.go           # 通用常量
│   ├── dal/                     # Data Access Layer (数据访问层)
│   │   └── mysql/              # MySQL 相关
│   │       ├── init.go         # 数据库初始化
│   │       └── user.go         # User DAO（示例）
│   ├── handler/                 # Handler 层（由 hz 生成）
│   │   ├── api/                # API handler
│   │   │   └── api_service.go
│   │   └── ping.go
│   ├── infra/                   # 基础设施代码
│   │   └── config/             # 配置管理
│   │       └── config.go
│   ├── mw/                      # Middleware（中间件）
│   │   ├── cors.go             # 跨域中间件
│   │   ├── logger.go           # 日志中间件
│   │   └── recovery.go         # 恢复中间件
│   ├── model/                   # 数据模型（由 hz 生成）
│   │   └── api/
│   │       └── api.go
│   ├── router/                  # 路由配置（由 hz 生成）
│   │   ├── api/
│   │   │   ├── api.go
│   │   │   └── middleware.go
│   │   └── register.go
│   ├── service/                 # Service 层（业务逻辑）
│   │   └── user_service.go     # 用户服务
│   └── utils/                   # 工具类
│       └── response.go         # 统一响应工具
├── conf/                        # 配置文件
│   └── config.yaml             # 主配置文件
├── idl/                         # IDL 定义文件
│   └── api.thrift              # API IDL 定义
├── script/                      # 脚本文件
│   └── bootstrap.sh
├── main.go                      # 主程序入口
├── router.go                    # 路由注册（由 hz 生成）
├── router_gen.go               # 路由生成器（由 hz 生成）
├── go.mod                       # Go 模块依赖
└── README.md                    # 本文件
```

## 🏗️ 分层架构

项目采用标准的分层架构，数据流向为：

```
Request → Router → Handler → Service → DAL → Database
         ↓                                      ↑
      Response ←――――――――――――――――――――――――――――――――┘
```

### 各层职责

1. **Router**: 路由配置，由 hz 工具根据 IDL 自动生成
2. **Handler**: 请求参数绑定和验证，调用 Service
3. **Service**: 业务逻辑处理层
4. **DAL (Data Access Layer)**: 数据访问层，封装数据库操作
5. **Model**: 数据模型，由 hz 根据 IDL 生成

## 🚀 快速开始

### 1. 安装依赖

```bash
cd orbia_api
go mod tidy
```

### 2. 配置数据库

编辑 `conf/config.yaml`，确保数据库配置正确：

```yaml
database:
  mysql:
    host: 127.0.0.1
    port: 3306
    database: orbia
    username: root
    password: root123
```

### 3. 创建数据库

```bash
mysql -uroot -proot123 -e "CREATE DATABASE IF NOT EXISTS orbia CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci;"
```

### 4. 启动服务

```bash
go run main.go
```

服务将在 `http://localhost:8888` 启动

### 5. 测试接口

```bash
# 健康检查
curl http://localhost:8888/health

# Hello Demo
curl "http://localhost:8888/api/v1/demo/hello?name=Orbia"

# 创建用户
curl -X POST http://localhost:8888/api/v1/users \
  -H "Content-Type: application/json" \
  -d '{"name":"张三","email":"zhangsan@example.com","phone":"13800138000"}'

# 获取用户列表
curl http://localhost:8888/api/v1/users
```

## 📝 使用 Hz 工具添加新功能

### 安装 Hz 工具

```bash
go install github.com/cloudwego/hertz/cmd/hz@latest
```

### 添加新的 API 接口

#### 步骤 1: 编辑 IDL 文件

编辑 `idl/api.thrift`，添加新的接口定义：

```thrift
// 添加新的请求和响应结构
struct NewFeatureReq {
    1: required string param1 (api.body="param1")
}

struct NewFeatureResp {
    1: string result
    2: BaseResp base_resp
}

// 在 service 中添加方法
service ApiService {
    // ... 现有方法 ...
    
    // 新方法
    NewFeatureResp NewFeature(1: NewFeatureReq req) (api.post="/api/v1/new-feature")
}
```

#### 步骤 2: 使用 Hz 更新代码

```bash
cd orbia_api

# 更新代码（会自动生成 handler、router、model）
hz update -idl idl/api.thrift
```

Hz 会自动更新：
- `biz/handler/api/api_service.go` - 添加新的 handler 函数框架
- `biz/model/api/api.go` - 添加新的数据模型
- `biz/router/api/api.go` - 添加新的路由

#### 步骤 3: 实现 Service 层

创建或编辑 `biz/service/your_service.go`：

```go
package service

import (
    "orbia_api/biz/consts"
    "orbia_api/biz/dal/mysql"
    "orbia_api/biz/model/api"
)

type YourService struct {
    dao *mysql.YourDAO
}

func NewYourService() *YourService {
    return &YourService{
        dao: mysql.NewYourDAO(),
    }
}

func (s *YourService) NewFeature(req *api.NewFeatureReq) (*api.NewFeatureResp, error) {
    // 实现业务逻辑
    
    return &api.NewFeatureResp{
        Result: "success",
        BaseResp: &api.BaseResp{
            Code:    consts.SuccessCode,
            Message: consts.SuccessMsg,
        },
    }, nil
}
```

#### 步骤 4: 在 Handler 中调用 Service

编辑 `biz/handler/api/api_service.go`：

```go
var (
    yourService = service.NewYourService()
)

func NewFeature(ctx context.Context, c *app.RequestContext) {
    var err error
    var req api.NewFeatureReq
    err = c.BindAndValidate(&req)
    if err != nil {
        utils.ParamError(c, err.Error())
        return
    }

    // 调用 service 层
    resp, err := yourService.NewFeature(&req)
    if err != nil {
        log.Printf("NewFeature service error: %v", err)
        utils.SystemError(c)
        return
    }

    c.JSON(consts.StatusOK, resp)
}
```

#### 步骤 5: (可选) 添加 DAL 层

如果需要数据库操作，创建 `biz/dal/mysql/your_model.go`：

```go
package mysql

import "gorm.io/gorm"

type YourModel struct {
    ID   uint   `gorm:"primaryKey"`
    Name string `gorm:"size:100"`
}

type YourDAO struct {
    db *gorm.DB
}

func NewYourDAO() *YourDAO {
    return &YourDAO{db: DB}
}

func (dao *YourDAO) Create(model *YourModel) error {
    return dao.db.Create(model).Error
}

// 添加其他 CRUD 方法...
```

记得在 `biz/dal/mysql/init.go` 的 `autoMigrate()` 函数中添加模型：

```go
func autoMigrate() error {
    return DB.AutoMigrate(
        &User{},
        &YourModel{}, // 添加新模型
    )
}
```

## 📖 开发规范

### 1. 代码组织

- **Handler**: 只负责参数绑定和调用 Service，不包含业务逻辑
- **Service**: 包含所有业务逻辑，可以调用多个 DAO
- **DAL**: 只负责数据库操作，不包含业务逻辑

### 2. 错误处理

统一使用 `biz/consts/errors.go` 中定义的错误码：

```go
return &YourResp{
    BaseResp: &api.BaseResp{
        Code:    consts.UserNotFoundCode,
        Message: consts.UserNotFoundMsg,
    },
}, nil
```

### 3. 响应格式

所有响应都包含 `BaseResp`：

```go
struct BaseResp {
    1: i32 code       // 0: 成功, 非0: 错误
    2: string message // 消息描述
}
```

### 4. 日志记录

使用标准 log 包记录关键信息：

```go
log.Printf("Operation failed: %v", err)
```

## 🔧 配置说明

### 环境变量

- `CONFIG_PATH`: 配置文件路径，默认为 `./conf/config.yaml`

### 配置文件 (conf/config.yaml)

- `server`: 服务器配置（地址、端口、超时）
- `database.mysql`: MySQL 数据库配置
- `redis`: Redis 配置（预留）
- `jwt`: JWT 认证配置（预留）
- `log`: 日志配置（预留）

## 📚 API 文档

### 基础接口

#### GET /health
健康检查

**响应示例**:
```json
{
  "status": "ok",
  "message": "Orbia API is running"
}
```

### Demo 接口

#### GET /api/v1/demo/hello?name=xxx
Hello 测试接口

**响应示例**:
```json
{
  "message": "Hello, Orbia! Welcome to Orbia API",
  "timestamp": 1697011200,
  "base_resp": {
    "code": 0,
    "message": "success"
  }
}
```

### 用户管理接口

#### POST /api/v1/users
创建用户

**请求体**:
```json
{
  "name": "张三",
  "email": "zhangsan@example.com",
  "phone": "13800138000"
}
```

**响应示例**:
```json
{
  "user_id": 1,
  "base_resp": {
    "code": 0,
    "message": "success"
  }
}
```

#### GET /api/v1/users/:user_id
获取用户信息

**响应示例**:
```json
{
  "user": {
    "id": 1,
    "name": "张三",
    "email": "zhangsan@example.com",
    "phone": "13800138000",
    "created_at": "2024-10-10 10:00:00",
    "updated_at": "2024-10-10 10:00:00"
  },
  "base_resp": {
    "code": 0,
    "message": "success"
  }
}
```

#### GET /api/v1/users?page=1&page_size=20
获取用户列表

**响应示例**:
```json
{
  "users": [...],
  "total": 100,
  "base_resp": {
    "code": 0,
    "message": "success"
  }
}
```

## 🛠️ 常用命令

```bash
# 下载依赖
go mod tidy

# 运行服务
go run main.go

# 更新 IDL 后重新生成代码
hz update -idl idl/api.thrift

# 构建
go build -o bin/orbia_api main.go

# 运行测试
go test ./...
```

## 🐛 故障排除

### 问题 1: 数据库连接失败
- 检查 MySQL 是否运行
- 验证 `conf/config.yaml` 中的数据库配置
- 确保数据库已创建

### 问题 2: Hz 命令找不到
```bash
# 重新安装 Hz
go install github.com/cloudwego/hertz/cmd/hz@latest

# 确保 $GOPATH/bin 在 PATH 中
export PATH=$PATH:$(go env GOPATH)/bin
```

### 问题 3: 端口被占用
修改 `conf/config.yaml` 中的端口配置

## 📖 参考文档

- [Hertz 官方文档](https://www.cloudwego.io/docs/hertz/)
- [Hz CLI 工具](https://www.cloudwego.io/docs/hertz/tutorials/toolkit/toolkit/)
- [GORM 文档](https://gorm.io/docs/)
- [Thrift IDL 语法](https://thrift.apache.org/docs/idl)

## 📝 更新日志

### v1.0.0 (2024-10-10)
- ✅ 初始化项目结构
- ✅ 实现用户管理功能
- ✅ 集成 MySQL 数据库
- ✅ 添加基础中间件
- ✅ 完善文档

