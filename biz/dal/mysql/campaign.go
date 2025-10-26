package mysql

import (
	"time"

	"gorm.io/gorm"
)

// Campaign 广告活动模型
type Campaign struct {
	ID                 int64          `gorm:"primaryKey;autoIncrement;column:id" json:"id"`
	CampaignID         string         `gorm:"uniqueIndex;column:campaign_id;size:64;not null" json:"campaign_id"`
	UserID             int64          `gorm:"index;column:user_id;not null" json:"user_id"`
	TeamID             int64          `gorm:"index;column:team_id;not null" json:"team_id"`
	CampaignName       string         `gorm:"column:campaign_name;size:200;not null" json:"campaign_name"`
	PromotionObjective string         `gorm:"index;column:promotion_objective;type:enum('awareness','consideration','conversion');not null" json:"promotion_objective"`
	OptimizationGoal   string         `gorm:"column:optimization_goal;size:50;not null" json:"optimization_goal"`
	Location           *string        `gorm:"column:location;type:text" json:"location"`
	Age                *int64         `gorm:"column:age" json:"age"`
	Gender             *int64         `gorm:"column:gender" json:"gender"`
	Languages          *string        `gorm:"column:languages;type:text" json:"languages"`
	SpendingPower      *int64         `gorm:"column:spending_power" json:"spending_power"`
	OperatingSystem    *int64         `gorm:"column:operating_system" json:"operating_system"`
	OSVersions         *string        `gorm:"column:os_versions;type:text" json:"os_versions"`
	DeviceModels       *string        `gorm:"column:device_models;type:text" json:"device_models"`
	ConnectionType     *int64         `gorm:"column:connection_type" json:"connection_type"`
	DevicePriceType    int8           `gorm:"column:device_price_type;default:0" json:"device_price_type"`
	DevicePriceMin     *float64       `gorm:"column:device_price_min;type:decimal(15,2)" json:"device_price_min"`
	DevicePriceMax     *float64       `gorm:"column:device_price_max;type:decimal(15,2)" json:"device_price_max"`
	PlannedStartTime   time.Time      `gorm:"index;column:planned_start_time;not null" json:"planned_start_time"`
	PlannedEndTime     time.Time      `gorm:"index;column:planned_end_time;not null" json:"planned_end_time"`
	TimeZone           *int64         `gorm:"column:time_zone" json:"time_zone"`
	DaypartingType     int8           `gorm:"column:dayparting_type;default:0" json:"dayparting_type"`
	DaypartingSchedule *string        `gorm:"column:dayparting_schedule;type:text" json:"dayparting_schedule"`
	FrequencyCapType   int8           `gorm:"column:frequency_cap_type;default:0" json:"frequency_cap_type"`
	FrequencyCapTimes  *int32         `gorm:"column:frequency_cap_times" json:"frequency_cap_times"`
	FrequencyCapDays   *int32         `gorm:"column:frequency_cap_days" json:"frequency_cap_days"`
	BudgetType         int8           `gorm:"column:budget_type;not null" json:"budget_type"`
	BudgetAmount       float64        `gorm:"column:budget_amount;type:decimal(15,2);not null" json:"budget_amount"`
	Website            *string        `gorm:"column:website;size:1000" json:"website"`
	IOSDownloadURL     *string        `gorm:"column:ios_download_url;size:1000" json:"ios_download_url"`
	AndroidDownloadURL *string        `gorm:"column:android_download_url;size:1000" json:"android_download_url"`
	Status             string         `gorm:"index;column:status;type:enum('pending','active','paused','ended');default:pending;not null" json:"status"`
	CreatedAt          time.Time      `gorm:"index;column:created_at;autoCreateTime" json:"created_at"`
	UpdatedAt          time.Time      `gorm:"column:updated_at;autoUpdateTime" json:"updated_at"`
	DeletedAt          gorm.DeletedAt `gorm:"index;column:deleted_at" json:"-"`
}

// TableName 指定表名
func (Campaign) TableName() string {
	return "orbia_campaign"
}

// CampaignAttachment Campaign附件模型
type CampaignAttachment struct {
	ID         int64          `gorm:"primaryKey;autoIncrement;column:id" json:"id"`
	CampaignID int64          `gorm:"index;column:campaign_id;not null" json:"campaign_id"`
	FileURL    string         `gorm:"column:file_url;size:1000;not null" json:"file_url"`
	FileName   string         `gorm:"column:file_name;size:500;not null" json:"file_name"`
	FileType   string         `gorm:"column:file_type;size:100;not null" json:"file_type"`
	FileSize   *int64         `gorm:"column:file_size" json:"file_size"`
	CreatedAt  time.Time      `gorm:"index;column:created_at;autoCreateTime" json:"created_at"`
	UpdatedAt  time.Time      `gorm:"column:updated_at;autoUpdateTime" json:"updated_at"`
	DeletedAt  gorm.DeletedAt `gorm:"index;column:deleted_at" json:"-"`
}

// TableName 指定表名
func (CampaignAttachment) TableName() string {
	return "orbia_campaign_attachment"
}

// CampaignRepository Campaign数据仓库接口
type CampaignRepository interface {
	// Campaign CRUD
	CreateCampaign(campaign *Campaign) error
	GetCampaignByID(id int64) (*Campaign, error)
	GetCampaignByCampaignID(campaignID string) (*Campaign, error)
	UpdateCampaign(campaign *Campaign) error
	DeleteCampaign(id int64) error

	// Campaign查询
	GetCampaignsByUserID(userID int64, keyword string, status string, promotionObjective string, offset int, limit int) ([]*Campaign, int64, error)
	GetCampaignsByTeamID(teamID int64, keyword string, status string, promotionObjective string, offset int, limit int) ([]*Campaign, int64, error)
	GetAllCampaigns(keyword string, status string, promotionObjective string, userID *int64, teamID *int64, offset int, limit int) ([]*Campaign, int64, error)

	// Attachment操作
	CreateAttachment(attachment *CampaignAttachment) error
	GetAttachmentsByCampaignID(campaignID int64) ([]*CampaignAttachment, error)
	DeleteAttachmentsByCampaignID(campaignID int64) error
}

