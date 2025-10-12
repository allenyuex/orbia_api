package utils

import (
	"crypto/rand"
	"fmt"
	"net/http"
	"path/filepath"
	"strings"
	"time"

	"orbia_api/biz/infra/config"
)

// SupportedImageExtensions 支持的图片格式
var SupportedImageExtensions = map[string]bool{
	".jpg":  true,
	".jpeg": true,
	".png":  true,
	".gif":  true,
	".svg":  true,
}

// ImageType 图片类型
type ImageType int

const (
	ImageTypeAvatar   ImageType = 1 // 头像
	ImageTypeTeamIcon ImageType = 2 // 团队图标
)

// ValidateImageExtension 验证图片扩展名
func ValidateImageExtension(extension string) bool {
	extension = strings.ToLower(extension)
	if !strings.HasPrefix(extension, ".") {
		extension = "." + extension
	}
	return SupportedImageExtensions[extension]
}

// ValidateFileSize 验证文件大小
func ValidateFileSize(size int64) bool {
	cfg := config.GlobalConfig.R2
	return size > 0 && size <= cfg.MaxFileSize
}

// GenerateImagePath 生成图片路径
func GenerateImagePath(imageType ImageType, extension string) string {
	// 确保扩展名格式正确
	extension = strings.ToLower(extension)
	if !strings.HasPrefix(extension, ".") {
		extension = "." + extension
	}

	// 生成随机文件名
	timestamp := time.Now().Unix()
	randomBytes := make([]byte, 8)
	rand.Read(randomBytes)
	filename := fmt.Sprintf("%d_%x%s", timestamp, randomBytes, extension)

	// 根据图片类型生成路径
	var basePath string
	switch imageType {
	case ImageTypeAvatar:
		basePath = "avatars"
	case ImageTypeTeamIcon:
		basePath = "team-icons"
	default:
		basePath = "images"
	}

	// 按年月分目录
	now := time.Now()
	yearMonth := now.Format("2006/01")
	
	return fmt.Sprintf("%s/%s/%s", basePath, yearMonth, filename)
}

// GeneratePublicURL 生成公开访问URL
func GeneratePublicURL(imagePath string) string {
	cfg := config.GlobalConfig.R2
	return fmt.Sprintf("%s/%s", strings.TrimRight(cfg.PublicURL, "/"), imagePath)
}

// ValidateImageURL 验证图片URL是否有效
func ValidateImageURL(imageURL string) (bool, string) {
	if imageURL == "" {
		return false, "image URL cannot be empty"
	}

	// 使用通用URL校验工具检查URL前缀
	isValid, errorMsg := ValidateURLPrefix(imageURL)
	if !isValid {
		return false, errorMsg
	}

	// 如果是白名单URL，直接返回成功（如GitHub等外部资源）
	if IsWhitelistedURL(imageURL) {
		return true, ""
	}

	// 对于配置的公共URL，需要进一步验证路径和扩展名
	cfg := config.GlobalConfig.R2
	expectedPrefix := strings.TrimRight(cfg.PublicURL, "/") + "/"
	
	if strings.HasPrefix(imageURL, expectedPrefix) {
		// 提取路径部分
		imagePath := strings.TrimPrefix(imageURL, expectedPrefix)
		if imagePath == "" {
			return false, "invalid image path"
		}

		// 验证路径格式
		if !isValidImagePath(imagePath) {
			return false, "invalid image path format"
		}

		// 验证扩展名
		extension := filepath.Ext(imagePath)
		if !ValidateImageExtension(extension) {
			return false, "unsupported image format"
		}
	}

	return true, ""
}

// CheckImageExists 检查图片是否存在（通过HTTP HEAD请求）
func CheckImageExists(imageURL string) bool {
	client := &http.Client{
		Timeout: 10 * time.Second,
	}

	req, err := http.NewRequest("HEAD", imageURL, nil)
	if err != nil {
		return false
	}

	resp, err := client.Do(req)
	if err != nil {
		return false
	}
	defer resp.Body.Close()

	return resp.StatusCode == http.StatusOK
}

// isValidImagePath 验证图片路径格式
func isValidImagePath(path string) bool {
	// 基本格式检查：不能包含..或其他危险字符
	if strings.Contains(path, "..") || strings.Contains(path, "//") {
		return false
	}

	// 检查是否符合预期的路径格式
	parts := strings.Split(path, "/")
	if len(parts) < 3 {
		return false
	}

	// 第一部分应该是类型目录
	validTypes := map[string]bool{
		"avatars":    true,
		"team-icons": true,
		"images":     true,
	}

	if !validTypes[parts[0]] {
		return false
	}

	return true
}

// GetImageTypeFromPath 从路径获取图片类型
func GetImageTypeFromPath(imagePath string) ImageType {
	if strings.HasPrefix(imagePath, "avatars/") {
		return ImageTypeAvatar
	}
	if strings.HasPrefix(imagePath, "team-icons/") {
		return ImageTypeTeamIcon
	}
	return ImageTypeAvatar // 默认返回头像类型
}