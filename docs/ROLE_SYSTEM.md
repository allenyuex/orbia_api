# 用户角色权限系统

## 概述

本系统实现了基于角色的访问控制（RBAC），目前支持两种角色：
- **普通用户（user）**：默认角色，拥有基本权限
- **管理员（admin）**：拥有管理权限，可以访问受限 API

## 数据库设计

### 用户表字段
```sql
role ENUM('user', 'admin') NOT NULL DEFAULT 'user' COMMENT '用户角色：user-普通用户，admin-管理员'
```

所有新注册的用户默认为普通用户角色。

## 核心组件

### 1. 角色常量定义
位置：`biz/consts/role.go`

```go
const (
    RoleUser  UserRole = "user"   // 普通用户
    RoleAdmin UserRole = "admin"  // 管理员
)
```

### 2. 认证中间件
位置：`biz/mw/auth.go`

功能：
- 验证 JWT Token
- 从数据库获取用户完整信息（包括角色）
- 将用户 ID 和角色存储到请求上下文中

### 3. 角色权限中间件
位置：`biz/mw/role.go`

提供了三个便捷的中间件：

#### AdminOnlyMiddleware
仅管理员可访问
```go
mw.AdminOnlyMiddleware()
```

#### UserOrAdminMiddleware
普通用户和管理员都可访问
```go
mw.UserOrAdminMiddleware()
```

#### RoleMiddleware
自定义角色组合（支持多角色）
```go
mw.RoleMiddleware(consts.RoleUser, consts.RoleAdmin)
```

## 使用示例

### 在路由中应用角色权限

#### 示例 1：仅管理员可访问的 API
```go
// 在路由注册中
adminGroup := h.Group("/api/v1/admin")
adminGroup.Use(mw.AuthMiddleware())          // 先验证身份
adminGroup.Use(mw.AdminOnlyMiddleware())     // 再验证权限
{
    adminGroup.POST("/users", handler.GetAllUsers)
    adminGroup.POST("/kol/approve", handler.ApproveKOL)
    adminGroup.DELETE("/user/:id", handler.DeleteUser)
}
```

#### 示例 2：普通用户和管理员都可访问
```go
userGroup := h.Group("/api/v1/user")
userGroup.Use(mw.AuthMiddleware())           // 验证身份
userGroup.Use(mw.UserOrAdminMiddleware())    // 验证权限
{
    userGroup.POST("/profile", handler.GetProfile)
    userGroup.POST("/update-profile", handler.UpdateProfile)
}
```

#### 示例 3：在 handler 中检查角色
```go
func SomeHandler(ctx context.Context, c *app.RequestContext) {
    // 获取用户角色
    role, exists := mw.GetAuthUserRole(c)
    if !exists {
        // 角色信息不存在
        return
    }

    // 根据角色执行不同逻辑
    if role.IsAdmin() {
        // 管理员特殊逻辑
    } else {
        // 普通用户逻辑
    }
}
```

## API 响应

所有用户相关的 API 现在都会返回角色信息：

```json
{
    "user": {
        "id": 1,
        "wallet_address": "0x...",
        "email": "user@example.com",
        "nickname": "Test User",
        "avatar_url": "https://...",
        "role": "user",  // 角色字段
        "created_at": "2024-01-01 00:00:00",
        "updated_at": "2024-01-01 00:00:00"
    },
    "base_resp": {
        "code": 200,
        "message": "Success"
    }
}
```

## 创建管理员账号

### 方法 1：通过 SQL 直接创建
```sql
-- 创建管理员账号（邮箱方式）
INSERT INTO orbia_user (email, password_hash, nickname, role) 
VALUES ('admin@orbia.com', 'hashed_password', 'Admin', 'admin');

-- 或者将现有用户提升为管理员
UPDATE orbia_user SET role = 'admin' WHERE email = 'user@example.com';
```

### 方法 2：使用脚本提升用户
```bash
# 创建一个临时脚本
mysql -h127.0.0.1 -P3306 -uroot -proot123 -e \
  "USE orbia; UPDATE orbia_user SET role = 'admin' WHERE id = 1;"
```

