package mw

import (
	"context"
	"fmt"
	"log"

	"github.com/cloudwego/hertz/pkg/app"
	"github.com/cloudwego/hertz/pkg/protocol/consts"
)

// Recovery 恢复中间件，捕获 panic
func Recovery() app.HandlerFunc {
	return func(ctx context.Context, c *app.RequestContext) {
		defer func() {
			if err := recover(); err != nil {
				log.Printf("❌ Panic recovered: %v", err)
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
