package team

import (
	"errors"
	"fmt"
	"strconv"
	"time"

	"gorm.io/gorm"
	"orbia_api/biz/dal/mysql"
	"orbia_api/biz/model/common"
	"orbia_api/biz/model/team"
	"orbia_api/biz/model/upload"
	"orbia_api/biz/utils"
)

// TeamService 团队服务接口
type TeamService interface {
	// 团队管理
	CreateTeam(userID int64, req *team.CreateTeamReq) (*team.CreateTeamResp, error)
	GetTeam(userID int64, req *team.GetTeamReq) (*team.GetTeamResp, error)
	UpdateTeam(userID int64, req *team.UpdateTeamReq) (*team.UpdateTeamResp, error)
	GetUserTeams(userID int64, req *team.GetUserTeamsReq) (*team.GetUserTeamsResp, error)

	// 成员管理
	InviteUser(userID int64, req *team.InviteUserReq) (*team.InviteUserResp, error)
	AcceptInvitation(userID int64, req *team.AcceptInvitationReq) (*team.AcceptInvitationResp, error)
	RejectInvitation(userID int64, req *team.RejectInvitationReq) (*team.RejectInvitationResp, error)
	GetTeamMembers(userID int64, req *team.GetTeamMembersReq) (*team.GetTeamMembersResp, error)
}

// teamService 团队服务实现
type teamService struct {
	teamRepo mysql.TeamRepository
	userRepo mysql.UserRepository
}

// NewTeamService 创建团队服务实例
func NewTeamService(teamRepo mysql.TeamRepository, userRepo mysql.UserRepository) TeamService {
	return &teamService{
		teamRepo: teamRepo,
		userRepo: userRepo,
	}
}

// CreateTeam 创建团队
func (s *teamService) CreateTeam(userID int64, req *team.CreateTeamReq) (*team.CreateTeamResp, error) {
	// 验证团队名称长度
	if len(req.Name) > 20 {
		return nil, errors.New("team name cannot exceed 20 characters")
	}

	// 创建团队
	t := &mysql.Team{
		Name:      req.Name,
		IconURL:   req.IconURL,
		CreatorID: userID,
	}

	err := s.teamRepo.CreateTeam(t)
	if err != nil {
		return nil, fmt.Errorf("failed to create team: %v", err)
	}

	// 将创建者添加为团队成员（角色为 creator）
	member := &mysql.TeamMember{
		TeamID: t.ID,
		UserID: userID,
		Role:   "creator",
	}

	err = s.teamRepo.AddTeamMember(member)
	if err != nil {
		return nil, fmt.Errorf("failed to add creator as member: %v", err)
	}

	return &team.CreateTeamResp{
		BaseResp: &common.BaseResp{Code: 0, Message: "success"},
		Team:     s.convertTeamToModel(t),
	}, nil
}

// GetTeam 获取团队信息
func (s *teamService) GetTeam(userID int64, req *team.GetTeamReq) (*team.GetTeamResp, error) {
	teamID, err := strconv.ParseInt(req.TeamID, 10, 64)
	if err != nil {
		return nil, fmt.Errorf("invalid team_id: %v", err)
	}

	t, err := s.teamRepo.GetTeamByID(teamID)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("team not found")
		}
		return nil, fmt.Errorf("failed to get team: %v", err)
	}

	// 检查用户是否是团队成员
	_, err = s.teamRepo.GetTeamMember(teamID, userID)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("access denied: user is not a team member")
		}
		return nil, fmt.Errorf("failed to check membership: %v", err)
	}

	return &team.GetTeamResp{
		BaseResp: &common.BaseResp{Code: 0, Message: "success"},
		Team:     s.convertTeamToModel(t),
	}, nil
}

