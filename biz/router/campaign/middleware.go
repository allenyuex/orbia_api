package campaign

import (
	"orbia_api/biz/consts"
	"orbia_api/biz/mw"

	"github.com/cloudwego/hertz/pkg/app"
)

func rootMw() []app.HandlerFunc {
	return nil
}

func _adminMw() []app.HandlerFunc {
	return nil
}

func _campaignMw() []app.HandlerFunc {
	return nil
}

func _adminlistcampaignsMw() []app.HandlerFunc {
	// 需要JWT认证，且仅管理员可访问
	return []app.HandlerFunc{mw.AuthMiddleware(consts.RoleAdmin)}
}

func _adminupdatecampaignstatusMw() []app.HandlerFunc {
	// 需要JWT认证，且仅管理员可访问
	return []app.HandlerFunc{mw.AuthMiddleware(consts.RoleAdmin)}
}

func _campaign0Mw() []app.HandlerFunc {
	return nil
}

func _createcampaignMw() []app.HandlerFunc {
	// 需要JWT认证，普通用户可访问
	return []app.HandlerFunc{mw.AuthMiddleware(consts.RoleNormal, consts.RoleAdmin)}
}

func _getcampaignMw() []app.HandlerFunc {
	// 需要JWT认证，普通用户可访问
	return []app.HandlerFunc{mw.AuthMiddleware(consts.RoleNormal, consts.RoleAdmin)}
}

func _listcampaignsMw() []app.HandlerFunc {
	// 需要JWT认证，普通用户可访问
	return []app.HandlerFunc{mw.AuthMiddleware(consts.RoleNormal, consts.RoleAdmin)}
}

func _updatecampaignstatusMw() []app.HandlerFunc {
	// 需要JWT认证，普通用户可访问
	return []app.HandlerFunc{mw.AuthMiddleware(consts.RoleNormal, consts.RoleAdmin)}
}

func _updatecampaignMw() []app.HandlerFunc {
	// 需要JWT认证，普通用户可访问
	return []app.HandlerFunc{mw.AuthMiddleware(consts.RoleNormal, consts.RoleAdmin)}
}

func _apiMw() []app.HandlerFunc {
	// your code...
	return nil
}

func _v1Mw() []app.HandlerFunc {
	// your code...
	return nil
}
