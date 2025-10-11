package auth

import (
	"errors"
	"fmt"
	"time"

	"gorm.io/gorm"
	"orbia_api/biz/consts"
	"orbia_api/biz/dal/mysql"
	"orbia_api/biz/utils"
)

// AuthService 认证服务接口
type AuthService interface {
	WalletLogin(walletAddress, signature, message string) (string, int64, error)
	EmailLogin(email, password string) (string, int64, error)
}

// authService 认证服务实现
type authService struct {
	userRepo mysql.UserRepository
}

// NewAuthService 创建认证服务实例
func NewAuthService(userRepo mysql.UserRepository) AuthService {
	return &authService{
		userRepo: userRepo,
	}
}

// WalletLogin 钱包登录
func (s *authService) WalletLogin(walletAddress, signature, message string) (string, int64, error) {
	// 验证钱包地址格式
	if !utils.ValidateWalletAddress(walletAddress) {
		return "", 0, errors.New("invalid wallet address format")
	}

	// 如果没有提供消息，生成默认消息
	if message == "" {
		message = utils.GenerateSignMessage(walletAddress, time.Now().Unix())
	}

	// 验证签名
	if err := utils.VerifySignature(walletAddress, message, signature); err != nil {
		return "", 0, fmt.Errorf("signature verification failed: %v", err)
	}

	// 查找用户
	user, err := s.userRepo.GetUserByWalletAddress(walletAddress)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			// 用户不存在，自动注册
			defaultAvatar := consts.DefaultAvatarURL
			user = &mysql.User{
				WalletAddress: &walletAddress,
				AvatarURL:     &defaultAvatar, // 设置默认头像
			}
			
			if err := s.userRepo.CreateUser(user); err != nil {
				return "", 0, fmt.Errorf("failed to create user: %v", err)
			}
		} else {
			return "", 0, fmt.Errorf("failed to query user: %v", err)
		}
	}

	// 生成JWT token
	token, expiresIn, err := utils.GenerateToken(user.ID)
	if err != nil {
		return "", 0, fmt.Errorf("failed to generate token: %v", err)
	}

	return token, expiresIn, nil
}

// EmailLogin 邮箱登录（预留实现）
func (s *authService) EmailLogin(email, password string) (string, int64, error) {
	// 查找用户
	user, err := s.userRepo.GetUserByEmail(email)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return "", 0, errors.New("user not found")
		}
		return "", 0, fmt.Errorf("failed to query user: %v", err)
	}

	// 验证密码
	if user.PasswordHash == nil {
		return "", 0, errors.New("password not set for this user")
	}

	if !utils.CheckPasswordHash(password, *user.PasswordHash) {
		return "", 0, errors.New("invalid password")
	}

	// 生成JWT token
	token, expiresIn, err := utils.GenerateToken(user.ID)
	if err != nil {
		return "", 0, fmt.Errorf("failed to generate token: %v", err)
	}

	return token, expiresIn, nil
}