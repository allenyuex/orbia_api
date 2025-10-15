package kol

import (
	"errors"
	"fmt"
	"time"

	"orbia_api/biz/dal/mysql"

	"gorm.io/gorm"
)

// KolService KOL服务接口
type KolService interface {
	// KOL申请和基本信息管理
	ApplyKol(userID int64, displayName, description, country string, avatarURL, tiktokURL, youtubeURL, xURL, discordURL *string, languageCodes, languageNames, tags []string) (int64, error)
	GetKolInfo(kolID *int64, userID *int64) (*mysql.Kol, []*mysql.KolLanguage, []*mysql.KolTag, *mysql.KolStats, error)
	UpdateKolInfo(userID int64, displayName, description, country, avatarURL, tiktokURL, youtubeURL, xURL, discordURL *string, languageCodes, languageNames, tags *[]string) error
	ReviewKol(kolID int64, status, rejectReason string) error
	GetKolList(status, country, tag *string, page, pageSize int) ([]*mysql.Kol, int64, error)

	// KOL统计数据管理
	UpdateKolStats(userID int64, totalFollowers, tiktokFollowers, youtubeSubscribers, xFollowers, discordMembers, tiktokAvgViews *int64, engagementRate *float64) error

	// KOL报价Plans管理
	SaveKolPlan(userID int64, planID *int64, title, description string, price float64, planType string) (int64, error)
	DeleteKolPlan(userID, planID int64) error
	GetKolPlans(kolID *int64, userID *int64) ([]*mysql.KolPlan, error)

	// KOL视频管理
	CreateKolVideo(userID int64, title, content, coverURL, videoURL, platform string, platformVideoID *string, likesCount, viewsCount, commentsCount, sharesCount *int64, publishedAt *string) (int64, error)
	UpdateKolVideo(userID, videoID int64, title, content, coverURL, videoURL *string, likesCount, viewsCount, commentsCount, sharesCount *int64) error
	DeleteKolVideo(userID, videoID int64) error
	GetKolVideos(kolID *int64, userID *int64, page, pageSize int) ([]*mysql.KolVideo, int64, error)
}

// kolService KOL服务实现
type kolService struct {
	kolRepo  mysql.KolRepository
	userRepo mysql.UserRepository
}

// NewKolService 创建KOL服务实例
func NewKolService(kolRepo mysql.KolRepository, userRepo mysql.UserRepository) KolService {
	return &kolService{
		kolRepo:  kolRepo,
		userRepo: userRepo,
	}
}

// ApplyKol 申请成为KOL
func (s *kolService) ApplyKol(userID int64, displayName, description, country string, avatarURL, tiktokURL, youtubeURL, xURL, discordURL *string, languageCodes, languageNames, tags []string) (int64, error) {
	// 验证用户是否存在
	_, err := s.userRepo.GetUserByID(userID)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return 0, errors.New("user not found")
		}
		return 0, fmt.Errorf("failed to get user: %v", err)
	}

	// 检查用户是否已经申请过KOL
	existingKol, err := s.kolRepo.GetKolByUserID(userID)
	if err != nil && !errors.Is(err, gorm.ErrRecordNotFound) {
		return 0, fmt.Errorf("failed to check existing KOL: %v", err)
	}
	if existingKol != nil {
		return 0, errors.New("user has already applied for KOL")
	}

	// 验证语言数组长度一致
	if len(languageCodes) != len(languageNames) {
		return 0, errors.New("language codes and names must have the same length")
	}

	// 创建KOL记录
	kol := &mysql.Kol{
		UserID:      userID,
		AvatarURL:   avatarURL,
		DisplayName: &displayName,
		Description: &description,
		Country:     &country,
		TiktokURL:   tiktokURL,
		YoutubeURL:  youtubeURL,
		XURL:        xURL,
		DiscordURL:  discordURL,
		Status:      "pending",
	}

	if err := s.kolRepo.CreateKol(kol); err != nil {
		return 0, fmt.Errorf("failed to create KOL: %v", err)
	}

	// 创建语言记录
	for i := range languageCodes {
		language := &mysql.KolLanguage{
			KolID:        kol.ID,
			LanguageCode: languageCodes[i],
			LanguageName: languageNames[i],
		}
		if err := s.kolRepo.CreateKolLanguage(language); err != nil {
			return 0, fmt.Errorf("failed to create KOL language: %v", err)
		}
	}

	// 创建标签记录
	for _, tag := range tags {
		kolTag := &mysql.KolTag{
			KolID: kol.ID,
			Tag:   tag,
		}
		if err := s.kolRepo.CreateKolTag(kolTag); err != nil {
			return 0, fmt.Errorf("failed to create KOL tag: %v", err)
		}
	}

	// 创建初始统计数据
	stats := &mysql.KolStats{
		KolID: kol.ID,
	}
	if err := s.kolRepo.CreateKolStats(stats); err != nil {
		return 0, fmt.Errorf("failed to create KOL stats: %v", err)
	}

	// 更新用户表中的 kol_id
	user, err := s.userRepo.GetUserByID(userID)
	if err != nil {
		return 0, fmt.Errorf("failed to get user: %v", err)
	}
	user.KolID = &kol.ID
	if err := s.userRepo.UpdateUser(user); err != nil {
		return 0, fmt.Errorf("failed to update user kol_id: %v", err)
	}

	return kol.ID, nil
}

