package wallet

import (
	"errors"
	"fmt"
	"time"

	"orbia_api/biz/dal/model"
	"orbia_api/biz/dal/mysql"
	"orbia_api/biz/utils"

	"gorm.io/gorm"
)

// WalletService 钱包服务接口
type WalletService interface {
	// CreateWallet 创建钱包
	CreateWallet(userID int64) error
	// GetWalletInfo 获取钱包信息
	GetWalletInfo(userID int64) (*model.OrbiaWallet, error)
	// CryptoRecharge 加密货币充值
	CryptoRecharge(userID int64, amount float64, cryptoCurrency, cryptoChain, cryptoAddress string) (*model.OrbiaTransaction, error)
	// OnlineRecharge 在线支付充值
	OnlineRecharge(userID int64, amount float64, platform string) (*model.OrbiaTransaction, string, error)
	// ConfirmCryptoRecharge 确认加密货币充值
	ConfirmCryptoRecharge(userID int64, transactionID, cryptoTxHash string) (*model.OrbiaTransaction, error)
	// GetTransactionList 获取交易记录列表
	GetTransactionList(userID int64, txType, status *string, page, pageSize int) ([]*model.OrbiaTransaction, int64, error)
	// GetTransactionDetail 获取交易详情
	GetTransactionDetail(userID int64, transactionID string) (*model.OrbiaTransaction, error)
}

// walletService 钱包服务实现
type walletService struct {
	db         *gorm.DB
	walletRepo mysql.WalletRepository
	txRepo     mysql.TransactionRepository
}

// NewWalletService 创建钱包服务实例
func NewWalletService(db *gorm.DB, walletRepo mysql.WalletRepository, txRepo mysql.TransactionRepository) WalletService {
	return &walletService{
		db:         db,
		walletRepo: walletRepo,
		txRepo:     txRepo,
	}
}

// CreateWallet 创建钱包
func (s *walletService) CreateWallet(userID int64) error {
	// 检查用户是否已有钱包
	_, err := s.walletRepo.GetWalletByUserID(userID)
	if err == nil {
		// 钱包已存在，直接返回
		return nil
	}
	if !errors.Is(err, gorm.ErrRecordNotFound) {
		return fmt.Errorf("failed to check wallet: %v", err)
	}

	// 创建新钱包
	wallet := &model.OrbiaWallet{
		UserID:        userID,
		Balance:       0,
		FrozenBalance: 0,
		TotalRecharge: 0,
		TotalConsume:  0,
	}

	if err := s.walletRepo.CreateWallet(wallet); err != nil {
		return fmt.Errorf("failed to create wallet: %v", err)
	}

	return nil
}

// GetWalletInfo 获取钱包信息
func (s *walletService) GetWalletInfo(userID int64) (*model.OrbiaWallet, error) {
	wallet, err := s.walletRepo.GetWalletByUserID(userID)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("wallet not found")
		}
		return nil, fmt.Errorf("failed to get wallet: %v", err)
	}

	return wallet, nil
}

// CryptoRecharge 加密货币充值
func (s *walletService) CryptoRecharge(userID int64, amount float64, cryptoCurrency, cryptoChain, cryptoAddress string) (*model.OrbiaTransaction, error) {
	if amount <= 0 {
		return nil, errors.New("invalid amount")
	}

	// 获取钱包信息
	wallet, err := s.walletRepo.GetWalletByUserID(userID)
	if err != nil {
		return nil, fmt.Errorf("failed to get wallet: %v", err)
	}

	// 生成交易ID
	transactionID, err := utils.GenerateID()
	if err != nil {
		return nil, fmt.Errorf("failed to generate transaction ID: %v", err)
	}

	// 创建交易记录
	paymentMethod := "crypto"
	txType := "recharge"
	status := "pending"
	transaction := &model.OrbiaTransaction{
		TransactionID:  fmt.Sprintf("TX%d", transactionID),
		UserID:         userID,
		Type:           txType,
		Amount:         amount,
		BalanceBefore:  wallet.Balance,
		BalanceAfter:   wallet.Balance, // 待确认前余额不变
		Status:         status,
		PaymentMethod:  &paymentMethod,
		CryptoCurrency: &cryptoCurrency,
		CryptoChain:    &cryptoChain,
		CryptoAddress:  &cryptoAddress,
	}

	// 保存交易记录
	if err := s.txRepo.CreateTransaction(nil, transaction); err != nil {
		return nil, fmt.Errorf("failed to create transaction: %v", err)
	}

	return transaction, nil
}

