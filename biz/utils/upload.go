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

// ValidateFileExtension 验证文件扩展名是否支持
func ValidateFileExtension(extension string) bool {
	extension = NormalizeExtension(extension)

	cfg := config.GlobalConfig.R2
	_, exists := cfg.AllowedExtensions[extension]
	return exists
}

// ValidateFileSize 验证文件大小
func ValidateFileSize(extension string, size int64) bool {
	if size <= 0 {
		return false
	}

	extension = NormalizeExtension(extension)

	cfg := config.GlobalConfig.R2

	// 获取该扩展名的配置
	if extConfig, exists := cfg.AllowedExtensions[extension]; exists {
		return size <= extConfig.MaxSize
	}

	// 如果没有配置，使用默认大小限制
	return size <= cfg.MaxFileSize
}

// NormalizeExtension 规范化文件扩展名（统一转为小写，确保有点号前缀）
func NormalizeExtension(extension string) string {
	extension = strings.ToLower(strings.TrimSpace(extension))
	if extension != "" && !strings.HasPrefix(extension, ".") {
		extension = "." + extension
	}
	return extension
}

// GenerateFilePath 生成文件路径
// 路径和文件名完全由后端控制，根据文件扩展名自动选择存储目录
func GenerateFilePath(extension string) string {
	// 规范化扩展名（统一小写）
	extension = NormalizeExtension(extension)

	// 生成随机文件名
	timestamp := time.Now().Unix()
	randomBytes := make([]byte, 8)
	rand.Read(randomBytes)
	filename := fmt.Sprintf("%d_%x%s", timestamp, randomBytes, extension)

	// 根据配置获取该扩展名的默认存储路径
	cfg := config.GlobalConfig.R2
	var basePath string
	if extConfig, exists := cfg.AllowedExtensions[extension]; exists {
		basePath = extConfig.DefaultPath
	} else {
		// 如果配置中没有该扩展名，使用通用的 files 目录
		basePath = "files"
	}

	// 按年月分目录
	now := time.Now()
	yearMonth := now.Format("2006/01")

	return fmt.Sprintf("%s/%s/%s", basePath, yearMonth, filename)
}

// GeneratePublicURL 生成公开访问URL
func GeneratePublicURL(filePath string) string {
	cfg := config.GlobalConfig.R2
	return fmt.Sprintf("%s/%s", strings.TrimRight(cfg.PublicURL, "/"), filePath)
}

// ValidateFileURL 验证文件URL是否有效
func ValidateFileURL(fileURL string) (bool, string) {
	if fileURL == "" {
		return false, "file URL cannot be empty"
	}

	// 使用通用URL校验工具检查URL前缀
	isValid, errorMsg := ValidateURLPrefix(fileURL)
	if !isValid {
		return false, errorMsg
	}

	// 如果是白名单URL，直接返回成功（如GitHub等外部资源）
	if IsWhitelistedURL(fileURL) {
		return true, ""
	}

	// 对于配置的公共URL，需要进一步验证路径和扩展名
	cfg := config.GlobalConfig.R2
	expectedPrefix := strings.TrimRight(cfg.PublicURL, "/") + "/"

	if strings.HasPrefix(fileURL, expectedPrefix) {
		// 提取路径部分
		filePath := strings.TrimPrefix(fileURL, expectedPrefix)
		if filePath == "" {
			return false, "invalid file path"
		}

		// 验证路径格式
		if !isValidFilePath(filePath) {
			return false, "invalid file path format"
		}

		// 验证扩展名
		extension := filepath.Ext(filePath)
		if !ValidateFileExtension(extension) {
			return false, "unsupported file format"
		}
	}

	return true, ""
}

// CheckFileExists 检查文件是否存在（通过HTTP HEAD请求）
func CheckFileExists(fileURL string) bool {
	client := &http.Client{
		Timeout: 10 * time.Second,
	}

	req, err := http.NewRequest("HEAD", fileURL, nil)
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

// isValidFilePath 验证文件路径格式
func isValidFilePath(path string) bool {
	// 基本格式检查：不能包含..或其他危险字符
	if strings.Contains(path, "..") || strings.Contains(path, "//") {
		return false
	}

	// 检查是否符合预期的路径格式
	parts := strings.Split(path, "/")
	if len(parts) < 3 {
		return false
	}

	// 第一部分应该是合法的目录名
	validPaths := map[string]bool{
		"avatars":          true,
		"team-icons":       true,
		"kol-video-covers": true,
		"images":           true,
		"documents":        true,
		"videos":           true,
		"attachments":      true,
		"files":            true,
	}

	if !validPaths[parts[0]] {
		return false
	}

	return true
}

// GetFileCategory 根据扩展名获取文件分类
func GetFileCategory(extension string) string {
	extension = strings.ToLower(extension)
	if !strings.HasPrefix(extension, ".") {
		extension = "." + extension
	}

	// 图片类型
	imageExts := map[string]bool{
		".jpg":  true,
		".jpeg": true,
		".png":  true,
		".gif":  true,
		".webp": true,
		".svg":  true,
	}
	if imageExts[extension] {
		return "image"
	}

	// 文档类型
	docExts := map[string]bool{
		".pdf":  true,
		".doc":  true,
		".docx": true,
		".xls":  true,
		".xlsx": true,
		".ppt":  true,
		".pptx": true,
		".txt":  true,
	}
	if docExts[extension] {
		return "document"
	}

	// 视频类型
	videoExts := map[string]bool{
		".mp4": true,
		".mov": true,
		".avi": true,
		".mkv": true,
	}
	if videoExts[extension] {
		return "video"
	}

	return "other"
}
