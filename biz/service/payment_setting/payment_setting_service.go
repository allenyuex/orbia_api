package payment_setting

import (
	"context"
	"errors"

	"orbia_api/biz/dal/model"
	"orbia_api/biz/dal/mysql"
	psModel "orbia_api/biz/model/payment_setting"
	"orbia_api/biz/utils"

	"github.com/cloudwego/hertz/pkg/common/hlog"
	"gorm.io/gorm"
)

// PaymentSettingService 收款钱包设置服务
type PaymentSettingService struct {
	repo mysql.PaymentSettingRepository
}

// NewPaymentSettingService 创建收款钱包设置服务实例
func NewPaymentSettingService(repo mysql.PaymentSettingRepository) *PaymentSettingService {
	return &PaymentSettingService{
		repo: repo,
	}
}

// ==================== 管理员接口 ====================

// GetPaymentSettingList 获取收款钱包设置列表
func (s *PaymentSettingService) GetPaymentSettingList(ctx context.Context, req *psModel.GetPaymentSettingListReq) (*psModel.GetPaymentSettingListResp, error) {
	// 设置默认分页参数
	page := int32(1)
	pageSize := int32(20)
	if req.Page != nil && *req.Page > 0 {
		page = *req.Page
	}
	if req.PageSize != nil && *req.PageSize > 0 {
		pageSize = *req.PageSize
	}

	// 计算偏移量
	offset := int((page - 1) * pageSize)
	limit := int(pageSize)

	// 查询列表
	network := ""
	if req.Network != nil {
		network = *req.Network
	}

	settings, total, err := s.repo.GetPaymentSettings(network, req.Status, offset, limit)
	if err != nil {
		hlog.Errorf("Failed to get payment settings: %v", err)
		return &psModel.GetPaymentSettingListResp{
			BaseResp: utils.BuildBaseResp(500, "获取收款钱包设置列表失败"),
		}, nil
	}

	// 构建响应
	list := make([]*psModel.PaymentSetting, 0, len(settings))
	for _, setting := range settings {
		list = append(list, s.buildPaymentSettingInfo(setting))
	}

	return &psModel.GetPaymentSettingListResp{
		List:     list,
		Total:    total,
		Page:     page,
		PageSize: pageSize,
		BaseResp: utils.BuildSuccessResp(),
	}, nil
}

// GetPaymentSettingDetail 获取收款钱包设置详情
func (s *PaymentSettingService) GetPaymentSettingDetail(ctx context.Context, req *psModel.GetPaymentSettingDetailReq) (*psModel.GetPaymentSettingDetailResp, error) {
	setting, err := s.repo.GetPaymentSettingByID(req.ID)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return &psModel.GetPaymentSettingDetailResp{
				BaseResp: utils.BuildBaseResp(404, "收款钱包设置不存在"),
			}, nil
		}
		hlog.Errorf("Failed to get payment setting: %v", err)
		return &psModel.GetPaymentSettingDetailResp{
			BaseResp: utils.BuildBaseResp(500, "获取收款钱包设置失败"),
		}, nil
	}

	return &psModel.GetPaymentSettingDetailResp{
		Setting:  s.buildPaymentSettingInfo(setting),
		BaseResp: utils.BuildSuccessResp(),
	}, nil
}

// CreatePaymentSetting 创建收款钱包设置
func (s *PaymentSettingService) CreatePaymentSetting(ctx context.Context, req *psModel.CreatePaymentSettingReq) (*psModel.CreatePaymentSettingResp, error) {
	// 参数验证
	if req.Network == "" {
		return &psModel.CreatePaymentSettingResp{
			BaseResp: utils.BuildBaseResp(400, "区块链网络不能为空"),
		}, nil
	}
	if req.Address == "" {
		return &psModel.CreatePaymentSettingResp{
			BaseResp: utils.BuildBaseResp(400, "钱包地址不能为空"),
		}, nil
	}
	if req.Label == "" {
		return &psModel.CreatePaymentSettingResp{
			BaseResp: utils.BuildBaseResp(400, "钱包标签不能为空"),
		}, nil
	}

	// 设置状态默认值
	status := int32(1)
	if req.Status != nil {
		status = *req.Status
	}

	// 创建设置
	setting := &model.OrbiaPaymentSetting{
		Network: req.Network,
		Address: req.Address,
		Label:   req.Label,
		Status:  status,
	}

	if err := s.repo.CreatePaymentSetting(setting); err != nil {
		hlog.Errorf("Failed to create payment setting: %v", err)
		return &psModel.CreatePaymentSettingResp{
			BaseResp: utils.BuildBaseResp(500, "创建收款钱包设置失败"),
		}, nil
	}

	return &psModel.CreatePaymentSettingResp{
		Setting:  s.buildPaymentSettingInfo(setting),
		BaseResp: utils.BuildSuccessResp(),
	}, nil
}

