package mysql

import (
	"orbia_api/biz/utils"
	"time"

	"gorm.io/gorm"
)

// Team 团队模型
type Team struct {
	ID        int64          `gorm:"primaryKey;autoIncrement;column:id" json:"id"`
	Name      string         `gorm:"column:name;size:20;not null" json:"name"`
	IconURL   *string        `gorm:"column:icon_url;size:500" json:"icon_url"`
	CreatorID int64          `gorm:"column:creator_id;not null;index" json:"creator_id"`
	CreatedAt time.Time      `gorm:"column:created_at;autoCreateTime" json:"created_at"`
	UpdatedAt time.Time      `gorm:"column:updated_at;autoUpdateTime" json:"updated_at"`
	DeletedAt gorm.DeletedAt `gorm:"index;column:deleted_at" json:"-"`
}

// TableName 指定表名
func (Team) TableName() string {
	return "orbia_team"
}

// TeamMember 团队成员模型
type TeamMember struct {
	ID       int64     `gorm:"primaryKey;autoIncrement;column:id" json:"id"`
	TeamID   int64     `gorm:"column:team_id;not null;index" json:"team_id"`
	UserID   int64     `gorm:"column:user_id;not null;index" json:"user_id"`
	Role     string    `gorm:"column:role;type:enum('creator','owner','member');default:'member'" json:"role"`
	JoinedAt time.Time `gorm:"column:joined_at;autoCreateTime" json:"joined_at"`
}

// TableName 指定表名
func (TeamMember) TableName() string {
	return "orbia_team_member"
}

// TeamInvitation 团队邀请模型
type TeamInvitation struct {
	ID             int64          `gorm:"primaryKey;autoIncrement;column:id" json:"id"`
	TeamID         int64          `gorm:"column:team_id;not null;index" json:"team_id"`
	InviterID      int64          `gorm:"column:inviter_id;not null;index" json:"inviter_id"`
	InviteeEmail   *string        `gorm:"column:invitee_email;size:255;index" json:"invitee_email"`
	InviteeWallet  *string        `gorm:"column:invitee_wallet;size:42;index" json:"invitee_wallet"`
	Role           string         `gorm:"column:role;type:enum('owner','member');default:'member'" json:"role"`
	Status         string         `gorm:"column:status;type:enum('pending','accepted','rejected','expired');default:'pending';index" json:"status"`
	InvitationCode string         `gorm:"column:invitation_code;size:32;not null;index" json:"invitation_code"`
	ExpiresAt      time.Time      `gorm:"column:expires_at;not null;index" json:"expires_at"`
	CreatedAt      time.Time      `gorm:"column:created_at;autoCreateTime" json:"created_at"`
	UpdatedAt      time.Time      `gorm:"column:updated_at;autoUpdateTime" json:"updated_at"`
	DeletedAt      gorm.DeletedAt `gorm:"index;column:deleted_at" json:"-"`
}

// TableName 指定表名
func (TeamInvitation) TableName() string {
	return "orbia_team_invitation"
}

// TeamRepository 团队仓储接口
type TeamRepository interface {
	// 团队相关
	CreateTeam(team *Team) error
	GetTeamByID(id int64) (*Team, error)
	UpdateTeam(team *Team) error
	DeleteTeam(id int64) error
	GetUserTeams(userID int64, offset, limit int) ([]*Team, int64, error)
	// 管理员功能
	GetAllTeams(keyword string, offset int, limit int) ([]*Team, int64, error)
	GetTeamMemberCount(teamID int64) (int64, error)

	// 团队成员相关
	AddTeamMember(member *TeamMember) error
	GetTeamMembers(teamID int64, offset, limit int) ([]*TeamMember, int64, error)
	GetTeamMember(teamID, userID int64) (*TeamMember, error)
	UpdateTeamMember(member *TeamMember) error
	RemoveTeamMember(teamID, userID int64) error
	GetUserRole(teamID, userID int64) (string, error)

	// 团队邀请相关
	CreateInvitation(invitation *TeamInvitation) error
	GetInvitationByCode(code string) (*TeamInvitation, error)
	GetUserInvitations(userID int64, offset, limit int) ([]*TeamInvitation, int64, error)
	UpdateInvitation(invitation *TeamInvitation) error
	ExpireInvitations() error
}

// teamRepository 团队仓储实现
type teamRepository struct {
	db *gorm.DB
}

// NewTeamRepository 创建团队仓储实例
func NewTeamRepository(db *gorm.DB) TeamRepository {
	return &teamRepository{db: db}
}

// CreateTeam 创建团队
func (r *teamRepository) CreateTeam(team *Team) error {
	return r.db.Create(team).Error
}

// GetTeamByID 根据ID获取团队
func (r *teamRepository) GetTeamByID(id int64) (*Team, error) {
	var team Team
	err := r.db.Where("id = ?", id).First(&team).Error
	if err != nil {
		return nil, err
	}
	return &team, nil
}

// UpdateTeam 更新团队
func (r *teamRepository) UpdateTeam(team *Team) error {
	return r.db.Save(team).Error
}

// DeleteTeam 删除团队
func (r *teamRepository) DeleteTeam(id int64) error {
	return r.db.Delete(&Team{}, id).Error
}

