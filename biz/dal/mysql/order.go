package mysql

import (
	"time"

	"gorm.io/gorm"
)

// KolOrder KOL订单模型
type KolOrder struct {
	ID              int64          `gorm:"primaryKey;autoIncrement;column:id" json:"id"`
	OrderID         string         `gorm:"uniqueIndex;column:order_id;size:64;not null" json:"order_id"`
	UserID          int64          `gorm:"column:user_id;not null" json:"user_id"`
	TeamID          *int64         `gorm:"column:team_id" json:"team_id"`
	KolID           int64          `gorm:"column:kol_id;not null" json:"kol_id"`
	PlanID          int64          `gorm:"column:plan_id;not null" json:"plan_id"`
	PlanTitle       string         `gorm:"column:plan_title;size:200;not null" json:"plan_title"`
	PlanDescription *string        `gorm:"column:plan_description;type:text" json:"plan_description"`
	PlanPrice       float64        `gorm:"column:plan_price;type:decimal(10,2);not null" json:"plan_price"`
	PlanType        string         `gorm:"column:plan_type;size:20;not null" json:"plan_type"`
	Description     string         `gorm:"column:description;type:text;not null" json:"description"`
	Status          string         `gorm:"column:status;type:enum('pending','confirmed','in_progress','completed','cancelled','refunded');default:pending;not null" json:"status"`
	RejectReason    *string        `gorm:"column:reject_reason;type:text" json:"reject_reason"`
	ConfirmedAt     *time.Time     `gorm:"column:confirmed_at" json:"confirmed_at"`
	CompletedAt     *time.Time     `gorm:"column:completed_at" json:"completed_at"`
	CancelledAt     *time.Time     `gorm:"column:cancelled_at" json:"cancelled_at"`
	CreatedAt       time.Time      `gorm:"column:created_at;autoCreateTime" json:"created_at"`
	UpdatedAt       time.Time      `gorm:"column:updated_at;autoUpdateTime" json:"updated_at"`
	DeletedAt       gorm.DeletedAt `gorm:"index;column:deleted_at" json:"-"`
}

// TableName 指定表名
func (KolOrder) TableName() string {
	return "orbia_kol_order"
}

// OrderWithKolInfo 订单和KOL信息的联合查询结果
type OrderWithKolInfo struct {
	KolOrder
	KolDisplayName *string `json:"kol_display_name"`
	KolAvatarURL   *string `json:"kol_avatar_url"`
}

// OrderRepository 订单仓储接口
type OrderRepository interface {
	// 创建订单
	CreateOrder(order *KolOrder) error

	// 根据订单ID获取订单
	GetOrderByID(orderID string) (*KolOrder, error)

	// 根据订单ID获取订单（包含KOL信息）
	GetOrderWithKolInfo(orderID string) (*OrderWithKolInfo, error)

	// 获取用户的订单列表
	GetUserOrders(userID int64, status *string, offset, limit int) ([]*OrderWithKolInfo, int64, error)

	// 获取KOL收到的订单列表
	GetKolOrders(kolID int64, status *string, offset, limit int) ([]*OrderWithKolInfo, int64, error)

	// 获取团队的订单列表
	GetTeamOrders(teamID int64, status *string, offset, limit int) ([]*OrderWithKolInfo, int64, error)

	// 更新订单状态
	UpdateOrderStatus(orderID string, status string, reason *string) error

	// 更新订单
	UpdateOrder(order *KolOrder) error

	// 检查用户是否拥有订单
	IsOrderOwner(orderID string, userID int64) (bool, error)

	// 管理员功能
	GetAllOrders(keyword string, status string, offset int, limit int) ([]*OrderWithKolInfo, int64, error)
}

// orderRepository 订单仓储实现
type orderRepository struct {
	db *gorm.DB
}

// NewOrderRepository 创建订单仓储实例
func NewOrderRepository(db *gorm.DB) OrderRepository {
	return &orderRepository{db: db}
}

// CreateOrder 创建订单
func (r *orderRepository) CreateOrder(order *KolOrder) error {
	return r.db.Create(order).Error
}

// GetOrderByID 根据订单ID获取订单
func (r *orderRepository) GetOrderByID(orderID string) (*KolOrder, error) {
	var order KolOrder
	err := r.db.Where("order_id = ?", orderID).First(&order).Error
	if err != nil {
		return nil, err
	}
	return &order, nil
}

// GetOrderWithKolInfo 根据订单ID获取订单（包含KOL信息）
func (r *orderRepository) GetOrderWithKolInfo(orderID string) (*OrderWithKolInfo, error) {
	var result OrderWithKolInfo
	err := r.db.Table("orbia_kol_order").
		Select("orbia_kol_order.*, orbia_kol.display_name as kol_display_name, orbia_kol.avatar_url as kol_avatar_url").
		Joins("LEFT JOIN orbia_kol ON orbia_kol_order.kol_id = orbia_kol.id").
		Where("orbia_kol_order.order_id = ?", orderID).
		First(&result).Error
	if err != nil {
		return nil, err
	}
	return &result, nil
}

