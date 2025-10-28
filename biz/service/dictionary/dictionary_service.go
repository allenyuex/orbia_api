package dictionary

import (
	"context"
	"errors"
	"fmt"

	"orbia_api/biz/dal/model"
	"orbia_api/biz/dal/mysql"
	dictModel "orbia_api/biz/model/dictionary"
	"orbia_api/biz/utils"

	"github.com/cloudwego/hertz/pkg/common/hlog"
	"gorm.io/gorm"
)

// DictionaryService 字典服务
type DictionaryService struct {
	dictRepo     mysql.DictionaryRepository
	dictItemRepo mysql.DictionaryItemRepository
}

// NewDictionaryService 创建字典服务实例
func NewDictionaryService(
	dictRepo mysql.DictionaryRepository,
	dictItemRepo mysql.DictionaryItemRepository,
) *DictionaryService {
	return &DictionaryService{
		dictRepo:     dictRepo,
		dictItemRepo: dictItemRepo,
	}
}

// ==================== 字典管理 ====================

// CreateDictionary 创建字典
func (s *DictionaryService) CreateDictionary(ctx context.Context, req *dictModel.CreateDictionaryReq) (*dictModel.CreateDictionaryResp, error) {
	// 验证字典编码格式
	if err := s.dictRepo.ValidateDictionaryCode(req.Code); err != nil {
		hlog.Errorf("Invalid dictionary code: %v", err)
		return &dictModel.CreateDictionaryResp{
			BaseResp: utils.BuildBaseResp(400, err.Error()),
		}, nil
	}

	// 检查字典编码是否已存在
	exists, err := s.dictRepo.CheckDictionaryCodeExists(req.Code, 0)
	if err != nil {
		hlog.Errorf("Failed to check dictionary code exists: %v", err)
		return &dictModel.CreateDictionaryResp{
			BaseResp: utils.BuildBaseResp(500, "检查字典编码失败"),
		}, nil
	}
	if exists {
		return &dictModel.CreateDictionaryResp{
			BaseResp: utils.BuildBaseResp(400, "字典编码已存在"),
		}, nil
	}

	// 创建字典
	dictionary := &model.OrbiaDictionary{
		Code:        req.Code,
		Name:        req.Name,
		Description: req.Description,
		Status:      1, // 默认启用
	}

	if err := s.dictRepo.CreateDictionary(dictionary); err != nil {
		hlog.Errorf("Failed to create dictionary: %v", err)
		return &dictModel.CreateDictionaryResp{
			BaseResp: utils.BuildBaseResp(500, "创建字典失败"),
		}, nil
	}

	// 构建响应
	dictInfo := s.buildDictionaryInfo(dictionary)
	return &dictModel.CreateDictionaryResp{
		BaseResp:   utils.BuildSuccessResp(),
		Dictionary: dictInfo,
	}, nil
}

// UpdateDictionary 更新字典（不能修改code）
func (s *DictionaryService) UpdateDictionary(ctx context.Context, req *dictModel.UpdateDictionaryReq) (*dictModel.UpdateDictionaryResp, error) {
	// 获取字典
	dictionary, err := s.dictRepo.GetDictionaryByID(req.ID)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return &dictModel.UpdateDictionaryResp{
				BaseResp: utils.BuildBaseResp(404, "字典不存在"),
			}, nil
		}
		hlog.Errorf("Failed to get dictionary: %v", err)
		return &dictModel.UpdateDictionaryResp{
			BaseResp: utils.BuildBaseResp(500, "获取字典失败"),
		}, nil
	}

	// 更新字段
	dictionary.Name = req.Name
	if req.Description != nil {
		dictionary.Description = req.Description
	}
	if req.Status != nil {
		dictionary.Status = *req.Status
	}

	// 更新字典
	if err := s.dictRepo.UpdateDictionary(dictionary); err != nil {
		hlog.Errorf("Failed to update dictionary: %v", err)
		return &dictModel.UpdateDictionaryResp{
			BaseResp: utils.BuildBaseResp(500, "更新字典失败"),
		}, nil
	}

	// 构建响应
	dictInfo := s.buildDictionaryInfo(dictionary)
	return &dictModel.UpdateDictionaryResp{
		BaseResp:   utils.BuildSuccessResp(),
		Dictionary: dictInfo,
	}, nil
}

