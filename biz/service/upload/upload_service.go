package upload

import (
	"fmt"

	"orbia_api/biz/model/common"
	"orbia_api/biz/model/upload"
	"orbia_api/biz/utils"
)

// UploadService 上传服务接口
type UploadService interface {
	GenerateUploadToken(userID int64, req *upload.GenerateUploadTokenReq) (*upload.GenerateUploadTokenResp, error)
	ValidateFileURL(userID int64, req *upload.ValidateFileURLReq) (*upload.ValidateFileURLResp, error)
}

// uploadService 上传服务实现
type uploadService struct{}

// NewUploadService 创建上传服务实例
func NewUploadService() UploadService {
	return &uploadService{}
}

// GenerateUploadToken 生成上传token
// 使用预签名 URL 方式，这是推荐的最佳实践
//
// 文件路径和名称完全由后端控制：
// - 扩展名会统一转为小写（.PNG -> .png）
// - 根据扩展名自动选择存储目录（配置文件中的 default_path）
// - 文件名使用时间戳+随机数保证唯一性
// - 按年月自动分目录存储
//
// 前端使用方法：
// 1. 调用此接口获取 upload_url 和 headers
// 2. 使用 HTTP PUT 方法上传文件到 upload_url，并携带 headers 中的所有请求头
// 3. 上传成功后，使用 public_url 访问文件
func (s *uploadService) GenerateUploadToken(userID int64, req *upload.GenerateUploadTokenReq) (*upload.GenerateUploadTokenResp, error) {
	// 规范化扩展名（统一转为小写）
	normalizedExt := utils.NormalizeExtension(req.FileExtension)

	// 验证文件扩展名
	if !utils.ValidateFileExtension(normalizedExt) {
		return &upload.GenerateUploadTokenResp{
			BaseResp: &common.BaseResp{
				Code:    400,
				Message: fmt.Sprintf("unsupported file format: %s", normalizedExt),
			},
		}, nil
	}

	// 验证文件大小
	if req.FileSize != nil && !utils.ValidateFileSize(normalizedExt, *req.FileSize) {
		return &upload.GenerateUploadTokenResp{
			BaseResp: &common.BaseResp{
				Code:    400,
				Message: "file size exceeds limit",
			},
		}, nil
	}

	// 生成文件路径（完全由后端控制）
	filePath := utils.GenerateFilePath(normalizedExt)

	// 生成预签名上传 URL（推荐方式）
	// 这种方式不会在响应中暴露凭证，更安全
	token, err := utils.GenerateS3UploadToken(filePath, req.FileSize)
	if err != nil {
		return &upload.GenerateUploadTokenResp{
			BaseResp: &common.BaseResp{
				Code:    500,
				Message: fmt.Sprintf("failed to generate upload token: %v", err),
			},
		}, nil
	}

	return &upload.GenerateUploadTokenResp{
		UploadURL: token.UploadURL,
		PublicURL: token.PublicURL,
		ExpiresIn: token.ExpiresIn,
		Headers:   token.Headers,
		BaseResp: &common.BaseResp{
			Code:    0,
			Message: "success",
		},
	}, nil
}

// ValidateFileURL 验证文件URL
func (s *uploadService) ValidateFileURL(userID int64, req *upload.ValidateFileURLReq) (*upload.ValidateFileURLResp, error) {
	// 基本URL格式验证
	isValid, errorMessage := utils.ValidateFileURL(req.FileURL)
	if !isValid {
		return &upload.ValidateFileURLResp{
			IsValid:      false,
			ErrorMessage: &errorMessage,
			BaseResp: &common.BaseResp{
				Code:    0,
				Message: "success",
			},
		}, nil
	}

	// 检查文件是否真实存在
	if !utils.CheckFileExists(req.FileURL) {
		errorMsg := "file does not exist or is not accessible"
		return &upload.ValidateFileURLResp{
			IsValid:      false,
			ErrorMessage: &errorMsg,
			BaseResp: &common.BaseResp{
				Code:    0,
				Message: "success",
			},
		}, nil
	}

	return &upload.ValidateFileURLResp{
		IsValid: true,
		BaseResp: &common.BaseResp{
			Code:    0,
			Message: "success",
		},
	}, nil
}
