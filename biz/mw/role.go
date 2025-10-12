package mw

import (
	"context"
	"net/http"

	"orbia_api/biz/consts"

	"github.com/cloudwego/hertz/pkg/app"
	"github.com/cloudwego/hertz/pkg/common/hlog"
)

const (
	AuthUserRoleKey = "auth_user_role"
)

// RoleMiddleware 角色权限中间件
// 支持多个角色，只要用户拥有其中任意一个角色即可访问
func RoleMiddleware(requiredRoles ...consts.UserRole) app.HandlerFunc {
	return func(ctx context.Context, c *app.RequestContext) {
		// 从上下文中获取用户角色
		role, exists := GetAuthUserRole(c)
		if !exists {
			hlog.Error("User role not found in context")
			c.JSON(http.StatusForbidden, map[string]interface{}{
				"code":    403,
				"message": "Access denied: role information not found",
			})
			c.Abort()
			return
		}

		// 检查用户是否拥有所需的角色
		hasPermission := false
		for _, requiredRole := range requiredRoles {
			if role == requiredRole {
				hasPermission = true
				break
			}
		}

		if !hasPermission {
			hlog.Warnf("User with role %s tried to access resource requiring roles: %v", role, requiredRoles)
			c.JSON(http.StatusForbidden, map[string]interface{}{
				"code":    403,
				"message": "Access denied: insufficient permissions",
			})
			c.Abort()
			return
		}

		c.Next(ctx)
	}
}

// AdminOnlyMiddleware 仅管理员可访问的中间件
func AdminOnlyMiddleware() app.HandlerFunc {
	return RoleMiddleware(consts.RoleAdmin)
}

// UserOrAdminMiddleware 普通用户和管理员都可访问的中间件
func UserOrAdminMiddleware() app.HandlerFunc {
	return RoleMiddleware(consts.RoleUser, consts.RoleAdmin)
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

// SetAuthUserRole 设置用户角色到上下文
func SetAuthUserRole(c *app.RequestContext, role consts.UserRole) {
	c.Set(AuthUserRoleKey, role)
}

