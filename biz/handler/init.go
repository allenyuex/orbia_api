package handler

import (
	"log"

	"orbia_api/biz/handler/admin"
	"orbia_api/biz/handler/auth"
	"orbia_api/biz/handler/kol"
	"orbia_api/biz/handler/order"
	"orbia_api/biz/handler/team"
	"orbia_api/biz/handler/user"
	"orbia_api/biz/handler/wallet"
)

// InitAllServices 初始化所有handler服务
func InitAllServices() {
	log.Println("🚀 Initializing all handler services...")

	// 初始化各个服务
	auth.InitAuthService()
	log.Println("  ✅ Auth service initialized")

	kol.InitKolService()
	log.Println("  ✅ KOL service initialized")

	user.InitUserService()
	log.Println("  ✅ User service initialized")

	team.InitTeamService()
	log.Println("  ✅ Team service initialized")

	wallet.InitWalletHandler()
	log.Println("  ✅ Wallet service initialized")

	order.InitOrderService()
	log.Println("  ✅ Order service initialized")

	admin.InitAdminService()
	log.Println("  ✅ Admin service initialized")

	log.Println("✅ All handler services initialized successfully")
}