// DeleteDictionary 删除字典（软删除）
func (s *DictionaryService) DeleteDictionary(ctx context.Context, req *dictModel.DeleteDictionaryReq) (*dictModel.DeleteDictionaryResp, error) {
	// 检查字典是否存在
	_, err := s.dictRepo.GetDictionaryByID(req.ID)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return &dictModel.DeleteDictionaryResp{
				BaseResp: utils.BuildBaseResp(404, "字典不存在"),
			}, nil
		}
		hlog.Errorf("Failed to get dictionary: %v", err)
		return &dictModel.DeleteDictionaryResp{
			BaseResp: utils.BuildBaseResp(500, "获取字典失败"),
		}, nil
	}

	// 删除字典（软删除）
	if err := s.dictRepo.DeleteDictionary(req.ID); err != nil {
		hlog.Errorf("Failed to delete dictionary: %v", err)
		return &dictModel.DeleteDictionaryResp{
			BaseResp: utils.BuildBaseResp(500, "删除字典失败"),
		}, nil
	}

	// 同时删除该字典下的所有字典项
	if err := s.dictItemRepo.DeleteDictionaryItemsByDictionaryID(req.ID); err != nil {
		hlog.Errorf("Failed to delete dictionary items: %v", err)
		// 不影响主流程，只记录日志
	}

	return &dictModel.DeleteDictionaryResp{
		BaseResp: utils.BuildSuccessResp(),
	}, nil
}

// GetDictionaries 获取字典列表
func (s *DictionaryService) GetDictionaries(ctx context.Context, req *dictModel.GetDictionariesReq) (*dictModel.GetDictionariesResp, error) {
	// 标准化分页参数
	params := utils.NormalizePaginationValue(req.Page, req.PageSize)
	offset := int((params.Page - 1) * params.PageSize)

	keyword := ""
	if req.Keyword != nil {
		keyword = *req.Keyword
	}

	// 查询字典列表
	dictionaries, total, err := s.dictRepo.GetDictionaries(keyword, req.Status, offset, int(params.PageSize))
	if err != nil {
		hlog.Errorf("Failed to get dictionaries: %v", err)
		return nil, err
	}

	// 构建响应
	dictList := make([]*dictModel.DictionaryInfo, 0, len(dictionaries))
	for _, dict := range dictionaries {
		dictList = append(dictList, s.buildDictionaryInfo(dict))
	}

	// 分页信息
	pageInfo := utils.BuildPageResp(params, total)

	return &dictModel.GetDictionariesResp{
		BaseResp:     utils.BuildSuccessResp(),
		Dictionaries: dictList,
		PageInfo:     pageInfo,
	}, nil
}

// GetDictionaryDetail 获取字典详情
func (s *DictionaryService) GetDictionaryDetail(ctx context.Context, req *dictModel.GetDictionaryDetailReq) (*dictModel.GetDictionaryDetailResp, error) {
	dictionary, err := s.dictRepo.GetDictionaryByID(req.ID)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return &dictModel.GetDictionaryDetailResp{
				BaseResp: utils.BuildBaseResp(404, "字典不存在"),
			}, nil
		}
		hlog.Errorf("Failed to get dictionary: %v", err)
		return &dictModel.GetDictionaryDetailResp{
			BaseResp: utils.BuildBaseResp(500, "获取字典失败"),
		}, nil
	}

	dictInfo := s.buildDictionaryInfo(dictionary)
	return &dictModel.GetDictionaryDetailResp{
		BaseResp:   utils.BuildSuccessResp(),
		Dictionary: dictInfo,
	}, nil
}

// GetDictionaryByCode 根据编码获取字典
func (s *DictionaryService) GetDictionaryByCode(ctx context.Context, req *dictModel.GetDictionaryByCodeReq) (*dictModel.GetDictionaryByCodeResp, error) {
	dict, err := s.dictRepo.GetDictionaryByCode(req.Code)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return &dictModel.GetDictionaryByCodeResp{
				BaseResp: utils.BuildBaseResp(404, "字典不存在"),
			}, nil
		}
		hlog.Errorf("Failed to get dictionary by code: %v", err)
		return &dictModel.GetDictionaryByCodeResp{
			BaseResp: utils.BuildBaseResp(500, "获取字典失败"),
		}, nil
	}

	dictInfo := s.buildDictionaryInfo(dict)
	return &dictModel.GetDictionaryByCodeResp{
		BaseResp:   utils.BuildSuccessResp(),
		Dictionary: dictInfo,
	}, nil
}

// ==================== 字典项管理 ====================

