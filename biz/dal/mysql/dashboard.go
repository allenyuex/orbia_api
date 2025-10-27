package mysql

import (
	"orbia_api/biz/dal/model"

	"gorm.io/gorm"
)

// ==================== 优秀广告案例仓储接口 ====================

// ExcellentCaseRepository 优秀广告案例仓储接口
type ExcellentCaseRepository interface {
	// 创建优秀案例
	CreateExcellentCase(excellentCase *model.OrbiaExcellentCase) error
	// 更新优秀案例
	UpdateExcellentCase(excellentCase *model.OrbiaExcellentCase) error
	// 删除优秀案例（软删除）
	DeleteExcellentCase(id int64) error
	// 根据ID获取优秀案例
	GetExcellentCaseByID(id int64) (*model.OrbiaExcellentCase, error)
	// 获取优秀案例列表
	GetExcellentCaseList(status *int32, offset int, limit int) ([]*model.OrbiaExcellentCase, int64, error)
	// 获取所有启用的优秀案例（用于前端展示）
	GetEnabledExcellentCases() ([]*model.OrbiaExcellentCase, error)
}

// excellentCaseRepository 优秀广告案例仓储实现
type excellentCaseRepository struct {
	db *gorm.DB
}

// NewExcellentCaseRepository 创建优秀广告案例仓储实例
func NewExcellentCaseRepository(db *gorm.DB) ExcellentCaseRepository {
	return &excellentCaseRepository{db: db}
}

// CreateExcellentCase 创建优秀案例
func (r *excellentCaseRepository) CreateExcellentCase(excellentCase *model.OrbiaExcellentCase) error {
	return r.db.Create(excellentCase).Error
}

// UpdateExcellentCase 更新优秀案例
func (r *excellentCaseRepository) UpdateExcellentCase(excellentCase *model.OrbiaExcellentCase) error {
	return r.db.Save(excellentCase).Error
}

// DeleteExcellentCase 删除优秀案例（软删除）
func (r *excellentCaseRepository) DeleteExcellentCase(id int64) error {
	return r.db.Delete(&model.OrbiaExcellentCase{}, id).Error
}

// GetExcellentCaseByID 根据ID获取优秀案例
func (r *excellentCaseRepository) GetExcellentCaseByID(id int64) (*model.OrbiaExcellentCase, error) {
	var excellentCase model.OrbiaExcellentCase
	err := r.db.Where("id = ?", id).First(&excellentCase).Error
	if err != nil {
		return nil, err
	}
	return &excellentCase, nil
}

// GetExcellentCaseList 获取优秀案例列表
func (r *excellentCaseRepository) GetExcellentCaseList(status *int32, offset int, limit int) ([]*model.OrbiaExcellentCase, int64, error) {
	var cases []*model.OrbiaExcellentCase
	var total int64

	query := r.db.Model(&model.OrbiaExcellentCase{})

	// 状态筛选
	if status != nil {
		query = query.Where("status = ?", *status)
	}

	// 获取总数
	if err := query.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	// 获取分页数据，按排序序号和创建时间排序
	err := query.Order("sort_order ASC, created_at ASC").
		Offset(offset).
		Limit(limit).
		Find(&cases).Error

	return cases, total, err
}

// GetEnabledExcellentCases 获取所有启用的优秀案例（用于前端展示）
func (r *excellentCaseRepository) GetEnabledExcellentCases() ([]*model.OrbiaExcellentCase, error) {
	var cases []*model.OrbiaExcellentCase
	err := r.db.Where("status = ?", 1).
		Order("sort_order ASC, created_at ASC").
		Find(&cases).Error
	return cases, err
}

// ==================== 内容趋势仓储接口 ====================

// ContentTrendRepository 内容趋势仓储接口
type ContentTrendRepository interface {
	// 创建内容趋势
	CreateContentTrend(trend *model.OrbiaContentTrend) error
	// 更新内容趋势
	UpdateContentTrend(trend *model.OrbiaContentTrend) error
	// 删除内容趋势（软删除）
	DeleteContentTrend(id int64) error
	// 根据ID获取内容趋势
	GetContentTrendByID(id int64) (*model.OrbiaContentTrend, error)
	// 获取内容趋势列表
	GetContentTrendList(status *int32, offset int, limit int) ([]*model.OrbiaContentTrend, int64, error)
	// 获取所有启用的内容趋势（用于前端展示，按排名排序）
	GetEnabledContentTrends() ([]*model.OrbiaContentTrend, error)
	// 检查排名是否已存在
	CheckRankingExists(ranking int32, excludeID int64) (bool, error)
}

// contentTrendRepository 内容趋势仓储实现
type contentTrendRepository struct {
	db *gorm.DB
}

// NewContentTrendRepository 创建内容趋势仓储实例
func NewContentTrendRepository(db *gorm.DB) ContentTrendRepository {
	return &contentTrendRepository{db: db}
}

