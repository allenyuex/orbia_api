package campaign

import (
	"context"
	"encoding/json"
	"math"

	"github.com/cloudwego/hertz/pkg/app"
	"github.com/cloudwego/hertz/pkg/protocol/consts"

	"orbia_api/biz/dal/mysql"
	campaignModel "orbia_api/biz/model/campaign"
	commonModel "orbia_api/biz/model/common"
	"orbia_api/biz/mw"
	campaignService "orbia_api/biz/service/campaign"
	"orbia_api/biz/utils"
)

var (
	svc campaignService.CampaignService
)

// InitCampaignService 初始化Campaign服务
func InitCampaignService() {
	campaignRepo := mysql.NewCampaignRepository(mysql.DB)
	userRepo := mysql.NewUserRepository(mysql.DB)
	teamRepo := mysql.NewTeamRepository(mysql.DB)
	svc = campaignService.NewCampaignService(campaignRepo, userRepo, teamRepo)
}

// CreateCampaign 创建Campaign
// @router /campaign/create [POST]
func CreateCampaign(ctx context.Context, c *app.RequestContext) {
	var req campaignModel.CreateCampaignReq
	if err := c.BindAndValidate(&req); err != nil {
		utils.ParamError(c, "Invalid request parameters: "+err.Error())
		return
	}

	// 从context获取用户ID和用户信息
	userID, exists := mw.GetAuthUserID(c)
	if !exists {
		utils.Error(c, 401, "User not authenticated")
		return
	}

	user, exists := mw.GetAuthUser(c)
	if !exists || user.CurrentTeamID == nil {
		utils.Error(c, 400, "User has no team")
		return
	}

	teamID := *user.CurrentTeamID

	// 构建service请求
	serviceReq := &campaignService.CreateCampaignRequest{
		CampaignName:       req.CampaignName,
		PromotionObjective: req.PromotionObjective,
		OptimizationGoal:   req.OptimizationGoal,
		Location:           req.Location,
		Age:                req.Age,
		Gender:             req.Gender,
		Languages:          req.Languages,
		SpendingPower:      req.SpendingPower,
		OperatingSystem:    req.OperatingSystem,
		OSVersions:         req.OsVersions,
		DeviceModels:       req.DeviceModels,
		ConnectionTypes:    req.ConnectionTypes,
		DevicePriceType:    int8(req.DevicePriceType),
		DevicePriceMin:     req.DevicePriceMin,
		DevicePriceMax:     req.DevicePriceMax,
		PlannedStartTime:   req.PlannedStartTime,
		PlannedEndTime:     req.PlannedEndTime,
		TimeZone:           req.TimeZone,
		DaypartingType:     int8(req.DaypartingType),
		DaypartingSchedule: req.DaypartingSchedule,
		FrequencyCapType:   int8(req.FrequencyCapType),
		FrequencyCapTimes:  req.FrequencyCapTimes,
		FrequencyCapDays:   req.FrequencyCapDays,
		BudgetType:         int8(req.BudgetType),
		BudgetAmount:       req.BudgetAmount,
		Website:            req.Website,
		IOSDownloadURL:     req.IosDownloadURL,
		AndroidDownloadURL: req.AndroidDownloadURL,
		AttachmentURLs:     req.AttachmentUrls,
	}

	// 调用service创建Campaign
	campaign, attachments, err := svc.CreateCampaign(userID, teamID, serviceReq)
	if err != nil {
		utils.Error(c, 500, err.Error())
		return
	}

	// 构建响应
	resp := &campaignModel.CreateCampaignResp{
		Campaign: convertToCampaignInfo(campaign, attachments),
		BaseResp: &commonModel.BaseResp{
			Code:    0,
			Message: "Campaign created successfully",
		},
	}

	c.JSON(consts.StatusOK, resp)
}

