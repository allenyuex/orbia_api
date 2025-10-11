package user

import (
	"errors"
	"fmt"

	"gorm.io/gorm"
	"orbia_api/biz/dal/mysql"
)

// UserService 用户服务接口
type UserService interface {
	GetProfile(userID int64) (*mysql.User, error)
	UpdateProfile(userID int64, nickname, avatarURL *string) error
	GetUserByID(userID int64) (*mysql.User, error)
}

// userService 用户服务实现
type userService struct {
	userRepo mysql.UserRepository
}

// NewUserService 创建用户服务实例
func NewUserService(userRepo mysql.UserRepository) UserService {
	return &userService{
		userRepo: userRepo,
	}
}

// GetProfile 获取用户资料
func (s *userService) GetProfile(userID int64) (*mysql.User, error) {
	user, err := s.userRepo.GetUserByID(userID)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("user not found")
		}
		return nil, fmt.Errorf("failed to get user profile: %v", err)
	}

	return user, nil
}

// UpdateProfile 更新用户资料
func (s *userService) UpdateProfile(userID int64, nickname, avatarURL *string) error {
	// 先获取用户信息
	user, err := s.userRepo.GetUserByID(userID)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return errors.New("user not found")
		}
		return fmt.Errorf("failed to get user: %v", err)
	}

	// 更新字段
	if nickname != nil {
		user.Nickname = nickname
	}
	if avatarURL != nil {
		user.AvatarURL = avatarURL
	}

	// 保存更新
	if err := s.userRepo.UpdateUser(user); err != nil {
		return fmt.Errorf("failed to update user profile: %v", err)
	}

	return nil
}

// GetUserByID 根据ID获取用户信息
func (s *userService) GetUserByID(userID int64) (*mysql.User, error) {
	user, err := s.userRepo.GetUserByID(userID)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("user not found")
		}
		return nil, fmt.Errorf("failed to get user: %v", err)
	}

	return user, nil
}