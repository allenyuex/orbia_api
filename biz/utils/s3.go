package utils

import (
	"fmt"
	"mime"
	"path/filepath"
	"time"

	"orbia_api/biz/infra/config"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/credentials"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/s3"
)

// S3UploadToken S3上传令牌信息
// 用于前端直接上传到 Cloudflare R2
type S3UploadToken struct {
	UploadURL       string            `json:"upload_url"`        // 预签名上传URL (用于PUT请求)
	AccessKeyID     string            `json:"access_key_id"`     // 访问密钥ID (可选，用于SDK上传)
	SecretAccessKey string            `json:"secret_access_key"` // 访问密钥 (可选，用于SDK上传)
	SessionToken    string            `json:"session_token"`     // 会话令牌 (R2不使用)
	Bucket          string            `json:"bucket"`            // 存储桶名称
	Key             string            `json:"key"`               // 对象键（文件路径）
	PublicURL       string            `json:"public_url"`        // 上传成功后的公开访问URL
	ExpiresIn       int64             `json:"expires_in"`        // 过期时间（秒）
	Headers         map[string]string `json:"headers"`           // 必需的请求头
}

// getS3Client 创建并返回配置好的 S3 客户端
func getS3Client() (*s3.S3, error) {
	cfg := config.GlobalConfig.R2

	// 创建 AWS 会话配置
	// Cloudflare R2 兼容 S3 API
	sess, err := session.NewSession(&aws.Config{
		Region:           aws.String("auto"), // R2 使用 "auto" 作为 region
		Endpoint:         aws.String(cfg.Endpoint),
		Credentials:      credentials.NewStaticCredentials(cfg.AccessKey, cfg.SecretKey, ""),
		S3ForcePathStyle: aws.Bool(true), // R2 需要使用路径样式访问
	})
	if err != nil {
		return nil, fmt.Errorf("failed to create AWS session: %w", err)
	}

	return s3.New(sess), nil
}

// getContentType 根据文件扩展名获取 Content-Type
func getContentType(fileExtension string) string {
	// 确保扩展名以点开头
	if fileExtension != "" && fileExtension[0] != '.' {
		fileExtension = "." + fileExtension
	}

	contentType := mime.TypeByExtension(fileExtension)
	if contentType == "" {
		// 默认使用 application/octet-stream
		contentType = "application/octet-stream"
	}
	return contentType
}

// GenerateS3UploadToken 生成 S3 上传令牌（推荐用于前端直接上传）
// 这个方法生成一个预签名的 PUT URL，前端可以直接使用这个 URL 上传文件
// 最佳实践：
// - 使用预签名 URL 避免在前端暴露凭证
// - 设置适当的过期时间
// - 添加 Content-Type 和 Content-Length 限制提高安全性
func GenerateS3UploadToken(imagePath string, fileSize *int64) (*S3UploadToken, error) {
	cfg := config.GlobalConfig.R2
	svc, err := getS3Client()
	if err != nil {
		return nil, err
	}

	// 设置过期时间
	expiration := time.Duration(cfg.UploadTokenExpireMinutes) * time.Minute

	// 获取文件类型
	fileExtension := filepath.Ext(imagePath)
	contentType := getContentType(fileExtension)

	// 创建 PutObject 请求
	putInput := &s3.PutObjectInput{
		Bucket:      aws.String(cfg.Bucket),
		Key:         aws.String(imagePath),
		ContentType: aws.String(contentType),
	}

	// 如果提供了文件大小，添加 Content-Length 限制
	if fileSize != nil && *fileSize > 0 {
		putInput.ContentLength = aws.Int64(*fileSize)
	}

	// 生成预签名 PUT URL
	req, _ := svc.PutObjectRequest(putInput)
	uploadURL, err := req.Presign(expiration)
	if err != nil {
		return nil, fmt.Errorf("failed to generate presigned URL: %w", err)
	}

	// 生成公开访问 URL
	publicURL := GeneratePublicURL(imagePath)

	// 准备必需的请求头
	headers := map[string]string{
		"Content-Type": contentType,
	}
	if fileSize != nil && *fileSize > 0 {
		headers["Content-Length"] = fmt.Sprintf("%d", *fileSize)
	}

	return &S3UploadToken{
		UploadURL:       uploadURL,
		AccessKeyID:     "", // 使用预签名 URL 时不需要暴露凭证
		SecretAccessKey: "", // 使用预签名 URL 时不需要暴露凭证
		SessionToken:    "", // Cloudflare R2 不使用 session token
		Bucket:          cfg.Bucket,
		Key:             imagePath,
		PublicURL:       publicURL,
		ExpiresIn:       int64(expiration.Seconds()),
		Headers:         headers,
	}, nil
}

// GenerateDirectUploadCredentials 生成直接上传凭证（用于需要 SDK 的场景）
// 注意：这种方式会暴露凭证给前端，不如预签名 URL 安全
// 推荐使用 GenerateS3UploadToken 方法
func GenerateDirectUploadCredentials(imagePath string) (*S3UploadToken, error) {
	cfg := config.GlobalConfig.R2

	// 生成公开访问 URL
	publicURL := GeneratePublicURL(imagePath)

	// 计算过期时间
	expiresIn := int64(cfg.UploadTokenExpireMinutes * 60)

	// 获取文件类型
	fileExtension := filepath.Ext(imagePath)
	contentType := getContentType(fileExtension)

	return &S3UploadToken{
		UploadURL:       cfg.Endpoint,
		AccessKeyID:     cfg.AccessKey,
		SecretAccessKey: cfg.SecretKey,
		SessionToken:    "",
		Bucket:          cfg.Bucket,
		Key:             imagePath,
		PublicURL:       publicURL,
		ExpiresIn:       expiresIn,
		Headers: map[string]string{
			"Content-Type": contentType,
		},
	}, nil
}