// GetKolInfo 获取KOL信息
func (s *kolService) GetKolInfo(kolID *int64, userID *int64) (*mysql.Kol, []*mysql.KolLanguage, []*mysql.KolTag, *mysql.KolStats, error) {
	var kol *mysql.Kol
	var err error

	if kolID != nil {
		kol, err = s.kolRepo.GetKolByID(*kolID)
	} else if userID != nil {
		kol, err = s.kolRepo.GetKolByUserID(*userID)
	} else {
		return nil, nil, nil, nil, errors.New("either kol_id or user_id must be provided")
	}

	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil, nil, nil, errors.New("KOL not found")
		}
		return nil, nil, nil, nil, fmt.Errorf("failed to get KOL: %v", err)
	}

	// 获取语言列表
	languages, err := s.kolRepo.GetKolLanguages(kol.ID)
	if err != nil {
		return nil, nil, nil, nil, fmt.Errorf("failed to get KOL languages: %v", err)
	}

	// 获取标签列表
	tags, err := s.kolRepo.GetKolTags(kol.ID)
	if err != nil {
		return nil, nil, nil, nil, fmt.Errorf("failed to get KOL tags: %v", err)
	}

	// 获取统计数据
	stats, err := s.kolRepo.GetKolStats(kol.ID)
	if err != nil && !errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, nil, nil, nil, fmt.Errorf("failed to get KOL stats: %v", err)
	}

	return kol, languages, tags, stats, nil
}

// UpdateKolInfo 更新KOL信息
func (s *kolService) UpdateKolInfo(userID int64, displayName, description, country, avatarURL, tiktokURL, youtubeURL, xURL, discordURL *string, languageCodes, languageNames, tags *[]string) error {
	// 获取KOL信息
	kol, err := s.kolRepo.GetKolByUserID(userID)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return errors.New("KOL not found")
		}
		return fmt.Errorf("failed to get KOL: %v", err)
	}

	// 更新基本信息
	if displayName != nil {
		kol.DisplayName = displayName
	}
	if description != nil {
		kol.Description = description
	}
	if country != nil {
		kol.Country = country
	}
	if avatarURL != nil {
		kol.AvatarURL = avatarURL
	}
	if tiktokURL != nil {
		kol.TiktokURL = tiktokURL
	}
	if youtubeURL != nil {
		kol.YoutubeURL = youtubeURL
	}
	if xURL != nil {
		kol.XURL = xURL
	}
	if discordURL != nil {
		kol.DiscordURL = discordURL
	}

	if err := s.kolRepo.UpdateKol(kol); err != nil {
		return fmt.Errorf("failed to update KOL: %v", err)
	}

	// 更新语言（如果提供）
	if languageCodes != nil && languageNames != nil {
		if len(*languageCodes) != len(*languageNames) {
			return errors.New("language codes and names must have the same length")
		}

		// 删除旧的语言记录
		if err := s.kolRepo.DeleteKolLanguages(kol.ID); err != nil {
			return fmt.Errorf("failed to delete old languages: %v", err)
		}

		// 创建新的语言记录
		for i := range *languageCodes {
			language := &mysql.KolLanguage{
				KolID:        kol.ID,
				LanguageCode: (*languageCodes)[i],
				LanguageName: (*languageNames)[i],
			}
			if err := s.kolRepo.CreateKolLanguage(language); err != nil {
				return fmt.Errorf("failed to create KOL language: %v", err)
			}
		}
	}

	// 更新标签（如果提供）
	if tags != nil {
		// 删除旧的标签记录
		if err := s.kolRepo.DeleteKolTags(kol.ID); err != nil {
			return fmt.Errorf("failed to delete old tags: %v", err)
		}

		// 创建新的标签记录
		for _, tag := range *tags {
			kolTag := &mysql.KolTag{
				KolID: kol.ID,
				Tag:   tag,
			}
			if err := s.kolRepo.CreateKolTag(kolTag); err != nil {
				return fmt.Errorf("failed to create KOL tag: %v", err)
			}
		}
	}

	return nil
}

