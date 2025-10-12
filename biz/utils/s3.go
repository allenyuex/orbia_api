package utils

import (
	"fmt"
	"time"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/credentials"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/s3"
	"orbia_api/biz/infra/config"
)

// S3UploadToken S3上传令牌信息
type S3UploadToken struct {
	UploadURL       string `json:"upload_url"`
	AccessKeyID     string `json:"access_key_id"`
	SecretAccessKey string `json:"secret_access_key"`
	SessionToken    string `json:"session_token"`
	Bucket          string `json:"bucket"`
	Key             string `json:"key"`
	PublicURL       string `json:"public_url"`
	ExpiresIn       int64  `json:"expires_in"`
}

// GenerateS3UploadToken 生成S3上传令牌
func GenerateS3UploadToken(imagePath string) (*S3UploadToken, error) {
	cfg := config.GlobalConfig.R2

	// 创建AWS会话
	sess, err := session.NewSession(&aws.Config{
		Region:      aws.String("auto"), // Cloudflare R2使用auto region
		Endpoint:    aws.String(cfg.Endpoint),
		Credentials: credentials.NewStaticCredentials(cfg.AccessKey, cfg.SecretKey, ""),
		S3ForcePathStyle: aws.Bool(true), // 强制使用路径样式
	})
	if err != nil {
		return nil, fmt.Errorf("failed to create AWS session: %v", err)
	}

	// 创建S3服务客户端
	svc := s3.New(sess)

	// 设置过期时间
	expiration := time.Duration(cfg.UploadTokenExpireMinutes) * time.Minute

	// 生成预签名PUT URL
	req, _ := svc.PutObjectRequest(&s3.PutObjectInput{
		Bucket: aws.String(cfg.Bucket),
		Key:    aws.String(imagePath),
	})

	uploadURL, err := req.Presign(expiration)
	if err != nil {
		return nil, fmt.Errorf("failed to generate presigned URL: %v", err)
	}

	// 生成公开访问URL
	publicURL := GeneratePublicURL(imagePath)

	return &S3UploadToken{
		UploadURL:       uploadURL,
		AccessKeyID:     cfg.AccessKey,
		SecretAccessKey: cfg.SecretKey,
		SessionToken:    "", // Cloudflare R2不使用session token
		Bucket:          cfg.Bucket,
		Key:             imagePath,
		PublicURL:       publicURL,
		ExpiresIn:       int64(expiration.Seconds()),
	}, nil
}

// GenerateDirectUploadCredentials 生成直接上传凭证（用于前端直接上传）
func GenerateDirectUploadCredentials(imagePath string) (*S3UploadToken, error) {
	cfg := config.GlobalConfig.R2

	// 生成公开访问URL
	publicURL := GeneratePublicURL(imagePath)

	// 计算过期时间
	expiresIn := int64(cfg.UploadTokenExpireMinutes * 60)

	return &S3UploadToken{
		UploadURL:       cfg.Endpoint + "/" + cfg.Bucket,
		AccessKeyID:     cfg.AccessKey,
		SecretAccessKey: cfg.SecretKey,
		SessionToken:    "",
		Bucket:          cfg.Bucket,
		Key:             imagePath,
		PublicURL:       publicURL,
		ExpiresIn:       expiresIn,
	}, nil
}