// CreateDictionaryItem 创建字典项
func (s *DictionaryService) CreateDictionaryItem(ctx context.Context, req *dictModel.CreateDictionaryItemReq) (*dictModel.CreateDictionaryItemResp, error) {
	// 检查字典是否存在
	_, err := s.dictRepo.GetDictionaryByID(req.DictionaryID)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return &dictModel.CreateDictionaryItemResp{
				BaseResp: utils.BuildBaseResp(404, "字典不存在"),
			}, nil
		}
		hlog.Errorf("Failed to get dictionary: %v", err)
		return &dictModel.CreateDictionaryItemResp{
			BaseResp: utils.BuildBaseResp(500, "获取字典失败"),
		}, nil
	}

	// 检查父节点是否存在（如果不是根节点）
	var parentItem *model.OrbiaDictionaryItem
	if req.ParentID > 0 {
		parentItem, err = s.dictItemRepo.GetDictionaryItemByID(req.ParentID)
		if err != nil {
			if errors.Is(err, gorm.ErrRecordNotFound) {
				return &dictModel.CreateDictionaryItemResp{
					BaseResp: utils.BuildBaseResp(404, "父节点不存在"),
				}, nil
			}
			hlog.Errorf("Failed to get parent dictionary item: %v", err)
			return &dictModel.CreateDictionaryItemResp{
				BaseResp: utils.BuildBaseResp(500, "获取父节点失败"),
			}, nil
		}

		// 检查父节点是否属于同一个字典
		if parentItem.DictionaryID != req.DictionaryID {
			return &dictModel.CreateDictionaryItemResp{
				BaseResp: utils.BuildBaseResp(400, "父节点不属于该字典"),
			}, nil
		}
	}

	// 检查字典项编码是否已存在（同一字典下同一父级下不能重复）
	exists, err := s.dictItemRepo.CheckDictionaryItemCodeExists(req.DictionaryID, req.ParentID, req.Code, 0)
	if err != nil {
		hlog.Errorf("Failed to check dictionary item code exists: %v", err)
		return &dictModel.CreateDictionaryItemResp{
			BaseResp: utils.BuildBaseResp(500, "检查字典项编码失败"),
		}, nil
	}
	if exists {
		return &dictModel.CreateDictionaryItemResp{
			BaseResp: utils.BuildBaseResp(400, "字典项编码在该父级下已存在"),
		}, nil
	}

	// 计算层级和路径
	level := int32(1)
	path := ""
	if parentItem != nil {
		level = parentItem.Level + 1
		path = parentItem.Path
	}

	// 创建字典项（先不设置完整path）
	sortOrder := int32(0)
	if req.SortOrder != nil {
		sortOrder = *req.SortOrder
	}

	item := &model.OrbiaDictionaryItem{
		DictionaryID: req.DictionaryID,
		ParentID:     req.ParentID,
		Code:         req.Code,
		Name:         req.Name,
		Description:  req.Description,
		IconURL:      req.IconURL,
		SortOrder:    sortOrder,
		Level:        level,
		Path:         path, // 暂时设置为父路径
		Status:       1,    // 默认启用
	}

	if err := s.dictItemRepo.CreateDictionaryItem(item); err != nil {
		hlog.Errorf("Failed to create dictionary item: %v", err)
		return &dictModel.CreateDictionaryItemResp{
			BaseResp: utils.BuildBaseResp(500, "创建字典项失败"),
		}, nil
	}

	// 更新path为完整路径
	if path != "" {
		item.Path = fmt.Sprintf("%s/%d", path, item.ID)
	} else {
		item.Path = fmt.Sprintf("%d", item.ID)
	}
	if err := s.dictItemRepo.UpdateDictionaryItem(item); err != nil {
		hlog.Errorf("Failed to update dictionary item path: %v", err)
		// 不影响主流程
	}

	// 构建响应
	itemInfo := s.buildDictionaryItemInfo(item)
	return &dictModel.CreateDictionaryItemResp{
		BaseResp: utils.BuildSuccessResp(),
		Item:     itemInfo,
	}, nil
}

