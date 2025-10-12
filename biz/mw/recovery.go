package mw

import (
	"context"
	"fmt"

	"github.com/cloudwego/hertz/pkg/app"
	"github.com/cloudwego/hertz/pkg/protocol/consts"
	"orbia_api/biz/utils"
)

// Recovery 恢复中间件，捕获 panic
func Recovery() app.HandlerFunc {
	return func(ctx context.Context, c *app.RequestContext) {
		defer func() {
			if err := recover(); err != nil {
				// 使用新的日志工具记录详细的panic信息
				utils.LogPanic(err, "Panic recovered in HTTP handler")
				
				c.JSON(consts.StatusInternalServerError, map[string]interface{}{
					"code":    500,
					"message": fmt.Sprintf("Internal Server Error: %v", err),
				})
				c.Abort()
			}
		}()
		c.Next(ctx)
	}
}
