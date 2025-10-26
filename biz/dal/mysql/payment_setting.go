package mysql

import (
	"orbia_api/biz/dal/model"

	"gorm.io/gorm"
)

// PaymentSettingRepository 收款钱包设置仓储接口
type PaymentSettingRepository interface {
	CreatePaymentSetting(setting *model.OrbiaPaymentSetting) error
	UpdatePaymentSetting(setting *model.OrbiaPaymentSetting) error
	DeletePaymentSetting(id int64) error
	GetPaymentSettingByID(id int64) (*model.OrbiaPaymentSetting, error)
	GetPaymentSettings(network string, status *int32, offset int, limit int) ([]*model.OrbiaPaymentSetting, int64, error)
	GetActivePaymentSettings(network string) ([]*model.OrbiaPaymentSetting, error)
}

// paymentSettingRepository 收款钱包设置仓储实现
type paymentSettingRepository struct {
	db *gorm.DB
}

// NewPaymentSettingRepository 创建收款钱包设置仓储实例
func NewPaymentSettingRepository(db *gorm.DB) PaymentSettingRepository {
	return &paymentSettingRepository{db: db}
}

// CreatePaymentSetting 创建收款钱包设置
func (r *paymentSettingRepository) CreatePaymentSetting(setting *model.OrbiaPaymentSetting) error {
	return r.db.Create(setting).Error
}

// UpdatePaymentSetting 更新收款钱包设置
func (r *paymentSettingRepository) UpdatePaymentSetting(setting *model.OrbiaPaymentSetting) error {
	return r.db.Save(setting).Error
}

// DeletePaymentSetting 删除收款钱包设置（软删除）
func (r *paymentSettingRepository) DeletePaymentSetting(id int64) error {
	return r.db.Delete(&model.OrbiaPaymentSetting{}, id).Error
}

// GetPaymentSettingByID 根据ID获取收款钱包设置
func (r *paymentSettingRepository) GetPaymentSettingByID(id int64) (*model.OrbiaPaymentSetting, error) {
	var setting model.OrbiaPaymentSetting
	err := r.db.Where("id = ?", id).First(&setting).Error
	if err != nil {
		return nil, err
	}
	return &setting, nil
}

// GetPaymentSettings 获取收款钱包设置列表
func (r *paymentSettingRepository) GetPaymentSettings(network string, status *int32, offset int, limit int) ([]*model.OrbiaPaymentSetting, int64, error) {
	var settings []*model.OrbiaPaymentSetting
	var total int64

	query := r.db.Model(&model.OrbiaPaymentSetting{})

	// 网络筛选
	if network != "" {
		query = query.Where("network = ?", network)
	}

	// 状态筛选
	if status != nil {
		query = query.Where("status = ?", *status)
	}

	// 获取总数
	if err := query.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	// 获取分页数据
	err := query.Order("created_at DESC").
		Offset(offset).
		Limit(limit).
		Find(&settings).Error

	return settings, total, err
}

// GetActivePaymentSettings 获取启用的收款钱包设置列表
func (r *paymentSettingRepository) GetActivePaymentSettings(network string) ([]*model.OrbiaPaymentSetting, error) {
	var settings []*model.OrbiaPaymentSetting

	query := r.db.Model(&model.OrbiaPaymentSetting{}).Where("status = ?", 1)

	// 网络筛选
	if network != "" {
		query = query.Where("network = ?", network)
	}

	err := query.Order("created_at DESC").Find(&settings).Error
	return settings, err
}
