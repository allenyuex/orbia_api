# 会话功能 API 文档

## 概述

会话功能允许用户在 KOL 订单进行过程中与 KOL 进行实时沟通。系统支持多人会话（群组），可以让客服等其他用户参与到对话中。

## 数据库设计

### 1. orbia_conversation（会话表）
- 存储会话的基本信息
- 支持多种会话类型：kol_order（KOL订单）、ad_order（广告订单）、general（普通）、support（客服）
- 关联订单信息

### 2. orbia_conversation_member（会话成员表）
- 存储会话的参与者
- 支持多人会话
- 记录未读消息数和最后阅读时间

### 3. orbia_message（消息表）
- 存储会话消息
- 支持多种消息类型：text（文本）、image（图片）、file（文件）、video（视频）、audio（音频）、system（系统消息）
- 使用毫秒时间戳支持精确的消息排序

## 自动创建会话

当 KOL 订单状态变更为 `confirmed`（已确认）时，系统会自动创建一个会话，参与者包括：
- 下单用户
- KOL 用户（使用 kol 表中的 user_id）

会话标题格式：`KOL订单: {订单标题}`

## API 接口

### 1. 发送消息

**路由**: `POST /api/v1/conversation/send_message`

**权限**: 需要 JWT 认证（普通用户和管理员）

**请求参数**:
```json
{
  "conversation_id": 123,          // 会话ID（必填）
  "message_type": "text",          // 消息类型：text, image, file, video, audio（必填）
  "content": "消息内容",            // 消息内容或文件URL（必填）
  "file_name": "example.pdf",      // 文件名（可选，文件类型消息时使用）
  "file_size": 1024,              // 文件大小（字节，可选）
  "file_type": "application/pdf"   // 文件MIME类型（可选）
}
```

**响应示例**:
```json
{
  "message": {
    "message_id": "MSG_1698765432000_12345",
    "conversation_id": 123,
    "sender_id": 1,
    "sender_nickname": "张三",
    "sender_avatar_url": "https://example.com/avatar.jpg",
    "message_type": "text",
    "content": "消息内容",
    "file_name": null,
    "file_size": null,
    "file_type": null,
    "status": "sent",
    "created_at": 1698765432000  // 毫秒时间戳
  },
  "base_resp": {
    "code": 200,
    "message": "success"
  }
}
```

### 2. 获取消息列表

**路由**: `POST /api/v1/conversation/get_messages`

**权限**: 需要 JWT 认证（普通用户和管理员）

**请求参数**:
```json
{
  "conversation_id": 123,           // 会话ID（必填）
  "before_timestamp": 1698765432000, // 毫秒时间戳，获取此时间之前的消息（可选）
  "limit": 20                        // 返回消息数量，默认20（可选）
}
```

**分页说明**:
- 首次加载：不传 `before_timestamp`，返回最新的 20 条消息
- 向上翻页：使用当前列表中最早消息的 `created_at` 作为 `before_timestamp`
- `has_more` 字段表示是否还有更多历史消息

**响应示例**:
```json
{
  "messages": [
    {
      "message_id": "MSG_1698765432000_12345",
      "conversation_id": 123,
      "sender_id": 1,
      "sender_nickname": "张三",
      "sender_avatar_url": "https://example.com/avatar.jpg",
      "message_type": "text",
      "content": "消息内容",
      "status": "sent",
      "created_at": 1698765432000
    }
  ],
  "has_more": true,  // 是否还有更多消息
  "base_resp": {
    "code": 200,
    "message": "success"
  }
}
```

### 3. 获取会话详情

**路由**: `POST /api/v1/conversation/get_conversation`

**权限**: 需要 JWT 认证（普通用户和管理员）

**请求参数**:
```json
{
  "conversation_id": 123  // 会话ID（必填）
}
```

**响应示例**:
```json
{
  "conversation": {
    "conversation_id": "CONV_1698765432000_12345",
    "title": "KOL订单: 产品推广视频",
    "type": "kol_order",
    "related_order_type": "kol_order",
    "related_order_id": "KORD_1698765432000_12345",
    "status": "active",
    "last_message_at": 1698765432000,
    "members": [
      {
        "user_id": 1,
        "nickname": "张三",
        "avatar_url": "https://example.com/avatar1.jpg",
        "role": "creator",
        "joined_at": "2024-01-01 10:00:00"
      },
      {
        "user_id": 2,
        "nickname": "李四 (KOL)",
        "avatar_url": "https://example.com/avatar2.jpg",
        "role": "member",
        "joined_at": "2024-01-01 10:00:00"
      }
    ],
    "unread_count": 5,
    "created_at": "2024-01-01 10:00:00"
  },
  "base_resp": {
    "code": 200,
    "message": "success"
  }
}
```

### 4. 获取会话列表

**路由**: `POST /api/v1/conversation/get_conversations`

**权限**: 需要 JWT 认证（普通用户和管理员）

**请求参数**:
```json
{
  "type": "kol_order",  // 筛选会话类型（可选）
  "page": 1,           // 页码，默认1（可选）
  "page_size": 20      // 每页数量，默认20（可选）
}
```

