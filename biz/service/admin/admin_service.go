package admin

import (
	"context"
	"errors"
	"fmt"
	"time"

	"orbia_api/biz/consts"
	"orbia_api/biz/dal/model"
	"orbia_api/biz/dal/mysql"
	adminmodel "orbia_api/biz/model/admin"
	"orbia_api/biz/model/common"
	"orbia_api/biz/utils"

	"github.com/cloudwego/hertz/pkg/common/hlog"
	"gorm.io/gorm"
)

// AdminService 管理员服务
type AdminService struct {
	userRepo     mysql.UserRepository
	kolRepo      mysql.KolRepository
	teamRepo     mysql.TeamRepository
	orderRepo    mysql.OrderRepository
	walletRepo   mysql.WalletRepository
	campaignRepo mysql.CampaignRepository
	txRepo       mysql.TransactionRepository
	db           *gorm.DB
}

// NewAdminService 创建管理员服务实例
func NewAdminService(
	userRepo mysql.UserRepository,
	kolRepo mysql.KolRepository,
	teamRepo mysql.TeamRepository,
	orderRepo mysql.OrderRepository,
	walletRepo mysql.WalletRepository,
	campaignRepo mysql.CampaignRepository,
	txRepo mysql.TransactionRepository,
	db *gorm.DB,
) *AdminService {
	return &AdminService{
		userRepo:     userRepo,
		kolRepo:      kolRepo,
		teamRepo:     teamRepo,
		orderRepo:    orderRepo,
		walletRepo:   walletRepo,
		campaignRepo: campaignRepo,
		txRepo:       txRepo,
		db:           db,
	}
}

// GetAllUsers 获取所有用户列表
func (s *AdminService) GetAllUsers(ctx context.Context, req *adminmodel.GetAllUsersReq) (*adminmodel.GetAllUsersResp, error) {
	// 标准化分页参数
	params := utils.NormalizePaginationValue(req.Page, req.PageSize)
	offset := int((params.Page - 1) * params.PageSize)

	keyword := ""
	if req.Keyword != nil {
		keyword = *req.Keyword
	}

	role := ""
	if req.Role != nil {
		role = *req.Role
	}

	status := ""
	if req.Status != nil {
		status = *req.Status
	}

	// 查询用户列表
	users, total, err := s.userRepo.GetAllUsers(keyword, role, status, offset, int(params.PageSize))
	if err != nil {
		hlog.Errorf("Failed to get all users: %v", err)
		return nil, err
	}

	// 构建响应
	userList := make([]*adminmodel.UserListItem, 0, len(users))
	for _, user := range users {
		item := &adminmodel.UserListItem{
			ID:        user.ID,
			Role:      user.Role,
			Status:    user.Status,
			CreatedAt: user.CreatedAt.Format("2006-01-02 15:04:05"),
			UpdatedAt: user.UpdatedAt.Format("2006-01-02 15:04:05"),
		}

		if user.WalletAddress != nil {
			item.WalletAddress = user.WalletAddress
		}
		if user.Email != nil {
			item.Email = user.Email
		}
		if user.Nickname != nil {
			item.Nickname = user.Nickname
		}
		if user.AvatarURL != nil {
			item.AvatarURL = user.AvatarURL
		}
		if user.KolID != nil {
			item.KolID = user.KolID
		}

		userList = append(userList, item)
	}

	pageInfo := utils.BuildPageResp(params, total)

	return &adminmodel.GetAllUsersResp{
		BaseResp: &common.BaseResp{
			Code:    0,
			Message: "success",
		},
		Users:    userList,
		PageInfo: pageInfo,
	}, nil
}

// SetUserStatus 设置用户状态
func (s *AdminService) SetUserStatus(ctx context.Context, req *adminmodel.SetUserStatusReq) (*adminmodel.SetUserStatusResp, error) {
	// 验证状态值
	if req.Status != "normal" && req.Status != "disabled" && req.Status != "deleted" {
		return nil, errors.New("invalid status value")
	}

	// 获取用户信息
	user, err := s.userRepo.GetUserByID(req.UserID)
	if err != nil {
		hlog.Errorf("Failed to get user: %v", err)
		return nil, errors.New("user not found")
	}

	// 检查是否是管理员
	if user.Role == string(consts.RoleAdmin) {
		return nil, errors.New("cannot modify admin user status")
	}

	// 更新用户状态
	if err := s.userRepo.UpdateUserStatus(req.UserID, req.Status); err != nil {
		hlog.Errorf("Failed to update user status: %v", err)
		return nil, err
	}

	return &adminmodel.SetUserStatusResp{
		BaseResp: &common.BaseResp{
			Code:    0,
			Message: "success",
		},
	}, nil
}

