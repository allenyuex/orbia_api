package mw

import (
	"time"

	"github.com/cloudwego/hertz/pkg/app"
	"github.com/hertz-contrib/cors"
)

// CORS 跨域中间件
func CORS() app.HandlerFunc {
	return cors.New(cors.Config{
		// 允许所有来源
		AllowOrigins: []string{"*"},

		// 允许的 HTTP 方法
		AllowMethods: []string{
			"GET", "POST", "PUT", "PATCH", "DELETE", "HEAD", "OPTIONS",
		},

		// 允许的请求头
		AllowHeaders: []string{
			"Origin", "Content-Length", "Content-Type",
			"Authorization",
			"X-Requested-With", "X-CSRF-Token",
		},

		// 暴露的响应头
		ExposeHeaders: []string{
			"Content-Length", "Access-Control-Allow-Origin",
			"Access-Control-Allow-Headers", "Content-Disposition",
		},

		// 允许携带认证信息
		AllowCredentials: true,

		// 预检请求的缓存时间
		MaxAge: 12 * time.Hour,
	})
}
