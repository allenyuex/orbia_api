namespace go conversation

include "common.thrift"

// 会话成员信息
struct ConversationMember {
    1: i64 user_id
    2: string nickname
    3: optional string avatar_url
    4: string role  // creator, member, admin
    5: string joined_at
}

// 消息信息
struct Message {
    1: string message_id
    2: string conversation_id
    3: i64 sender_id
    4: string sender_nickname
    5: optional string sender_avatar_url
    6: string message_type  // text, image, file, video, audio, system
    7: string content
    8: optional string file_name
    9: optional i64 file_size
    10: optional string file_type
    11: string status  // sent, delivered, read, failed
    12: i64 created_at  // 毫秒时间戳
}

// 会话详情
struct ConversationInfo {
    1: string conversation_id
    2: optional string title
    3: string type  // kol_order, ad_order, general, support
    4: optional string related_order_type
    5: optional string related_order_id
    6: string status  // active, archived, closed
    7: optional i64 last_message_at
    8: list<ConversationMember> members
    9: i32 unread_count
    10: string created_at
}

// 会话列表项
struct ConversationItem {
    1: string conversation_id
    2: optional string title
    3: string type
    4: optional string related_order_type
    5: optional string related_order_id
    6: string status
    7: optional Message last_message
    8: i32 unread_count
    9: list<ConversationMember> members
    10: string created_at
    11: optional i64 last_message_at
}

// 发送消息请求
struct SendMessageReq {
    1: string conversation_id (api.body="conversation_id")
    2: string message_type (api.body="message_type")  // text, image, file, video, audio
    3: string content (api.body="content")
    4: optional string file_name (api.body="file_name")
    5: optional i64 file_size (api.body="file_size")
    6: optional string file_type (api.body="file_type")
}

// 发送消息响应
struct SendMessageResp {
    1: Message message
    2: common.BaseResp base_resp
}

// 获取消息列表请求
struct GetMessagesReq {
    1: string conversation_id (api.body="conversation_id")
    2: optional i64 before_timestamp (api.body="before_timestamp")  // 毫秒时间戳，获取此时间之前的消息
    3: optional i32 limit = 20 (api.body="limit")  // 默认返回20条
}

// 获取消息列表响应
struct GetMessagesResp {
    1: list<Message> messages
    2: bool has_more  // 是否还有更多消息
    3: common.BaseResp base_resp
}

// 获取会话详情请求
struct GetConversationReq {
    1: string conversation_id (api.body="conversation_id")
}

// 获取会话详情响应
struct GetConversationResp {
    1: ConversationInfo conversation
    2: common.BaseResp base_resp
}

// 获取会话列表请求
struct GetConversationsReq {
    1: optional string type (api.body="type")  // 筛选会话类型
    2: optional i32 page = 1 (api.body="page")
    3: optional i32 page_size = 20 (api.body="page_size")
}

// 获取会话列表响应
struct GetConversationsResp {
    1: list<ConversationItem> conversations
    2: common.PageResp page_info
    3: common.BaseResp base_resp
}

// 标记消息已读请求
struct MarkMessagesReadReq {
    1: string conversation_id (api.body="conversation_id")
}

// 标记消息已读响应
struct MarkMessagesReadResp {
    1: common.BaseResp base_resp
}

// 会话服务定义
service ConversationService {
    // 发送消息
    SendMessageResp SendMessage(1: SendMessageReq req) (api.post="/api/v1/conversation/send_message")
    
    // 获取消息列表
    GetMessagesResp GetMessages(1: GetMessagesReq req) (api.post="/api/v1/conversation/get_messages")
    
    // 获取会话详情
    GetConversationResp GetConversation(1: GetConversationReq req) (api.post="/api/v1/conversation/get_conversation")
    
    // 获取会话列表
    GetConversationsResp GetConversations(1: GetConversationsReq req) (api.post="/api/v1/conversation/get_conversations")
    
    // 标记消息已读
    MarkMessagesReadResp MarkMessagesRead(1: MarkMessagesReadReq req) (api.post="/api/v1/conversation/mark_read")
}

