package mysql

import (
	"fmt"
	"orbia_api/biz/dal/model"
	"regexp"

	"gorm.io/gorm"
)

// DictionaryRepository 字典仓储接口
type DictionaryRepository interface {
	// 字典管理
	CreateDictionary(dictionary *model.OrbiaDictionary) error
	UpdateDictionary(dictionary *model.OrbiaDictionary) error
	DeleteDictionary(id int64) error
	GetDictionaryByID(id int64) (*model.OrbiaDictionary, error)
	GetDictionaryByCode(code string) (*model.OrbiaDictionary, error)
	GetDictionaries(keyword string, status *int32, offset int, limit int) ([]*model.OrbiaDictionary, int64, error)
	ValidateDictionaryCode(code string) error
	CheckDictionaryCodeExists(code string, excludeID int64) (bool, error)
}

// DictionaryItemRepository 字典项仓储接口
type DictionaryItemRepository interface {
	// 字典项管理
	CreateDictionaryItem(item *model.OrbiaDictionaryItem) error
	UpdateDictionaryItem(item *model.OrbiaDictionaryItem) error
	DeleteDictionaryItem(id int64) error
	DeleteDictionaryItemsByDictionaryID(dictionaryID int64) error
	GetDictionaryItemByID(id int64) (*model.OrbiaDictionaryItem, error)
	GetDictionaryItems(dictionaryID int64, parentID *int64, status *int32, offset int, limit int) ([]*model.OrbiaDictionaryItem, int64, error)
	GetDictionaryItemsByDictionaryCode(dictionaryCode string, onlyEnabled bool) ([]*model.OrbiaDictionaryItem, error)
	CheckDictionaryItemCodeExists(dictionaryID, parentID int64, code string, excludeID int64) (bool, error)
	GetChildrenByParentID(parentID int64) ([]*model.OrbiaDictionaryItem, error)
}

// dictionaryRepository 字典仓储实现
type dictionaryRepository struct {
	db *gorm.DB
}

// dictionaryItemRepository 字典项仓储实现
type dictionaryItemRepository struct {
	db *gorm.DB
}

// NewDictionaryRepository 创建字典仓储实例
func NewDictionaryRepository(db *gorm.DB) DictionaryRepository {
	return &dictionaryRepository{db: db}
}

// NewDictionaryItemRepository 创建字典项仓储实例
func NewDictionaryItemRepository(db *gorm.DB) DictionaryItemRepository {
	return &dictionaryItemRepository{db: db}
}

// ==================== 字典管理实现 ====================

// CreateDictionary 创建字典
func (r *dictionaryRepository) CreateDictionary(dictionary *model.OrbiaDictionary) error {
	return r.db.Create(dictionary).Error
}

// UpdateDictionary 更新字典
func (r *dictionaryRepository) UpdateDictionary(dictionary *model.OrbiaDictionary) error {
	return r.db.Save(dictionary).Error
}

// DeleteDictionary 删除字典（软删除）
func (r *dictionaryRepository) DeleteDictionary(id int64) error {
	return r.db.Delete(&model.OrbiaDictionary{}, id).Error
}

// GetDictionaryByID 根据ID获取字典
func (r *dictionaryRepository) GetDictionaryByID(id int64) (*model.OrbiaDictionary, error) {
	var dictionary model.OrbiaDictionary
	err := r.db.Where("id = ?", id).First(&dictionary).Error
	if err != nil {
		return nil, err
	}
	return &dictionary, nil
}

// GetDictionaryByCode 根据编码获取字典
func (r *dictionaryRepository) GetDictionaryByCode(code string) (*model.OrbiaDictionary, error) {
	var dictionary model.OrbiaDictionary
	err := r.db.Where("code = ?", code).First(&dictionary).Error
	if err != nil {
		return nil, err
	}
	return &dictionary, nil
}