// ReviewKol 审核KOL（管理员使用）
func (s *kolService) ReviewKol(kolID int64, status, rejectReason string) error {
	// 验证状态
	if status != "approved" && status != "rejected" {
		return errors.New("invalid status, must be approved or rejected")
	}

	// 获取KOL信息
	kol, err := s.kolRepo.GetKolByID(kolID)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return errors.New("KOL not found")
		}
		return fmt.Errorf("failed to get KOL: %v", err)
	}

	// 更新状态
	kol.Status = status
	if status == "rejected" && rejectReason != "" {
		kol.RejectReason = &rejectReason
	}
	if status == "approved" {
		now := time.Now()
		kol.ApprovedAt = &now
	}

	if err := s.kolRepo.UpdateKol(kol); err != nil {
		return fmt.Errorf("failed to update KOL status: %v", err)
	}

	return nil
}

// GetKolList 获取KOL列表
func (s *kolService) GetKolList(status, country, tag *string, page, pageSize int) ([]*mysql.Kol, int64, error) {
	if page <= 0 {
		page = 1
	}
	if pageSize <= 0 {
		pageSize = 10
	}

	offset := (page - 1) * pageSize
	kols, total, err := s.kolRepo.GetKolList(status, country, tag, offset, pageSize)
	if err != nil {
		return nil, 0, fmt.Errorf("failed to get KOL list: %v", err)
	}

	return kols, total, nil
}

// UpdateKolStats 更新KOL统计数据
func (s *kolService) UpdateKolStats(userID int64, totalFollowers, tiktokFollowers, youtubeSubscribers, xFollowers, discordMembers, tiktokAvgViews *int64, engagementRate *float64) error {
	// 获取KOL信息
	kol, err := s.kolRepo.GetKolByUserID(userID)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return errors.New("KOL not found")
		}
		return fmt.Errorf("failed to get KOL: %v", err)
	}

	// 获取统计数据
	stats, err := s.kolRepo.GetKolStats(kol.ID)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			// 如果不存在，创建新的统计数据
			stats = &mysql.KolStats{
				KolID: kol.ID,
			}
		} else {
			return fmt.Errorf("failed to get KOL stats: %v", err)
		}
	}

	// 更新统计数据
	if totalFollowers != nil {
		stats.TotalFollowers = *totalFollowers
	}
	if tiktokFollowers != nil {
		stats.TiktokFollowers = *tiktokFollowers
	}
	if youtubeSubscribers != nil {
		stats.YoutubeSubscribers = *youtubeSubscribers
	}
	if xFollowers != nil {
		stats.XFollowers = *xFollowers
	}
	if discordMembers != nil {
		stats.DiscordMembers = *discordMembers
	}
	if tiktokAvgViews != nil {
		stats.TiktokAvgViews = *tiktokAvgViews
	}
	if engagementRate != nil {
		stats.EngagementRate = *engagementRate
	}

	if stats.ID == 0 {
		// 创建
		if err := s.kolRepo.CreateKolStats(stats); err != nil {
			return fmt.Errorf("failed to create KOL stats: %v", err)
		}
	} else {
		// 更新
		if err := s.kolRepo.UpdateKolStats(stats); err != nil {
			return fmt.Errorf("failed to update KOL stats: %v", err)
		}
	}

	return nil
}