// UpdateTeam 更新团队信息
func (s *teamService) UpdateTeam(userID int64, req *team.UpdateTeamReq) (*team.UpdateTeamResp, error) {
	// 验证团队名称长度
	if req.Name != nil && len(*req.Name) > 20 {
		return nil, errors.New("team name cannot exceed 20 characters")
	}

	teamID, err := strconv.ParseInt(req.TeamID, 10, 64)
	if err != nil {
		return nil, fmt.Errorf("invalid team_id: %v", err)
	}

	// 检查用户权限（只有 creator 和 owner 可以修改团队信息）
	hasPermission, err := s.checkTeamEditPermission(teamID, userID)
	if err != nil {
		return nil, err
	}
	if !hasPermission {
		return nil, errors.New("access denied: insufficient permissions")
	}

	// 获取团队
	t, err := s.teamRepo.GetTeamByID(teamID)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("team not found")
		}
		return nil, fmt.Errorf("failed to get team: %v", err)
	}

	// 更新团队信息
	if req.Name != nil {
		t.Name = *req.Name
	}
	if req.IconURL != nil {
		// 验证团队图标URL（如果提供）
		if *req.IconURL != "" {
			isValid, errorMessage := utils.ValidateImageURL(*req.IconURL)
			if !isValid {
				return nil, fmt.Errorf("invalid icon URL: %s", errorMessage)
			}

			// 检查图片是否存在
			if !utils.CheckImageExists(*req.IconURL) {
				return nil, errors.New("icon image does not exist or is not accessible")
			}

			// 验证图片类型是否为团队图标
			imagePath := (*req.IconURL)[len(utils.GeneratePublicURL("")):]
			expectedType := utils.GetImageTypeFromPath(imagePath)
			if expectedType != utils.ImageType(upload.ImageType_TEAM_ICON) {
				return nil, errors.New("image type must be team icon")
			}
		}
		t.IconURL = req.IconURL
	}

	err = s.teamRepo.UpdateTeam(t)
	if err != nil {
		return nil, fmt.Errorf("failed to update team: %v", err)
	}

	return &team.UpdateTeamResp{
		BaseResp: &common.BaseResp{Code: 0, Message: "success"},
		Team:     s.convertTeamToModel(t),
	}, nil
}

// GetUserTeams 获取用户团队列表
func (s *teamService) GetUserTeams(userID int64, req *team.GetUserTeamsReq) (*team.GetUserTeamsResp, error) {
	utils.LogDebug("GetUserTeams service called", map[string]interface{}{
		"user_id": userID,
	})
	
	// 检查 teamRepo 是否为 nil
	if s.teamRepo == nil {
		utils.LogError(nil, "teamRepo is nil in GetUserTeams")
		return nil, errors.New("team repository is not initialized")
	}
	
	teams, _, err := s.teamRepo.GetUserTeams(userID, 0, 100) // 暂时设置固定分页
	if err != nil {
		utils.LogError(err, "failed to get user teams from repository")
		return nil, fmt.Errorf("failed to get user teams: %v", err)
	}

	var teamList []*team.Team
	for _, t := range teams {
		if t == nil {
			utils.LogError(nil, "found nil team in teams list")
			continue
		}
		teamList = append(teamList, s.convertTeamToModel(t))
	}

	utils.LogDebug("GetUserTeams service completed", map[string]interface{}{
		"user_id":    userID,
		"team_count": len(teamList),
	})

	return &team.GetUserTeamsResp{
		BaseResp: &common.BaseResp{Code: 0, Message: "success"},
		Teams:    teamList,
	}, nil
}

// InviteUser 邀请用户加入团队
func (s *teamService) InviteUser(userID int64, req *team.InviteUserReq) (*team.InviteUserResp, error) {
	teamID, err := strconv.ParseInt(req.TeamID, 10, 64)
	if err != nil {
		return nil, fmt.Errorf("invalid team_id: %v", err)
	}

	// 检查邀请者权限（creator 和 owner 可以邀请用户）
	hasPermission, err := s.checkTeamEditPermission(teamID, userID)
	if err != nil {
		return nil, err
	}
	if !hasPermission {
		return nil, errors.New("access denied: insufficient permissions to invite users")
	}

	// 检查被邀请用户是否已经是团队成员
	var inviteeID int64
	if req.Email != nil {
		user, err := s.userRepo.GetUserByEmail(*req.Email)
		if err != nil {
			if errors.Is(err, gorm.ErrRecordNotFound) {
				return nil, errors.New("user not found with this email")
			}
			return nil, fmt.Errorf("failed to get user by email: %v", err)
		}
		inviteeID = user.ID
	} else if req.WalletAddress != nil {
		user, err := s.userRepo.GetUserByWalletAddress(*req.WalletAddress)
		if err != nil {
			if errors.Is(err, gorm.ErrRecordNotFound) {
				return nil, errors.New("user not found with this wallet address")
			}
			return nil, fmt.Errorf("failed to get user by wallet address: %v", err)
		}
		inviteeID = user.ID
	} else {
		return nil, errors.New("either email or wallet address must be provided")
	}

	// 检查用户是否已经是团队成员
	_, err = s.teamRepo.GetTeamMember(teamID, inviteeID)
	if err == nil {
		return nil, errors.New("user is already a team member")
	}

	if !errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, fmt.Errorf("failed to check membership: %v", err)
	}

	// 创建邀请
	invitation := &mysql.TeamInvitation{
		TeamID:         teamID,
		InviterID:      userID,
		InviteeEmail:   req.Email,
		InviteeWallet:  req.WalletAddress,
		Role:           s.convertRoleToString(req.Role),
		Status:         "pending",
		InvitationCode: s.generateInvitationCode(),
		ExpiresAt:      time.Now().Add(7 * 24 * time.Hour), // 7天后过期
	}

	err = s.teamRepo.CreateInvitation(invitation)
	if err != nil {
		return nil, fmt.Errorf("failed to create invitation: %v", err)
	}

	return &team.InviteUserResp{
		BaseResp:   &common.BaseResp{Code: 0, Message: "success"},
		Invitation: s.convertInvitationToModel(invitation),
	}, nil
}

