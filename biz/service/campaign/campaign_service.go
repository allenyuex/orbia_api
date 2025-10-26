package campaign

import (
	"encoding/json"
	"errors"
	"fmt"
	"path/filepath"
	"strings"
	"time"

	"orbia_api/biz/dal/mysql"
	"orbia_api/biz/utils"

	"gorm.io/gorm"
)

// CampaignService Campaign服务接口
type CampaignService interface {
	// 普通用户接口
	CreateCampaign(userID int64, teamID int64, req *CreateCampaignRequest) (*mysql.Campaign, []*mysql.CampaignAttachment, error)
	UpdateCampaign(userID int64, campaignID string, req *UpdateCampaignRequest) (*mysql.Campaign, []*mysql.CampaignAttachment, error)
	UpdateCampaignStatus(userID int64, campaignID string, status string) error
	GetCampaign(userID int64, campaignID string) (*mysql.Campaign, []*mysql.CampaignAttachment, error)
	ListCampaigns(userID int64, teamID int64, keyword string, status string, promotionObjective string, page int, pageSize int) ([]*mysql.Campaign, int64, error)

	// 管理员接口
	AdminListCampaigns(keyword string, status string, promotionObjective string, userID *int64, teamID *int64, page int, pageSize int) ([]*mysql.Campaign, int64, error)
	AdminUpdateCampaignStatus(campaignID string, status string) error
}

// CreateCampaignRequest 创建Campaign请求
type CreateCampaignRequest struct {
	CampaignName       string
	PromotionObjective string
	OptimizationGoal   string
	Location           []int64
	Age                *int64
	Gender             *int64
	Languages          []int64
	SpendingPower      *int64
	OperatingSystem    *int64
	OSVersions         []int64
	DeviceModels       []int64
	ConnectionType     *int64
	DevicePriceType    int8
	DevicePriceMin     *float64
	DevicePriceMax     *float64
	PlannedStartTime   string
	PlannedEndTime     string
	TimeZone           *int64
	DaypartingType     int8
	DaypartingSchedule *string
	FrequencyCapType   int8
	FrequencyCapTimes  *int32
	FrequencyCapDays   *int32
	BudgetType         int8
	BudgetAmount       float64
	Website            *string
	IOSDownloadURL     *string
	AndroidDownloadURL *string
	AttachmentURLs     []string
}

// UpdateCampaignRequest 更新Campaign请求
type UpdateCampaignRequest struct {
	CampaignName       *string
	PromotionObjective *string
	OptimizationGoal   *string
	Location           []int64
	Age                *int64
	Gender             *int64
	Languages          []int64
	SpendingPower      *int64
	OperatingSystem    *int64
	OSVersions         []int64
	DeviceModels       []int64
	ConnectionType     *int64
	DevicePriceType    *int8
	DevicePriceMin     *float64
	DevicePriceMax     *float64
	PlannedStartTime   *string
	PlannedEndTime     *string
	TimeZone           *int64
	DaypartingType     *int8
	DaypartingSchedule *string
	FrequencyCapType   *int8
	FrequencyCapTimes  *int32
	FrequencyCapDays   *int32
	BudgetType         *int8
	BudgetAmount       *float64
	Website            *string
	IOSDownloadURL     *string
	AndroidDownloadURL *string
	AttachmentURLs     []string
}

// campaignService Campaign服务实现
type campaignService struct {
	campaignRepo mysql.CampaignRepository
	userRepo     mysql.UserRepository
	teamRepo     mysql.TeamRepository
}

// NewCampaignService 创建Campaign服务实例
func NewCampaignService(campaignRepo mysql.CampaignRepository, userRepo mysql.UserRepository, teamRepo mysql.TeamRepository) CampaignService {
	return &campaignService{
		campaignRepo: campaignRepo,
		userRepo:     userRepo,
		teamRepo:     teamRepo,
	}
}

