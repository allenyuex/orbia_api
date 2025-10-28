package mysql

import (
	"time"

	"gorm.io/gorm"
)

// Conversation 会话模型
type Conversation struct {
	ConversationID   string         `gorm:"primaryKey;column:conversation_id;size:64" json:"conversation_id"`
	Title            *string        `gorm:"column:title;size:200" json:"title"`
	Type             string         `gorm:"column:type;type:enum('kol_order','ad_order','general','support');default:'general';not null" json:"type"`
	RelatedOrderType *string        `gorm:"column:related_order_type;size:50" json:"related_order_type"`
	RelatedOrderID   *string        `gorm:"column:related_order_id;size:64" json:"related_order_id"`
	Status           string         `gorm:"column:status;type:enum('active','archived','closed');default:'active';not null" json:"status"`
	LastMessageAt    *time.Time     `gorm:"column:last_message_at" json:"last_message_at"`
	CreatedAt        time.Time      `gorm:"column:created_at;autoCreateTime" json:"created_at"`
	UpdatedAt        time.Time      `gorm:"column:updated_at;autoUpdateTime" json:"updated_at"`
	DeletedAt        gorm.DeletedAt `gorm:"index;column:deleted_at" json:"-"`
}

// TableName 指定表名
func (Conversation) TableName() string {
	return "orbia_conversation"
}

// ConversationMember 会话成员模型
type ConversationMember struct {
	ConversationID string     `gorm:"primaryKey;column:conversation_id;size:64" json:"conversation_id"`
	UserID         int64      `gorm:"primaryKey;column:user_id" json:"user_id"`
	Role           string     `gorm:"column:role;type:enum('creator','member','admin');default:'member';not null" json:"role"`
	UnreadCount    int        `gorm:"column:unread_count;default:0;not null" json:"unread_count"`
	LastReadAt     *time.Time `gorm:"column:last_read_at" json:"last_read_at"`
	JoinedAt       time.Time  `gorm:"column:joined_at;autoCreateTime" json:"joined_at"`
}

// TableName 指定表名
func (ConversationMember) TableName() string {
	return "orbia_conversation_member"
}

// Message 消息模型
type Message struct {
	MessageID      string         `gorm:"primaryKey;column:message_id;size:64" json:"message_id"`
	ConversationID string         `gorm:"column:conversation_id;size:64;not null;index" json:"conversation_id"`
	SenderID       int64          `gorm:"column:sender_id;not null;index" json:"sender_id"`
	MessageType    string         `gorm:"column:message_type;type:enum('text','image','file','video','audio','system');default:'text';not null" json:"message_type"`
	Content        string         `gorm:"column:content;type:text;not null" json:"content"`
	FileName       *string        `gorm:"column:file_name;size:500" json:"file_name"`
	FileSize       *int64         `gorm:"column:file_size" json:"file_size"`
	FileType       *string        `gorm:"column:file_type;size:100" json:"file_type"`
	Status         string         `gorm:"column:status;type:enum('sent','delivered','read','failed');default:'sent';not null" json:"status"`
	CreatedAt      time.Time      `gorm:"column:created_at;type:timestamp(3);autoCreateTime:milli" json:"created_at"`
	UpdatedAt      time.Time      `gorm:"column:updated_at;autoUpdateTime" json:"updated_at"`
	DeletedAt      gorm.DeletedAt `gorm:"index;column:deleted_at" json:"-"`
}

// TableName 指定表名
func (Message) TableName() string {
	return "orbia_message"
}

// ConversationRepository 会话仓储接口
type ConversationRepository interface {
	// 会话相关
	CreateConversation(conversation *Conversation) error
	GetConversationByConversationID(conversationID string) (*Conversation, error)
	GetConversationByOrderID(orderType, orderID string) (*Conversation, error)
	UpdateConversation(conversation *Conversation) error
	GetUserConversations(userID int64, conversationType *string, offset, limit int) ([]*Conversation, int64, error)

	// 会话成员相关
	AddConversationMember(member *ConversationMember) error
	GetConversationMembers(conversationID string) ([]*ConversationMember, error)
	GetConversationMember(conversationID string, userID int64) (*ConversationMember, error)
	UpdateConversationMember(member *ConversationMember) error
	IsConversationMember(conversationID string, userID int64) (bool, error)

	// 消息相关
	CreateMessage(message *Message) error
	GetMessageByMessageID(messageID string) (*Message, error)
	GetMessages(conversationID string, beforeTimestamp *int64, limit int) ([]*Message, error)
	UpdateMessage(message *Message) error

	// 未读消息相关
	IncrementUnreadCount(conversationID string, userID int64) error
	ResetUnreadCount(conversationID string, userID int64) error
}

// conversationRepository 会话仓储实现
type conversationRepository struct {
	db *gorm.DB
}

// NewConversationRepository 创建会话仓储实例
func NewConversationRepository(db *gorm.DB) ConversationRepository {
	return &conversationRepository{db: db}
}

// CreateConversation 创建会话
func (r *conversationRepository) CreateConversation(conversation *Conversation) error {
	return r.db.Create(conversation).Error
}

// GetConversationByConversationID 根据会话ID获取会话
func (r *conversationRepository) GetConversationByConversationID(conversationID string) (*Conversation, error) {
	var conversation Conversation
	err := r.db.Where("conversation_id = ?", conversationID).First(&conversation).Error
	if err != nil {
		return nil, err
	}
	return &conversation, nil
}