**响应示例**:
```json
{
  "conversations": [
    {
      "conversation_id": "CONV_1698765432000_12345",
      "title": "KOL订单: 产品推广视频",
      "type": "kol_order",
      "related_order_type": "kol_order",
      "related_order_id": "KORD_1698765432000_12345",
      "status": "active",
      "last_message": {
        "message_id": "MSG_1698765432000_12345",
        "sender_nickname": "张三",
        "message_type": "text",
        "content": "最后一条消息内容",
        "created_at": 1698765432000
      },
      "unread_count": 5,
      "members": [
        {
          "user_id": 1,
          "nickname": "张三",
          "avatar_url": "https://example.com/avatar1.jpg",
          "role": "creator",
          "joined_at": "2024-01-01 10:00:00"
        },
        {
          "user_id": 2,
          "nickname": "李四 (KOL)",
          "avatar_url": "https://example.com/avatar2.jpg",
          "role": "member",
          "joined_at": "2024-01-01 10:00:00"
        }
      ],
      "created_at": "2024-01-01 10:00:00",
      "last_message_at": 1698765432000
    }
  ],
  "page_info": {
    "page": 1,
    "page_size": 20,
    "total": 100,
    "total_pages": 5
  },
  "base_resp": {
    "code": 200,
    "message": "success"
  }
}
```

### 5. 标记消息已读

**路由**: `POST /api/v1/conversation/mark_read`

**权限**: 需要 JWT 认证（普通用户和管理员）

**请求参数**:
```json
{
  "conversation_id": 123  // 会话ID（必填）
}
```

**响应示例**:
```json
{
  "base_resp": {
    "code": 200,
    "message": "success"
  }
}
```

**说明**:
- 调用此接口后，会重置当前用户在该会话中的未读消息数
- 更新 `last_read_at` 时间戳

## 使用场景

### 场景1：订单沟通流程

1. 用户创建 KOL 订单
2. 用户支付订单（订单状态：pending_payment → pending）
3. KOL 确认订单（订单状态：pending → confirmed）
   - **系统自动创建会话**
4. 用户和 KOL 在会话中沟通需求细节
5. KOL 开始制作（订单状态：confirmed → in_progress）
6. 持续沟通直到完成
7. 订单完成（订单状态：in_progress → completed）

### 场景2：客服介入

当订单出现问题时，管理员可以：
1. 通过数据库直接将客服用户添加到会话成员表
2. 客服可以在会话中查看历史消息
3. 客服可以参与沟通，帮助解决问题

### 场景3：文件分享

用户或 KOL 可以在会话中分享：
- 需求文档（file 类型）
- 参考图片（image 类型）
- 视频样本（video 类型）
- 音频文件（audio 类型）

## 前端实现建议

### 1. 消息列表实现

```javascript
// 首次加载
const loadMessages = async (conversationId) => {
  const response = await api.post('/api/v1/conversation/get_messages', {
    conversation_id: conversationId,
    limit: 20
  });
  
  // 显示消息（最新的在底部）
  displayMessages(response.messages);
  
  // 如果还有更多消息，显示"加载更多"按钮
  if (response.has_more) {
    showLoadMoreButton();
  }
};

// 向上翻页加载历史消息
const loadMoreMessages = async (conversationId) => {
  const oldestMessage = messages[0]; // 当前列表中最早的消息
  
  const response = await api.post('/api/v1/conversation/get_messages', {
    conversation_id: conversationId,
    before_timestamp: oldestMessage.created_at,
    limit: 20
  });
  
  // 将新消息插入到列表顶部
  prependMessages(response.messages);
};
```

### 2. 实时性建议

由于不使用 WebSocket 或 SSE，建议：
- 使用轮询（Polling）定期获取新消息
- 建议轮询间隔：3-5 秒
- 在用户发送消息后立即刷新一次

```javascript
// 轮询获取新消息
let pollingInterval;

const startPolling = (conversationId) => {
  pollingInterval = setInterval(async () => {
    const latestMessage = messages[messages.length - 1];
    
    const response = await api.post('/api/v1/conversation/get_messages', {
      conversation_id: conversationId,
      // 不传 before_timestamp，获取最新消息
      limit: 20
    });
    
    // 找出新消息并添加到列表
    const newMessages = response.messages.filter(msg => 
      msg.created_at > latestMessage.created_at
    );
    
    if (newMessages.length > 0) {
      appendMessages(newMessages);
      playNotificationSound();
    }
  }, 5000); // 5秒轮询一次
};

const stopPolling = () => {
  clearInterval(pollingInterval);
};
```

### 3. 未读消息提示

```javascript
// 在会话列表页面显示未读数
const renderConversationList = (conversations) => {
  conversations.forEach(conv => {
    if (conv.unread_count > 0) {
      showUnreadBadge(conv.conversation_id, conv.unread_count);
    }
  });
};

// 进入会话时标记已读
const enterConversation = async (conversationId) => {
  await api.post('/api/v1/conversation/mark_read', {
    conversation_id: conversationId
  });
  
  clearUnreadBadge(conversationId);
  startPolling(conversationId);
};

// 离开会话时停止轮询
const leaveConversation = () => {
  stopPolling();
};
```

## 注意事项

1. **权限验证**: 所有接口都会验证用户是否是会话成员
2. **消息顺序**: 使用毫秒时间戳确保消息顺序准确
3. **文件上传**: 需要先使用上传接口获取文件 URL，然后在消息中发送 URL
4. **会话状态**: 
   - `active`: 活跃中，可以正常收发消息
   - `archived`: 已归档，仍可查看但不建议继续发送
   - `closed`: 已关闭，不能发送新消息
5. **性能优化**: 
   - 消息列表默认返回 20 条
   - 建议前端实现虚拟滚动以处理大量历史消息
   - 轮询间隔不要设置太短，避免服务器压力

## 错误处理

常见错误码：
- `400`: 参数错误或业务逻辑错误
- `401`: 未认证或认证失败
- `403`: 无权访问（不是会话成员）
- `404`: 会话或消息不存在

示例错误响应：
```json
{
  "base_resp": {
    "code": 400,
    "message": "user is not a member of this conversation"
  }
}
```