// GetAllKols 获取所有KOL列表
func (s *AdminService) GetAllKols(ctx context.Context, req *adminmodel.GetAllKolsReq) (*adminmodel.GetAllKolsResp, error) {
	// 标准化分页参数
	params := utils.NormalizePaginationValue(req.Page, req.PageSize)
	offset := int((params.Page - 1) * params.PageSize)

	keyword := ""
	if req.Keyword != nil {
		keyword = *req.Keyword
	}

	status := ""
	if req.Status != nil {
		status = *req.Status
	}

	country := ""
	if req.Country != nil {
		country = *req.Country
	}

	tag := ""
	if req.Tag != nil {
		tag = *req.Tag
	}

	// 查询KOL列表
	kols, total, err := s.kolRepo.GetAllKols(keyword, status, country, tag, offset, int(params.PageSize))
	if err != nil {
		hlog.Errorf("Failed to get all kols: %v", err)
		return nil, err
	}

	// 构建响应
	kolList := make([]*adminmodel.KolListItem, 0, len(kols))
	for _, kol := range kols {
		// 获取统计数据
		stats, _ := s.kolRepo.GetKolStats(kol.ID)

		item := &adminmodel.KolListItem{
			ID:          kol.ID,
			UserID:      kol.UserID,
			DisplayName: getStringValue(kol.DisplayName),
			AvatarURL:   getStringValue(kol.AvatarURL),
			Country:     getStringValue(kol.Country),
			Status:      kol.Status,
			CreatedAt:   kol.CreatedAt.Format("2006-01-02 15:04:05"),
			UpdatedAt:   kol.UpdatedAt.Format("2006-01-02 15:04:05"),
		}

		if stats != nil {
			totalFollowers := stats.TotalFollowers
			item.TotalFollowers = &totalFollowers
		}

		kolList = append(kolList, item)
	}

	pageInfo := utils.BuildPageResp(params, total)

	return &adminmodel.GetAllKolsResp{
		BaseResp: &common.BaseResp{
			Code:    0,
			Message: "success",
		},
		Kols:     kolList,
		PageInfo: pageInfo,
	}, nil
}

// AdminReviewKol 管理员审核KOL
func (s *AdminService) AdminReviewKol(ctx context.Context, req *adminmodel.AdminReviewKolReq) (*adminmodel.AdminReviewKolResp, error) {
	// 验证状态值
	if req.Status != "approved" && req.Status != "rejected" {
		return nil, errors.New("invalid status value")
	}

	// 获取KOL信息
	kol, err := s.kolRepo.GetKolByID(req.KolID)
	if err != nil {
		hlog.Errorf("Failed to get kol: %v", err)
		return nil, errors.New("kol not found")
	}

	// 更新KOL状态
	kol.Status = req.Status
	if req.Status == "approved" {
		now := time.Now()
		kol.ApprovedAt = &now
		kol.RejectReason = nil
	} else if req.Status == "rejected" && req.RejectReason != nil {
		kol.RejectReason = req.RejectReason
	}

	if err := s.kolRepo.UpdateKol(kol); err != nil {
		hlog.Errorf("Failed to update kol: %v", err)
		return nil, err
	}

	return &adminmodel.AdminReviewKolResp{
		BaseResp: &common.BaseResp{
			Code:    0,
			Message: "success",
		},
	}, nil
}

// GetAllTeams 获取所有团队列表
func (s *AdminService) GetAllTeams(ctx context.Context, req *adminmodel.GetAllTeamsReq) (*adminmodel.GetAllTeamsResp, error) {
	// 标准化分页参数
	params := utils.NormalizePaginationValue(req.Page, req.PageSize)
	offset := int((params.Page - 1) * params.PageSize)

	keyword := ""
	if req.Keyword != nil {
		keyword = *req.Keyword
	}

	// 查询团队列表
	teams, total, err := s.teamRepo.GetAllTeams(keyword, offset, int(params.PageSize))
	if err != nil {
		hlog.Errorf("Failed to get all teams: %v", err)
		return nil, err
	}

	// 构建响应
	teamList := make([]*adminmodel.TeamListItem, 0, len(teams))
	for _, team := range teams {
		// 获取成员数量
		memberCount, _ := s.teamRepo.GetTeamMemberCount(team.ID)

		// 获取创建者信息
		creator, _ := s.userRepo.GetUserByID(team.CreatorID)

		item := &adminmodel.TeamListItem{
			ID:          team.ID,
			Name:        team.Name,
			CreatorID:   team.CreatorID,
			MemberCount: memberCount,
			CreatedAt:   team.CreatedAt.Format("2006-01-02 15:04:05"),
		}

		if team.IconURL != nil {
			item.IconURL = team.IconURL
		}

		if creator != nil && creator.Nickname != nil {
			item.CreatorName = creator.Nickname
		}

		teamList = append(teamList, item)
	}

	pageInfo := utils.BuildPageResp(params, total)

	return &adminmodel.GetAllTeamsResp{
		BaseResp: &common.BaseResp{
			Code:    0,
			Message: "success",
		},
		Teams:    teamList,
		PageInfo: pageInfo,
	}, nil
}

