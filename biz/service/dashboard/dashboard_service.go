package dashboard

import (
	"context"
	"errors"
	"time"

	"orbia_api/biz/dal/model"
	"orbia_api/biz/dal/mysql"
	common "orbia_api/biz/model/common"
	dashboardModel "orbia_api/biz/model/dashboard"
	"orbia_api/biz/utils"

	"github.com/cloudwego/hertz/pkg/common/hlog"
	"gorm.io/gorm"
)

// DashboardService Dashboard服务
type DashboardService struct {
	excellentCaseRepo mysql.ExcellentCaseRepository
	contentTrendRepo  mysql.ContentTrendRepository
	platformStatsRepo mysql.PlatformStatsRepository
}

// NewDashboardService 创建Dashboard服务实例
func NewDashboardService(
	excellentCaseRepo mysql.ExcellentCaseRepository,
	contentTrendRepo mysql.ContentTrendRepository,
	platformStatsRepo mysql.PlatformStatsRepository,
) *DashboardService {
	return &DashboardService{
		excellentCaseRepo: excellentCaseRepo,
		contentTrendRepo:  contentTrendRepo,
		platformStatsRepo: platformStatsRepo,
	}
}

// ==================== 优秀广告案例管理 ====================

// CreateExcellentCase 创建优秀案例
func (s *DashboardService) CreateExcellentCase(ctx context.Context, req *dashboardModel.CreateExcellentCaseReq) (*dashboardModel.CreateExcellentCaseResp, error) {
	// 参数校验
	if req.VideoURL == "" {
		return &dashboardModel.CreateExcellentCaseResp{
			BaseResp: utils.BuildBaseResp(400, "视频URL不能为空"),
		}, nil
	}
	if req.CoverURL == "" {
		return &dashboardModel.CreateExcellentCaseResp{
			BaseResp: utils.BuildBaseResp(400, "封面URL不能为空"),
		}, nil
	}
	if req.Title == "" {
		return &dashboardModel.CreateExcellentCaseResp{
			BaseResp: utils.BuildBaseResp(400, "案例标题不能为空"),
		}, nil
	}

	// 构建数据库模型
	sortOrder := int32(0)
	if req.SortOrder != nil {
		sortOrder = *req.SortOrder
	}

	excellentCase := &model.OrbiaExcellentCase{
		VideoURL:    req.VideoURL,
		CoverURL:    req.CoverURL,
		Title:       req.Title,
		Description: req.Description,
		SortOrder:   sortOrder,
		Status:      1, // 默认启用
	}

	// 创建优秀案例
	if err := s.excellentCaseRepo.CreateExcellentCase(excellentCase); err != nil {
		hlog.CtxErrorf(ctx, "Failed to create excellent case: %v", err)
		return &dashboardModel.CreateExcellentCaseResp{
			BaseResp: utils.BuildBaseResp(500, "创建优秀案例失败"),
		}, nil
	}

	return &dashboardModel.CreateExcellentCaseResp{
		BaseResp: utils.BuildSuccessResp(),
		ID:       &excellentCase.ID,
	}, nil
}

// UpdateExcellentCase 更新优秀案例
func (s *DashboardService) UpdateExcellentCase(ctx context.Context, req *dashboardModel.UpdateExcellentCaseReq) (*dashboardModel.UpdateExcellentCaseResp, error) {
	// 获取现有案例
	existingCase, err := s.excellentCaseRepo.GetExcellentCaseByID(req.ID)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return &dashboardModel.UpdateExcellentCaseResp{
				BaseResp: utils.BuildBaseResp(404, "优秀案例不存在"),
			}, nil
		}
		hlog.CtxErrorf(ctx, "Failed to get excellent case: %v", err)
		return &dashboardModel.UpdateExcellentCaseResp{
			BaseResp: utils.BuildBaseResp(500, "获取优秀案例失败"),
		}, nil
	}

	// 更新字段
	if req.VideoURL != nil {
		existingCase.VideoURL = *req.VideoURL
	}
	if req.CoverURL != nil {
		existingCase.CoverURL = *req.CoverURL
	}
	if req.Title != nil {
		existingCase.Title = *req.Title
	}
	if req.Description != nil {
		existingCase.Description = req.Description
	}
	if req.SortOrder != nil {
		existingCase.SortOrder = *req.SortOrder
	}
	if req.Status != nil {
		existingCase.Status = *req.Status
	}

	// 保存更新
	if err := s.excellentCaseRepo.UpdateExcellentCase(existingCase); err != nil {
		hlog.CtxErrorf(ctx, "Failed to update excellent case: %v", err)
		return &dashboardModel.UpdateExcellentCaseResp{
			BaseResp: utils.BuildBaseResp(500, "更新优秀案例失败"),
		}, nil
	}

	return &dashboardModel.UpdateExcellentCaseResp{
		BaseResp: utils.BuildSuccessResp(),
	}, nil
}