## 未来扩展

系统架构支持轻松扩展更多角色和细粒度权限：

### 扩展角色
1. 在数据库中添加新的 ENUM 值：
```sql
ALTER TABLE orbia_user MODIFY role ENUM('user', 'admin', 'moderator', 'premium') NOT NULL DEFAULT 'user';
```

2. 在 `biz/consts/role.go` 中添加新常量：
```go
const (
    RoleUser      UserRole = "user"
    RoleAdmin     UserRole = "admin"
    RoleModerator UserRole = "moderator"  // 新角色
    RolePremium   UserRole = "premium"    // 新角色
)
```

3. 创建新的中间件或使用 RoleMiddleware：
```go
mw.RoleMiddleware(consts.RoleModerator, consts.RoleAdmin)
```

### 细粒度权限管理

如果未来需要更细粒度的权限管理，可以：

1. 创建权限表：
```sql
CREATE TABLE orbia_permission (
    id BIGINT PRIMARY KEY AUTO_INCREMENT,
    name VARCHAR(100) NOT NULL,
    description TEXT
);

CREATE TABLE orbia_role_permission (
    role ENUM('user', 'admin', ...) NOT NULL,
    permission_id BIGINT NOT NULL,
    FOREIGN KEY (permission_id) REFERENCES orbia_permission(id)
);
```

2. 实现基于权限的中间件：
```go
func PermissionMiddleware(requiredPermissions ...string) app.HandlerFunc {
    // 检查用户是否拥有所需权限
}
```

## 安全注意事项

1. **默认最小权限原则**：所有新用户默认为普通用户
2. **角色信息存储在 JWT 外**：角色信息从数据库实时获取，防止 JWT 中的角色信息被篡改
3. **中间件顺序**：务必先应用 `AuthMiddleware()`，再应用角色权限中间件
4. **审计日志**：建议在权限验证失败时记录审计日志

## 测试

### 测试角色权限
```bash
# 1. 注册/登录获取 token（默认为普通用户）
TOKEN=$(curl -s -X POST http://localhost:8080/api/v1/auth/email-login \
  -H "Content-Type: application/json" \
  -d '{"email":"user@test.com","password":"password"}' | jq -r '.token')

# 2. 访问普通用户可访问的 API
curl -X POST http://localhost:8080/api/v1/user/profile \
  -H "Authorization: Bearer $TOKEN"

# 3. 尝试访问仅管理员可访问的 API（应该返回 403）
curl -X POST http://localhost:8080/api/v1/admin/users \
  -H "Authorization: Bearer $TOKEN"

# 4. 将用户提升为管理员
mysql -h127.0.0.1 -uroot -proot123 -e \
  "USE orbia; UPDATE orbia_user SET role = 'admin' WHERE email = 'user@test.com';"

# 5. 再次尝试访问管理员 API（应该成功）
curl -X POST http://localhost:8080/api/v1/admin/users \
  -H "Authorization: Bearer $TOKEN"
```

## 常见问题

### Q: 修改用户角色后需要重新登录吗？
A: 需要。因为角色信息是在每次请求时从数据库读取的，但 JWT token 包含的用户 ID 不会变化，所以理论上不需要重新登录。但为了安全考虑，建议在角色变更后要求用户重新登录。

### Q: 如何批量创建管理员？
A: 使用 SQL 批量更新：
```sql
UPDATE orbia_user 
SET role = 'admin' 
WHERE email IN ('admin1@orbia.com', 'admin2@orbia.com');
```

### Q: 如何查询所有管理员？
A: 使用 SQL 查询：
```sql
SELECT id, email, nickname, role, created_at 
FROM orbia_user 
WHERE role = 'admin';
```

## 总结

本角色权限系统的特点：
- ✅ 简单易用：提供便捷的中间件函数
- ✅ 易于扩展：支持添加更多角色和权限
- ✅ 安全可靠：角色信息实时从数据库获取
- ✅ 灵活配置：支持多角色组合验证
- ✅ 符合最佳实践：基于 RBAC 模型设计