// SaveKolPlan 创建或更新KOL报价Plan
func (s *kolService) SaveKolPlan(userID int64, planID *int64, title, description string, price float64, planType string) (int64, error) {
	// 验证planType
	if planType != "basic" && planType != "standard" && planType != "premium" {
		return 0, errors.New("invalid plan type, must be basic, standard, or premium")
	}

	// 获取KOL信息
	kol, err := s.kolRepo.GetKolByUserID(userID)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return 0, errors.New("KOL not found")
		}
		return 0, fmt.Errorf("failed to get KOL: %v", err)
	}

	if planID != nil && *planID > 0 {
		// 更新现有Plan
		plan, err := s.kolRepo.GetKolPlanByID(*planID)
		if err != nil {
			if errors.Is(err, gorm.ErrRecordNotFound) {
				return 0, errors.New("plan not found")
			}
			return 0, fmt.Errorf("failed to get plan: %v", err)
		}

		// 验证Plan是否属于该KOL
		if plan.KolID != kol.ID {
			return 0, errors.New("plan does not belong to this KOL")
		}

		// 更新Plan信息
		plan.Title = title
		plan.Description = &description
		plan.Price = price
		plan.PlanType = planType

		if err := s.kolRepo.UpdateKolPlan(plan); err != nil {
			return 0, fmt.Errorf("failed to update plan: %v", err)
		}

		return plan.ID, nil
	} else {
		// 创建新Plan
		plan := &mysql.KolPlan{
			KolID:       kol.ID,
			Title:       title,
			Description: &description,
			Price:       price,
			PlanType:    planType,
		}

		if err := s.kolRepo.CreateKolPlan(plan); err != nil {
			return 0, fmt.Errorf("failed to create plan: %v", err)
		}

		return plan.ID, nil
	}
}

// DeleteKolPlan 删除KOL报价Plan
func (s *kolService) DeleteKolPlan(userID, planID int64) error {
	// 获取KOL信息
	kol, err := s.kolRepo.GetKolByUserID(userID)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return errors.New("KOL not found")
		}
		return fmt.Errorf("failed to get KOL: %v", err)
	}

	// 获取Plan信息
	plan, err := s.kolRepo.GetKolPlanByID(planID)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return errors.New("plan not found")
		}
		return fmt.Errorf("failed to get plan: %v", err)
	}

	// 验证Plan是否属于该KOL
	if plan.KolID != kol.ID {
		return errors.New("plan does not belong to this KOL")
	}

	// 删除Plan
	if err := s.kolRepo.DeleteKolPlan(planID); err != nil {
		return fmt.Errorf("failed to delete plan: %v", err)
	}

	return nil
}

// GetKolPlans 获取KOL报价Plans列表
func (s *kolService) GetKolPlans(kolID *int64, userID *int64) ([]*mysql.KolPlan, error) {
	var kol *mysql.Kol
	var err error

	if kolID != nil {
		kol, err = s.kolRepo.GetKolByID(*kolID)
	} else if userID != nil {
		kol, err = s.kolRepo.GetKolByUserID(*userID)
	} else {
		return nil, errors.New("either kol_id or user_id must be provided")
	}

	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("KOL not found")
		}
		return nil, fmt.Errorf("failed to get KOL: %v", err)
	}

	plans, err := s.kolRepo.GetKolPlans(kol.ID)
	if err != nil {
		return nil, fmt.Errorf("failed to get plans: %v", err)
	}

	return plans, nil
}

// CreateKolVideo 创建KOL视频
func (s *kolService) CreateKolVideo(userID int64, title, content, coverURL, videoURL, platform string, platformVideoID *string, likesCount, viewsCount, commentsCount, sharesCount *int64, publishedAt *string) (int64, error) {
	// 验证平台
	if platform != "tiktok" && platform != "youtube" {
		return 0, errors.New("invalid platform, must be tiktok or youtube")
	}

	// 获取KOL信息
	kol, err := s.kolRepo.GetKolByUserID(userID)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return 0, errors.New("KOL not found")
		}
		return 0, fmt.Errorf("failed to get KOL: %v", err)
	}

	// 解析发布时间
	var publishedTime *time.Time
	if publishedAt != nil && *publishedAt != "" {
		t, err := time.Parse(time.RFC3339, *publishedAt)
		if err != nil {
			return 0, fmt.Errorf("invalid published_at format, must be RFC3339: %v", err)
		}
		publishedTime = &t
	}

	// 创建视频记录
	video := &mysql.KolVideo{
		KolID:           kol.ID,
		Title:           title,
		Content:         &content,
		CoverURL:        &coverURL,
		VideoURL:        &videoURL,
		Platform:        platform,
		PlatformVideoID: platformVideoID,
		PublishedAt:     publishedTime,
	}

	if likesCount != nil {
		video.LikesCount = *likesCount
	}
	if viewsCount != nil {
		video.ViewsCount = *viewsCount
	}
	if commentsCount != nil {
		video.CommentsCount = *commentsCount
	}
	if sharesCount != nil {
		video.SharesCount = *sharesCount
	}

	if err := s.kolRepo.CreateKolVideo(video); err != nil {
		return 0, fmt.Errorf("failed to create video: %v", err)
	}

	return video.ID, nil
}

