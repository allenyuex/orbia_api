package conversation

import (
	"errors"
	"fmt"
	"time"

	"orbia_api/biz/dal/mysql"
	"orbia_api/biz/utils"

	"gorm.io/gorm"
)

// ConversationService 会话服务接口
type ConversationService interface {
	// 创建会话（用于订单创建时自动创建）
	CreateConversation(conversationType, relatedOrderType, relatedOrderID string, title *string, memberUserIDs []int64) (*mysql.Conversation, error)

	// 发送消息
	SendMessage(userID int64, conversationID string, messageType, content string, fileName *string, fileSize *int64, fileType *string) (*MessageWithSender, error)

	// 获取消息列表
	GetMessages(userID int64, conversationID string, beforeTimestamp *int64, limit int) ([]*MessageWithSender, bool, error)

	// 获取会话详情
	GetConversation(userID int64, conversationID string) (*ConversationDetail, error)

	// 获取用户的会话列表
	GetConversations(userID int64, conversationType *string, page, pageSize int) ([]*ConversationItem, int64, error)

	// 标记消息已读
	MarkMessagesRead(userID int64, conversationID string) error
}

// MessageWithSender 带发送者信息的消息
type MessageWithSender struct {
	MessageID       string
	ConversationID  string
	SenderID        int64
	SenderNickname  string
	SenderAvatarURL *string
	MessageType     string
	Content         string
	FileName        *string
	FileSize        *int64
	FileType        *string
	Status          string
	CreatedAt       int64 // 毫秒时间戳
}

// MemberInfo 会话成员信息
type MemberInfo struct {
	UserID    int64
	Nickname  string
	AvatarURL *string
	Role      string
	JoinedAt  time.Time
}

// ConversationDetail 会话详情
type ConversationDetail struct {
	ConversationID   string
	Title            *string
	Type             string
	RelatedOrderType *string
	RelatedOrderID   *string
	Status           string
	LastMessageAt    *time.Time
	Members          []*MemberInfo
	UnreadCount      int
	CreatedAt        time.Time
}

// ConversationItem 会话列表项
type ConversationItem struct {
	ConversationID   string
	Title            *string
	Type             string
	RelatedOrderType *string
	RelatedOrderID   *string
	Status           string
	LastMessage      *MessageWithSender
	UnreadCount      int
	Members          []*MemberInfo
	CreatedAt        time.Time
	LastMessageAt    *time.Time
}

// conversationService 会话服务实现
type conversationService struct {
	convRepo mysql.ConversationRepository
	userRepo mysql.UserRepository
}

// NewConversationService 创建会话服务实例
func NewConversationService(convRepo mysql.ConversationRepository, userRepo mysql.UserRepository) ConversationService {
	return &conversationService{
		convRepo: convRepo,
		userRepo: userRepo,
	}
}

// CreateConversation 创建会话
func (s *conversationService) CreateConversation(conversationType, relatedOrderType, relatedOrderID string, title *string, memberUserIDs []int64) (*mysql.Conversation, error) {
	// 生成会话ID
	conversationID := utils.GenerateConversationID()

	// 创建会话
	conversation := &mysql.Conversation{
		ConversationID:   conversationID,
		Title:            title,
		Type:             conversationType,
		RelatedOrderType: &relatedOrderType,
		RelatedOrderID:   &relatedOrderID,
		Status:           "active",
	}

	if err := s.convRepo.CreateConversation(conversation); err != nil {
		return nil, fmt.Errorf("failed to create conversation: %v", err)
	}

	// 添加会话成员
	for i, userID := range memberUserIDs {
		role := "member"
		if i == 0 {
			role = "creator"
		}

		member := &mysql.ConversationMember{
			ConversationID: conversation.ConversationID,
			UserID:         userID,
			Role:           role,
			UnreadCount:    0,
		}

		if err := s.convRepo.AddConversationMember(member); err != nil {
			return nil, fmt.Errorf("failed to add conversation member: %v", err)
		}
	}

	return conversation, nil
}

