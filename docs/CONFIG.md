# 配置文件说明

## 概述

Orbia API 支持多环境配置，通过环境变量 `ORBIA_ENV` 来切换不同的配置环境。

## 配置目录结构

```
conf/
├── dev/          # 开发环境配置
│   └── config.yaml
└── prod/         # 生产环境配置
    └── config.yaml
```

## 环境切换

### 方法一：设置环境变量

```bash
# 使用开发环境（默认）
export ORBIA_ENV=dev
go run .

# 使用生产环境
export ORBIA_ENV=prod
go run .
```

### 方法二：启动时指定

```bash
# 开发环境
ORBIA_ENV=dev go run .

# 生产环境
ORBIA_ENV=prod go run .
```

### 方法三：使用 start.sh 脚本

```bash
# 默认使用 dev 环境
./start.sh

# 使用 prod 环境
ORBIA_ENV=prod ./start.sh
```

## 配置文件说明

### 开发环境 (dev)

开发环境配置文件 `conf/dev/config.yaml` 包含所有开发环境所需的配置，如本地数据库、Redis 等。

配置值直接写在配置文件中，便于本地开发调试。

### 生产环境 (prod)

生产环境配置文件 `conf/prod/config.yaml` 支持通过环境变量来设置敏感信息。

#### 环境变量格式

配置文件支持以下格式的环境变量替换：

```yaml
# 格式1: ${VAR_NAME:default_value}
# 如果环境变量 VAR_NAME 存在，使用其值；否则使用 default_value
database:
  mysql:
    host: ${MYSQL_HOST:127.0.0.1}
    
# 格式2: ${VAR_NAME}
# 如果环境变量 VAR_NAME 存在，使用其值；否则为空字符串
jwt:
  secret: ${JWT_SECRET}
```

#### 生产环境推荐设置

生产环境建议通过环境变量设置以下敏感信息：

```bash
# 数据库配置
export MYSQL_HOST=your_db_host
export MYSQL_PORT=3306
export MYSQL_DATABASE=orbia
export MYSQL_USERNAME=your_username
export MYSQL_PASSWORD=your_password

# Redis 配置
export REDIS_HOST=your_redis_host
export REDIS_PORT=6379
export REDIS_PASSWORD=your_redis_password

# JWT 配置
export JWT_SECRET=your_jwt_secret_key

# R2 存储配置
export R2_ENDPOINT=your_r2_endpoint
export R2_ACCESS_KEY=your_r2_access_key
export R2_SECRET_KEY=your_r2_secret_key
export R2_BUCKET=your_bucket_name
export R2_PUBLIC_URL=your_public_url

# SMTP 邮件配置
export SMTP_SERVER=your_smtp_server
export SMTP_PORT=your_smtp_port
export SMTP_USERNAME=your_smtp_username
export SMTP_PASSWORD=your_smtp_password
export SMTP_EMAIL=your_email
export SMTP_FROM_NAME=your_from_name
```

## 配置结构

配置文件包含以下主要部分：

1. **server**: 服务器配置（端口、超时等）
2. **database**: 数据库配置（MySQL）
3. **redis**: Redis 配置
4. **jwt**: JWT 认证配置
5. **log**: 日志配置
6. **r2**: Cloudflare R2 对象存储配置
7. **smtp**: 邮件服务配置
8. **verification_code**: 验证码配置

详细配置项请参考配置文件中的注释。

## 最佳实践

### 开发环境

1. 直接修改 `conf/dev/config.yaml` 文件
2. 不要提交敏感信息到 git 仓库
3. 使用本地数据库和服务进行开发

### 生产环境

1. 不要在配置文件中直接写入敏感信息
2. 通过环境变量或密钥管理服务设置敏感配置
3. 定期更新 JWT secret 和其他密钥
4. 使用强密码保护数据库和 Redis

## 添加新环境

如需添加新的环境（如 staging），只需：

1. 创建新的配置目录：
   ```bash
   mkdir -p conf/staging
   cp conf/dev/config.yaml conf/staging/config.yaml
   ```

2. 修改 `conf/staging/config.yaml` 中的配置

3. 使用新环境：
   ```bash
   ORBIA_ENV=staging go run .
   ```

## 故障排查

### 配置文件找不到

错误信息：
```
❌ Failed to load config: 读取配置文件失败 (conf/xxx/config.yaml): ...
```

解决方法：
1. 检查 `ORBIA_ENV` 环境变量是否设置正确
2. 确认对应的配置文件是否存在
3. 检查文件路径和权限

### 配置解析失败

错误信息：
```
❌ Failed to load config: 解析配置文件失败: ...
```

解决方法：
1. 检查 YAML 格式是否正确（注意缩进）
2. 确认所有必需的配置项都已设置
3. 检查环境变量格式是否正确

### 数据库连接失败

确保：
1. 数据库配置正确
2. 环境变量已正确设置（生产环境）
3. 数据库服务正在运行
4. 网络连接正常