// GetDictionaries 获取字典列表
func (r *dictionaryRepository) GetDictionaries(keyword string, status *int32, offset int, limit int) ([]*model.OrbiaDictionary, int64, error) {
	var dictionaries []*model.OrbiaDictionary
	var total int64

	query := r.db.Model(&model.OrbiaDictionary{})

	// 关键字搜索（字典编码、名称）
	if keyword != "" {
		query = query.Where("code LIKE ? OR name LIKE ?", "%"+keyword+"%", "%"+keyword+"%")
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
		Find(&dictionaries).Error

	return dictionaries, total, err
}

// ValidateDictionaryCode 验证字典编码格式（只能包含大小写字母）
func (r *dictionaryRepository) ValidateDictionaryCode(code string) error {
	matched, err := regexp.MatchString("^[a-zA-Z]+$", code)
	if err != nil {
		return err
	}
	if !matched {
		return fmt.Errorf("字典编码只能包含大小写字母")
	}
	return nil
}

// CheckDictionaryCodeExists 检查字典编码是否存在
func (r *dictionaryRepository) CheckDictionaryCodeExists(code string, excludeID int64) (bool, error) {
	var count int64
	query := r.db.Model(&model.OrbiaDictionary{}).Where("code = ?", code)
	if excludeID > 0 {
		query = query.Where("id != ?", excludeID)
	}
	err := query.Count(&count).Error
	return count > 0, err
}

// ==================== 字典项管理实现 ====================

// CreateDictionaryItem 创建字典项
func (r *dictionaryItemRepository) CreateDictionaryItem(item *model.OrbiaDictionaryItem) error {
	return r.db.Create(item).Error
}

// UpdateDictionaryItem 更新字典项
func (r *dictionaryItemRepository) UpdateDictionaryItem(item *model.OrbiaDictionaryItem) error {
	return r.db.Save(item).Error
}

// DeleteDictionaryItem 删除字典项（软删除）
func (r *dictionaryItemRepository) DeleteDictionaryItem(id int64) error {
	return r.db.Delete(&model.OrbiaDictionaryItem{}, id).Error
}

// DeleteDictionaryItemsByDictionaryID 删除字典下的所有字典项
func (r *dictionaryItemRepository) DeleteDictionaryItemsByDictionaryID(dictionaryID int64) error {
	return r.db.Where("dictionary_id = ?", dictionaryID).Delete(&model.OrbiaDictionaryItem{}).Error
}

// GetDictionaryItemByID 根据ID获取字典项
func (r *dictionaryItemRepository) GetDictionaryItemByID(id int64) (*model.OrbiaDictionaryItem, error) {
	var item model.OrbiaDictionaryItem
	err := r.db.Where("id = ?", id).First(&item).Error
	if err != nil {
		return nil, err
	}
	return &item, nil
}

// GetDictionaryItems 获取字典项列表
func (r *dictionaryItemRepository) GetDictionaryItems(dictionaryID int64, parentID *int64, status *int32, offset int, limit int) ([]*model.OrbiaDictionaryItem, int64, error) {
	var items []*model.OrbiaDictionaryItem
	var total int64

	query := r.db.Model(&model.OrbiaDictionaryItem{}).Where("dictionary_id = ?", dictionaryID)

	// 父级筛选
	if parentID != nil {
		query = query.Where("parent_id = ?", *parentID)
	}

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
		Find(&items).Error

	return items, total, err
}

// GetDictionaryItemsByDictionaryCode 根据字典编码获取所有字典项
func (r *dictionaryItemRepository) GetDictionaryItemsByDictionaryCode(dictionaryCode string, onlyEnabled bool) ([]*model.OrbiaDictionaryItem, error) {
	var items []*model.OrbiaDictionaryItem

	query := r.db.Model(&model.OrbiaDictionaryItem{}).
		Joins("JOIN orbia_dictionary ON orbia_dictionary.id = orbia_dictionary_item.dictionary_id").
		Where("orbia_dictionary.code = ?", dictionaryCode)

	// 只获取启用的字典和字典项
	if onlyEnabled {
		query = query.Where("orbia_dictionary.status = 1 AND orbia_dictionary_item.status = 1")
	}

	err := query.Order("orbia_dictionary_item.sort_order ASC, orbia_dictionary_item.created_at ASC").
		Find(&items).Error

	return items, err
}

// CheckDictionaryItemCodeExists 检查字典项编码是否存在（同一字典下同一父级下不能重复）
func (r *dictionaryItemRepository) CheckDictionaryItemCodeExists(dictionaryID, parentID int64, code string, excludeID int64) (bool, error) {
	var count int64
	query := r.db.Model(&model.OrbiaDictionaryItem{}).
		Where("dictionary_id = ? AND parent_id = ? AND code = ?", dictionaryID, parentID, code)
	if excludeID > 0 {
		query = query.Where("id != ?", excludeID)
	}
	err := query.Count(&count).Error
	return count > 0, err
}

// GetChildrenByParentID 获取某个节点的所有子节点
func (r *dictionaryItemRepository) GetChildrenByParentID(parentID int64) ([]*model.OrbiaDictionaryItem, error) {
	var items []*model.OrbiaDictionaryItem
	err := r.db.Where("parent_id = ?", parentID).
		Order("sort_order ASC, created_at ASC").
		Find(&items).Error
	return items, err
}