// OnlineRecharge 在线支付充值
func (s *walletService) OnlineRecharge(userID int64, amount float64, platform string) (*model.OrbiaTransaction, string, error) {
	if amount <= 0 {
		return nil, "", errors.New("invalid amount")
	}

	// 获取钱包信息
	wallet, err := s.walletRepo.GetWalletByUserID(userID)
	if err != nil {
		return nil, "", fmt.Errorf("failed to get wallet: %v", err)
	}

	// 生成交易ID
	transactionID, err := utils.GenerateID()
	if err != nil {
		return nil, "", fmt.Errorf("failed to generate transaction ID: %v", err)
	}

	// 生成支付订单ID（这里简单使用交易ID，实际应该调用支付平台API）
	paymentOrderID := fmt.Sprintf("PAY%d", transactionID)

	// 生成支付URL（这里是模拟，实际应该调用支付平台API）
	paymentURL := fmt.Sprintf("https://payment.example.com/%s/%s", platform, paymentOrderID)

	// 创建交易记录
	paymentMethod := "online"
	txType := "recharge"
	status := "pending"
	transaction := &model.OrbiaTransaction{
		TransactionID:         fmt.Sprintf("TX%d", transactionID),
		UserID:                userID,
		Type:                  txType,
		Amount:                amount,
		BalanceBefore:         wallet.Balance,
		BalanceAfter:          wallet.Balance, // 待确认前余额不变
		Status:                status,
		PaymentMethod:         &paymentMethod,
		OnlinePaymentPlatform: &platform,
		OnlinePaymentOrderID:  &paymentOrderID,
		OnlinePaymentURL:      &paymentURL,
	}

	// 保存交易记录
	if err := s.txRepo.CreateTransaction(nil, transaction); err != nil {
		return nil, "", fmt.Errorf("failed to create transaction: %v", err)
	}

	return transaction, paymentURL, nil
}

// ConfirmCryptoRecharge 确认加密货币充值
func (s *walletService) ConfirmCryptoRecharge(userID int64, transactionID, cryptoTxHash string) (*model.OrbiaTransaction, error) {
	// 获取交易记录
	transaction, err := s.txRepo.GetTransactionByID(transactionID)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("transaction not found")
		}
		return nil, fmt.Errorf("failed to get transaction: %v", err)
	}

	// 验证交易所属用户
	if transaction.UserID != userID {
		return nil, errors.New("transaction does not belong to this user")
	}

	// 验证交易状态
	if transaction.Status != "pending" {
		return nil, errors.New("transaction is not in pending status")
	}

	// 验证交易类型
	if transaction.Type != "recharge" {
		return nil, errors.New("transaction is not a recharge")
	}

	// TODO: 实际应该验证链上交易哈希的真实性

	// 开始事务
	return transaction, s.db.Transaction(func(tx *gorm.DB) error {
		// 更新余额
		if err := s.walletRepo.UpdateBalance(tx, userID, transaction.Amount, 0); err != nil {
			return fmt.Errorf("failed to update balance: %v", err)
		}

		// 获取更新后的余额
		wallet, err := s.walletRepo.GetWalletByUserID(userID)
		if err != nil {
			return fmt.Errorf("failed to get wallet: %v", err)
		}

		// 更新交易状态
		now := time.Now()
		transaction.Status = "completed"
		transaction.CryptoTxHash = &cryptoTxHash
		transaction.BalanceAfter = wallet.Balance
		transaction.CompletedAt = &now

		if err := s.txRepo.UpdateTransaction(transaction); err != nil {
			return fmt.Errorf("failed to update transaction: %v", err)
		}

		return nil
	})
}

// GetTransactionList 获取交易记录列表
func (s *walletService) GetTransactionList(userID int64, txType, status *string, page, pageSize int) ([]*model.OrbiaTransaction, int64, error) {
	return s.txRepo.GetTransactionsByUserID(userID, txType, status, page, pageSize)
}

// GetTransactionDetail 获取交易详情
func (s *walletService) GetTransactionDetail(userID int64, transactionID string) (*model.OrbiaTransaction, error) {
	transaction, err := s.txRepo.GetTransactionByID(transactionID)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("transaction not found")
		}
		return nil, fmt.Errorf("failed to get transaction: %v", err)
	}

	// 验证交易所属用户
	if transaction.UserID != userID {
		return nil, errors.New("transaction does not belong to this user")
	}

	return transaction, nil
}