// SendMessage 发送消息
func (s *conversationService) SendMessage(userID int64, conversationID string, messageType, content string, fileName *string, fileSize *int64, fileType *string) (*MessageWithSender, error) {
	// 验证用户是否是会话成员
	isMember, err := s.convRepo.IsConversationMember(conversationID, userID)
	if err != nil {
		return nil, fmt.Errorf("failed to check conversation member: %v", err)
	}
	if !isMember {
		return nil, errors.New("user is not a member of this conversation")
	}

	// 获取会话信息
	conversation, err := s.convRepo.GetConversationByConversationID(conversationID)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("conversation not found")
		}
		return nil, fmt.Errorf("failed to get conversation: %v", err)
	}

	// 检查会话状态
	if conversation.Status == "closed" {
		return nil, errors.New("conversation is closed")
	}

	// 生成消息ID
	messageID := utils.GenerateMessageID()

	// 创建消息
	message := &mysql.Message{
		MessageID:      messageID,
		ConversationID: conversationID,
		SenderID:       userID,
		MessageType:    messageType,
		Content:        content,
		FileName:       fileName,
		FileSize:       fileSize,
		FileType:       fileType,
		Status:         "sent",
	}

	if err := s.convRepo.CreateMessage(message); err != nil {
		return nil, fmt.Errorf("failed to create message: %v", err)
	}

	// 更新会话的最后消息时间
	now := time.Now()
	conversation.LastMessageAt = &now
	if err := s.convRepo.UpdateConversation(conversation); err != nil {
		return nil, fmt.Errorf("failed to update conversation: %v", err)
	}

	// 增加其他成员的未读消息数
	members, err := s.convRepo.GetConversationMembers(conversationID)
	if err != nil {
		return nil, fmt.Errorf("failed to get conversation members: %v", err)
	}

	for _, member := range members {
		if member.UserID != userID {
			if err := s.convRepo.IncrementUnreadCount(conversationID, member.UserID); err != nil {
				// 记录错误但不中断流程
				fmt.Printf("failed to increment unread count for user %d: %v\n", member.UserID, err)
			}
		}
	}

	// 获取发送者信息
	sender, err := s.userRepo.GetUserByID(userID)
	if err != nil {
		return nil, fmt.Errorf("failed to get sender info: %v", err)
	}

	senderNickname := "Unknown"
	if sender.Nickname != nil {
		senderNickname = *sender.Nickname
	} else if sender.Email != nil {
		senderNickname = *sender.Email
	}

	// 构建返回结果
	result := &MessageWithSender{
		MessageID:       message.MessageID,
		ConversationID:  message.ConversationID,
		SenderID:        message.SenderID,
		SenderNickname:  senderNickname,
		SenderAvatarURL: sender.AvatarURL,
		MessageType:     message.MessageType,
		Content:         message.Content,
		FileName:        message.FileName,
		FileSize:        message.FileSize,
		FileType:        message.FileType,
		Status:          message.Status,
		CreatedAt:       message.CreatedAt.UnixMilli(),
	}

	return result, nil
}

// GetMessages 获取消息列表
func (s *conversationService) GetMessages(userID int64, conversationID string, beforeTimestamp *int64, limit int) ([]*MessageWithSender, bool, error) {
	// 验证用户是否是会话成员
	isMember, err := s.convRepo.IsConversationMember(conversationID, userID)
	if err != nil {
		return nil, false, fmt.Errorf("failed to check conversation member: %v", err)
	}
	if !isMember {
		return nil, false, errors.New("user is not a member of this conversation")
	}

	// 获取消息列表（多查询一条来判断是否还有更多）
	messages, err := s.convRepo.GetMessages(conversationID, beforeTimestamp, limit+1)
	if err != nil {
		return nil, false, fmt.Errorf("failed to get messages: %v", err)
	}

	hasMore := len(messages) > limit
	if hasMore {
		messages = messages[:limit]
	}

	// 获取所有发送者的用户信息
	senderMap := make(map[int64]*mysql.User)
	for _, msg := range messages {
		if _, exists := senderMap[msg.SenderID]; !exists {
			sender, err := s.userRepo.GetUserByID(msg.SenderID)
			if err != nil {
				// 如果获取用户信息失败，使用默认值
				senderMap[msg.SenderID] = &mysql.User{}
			} else {
				senderMap[msg.SenderID] = sender
			}
		}
	}

	// 构建返回结果
	result := make([]*MessageWithSender, 0, len(messages))
	for _, msg := range messages {
		sender := senderMap[msg.SenderID]
		senderNickname := "Unknown"
		if sender.Nickname != nil {
			senderNickname = *sender.Nickname
		} else if sender.Email != nil {
			senderNickname = *sender.Email
		}

		result = append(result, &MessageWithSender{
			MessageID:       msg.MessageID,
			ConversationID:  msg.ConversationID,
			SenderID:        msg.SenderID,
			SenderNickname:  senderNickname,
			SenderAvatarURL: sender.AvatarURL,
			MessageType:     msg.MessageType,
			Content:         msg.Content,
			FileName:        msg.FileName,
			FileSize:        msg.FileSize,
			FileType:        msg.FileType,
			Status:          msg.Status,
			CreatedAt:       msg.CreatedAt.UnixMilli(),
		})
	}

	return result, hasMore, nil
}