// DeleteExcellentCase 删除优秀案例
func (s *DashboardService) DeleteExcellentCase(ctx context.Context, req *dashboardModel.DeleteExcellentCaseReq) (*dashboardModel.DeleteExcellentCaseResp, error) {
	// 检查案例是否存在
	_, err := s.excellentCaseRepo.GetExcellentCaseByID(req.ID)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return &dashboardModel.DeleteExcellentCaseResp{
				BaseResp: utils.BuildBaseResp(404, "优秀案例不存在"),
			}, nil
		}
		hlog.CtxErrorf(ctx, "Failed to get excellent case: %v", err)
		return &dashboardModel.DeleteExcellentCaseResp{
			BaseResp: utils.BuildBaseResp(500, "获取优秀案例失败"),
		}, nil
	}

	// 删除案例
	if err := s.excellentCaseRepo.DeleteExcellentCase(req.ID); err != nil {
		hlog.CtxErrorf(ctx, "Failed to delete excellent case: %v", err)
		return &dashboardModel.DeleteExcellentCaseResp{
			BaseResp: utils.BuildBaseResp(500, "删除优秀案例失败"),
		}, nil
	}

	return &dashboardModel.DeleteExcellentCaseResp{
		BaseResp: utils.BuildSuccessResp(),
	}, nil
}

// GetExcellentCaseList 获取优秀案例列表
func (s *DashboardService) GetExcellentCaseList(ctx context.Context, req *dashboardModel.GetExcellentCaseListReq) (*dashboardModel.GetExcellentCaseListResp, error) {
	// 参数校验
	page := req.Page
	if page <= 0 {
		page = 1
	}
	pageSize := req.PageSize
	if pageSize <= 0 {
		pageSize = 10
	}

	offset := (page - 1) * pageSize

	// 获取列表
	cases, total, err := s.excellentCaseRepo.GetExcellentCaseList(req.Status, int(offset), int(pageSize))
	if err != nil {
		hlog.CtxErrorf(ctx, "Failed to get excellent case list: %v", err)
		return &dashboardModel.GetExcellentCaseListResp{
			BaseResp: utils.BuildBaseResp(500, "获取优秀案例列表失败"),
		}, nil
	}

	// 转换为响应模型
	caseItems := make([]*dashboardModel.ExcellentCaseItem, 0, len(cases))
	for _, c := range cases {
		caseItems = append(caseItems, convertToExcellentCaseItem(c))
	}

	return &dashboardModel.GetExcellentCaseListResp{
		BaseResp: utils.BuildSuccessResp(),
		Cases:    caseItems,
		PageInfo: &common.PageResp{
			Page:       page,
			PageSize:   pageSize,
			Total:      total,
			TotalPages: int32((total + int64(pageSize) - 1) / int64(pageSize)),
		},
	}, nil
}