// CreateCampaign 创建Campaign
func (s *campaignService) CreateCampaign(userID int64, teamID int64, req *CreateCampaignRequest) (*mysql.Campaign, []*mysql.CampaignAttachment, error) {
	// 验证用户存在
	_, err := s.userRepo.GetUserByID(userID)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil, errors.New("user not found")
		}
		return nil, nil, fmt.Errorf("failed to get user: %v", err)
	}

	// 验证团队存在
	_, err = s.teamRepo.GetTeamByID(teamID)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil, errors.New("team not found")
		}
		return nil, nil, fmt.Errorf("failed to get team: %v", err)
	}

	// 验证用户是团队成员
	isMember, err := s.teamRepo.IsTeamMember(teamID, userID)
	if err != nil {
		return nil, nil, fmt.Errorf("failed to check team membership: %v", err)
	}
	if !isMember {
		return nil, nil, errors.New("user is not a member of the team")
	}

	// 验证参数
	if err := validatePromotionObjective(req.PromotionObjective, req.OptimizationGoal); err != nil {
		return nil, nil, err
	}

	// 解析时间
	startTime, err := time.Parse(time.RFC3339, req.PlannedStartTime)
	if err != nil {
		return nil, nil, fmt.Errorf("invalid planned_start_time format: %v", err)
	}

	endTime, err := time.Parse(time.RFC3339, req.PlannedEndTime)
	if err != nil {
		return nil, nil, fmt.Errorf("invalid planned_end_time format: %v", err)
	}

	if endTime.Before(startTime) {
		return nil, nil, errors.New("planned_end_time must be after planned_start_time")
	}

	// 生成Campaign ID
	campaignID := utils.GenerateCampaignID()

	// 转换数组为JSON
	locationJSON := arrayToJSON(req.Location)
	languagesJSON := arrayToJSON(req.Languages)
	osVersionsJSON := arrayToJSON(req.OSVersions)
	deviceModelsJSON := arrayToJSON(req.DeviceModels)

	// 创建Campaign
	campaign := &mysql.Campaign{
		CampaignID:         campaignID,
		UserID:             userID,
		TeamID:             teamID,
		CampaignName:       req.CampaignName,
		PromotionObjective: req.PromotionObjective,
		OptimizationGoal:   req.OptimizationGoal,
		Location:           locationJSON,
		Age:                req.Age,
		Gender:             req.Gender,
		Languages:          languagesJSON,
		SpendingPower:      req.SpendingPower,
		OperatingSystem:    req.OperatingSystem,
		OSVersions:         osVersionsJSON,
		DeviceModels:       deviceModelsJSON,
		ConnectionType:     req.ConnectionType,
		DevicePriceType:    req.DevicePriceType,
		DevicePriceMin:     req.DevicePriceMin,
		DevicePriceMax:     req.DevicePriceMax,
		PlannedStartTime:   startTime,
		PlannedEndTime:     endTime,
		TimeZone:           req.TimeZone,
		DaypartingType:     req.DaypartingType,
		DaypartingSchedule: req.DaypartingSchedule,
		FrequencyCapType:   req.FrequencyCapType,
		FrequencyCapTimes:  req.FrequencyCapTimes,
		FrequencyCapDays:   req.FrequencyCapDays,
		BudgetType:         req.BudgetType,
		BudgetAmount:       req.BudgetAmount,
		Website:            req.Website,
		IOSDownloadURL:     req.IOSDownloadURL,
		AndroidDownloadURL: req.AndroidDownloadURL,
		Status:             "pending",
	}

	if err := s.campaignRepo.CreateCampaign(campaign); err != nil {
		return nil, nil, fmt.Errorf("failed to create campaign: %v", err)
	}

	// 创建附件
	var attachments []*mysql.CampaignAttachment
	for _, url := range req.AttachmentURLs {
		fileName := filepath.Base(url)
		fileExt := strings.ToLower(filepath.Ext(fileName))
		fileType := getFileType(fileExt)

		attachment := &mysql.CampaignAttachment{
			CampaignID: campaign.ID,
			FileURL:    url,
			FileName:   fileName,
			FileType:   fileType,
		}

		if err := s.campaignRepo.CreateAttachment(attachment); err != nil {
			return nil, nil, fmt.Errorf("failed to create attachment: %v", err)
		}

		attachments = append(attachments, attachment)
	}

	return campaign, attachments, nil
}

