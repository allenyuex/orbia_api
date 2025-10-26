package recharge_order

import (
	"errors"
	"fmt"
	"time"

	"orbia_api/biz/dal/model"
	"orbia_api/biz/dal/mysql"
	"orbia_api/biz/utils"

	"gorm.io/gorm"
)

// RechargeOrderService 充值订单服务接口
type RechargeOrderService interface {
	// CreateCryptoRechargeOrder 创建加密货币充值订单
	CreateCryptoRechargeOrder(userID int64, amount float64, paymentSettingID int64, userCryptoAddress string, cryptoTxHash, remark *string) (*model.OrbiaRechargeOrder, error)
	// CreateOnlineRechargeOrder 创建在线支付充值订单
	CreateOnlineRechargeOrder(userID int64, amount float64, platform string) (*model.OrbiaRechargeOrder, string, error)
	// GetMyRechargeOrders 获取用户自己的充值订单列表
	GetMyRechargeOrders(userID int64, status *string, page, pageSize int) ([]*model.OrbiaRechargeOrder, int64, error)
	// GetRechargeOrderDetail 获取充值订单详情
	GetRechargeOrderDetail(userID int64, orderID string, isAdmin bool) (*model.OrbiaRechargeOrder, error)
	// GetAllRechargeOrders 获取所有充值订单列表（管理员）
	GetAllRechargeOrders(userID *int64, status, paymentType *string, page, pageSize int) ([]*model.OrbiaRechargeOrder, int64, error)
	// ConfirmRechargeOrder 确认充值订单（管理员）
	ConfirmRechargeOrder(adminUserID int64, orderID string, cryptoTxHash, remark *string) (*model.OrbiaRechargeOrder, error)
	// RejectRechargeOrder 拒绝充值订单（管理员）
	RejectRechargeOrder(adminUserID int64, orderID string, failedReason string) (*model.OrbiaRechargeOrder, error)
}

// rechargeOrderService 充值订单服务实现
type rechargeOrderService struct {
	db                 *gorm.DB
	rechargeOrderRepo  mysql.RechargeOrderRepository
	paymentSettingRepo mysql.PaymentSettingRepository
	walletRepo         mysql.WalletRepository
}

// NewRechargeOrderService 创建充值订单服务实例
func NewRechargeOrderService(
	db *gorm.DB,
	rechargeOrderRepo mysql.RechargeOrderRepository,
	paymentSettingRepo mysql.PaymentSettingRepository,
	walletRepo mysql.WalletRepository,
) RechargeOrderService {
	return &rechargeOrderService{
		db:                 db,
		rechargeOrderRepo:  rechargeOrderRepo,
		paymentSettingRepo: paymentSettingRepo,
		walletRepo:         walletRepo,
	}
}

// CreateCryptoRechargeOrder 创建加密货币充值订单
func (s *rechargeOrderService) CreateCryptoRechargeOrder(
	userID int64,
	amount float64,
	paymentSettingID int64,
	userCryptoAddress string,
	cryptoTxHash, remark *string,
) (*model.OrbiaRechargeOrder, error) {
	if amount <= 0 {
		return nil, errors.New("invalid amount")
	}

	if userCryptoAddress == "" {
		return nil, errors.New("user crypto address is required")
	}

	// 获取payment_setting信息并快照
	paymentSetting, err := s.paymentSettingRepo.GetPaymentSettingByID(paymentSettingID)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("payment setting not found")
		}
		return nil, fmt.Errorf("failed to get payment setting: %v", err)
	}

	// 检查payment_setting是否启用
	if paymentSetting.Status != 1 {
		return nil, errors.New("payment setting is not active")
	}

	// 生成订单ID
	timestamp := time.Now().Unix()
	id, err := utils.GetDefaultGenerator().NextID()
	if err != nil {
		return nil, fmt.Errorf("failed to generate order ID: %v", err)
	}

	// 创建充值订单
	paymentType := "crypto"
	status := "pending"
	order := &model.OrbiaRechargeOrder{
		OrderID:           fmt.Sprintf("RCHORD_%d_%d", timestamp, id),
		UserID:            userID,
		Amount:            amount,
		PaymentType:       paymentType,
		PaymentSettingID:  &paymentSettingID,
		PaymentNetwork:    &paymentSetting.Network,
		PaymentAddress:    &paymentSetting.Address,
		PaymentLabel:      &paymentSetting.Label,
		UserCryptoAddress: &userCryptoAddress,
		CryptoTxHash:      cryptoTxHash,
		Status:            status,
		Remark:            remark,
	}

	// 保存充值订单
	if err := s.rechargeOrderRepo.CreateRechargeOrder(order); err != nil {
		return nil, fmt.Errorf("failed to create recharge order: %v", err)
	}

	return order, nil
}

// CreateOnlineRechargeOrder 创建在线支付充值订单
func (s *rechargeOrderService) CreateOnlineRechargeOrder(
	userID int64,
	amount float64,
	platform string,
) (*model.OrbiaRechargeOrder, string, error) {
	if amount <= 0 {
		return nil, "", errors.New("invalid amount")
	}

	// 生成订单ID
	timestamp := time.Now().Unix()
	id, err := utils.GetDefaultGenerator().NextID()
	if err != nil {
		return nil, "", fmt.Errorf("failed to generate order ID: %v", err)
	}

	// 创建充值订单
	paymentType := "online"
	status := "pending"
	order := &model.OrbiaRechargeOrder{
		OrderID:               fmt.Sprintf("RCHORD_%d_%d", timestamp, id),
		UserID:                userID,
		Amount:                amount,
		PaymentType:           paymentType,
		OnlinePaymentPlatform: &platform,
		Status:                status,
	}

	// TODO: 这里应该调用第三方支付平台API创建支付订单
	// 目前只是预留字段，返回空URL
	paymentURL := ""
	// 示例：如果是Stripe，这里会返回Stripe的支付URL
	// if platform == "stripe" {
	//     paymentURL = createStripePaymentSession(order.OrderID, amount)
	// }

	// 保存充值订单
	if err := s.rechargeOrderRepo.CreateRechargeOrder(order); err != nil {
		return nil, "", fmt.Errorf("failed to create recharge order: %v", err)
	}

	return order, paymentURL, nil
}

