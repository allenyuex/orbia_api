package mysql

import (
	"time"

	"gorm.io/gorm"
)

// Kol KOL信息模型
type Kol struct {
	ID           int64          `gorm:"primaryKey;autoIncrement;column:id" json:"id"`
	UserID       int64          `gorm:"uniqueIndex;column:user_id;not null" json:"user_id"`
	AvatarURL    *string        `gorm:"column:avatar_url;size:500" json:"avatar_url"`
	DisplayName  *string        `gorm:"column:display_name;size:100" json:"display_name"`
	Description  *string        `gorm:"column:description;type:text" json:"description"`
	Country      *string        `gorm:"column:country;size:50" json:"country"`
	TiktokURL    *string        `gorm:"column:tiktok_url;size:500" json:"tiktok_url"`
	YoutubeURL   *string        `gorm:"column:youtube_url;size:500" json:"youtube_url"`
	XURL         *string        `gorm:"column:x_url;size:500" json:"x_url"`
	DiscordURL   *string        `gorm:"column:discord_url;size:500" json:"discord_url"`
	Status       string         `gorm:"column:status;type:enum('pending','approved','rejected');default:pending;not null" json:"status"`
	RejectReason *string        `gorm:"column:reject_reason;type:text" json:"reject_reason"`
	ApprovedAt   *time.Time     `gorm:"column:approved_at" json:"approved_at"`
	CreatedAt    time.Time      `gorm:"column:created_at;autoCreateTime" json:"created_at"`
	UpdatedAt    time.Time      `gorm:"column:updated_at;autoUpdateTime" json:"updated_at"`
	DeletedAt    gorm.DeletedAt `gorm:"index;column:deleted_at" json:"-"`
}

// TableName 指定表名
func (Kol) TableName() string {
	return "orbia_kol"
}

// KolLanguage KOL语言模型
type KolLanguage struct {
	ID           int64     `gorm:"primaryKey;autoIncrement;column:id" json:"id"`
	KolID        int64     `gorm:"column:kol_id;not null;uniqueIndex:uk_kol_language,priority:1" json:"kol_id"`
	LanguageCode string    `gorm:"column:language_code;size:10;not null;uniqueIndex:uk_kol_language,priority:2" json:"language_code"`
	LanguageName string    `gorm:"column:language_name;size:50;not null" json:"language_name"`
	CreatedAt    time.Time `gorm:"column:created_at;autoCreateTime" json:"created_at"`
}

// TableName 指定表名
func (KolLanguage) TableName() string {
	return "orbia_kol_language"
}

// KolTag KOL标签模型
type KolTag struct {
	ID        int64     `gorm:"primaryKey;autoIncrement;column:id" json:"id"`
	KolID     int64     `gorm:"column:kol_id;not null;uniqueIndex:uk_kol_tag,priority:1" json:"kol_id"`
	Tag       string    `gorm:"column:tag;size:50;not null;uniqueIndex:uk_kol_tag,priority:2;index" json:"tag"`
	CreatedAt time.Time `gorm:"column:created_at;autoCreateTime" json:"created_at"`
}

// TableName 指定表名
func (KolTag) TableName() string {
	return "orbia_kol_tag"
}

// KolStats KOL数据统计模型
type KolStats struct {
	ID                 int64     `gorm:"primaryKey;autoIncrement;column:id" json:"id"`
	KolID              int64     `gorm:"uniqueIndex;column:kol_id;not null" json:"kol_id"`
	TotalFollowers     int64     `gorm:"column:total_followers;default:0" json:"total_followers"`
	TiktokFollowers    int64     `gorm:"column:tiktok_followers;default:0" json:"tiktok_followers"`
	YoutubeSubscribers int64     `gorm:"column:youtube_subscribers;default:0" json:"youtube_subscribers"`
	XFollowers         int64     `gorm:"column:x_followers;default:0" json:"x_followers"`
	DiscordMembers     int64     `gorm:"column:discord_members;default:0" json:"discord_members"`
	TiktokAvgViews     int64     `gorm:"column:tiktok_avg_views;default:0" json:"tiktok_avg_views"`
	EngagementRate     float64   `gorm:"column:engagement_rate;type:decimal(10,2);default:0" json:"engagement_rate"`
	CreatedAt          time.Time `gorm:"column:created_at;autoCreateTime" json:"created_at"`
	UpdatedAt          time.Time `gorm:"column:updated_at;autoUpdateTime" json:"updated_at"`
}

