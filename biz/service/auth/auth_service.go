package auth

import (
	"errors"
	"fmt"
	"time"

	"orbia_api/biz/consts"
	"orbia_api/biz/dal/mysql"
	"orbia_api/biz/infra/config"
	walletService "orbia_api/biz/service/wallet"
	"orbia_api/biz/utils"

	"gorm.io/gorm"
)

// AuthService 认证服务接口
type AuthService interface {
	WalletLogin(walletAddress, signature, message string) (string, int64, error)
	SendVerificationCode(email, codeType string) error
	EmailLogin(email, code string) (string, int64, error)
}

// authService 认证服务实现
type authService struct {
	userRepo         mysql.UserRepository
	teamRepo         mysql.TeamRepository
	walletSvc        walletService.WalletService
	verificationRepo mysql.VerificationCodeRepository
}

// NewAuthService 创建认证服务实例
func NewAuthService(userRepo mysql.UserRepository, teamRepo mysql.TeamRepository, walletSvc walletService.WalletService, verificationRepo mysql.VerificationCodeRepository) AuthService {
	return &authService{
		userRepo:         userRepo,
		teamRepo:         teamRepo,
		walletSvc:        walletSvc,
		verificationRepo: verificationRepo,
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
				Nickname:      &walletAddress,            // 使用钱包地址作为 nickname
				AvatarURL:     &defaultAvatar,            // 设置默认头像
				Role:          string(consts.RoleNormal), // 设置为普通用户角色
			}

			if err := s.userRepo.CreateUser(user); err != nil {
				return "", 0, fmt.Errorf("failed to create user: %v", err)
			}

			// 为新用户初始化账户
			if err := s.initializeNewUser(user.ID); err != nil {
				return "", 0, err
			}
		} else {
			return "", 0, fmt.Errorf("failed to query user: %v", err)
		}
	}

	// 生成JWT token并返回
	return s.generateTokenForUser(user.ID)
}

// initializeNewUser 为新用户初始化账户（创建默认团队和钱包）
func (s *authService) initializeNewUser(userID int64) error {
	// 为新用户创建默认团队
	if err := s.createDefaultTeam(userID); err != nil {
		return fmt.Errorf("failed to create default team: %v", err)
	}

	// 为新用户创建钱包
	if err := s.walletSvc.CreateWallet(userID); err != nil {
		return fmt.Errorf("failed to create wallet: %v", err)
	}

	return nil
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

// generateTokenForUser 为用户生成JWT token
func (s *authService) generateTokenForUser(userID int64) (string, int64, error) {
	token, expiresIn, err := utils.GenerateToken(userID)
	if err != nil {
		return "", 0, fmt.Errorf("failed to generate token: %v", err)
	}
	return token, expiresIn, nil
}

// SendVerificationCode 发送验证码
func (s *authService) SendVerificationCode(email, codeType string) error {
	// 验证邮箱格式
	if !utils.ValidateEmail(email) {
		return errors.New("invalid email format")
	}

	// 如果没有指定类型，默认为登录
	if codeType == "" {
		codeType = "login"
	}

	// 验证codeType合法性
	if codeType != "login" && codeType != "register" && codeType != "reset_password" {
		return errors.New("invalid code type")
	}

	// 生成验证码
	cfg := config.GlobalConfig.VerificationCode
	code := utils.GenerateVerificationCode(cfg.Length)

	// 计算过期时间
	expiresAt := time.Now().Add(time.Duration(cfg.ExpireMinutes) * time.Minute)

	// 保存验证码到数据库
	verificationCode := &mysql.VerificationCode{
		Email:     email,
		Code:      code,
		CodeType:  codeType,
		Status:    "unused",
		ExpiresAt: expiresAt,
	}

	if err := s.verificationRepo.CreateVerificationCode(verificationCode); err != nil {
		return fmt.Errorf("failed to save verification code: %v", err)
	}

	// 发送验证码邮件
	//if err := utils.SendVerificationEmail(email, code, cfg.ExpireMinutes); err != nil {
	//	return fmt.Errorf("failed to send verification email: %v", err)
	//}

	return nil
}

// EmailLogin 邮箱验证码登录
func (s *authService) EmailLogin(email, code string) (string, int64, error) {
	// 验证邮箱格式
	if !utils.ValidateEmail(email) {
		return "", 0, errors.New("invalid email format")
	}

	// 验证码为空检查
	if code == "" {
		return "", 0, errors.New("verification code is required")
	}

	// Debug模式：如果验证码是888888，直接通过验证
	if code != "888888" {
		// 验证验证码
		verificationCode, err := s.verificationRepo.GetValidVerificationCode(email, code, "login")
		if err != nil {
			if errors.Is(err, gorm.ErrRecordNotFound) {
				return "", 0, errors.New("invalid or expired verification code")
			}
			return "", 0, fmt.Errorf("failed to verify code: %v", err)
		}

		// 标记验证码为已使用
		if err := s.verificationRepo.MarkAsUsed(verificationCode.ID); err != nil {
			return "", 0, fmt.Errorf("failed to mark code as used: %v", err)
		}
	}

	// 查找用户
	user, err := s.userRepo.GetUserByEmail(email)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			// 用户不存在，自动注册
			defaultAvatar := consts.DefaultAvatarURL
			user = &mysql.User{
				Email:     &email,
				Nickname:  &email,                    // 使用邮箱作为 nickname
				AvatarURL: &defaultAvatar,            // 设置默认头像
				Role:      string(consts.RoleNormal), // 设置为普通用户角色
			}

			if err := s.userRepo.CreateUser(user); err != nil {
				return "", 0, fmt.Errorf("failed to create user: %v", err)
			}

			// 为新用户初始化账户
			if err := s.initializeNewUser(user.ID); err != nil {
				return "", 0, err
			}
		} else {
			return "", 0, fmt.Errorf("failed to query user: %v", err)
		}
	}

	// 生成JWT token并返回
	return s.generateTokenForUser(user.ID)
}
