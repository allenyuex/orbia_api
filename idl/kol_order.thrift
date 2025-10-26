namespace go kol_order

include "common.thrift"

// KOL订单信息
struct KolOrderInfo {
    1: string order_id  // 订单ID（字符串，格式：KORD_{timestamp}_{random}）
    2: i64 user_id
    3: string user_nickname  // 用户昵称（关联查询）
    4: optional i64 team_id
    5: optional string team_name  // 团队名称（关联查询）
    6: i64 kol_id
    7: string kol_display_name  // KOL显示名称（关联查询）
    8: string kol_avatar_url  // KOL头像（关联查询）
    9: i64 plan_id
    10: string plan_title  // Plan标题（快照）
    11: string plan_description  // Plan描述（快照）
    12: double plan_price  // Plan价格（快照，美元）
    13: string plan_type  // Plan类型（快照）：basic, standard, premium
    14: string title  // 订单标题
    15: string requirement_description  // 合作需求描述
    16: string video_type  // 视频类型（用户手动输入）
    17: i32 video_duration  // 视频预计时长（秒数）
    18: string target_audience  // 目标受众
    19: string expected_delivery_date  // 期望交付日期（YYYY-MM-DD）
    20: optional string additional_requirements  // 额外要求
    21: string status  // pending_payment-待支付, pending-待确认, confirmed-已确认, in_progress-进行中, completed-已完成, cancelled-已取消, refunded-已退款
    22: optional string reject_reason  // 拒绝/取消原因
    23: optional string confirmed_at
    24: optional string completed_at
    25: optional string cancelled_at
    26: string created_at
    27: string updated_at
}

// 创建KOL订单请求
struct CreateKolOrderReq {
    1: i64 kol_id (api.body="kol_id")
    2: i64 plan_id (api.body="plan_id")
    3: string title (api.body="title")  // 订单标题
    4: string requirement_description (api.body="requirement_description")  // 合作需求描述
    5: string video_type (api.body="video_type")  // 视频类型
    6: i32 video_duration (api.body="video_duration")  // 视频预计时长（秒数）
    7: string target_audience (api.body="target_audience")  // 目标受众
    8: string expected_delivery_date (api.body="expected_delivery_date")  // 期望交付日期（YYYY-MM-DD）
    9: optional string additional_requirements (api.body="additional_requirements")  // 额外要求
    10: optional i64 team_id (api.body="team_id")  // 如果是团队下单，传递团队ID
}

// 创建KOL订单响应
struct CreateKolOrderResp {
    1: common.BaseResp base_resp
    2: optional string order_id
}

// 获取KOL订单详情请求
struct GetKolOrderReq {
    1: string order_id (api.body="order_id")
}

// 获取KOL订单详情响应
struct GetKolOrderResp {
    1: common.BaseResp base_resp
    2: optional KolOrderInfo order
}

// 获取用户自己的KOL订单列表请求
struct GetUserKolOrderListReq {
    1: optional string status (api.body="status")  // 订单状态筛选
    2: optional string keyword (api.body="keyword")  // 模糊搜索关键词（搜索标题、KOL名称等）
    3: optional i64 kol_id (api.body="kol_id")  // 筛选指定KOL的订单
    4: optional i64 team_id (api.body="team_id")  // 筛选指定团队的订单
    5: optional i32 page (api.body="page")  // 默认1
    6: optional i32 page_size (api.body="page_size")  // 默认10
}

// 获取用户自己的KOL订单列表响应
struct GetUserKolOrderListResp {
    1: common.BaseResp base_resp
    2: list<KolOrderInfo> orders
    3: i64 total
}

// 获取KOL收到的订单列表请求（KOL端使用）
struct GetKolReceivedOrderListReq {
    1: optional string status (api.body="status")  // 订单状态筛选
    2: optional string keyword (api.body="keyword")  // 模糊搜索关键词（搜索标题、用户名称等）
    3: optional i32 page (api.body="page")  // 默认1
    4: optional i32 page_size (api.body="page_size")  // 默认10
}

// 获取KOL收到的订单列表响应
struct GetKolReceivedOrderListResp {
    1: common.BaseResp base_resp
    2: list<KolOrderInfo> orders
    3: i64 total
}

// 更新KOL订单状态请求（KOL使用）
struct UpdateKolOrderStatusReq {
    1: string order_id (api.body="order_id")
    2: string status (api.body="status")  // pending, confirmed, in_progress, completed, cancelled, refunded
    3: optional string reject_reason (api.body="reject_reason")  // 取消时需要提供原因
}

// 更新KOL订单状态响应
struct UpdateKolOrderStatusResp {
    1: common.BaseResp base_resp
}

// 取消KOL订单请求（用户使用）
struct CancelKolOrderReq {
    1: string order_id (api.body="order_id")
    2: string reason (api.body="reason")  // 取消原因
}

// 取消KOL订单响应
struct CancelKolOrderResp {
    1: common.BaseResp base_resp
}

// 确认KOL订单支付请求（用户支付完成后调用）
struct ConfirmKolOrderPaymentReq {
    1: string order_id (api.body="order_id")
}

// 确认KOL订单支付响应
struct ConfirmKolOrderPaymentResp {
    1: common.BaseResp base_resp
}

// KOL订单服务
service KolOrderService {
    // 用户订单管理
    CreateKolOrderResp CreateKolOrder(1: CreateKolOrderReq req) (api.post="/api/v1/kol-order/create")
    GetKolOrderResp GetKolOrder(1: GetKolOrderReq req) (api.post="/api/v1/kol-order/detail")
    GetUserKolOrderListResp GetUserKolOrderList(1: GetUserKolOrderListReq req) (api.post="/api/v1/kol-order/user/list")
    CancelKolOrderResp CancelKolOrder(1: CancelKolOrderReq req) (api.post="/api/v1/kol-order/cancel")
    ConfirmKolOrderPaymentResp ConfirmKolOrderPayment(1: ConfirmKolOrderPaymentReq req) (api.post="/api/v1/kol-order/payment/confirm")
    
    // KOL订单管理
    GetKolReceivedOrderListResp GetKolReceivedOrderList(1: GetKolReceivedOrderListReq req) (api.post="/api/v1/kol-order/kol/list")
    UpdateKolOrderStatusResp UpdateKolOrderStatus(1: UpdateKolOrderStatusReq req) (api.post="/api/v1/kol-order/status/update")
}

