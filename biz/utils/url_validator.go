package utils

import (
	"strings"

	"orbia_api/biz/infra/config"
)

// URLWhitelist URL白名单配置
var URLWhitelist = []string{
	"https://raw.githubusercontent.com/",
}

// ValidateURLPrefix 校验URL是否符合允许的前缀
// 支持配置的公共URL和白名单URL
func ValidateURLPrefix(url string) (bool, string) {
	if url == "" {
		return false, "URL cannot be empty"
	}

	// 检查是否匹配配置的公共URL
	cfg := config.GlobalConfig.R2
	if cfg.PublicURL != "" {
		expectedPrefix := strings.TrimRight(cfg.PublicURL, "/") + "/"
		if strings.HasPrefix(url, expectedPrefix) {
			return true, ""
		}
	}

	// 检查是否匹配白名单中的URL
	for _, whitelistURL := range URLWhitelist {
		if strings.HasPrefix(url, whitelistURL) {
			return true, ""
		}
	}

	return false, "URL domain not allowed"
}

// AddURLToWhitelist 添加URL到白名单
func AddURLToWhitelist(url string) {
	// 确保URL以/结尾，便于前缀匹配
	if !strings.HasSuffix(url, "/") {
		url += "/"
	}

	// 检查是否已存在
	for _, existing := range URLWhitelist {
		if existing == url {
			return
		}
	}

	URLWhitelist = append(URLWhitelist, url)
}

// RemoveURLFromWhitelist 从白名单中移除URL
func RemoveURLFromWhitelist(url string) {
	// 确保URL以/结尾，便于匹配
	if !strings.HasSuffix(url, "/") {
		url += "/"
	}

	for i, existing := range URLWhitelist {
		if existing == url {
			URLWhitelist = append(URLWhitelist[:i], URLWhitelist[i+1:]...)
			return
		}
	}
}

// GetURLWhitelist 获取当前白名单
func GetURLWhitelist() []string {
	return URLWhitelist
}

// IsWhitelistedURL 检查URL是否在白名单中
func IsWhitelistedURL(url string) bool {
	for _, whitelistURL := range URLWhitelist {
		if strings.HasPrefix(url, whitelistURL) {
			return true
		}
	}
	return false
}
