package mysql

import (
	"errors"
	"fmt"

	"orbia_api/biz/dal/model"

	"gorm.io/gorm"
)

// WalletRepository 钱包仓库接口
type WalletRepository interface {
	CreateWallet(wallet *model.OrbiaWallet) error
	GetWalletByUserID(userID int64) (*model.OrbiaWallet, error)
	UpdateWallet(wallet *model.OrbiaWallet) error
	UpdateBalance(tx *gorm.DB, userID int64, balanceDelta float64, frozenDelta float64) error
}

// TransactionRepository 交易记录仓库接口
type TransactionRepository interface {
	CreateTransaction(tx *gorm.DB, transaction *model.OrbiaTransaction) error
	GetTransactionByID(transactionID string) (*model.OrbiaTransaction, error)
	UpdateTransaction(transaction *model.OrbiaTransaction) error
	GetTransactionsByUserID(userID int64, txType, status *string, page, pageSize int) ([]*model.OrbiaTransaction, int64, error)
}

// walletRepository 钱包仓库实现
type walletRepository struct {
	db *gorm.DB
}

// NewWalletRepository 创建钱包仓库实例
func NewWalletRepository(db *gorm.DB) WalletRepository {
	return &walletRepository{db: db}
}

// CreateWallet 创建钱包
func (r *walletRepository) CreateWallet(wallet *model.OrbiaWallet) error {
	return r.db.Create(wallet).Error
}

// GetWalletByUserID 根据用户ID获取钱包
func (r *walletRepository) GetWalletByUserID(userID int64) (*model.OrbiaWallet, error) {
	var wallet model.OrbiaWallet
	err := r.db.Where("user_id = ?", userID).First(&wallet).Error
	if err != nil {
		return nil, err
	}
	return &wallet, nil
}

// UpdateWallet 更新钱包
func (r *walletRepository) UpdateWallet(wallet *model.OrbiaWallet) error {
	return r.db.Save(wallet).Error
}

// UpdateBalance 更新余额（在事务中执行）
func (r *walletRepository) UpdateBalance(tx *gorm.DB, userID int64, balanceDelta float64, frozenDelta float64) error {
	if tx == nil {
		tx = r.db
	}

	// 使用乐观锁更新余额
	result := tx.Model(&model.OrbiaWallet{}).
		Where("user_id = ?", userID).
		Where("balance + ? >= 0", balanceDelta).       // 确保余额不会为负
		Where("frozen_balance + ? >= 0", frozenDelta). // 确保冻结余额不会为负
		Updates(map[string]interface{}{
			"balance":        gorm.Expr("balance + ?", balanceDelta),
			"frozen_balance": gorm.Expr("frozen_balance + ?", frozenDelta),
		})

	if result.Error != nil {
		return result.Error
	}

	if result.RowsAffected == 0 {
		return errors.New("insufficient balance or wallet not found")
	}

	return nil
}

// transactionRepository 交易记录仓库实现
type transactionRepository struct {
	db *gorm.DB
}

// NewTransactionRepository 创建交易记录仓库实例
func NewTransactionRepository(db *gorm.DB) TransactionRepository {
	return &transactionRepository{db: db}
}

// CreateTransaction 创建交易记录（在事务中执行）
func (r *transactionRepository) CreateTransaction(tx *gorm.DB, transaction *model.OrbiaTransaction) error {
	if tx == nil {
		tx = r.db
	}
	return tx.Create(transaction).Error
}

// GetTransactionByID 根据交易ID获取交易记录
func (r *transactionRepository) GetTransactionByID(transactionID string) (*model.OrbiaTransaction, error) {
	var transaction model.OrbiaTransaction
	err := r.db.Where("transaction_id = ?", transactionID).First(&transaction).Error
	if err != nil {
		return nil, err
	}
	return &transaction, nil
}

// UpdateTransaction 更新交易记录
func (r *transactionRepository) UpdateTransaction(transaction *model.OrbiaTransaction) error {
	return r.db.Save(transaction).Error
}

// GetTransactionsByUserID 根据用户ID获取交易记录列表
func (r *transactionRepository) GetTransactionsByUserID(userID int64, txType, status *string, page, pageSize int) ([]*model.OrbiaTransaction, int64, error) {
	var transactions []*model.OrbiaTransaction
	var total int64

	query := r.db.Model(&model.OrbiaTransaction{}).Where("user_id = ?", userID)

	if txType != nil && *txType != "" {
		query = query.Where("type = ?", *txType)
	}

	if status != nil && *status != "" {
		query = query.Where("status = ?", *status)
	}

	// 获取总数
	if err := query.Count(&total).Error; err != nil {
		return nil, 0, fmt.Errorf("failed to count transactions: %v", err)
	}

	// 分页查询
	offset := (page - 1) * pageSize
	if err := query.Order("created_at DESC").Limit(pageSize).Offset(offset).Find(&transactions).Error; err != nil {
		return nil, 0, fmt.Errorf("failed to query transactions: %v", err)
	}

	return transactions, total, nil
}