// GetUserOrders 获取用户的订单列表
func (r *orderRepository) GetUserOrders(userID int64, status *string, offset, limit int) ([]*OrderWithKolInfo, int64, error) {
	var orders []*OrderWithKolInfo
	var total int64

	query := r.db.Table("orbia_kol_order").
		Select("orbia_kol_order.*, orbia_kol.display_name as kol_display_name, orbia_kol.avatar_url as kol_avatar_url").
		Joins("LEFT JOIN orbia_kol ON orbia_kol_order.kol_id = orbia_kol.id").
		Where("orbia_kol_order.user_id = ?", userID)

	if status != nil && *status != "" {
		query = query.Where("orbia_kol_order.status = ?", *status)
	}

	// 获取总数
	if err := query.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	// 获取列表，按创建时间倒序
	if err := query.Order("orbia_kol_order.created_at DESC").
		Offset(offset).
		Limit(limit).
		Find(&orders).Error; err != nil {
		return nil, 0, err
	}

	return orders, total, nil
}

// GetKolOrders 获取KOL收到的订单列表
func (r *orderRepository) GetKolOrders(kolID int64, status *string, offset, limit int) ([]*OrderWithKolInfo, int64, error) {
	var orders []*OrderWithKolInfo
	var total int64

	query := r.db.Table("orbia_kol_order").
		Select("orbia_kol_order.*, orbia_kol.display_name as kol_display_name, orbia_kol.avatar_url as kol_avatar_url").
		Joins("LEFT JOIN orbia_kol ON orbia_kol_order.kol_id = orbia_kol.id").
		Where("orbia_kol_order.kol_id = ?", kolID)

	if status != nil && *status != "" {
		query = query.Where("orbia_kol_order.status = ?", *status)
	}

	// 获取总数
	if err := query.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	// 获取列表，按创建时间倒序
	if err := query.Order("orbia_kol_order.created_at DESC").
		Offset(offset).
		Limit(limit).
		Find(&orders).Error; err != nil {
		return nil, 0, err
	}

	return orders, total, nil
}

// GetTeamOrders 获取团队的订单列表
func (r *orderRepository) GetTeamOrders(teamID int64, status *string, offset, limit int) ([]*OrderWithKolInfo, int64, error) {
	var orders []*OrderWithKolInfo
	var total int64

	query := r.db.Table("orbia_kol_order").
		Select("orbia_kol_order.*, orbia_kol.display_name as kol_display_name, orbia_kol.avatar_url as kol_avatar_url").
		Joins("LEFT JOIN orbia_kol ON orbia_kol_order.kol_id = orbia_kol.id").
		Where("orbia_kol_order.team_id = ?", teamID)

	if status != nil && *status != "" {
		query = query.Where("orbia_kol_order.status = ?", *status)
	}

	// 获取总数
	if err := query.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	// 获取列表，按创建时间倒序
	if err := query.Order("orbia_kol_order.created_at DESC").
		Offset(offset).
		Limit(limit).
		Find(&orders).Error; err != nil {
		return nil, 0, err
	}

	return orders, total, nil
}

// UpdateOrderStatus 更新订单状态
func (r *orderRepository) UpdateOrderStatus(orderID string, status string, reason *string) error {
	updates := map[string]interface{}{
		"status": status,
	}

	now := time.Now()
	switch status {
	case "confirmed":
		updates["confirmed_at"] = now
	case "completed":
		updates["completed_at"] = now
	case "cancelled", "refunded":
		updates["cancelled_at"] = now
		if reason != nil {
			updates["reject_reason"] = *reason
		}
	}

	return r.db.Model(&KolOrder{}).
		Where("order_id = ?", orderID).
		Updates(updates).Error
}

// UpdateOrder 更新订单
func (r *orderRepository) UpdateOrder(order *KolOrder) error {
	return r.db.Save(order).Error
}

// IsOrderOwner 检查用户是否拥有订单
func (r *orderRepository) IsOrderOwner(orderID string, userID int64) (bool, error) {
	var count int64
	err := r.db.Model(&KolOrder{}).
		Where("order_id = ? AND user_id = ?", orderID, userID).
		Count(&count).Error
	if err != nil {
		return false, err
	}
	return count > 0, nil
}

// GetAllOrders 获取所有订单列表（管理员功能）
func (r *orderRepository) GetAllOrders(keyword string, status string, offset int, limit int) ([]*OrderWithKolInfo, int64, error) {
	var orders []*OrderWithKolInfo
	var total int64

	query := r.db.Table("orbia_kol_order").
		Select("orbia_kol_order.*, orbia_kol.display_name as kol_display_name, orbia_kol.avatar_url as kol_avatar_url").
		Joins("LEFT JOIN orbia_kol ON orbia_kol_order.kol_id = orbia_kol.id").
		Joins("LEFT JOIN orbia_user ON orbia_kol_order.user_id = orbia_user.id")

	// 关键字搜索（订单ID、用户名、邮箱、钱包地址）
	if keyword != "" {
		query = query.Where("orbia_kol_order.order_id LIKE ? OR orbia_user.nickname LIKE ? OR orbia_user.email LIKE ? OR orbia_user.wallet_address LIKE ?",
			"%"+keyword+"%", "%"+keyword+"%", "%"+keyword+"%", "%"+keyword+"%")
	}

	// 状态筛选
	if status != "" {
		query = query.Where("orbia_kol_order.status = ?", status)
	}

	// 获取总数
	if err := query.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	// 获取分页数据
	err := query.Order("orbia_kol_order.created_at DESC").
		Offset(offset).
		Limit(limit).
		Find(&orders).Error

	return orders, total, err
}