// GetExcellentCaseDetail 获取优秀案例详情
func (s *DashboardService) GetExcellentCaseDetail(ctx context.Context, req *dashboardModel.GetExcellentCaseDetailReq) (*dashboardModel.GetExcellentCaseDetailResp, error) {
	excellentCase, err := s.excellentCaseRepo.GetExcellentCaseByID(req.ID)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return &dashboardModel.GetExcellentCaseDetailResp{
				BaseResp: utils.BuildBaseResp(404, "优秀案例不存在"),
			}, nil
		}
		hlog.CtxErrorf(ctx, "Failed to get excellent case: %v", err)
		return &dashboardModel.GetExcellentCaseDetailResp{
			BaseResp: utils.BuildBaseResp(500, "获取优秀案例失败"),
		}, nil
	}

	return &dashboardModel.GetExcellentCaseDetailResp{
		BaseResp:   utils.BuildSuccessResp(),
		CaseDetail: convertToExcellentCaseItem(excellentCase),
	}, nil
}

// ==================== 内容趋势管理 ====================

// CreateContentTrend 创建内容趋势
func (s *DashboardService) CreateContentTrend(ctx context.Context, req *dashboardModel.CreateContentTrendReq) (*dashboardModel.CreateContentTrendResp, error) {
	// 参数校验
	if req.HotKeyword == "" {
		return &dashboardModel.CreateContentTrendResp{
			BaseResp: utils.BuildBaseResp(400, "热点词不能为空"),
		}, nil
	}
	if req.ValueLevel != "low" && req.ValueLevel != "medium" && req.ValueLevel != "high" {
		return &dashboardModel.CreateContentTrendResp{
			BaseResp: utils.BuildBaseResp(400, "价值等级必须是low、medium或high"),
		}, nil
	}

	// 检查排名是否已存在
	exists, err := s.contentTrendRepo.CheckRankingExists(req.Ranking, 0)
	if err != nil {
		hlog.CtxErrorf(ctx, "Failed to check ranking exists: %v", err)
		return &dashboardModel.CreateContentTrendResp{
			BaseResp: utils.BuildBaseResp(500, "检查排名失败"),
		}, nil
	}
	if exists {
		return &dashboardModel.CreateContentTrendResp{
			BaseResp: utils.BuildBaseResp(400, "该排名已存在"),
		}, nil
	}

	// 构建数据库模型
	trend := &model.OrbiaContentTrend{
		Ranking:    req.Ranking,
		HotKeyword: req.HotKeyword,
		ValueLevel: req.ValueLevel,
		Heat:       req.Heat,
		GrowthRate: req.GrowthRate,
		Status:     1, // 默认启用
	}

	// 创建内容趋势
	if err := s.contentTrendRepo.CreateContentTrend(trend); err != nil {
		hlog.CtxErrorf(ctx, "Failed to create content trend: %v", err)
		return &dashboardModel.CreateContentTrendResp{
			BaseResp: utils.BuildBaseResp(500, "创建内容趋势失败"),
		}, nil
	}

	return &dashboardModel.CreateContentTrendResp{
		BaseResp: utils.BuildSuccessResp(),
		ID:       &trend.ID,
	}, nil
}

// UpdateContentTrend 更新内容趋势
func (s *DashboardService) UpdateContentTrend(ctx context.Context, req *dashboardModel.UpdateContentTrendReq) (*dashboardModel.UpdateContentTrendResp, error) {
	// 获取现有趋势
	existingTrend, err := s.contentTrendRepo.GetContentTrendByID(req.ID)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return &dashboardModel.UpdateContentTrendResp{
				BaseResp: utils.BuildBaseResp(404, "内容趋势不存在"),
			}, nil
		}
		hlog.CtxErrorf(ctx, "Failed to get content trend: %v", err)
		return &dashboardModel.UpdateContentTrendResp{
			BaseResp: utils.BuildBaseResp(500, "获取内容趋势失败"),
		}, nil
	}

	// 如果要更新排名，检查新排名是否已存在
	if req.Ranking != nil && *req.Ranking != existingTrend.Ranking {
		exists, err := s.contentTrendRepo.CheckRankingExists(*req.Ranking, req.ID)
		if err != nil {
			hlog.CtxErrorf(ctx, "Failed to check ranking exists: %v", err)
			return &dashboardModel.UpdateContentTrendResp{
				BaseResp: utils.BuildBaseResp(500, "检查排名失败"),
			}, nil
		}
		if exists {
			return &dashboardModel.UpdateContentTrendResp{
				BaseResp: utils.BuildBaseResp(400, "该排名已存在"),
			}, nil
		}
		existingTrend.Ranking = *req.Ranking
	}

	// 更新字段
	if req.HotKeyword != nil {
		existingTrend.HotKeyword = *req.HotKeyword
	}
	if req.ValueLevel != nil {
		if *req.ValueLevel != "low" && *req.ValueLevel != "medium" && *req.ValueLevel != "high" {
			return &dashboardModel.UpdateContentTrendResp{
				BaseResp: utils.BuildBaseResp(400, "价值等级必须是low、medium或high"),
			}, nil
		}
		existingTrend.ValueLevel = *req.ValueLevel
	}
	if req.Heat != nil {
		existingTrend.Heat = *req.Heat
	}
	if req.GrowthRate != nil {
		existingTrend.GrowthRate = *req.GrowthRate
	}
	if req.Status != nil {
		existingTrend.Status = *req.Status
	}

	// 保存更新
	if err := s.contentTrendRepo.UpdateContentTrend(existingTrend); err != nil {
		hlog.CtxErrorf(ctx, "Failed to update content trend: %v", err)
		return &dashboardModel.UpdateContentTrendResp{
			BaseResp: utils.BuildBaseResp(500, "更新内容趋势失败"),
		}, nil
	}

	return &dashboardModel.UpdateContentTrendResp{
		BaseResp: utils.BuildSuccessResp(),
	}, nil
}

