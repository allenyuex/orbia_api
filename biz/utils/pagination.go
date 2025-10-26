package utils

import (
	"orbia_api/biz/model/common"

	"gorm.io/gorm"
)

// PaginationParams 分页参数
type PaginationParams struct {
	Page     int32
	PageSize int32
}

// PaginationResult 分页结果
type PaginationResult struct {
	Total      int64
	TotalPages int32
}

// NormalizePagination 标准化分页参数
func NormalizePagination(page, pageSize *int32) PaginationParams {
	var p, ps int32

	// 默认值
	if page == nil || *page <= 0 {
		p = 1
	} else {
		p = *page
	}

	if pageSize == nil || *pageSize <= 0 {
		ps = 10
	} else if *pageSize > 100 {
		// 限制最大每页数量
		ps = 100
	} else {
		ps = *pageSize
	}

	return PaginationParams{
		Page:     p,
		PageSize: ps,
	}
}

// NormalizePaginationValue 标准化分页参数（值类型）
func NormalizePaginationValue(page, pageSize int32) PaginationParams {
	return NormalizePagination(&page, &pageSize)
}

// ApplyPagination 应用分页到查询
func ApplyPagination(db *gorm.DB, params PaginationParams) *gorm.DB {
	offset := (params.Page - 1) * params.PageSize
	return db.Offset(int(offset)).Limit(int(params.PageSize))
}

// BuildPageResp 构建分页响应
func BuildPageResp(params PaginationParams, total int64) *common.PageResp {
	totalPages := int32(total / int64(params.PageSize))
	if total%int64(params.PageSize) != 0 {
		totalPages++
	}

	return &common.PageResp{
		Page:       params.Page,
		PageSize:   params.PageSize,
		Total:      total,
		TotalPages: totalPages,
	}
}

// GetPaginationInfo 获取分页信息（用于查询和构建响应）
func GetPaginationInfo(page, pageSize *int32) (params PaginationParams, offset int) {
	params = NormalizePagination(page, pageSize)
	offset = int((params.Page - 1) * params.PageSize)
	return
}

// GetPageParams 获取分页参数，处理默认值
func GetPageParams(page, pageSize *int32) (int, int) {
	p := 1
	ps := 10

	if page != nil && *page > 0 {
		p = int(*page)
	}

	if pageSize != nil && *pageSize > 0 {
		ps = int(*pageSize)
		if ps > 100 {
			ps = 100 // 限制最大每页数量
		}
	}

	return p, ps
}