// AcceptInvitation 接受邀请
func (s *teamService) AcceptInvitation(userID int64, req *team.AcceptInvitationReq) (*team.AcceptInvitationResp, error) {
	// 获取邀请信息
	invitation, err := s.teamRepo.GetInvitationByCode(req.InvitationCode)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("invitation not found")
		}
		return nil, fmt.Errorf("failed to get invitation: %v", err)
	}

	// 检查邀请状态
	if invitation.Status != "pending" {
		return nil, errors.New("invitation is not pending")
	}

	// 检查邀请是否过期
	if time.Now().After(invitation.ExpiresAt) {
		return nil, errors.New("invitation has expired")
	}

	// 验证用户身份（通过邮箱或钱包地址）
	user, err := s.userRepo.GetUserByID(userID)
	if err != nil {
		return nil, fmt.Errorf("failed to get user: %v", err)
	}

	if invitation.InviteeEmail != nil && (user.Email == nil || *user.Email != *invitation.InviteeEmail) {
		return nil, errors.New("invitation email does not match user email")
	}

	if invitation.InviteeWallet != nil && (user.WalletAddress == nil || *user.WalletAddress != *invitation.InviteeWallet) {
		return nil, errors.New("invitation wallet address does not match user wallet address")
	}

	// 添加用户为团队成员
	member := &mysql.TeamMember{
		TeamID: invitation.TeamID,
		UserID: userID,
		Role:   invitation.Role,
	}

	err = s.teamRepo.AddTeamMember(member)
	if err != nil {
		return nil, fmt.Errorf("failed to add team member: %v", err)
	}

	// 更新邀请状态
	invitation.Status = "accepted"
	err = s.teamRepo.UpdateInvitation(invitation)
	if err != nil {
		return nil, fmt.Errorf("failed to update invitation status: %v", err)
	}

	return &team.AcceptInvitationResp{
		BaseResp: &common.BaseResp{Code: 0, Message: "success"},
	}, nil
}

// RejectInvitation 拒绝邀请
func (s *teamService) RejectInvitation(userID int64, req *team.RejectInvitationReq) (*team.RejectInvitationResp, error) {
	// 获取邀请信息
	invitation, err := s.teamRepo.GetInvitationByCode(req.InvitationCode)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("invitation not found")
		}
		return nil, fmt.Errorf("failed to get invitation: %v", err)
	}

	// 检查邀请状态
	if invitation.Status != "pending" {
		return nil, errors.New("invitation is not pending")
	}

	// 更新邀请状态
	invitation.Status = "rejected"
	err = s.teamRepo.UpdateInvitation(invitation)
	if err != nil {
		return nil, fmt.Errorf("failed to update invitation status: %v", err)
	}

	return &team.RejectInvitationResp{
		BaseResp: &common.BaseResp{Code: 0, Message: "success"},
	}, nil
}

// GetTeamMembers 获取团队成员列表
func (s *teamService) GetTeamMembers(userID int64, req *team.GetTeamMembersReq) (*team.GetTeamMembersResp, error) {
	teamID, err := strconv.ParseInt(req.TeamID, 10, 64)
	if err != nil {
		return nil, fmt.Errorf("invalid team_id: %v", err)
	}

	// 检查用户是否是团队成员
	_, err = s.teamRepo.GetTeamMember(teamID, userID)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("access denied: user is not a team member")
		}
		return nil, fmt.Errorf("failed to check membership: %v", err)
	}

	// 获取团队成员列表
	members, _, err := s.teamRepo.GetTeamMembers(teamID, 0, 100) // 暂时设置固定分页
	if err != nil {
		return nil, fmt.Errorf("failed to get team members: %v", err)
	}

	var memberList []*team.TeamMember
	for _, member := range members {
		memberList = append(memberList, s.convertMemberToModel(member))
	}

	return &team.GetTeamMembersResp{
		BaseResp: &common.BaseResp{Code: 0, Message: "success"},
		Members:  memberList,
	}, nil
}