// UpdatePaymentSetting 更新收款钱包设置
func (s *PaymentSettingService) UpdatePaymentSetting(ctx context.Context, req *psModel.UpdatePaymentSettingReq) (*psModel.UpdatePaymentSettingResp, error) {
	// 获取设置
	setting, err := s.repo.GetPaymentSettingByID(req.ID)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return &psModel.UpdatePaymentSettingResp{
				BaseResp: utils.BuildBaseResp(404, "收款钱包设置不存在"),
			}, nil
		}
		hlog.Errorf("Failed to get payment setting: %v", err)
		return &psModel.UpdatePaymentSettingResp{
			BaseResp: utils.BuildBaseResp(500, "获取收款钱包设置失败"),
		}, nil
	}

	// 更新字段
	if req.Network != nil {
		setting.Network = *req.Network
	}
	if req.Address != nil {
		setting.Address = *req.Address
	}
	if req.Label != nil {
		setting.Label = *req.Label
	}
	if req.Status != nil {
		setting.Status = *req.Status
	}

	// 更新设置
	if err := s.repo.UpdatePaymentSetting(setting); err != nil {
		hlog.Errorf("Failed to update payment setting: %v", err)
		return &psModel.UpdatePaymentSettingResp{
			BaseResp: utils.BuildBaseResp(500, "更新收款钱包设置失败"),
		}, nil
	}

	return &psModel.UpdatePaymentSettingResp{
		Setting:  s.buildPaymentSettingInfo(setting),
		BaseResp: utils.BuildSuccessResp(),
	}, nil
}

// DeletePaymentSetting 删除收款钱包设置
func (s *PaymentSettingService) DeletePaymentSetting(ctx context.Context, req *psModel.DeletePaymentSettingReq) (*psModel.DeletePaymentSettingResp, error) {
	// 检查设置是否存在
	_, err := s.repo.GetPaymentSettingByID(req.ID)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return &psModel.DeletePaymentSettingResp{
				BaseResp: utils.BuildBaseResp(404, "收款钱包设置不存在"),
			}, nil
		}
		hlog.Errorf("Failed to get payment setting: %v", err)
		return &psModel.DeletePaymentSettingResp{
			BaseResp: utils.BuildBaseResp(500, "获取收款钱包设置失败"),
		}, nil
	}

	// 删除设置（软删除）
	if err := s.repo.DeletePaymentSetting(req.ID); err != nil {
		hlog.Errorf("Failed to delete payment setting: %v", err)
		return &psModel.DeletePaymentSettingResp{
			BaseResp: utils.BuildBaseResp(500, "删除收款钱包设置失败"),
		}, nil
	}

	return &psModel.DeletePaymentSettingResp{
		BaseResp: utils.BuildSuccessResp(),
	}, nil
}

// ==================== 用户接口 ====================

// GetActivePaymentSettings 获取启用的收款钱包设置列表
func (s *PaymentSettingService) GetActivePaymentSettings(ctx context.Context, req *psModel.GetActivePaymentSettingsReq) (*psModel.GetActivePaymentSettingsResp, error) {
	network := ""
	if req.Network != nil {
		network = *req.Network
	}

	settings, err := s.repo.GetActivePaymentSettings(network)
	if err != nil {
		hlog.Errorf("Failed to get active payment settings: %v", err)
		return &psModel.GetActivePaymentSettingsResp{
			BaseResp: utils.BuildBaseResp(500, "获取启用的收款钱包设置失败"),
		}, nil
	}

	// 构建响应
	list := make([]*psModel.PaymentSetting, 0, len(settings))
	for _, setting := range settings {
		list = append(list, s.buildPaymentSettingInfo(setting))
	}

	return &psModel.GetActivePaymentSettingsResp{
		List:     list,
		BaseResp: utils.BuildSuccessResp(),
	}, nil
}

// ==================== 辅助方法 ====================

// buildPaymentSettingInfo 构建收款钱包设置信息
func (s *PaymentSettingService) buildPaymentSettingInfo(setting *model.OrbiaPaymentSetting) *psModel.PaymentSetting {
	info := &psModel.PaymentSetting{
		ID:      setting.ID,
		Network: setting.Network,
		Address: setting.Address,
		Label:   setting.Label,
		Status:  setting.Status,
	}

	if setting.CreatedAt != nil {
		info.CreatedAt = setting.CreatedAt.Format("2006-01-02 15:04:05")
	}
	if setting.UpdatedAt != nil {
		info.UpdatedAt = setting.UpdatedAt.Format("2006-01-02 15:04:05")
	}

	return info
}