// TableName 指定表名
func (KolStats) TableName() string {
	return "orbia_kol_stats"
}

// KolPlan KOL报价Plan模型
type KolPlan struct {
	ID          int64          `gorm:"primaryKey;autoIncrement;column:id" json:"id"`
	KolID       int64          `gorm:"column:kol_id;not null" json:"kol_id"`
	Title       string         `gorm:"column:title;size:200;not null" json:"title"`
	Description *string        `gorm:"column:description;type:text" json:"description"`
	Price       float64        `gorm:"column:price;type:decimal(10,2);not null" json:"price"`
	PlanType    string         `gorm:"column:plan_type;type:enum('basic','standard','premium');not null" json:"plan_type"`
	CreatedAt   time.Time      `gorm:"column:created_at;autoCreateTime" json:"created_at"`
	UpdatedAt   time.Time      `gorm:"column:updated_at;autoUpdateTime" json:"updated_at"`
	DeletedAt   gorm.DeletedAt `gorm:"index;column:deleted_at" json:"-"`
}

// TableName 指定表名
func (KolPlan) TableName() string {
	return "orbia_kol_plan"
}

// KolVideo KOL视频模型
type KolVideo struct {
	ID              int64          `gorm:"primaryKey;autoIncrement;column:id" json:"id"`
	KolID           int64          `gorm:"column:kol_id;not null" json:"kol_id"`
	Title           string         `gorm:"column:title;size:500;not null" json:"title"`
	Content         *string        `gorm:"column:content;type:text" json:"content"`
	CoverURL        *string        `gorm:"column:cover_url;size:500" json:"cover_url"`
	VideoURL        *string        `gorm:"column:video_url;size:500" json:"video_url"`
	Platform        string         `gorm:"column:platform;size:50;not null" json:"platform"`
	PlatformVideoID *string        `gorm:"column:platform_video_id;size:200" json:"platform_video_id"`
	LikesCount      int64          `gorm:"column:likes_count;default:0" json:"likes_count"`
	ViewsCount      int64          `gorm:"column:views_count;default:0" json:"views_count"`
	CommentsCount   int64          `gorm:"column:comments_count;default:0" json:"comments_count"`
	SharesCount     int64          `gorm:"column:shares_count;default:0" json:"shares_count"`
	PublishedAt     *time.Time     `gorm:"column:published_at" json:"published_at"`
	CreatedAt       time.Time      `gorm:"column:created_at;autoCreateTime" json:"created_at"`
	UpdatedAt       time.Time      `gorm:"column:updated_at;autoUpdateTime" json:"updated_at"`
	DeletedAt       gorm.DeletedAt `gorm:"index;column:deleted_at" json:"-"`
}

// TableName 指定表名
func (KolVideo) TableName() string {
	return "orbia_kol_video"
}