// UpdateCampaign 更新Campaign
func (s *campaignService) UpdateCampaign(userID int64, campaignID string, req *UpdateCampaignRequest) (*mysql.Campaign, []*mysql.CampaignAttachment, error) {
	// 获取Campaign
	campaign, err := s.campaignRepo.GetCampaignByCampaignID(campaignID)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil, errors.New("campaign not found")
		}
		return nil, nil, fmt.Errorf("failed to get campaign: %v", err)
	}

	// 验证权限
	if campaign.UserID != userID {
		return nil, nil, errors.New("permission denied")
	}

	// 只能在pending状态下修改
	if campaign.Status != "pending" {
		return nil, nil, errors.New("only pending campaigns can be updated")
	}

	// 更新字段
	if req.CampaignName != nil {
		campaign.CampaignName = *req.CampaignName
	}

	if req.PromotionObjective != nil && req.OptimizationGoal != nil {
		if err := validatePromotionObjective(*req.PromotionObjective, *req.OptimizationGoal); err != nil {
			return nil, nil, err
		}
		campaign.PromotionObjective = *req.PromotionObjective
		campaign.OptimizationGoal = *req.OptimizationGoal
	}

	if len(req.Location) > 0 {
		campaign.Location = arrayToJSON(req.Location)
	}

	if req.Age != nil {
		campaign.Age = req.Age
	}

	if req.Gender != nil {
		campaign.Gender = req.Gender
	}

	if len(req.Languages) > 0 {
		campaign.Languages = arrayToJSON(req.Languages)
	}

	if req.SpendingPower != nil {
		campaign.SpendingPower = req.SpendingPower
	}

	if req.OperatingSystem != nil {
		campaign.OperatingSystem = req.OperatingSystem
	}

	if len(req.OSVersions) > 0 {
		campaign.OSVersions = arrayToJSON(req.OSVersions)
	}

	if len(req.DeviceModels) > 0 {
		campaign.DeviceModels = arrayToJSON(req.DeviceModels)
	}

	if req.ConnectionType != nil {
		campaign.ConnectionType = req.ConnectionType
	}

	if req.DevicePriceType != nil {
		campaign.DevicePriceType = *req.DevicePriceType
	}

	if req.DevicePriceMin != nil {
		campaign.DevicePriceMin = req.DevicePriceMin
	}

	if req.DevicePriceMax != nil {
		campaign.DevicePriceMax = req.DevicePriceMax
	}

	if req.PlannedStartTime != nil {
		startTime, err := time.Parse(time.RFC3339, *req.PlannedStartTime)
		if err != nil {
			return nil, nil, fmt.Errorf("invalid planned_start_time format: %v", err)
		}
		campaign.PlannedStartTime = startTime
	}

	if req.PlannedEndTime != nil {
		endTime, err := time.Parse(time.RFC3339, *req.PlannedEndTime)
		if err != nil {
			return nil, nil, fmt.Errorf("invalid planned_end_time format: %v", err)
		}
		campaign.PlannedEndTime = endTime
	}

	if campaign.PlannedEndTime.Before(campaign.PlannedStartTime) {
		return nil, nil, errors.New("planned_end_time must be after planned_start_time")
	}

	if req.TimeZone != nil {
		campaign.TimeZone = req.TimeZone
	}

	if req.DaypartingType != nil {
		campaign.DaypartingType = *req.DaypartingType
	}

	if req.DaypartingSchedule != nil {
		campaign.DaypartingSchedule = req.DaypartingSchedule
	}

	if req.FrequencyCapType != nil {
		campaign.FrequencyCapType = *req.FrequencyCapType
	}

	if req.FrequencyCapTimes != nil {
		campaign.FrequencyCapTimes = req.FrequencyCapTimes
	}

	if req.FrequencyCapDays != nil {
		campaign.FrequencyCapDays = req.FrequencyCapDays
	}

	if req.BudgetType != nil {
		campaign.BudgetType = *req.BudgetType
	}

	if req.BudgetAmount != nil {
		campaign.BudgetAmount = *req.BudgetAmount
	}

	if req.Website != nil {
		campaign.Website = req.Website
	}

	if req.IOSDownloadURL != nil {
		campaign.IOSDownloadURL = req.IOSDownloadURL
	}

	if req.AndroidDownloadURL != nil {
		campaign.AndroidDownloadURL = req.AndroidDownloadURL
	}

	// 更新附件
	if len(req.AttachmentURLs) > 0 {
		// 删除旧附件
		if err := s.campaignRepo.DeleteAttachmentsByCampaignID(campaign.ID); err != nil {
			return nil, nil, fmt.Errorf("failed to delete old attachments: %v", err)
		}

		// 创建新附件
		for _, url := range req.AttachmentURLs {
			fileName := filepath.Base(url)
			fileExt := strings.ToLower(filepath.Ext(fileName))
			fileType := getFileType(fileExt)

			attachment := &mysql.CampaignAttachment{
				CampaignID: campaign.ID,
				FileURL:    url,
				FileName:   fileName,
				FileType:   fileType,
			}

			if err := s.campaignRepo.CreateAttachment(attachment); err != nil {
				return nil, nil, fmt.Errorf("failed to create attachment: %v", err)
			}
		}
	}

	// 保存Campaign
	if err := s.campaignRepo.UpdateCampaign(campaign); err != nil {
		return nil, nil, fmt.Errorf("failed to update campaign: %v", err)
	}

	// 获取附件列表
	attachments, err := s.campaignRepo.GetAttachmentsByCampaignID(campaign.ID)
	if err != nil {
		return nil, nil, fmt.Errorf("failed to get attachments: %v", err)
	}

	return campaign, attachments, nil
}