// GetUserTeams 获取用户的团队列表
func (r *teamRepository) GetUserTeams(userID int64, offset, limit int) ([]*Team, int64, error) {
	utils.LogDebug("GetUserTeams repository called", map[string]interface{}{
		"user_id": userID,
		"offset":  offset,
		"limit":   limit,
	})

	// 检查 db 是否为 nil
	if r.db == nil {
		utils.LogError(nil, "database connection is nil in GetUserTeams")
		return nil, 0, gorm.ErrInvalidDB
	}

	var teams []*Team
	var total int64

	// 通过团队成员表关联查询用户的团队
	query := r.db.Table("orbia_team t").
		Joins("JOIN orbia_team_member tm ON t.id = tm.team_id").
		Where("tm.user_id = ?", userID)

	utils.LogDebug("Executing count query", map[string]interface{}{
		"user_id": userID,
	})

	// 获取总数
	err := query.Count(&total).Error
	if err != nil {
		utils.LogError(err, "failed to count user teams")
		return nil, 0, err
	}

	utils.LogDebug("Count query completed", map[string]interface{}{
		"user_id": userID,
		"total":   total,
	})

	// 获取分页数据
	err = query.Offset(offset).Limit(limit).Find(&teams).Error
	if err != nil {
		utils.LogError(err, "failed to find user teams")
		return nil, 0, err
	}

	utils.LogDebug("GetUserTeams repository completed", map[string]interface{}{
		"user_id":    userID,
		"total":      total,
		"team_count": len(teams),
	})

	return teams, total, nil
}

// AddTeamMember 添加团队成员
func (r *teamRepository) AddTeamMember(member *TeamMember) error {
	return r.db.Create(member).Error
}

// GetTeamMembers 获取团队成员列表
func (r *teamRepository) GetTeamMembers(teamID int64, offset, limit int) ([]*TeamMember, int64, error) {
	var members []*TeamMember
	var total int64

	query := r.db.Where("team_id = ?", teamID)

	// 获取总数
	err := query.Model(&TeamMember{}).Count(&total).Error
	if err != nil {
		return nil, 0, err
	}

	// 获取分页数据
	err = query.Offset(offset).Limit(limit).Find(&members).Error
	if err != nil {
		return nil, 0, err
	}

	return members, total, nil
}

// GetTeamMember 获取团队成员
func (r *teamRepository) GetTeamMember(teamID, userID int64) (*TeamMember, error) {
	var member TeamMember
	err := r.db.Where("team_id = ? AND user_id = ?", teamID, userID).First(&member).Error
	if err != nil {
		return nil, err
	}
	return &member, nil
}

// UpdateTeamMember 更新团队成员
func (r *teamRepository) UpdateTeamMember(member *TeamMember) error {
	return r.db.Save(member).Error
}

// RemoveTeamMember 移除团队成员
func (r *teamRepository) RemoveTeamMember(teamID, userID int64) error {
	return r.db.Where("team_id = ? AND user_id = ?", teamID, userID).Delete(&TeamMember{}).Error
}

// GetUserRole 获取用户在团队中的角色
func (r *teamRepository) GetUserRole(teamID, userID int64) (string, error) {
	var member TeamMember
	err := r.db.Select("role").Where("team_id = ? AND user_id = ?", teamID, userID).First(&member).Error
	if err != nil {
		return "", err
	}
	return member.Role, nil
}

// CreateInvitation 创建邀请
func (r *teamRepository) CreateInvitation(invitation *TeamInvitation) error {
	return r.db.Create(invitation).Error
}

// GetInvitationByCode 根据邀请码获取邀请
func (r *teamRepository) GetInvitationByCode(code string) (*TeamInvitation, error) {
	var invitation TeamInvitation
	err := r.db.Where("invitation_code = ?", code).First(&invitation).Error
	if err != nil {
		return nil, err
	}
	return &invitation, nil
}

// GetUserInvitations 获取用户的邀请列表
func (r *teamRepository) GetUserInvitations(userID int64, offset, limit int) ([]*TeamInvitation, int64, error) {
	var invitations []*TeamInvitation
	var total int64

	// 根据邮箱或钱包地址查询邀请
	var user User
	err := r.db.Where("id = ?", userID).First(&user).Error
	if err != nil {
		return nil, 0, err
	}

	query := r.db.Where("status = 'pending' AND expires_at > ?", time.Now())
	if user.Email != nil {
		query = query.Where("invitee_email = ?", *user.Email)
	}
	if user.WalletAddress != nil {
		query = query.Or("invitee_wallet = ?", *user.WalletAddress)
	}

	// 获取总数
	err = query.Model(&TeamInvitation{}).Count(&total).Error
	if err != nil {
		return nil, 0, err
	}

	// 获取分页数据
	err = query.Offset(offset).Limit(limit).Find(&invitations).Error
	if err != nil {
		return nil, 0, err
	}

	return invitations, total, nil
}

// UpdateInvitation 更新邀请
func (r *teamRepository) UpdateInvitation(invitation *TeamInvitation) error {
	return r.db.Save(invitation).Error
}

// ExpireInvitations 过期邀请
func (r *teamRepository) ExpireInvitations() error {
	return r.db.Model(&TeamInvitation{}).
		Where("status = 'pending' AND expires_at < ?", time.Now()).
		Update("status", "expired").Error
}

// GetAllTeams 获取所有团队列表（管理员功能）
func (r *teamRepository) GetAllTeams(keyword string, offset int, limit int) ([]*Team, int64, error) {
	var teams []*Team
	var total int64

	query := r.db.Model(&Team{})

	// 关键字搜索（团队名称）
	if keyword != "" {
		query = query.Where("name LIKE ?", "%"+keyword+"%")
	}

	// 获取总数
	if err := query.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	// 获取分页数据
	err := query.Order("created_at DESC").
		Offset(offset).
		Limit(limit).
		Find(&teams).Error

	return teams, total, err
}

// GetTeamMemberCount 获取团队成员数量
func (r *teamRepository) GetTeamMemberCount(teamID int64) (int64, error) {
	var count int64
	err := r.db.Model(&TeamMember{}).
		Where("team_id = ?", teamID).
		Count(&count).Error
	return count, err
}