// KolRepository KOL仓储接口
type KolRepository interface {
	// KOL基本信息
	CreateKol(kol *Kol) error
	GetKolByID(id int64) (*Kol, error)
	GetKolByUserID(userID int64) (*Kol, error)
	UpdateKol(kol *Kol) error
	DeleteKol(id int64) error
	GetKolList(status *string, country *string, tag *string, offset, limit int) ([]*Kol, int64, error)

	// KOL语言
	CreateKolLanguage(language *KolLanguage) error
	GetKolLanguages(kolID int64) ([]*KolLanguage, error)
	DeleteKolLanguages(kolID int64) error

	// KOL标签
	CreateKolTag(tag *KolTag) error
	GetKolTags(kolID int64) ([]*KolTag, error)
	DeleteKolTags(kolID int64) error

	// KOL统计数据
	CreateKolStats(stats *KolStats) error
	GetKolStats(kolID int64) (*KolStats, error)
	UpdateKolStats(stats *KolStats) error

	// KOL报价Plan
	CreateKolPlan(plan *KolPlan) error
	GetKolPlanByID(id int64) (*KolPlan, error)
	GetKolPlans(kolID int64) ([]*KolPlan, error)
	UpdateKolPlan(plan *KolPlan) error
	DeleteKolPlan(id int64) error

	// KOL视频
	CreateKolVideo(video *KolVideo) error
	GetKolVideoByID(id int64) (*KolVideo, error)
	GetKolVideos(kolID int64, offset, limit int) ([]*KolVideo, int64, error)
	UpdateKolVideo(video *KolVideo) error
	DeleteKolVideo(id int64) error
}

// kolRepository KOL仓储实现
type kolRepository struct {
	db *gorm.DB
}

// NewKolRepository 创建KOL仓储实例
func NewKolRepository(db *gorm.DB) KolRepository {
	return &kolRepository{db: db}
}

// CreateKol 创建KOL
func (r *kolRepository) CreateKol(kol *Kol) error {
	return r.db.Create(kol).Error
}

// GetKolByID 根据ID获取KOL
func (r *kolRepository) GetKolByID(id int64) (*Kol, error) {
	var kol Kol
	err := r.db.Where("id = ?", id).First(&kol).Error
	if err != nil {
		return nil, err
	}
	return &kol, nil
}

// GetKolByUserID 根据用户ID获取KOL
func (r *kolRepository) GetKolByUserID(userID int64) (*Kol, error) {
	var kol Kol
	err := r.db.Where("user_id = ?", userID).First(&kol).Error
	if err != nil {
		return nil, err
	}
	return &kol, nil
}

// UpdateKol 更新KOL信息
func (r *kolRepository) UpdateKol(kol *Kol) error {
	return r.db.Save(kol).Error
}

// DeleteKol 删除KOL（软删除）
func (r *kolRepository) DeleteKol(id int64) error {
	return r.db.Delete(&Kol{}, id).Error
}

// GetKolList 获取KOL列表
func (r *kolRepository) GetKolList(status *string, country *string, tag *string, offset, limit int) ([]*Kol, int64, error) {
	var kols []*Kol
	var total int64

	query := r.db.Model(&Kol{})

	// 如果status为空，默认只查询已审核通过的KOL
	if status != nil && *status != "" {
		query = query.Where("status = ?", *status)
	} else {
		query = query.Where("status = ?", "approved")
	}

	if country != nil && *country != "" {
		query = query.Where("country = ?", *country)
	}

	// 如果有标签过滤，需要join KolTag表
	if tag != nil && *tag != "" {
		query = query.Joins("INNER JOIN orbia_kol_tag ON orbia_kol.id = orbia_kol_tag.kol_id").
			Where("orbia_kol_tag.tag = ?", *tag)
	}

	// 获取总数
	if err := query.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	// 获取列表
	if err := query.Offset(offset).Limit(limit).Find(&kols).Error; err != nil {
		return nil, 0, err
	}

	return kols, total, nil
}

// CreateKolLanguage 创建KOL语言
func (r *kolRepository) CreateKolLanguage(language *KolLanguage) error {
	return r.db.Create(language).Error
}

// GetKolLanguages 获取KOL语言列表
func (r *kolRepository) GetKolLanguages(kolID int64) ([]*KolLanguage, error) {
	var languages []*KolLanguage
	err := r.db.Where("kol_id = ?", kolID).Find(&languages).Error
	if err != nil {
		return nil, err
	}
	return languages, nil
}

// DeleteKolLanguages 删除KOL的所有语言
func (r *kolRepository) DeleteKolLanguages(kolID int64) error {
	return r.db.Where("kol_id = ?", kolID).Delete(&KolLanguage{}).Error
}