// UpdateCampaign 更新Campaign
// @router /campaign/update [POST]
func UpdateCampaign(ctx context.Context, c *app.RequestContext) {
	var req campaignModel.UpdateCampaignReq
	if err := c.BindAndValidate(&req); err != nil {
		utils.ParamError(c, "Invalid request parameters: "+err.Error())
		return
	}

	// 从context获取用户ID
	userID, exists := mw.GetAuthUserID(c)
	if !exists {
		utils.Error(c, 401, "User not authenticated")
		return
	}

	// 构建service请求
	serviceReq := &campaignService.UpdateCampaignRequest{
		CampaignName:       req.CampaignName,
		PromotionObjective: req.PromotionObjective,
		OptimizationGoal:   req.OptimizationGoal,
		Location:           req.Location,
		Age:                req.Age,
		Gender:             req.Gender,
		Languages:          req.Languages,
		SpendingPower:      req.SpendingPower,
		OperatingSystem:    req.OperatingSystem,
		OSVersions:         req.OsVersions,
		DeviceModels:       req.DeviceModels,
		ConnectionTypes:    req.ConnectionTypes,
		Website:            req.Website,
		IOSDownloadURL:     req.IosDownloadURL,
		AndroidDownloadURL: req.AndroidDownloadURL,
		AttachmentURLs:     req.AttachmentUrls,
	}

	if req.DevicePriceType != nil {
		devicePriceType := int8(*req.DevicePriceType)
		serviceReq.DevicePriceType = &devicePriceType
	}
	if req.DevicePriceMin != nil {
		serviceReq.DevicePriceMin = req.DevicePriceMin
	}
	if req.DevicePriceMax != nil {
		serviceReq.DevicePriceMax = req.DevicePriceMax
	}
	if req.PlannedStartTime != nil {
		serviceReq.PlannedStartTime = req.PlannedStartTime
	}
	if req.PlannedEndTime != nil {
		serviceReq.PlannedEndTime = req.PlannedEndTime
	}
	if req.TimeZone != nil {
		serviceReq.TimeZone = req.TimeZone
	}
	if req.DaypartingType != nil {
		daypartingType := int8(*req.DaypartingType)
		serviceReq.DaypartingType = &daypartingType
	}
	if req.DaypartingSchedule != nil {
		serviceReq.DaypartingSchedule = req.DaypartingSchedule
	}
	if req.FrequencyCapType != nil {
		frequencyCapType := int8(*req.FrequencyCapType)
		serviceReq.FrequencyCapType = &frequencyCapType
	}
	if req.FrequencyCapTimes != nil {
		serviceReq.FrequencyCapTimes = req.FrequencyCapTimes
	}
	if req.FrequencyCapDays != nil {
		serviceReq.FrequencyCapDays = req.FrequencyCapDays
	}
	if req.BudgetType != nil {
		budgetType := int8(*req.BudgetType)
		serviceReq.BudgetType = &budgetType
	}
	if req.BudgetAmount != nil {
		serviceReq.BudgetAmount = req.BudgetAmount
	}

	// 调用service更新Campaign
	campaign, attachments, err := svc.UpdateCampaign(userID, req.CampaignID, serviceReq)
	if err != nil {
		utils.Error(c, 500, err.Error())
		return
	}

	// 构建响应
	resp := &campaignModel.UpdateCampaignResp{
		Campaign: convertToCampaignInfo(campaign, attachments),
		BaseResp: &commonModel.BaseResp{
			Code:    0,
			Message: "Campaign updated successfully",
		},
	}

	c.JSON(consts.StatusOK, resp)
}

// UpdateCampaignStatus 更新Campaign状态
// @router /campaign/status [POST]
func UpdateCampaignStatus(ctx context.Context, c *app.RequestContext) {
	var req campaignModel.UpdateCampaignStatusReq
	if err := c.BindAndValidate(&req); err != nil {
		utils.ParamError(c, "Invalid request parameters: "+err.Error())
		return
	}

	// 从context获取用户ID
	userID, exists := mw.GetAuthUserID(c)
	if !exists {
		utils.Error(c, 401, "User not authenticated")
		return
	}

	// 调用service更新状态
	if err := svc.UpdateCampaignStatus(userID, req.CampaignID, req.Status); err != nil {
		utils.Error(c, 500, err.Error())
		return
	}

	resp := &campaignModel.UpdateCampaignStatusResp{
		BaseResp: &commonModel.BaseResp{
			Code:    0,
			Message: "Campaign status updated successfully",
		},
	}

	c.JSON(consts.StatusOK, resp)
}

// ListCampaigns 获取Campaign列表
// @router /campaign/list [POST]
func ListCampaigns(ctx context.Context, c *app.RequestContext) {
	var req campaignModel.ListCampaignsReq
	if err := c.BindAndValidate(&req); err != nil {
		utils.ParamError(c, "Invalid request parameters: "+err.Error())
		return
	}

	// 从context获取用户ID和用户信息
	userID, exists := mw.GetAuthUserID(c)
	if !exists {
		utils.Error(c, 401, "User not authenticated")
		return
	}

	user, exists := mw.GetAuthUser(c)
	if !exists || user.CurrentTeamID == nil {
		utils.Error(c, 400, "User has no team")
		return
	}

	teamID := *user.CurrentTeamID

	keyword := ""
	if req.Keyword != nil {
		keyword = *req.Keyword
	}

	status := ""
	if req.Status != nil {
		status = *req.Status
	}

	promotionObjective := ""
	if req.PromotionObjective != nil {
		promotionObjective = *req.PromotionObjective
	}

	// 调用service获取列表
	campaigns, total, err := svc.ListCampaigns(userID, teamID, keyword, status, promotionObjective, int(req.Page), int(req.PageSize))
	if err != nil {
		utils.Error(c, 500, err.Error())
		return
	}

	// 构建响应
	campaignInfos := make([]*campaignModel.CampaignInfo, 0, len(campaigns))
	for _, campaign := range campaigns {
		// 获取附件
		attachments, _ := mysql.NewCampaignRepository(mysql.DB).GetAttachmentsByCampaignID(campaign.ID)
		campaignInfos = append(campaignInfos, convertToCampaignInfo(campaign, attachments))
	}

	totalPages := int32(math.Ceil(float64(total) / float64(req.PageSize)))

	resp := &campaignModel.ListCampaignsResp{
		Campaigns: campaignInfos,
		PageInfo: &commonModel.PageResp{
			Page:       req.Page,
			PageSize:   req.PageSize,
			Total:      total,
			TotalPages: totalPages,
		},
		BaseResp: &commonModel.BaseResp{
			Code:    0,
			Message: "Success",
		},
	}

	c.JSON(consts.StatusOK, resp)
}