// UpdateCampaignStatus 更新Campaign状态
func (s *campaignService) UpdateCampaignStatus(userID int64, campaignID string, status string) error {
	// 获取Campaign
	campaign, err := s.campaignRepo.GetCampaignByCampaignID(campaignID)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return errors.New("campaign not found")
		}
		return fmt.Errorf("failed to get campaign: %v", err)
	}

	// 验证权限
	if campaign.UserID != userID {
		return errors.New("permission denied")
	}

	// 验证状态转换
	if err := validateStatusTransition(campaign.Status, status, false); err != nil {
		return err
	}

	campaign.Status = status

	if err := s.campaignRepo.UpdateCampaign(campaign); err != nil {
		return fmt.Errorf("failed to update campaign status: %v", err)
	}

	return nil
}

// GetCampaign 获取Campaign详情
func (s *campaignService) GetCampaign(userID int64, campaignID string) (*mysql.Campaign, []*mysql.CampaignAttachment, error) {
	// 获取Campaign
	campaign, err := s.campaignRepo.GetCampaignByCampaignID(campaignID)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil, errors.New("campaign not found")
		}
		return nil, nil, fmt.Errorf("failed to get campaign: %v", err)
	}

	// 验证权限 - 只能查看自己创建的Campaign
	if campaign.UserID != userID {
		return nil, nil, errors.New("permission denied")
	}

	// 获取附件列表
	attachments, err := s.campaignRepo.GetAttachmentsByCampaignID(campaign.ID)
	if err != nil {
		return nil, nil, fmt.Errorf("failed to get attachments: %v", err)
	}

	return campaign, attachments, nil
}

// ListCampaigns 获取Campaign列表
func (s *campaignService) ListCampaigns(userID int64, teamID int64, keyword string, status string, promotionObjective string, page int, pageSize int) ([]*mysql.Campaign, int64, error) {
	// 验证用户是团队成员
	isMember, err := s.teamRepo.IsTeamMember(teamID, userID)
	if err != nil {
		return nil, 0, fmt.Errorf("failed to check team membership: %v", err)
	}
	if !isMember {
		return nil, 0, errors.New("user is not a member of the team")
	}

	offset := (page - 1) * pageSize

	campaigns, total, err := s.campaignRepo.GetCampaignsByTeamID(teamID, keyword, status, promotionObjective, offset, pageSize)
	if err != nil {
		return nil, 0, fmt.Errorf("failed to get campaigns: %v", err)
	}

	return campaigns, total, nil
}