// 辅助方法

// checkTeamEditPermission 检查团队编辑权限
func (s *teamService) checkTeamEditPermission(teamID, userID int64) (bool, error) {
	member, err := s.teamRepo.GetTeamMember(teamID, userID)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return false, nil
		}
		return false, fmt.Errorf("failed to get member info: %v", err)
	}

	return member.Role == "creator" || member.Role == "owner", nil
}

// generateInvitationCode 生成邀请码
func (s *teamService) generateInvitationCode() string {
	return fmt.Sprintf("inv_%d", time.Now().UnixNano())
}

// convertRoleToString 将角色枚举转换为字符串
func (s *teamService) convertRoleToString(role team.TeamRole) string {
	switch role {
	case team.TeamRole_CREATOR:
		return "creator"
	case team.TeamRole_OWNER:
		return "owner"
	case team.TeamRole_MEMBER:
		return "member"
	default:
		return "member"
	}
}

// convertTeamToModel 将数据库团队模型转换为API模型
func (s *teamService) convertTeamToModel(t *mysql.Team) *team.Team {
	return &team.Team{
		ID:        t.ID,
		Name:      t.Name,
		IconURL:   t.IconURL,
		CreatorID: t.CreatorID,
		CreatedAt: t.CreatedAt.Format(time.RFC3339),
		UpdatedAt: t.UpdatedAt.Format(time.RFC3339),
	}
}

// convertMemberToModel 将数据库成员模型转换为API模型
func (s *teamService) convertMemberToModel(member *mysql.TeamMember) *team.TeamMember {
	// 获取用户信息
	user, err := s.userRepo.GetUserByID(member.UserID)
	if err != nil {
		// 如果获取用户信息失败，返回基本信息
		return &team.TeamMember{
			ID:       member.ID,
			TeamID:   member.TeamID,
			UserID:   member.UserID,
			Role:     s.convertStringToRole(member.Role),
			JoinedAt: member.JoinedAt.Format(time.RFC3339),
		}
	}

	return &team.TeamMember{
		ID:                member.ID,
		TeamID:            member.TeamID,
		UserID:            member.UserID,
		Role:              s.convertStringToRole(member.Role),
		JoinedAt:          member.JoinedAt.Format(time.RFC3339),
		UserNickname:      user.Nickname,
		UserAvatarURL:     user.AvatarURL,
		UserEmail:         user.Email,
		UserWalletAddress: user.WalletAddress,
	}
}

// convertInvitationToModel 将数据库邀请模型转换为API模型
func (s *teamService) convertInvitationToModel(invitation *mysql.TeamInvitation) *team.TeamInvitation {
	return &team.TeamInvitation{
		ID:             invitation.ID,
		TeamID:         invitation.TeamID,
		InviterID:      invitation.InviterID,
		InviteeEmail:   invitation.InviteeEmail,
		InviteeWallet:  invitation.InviteeWallet,
		Role:           s.convertStringToRole(invitation.Role),
		Status:         s.convertStringToInvitationStatus(invitation.Status),
		InvitationCode: invitation.InvitationCode,
		ExpiresAt:      invitation.ExpiresAt.Format(time.RFC3339),
		CreatedAt:      invitation.CreatedAt.Format(time.RFC3339),
	}
}

// convertStringToRole 将字符串转换为角色枚举
func (s *teamService) convertStringToRole(role string) team.TeamRole {
	switch role {
	case "creator":
		return team.TeamRole_CREATOR
	case "owner":
		return team.TeamRole_OWNER
	case "member":
		return team.TeamRole_MEMBER
	default:
		return team.TeamRole_MEMBER
	}
}

// convertStringToInvitationStatus 将字符串转换为邀请状态枚举
func (s *teamService) convertStringToInvitationStatus(status string) team.InvitationStatus {
	switch status {
	case "pending":
		return team.InvitationStatus_PENDING
	case "accepted":
		return team.InvitationStatus_ACCEPTED
	case "rejected":
		return team.InvitationStatus_REJECTED
	case "expired":
		return team.InvitationStatus_EXPIRED
	default:
		return team.InvitationStatus_PENDING
	}
}