// GetTeamMembers 获取特定团队的所有用户
func (s *AdminService) GetTeamMembers(ctx context.Context, req *adminmodel.GetTeamMembersReq) (*adminmodel.GetTeamMembersResp, error) {
	// 验证团队是否存在
	_, err := s.teamRepo.GetTeamByID(req.TeamID)
	if err != nil {
		hlog.Errorf("Failed to get team: %v", err)
		return nil, errors.New("team not found")
	}

	// 标准化分页参数
	params := utils.NormalizePaginationValue(req.Page, req.PageSize)
	offset := int((params.Page - 1) * params.PageSize)

	// 查询团队成员
	members, total, err := s.teamRepo.GetTeamMembers(req.TeamID, offset, int(params.PageSize))
	if err != nil {
		hlog.Errorf("Failed to get team members: %v", err)
		return nil, err
	}

	// 构建响应
	memberList := make([]*adminmodel.TeamMemberItem, 0, len(members))
	for _, member := range members {
		// 获取用户信息
		user, err := s.userRepo.GetUserByID(member.UserID)
		if err != nil {
			continue
		}

		item := &adminmodel.TeamMemberItem{
			UserID:   member.UserID,
			Role:     member.Role,
			JoinedAt: member.JoinedAt.Format("2006-01-02 15:04:05"),
		}

		if user.Nickname != nil {
			item.Nickname = user.Nickname
		}
		if user.Email != nil {
			item.Email = user.Email
		}
		if user.AvatarURL != nil {
			item.AvatarURL = user.AvatarURL
		}

		memberList = append(memberList, item)
	}

	pageInfo := utils.BuildPageResp(params, total)

	return &adminmodel.GetTeamMembersResp{
		BaseResp: &common.BaseResp{
			Code:    0,
			Message: "success",
		},
		Members:  memberList,
		PageInfo: pageInfo,
	}, nil
}

// GetAllOrders 获取所有订单列表
func (s *AdminService) GetAllOrders(ctx context.Context, req *adminmodel.GetAllOrdersReq) (*adminmodel.GetAllOrdersResp, error) {
	// 标准化分页参数
	params := utils.NormalizePaginationValue(req.Page, req.PageSize)
	offset := int((params.Page - 1) * params.PageSize)

	keyword := ""
	if req.Keyword != nil {
		keyword = *req.Keyword
	}

	status := ""
	if req.Status != nil {
		status = *req.Status
	}

	// 查询订单列表
	orders, total, err := s.orderRepo.GetAllOrders(keyword, status, offset, int(params.PageSize))
	if err != nil {
		hlog.Errorf("Failed to get all orders: %v", err)
		return nil, err
	}

	// 构建响应
	orderList := make([]*adminmodel.OrderListItem, 0, len(orders))
	for _, order := range orders {
		// 获取用户信息
		user, _ := s.userRepo.GetUserByID(order.UserID)

		item := &adminmodel.OrderListItem{
			OrderID:   order.OrderID,
			UserID:    order.UserID,
			KolID:     order.KolID,
			PlanTitle: order.PlanTitle,
			PlanPrice: order.PlanPrice,
			Status:    order.Status,
			CreatedAt: order.CreatedAt.Format("2006-01-02 15:04:05"),
		}

		if user != nil {
			if user.Nickname != nil {
				item.UserName = user.Nickname
			}
			if user.Email != nil {
				item.UserEmail = user.Email
			}
		}

		if order.KolDisplayName != nil {
			item.KolName = order.KolDisplayName
		}

		if order.CompletedAt != nil {
			completedAt := order.CompletedAt.Format("2006-01-02 15:04:05")
			item.CompletedAt = &completedAt
		}

		orderList = append(orderList, item)
	}

	pageInfo := utils.BuildPageResp(params, total)

	return &adminmodel.GetAllOrdersResp{
		BaseResp: &common.BaseResp{
			Code:    0,
			Message: "success",
		},
		Orders:   orderList,
		PageInfo: pageInfo,
	}, nil
}