// GetCampaign 获取Campaign详情
// @router /campaign/detail [POST]
func GetCampaign(ctx context.Context, c *app.RequestContext) {
	var req campaignModel.GetCampaignReq
	if err := c.BindAndValidate(&req); err != nil {
		utils.ParamError(c, "Invalid request parameters: "+err.Error())
		return
	}

	// 从context获取用户ID
	userID, exists := mw.GetAuthUserID(c)
	if !exists {
		utils.Error(c, 401, "User not authenticated")
		return
	}

	// 调用service获取详情
	campaign, attachments, err := svc.GetCampaign(userID, req.CampaignID)
	if err != nil {
		utils.Error(c, 500, err.Error())
		return
	}

	resp := &campaignModel.GetCampaignResp{
		Campaign: convertToCampaignInfo(campaign, attachments),
		BaseResp: &commonModel.BaseResp{
			Code:    0,
			Message: "Success",
		},
	}

	c.JSON(consts.StatusOK, resp)
}

// AdminListCampaigns 管理员获取所有Campaign列表
// @router /admin/campaign/list [POST]
func AdminListCampaigns(ctx context.Context, c *app.RequestContext) {
	var req campaignModel.AdminListCampaignsReq
	if err := c.BindAndValidate(&req); err != nil {
		utils.ParamError(c, "Invalid request parameters: "+err.Error())
		return
	}

	keyword := ""
	if req.Keyword != nil {
		keyword = *req.Keyword
	}

	status := ""
	if req.Status != nil {
		status = *req.Status
	}

	promotionObjective := ""
	if req.PromotionObjective != nil {
		promotionObjective = *req.PromotionObjective
	}

	// 调用service获取列表
	campaigns, total, err := svc.AdminListCampaigns(keyword, status, promotionObjective, req.UserID, req.TeamID, int(req.Page), int(req.PageSize))
	if err != nil {
		utils.Error(c, 500, err.Error())
		return
	}

	// 构建响应
	campaignInfos := make([]*campaignModel.CampaignInfo, 0, len(campaigns))
	for _, campaign := range campaigns {
		// 获取附件
		attachments, _ := mysql.NewCampaignRepository(mysql.DB).GetAttachmentsByCampaignID(campaign.ID)
		campaignInfos = append(campaignInfos, convertToCampaignInfo(campaign, attachments))
	}

	totalPages := int32(math.Ceil(float64(total) / float64(req.PageSize)))

	resp := &campaignModel.AdminListCampaignsResp{
		Campaigns: campaignInfos,
		PageInfo: &commonModel.PageResp{
			Page:       req.Page,
			PageSize:   req.PageSize,
			Total:      total,
			TotalPages: totalPages,
		},
		BaseResp: &commonModel.BaseResp{
			Code:    0,
			Message: "Success",
		},
	}

	c.JSON(consts.StatusOK, resp)
}

// AdminUpdateCampaignStatus 管理员更新Campaign状态
// @router /admin/campaign/status [POST]
func AdminUpdateCampaignStatus(ctx context.Context, c *app.RequestContext) {
	var req campaignModel.AdminUpdateCampaignStatusReq
	if err := c.BindAndValidate(&req); err != nil {
		utils.ParamError(c, "Invalid request parameters: "+err.Error())
		return
	}

	// 调用service更新状态
	if err := svc.AdminUpdateCampaignStatus(req.CampaignID, req.Status); err != nil {
		utils.Error(c, 500, err.Error())
		return
	}

	resp := &campaignModel.AdminUpdateCampaignStatusResp{
		BaseResp: &commonModel.BaseResp{
			Code:    0,
			Message: "Campaign status updated successfully",
		},
	}

	c.JSON(consts.StatusOK, resp)
}

