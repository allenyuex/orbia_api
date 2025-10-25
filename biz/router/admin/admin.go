package admin

import (
	"orbia_api/biz/handler/admin"
	"orbia_api/biz/mw"

	"github.com/cloudwego/hertz/pkg/app/server"
)

// Register 注册admin路由
func Register(h *server.Hertz) {
	adminGroup := h.Group("/api/v1/admin")

	// 应用鉴权中间件和管理员角色中间件
	adminGroup.Use(mw.AuthMiddleware())
	adminGroup.Use(mw.AdminOnlyMiddleware())

	// 用户管理
	adminGroup.POST("/users", admin.GetAllUsers)
	adminGroup.POST("/user/status", admin.SetUserStatus)

	// KOL管理
	adminGroup.POST("/kols", admin.GetAllKols)
	adminGroup.POST("/kol/review", admin.AdminReviewKol)

	// 团队管理
	adminGroup.POST("/teams", admin.GetAllTeams)
	adminGroup.POST("/team/:team_id/members", admin.GetTeamMembers)

	// 订单管理
	adminGroup.POST("/orders", admin.GetAllOrders)

	// 钱包管理
	adminGroup.POST("/user/:user_id/wallet", admin.GetUserWallet)
}