// GetConversation 获取会话详情
func (s *conversationService) GetConversation(userID int64, conversationID string) (*ConversationDetail, error) {
	// 验证用户是否是会话成员
	isMember, err := s.convRepo.IsConversationMember(conversationID, userID)
	if err != nil {
		return nil, fmt.Errorf("failed to check conversation member: %v", err)
	}
	if !isMember {
		return nil, errors.New("user is not a member of this conversation")
	}

	// 获取会话信息
	conversation, err := s.convRepo.GetConversationByConversationID(conversationID)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("conversation not found")
		}
		return nil, fmt.Errorf("failed to get conversation: %v", err)
	}

	// 获取会话成员
	members, err := s.convRepo.GetConversationMembers(conversationID)
	if err != nil {
		return nil, fmt.Errorf("failed to get conversation members: %v", err)
	}

	// 获取成员用户信息
	memberInfos := make([]*MemberInfo, 0, len(members))
	var unreadCount int
	for _, member := range members {
		if member.UserID == userID {
			unreadCount = member.UnreadCount
		}

		user, err := s.userRepo.GetUserByID(member.UserID)
		if err != nil {
			continue
		}

		nickname := "Unknown"
		if user.Nickname != nil {
			nickname = *user.Nickname
		} else if user.Email != nil {
			nickname = *user.Email
		}

		memberInfos = append(memberInfos, &MemberInfo{
			UserID:    member.UserID,
			Nickname:  nickname,
			AvatarURL: user.AvatarURL,
			Role:      member.Role,
			JoinedAt:  member.JoinedAt,
		})
	}

	// 构建返回结果
	result := &ConversationDetail{
		ConversationID:   conversation.ConversationID,
		Title:            conversation.Title,
		Type:             conversation.Type,
		RelatedOrderType: conversation.RelatedOrderType,
		RelatedOrderID:   conversation.RelatedOrderID,
		Status:           conversation.Status,
		LastMessageAt:    conversation.LastMessageAt,
		Members:          memberInfos,
		UnreadCount:      unreadCount,
		CreatedAt:        conversation.CreatedAt,
	}

	return result, nil
}

// GetConversations 获取用户的会话列表
func (s *conversationService) GetConversations(userID int64, conversationType *string, page, pageSize int) ([]*ConversationItem, int64, error) {
	offset := (page - 1) * pageSize

	// 获取会话列表
	conversations, total, err := s.convRepo.GetUserConversations(userID, conversationType, offset, pageSize)
	if err != nil {
		return nil, 0, fmt.Errorf("failed to get conversations: %v", err)
	}

	// 构建返回结果
	result := make([]*ConversationItem, 0, len(conversations))
	for _, conv := range conversations {
		// 获取会话成员
		members, err := s.convRepo.GetConversationMembers(conv.ConversationID)
		if err != nil {
			continue
		}

		// 获取成员用户信息
		memberInfos := make([]*MemberInfo, 0, len(members))
		var unreadCount int
		for _, member := range members {
			if member.UserID == userID {
				unreadCount = member.UnreadCount
			}

			user, err := s.userRepo.GetUserByID(member.UserID)
			if err != nil {
				continue
			}

			nickname := "Unknown"
			if user.Nickname != nil {
				nickname = *user.Nickname
			} else if user.Email != nil {
				nickname = *user.Email
			}

			memberInfos = append(memberInfos, &MemberInfo{
				UserID:    member.UserID,
				Nickname:  nickname,
				AvatarURL: user.AvatarURL,
				Role:      member.Role,
				JoinedAt:  member.JoinedAt,
			})
		}

		// 获取最后一条消息
		var lastMessage *MessageWithSender
		messages, err := s.convRepo.GetMessages(conv.ConversationID, nil, 1)
		if err == nil && len(messages) > 0 {
			msg := messages[0]
			sender, err := s.userRepo.GetUserByID(msg.SenderID)
			if err == nil {
				senderNickname := "Unknown"
				if sender.Nickname != nil {
					senderNickname = *sender.Nickname
				} else if sender.Email != nil {
					senderNickname = *sender.Email
				}

				lastMessage = &MessageWithSender{
					MessageID:       msg.MessageID,
					ConversationID:  msg.ConversationID,
					SenderID:        msg.SenderID,
					SenderNickname:  senderNickname,
					SenderAvatarURL: sender.AvatarURL,
					MessageType:     msg.MessageType,
					Content:         msg.Content,
					FileName:        msg.FileName,
					FileSize:        msg.FileSize,
					FileType:        msg.FileType,
					Status:          msg.Status,
					CreatedAt:       msg.CreatedAt.UnixMilli(),
				}
			}
		}

		result = append(result, &ConversationItem{
			ConversationID:   conv.ConversationID,
			Title:            conv.Title,
			Type:             conv.Type,
			RelatedOrderType: conv.RelatedOrderType,
			RelatedOrderID:   conv.RelatedOrderID,
			Status:           conv.Status,
			LastMessage:      lastMessage,
			UnreadCount:      unreadCount,
			Members:          memberInfos,
			CreatedAt:        conv.CreatedAt,
			LastMessageAt:    conv.LastMessageAt,
		})
	}

	return result, total, nil
}

// MarkMessagesRead 标记消息已读
func (s *conversationService) MarkMessagesRead(userID int64, conversationID string) error {
	// 验证用户是否是会话成员
	isMember, err := s.convRepo.IsConversationMember(conversationID, userID)
	if err != nil {
		return fmt.Errorf("failed to check conversation member: %v", err)
	}
	if !isMember {
		return errors.New("user is not a member of this conversation")
	}

	// 重置未读消息数
	if err := s.convRepo.ResetUnreadCount(conversationID, userID); err != nil {
		return fmt.Errorf("failed to reset unread count: %v", err)
	}

	return nil
}
