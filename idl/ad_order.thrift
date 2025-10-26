namespace go ad_order

include "common.thrift"

// 广告订单信息
struct AdOrderInfo {
    1: string order_id  // 订单ID（字符串，格式：ADORD_{timestamp}_{random}）
    2: i64 user_id
    3: string user_nickname  // 用户昵称（关联查询）
    4: optional i64 team_id
    5: optional string team_name  // 团队名称（关联查询）
    6: string title  // 广告订单标题
    7: string description  // 广告订单描述
    8: double budget  // 广告预算（美元）
    9: string ad_type  // 广告类型：banner, video, social_media, influencer
    10: string target_audience  // 目标受众
    11: string start_date  // 开始日期（YYYY-MM-DD）
    12: string end_date  // 结束日期（YYYY-MM-DD）
    13: string status  // pending-待审核, approved-已批准, in_progress-进行中, completed-已完成, cancelled-已取消
    14: optional string reject_reason  // 拒绝/取消原因
    15: optional string approved_at
    16: optional string completed_at
    17: optional string cancelled_at
    18: string created_at
    19: string updated_at
}

// 创建广告订单请求
struct CreateAdOrderReq {
    1: string title (api.body="title")  // 广告订单标题
    2: string description (api.body="description")  // 广告订单描述
    3: double budget (api.body="budget")  // 广告预算（美元）
    4: string ad_type (api.body="ad_type")  // 广告类型
    5: string target_audience (api.body="target_audience")  // 目标受众
    6: string start_date (api.body="start_date")  // 开始日期（YYYY-MM-DD）
    7: string end_date (api.body="end_date")  // 结束日期（YYYY-MM-DD）
    8: optional i64 team_id (api.body="team_id")  // 如果是团队下单，传递团队ID
}

// 创建广告订单响应
struct CreateAdOrderResp {
    1: common.BaseResp base_resp
    2: optional string order_id
}

// 获取广告订单详情请求
struct GetAdOrderReq {
    1: string order_id (api.body="order_id")
}

// 获取广告订单详情响应
struct GetAdOrderResp {
    1: common.BaseResp base_resp
    2: optional AdOrderInfo order
}

// 获取用户自己的广告订单列表请求
struct GetUserAdOrderListReq {
    1: optional string status (api.body="status")  // 订单状态筛选
    2: optional string keyword (api.body="keyword")  // 模糊搜索关键词（搜索标题、描述等）
    3: optional string ad_type (api.body="ad_type")  // 广告类型筛选
    4: optional i64 team_id (api.body="team_id")  // 筛选指定团队的订单
    5: optional i32 page (api.body="page")  // 默认1
    6: optional i32 page_size (api.body="page_size")  // 默认10
}

// 获取用户自己的广告订单列表响应
struct GetUserAdOrderListResp {
    1: common.BaseResp base_resp
    2: list<AdOrderInfo> orders
    3: i64 total
}

// 获取所有广告订单列表请求（管理员使用）
struct GetAdOrderListReq {
    1: optional string status (api.body="status")  // 订单状态筛选
    2: optional string keyword (api.body="keyword")  // 模糊搜索关键词
    3: optional string ad_type (api.body="ad_type")  // 广告类型筛选
    4: optional i32 page (api.body="page")  // 默认1
    5: optional i32 page_size (api.body="page_size")  // 默认10
}

// 获取所有广告订单列表响应
struct GetAdOrderListResp {
    1: common.BaseResp base_resp
    2: list<AdOrderInfo> orders
    3: i64 total
}

// 更新广告订单状态请求（管理员使用）
struct UpdateAdOrderStatusReq {
    1: string order_id (api.body="order_id")
    2: string status (api.body="status")  // approved, in_progress, completed, cancelled
    3: optional string reject_reason (api.body="reject_reason")  // 拒绝时需要提供原因
}

// 更新广告订单状态响应
struct UpdateAdOrderStatusResp {
    1: common.BaseResp base_resp
}

// 取消广告订单请求（用户使用）
struct CancelAdOrderReq {
    1: string order_id (api.body="order_id")
    2: string reason (api.body="reason")  // 取消原因
}

// 取消广告订单响应
struct CancelAdOrderResp {
    1: common.BaseResp base_resp
}

// 广告订单服务
service AdOrderService {
    // 用户订单管理
    CreateAdOrderResp CreateAdOrder(1: CreateAdOrderReq req) (api.post="/api/v1/ad-order/create")
    GetAdOrderResp GetAdOrder(1: GetAdOrderReq req) (api.post="/api/v1/ad-order/detail")
    GetUserAdOrderListResp GetUserAdOrderList(1: GetUserAdOrderListReq req) (api.post="/api/v1/ad-order/user/list")
    CancelAdOrderResp CancelAdOrder(1: CancelAdOrderReq req) (api.post="/api/v1/ad-order/cancel")
    
    // 管理员订单管理
    GetAdOrderListResp GetAdOrderList(1: GetAdOrderListReq req) (api.post="/api/v1/ad-order/admin/list")
    UpdateAdOrderStatusResp UpdateAdOrderStatus(1: UpdateAdOrderStatusReq req) (api.post="/api/v1/ad-order/status/update")
}