// campaignRepository Campaign数据仓库实现
type campaignRepository struct {
	db *gorm.DB
}

// NewCampaignRepository 创建Campaign数据仓库实例
func NewCampaignRepository(db *gorm.DB) CampaignRepository {
	return &campaignRepository{db: db}
}

// CreateCampaign 创建Campaign
func (r *campaignRepository) CreateCampaign(campaign *Campaign) error {
	return r.db.Create(campaign).Error
}

// GetCampaignByID 根据ID获取Campaign
func (r *campaignRepository) GetCampaignByID(id int64) (*Campaign, error) {
	var campaign Campaign
	err := r.db.Where("id = ?", id).First(&campaign).Error
	if err != nil {
		return nil, err
	}
	return &campaign, nil
}

// GetCampaignByCampaignID 根据业务ID获取Campaign
func (r *campaignRepository) GetCampaignByCampaignID(campaignID string) (*Campaign, error) {
	var campaign Campaign
	err := r.db.Where("campaign_id = ?", campaignID).First(&campaign).Error
	if err != nil {
		return nil, err
	}
	return &campaign, nil
}

// UpdateCampaign 更新Campaign
func (r *campaignRepository) UpdateCampaign(campaign *Campaign) error {
	return r.db.Save(campaign).Error
}

// DeleteCampaign 删除Campaign（软删除）
func (r *campaignRepository) DeleteCampaign(id int64) error {
	return r.db.Delete(&Campaign{}, id).Error
}

// GetCampaignsByUserID 获取用户的Campaign列表
func (r *campaignRepository) GetCampaignsByUserID(userID int64, keyword string, status string, promotionObjective string, offset int, limit int) ([]*Campaign, int64, error) {
	var campaigns []*Campaign
	var total int64

	query := r.db.Model(&Campaign{}).Where("user_id = ?", userID)

	// 关键字搜索
	if keyword != "" {
		query = query.Where("campaign_name LIKE ?", "%"+keyword+"%")
	}

	// 状态筛选
	if status != "" {
		query = query.Where("status = ?", status)
	}

	// 推广目标筛选
	if promotionObjective != "" {
		query = query.Where("promotion_objective = ?", promotionObjective)
	}

	// 获取总数
	if err := query.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	// 获取分页数据
	err := query.Order("created_at DESC").
		Offset(offset).
		Limit(limit).
		Find(&campaigns).Error

	return campaigns, total, err
}

// GetCampaignsByTeamID 获取团队的Campaign列表
func (r *campaignRepository) GetCampaignsByTeamID(teamID int64, keyword string, status string, promotionObjective string, offset int, limit int) ([]*Campaign, int64, error) {
	var campaigns []*Campaign
	var total int64

	query := r.db.Model(&Campaign{}).Where("team_id = ?", teamID)

	// 关键字搜索
	if keyword != "" {
		query = query.Where("campaign_name LIKE ?", "%"+keyword+"%")
	}

	// 状态筛选
	if status != "" {
		query = query.Where("status = ?", status)
	}

	// 推广目标筛选
	if promotionObjective != "" {
		query = query.Where("promotion_objective = ?", promotionObjective)
	}

	// 获取总数
	if err := query.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	// 获取分页数据
	err := query.Order("created_at DESC").
		Offset(offset).
		Limit(limit).
		Find(&campaigns).Error

	return campaigns, total, err
}

// GetAllCampaigns 获取所有Campaign列表（管理员）
func (r *campaignRepository) GetAllCampaigns(keyword string, status string, promotionObjective string, userID *int64, teamID *int64, offset int, limit int) ([]*Campaign, int64, error) {
	var campaigns []*Campaign
	var total int64

	query := r.db.Model(&Campaign{})

	// 关键字搜索
	if keyword != "" {
		query = query.Where("campaign_name LIKE ?", "%"+keyword+"%")
	}

	// 状态筛选
	if status != "" {
		query = query.Where("status = ?", status)
	}

	// 推广目标筛选
	if promotionObjective != "" {
		query = query.Where("promotion_objective = ?", promotionObjective)
	}

	// 用户筛选
	if userID != nil {
		query = query.Where("user_id = ?", *userID)
	}

	// 团队筛选
	if teamID != nil {
		query = query.Where("team_id = ?", *teamID)
	}

	// 获取总数
	if err := query.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	// 获取分页数据
	err := query.Order("created_at DESC").
		Offset(offset).
		Limit(limit).
		Find(&campaigns).Error

	return campaigns, total, err
}

// CreateAttachment 创建附件
func (r *campaignRepository) CreateAttachment(attachment *CampaignAttachment) error {
	return r.db.Create(attachment).Error
}

// GetAttachmentsByCampaignID 获取Campaign的所有附件
func (r *campaignRepository) GetAttachmentsByCampaignID(campaignID int64) ([]*CampaignAttachment, error) {
	var attachments []*CampaignAttachment
	err := r.db.Where("campaign_id = ?", campaignID).Find(&attachments).Error
	return attachments, err
}

// DeleteAttachmentsByCampaignID 删除Campaign的所有附件
func (r *campaignRepository) DeleteAttachmentsByCampaignID(campaignID int64) error {
	return r.db.Where("campaign_id = ?", campaignID).Delete(&CampaignAttachment{}).Error
}