// DeleteContentTrend 删除内容趋势
func (s *DashboardService) DeleteContentTrend(ctx context.Context, req *dashboardModel.DeleteContentTrendReq) (*dashboardModel.DeleteContentTrendResp, error) {
	// 检查趋势是否存在
	_, err := s.contentTrendRepo.GetContentTrendByID(req.ID)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return &dashboardModel.DeleteContentTrendResp{
				BaseResp: utils.BuildBaseResp(404, "内容趋势不存在"),
			}, nil
		}
		hlog.CtxErrorf(ctx, "Failed to get content trend: %v", err)
		return &dashboardModel.DeleteContentTrendResp{
			BaseResp: utils.BuildBaseResp(500, "获取内容趋势失败"),
		}, nil
	}

	// 删除趋势
	if err := s.contentTrendRepo.DeleteContentTrend(req.ID); err != nil {
		hlog.CtxErrorf(ctx, "Failed to delete content trend: %v", err)
		return &dashboardModel.DeleteContentTrendResp{
			BaseResp: utils.BuildBaseResp(500, "删除内容趋势失败"),
		}, nil
	}

	return &dashboardModel.DeleteContentTrendResp{
		BaseResp: utils.BuildSuccessResp(),
	}, nil
}

// GetContentTrendList 获取内容趋势列表
func (s *DashboardService) GetContentTrendList(ctx context.Context, req *dashboardModel.GetContentTrendListReq) (*dashboardModel.GetContentTrendListResp, error) {
	// 参数校验
	page := req.Page
	if page <= 0 {
		page = 1
	}
	pageSize := req.PageSize
	if pageSize <= 0 {
		pageSize = 10
	}

	offset := (page - 1) * pageSize

	// 获取列表
	trends, total, err := s.contentTrendRepo.GetContentTrendList(req.Status, int(offset), int(pageSize))
	if err != nil {
		hlog.CtxErrorf(ctx, "Failed to get content trend list: %v", err)
		return &dashboardModel.GetContentTrendListResp{
			BaseResp: utils.BuildBaseResp(500, "获取内容趋势列表失败"),
		}, nil
	}

	// 转换为响应模型
	trendItems := make([]*dashboardModel.ContentTrendItem, 0, len(trends))
	for _, t := range trends {
		trendItems = append(trendItems, convertToContentTrendItem(t))
	}

	return &dashboardModel.GetContentTrendListResp{
		BaseResp: utils.BuildSuccessResp(),
		Trends:   trendItems,
		PageInfo: &common.PageResp{
			Page:       page,
			PageSize:   pageSize,
			Total:      total,
			TotalPages: int32((total + int64(pageSize) - 1) / int64(pageSize)),
		},
	}, nil
}

