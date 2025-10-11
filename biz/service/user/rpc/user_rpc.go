package rpc

import (
	"errors"
	"fmt"

	"gorm.io/gorm"
	"orbia_api/biz/dal/mysql"
)

// UserRPC 用户RPC服务
type UserRPC struct {
	userRepo mysql.UserRepository
}

// NewUserRPC 创建用户RPC服务实例
func NewUserRPC(userRepo mysql.UserRepository) *UserRPC {
	return &UserRPC{
		userRepo: userRepo,
	}
}

// GetUserByID 根据ID获取用户信息
func (r *UserRPC) GetUserByID(userID int64) (*mysql.User, error) {
	user, err := r.userRepo.GetUserByID(userID)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("user not found")
		}
		return nil, fmt.Errorf("failed to get user: %v", err)
	}

	return user, nil
}

// GetUserByWalletAddress 根据钱包地址获取用户信息
func (r *UserRPC) GetUserByWalletAddress(walletAddress string) (*mysql.User, error) {
	user, err := r.userRepo.GetUserByWalletAddress(walletAddress)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("user not found")
		}
		return nil, fmt.Errorf("failed to get user: %v", err)
	}

	return user, nil
}

// GetUserByEmail 根据邮箱获取用户信息
func (r *UserRPC) GetUserByEmail(email string) (*mysql.User, error) {
	user, err := r.userRepo.GetUserByEmail(email)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("user not found")
		}
		return nil, fmt.Errorf("failed to get user: %v", err)
	}

	return user, nil
}

// CreateUser 创建用户
func (r *UserRPC) CreateUser(user *mysql.User) error {
	if err := r.userRepo.CreateUser(user); err != nil {
		return fmt.Errorf("failed to create user: %v", err)
	}

	return nil
}

// UpdateUser 更新用户信息
func (r *UserRPC) UpdateUser(user *mysql.User) error {
	if err := r.userRepo.UpdateUser(user); err != nil {
		return fmt.Errorf("failed to update user: %v", err)
	}

	return nil
}

// DeleteUser 删除用户
func (r *UserRPC) DeleteUser(userID int64) error {
	if err := r.userRepo.DeleteUser(userID); err != nil {
		return fmt.Errorf("failed to delete user: %v", err)
	}

	return nil
}