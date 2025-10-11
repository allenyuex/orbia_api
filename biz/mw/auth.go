package mw

import (
	"context"
	"net/http"
	"strings"

	"github.com/cloudwego/hertz/pkg/app"
	"github.com/cloudwego/hertz/pkg/common/hlog"
	"orbia_api/biz/utils"
)

const (
	AuthUserIDKey = "auth_user_id"
)

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