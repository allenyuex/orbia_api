package handler

import (
	"log"

	"orbia_api/biz/handler/admin"
	"orbia_api/biz/handler/auth"
	"orbia_api/biz/handler/campaign"
	"orbia_api/biz/handler/conversation"
	"orbia_api/biz/handler/dashboard"
	"orbia_api/biz/handler/dictionary"
	"orbia_api/biz/handler/kol"
	"orbia_api/biz/handler/payment_setting"
	"orbia_api/biz/handler/recharge_order"
	"orbia_api/biz/handler/team"
	"orbia_api/biz/handler/user"
	"orbia_api/biz/handler/wallet"
	adOrderService "orbia_api/biz/service/ad_order"
	kolOrderService "orbia_api/biz/service/kol_order"
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

	kolOrderService.InitKolOrderService()
	log.Println("  ✅ KOL Order service initialized")

	adOrderService.InitAdOrderService()
	log.Println("  ✅ Ad Order service initialized")

	admin.InitAdminService()
	log.Println("  ✅ Admin service initialized")

	dictionary.InitDictionaryService()
	log.Println("  ✅ Dictionary service initialized")

	payment_setting.InitPaymentSettingService()
	log.Println("  ✅ Payment Setting service initialized")

	recharge_order.InitRechargeOrderHandler()
	log.Println("  ✅ Recharge Order service initialized")

	conversation.InitConversationService()
	log.Println("  ✅ Conversation service initialized")

	campaign.InitCampaignService()
	log.Println("  ✅ Campaign service initialized")

	dashboard.InitDashboardService()
	log.Println("  ✅ Dashboard service initialized")

	log.Println("✅ All handler services initialized successfully")
}