// Helper functions

// convertToCampaignInfo 转换为CampaignInfo
func convertToCampaignInfo(campaign *mysql.Campaign, attachments []*mysql.CampaignAttachment) *campaignModel.CampaignInfo {
	info := &campaignModel.CampaignInfo{
		ID:                 campaign.ID,
		CampaignID:         campaign.CampaignID,
		UserID:             campaign.UserID,
		TeamID:             campaign.TeamID,
		CampaignName:       campaign.CampaignName,
		PromotionObjective: campaign.PromotionObjective,
		OptimizationGoal:   campaign.OptimizationGoal,
		DevicePriceType:    int32(campaign.DevicePriceType),
		PlannedStartTime:   campaign.PlannedStartTime.Format("2006-01-02T15:04:05Z07:00"),
		PlannedEndTime:     campaign.PlannedEndTime.Format("2006-01-02T15:04:05Z07:00"),
		DaypartingType:     int32(campaign.DaypartingType),
		FrequencyCapType:   int32(campaign.FrequencyCapType),
		BudgetType:         int32(campaign.BudgetType),
		BudgetAmount:       campaign.BudgetAmount,
		Status:             campaign.Status,
		CreatedAt:          campaign.CreatedAt.Format("2006-01-02T15:04:05Z07:00"),
		UpdatedAt:          campaign.UpdatedAt.Format("2006-01-02T15:04:05Z07:00"),
	}

	// 解析JSON字段
	if campaign.Location != nil && *campaign.Location != "" {
		var location []int64
		json.Unmarshal([]byte(*campaign.Location), &location)
		info.Location = location
	}

	if campaign.Age != nil {
		info.Age = campaign.Age
	}

	if campaign.Gender != nil {
		info.Gender = campaign.Gender
	}

	if campaign.Languages != nil && *campaign.Languages != "" {
		var languages []int64
		json.Unmarshal([]byte(*campaign.Languages), &languages)
		info.Languages = languages
	}

	if campaign.SpendingPower != nil {
		info.SpendingPower = campaign.SpendingPower
	}

	if campaign.OperatingSystem != nil {
		info.OperatingSystem = campaign.OperatingSystem
	}

	if campaign.OSVersions != nil && *campaign.OSVersions != "" {
		var osVersions []int64
		json.Unmarshal([]byte(*campaign.OSVersions), &osVersions)
		info.OsVersions = osVersions
	}

	if campaign.DeviceModels != nil && *campaign.DeviceModels != "" {
		var deviceModels []int64
		json.Unmarshal([]byte(*campaign.DeviceModels), &deviceModels)
		info.DeviceModels = deviceModels
	}

	if campaign.ConnectionTypes != nil && *campaign.ConnectionTypes != "" {
		var connectionTypes []int64
		json.Unmarshal([]byte(*campaign.ConnectionTypes), &connectionTypes)
		info.ConnectionTypes = connectionTypes
	}

	if campaign.DevicePriceMin != nil {
		info.DevicePriceMin = campaign.DevicePriceMin
	}

	if campaign.DevicePriceMax != nil {
		info.DevicePriceMax = campaign.DevicePriceMax
	}

	if campaign.TimeZone != nil {
		info.TimeZone = campaign.TimeZone
	}

	if campaign.DaypartingSchedule != nil {
		info.DaypartingSchedule = campaign.DaypartingSchedule
	}

	if campaign.FrequencyCapTimes != nil {
		info.FrequencyCapTimes = campaign.FrequencyCapTimes
	}

	if campaign.FrequencyCapDays != nil {
		info.FrequencyCapDays = campaign.FrequencyCapDays
	}

	if campaign.Website != nil {
		info.Website = campaign.Website
	}

	if campaign.IOSDownloadURL != nil {
		info.IosDownloadURL = campaign.IOSDownloadURL
	}

	if campaign.AndroidDownloadURL != nil {
		info.AndroidDownloadURL = campaign.AndroidDownloadURL
	}

	// 转换附件
	attachmentInfos := make([]*campaignModel.CampaignAttachment, 0, len(attachments))
	for _, attachment := range attachments {
		var fileSize int64
		if attachment.FileSize != nil {
			fileSize = *attachment.FileSize
		}
		attachmentInfos = append(attachmentInfos, &campaignModel.CampaignAttachment{
			ID:        attachment.ID,
			FileURL:   attachment.FileURL,
			FileName:  attachment.FileName,
			FileType:  attachment.FileType,
			FileSize:  fileSize,
			CreatedAt: attachment.CreatedAt.Format("2006-01-02T15:04:05Z07:00"),
		})
	}
	info.Attachments = attachmentInfos

	return info
}
