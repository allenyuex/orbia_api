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
	ValidateImageURL(userID int64, req *upload.ValidateImageURLReq) (*upload.ValidateImageURLResp, error)
}

// uploadService 上传服务实现
type uploadService struct{}

// NewUploadService 创建上传服务实例
func NewUploadService() UploadService {
	return &uploadService{}
}

// GenerateUploadToken 生成上传token
func (s *uploadService) GenerateUploadToken(userID int64, req *upload.GenerateUploadTokenReq) (*upload.GenerateUploadTokenResp, error) {
	// 验证文件扩展名
	if !utils.ValidateImageExtension(req.FileExtension) {
		return &upload.GenerateUploadTokenResp{
			BaseResp: &common.BaseResp{
				Code:    400,
				Message: "unsupported image format",
			},
		}, nil
	}

	// 验证文件大小
	if req.FileSize != nil && !utils.ValidateFileSize(*req.FileSize) {
		return &upload.GenerateUploadTokenResp{
			BaseResp: &common.BaseResp{
				Code:    400,
				Message: "file size exceeds limit (10MB)",
			},
		}, nil
	}

	// 转换图片类型
	imageType := utils.ImageType(req.ImageType)

	// 生成图片路径
	imagePath := utils.GenerateImagePath(imageType, req.FileExtension)

	// 生成上传凭证
	token, err := utils.GenerateDirectUploadCredentials(imagePath)
	if err != nil {
		return &upload.GenerateUploadTokenResp{
			BaseResp: &common.BaseResp{
				Code:    500,
				Message: fmt.Sprintf("failed to generate upload token: %v", err),
			},
		}, nil
	}

	return &upload.GenerateUploadTokenResp{
		UploadURL:       token.UploadURL,
		AccessKeyID:     token.AccessKeyID,
		SecretAccessKey: token.SecretAccessKey,
		SessionToken:    token.SessionToken,
		Bucket:          token.Bucket,
		Key:             token.Key,
		PublicURL:       token.PublicURL,
		ExpiresIn:       token.ExpiresIn,
		BaseResp: &common.BaseResp{
			Code:    0,
			Message: "success",
		},
	}, nil
}

// ValidateImageURL 验证图片URL
func (s *uploadService) ValidateImageURL(userID int64, req *upload.ValidateImageURLReq) (*upload.ValidateImageURLResp, error) {
	// 基本URL格式验证
	isValid, errorMessage := utils.ValidateImageURL(req.ImageURL)
	if !isValid {
		return &upload.ValidateImageURLResp{
			IsValid:      false,
			ErrorMessage: &errorMessage,
			BaseResp: &common.BaseResp{
				Code:    0,
				Message: "success",
			},
		}, nil
	}

	// 检查图片是否真实存在
	if !utils.CheckImageExists(req.ImageURL) {
		errorMsg := "image does not exist or is not accessible"
		return &upload.ValidateImageURLResp{
			IsValid:      false,
			ErrorMessage: &errorMsg,
			BaseResp: &common.BaseResp{
				Code:    0,
				Message: "success",
			},
		}, nil
	}

	// 验证图片类型是否匹配路径
	imagePath := req.ImageURL[len(utils.GeneratePublicURL("")):]
	expectedType := utils.GetImageTypeFromPath(imagePath)
	if utils.ImageType(req.ImageType) != expectedType {
		errorMsg := "image type does not match the path"
		return &upload.ValidateImageURLResp{
			IsValid:      false,
			ErrorMessage: &errorMsg,
			BaseResp: &common.BaseResp{
				Code:    0,
				Message: "success",
			},
		}, nil
	}

	return &upload.ValidateImageURLResp{
		IsValid: true,
		BaseResp: &common.BaseResp{
			Code:    0,
			Message: "success",
		},
	}, nil
}