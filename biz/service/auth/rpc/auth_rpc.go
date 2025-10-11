package rpc

import (
	"errors"
	"fmt"

	"gorm.io/gorm"
	"orbia_api/biz/dal/mysql"
	"orbia_api/biz/utils"
)

// AuthRPC 认证RPC服务
type AuthRPC struct {
	userRepo mysql.UserRepository
}

// NewAuthRPC 创建认证RPC服务实例
func NewAuthRPC(userRepo mysql.UserRepository) *AuthRPC {
	return &AuthRPC{
		userRepo: userRepo,
	}
}

// ValidateUserToken 验证用户token并返回用户信息
func (r *AuthRPC) ValidateUserToken(token string) (*mysql.User, error) {
	// 验证token
	userID, err := utils.ValidateToken(token)
	if err != nil {
		return nil, fmt.Errorf("invalid token: %v", err)
	}

	// 获取用户信息
	user, err := r.userRepo.GetUserByID(userID)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("user not found")
		}
		return nil, fmt.Errorf("failed to get user: %v", err)
	}

	return user, nil
}

// GetUserByID 根据ID获取用户信息
func (r *AuthRPC) GetUserByID(userID int64) (*mysql.User, error) {
	user, err := r.userRepo.GetUserByID(userID)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("user not found")
		}
		return nil, fmt.Errorf("failed to get user: %v", err)
	}

	return user, nil
}

// CheckWalletExists 检查钱包地址是否已存在
func (r *AuthRPC) CheckWalletExists(walletAddress string) (bool, error) {
	_, err := r.userRepo.GetUserByWalletAddress(walletAddress)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return false, nil
		}
		return false, fmt.Errorf("failed to check wallet: %v", err)
	}
	return true, nil
}

// CheckEmailExists 检查邮箱是否已存在
func (r *AuthRPC) CheckEmailExists(email string) (bool, error) {
	_, err := r.userRepo.GetUserByEmail(email)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return false, nil
		}
		return false, fmt.Errorf("failed to check email: %v", err)
	}
	return true, nil
}