// UpdateDictionaryItem 更新字典项
func (s *DictionaryService) UpdateDictionaryItem(ctx context.Context, req *dictModel.UpdateDictionaryItemReq) (*dictModel.UpdateDictionaryItemResp, error) {
	// 获取字典项
	item, err := s.dictItemRepo.GetDictionaryItemByID(req.ID)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return &dictModel.UpdateDictionaryItemResp{
				BaseResp: utils.BuildBaseResp(404, "字典项不存在"),
			}, nil
		}
		hlog.Errorf("Failed to get dictionary item: %v", err)
		return &dictModel.UpdateDictionaryItemResp{
			BaseResp: utils.BuildBaseResp(500, "获取字典项失败"),
		}, nil
	}

	// 更新字段
	if req.Name != nil {
		item.Name = *req.Name
	}
	if req.Description != nil {
		item.Description = req.Description
	}
	if req.IconURL != nil {
		item.IconURL = req.IconURL
	}
	if req.SortOrder != nil {
		item.SortOrder = *req.SortOrder
	}
	if req.Status != nil {
		item.Status = *req.Status
	}

	// 更新字典项
	if err := s.dictItemRepo.UpdateDictionaryItem(item); err != nil {
		hlog.Errorf("Failed to update dictionary item: %v", err)
		return &dictModel.UpdateDictionaryItemResp{
			BaseResp: utils.BuildBaseResp(500, "更新字典项失败"),
		}, nil
	}

	// 构建响应
	itemInfo := s.buildDictionaryItemInfo(item)
	return &dictModel.UpdateDictionaryItemResp{
		BaseResp: utils.BuildSuccessResp(),
		Item:     itemInfo,
	}, nil
}

// DeleteDictionaryItem 删除字典项（软删除，递归删除所有子节点）
func (s *DictionaryService) DeleteDictionaryItem(ctx context.Context, req *dictModel.DeleteDictionaryItemReq) (*dictModel.DeleteDictionaryItemResp, error) {
	// 检查字典项是否存在
	item, err := s.dictItemRepo.GetDictionaryItemByID(req.ID)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return &dictModel.DeleteDictionaryItemResp{
				BaseResp: utils.BuildBaseResp(404, "字典项不存在"),
			}, nil
		}
		hlog.Errorf("Failed to get dictionary item: %v", err)
		return &dictModel.DeleteDictionaryItemResp{
			BaseResp: utils.BuildBaseResp(500, "获取字典项失败"),
		}, nil
	}

	// 递归删除该节点及其所有子节点
	if err := s.deleteDictionaryItemRecursive(item.ID); err != nil {
		hlog.Errorf("Failed to delete dictionary item: %v", err)
		return &dictModel.DeleteDictionaryItemResp{
			BaseResp: utils.BuildBaseResp(500, "删除字典项失败"),
		}, nil
	}

	return &dictModel.DeleteDictionaryItemResp{
		BaseResp: utils.BuildSuccessResp(),
	}, nil
}

// GetDictionaryItems 获取字典项列表
func (s *DictionaryService) GetDictionaryItems(ctx context.Context, req *dictModel.GetDictionaryItemsReq) (*dictModel.GetDictionaryItemsResp, error) {
	// 检查字典是否存在
	_, err := s.dictRepo.GetDictionaryByID(req.DictionaryID)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return &dictModel.GetDictionaryItemsResp{
				BaseResp: utils.BuildBaseResp(404, "字典不存在"),
			}, nil
		}
		hlog.Errorf("Failed to get dictionary: %v", err)
		return &dictModel.GetDictionaryItemsResp{
			BaseResp: utils.BuildBaseResp(500, "获取字典失败"),
		}, nil
	}

	// 标准化分页参数
	params := utils.NormalizePaginationValue(req.Page, req.PageSize)
	offset := int((params.Page - 1) * params.PageSize)

	// 查询字典项列表
	items, total, err := s.dictItemRepo.GetDictionaryItems(req.DictionaryID, req.ParentID, req.Status, offset, int(params.PageSize))
	if err != nil {
		hlog.Errorf("Failed to get dictionary items: %v", err)
		return nil, err
	}

	// 构建响应
	itemList := make([]*dictModel.DictionaryItemInfo, 0, len(items))
	for _, item := range items {
		itemList = append(itemList, s.buildDictionaryItemInfo(item))
	}

	// 分页信息
	pageInfo := utils.BuildPageResp(params, total)

	return &dictModel.GetDictionaryItemsResp{
		BaseResp: utils.BuildSuccessResp(),
		Items:    itemList,
		PageInfo: pageInfo,
	}, nil
}

