package mw

import (
	"context"
	"net/http"
	"strings"

	"orbia_api/biz/consts"
	"orbia_api/biz/dal/mysql"
	"orbia_api/biz/service/user/rpc"
	"orbia_api/biz/utils"

	"github.com/cloudwego/hertz/pkg/app"
	"github.com/cloudwego/hertz/pkg/common/hlog"
)

const (
	AuthUserIDKey   = "auth_user_id"
	AuthUserKey     = "auth_user"
	AuthUserRoleKey = "auth_user_role"
)

var userRPC *rpc.UserRPC

// InitAuthMiddleware 初始化认证中间件
func InitAuthMiddleware(userRepo mysql.UserRepository) {
	userRPC = rpc.NewUserRPC(userRepo)
}

// AuthMiddleware JWT认证和角色鉴权中间件
// 支持多个角色参数，只要用户拥有其中任意一个角色即可访问
// 使用示例:
//   - AuthMiddleware(consts.RoleUser) - 仅普通用户可访问
//   - AuthMiddleware(consts.RoleAdmin) - 仅管理员可访问
//   - AuthMiddleware(consts.RoleUser, consts.RoleAdmin) - 普通用户和管理员都可访问
func AuthMiddleware(allowedRoles ...consts.UserRole) app.HandlerFunc {
	return func(ctx context.Context, c *app.RequestContext) {
		// 从Header中获取Authorization
		authHeader := string(c.GetHeader("Authorization"))
		if authHeader == "" {
			c.JSON(http.StatusUnauthorized, map[string]interface{}{
				"code":    401,
				"message": "Authorization header is required",
			})
			c.Abort()
			return
		}

		// 检查Bearer前缀
		if !strings.HasPrefix(authHeader, "Bearer ") {
			c.JSON(http.StatusUnauthorized, map[string]interface{}{
				"code":    401,
				"message": "Authorization header must start with Bearer",
			})
			c.Abort()
			return
		}

		// 提取token
		token := strings.TrimPrefix(authHeader, "Bearer ")
		if token == "" {
			c.JSON(http.StatusUnauthorized, map[string]interface{}{
				"code":    401,
				"message": "Token is required",
			})
			c.Abort()
			return
		}

		// 验证token
		userID, err := utils.ValidateToken(token)
		if err != nil {
			hlog.Errorf("JWT validation failed: %v", err)
			c.JSON(http.StatusUnauthorized, map[string]interface{}{
				"code":    401,
				"message": "Invalid or expired token",
			})
			c.Abort()
			return
		}

		// 从数据库获取用户信息
		if userRPC == nil {
			hlog.Error("userRPC is not initialized")
			c.JSON(http.StatusInternalServerError, map[string]interface{}{
				"code":    500,
				"message": "Internal server error",
			})
			c.Abort()
			return
		}

		user, err := userRPC.GetUserByID(userID)
		if err != nil {
			hlog.Errorf("Failed to get user info: %v", err)
			c.JSON(http.StatusUnauthorized, map[string]interface{}{
				"code":    401,
				"message": "User not found",
			})
			c.Abort()
			return
		}

		// 检查用户状态
		if user.Status != "normal" {
			hlog.Warnf("User %d with status '%s' tried to access resource", userID, user.Status)
			c.JSON(http.StatusForbidden, map[string]interface{}{
				"code":    403,
				"message": "User account is not in normal status",
			})
			c.Abort()
			return
		}

		// 将用户信息存储到上下文中
		userRole := consts.UserRole(user.Role)
		c.Set(AuthUserIDKey, userID)
		c.Set(AuthUserKey, user)
		c.Set(AuthUserRoleKey, userRole)

		// 如果指定了角色要求，进行角色验证
		if len(allowedRoles) > 0 {
			hasPermission := false
			for _, allowedRole := range allowedRoles {
				if userRole == allowedRole {
					hasPermission = true
					break
				}
			}

			if !hasPermission {
				hlog.Warnf("User %d with role '%s' tried to access resource requiring roles: %v", userID, userRole, allowedRoles)
				c.JSON(http.StatusForbidden, map[string]interface{}{
					"code":    403,
					"message": "Access denied: insufficient permissions",
				})
				c.Abort()
				return
			}
		}

		c.Next(ctx)
	}
}

// GetAuthUserID 从上下文中获取用户ID
func GetAuthUserID(c *app.RequestContext) (int64, bool) {
	userID, exists := c.Get(AuthUserIDKey)
	if !exists {
		return 0, false
	}

	if id, ok := userID.(int64); ok {
		return id, true
	}

	return 0, false
}

// GetAuthUser 从上下文中获取用户信息
func GetAuthUser(c *app.RequestContext) (*mysql.User, bool) {
	user, exists := c.Get(AuthUserKey)
	if !exists {
		return nil, false
	}

	if u, ok := user.(*mysql.User); ok {
		return u, true
	}

	return nil, false
}

// GetAuthUserRole 从上下文中获取用户角色
func GetAuthUserRole(c *app.RequestContext) (consts.UserRole, bool) {
	role, exists := c.Get(AuthUserRoleKey)
	if !exists {
		return "", false
	}

	if r, ok := role.(consts.UserRole); ok {
		return r, true
	}

	return "", false
}
