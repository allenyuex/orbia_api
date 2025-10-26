package mysql

import (
	"fmt"

	"orbia_api/biz/dal/model"

	"gorm.io/gorm"
)

// RechargeOrderRepository 充值订单仓库接口
type RechargeOrderRepository interface {
	CreateRechargeOrder(order *model.OrbiaRechargeOrder) error
	GetRechargeOrderByID(orderID string) (*model.OrbiaRechargeOrder, error)
	GetRechargeOrderByOrderID(orderID string) (*model.OrbiaRechargeOrder, error)
	UpdateRechargeOrder(order *model.OrbiaRechargeOrder) error
	GetRechargeOrdersByUserID(userID int64, status *string, page, pageSize int) ([]*model.OrbiaRechargeOrder, int64, error)
	GetAllRechargeOrders(userID *int64, status, paymentType *string, page, pageSize int) ([]*model.OrbiaRechargeOrder, int64, error)
}

// rechargeOrderRepository 充值订单仓库实现
type rechargeOrderRepository struct {
	db *gorm.DB
}

// NewRechargeOrderRepository 创建充值订单仓库实例
func NewRechargeOrderRepository(db *gorm.DB) RechargeOrderRepository {
	return &rechargeOrderRepository{db: db}
}

// CreateRechargeOrder 创建充值订单
func (r *rechargeOrderRepository) CreateRechargeOrder(order *model.OrbiaRechargeOrder) error {
	return r.db.Create(order).Error
}

// GetRechargeOrderByID 根据内部ID获取充值订单
func (r *rechargeOrderRepository) GetRechargeOrderByID(orderID string) (*model.OrbiaRechargeOrder, error) {
	var order model.OrbiaRechargeOrder
	err := r.db.Where("id = ?", orderID).First(&order).Error
	if err != nil {
		return nil, err
	}
	return &order, nil
}

// GetRechargeOrderByOrderID 根据订单ID获取充值订单
func (r *rechargeOrderRepository) GetRechargeOrderByOrderID(orderID string) (*model.OrbiaRechargeOrder, error) {
	var order model.OrbiaRechargeOrder
	err := r.db.Where("order_id = ?", orderID).First(&order).Error
	if err != nil {
		return nil, err
	}
	return &order, nil
}

// UpdateRechargeOrder 更新充值订单
func (r *rechargeOrderRepository) UpdateRechargeOrder(order *model.OrbiaRechargeOrder) error {
	return r.db.Save(order).Error
}

// GetRechargeOrdersByUserID 根据用户ID获取充值订单列表
func (r *rechargeOrderRepository) GetRechargeOrdersByUserID(userID int64, status *string, page, pageSize int) ([]*model.OrbiaRechargeOrder, int64, error) {
	var orders []*model.OrbiaRechargeOrder
	var total int64

	query := r.db.Model(&model.OrbiaRechargeOrder{}).Where("user_id = ?", userID)

	if status != nil && *status != "" {
		query = query.Where("status = ?", *status)
	}

	// 获取总数
	if err := query.Count(&total).Error; err != nil {
		return nil, 0, fmt.Errorf("failed to count recharge orders: %v", err)
	}

	// 分页查询
	offset := (page - 1) * pageSize
	if err := query.Order("created_at DESC").Limit(pageSize).Offset(offset).Find(&orders).Error; err != nil {
		return nil, 0, fmt.Errorf("failed to query recharge orders: %v", err)
	}

	return orders, total, nil
}

// GetAllRechargeOrders 获取所有充值订单列表（管理员）
func (r *rechargeOrderRepository) GetAllRechargeOrders(userID *int64, status, paymentType *string, page, pageSize int) ([]*model.OrbiaRechargeOrder, int64, error) {
	var orders []*model.OrbiaRechargeOrder
	var total int64

	query := r.db.Model(&model.OrbiaRechargeOrder{})

	if userID != nil && *userID > 0 {
		query = query.Where("user_id = ?", *userID)
	}

	if status != nil && *status != "" {
		query = query.Where("status = ?", *status)
	}

	if paymentType != nil && *paymentType != "" {
		query = query.Where("payment_type = ?", *paymentType)
	}

	// 获取总数
	if err := query.Count(&total).Error; err != nil {
		return nil, 0, fmt.Errorf("failed to count all recharge orders: %v", err)
	}

	// 分页查询，按创建时间倒序
	offset := (page - 1) * pageSize
	if err := query.Order("created_at DESC").Limit(pageSize).Offset(offset).Find(&orders).Error; err != nil {
		return nil, 0, fmt.Errorf("failed to query all recharge orders: %v", err)
	}

	return orders, total, nil
}
