package user

import (
	"errors"
	"fmt"

	"orbia_api/biz/dal/mysql"
	"orbia_api/biz/utils"

	"gorm.io/gorm"
)

// KolInfo KOL完整信息
type KolInfo struct {
	Kol       *mysql.Kol
	Languages []*mysql.KolLanguage
	Tags      []*mysql.KolTag
	Stats     *mysql.KolStats
}

// UserService 用户服务接口
type UserService interface {
	GetProfile(userID int64) (*mysql.User, *mysql.Team, *KolInfo, error)
	UpdateProfile(userID int64, nickname, avatarURL *string) error
	GetUserByID(userID int64) (*mysql.User, error)
	SwitchCurrentTeam(userID int64, teamID int64) (*mysql.Team, error)
}

// userService 用户服务实现
type userService struct {
	userRepo mysql.UserRepository
	teamRepo mysql.TeamRepository
	kolRepo  mysql.KolRepository
}

// NewUserService 创建用户服务实例
func NewUserService(userRepo mysql.UserRepository, teamRepo mysql.TeamRepository, kolRepo mysql.KolRepository) UserService {
	return &userService{
		userRepo: userRepo,
		teamRepo: teamRepo,
		kolRepo:  kolRepo,
	}
}

// GetProfile 获取用户资料
func (s *userService) GetProfile(userID int64) (*mysql.User, *mysql.Team, *KolInfo, error) {
	user, err := s.userRepo.GetUserByID(userID)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil, nil, errors.New("user not found")
		}
		return nil, nil, nil, fmt.Errorf("failed to get user profile: %v", err)
	}

	// 获取用户的当前团队信息
	var currentTeam *mysql.Team
	if user.CurrentTeamID != nil {
		team, err := s.teamRepo.GetTeamByID(*user.CurrentTeamID)
		if err != nil {
			// 如果团队不存在，记录日志但不返回错误
			if !errors.Is(err, gorm.ErrRecordNotFound) {
				return user, nil, nil, fmt.Errorf("failed to get current team: %v", err)
			}
		} else {
			currentTeam = team
		}
	}

	// 获取用户的 KOL 信息
	var kolInfo *KolInfo
	if user.KolID != nil {
		kol, err := s.kolRepo.GetKolByID(*user.KolID)
		if err != nil {
			// 如果KOL不存在，记录日志但不返回错误
			if !errors.Is(err, gorm.ErrRecordNotFound) {
				return user, currentTeam, nil, fmt.Errorf("failed to get kol info: %v", err)
			}
		} else {
			// 获取KOL的详细信息
			languages, _ := s.kolRepo.GetKolLanguages(kol.ID)
			tags, _ := s.kolRepo.GetKolTags(kol.ID)
			stats, _ := s.kolRepo.GetKolStats(kol.ID)

			kolInfo = &KolInfo{
				Kol:       kol,
				Languages: languages,
				Tags:      tags,
				Stats:     stats,
			}
		}
	}

	return user, currentTeam, kolInfo, nil
}

// UpdateProfile 更新用户资料
func (s *userService) UpdateProfile(userID int64, nickname, avatarURL *string) error {
	// 获取用户信息
	user, err := s.userRepo.GetUserByID(userID)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return errors.New("user not found")
		}
		return fmt.Errorf("failed to get user: %v", err)
	}

	// 验证头像URL（如果提供）
	if avatarURL != nil && *avatarURL != "" {
		isValid, errorMessage := utils.ValidateFileURL(*avatarURL)
		if !isValid {
			return fmt.Errorf("invalid avatar URL: %s", errorMessage)
		}

		// 检查图片是否存在
		if !utils.CheckFileExists(*avatarURL) {
			return errors.New("avatar image does not exist or is not accessible")
		}
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

// SwitchCurrentTeam 切换用户当前团队
func (s *userService) SwitchCurrentTeam(userID int64, teamID int64) (*mysql.Team, error) {
	// 验证用户是否存在
	user, err := s.userRepo.GetUserByID(userID)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("user not found")
		}
		return nil, fmt.Errorf("failed to get user: %v", err)
	}

	// 验证团队是否存在
	team, err := s.teamRepo.GetTeamByID(teamID)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("team not found")
		}
		return nil, fmt.Errorf("failed to get team: %v", err)
	}

	// 验证用户是否是团队成员
	member, err := s.teamRepo.GetTeamMember(teamID, userID)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("user is not a member of this team")
		}
		return nil, fmt.Errorf("failed to check team membership: %v", err)
	}
	if member == nil {
		return nil, errors.New("user is not a member of this team")
	}

	// 更新用户的当前团队ID
	user.CurrentTeamID = &teamID
	if err := s.userRepo.UpdateUser(user); err != nil {
		return nil, fmt.Errorf("failed to update user current team: %v", err)
	}

	return team, nil
}
