package mysql

import (
	"time"

	"gorm.io/gorm"
)

// VerificationCode 验证码模型
type VerificationCode struct {
	ID        int64      `gorm:"primaryKey;autoIncrement;column:id" json:"id"`
	Email     string     `gorm:"column:email;size:255;not null" json:"email"`
	Code      string     `gorm:"column:code;size:10;not null" json:"code"`
	CodeType  string     `gorm:"column:code_type;type:enum('login','register','reset_password');not null;default:'login'" json:"code_type"`
	Status    string     `gorm:"column:status;type:enum('unused','used','expired');not null;default:'unused'" json:"status"`
	UsedAt    *time.Time `gorm:"column:used_at" json:"used_at"`
	ExpiresAt time.Time  `gorm:"column:expires_at;not null" json:"expires_at"`
	CreatedAt time.Time  `gorm:"column:created_at;autoCreateTime" json:"created_at"`
	UpdatedAt time.Time  `gorm:"column:updated_at;autoUpdateTime" json:"updated_at"`
}

// TableName 指定表名
func (VerificationCode) TableName() string {
	return "orbia_verification_code"
}

// VerificationCodeRepository 验证码仓储接口
type VerificationCodeRepository interface {
	CreateVerificationCode(code *VerificationCode) error
	GetValidVerificationCode(email, code, codeType string) (*VerificationCode, error)
	MarkAsUsed(id int64) error
	CleanExpiredCodes() error
}

// verificationCodeRepository 验证码仓储实现
type verificationCodeRepository struct {
	db *gorm.DB
}

// NewVerificationCodeRepository 创建验证码仓储实例
func NewVerificationCodeRepository(db *gorm.DB) VerificationCodeRepository {
	return &verificationCodeRepository{db: db}
}

// CreateVerificationCode 创建验证码
func (r *verificationCodeRepository) CreateVerificationCode(code *VerificationCode) error {
	return r.db.Create(code).Error
}

// GetValidVerificationCode 获取有效的验证码
func (r *verificationCodeRepository) GetValidVerificationCode(email, code, codeType string) (*VerificationCode, error) {
	var verificationCode VerificationCode
	err := r.db.Where("email = ? AND code = ? AND code_type = ? AND status = ? AND expires_at > ?",
		email, code, codeType, "unused", time.Now()).
		Order("created_at DESC").
		First(&verificationCode).Error
	if err != nil {
		return nil, err
	}
	return &verificationCode, nil
}

// MarkAsUsed 标记验证码为已使用
func (r *verificationCodeRepository) MarkAsUsed(id int64) error {
	now := time.Now()
	return r.db.Model(&VerificationCode{}).
		Where("id = ?", id).
		Updates(map[string]interface{}{
			"status":  "used",
			"used_at": &now,
		}).Error
}

// CleanExpiredCodes 清理过期的验证码
func (r *verificationCodeRepository) CleanExpiredCodes() error {
	return r.db.Model(&VerificationCode{}).
		Where("status = ? AND expires_at < ?", "unused", time.Now()).
		Update("status", "expired").Error
}