// GetConversationByOrderID 根据订单ID获取会话
func (r *conversationRepository) GetConversationByOrderID(orderType, orderID string) (*Conversation, error) {
	var conversation Conversation
	err := r.db.Where("related_order_type = ? AND related_order_id = ?", orderType, orderID).First(&conversation).Error
	if err != nil {
		return nil, err
	}
	return &conversation, nil
}

// UpdateConversation 更新会话
func (r *conversationRepository) UpdateConversation(conversation *Conversation) error {
	return r.db.Save(conversation).Error
}

// GetUserConversations 获取用户的会话列表
func (r *conversationRepository) GetUserConversations(userID int64, conversationType *string, offset, limit int) ([]*Conversation, int64, error) {
	var conversations []*Conversation
	var total int64

	// 先通过 conversation_member 表找到用户参与的会话ID
	query := r.db.Model(&Conversation{}).
		Joins("JOIN orbia_conversation_member ON orbia_conversation.conversation_id = orbia_conversation_member.conversation_id").
		Where("orbia_conversation_member.user_id = ?", userID)

	if conversationType != nil && *conversationType != "" {
		query = query.Where("orbia_conversation.type = ?", *conversationType)
	}

	// 计数
	if err := query.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	// 查询
	err := query.Order("orbia_conversation.last_message_at DESC, orbia_conversation.created_at DESC").
		Offset(offset).
		Limit(limit).
		Find(&conversations).Error

	if err != nil {
		return nil, 0, err
	}

	return conversations, total, nil
}

// AddConversationMember 添加会话成员
func (r *conversationRepository) AddConversationMember(member *ConversationMember) error {
	return r.db.Create(member).Error
}

// GetConversationMembers 获取会话成员列表
func (r *conversationRepository) GetConversationMembers(conversationID string) ([]*ConversationMember, error) {
	var members []*ConversationMember
	err := r.db.Where("conversation_id = ?", conversationID).Find(&members).Error
	if err != nil {
		return nil, err
	}
	return members, nil
}

// GetConversationMember 获取会话成员
func (r *conversationRepository) GetConversationMember(conversationID string, userID int64) (*ConversationMember, error) {
	var member ConversationMember
	err := r.db.Where("conversation_id = ? AND user_id = ?", conversationID, userID).First(&member).Error
	if err != nil {
		return nil, err
	}
	return &member, nil
}

// UpdateConversationMember 更新会话成员
func (r *conversationRepository) UpdateConversationMember(member *ConversationMember) error {
	return r.db.Save(member).Error
}

// IsConversationMember 检查用户是否是会话成员
func (r *conversationRepository) IsConversationMember(conversationID string, userID int64) (bool, error) {
	var count int64
	err := r.db.Model(&ConversationMember{}).
		Where("conversation_id = ? AND user_id = ?", conversationID, userID).
		Count(&count).Error
	if err != nil {
		return false, err
	}
	return count > 0, nil
}

// CreateMessage 创建消息
func (r *conversationRepository) CreateMessage(message *Message) error {
	return r.db.Create(message).Error
}

// GetMessageByMessageID 根据消息ID获取消息
func (r *conversationRepository) GetMessageByMessageID(messageID string) (*Message, error) {
	var message Message
	err := r.db.Where("message_id = ?", messageID).First(&message).Error
	if err != nil {
		return nil, err
	}
	return &message, nil
}

// GetMessages 获取消息列表
func (r *conversationRepository) GetMessages(conversationID string, beforeTimestamp *int64, limit int) ([]*Message, error) {
	var messages []*Message
	query := r.db.Where("conversation_id = ?", conversationID)

	// 如果提供了 beforeTimestamp，则查询此时间之前的消息
	if beforeTimestamp != nil && *beforeTimestamp > 0 {
		beforeTime := time.UnixMilli(*beforeTimestamp)
		query = query.Where("created_at < ?", beforeTime)
	}

	// 按时间倒序排列，获取最新的 limit 条消息
	err := query.Order("created_at DESC").Limit(limit).Find(&messages).Error
	if err != nil {
		return nil, err
	}

	// 反转切片，使消息按时间正序排列（最早的在前）
	for i, j := 0, len(messages)-1; i < j; i, j = i+1, j-1 {
		messages[i], messages[j] = messages[j], messages[i]
	}

	return messages, nil
}

// UpdateMessage 更新消息
func (r *conversationRepository) UpdateMessage(message *Message) error {
	return r.db.Save(message).Error
}

// IncrementUnreadCount 增加未读消息数
func (r *conversationRepository) IncrementUnreadCount(conversationID string, userID int64) error {
	return r.db.Model(&ConversationMember{}).
		Where("conversation_id = ? AND user_id = ?", conversationID, userID).
		UpdateColumn("unread_count", gorm.Expr("unread_count + ?", 1)).Error
}

// ResetUnreadCount 重置未读消息数
func (r *conversationRepository) ResetUnreadCount(conversationID string, userID int64) error {
	now := time.Now()
	return r.db.Model(&ConversationMember{}).
		Where("conversation_id = ? AND user_id = ?", conversationID, userID).
		Updates(map[string]interface{}{
			"unread_count": 0,
			"last_read_at": now,
		}).Error
}
