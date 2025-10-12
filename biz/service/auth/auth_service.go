package auth

import (
	"errors"
	"fmt"
	"time"

	"orbia_api/biz/consts"
	"orbia_api/biz/dal/mysql"
	"orbia_api/biz/utils"

	"gorm.io/gorm"
)

// AuthService 认证服务接口
type AuthService interface {
	WalletLogin(walletAddress, signature, message string) (string, int64, error)
	EmailLogin(email, password string) (string, int64, error)
}

// authService 认证服务实现
type authService struct {
	userRepo mysql.UserRepository
	teamRepo mysql.TeamRepository
}

// NewAuthService 创建认证服务实例
func NewAuthService(userRepo mysql.UserRepository, teamRepo mysql.TeamRepository) AuthService {
	return &authService{
		userRepo: userRepo,
		teamRepo: teamRepo,
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
				Nickname:      &walletAddress,          // 使用钱包地址作为 nickname
				AvatarURL:     &defaultAvatar,          // 设置默认头像
				Role:          string(consts.RoleUser), // 设置为普通用户角色
			}

			if err := s.userRepo.CreateUser(user); err != nil {
				return "", 0, fmt.Errorf("failed to create user: %v", err)
			}

			// 为新用户创建默认项目
			if err := s.createDefaultTeam(user.ID); err != nil {
				return "", 0, fmt.Errorf("failed to create default team: %v", err)
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

// createDefaultTeam 为新用户创建默认团队
func (s *authService) createDefaultTeam(userID int64) error {
	defaultTeamName := "default team"
	defaultTeamIcon := "https://raw.githubusercontent.com/Tarikul-Islam-Anik/Animated-Fluent-Emojis/master/Emojis/Travel%20and%20places/Star.png"

	// 创建默认团队
	team := &mysql.Team{
		Name:      defaultTeamName,
		IconURL:   &defaultTeamIcon,
		CreatorID: userID,
	}

	if err := s.teamRepo.CreateTeam(team); err != nil {
		return fmt.Errorf("failed to create default team: %v", err)
	}

	// 创建团队成员关系，设置用户为创建者
	member := &mysql.TeamMember{
		TeamID: team.ID,
		UserID: userID,
		Role:   "creator",
	}

	if err := s.teamRepo.AddTeamMember(member); err != nil {
		return fmt.Errorf("failed to create team member: %v", err)
	}

	// 更新用户的当前团队ID
	user, err := s.userRepo.GetUserByID(userID)
	if err != nil {
		return fmt.Errorf("failed to get user: %v", err)
	}

	user.CurrentTeamID = &team.ID
	if err := s.userRepo.UpdateUser(user); err != nil {
		return fmt.Errorf("failed to update user current team: %v", err)
	}

	return nil
}

// EmailLogin 邮箱登录
func (s *authService) EmailLogin(email, password string) (string, int64, error) {
	// 查找用户
	user, err := s.userRepo.GetUserByEmail(email)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			// 用户不存在，自动注册
			passwordHash, err := utils.HashPassword(password)
			if err != nil {
				return "", 0, fmt.Errorf("failed to hash password: %v", err)
			}

			defaultAvatar := consts.DefaultAvatarURL
			user = &mysql.User{
				Email:        &email,
				PasswordHash: &passwordHash,
				Nickname:     &email,                  // 使用邮箱作为 nickname
				AvatarURL:    &defaultAvatar,          // 设置默认头像
				Role:         string(consts.RoleUser), // 设置为普通用户角色
			}

			if err := s.userRepo.CreateUser(user); err != nil {
				return "", 0, fmt.Errorf("failed to create user: %v", err)
			}

			// 为新用户创建默认项目
			if err := s.createDefaultTeam(user.ID); err != nil {
				return "", 0, fmt.Errorf("failed to create default team: %v", err)
			}
		} else {
			return "", 0, fmt.Errorf("failed to query user: %v", err)
		}
	} else {
		// 用户存在，验证密码
		if user.PasswordHash == nil {
			return "", 0, errors.New("password not set for this user")
		}

		if !utils.CheckPasswordHash(password, *user.PasswordHash) {
			return "", 0, errors.New("invalid password")
		}
	}

	// 生成JWT token
	token, expiresIn, err := utils.GenerateToken(user.ID)
	if err != nil {
		return "", 0, fmt.Errorf("failed to generate token: %v", err)
	}

	return token, expiresIn, nil
}
