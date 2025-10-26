package wallet

import (
	"errors"
	"fmt"

	"orbia_api/biz/dal/model"
	"orbia_api/biz/dal/mysql"

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

// CryptoRecharge 加密货币充值（已废弃，请使用充值订单接口）
func (s *walletService) CryptoRecharge(userID int64, amount float64, cryptoCurrency, cryptoChain, cryptoAddress string) (*model.OrbiaTransaction, error) {
	return nil, errors.New("this API is deprecated, please use /api/v1/recharge/create/crypto instead")
}

// OnlineRecharge 在线支付充值（已废弃，请使用充值订单接口）
func (s *walletService) OnlineRecharge(userID int64, amount float64, platform string) (*model.OrbiaTransaction, string, error) {
	return nil, "", errors.New("this API is deprecated, please use /api/v1/recharge/create/online instead")
}

// ConfirmCryptoRecharge 确认加密货币充值（已废弃）
func (s *walletService) ConfirmCryptoRecharge(userID int64, transactionID, cryptoTxHash string) (*model.OrbiaTransaction, error) {
	return nil, errors.New("this API is deprecated, recharge is now managed through recharge orders")
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
