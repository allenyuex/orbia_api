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
	AuthUserIDKey = "auth_user_id"
)

var userRPC *rpc.UserRPC

// InitAuthMiddleware 初始化认证中间件
func InitAuthMiddleware(userRepo mysql.UserRepository) {
	userRPC = rpc.NewUserRPC(userRepo)
}

// AuthMiddleware JWT认证中间件
func AuthMiddleware() app.HandlerFunc {
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

		// 获取用户信息（包括角色）
		if userRPC != nil {
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

			// 将用户角色存储到上下文中
			SetAuthUserRole(c, consts.UserRole(user.Role))
		}

		// 将用户ID存储到上下文中
		c.Set(AuthUserIDKey, userID)
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
