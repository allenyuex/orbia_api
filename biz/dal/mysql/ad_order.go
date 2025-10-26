package mysql

import (
	"time"

	"gorm.io/gorm"
)

// AdOrder 广告订单模型
type AdOrder struct {
	ID             int64          `gorm:"primaryKey;autoIncrement;column:id" json:"id"`
	OrderID        string         `gorm:"uniqueIndex;column:order_id;size:64;not null" json:"order_id"`
	UserID         int64          `gorm:"column:user_id;not null" json:"user_id"`
	TeamID         *int64         `gorm:"column:team_id" json:"team_id"`
	Title          string         `gorm:"column:title;size:200;not null" json:"title"`
	Description    string         `gorm:"column:description;type:text;not null" json:"description"`
	Budget         float64        `gorm:"column:budget;type:decimal(12,2);not null" json:"budget"`
	AdType         string         `gorm:"column:ad_type;size:50;not null" json:"ad_type"`
	TargetAudience string         `gorm:"column:target_audience;size:500;not null" json:"target_audience"`
	StartDate      string         `gorm:"column:start_date;type:date;not null" json:"start_date"`
	EndDate        string         `gorm:"column:end_date;type:date;not null" json:"end_date"`
	Status         string         `gorm:"column:status;type:enum('pending','approved','in_progress','completed','cancelled');default:pending;not null" json:"status"`
	RejectReason   *string        `gorm:"column:reject_reason;type:text" json:"reject_reason"`
	ApprovedAt     *time.Time     `gorm:"column:approved_at" json:"approved_at"`
	CompletedAt    *time.Time     `gorm:"column:completed_at" json:"completed_at"`
	CancelledAt    *time.Time     `gorm:"column:cancelled_at" json:"cancelled_at"`
	CreatedAt      time.Time      `gorm:"column:created_at;autoCreateTime" json:"created_at"`
	UpdatedAt      time.Time      `gorm:"column:updated_at;autoUpdateTime" json:"updated_at"`
	DeletedAt      gorm.DeletedAt `gorm:"index;column:deleted_at" json:"-"`
}

// TableName 指定表名
func (AdOrder) TableName() string {
	return "orbia_ad_order"
}

// AdOrderWithUserInfo 广告订单和用户信息的联合查询结果
type AdOrderWithUserInfo struct {
	AdOrder
	UserNickname *string `json:"user_nickname"`
	TeamName     *string `json:"team_name"`
}

// AdOrderRepository 广告订单仓储接口
type AdOrderRepository interface {
	// 创建广告订单
	CreateAdOrder(order *AdOrder) error

	// 根据订单ID获取广告订单
	GetAdOrderByID(orderID string) (*AdOrder, error)

	// 根据订单ID获取广告订单（包含用户信息）
	GetAdOrderWithUserInfo(orderID string) (*AdOrderWithUserInfo, error)

	// 获取用户的广告订单列表（支持模糊搜索）
	GetUserAdOrders(userID int64, status *string, keyword *string, adType *string, offset, limit int) ([]*AdOrderWithUserInfo, int64, error)

	// 获取团队的广告订单列表
	GetTeamAdOrders(teamID int64, status *string, offset, limit int) ([]*AdOrderWithUserInfo, int64, error)

	// 获取所有广告订单列表（管理员）（支持模糊搜索）
	GetAllAdOrders(status *string, keyword *string, adType *string, offset, limit int) ([]*AdOrderWithUserInfo, int64, error)

	// 更新广告订单状态
	UpdateAdOrderStatus(orderID string, status string, reason *string) error

	// 更新广告订单
	UpdateAdOrder(order *AdOrder) error

	// 检查用户是否拥有广告订单
	IsAdOrderOwner(orderID string, userID int64) (bool, error)
}

// adOrderRepository 广告订单仓储实现
type adOrderRepository struct {
	db *gorm.DB
}

// NewAdOrderRepository 创建广告订单仓储实例
func NewAdOrderRepository(db *gorm.DB) AdOrderRepository {
	return &adOrderRepository{db: db}
}

// CreateAdOrder 创建广告订单
func (r *adOrderRepository) CreateAdOrder(order *AdOrder) error {
	return r.db.Create(order).Error
}

// GetAdOrderByID 根据订单ID获取广告订单
func (r *adOrderRepository) GetAdOrderByID(orderID string) (*AdOrder, error) {
	var order AdOrder
	err := r.db.Where("order_id = ?", orderID).First(&order).Error
	if err != nil {
		return nil, err
	}
	return &order, nil
}

// GetAdOrderWithUserInfo 根据订单ID获取广告订单（包含用户信息）
func (r *adOrderRepository) GetAdOrderWithUserInfo(orderID string) (*AdOrderWithUserInfo, error) {
	var result AdOrderWithUserInfo
	err := r.db.Table("orbia_ad_order").
		Select("orbia_ad_order.*, orbia_user.nickname as user_nickname, orbia_team.name as team_name").
		Joins("LEFT JOIN orbia_user ON orbia_ad_order.user_id = orbia_user.id").
		Joins("LEFT JOIN orbia_team ON orbia_ad_order.team_id = orbia_team.id").
		Where("orbia_ad_order.order_id = ?", orderID).
		First(&result).Error
	if err != nil {
		return nil, err
	}
	return &result, nil
}