// GetUserWallet 获取特定用户的钱包信息
func (s *AdminService) GetUserWallet(ctx context.Context, req *adminmodel.GetUserWalletReq) (*adminmodel.GetUserWalletResp, error) {
	// 获取用户信息
	user, err := s.userRepo.GetUserByID(req.UserID)
	if err != nil {
		hlog.Errorf("Failed to get user: %v", err)
		return nil, errors.New("user not found")
	}

	// 获取钱包信息
	wallet, err := s.walletRepo.GetWalletByUserID(req.UserID)
	if err != nil {
		hlog.Errorf("Failed to get wallet: %v", err)
		return nil, errors.New("wallet not found")
	}

	walletInfo := &adminmodel.UserWalletInfo{
		UserID:        user.ID,
		Balance:       wallet.Balance,
		FrozenBalance: wallet.FrozenBalance,
		TotalRecharge: wallet.TotalRecharge,
		TotalConsume:  wallet.TotalConsume,
		CreatedAt:     wallet.CreatedAt.Format("2006-01-02 15:04:05"),
		UpdatedAt:     wallet.UpdatedAt.Format("2006-01-02 15:04:05"),
	}

	if user.Nickname != nil {
		walletInfo.UserName = user.Nickname
	}
	if user.Email != nil {
		walletInfo.UserEmail = user.Email
	}

	return &adminmodel.GetUserWalletResp{
		BaseResp: &common.BaseResp{
			Code:    0,
			Message: "success",
		},
		Wallet: walletInfo,
	}, nil
}

// getStringValue 辅助函数：获取字符串指针的值
func getStringValue(s *string) string {
	if s == nil {
		return ""
	}
	return *s
}

// AddCampaignConsume 管理员给Campaign添加消费账单
func (s *AdminService) AddCampaignConsume(ctx context.Context, req *adminmodel.AddCampaignConsumeReq) (*adminmodel.AddCampaignConsumeResp, error) {
	// 1. 验证金额
	if req.Amount <= 0 {
		return nil, errors.New("amount must be greater than 0")
	}

	// 2. 获取 Campaign 信息
	campaign, err := s.campaignRepo.GetCampaignByCampaignID(req.CampaignID)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("campaign not found")
		}
		hlog.Errorf("Failed to get campaign: %v", err)
		return nil, fmt.Errorf("failed to get campaign: %v", err)
	}

	// 3. 获取 Campaign 创建用户的钱包
	wallet, err := s.walletRepo.GetWalletByUserID(campaign.UserID)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("user wallet not found")
		}
		hlog.Errorf("Failed to get wallet: %v", err)
		return nil, fmt.Errorf("failed to get wallet: %v", err)
	}

	// 4. 检查余额是否足够
	if wallet.Balance < req.Amount {
		return nil, fmt.Errorf("insufficient balance: current balance %.2f, required %.2f", wallet.Balance, req.Amount)
	}

	// 5. 开始事务
	var transactionID string
	err = s.db.Transaction(func(tx *gorm.DB) error {
		// 5.1 扣除用户钱包余额
		balanceDelta := -req.Amount // 负数表示扣款
		if err := s.walletRepo.UpdateBalance(tx, campaign.UserID, balanceDelta, 0); err != nil {
			hlog.Errorf("Failed to update wallet balance: %v", err)
			return fmt.Errorf("failed to update wallet balance: %v", err)
		}

		// 5.2 更新累计消费金额
		err := tx.Model(&model.OrbiaWallet{}).
			Where("user_id = ?", campaign.UserID).
			Update("total_consume", gorm.Expr("total_consume + ?", req.Amount)).Error
		if err != nil {
			hlog.Errorf("Failed to update total_consume: %v", err)
			return fmt.Errorf("failed to update total_consume: %v", err)
		}

		// 5.3 获取更新后的钱包信息
		updatedWallet, err := s.walletRepo.GetWalletByUserID(campaign.UserID)
		if err != nil {
			hlog.Errorf("Failed to get updated wallet: %v", err)
			return fmt.Errorf("failed to get updated wallet: %v", err)
		}

		// 5.4 创建交易记录
		transactionID = utils.GenerateTransactionID()
		now := time.Now()

		remark := ""
		if req.Remark != nil {
			remark = *req.Remark
		}

		transaction := &model.OrbiaTransaction{
			TransactionID:    transactionID,
			UserID:           campaign.UserID,
			Type:             "consume",
			Amount:           req.Amount,
			BalanceBefore:    wallet.Balance,
			BalanceAfter:     updatedWallet.Balance,
			Status:           "completed",
			RelatedOrderType: stringPtr("campaign"),
			RelatedOrderID:   &req.CampaignID,
			Remark:           stringPtr(remark),
			CompletedAt:      &now,
		}

		if err := s.txRepo.CreateTransaction(tx, transaction); err != nil {
			hlog.Errorf("Failed to create transaction: %v", err)
			return fmt.Errorf("failed to create transaction: %v", err)
		}

		return nil
	})

	if err != nil {
		hlog.Errorf("Transaction failed: %v", err)
		return nil, err
	}

	return &adminmodel.AddCampaignConsumeResp{
		BaseResp: &common.BaseResp{
			Code:    0,
			Message: "success",
		},
		TransactionID: &transactionID,
	}, nil
}

// stringPtr 辅助函数：将字符串转换为指针
func stringPtr(s string) *string {
	if s == "" {
		return nil
	}
	return &s
}