// GetContentTrendDetail 获取内容趋势详情
func (s *DashboardService) GetContentTrendDetail(ctx context.Context, req *dashboardModel.GetContentTrendDetailReq) (*dashboardModel.GetContentTrendDetailResp, error) {
	trend, err := s.contentTrendRepo.GetContentTrendByID(req.ID)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return &dashboardModel.GetContentTrendDetailResp{
				BaseResp: utils.BuildBaseResp(404, "内容趋势不存在"),
			}, nil
		}
		hlog.CtxErrorf(ctx, "Failed to get content trend: %v", err)
		return &dashboardModel.GetContentTrendDetailResp{
			BaseResp: utils.BuildBaseResp(500, "获取内容趋势失败"),
		}, nil
	}

	return &dashboardModel.GetContentTrendDetailResp{
		BaseResp:    utils.BuildSuccessResp(),
		TrendDetail: convertToContentTrendItem(trend),
	}, nil
}

// ==================== 平台数据统计管理 ====================

// UpdatePlatformStats 更新平台数据
func (s *DashboardService) UpdatePlatformStats(ctx context.Context, req *dashboardModel.UpdatePlatformStatsReq) (*dashboardModel.UpdatePlatformStatsResp, error) {
	// 获取或创建平台数据统计
	stats, err := s.platformStatsRepo.GetOrCreatePlatformStats()
	if err != nil {
		hlog.CtxErrorf(ctx, "Failed to get or create platform stats: %v", err)
		return &dashboardModel.UpdatePlatformStatsResp{
			BaseResp: utils.BuildBaseResp(500, "获取平台数据失败"),
		}, nil
	}

	// 更新字段
	if req.ActiveKols != nil {
		stats.ActiveKols = *req.ActiveKols
	}
	if req.TotalCoverage != nil {
		stats.TotalCoverage = *req.TotalCoverage
	}
	if req.TotalAdImpressions != nil {
		stats.TotalAdImpressions = *req.TotalAdImpressions
	}
	if req.TotalTransactionAmount != nil {
		stats.TotalTransactionAmount = *req.TotalTransactionAmount
	}
	if req.AverageRoi != nil {
		stats.AverageRoi = *req.AverageRoi
	}
	if req.AverageCpm != nil {
		stats.AverageCpm = *req.AverageCpm
	}
	if req.Web3BrandCount != nil {
		stats.Web3BrandCount = *req.Web3BrandCount
	}

	// 保存更新
	if err := s.platformStatsRepo.UpdatePlatformStats(stats); err != nil {
		hlog.CtxErrorf(ctx, "Failed to update platform stats: %v", err)
		return &dashboardModel.UpdatePlatformStatsResp{
			BaseResp: utils.BuildBaseResp(500, "更新平台数据失败"),
		}, nil
	}

	return &dashboardModel.UpdatePlatformStatsResp{
		BaseResp: utils.BuildSuccessResp(),
	}, nil
}

// GetPlatformStats 获取平台数据
func (s *DashboardService) GetPlatformStats(ctx context.Context, req *dashboardModel.GetPlatformStatsReq) (*dashboardModel.GetPlatformStatsResp, error) {
	stats, err := s.platformStatsRepo.GetOrCreatePlatformStats()
	if err != nil {
		hlog.CtxErrorf(ctx, "Failed to get or create platform stats: %v", err)
		return &dashboardModel.GetPlatformStatsResp{
			BaseResp: utils.BuildBaseResp(500, "获取平台数据失败"),
		}, nil
	}

	return &dashboardModel.GetPlatformStatsResp{
		BaseResp: utils.BuildSuccessResp(),
		Stats:    convertToPlatformStatsData(stats),
	}, nil
}

// ==================== Dashboard 数据（普通用户接口） ====================