// GetUserAdOrders 获取用户的广告订单列表（支持模糊搜索）
func (r *adOrderRepository) GetUserAdOrders(userID int64, status *string, keyword *string, adType *string, offset, limit int) ([]*AdOrderWithUserInfo, int64, error) {
	var orders []*AdOrderWithUserInfo
	var total int64

	query := r.db.Table("orbia_ad_order").
		Select("orbia_ad_order.*, orbia_user.nickname as user_nickname, orbia_team.name as team_name").
		Joins("LEFT JOIN orbia_user ON orbia_ad_order.user_id = orbia_user.id").
		Joins("LEFT JOIN orbia_team ON orbia_ad_order.team_id = orbia_team.id").
		Where("orbia_ad_order.user_id = ?", userID)

	if status != nil && *status != "" {
		query = query.Where("orbia_ad_order.status = ?", *status)
	}

	if adType != nil && *adType != "" {
		query = query.Where("orbia_ad_order.ad_type = ?", *adType)
	}

	if keyword != nil && *keyword != "" {
		likeKeyword := "%" + *keyword + "%"
		query = query.Where("orbia_ad_order.title LIKE ? OR orbia_ad_order.order_id LIKE ? OR orbia_ad_order.description LIKE ?",
			likeKeyword, likeKeyword, likeKeyword)
	}

	// 获取总数
	if err := query.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	// 获取列表，按创建时间倒序
	if err := query.Order("orbia_ad_order.created_at DESC").
		Offset(offset).
		Limit(limit).
		Find(&orders).Error; err != nil {
		return nil, 0, err
	}

	return orders, total, nil
}

// GetTeamAdOrders 获取团队的广告订单列表
func (r *adOrderRepository) GetTeamAdOrders(teamID int64, status *string, offset, limit int) ([]*AdOrderWithUserInfo, int64, error) {
	var orders []*AdOrderWithUserInfo
	var total int64

	query := r.db.Table("orbia_ad_order").
		Select("orbia_ad_order.*, orbia_user.nickname as user_nickname, orbia_team.name as team_name").
		Joins("LEFT JOIN orbia_user ON orbia_ad_order.user_id = orbia_user.id").
		Joins("LEFT JOIN orbia_team ON orbia_ad_order.team_id = orbia_team.id").
		Where("orbia_ad_order.team_id = ?", teamID)

	if status != nil && *status != "" {
		query = query.Where("orbia_ad_order.status = ?", *status)
	}

	// 获取总数
	if err := query.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	// 获取列表，按创建时间倒序
	if err := query.Order("orbia_ad_order.created_at DESC").
		Offset(offset).
		Limit(limit).
		Find(&orders).Error; err != nil {
		return nil, 0, err
	}

	return orders, total, nil
}

// GetAllAdOrders 获取所有广告订单列表（管理员）（支持模糊搜索）
func (r *adOrderRepository) GetAllAdOrders(status *string, keyword *string, adType *string, offset, limit int) ([]*AdOrderWithUserInfo, int64, error) {
	var orders []*AdOrderWithUserInfo
	var total int64

	query := r.db.Table("orbia_ad_order").
		Select("orbia_ad_order.*, orbia_user.nickname as user_nickname, orbia_team.name as team_name").
		Joins("LEFT JOIN orbia_user ON orbia_ad_order.user_id = orbia_user.id").
		Joins("LEFT JOIN orbia_team ON orbia_ad_order.team_id = orbia_team.id")

	if status != nil && *status != "" {
		query = query.Where("orbia_ad_order.status = ?", *status)
	}

	if adType != nil && *adType != "" {
		query = query.Where("orbia_ad_order.ad_type = ?", *adType)
	}

	if keyword != nil && *keyword != "" {
		likeKeyword := "%" + *keyword + "%"
		query = query.Where("orbia_ad_order.title LIKE ? OR orbia_ad_order.order_id LIKE ? OR orbia_ad_order.description LIKE ? OR orbia_user.nickname LIKE ? OR orbia_user.email LIKE ?",
			likeKeyword, likeKeyword, likeKeyword, likeKeyword, likeKeyword)
	}

	// 获取总数
	if err := query.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	// 获取列表，按创建时间倒序
	if err := query.Order("orbia_ad_order.created_at DESC").
		Offset(offset).
		Limit(limit).
		Find(&orders).Error; err != nil {
		return nil, 0, err
	}

	return orders, total, nil
}

// UpdateAdOrderStatus 更新广告订单状态
func (r *adOrderRepository) UpdateAdOrderStatus(orderID string, status string, reason *string) error {
	updates := map[string]interface{}{
		"status": status,
	}

	now := time.Now()
	switch status {
	case "approved":
		updates["approved_at"] = now
	case "completed":
		updates["completed_at"] = now
	case "cancelled":
		updates["cancelled_at"] = now
		if reason != nil {
			updates["reject_reason"] = *reason
		}
	}

	return r.db.Model(&AdOrder{}).
		Where("order_id = ?", orderID).
		Updates(updates).Error
}

// UpdateAdOrder 更新广告订单
func (r *adOrderRepository) UpdateAdOrder(order *AdOrder) error {
	return r.db.Save(order).Error
}

// IsAdOrderOwner 检查用户是否拥有广告订单
func (r *adOrderRepository) IsAdOrderOwner(orderID string, userID int64) (bool, error) {
	var count int64
	err := r.db.Model(&AdOrder{}).
		Where("order_id = ? AND user_id = ?", orderID, userID).
		Count(&count).Error
	if err != nil {
		return false, err
	}
	return count > 0, nil
}