// GetMyRechargeOrders 获取用户自己的充值订单列表
func (s *rechargeOrderService) GetMyRechargeOrders(
	userID int64,
	status *string,
	page, pageSize int,
) ([]*model.OrbiaRechargeOrder, int64, error) {
	if page < 1 {
		page = 1
	}
	if pageSize < 1 {
		pageSize = 10
	}
	if pageSize > 100 {
		pageSize = 100
	}

	return s.rechargeOrderRepo.GetRechargeOrdersByUserID(userID, status, page, pageSize)
}

// GetRechargeOrderDetail 获取充值订单详情
func (s *rechargeOrderService) GetRechargeOrderDetail(
	userID int64,
	orderID string,
	isAdmin bool,
) (*model.OrbiaRechargeOrder, error) {
	order, err := s.rechargeOrderRepo.GetRechargeOrderByOrderID(orderID)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("recharge order not found")
		}
		return nil, fmt.Errorf("failed to get recharge order: %v", err)
	}

	// 如果不是管理员，检查是否是本人的订单
	if !isAdmin && order.UserID != userID {
		return nil, errors.New("permission denied")
	}

	return order, nil
}

// GetAllRechargeOrders 获取所有充值订单列表（管理员）
func (s *rechargeOrderService) GetAllRechargeOrders(
	userID *int64,
	status, paymentType *string,
	page, pageSize int,
) ([]*model.OrbiaRechargeOrder, int64, error) {
	if page < 1 {
		page = 1
	}
	if pageSize < 1 {
		pageSize = 10
	}
	if pageSize > 100 {
		pageSize = 100
	}

	return s.rechargeOrderRepo.GetAllRechargeOrders(userID, status, paymentType, page, pageSize)
}

// ConfirmRechargeOrder 确认充值订单（管理员）
func (s *rechargeOrderService) ConfirmRechargeOrder(
	adminUserID int64,
	orderID string,
	cryptoTxHash, remark *string,
) (*model.OrbiaRechargeOrder, error) {
	// 开始事务
	tx := s.db.Begin()
	if tx.Error != nil {
		return nil, fmt.Errorf("failed to begin transaction: %v", tx.Error)
	}
	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
		}
	}()

	// 获取充值订单
	order, err := s.rechargeOrderRepo.GetRechargeOrderByOrderID(orderID)
	if err != nil {
		tx.Rollback()
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("recharge order not found")
		}
		return nil, fmt.Errorf("failed to get recharge order: %v", err)
	}

	// 检查订单状态
	if order.Status != "pending" {
		tx.Rollback()
		return nil, fmt.Errorf("recharge order is not in pending status, current status: %s", order.Status)
	}

	// 更新订单状态
	now := time.Now()
	order.Status = "confirmed"
	order.ConfirmedBy = &adminUserID
	order.ConfirmedAt = &now

	if cryptoTxHash != nil && *cryptoTxHash != "" {
		order.CryptoTxHash = cryptoTxHash
	}

	if remark != nil && *remark != "" {
		if order.Remark != nil && *order.Remark != "" {
			combinedRemark := *order.Remark + "\n[Admin]: " + *remark
			order.Remark = &combinedRemark
		} else {
			adminRemark := "[Admin]: " + *remark
			order.Remark = &adminRemark
		}
	}

	if err := s.rechargeOrderRepo.UpdateRechargeOrder(order); err != nil {
		tx.Rollback()
		return nil, fmt.Errorf("failed to update recharge order: %v", err)
	}

	// 更新用户钱包余额
	wallet, err := s.walletRepo.GetWalletByUserID(order.UserID)
	if err != nil {
		tx.Rollback()
		return nil, fmt.Errorf("failed to get user wallet: %v", err)
	}

	wallet.Balance += order.Amount
	wallet.TotalRecharge += order.Amount

	if err := s.walletRepo.UpdateWallet(wallet); err != nil {
		tx.Rollback()
		return nil, fmt.Errorf("failed to update wallet balance: %v", err)
	}

	// 提交事务
	if err := tx.Commit().Error; err != nil {
		return nil, fmt.Errorf("failed to commit transaction: %v", err)
	}

	return order, nil
}

// RejectRechargeOrder 拒绝充值订单（管理员）
func (s *rechargeOrderService) RejectRechargeOrder(
	adminUserID int64,
	orderID string,
	failedReason string,
) (*model.OrbiaRechargeOrder, error) {
	if failedReason == "" {
		return nil, errors.New("failed reason is required")
	}

	// 获取充值订单
	order, err := s.rechargeOrderRepo.GetRechargeOrderByOrderID(orderID)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("recharge order not found")
		}
		return nil, fmt.Errorf("failed to get recharge order: %v", err)
	}

	// 检查订单状态
	if order.Status != "pending" {
		return nil, fmt.Errorf("recharge order is not in pending status, current status: %s", order.Status)
	}

	// 更新订单状态
	order.Status = "failed"
	order.FailedReason = &failedReason

	if err := s.rechargeOrderRepo.UpdateRechargeOrder(order); err != nil {
		return nil, fmt.Errorf("failed to update recharge order: %v", err)
	}

	return order, nil
}