// CreateContentTrend 创建内容趋势
func (r *contentTrendRepository) CreateContentTrend(trend *model.OrbiaContentTrend) error {
	return r.db.Create(trend).Error
}

// UpdateContentTrend 更新内容趋势
func (r *contentTrendRepository) UpdateContentTrend(trend *model.OrbiaContentTrend) error {
	return r.db.Save(trend).Error
}

// DeleteContentTrend 删除内容趋势（软删除）
func (r *contentTrendRepository) DeleteContentTrend(id int64) error {
	return r.db.Delete(&model.OrbiaContentTrend{}, id).Error
}

// GetContentTrendByID 根据ID获取内容趋势
func (r *contentTrendRepository) GetContentTrendByID(id int64) (*model.OrbiaContentTrend, error) {
	var trend model.OrbiaContentTrend
	err := r.db.Where("id = ?", id).First(&trend).Error
	if err != nil {
		return nil, err
	}
	return &trend, nil
}

// GetContentTrendList 获取内容趋势列表
func (r *contentTrendRepository) GetContentTrendList(status *int32, offset int, limit int) ([]*model.OrbiaContentTrend, int64, error) {
	var trends []*model.OrbiaContentTrend
	var total int64

	query := r.db.Model(&model.OrbiaContentTrend{})

	// 状态筛选
	if status != nil {
		query = query.Where("status = ?", *status)
	}

	// 获取总数
	if err := query.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	// 获取分页数据，按排名排序
	err := query.Order("ranking ASC, created_at ASC").
		Offset(offset).
		Limit(limit).
		Find(&trends).Error

	return trends, total, err
}

// GetEnabledContentTrends 获取所有启用的内容趋势（用于前端展示，按排名排序）
func (r *contentTrendRepository) GetEnabledContentTrends() ([]*model.OrbiaContentTrend, error) {
	var trends []*model.OrbiaContentTrend
	err := r.db.Where("status = ?", 1).
		Order("ranking ASC, created_at ASC").
		Find(&trends).Error
	return trends, err
}

// CheckRankingExists 检查排名是否已存在
func (r *contentTrendRepository) CheckRankingExists(ranking int32, excludeID int64) (bool, error) {
	var count int64
	query := r.db.Model(&model.OrbiaContentTrend{}).Where("ranking = ?", ranking)
	if excludeID > 0 {
		query = query.Where("id != ?", excludeID)
	}
	err := query.Count(&count).Error
	return count > 0, err
}

// ==================== 平台数据统计仓储接口 ====================

// PlatformStatsRepository 平台数据统计仓储接口
type PlatformStatsRepository interface {
	// 获取平台数据统计（只有一行数据）
	GetPlatformStats() (*model.OrbiaPlatformStat, error)
	// 更新平台数据统计
	UpdatePlatformStats(stats *model.OrbiaPlatformStat) error
	// 创建平台数据统计（首次创建）
	CreatePlatformStats(stats *model.OrbiaPlatformStat) error
	// 获取或创建平台数据统计
	GetOrCreatePlatformStats() (*model.OrbiaPlatformStat, error)
}

// platformStatsRepository 平台数据统计仓储实现
type platformStatsRepository struct {
	db *gorm.DB
}

// NewPlatformStatsRepository 创建平台数据统计仓储实例
func NewPlatformStatsRepository(db *gorm.DB) PlatformStatsRepository {
	return &platformStatsRepository{db: db}
}

// GetPlatformStats 获取平台数据统计（只有一行数据）
func (r *platformStatsRepository) GetPlatformStats() (*model.OrbiaPlatformStat, error) {
	var stats model.OrbiaPlatformStat
	err := r.db.First(&stats).Error
	if err != nil {
		return nil, err
	}
	return &stats, nil
}

// UpdatePlatformStats 更新平台数据统计
func (r *platformStatsRepository) UpdatePlatformStats(stats *model.OrbiaPlatformStat) error {
	return r.db.Save(stats).Error
}

// CreatePlatformStats 创建平台数据统计（首次创建）
func (r *platformStatsRepository) CreatePlatformStats(stats *model.OrbiaPlatformStat) error {
	return r.db.Create(stats).Error
}

// GetOrCreatePlatformStats 获取或创建平台数据统计
func (r *platformStatsRepository) GetOrCreatePlatformStats() (*model.OrbiaPlatformStat, error) {
	stats, err := r.GetPlatformStats()
	if err == gorm.ErrRecordNotFound {
		// 如果没有数据，创建一条默认数据
		newStats := &model.OrbiaPlatformStat{
			ActiveKols:             0,
			TotalCoverage:          0,
			TotalAdImpressions:     0,
			TotalTransactionAmount: 0,
			AverageRoi:             0,
			AverageCpm:             0,
			Web3BrandCount:         0,
		}
		if err := r.CreatePlatformStats(newStats); err != nil {
			return nil, err
		}
		return newStats, nil
	}
	if err != nil {
		return nil, err
	}
	return stats, nil
}