// CreateKolTag 创建KOL标签
func (r *kolRepository) CreateKolTag(tag *KolTag) error {
	return r.db.Create(tag).Error
}

// GetKolTags 获取KOL标签列表
func (r *kolRepository) GetKolTags(kolID int64) ([]*KolTag, error) {
	var tags []*KolTag
	err := r.db.Where("kol_id = ?", kolID).Find(&tags).Error
	if err != nil {
		return nil, err
	}
	return tags, nil
}

// DeleteKolTags 删除KOL的所有标签
func (r *kolRepository) DeleteKolTags(kolID int64) error {
	return r.db.Where("kol_id = ?", kolID).Delete(&KolTag{}).Error
}

// CreateKolStats 创建KOL统计数据
func (r *kolRepository) CreateKolStats(stats *KolStats) error {
	return r.db.Create(stats).Error
}

// GetKolStats 获取KOL统计数据
func (r *kolRepository) GetKolStats(kolID int64) (*KolStats, error) {
	var stats KolStats
	err := r.db.Where("kol_id = ?", kolID).First(&stats).Error
	if err != nil {
		return nil, err
	}
	return &stats, nil
}

// UpdateKolStats 更新KOL统计数据
func (r *kolRepository) UpdateKolStats(stats *KolStats) error {
	return r.db.Save(stats).Error
}

// CreateKolPlan 创建KOL报价Plan
func (r *kolRepository) CreateKolPlan(plan *KolPlan) error {
	return r.db.Create(plan).Error
}

// GetKolPlanByID 根据ID获取KOL报价Plan
func (r *kolRepository) GetKolPlanByID(id int64) (*KolPlan, error) {
	var plan KolPlan
	err := r.db.Where("id = ?", id).First(&plan).Error
	if err != nil {
		return nil, err
	}
	return &plan, nil
}

// GetKolPlans 获取KOL报价Plans列表
func (r *kolRepository) GetKolPlans(kolID int64) ([]*KolPlan, error) {
	var plans []*KolPlan
	err := r.db.Where("kol_id = ?", kolID).Find(&plans).Error
	if err != nil {
		return nil, err
	}
	return plans, nil
}

// UpdateKolPlan 更新KOL报价Plan
func (r *kolRepository) UpdateKolPlan(plan *KolPlan) error {
	return r.db.Save(plan).Error
}

// DeleteKolPlan 删除KOL报价Plan（软删除）
func (r *kolRepository) DeleteKolPlan(id int64) error {
	return r.db.Delete(&KolPlan{}, id).Error
}

// CreateKolVideo 创建KOL视频
func (r *kolRepository) CreateKolVideo(video *KolVideo) error {
	return r.db.Create(video).Error
}

// GetKolVideoByID 根据ID获取KOL视频
func (r *kolRepository) GetKolVideoByID(id int64) (*KolVideo, error) {
	var video KolVideo
	err := r.db.Where("id = ?", id).First(&video).Error
	if err != nil {
		return nil, err
	}
	return &video, nil
}

// GetKolVideos 获取KOL视频列表
func (r *kolRepository) GetKolVideos(kolID int64, offset, limit int) ([]*KolVideo, int64, error) {
	var videos []*KolVideo
	var total int64

	// 获取总数
	if err := r.db.Model(&KolVideo{}).Where("kol_id = ?", kolID).Count(&total).Error; err != nil {
		return nil, 0, err
	}

	// 获取列表，按发布时间倒序
	if err := r.db.Where("kol_id = ?", kolID).
		Order("published_at DESC").
		Offset(offset).
		Limit(limit).
		Find(&videos).Error; err != nil {
		return nil, 0, err
	}

	return videos, total, nil
}

// UpdateKolVideo 更新KOL视频
func (r *kolRepository) UpdateKolVideo(video *KolVideo) error {
	return r.db.Save(video).Error
}

// DeleteKolVideo 删除KOL视频（软删除）
func (r *kolRepository) DeleteKolVideo(id int64) error {
	return r.db.Delete(&KolVideo{}, id).Error
}