// GetDashboardData 获取Dashboard数据
func (s *DashboardService) GetDashboardData(ctx context.Context, req *dashboardModel.GetDashboardDataReq) (*dashboardModel.GetDashboardDataResp, error) {
	// 获取所有启用的优秀广告案例
	cases, err := s.excellentCaseRepo.GetEnabledExcellentCases()
	if err != nil {
		hlog.CtxErrorf(ctx, "Failed to get enabled excellent cases: %v", err)
		return &dashboardModel.GetDashboardDataResp{
			BaseResp: utils.BuildBaseResp(500, "获取优秀案例失败"),
		}, nil
	}

	// 获取所有启用的内容趋势
	trends, err := s.contentTrendRepo.GetEnabledContentTrends()
	if err != nil {
		hlog.CtxErrorf(ctx, "Failed to get enabled content trends: %v", err)
		return &dashboardModel.GetDashboardDataResp{
			BaseResp: utils.BuildBaseResp(500, "获取内容趋势失败"),
		}, nil
	}

	// 获取平台数据统计
	stats, err := s.platformStatsRepo.GetOrCreatePlatformStats()
	if err != nil {
		hlog.CtxErrorf(ctx, "Failed to get or create platform stats: %v", err)
		return &dashboardModel.GetDashboardDataResp{
			BaseResp: utils.BuildBaseResp(500, "获取平台数据失败"),
		}, nil
	}

	// 转换为响应模型
	caseItems := make([]*dashboardModel.ExcellentCaseItem, 0, len(cases))
	for _, c := range cases {
		caseItems = append(caseItems, convertToExcellentCaseItem(c))
	}

	trendItems := make([]*dashboardModel.ContentTrendItem, 0, len(trends))
	for _, t := range trends {
		trendItems = append(trendItems, convertToContentTrendItem(t))
	}

	return &dashboardModel.GetDashboardDataResp{
		BaseResp:       utils.BuildSuccessResp(),
		ExcellentCases: caseItems,
		ContentTrends:  trendItems,
		PlatformStats:  convertToPlatformStatsData(stats),
	}, nil
}

// ==================== 辅助转换函数 ====================

// convertToExcellentCaseItem 转换为优秀案例项
func convertToExcellentCaseItem(c *model.OrbiaExcellentCase) *dashboardModel.ExcellentCaseItem {
	return &dashboardModel.ExcellentCaseItem{
		ID:          c.ID,
		VideoURL:    c.VideoURL,
		CoverURL:    c.CoverURL,
		Title:       c.Title,
		Description: c.Description,
		SortOrder:   c.SortOrder,
		Status:      c.Status,
		CreatedAt:   c.CreatedAt.Format(time.RFC3339),
		UpdatedAt:   c.UpdatedAt.Format(time.RFC3339),
	}
}

// convertToContentTrendItem 转换为内容趋势项
func convertToContentTrendItem(t *model.OrbiaContentTrend) *dashboardModel.ContentTrendItem {
	return &dashboardModel.ContentTrendItem{
		ID:         t.ID,
		Ranking:    t.Ranking,
		HotKeyword: t.HotKeyword,
		ValueLevel: t.ValueLevel,
		Heat:       t.Heat,
		GrowthRate: t.GrowthRate,
		Status:     t.Status,
		CreatedAt:  t.CreatedAt.Format(time.RFC3339),
		UpdatedAt:  t.UpdatedAt.Format(time.RFC3339),
	}
}

// convertToPlatformStatsData 转换为平台数据统计
func convertToPlatformStatsData(s *model.OrbiaPlatformStat) *dashboardModel.PlatformStatsData {
	createdAt := ""
	if s.CreatedAt != nil {
		createdAt = s.CreatedAt.Format(time.RFC3339)
	}
	updatedAt := ""
	if s.UpdatedAt != nil {
		updatedAt = s.UpdatedAt.Format(time.RFC3339)
	}

	return &dashboardModel.PlatformStatsData{
		ID:                     s.ID,
		ActiveKols:             s.ActiveKols,
		TotalCoverage:          s.TotalCoverage,
		TotalAdImpressions:     s.TotalAdImpressions,
		TotalTransactionAmount: s.TotalTransactionAmount,
		AverageRoi:             s.AverageRoi,
		AverageCpm:             s.AverageCpm,
		Web3BrandCount:         s.Web3BrandCount,
		CreatedAt:              createdAt,
		UpdatedAt:              updatedAt,
	}
}