// AdminListCampaigns 管理员获取所有Campaign列表
func (s *campaignService) AdminListCampaigns(keyword string, status string, promotionObjective string, userID *int64, teamID *int64, page int, pageSize int) ([]*mysql.Campaign, int64, error) {
	offset := (page - 1) * pageSize

	campaigns, total, err := s.campaignRepo.GetAllCampaigns(keyword, status, promotionObjective, userID, teamID, offset, pageSize)
	if err != nil {
		return nil, 0, fmt.Errorf("failed to get campaigns: %v", err)
	}

	return campaigns, total, nil
}

// AdminUpdateCampaignStatus 管理员更新Campaign状态
func (s *campaignService) AdminUpdateCampaignStatus(campaignID string, status string) error {
	// 获取Campaign
	campaign, err := s.campaignRepo.GetCampaignByCampaignID(campaignID)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return errors.New("campaign not found")
		}
		return fmt.Errorf("failed to get campaign: %v", err)
	}

	// 验证状态转换
	if err := validateStatusTransition(campaign.Status, status, true); err != nil {
		return err
	}

	campaign.Status = status

	if err := s.campaignRepo.UpdateCampaign(campaign); err != nil {
		return fmt.Errorf("failed to update campaign status: %v", err)
	}

	return nil
}

// Helper functions

// validatePromotionObjective 验证推广目标和优化目标的组合
func validatePromotionObjective(objective string, goal string) error {
	validGoals := map[string][]string{
		"awareness":     {"reach"},
		"consideration": {"website", "app"},
		"conversion":    {"app_promotion", "lead_generation"},
	}

	goals, ok := validGoals[objective]
	if !ok {
		return errors.New("invalid promotion_objective")
	}

	for _, validGoal := range goals {
		if goal == validGoal {
			return nil
		}
	}

	return fmt.Errorf("invalid optimization_goal for promotion_objective %s", objective)
}

// validateStatusTransition 验证状态转换
func validateStatusTransition(currentStatus string, newStatus string, isAdmin bool) error {
	if currentStatus == newStatus {
		return errors.New("status is already " + newStatus)
	}

	// 普通用户只能在active和paused之间切换
	if !isAdmin {
		if currentStatus == "active" && newStatus == "paused" {
			return nil
		}
		if currentStatus == "paused" && newStatus == "active" {
			return nil
		}
		return errors.New("invalid status transition")
	}

	// 管理员可以设置为active、paused、ended
	if newStatus == "active" || newStatus == "paused" || newStatus == "ended" {
		return nil
	}

	return errors.New("invalid status")
}

// arrayToJSON 将数组转换为JSON字符串
func arrayToJSON(arr []int64) *string {
	if len(arr) == 0 {
		return nil
	}
	data, _ := json.Marshal(arr)
	str := string(data)
	return &str
}

// getFileType 根据文件扩展名获取MIME类型
func getFileType(ext string) string {
	mimeTypes := map[string]string{
		".pdf":  "application/pdf",
		".jpg":  "image/jpeg",
		".jpeg": "image/jpeg",
		".png":  "image/png",
		".gif":  "image/gif",
		".mp4":  "video/mp4",
		".mov":  "video/quicktime",
		".avi":  "video/x-msvideo",
		".zip":  "application/zip",
		".doc":  "application/msword",
		".docx": "application/vnd.openxmlformats-officedocument.wordprocessingml.document",
		".xls":  "application/vnd.ms-excel",
		".xlsx": "application/vnd.openxmlformats-officedocument.spreadsheetml.sheet",
	}

	if mimeType, ok := mimeTypes[ext]; ok {
		return mimeType
	}

	return "application/octet-stream"
}