// UpdateKolVideo 更新KOL视频
func (s *kolService) UpdateKolVideo(userID, videoID int64, title, content, coverURL, videoURL *string, likesCount, viewsCount, commentsCount, sharesCount *int64) error {
	// 获取KOL信息
	kol, err := s.kolRepo.GetKolByUserID(userID)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return errors.New("KOL not found")
		}
		return fmt.Errorf("failed to get KOL: %v", err)
	}

	// 获取视频信息
	video, err := s.kolRepo.GetKolVideoByID(videoID)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return errors.New("video not found")
		}
		return fmt.Errorf("failed to get video: %v", err)
	}

	// 验证视频是否属于该KOL
	if video.KolID != kol.ID {
		return errors.New("video does not belong to this KOL")
	}

	// 更新视频信息
	if title != nil {
		video.Title = *title
	}
	if content != nil {
		video.Content = content
	}
	if coverURL != nil {
		video.CoverURL = coverURL
	}
	if videoURL != nil {
		video.VideoURL = videoURL
	}
	if likesCount != nil {
		video.LikesCount = *likesCount
	}
	if viewsCount != nil {
		video.ViewsCount = *viewsCount
	}
	if commentsCount != nil {
		video.CommentsCount = *commentsCount
	}
	if sharesCount != nil {
		video.SharesCount = *sharesCount
	}

	if err := s.kolRepo.UpdateKolVideo(video); err != nil {
		return fmt.Errorf("failed to update video: %v", err)
	}

	return nil
}

// DeleteKolVideo 删除KOL视频
func (s *kolService) DeleteKolVideo(userID, videoID int64) error {
	// 获取KOL信息
	kol, err := s.kolRepo.GetKolByUserID(userID)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return errors.New("KOL not found")
		}
		return fmt.Errorf("failed to get KOL: %v", err)
	}

	// 获取视频信息
	video, err := s.kolRepo.GetKolVideoByID(videoID)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return errors.New("video not found")
		}
		return fmt.Errorf("failed to get video: %v", err)
	}

	// 验证视频是否属于该KOL
	if video.KolID != kol.ID {
		return errors.New("video does not belong to this KOL")
	}

	// 删除视频
	if err := s.kolRepo.DeleteKolVideo(videoID); err != nil {
		return fmt.Errorf("failed to delete video: %v", err)
	}

	return nil
}

// GetKolVideos 获取KOL视频列表
func (s *kolService) GetKolVideos(kolID *int64, userID *int64, page, pageSize int) ([]*mysql.KolVideo, int64, error) {
	var kol *mysql.Kol
	var err error

	if kolID != nil {
		kol, err = s.kolRepo.GetKolByID(*kolID)
	} else if userID != nil {
		kol, err = s.kolRepo.GetKolByUserID(*userID)
	} else {
		return nil, 0, errors.New("either kol_id or user_id must be provided")
	}

	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, 0, errors.New("KOL not found")
		}
		return nil, 0, fmt.Errorf("failed to get KOL: %v", err)
	}

	if page <= 0 {
		page = 1
	}
	if pageSize <= 0 {
		pageSize = 10
	}

	offset := (page - 1) * pageSize
	videos, total, err := s.kolRepo.GetKolVideos(kol.ID, offset, pageSize)
	if err != nil {
		return nil, 0, fmt.Errorf("failed to get videos: %v", err)
	}

	return videos, total, nil
}