// GetAllDictionariesWithItems 批量获取所有字典和字典项（用于前端冷启动）
func (s *DictionaryService) GetAllDictionariesWithItems(ctx context.Context, req *dictModel.GetAllDictionariesWithItemsReq) (*dictModel.GetAllDictionariesWithItemsResp, error) {
	// 标准化分页参数，最大页面大小为20
	params := utils.NormalizePaginationValue(req.Page, req.PageSize)
	if params.PageSize > 20 {
		params.PageSize = 20
	}
	offset := int((params.Page - 1) * params.PageSize)

	// 只获取启用的字典
	status := int32(1)
	dictionaries, total, err := s.dictRepo.GetDictionaries("", &status, offset, int(params.PageSize))
	if err != nil {
		hlog.Errorf("Failed to get dictionaries: %v", err)
		return &dictModel.GetAllDictionariesWithItemsResp{
			BaseResp: utils.BuildBaseResp(500, "获取字典列表失败"),
		}, nil
	}

	// 构建响应
	dictWithTreeList := make([]*dictModel.DictionaryWithTree, 0, len(dictionaries))

	for _, dict := range dictionaries {
		// 获取该字典的所有启用的字典项
		items, err := s.dictItemRepo.GetDictionaryItemsByDictionaryCode(dict.Code, true)
		if err != nil {
			hlog.Errorf("Failed to get dictionary items for dictionary %s: %v", dict.Code, err)
			continue
		}

		// 构建字典信息
		dictInfo := s.buildDictionaryInfo(dict)

		// 构建树形结构
		tree := s.buildTree(items, 0)

		dictWithTreeList = append(dictWithTreeList, &dictModel.DictionaryWithTree{
			Dictionary: dictInfo,
			Tree:       tree,
		})
	}

	// 分页信息
	pageInfo := utils.BuildPageResp(params, total)

	return &dictModel.GetAllDictionariesWithItemsResp{
		BaseResp:     utils.BuildSuccessResp(),
		Dictionaries: dictWithTreeList,
		PageInfo:     pageInfo,
	}, nil
}

// ==================== 辅助方法 ====================

// buildDictionaryInfo 构建字典信息
func (s *DictionaryService) buildDictionaryInfo(dictionary *model.OrbiaDictionary) *dictModel.DictionaryInfo {
	info := &dictModel.DictionaryInfo{
		ID:     dictionary.ID,
		Code:   dictionary.Code,
		Name:   dictionary.Name,
		Status: dictionary.Status,
	}

	if dictionary.Description != nil {
		info.Description = dictionary.Description
	}

	if dictionary.CreatedAt != nil {
		info.CreatedAt = dictionary.CreatedAt.Format("2006-01-02 15:04:05")
	}

	if dictionary.UpdatedAt != nil {
		info.UpdatedAt = dictionary.UpdatedAt.Format("2006-01-02 15:04:05")
	}

	return info
}

// buildDictionaryItemInfo 构建字典项信息
func (s *DictionaryService) buildDictionaryItemInfo(item *model.OrbiaDictionaryItem) *dictModel.DictionaryItemInfo {
	info := &dictModel.DictionaryItemInfo{
		ID:           item.ID,
		DictionaryID: item.DictionaryID,
		ParentID:     item.ParentID,
		Code:         item.Code,
		Name:         item.Name,
		SortOrder:    item.SortOrder,
		Level:        item.Level,
		Path:         item.Path,
		Status:       item.Status,
	}

	if item.Description != nil {
		info.Description = item.Description
	}

	if item.IconURL != nil {
		info.IconURL = item.IconURL
	}

	if item.CreatedAt != nil {
		info.CreatedAt = item.CreatedAt.Format("2006-01-02 15:04:05")
	}

	if item.UpdatedAt != nil {
		info.UpdatedAt = item.UpdatedAt.Format("2006-01-02 15:04:05")
	}

	return info
}

// buildTree 构建树形结构
func (s *DictionaryService) buildTree(items []*model.OrbiaDictionaryItem, parentID int64) []*dictModel.DictionaryItemTreeNode {
	var nodes []*dictModel.DictionaryItemTreeNode

	for _, item := range items {
		if item.ParentID == parentID {
			node := &dictModel.DictionaryItemTreeNode{
				ID:        item.ID,
				Code:      item.Code,
				Name:      item.Name,
				SortOrder: item.SortOrder,
				Level:     item.Level,
				Status:    item.Status,
			}

			if item.Description != nil {
				node.Description = item.Description
			}

			if item.IconURL != nil {
				node.IconURL = item.IconURL
			}

			// 递归构建子节点
			node.Children = s.buildTree(items, item.ID)

			nodes = append(nodes, node)
		}
	}

	return nodes
}

// deleteDictionaryItemRecursive 递归删除字典项及其所有子节点
func (s *DictionaryService) deleteDictionaryItemRecursive(itemID int64) error {
	// 获取所有子节点
	children, err := s.dictItemRepo.GetChildrenByParentID(itemID)
	if err != nil {
		return err
	}

	// 递归删除所有子节点
	for _, child := range children {
		if err := s.deleteDictionaryItemRecursive(child.ID); err != nil {
			return err
		}
	}

	// 删除当前节点
	return s.dictItemRepo.DeleteDictionaryItem(itemID)
}
