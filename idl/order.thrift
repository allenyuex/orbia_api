namespace go order

include "common.thrift"

// 订单信息
struct OrderInfo {
    1: i64 id
    2: string order_id  // 订单ID（格式：ORD{snowflake_id}）
    3: i64 user_id
    4: optional i64 team_id
    5: i64 kol_id
    6: string kol_display_name  // KOL显示名称（关联查询）
    7: string kol_avatar_url  // KOL头像（关联查询）
    8: i64 plan_id
    9: string plan_title
    10: string plan_description
    11: double plan_price
    12: string plan_type  // basic, standard, premium
    13: string description  // 订单描述
    14: string status  // pending, confirmed, in_progress, completed, cancelled, refunded
    15: optional string reject_reason
    16: optional string confirmed_at
    17: optional string completed_at
    18: optional string cancelled_at
    19: string created_at
    20: string updated_at
}

// 创建订单请求
struct CreateOrderReq {
    1: i64 kol_id (api.body="kol_id")
    2: i64 plan_id (api.body="plan_id")
    3: string description (api.body="description")
    4: optional i64 team_id (api.body="team_id")  // 如果是团队下单，传递团队ID
}

// 创建订单响应
struct CreateOrderResp {
    1: common.BaseResp base_resp
    2: optional string order_id
}

// 获取订单详情请求
struct GetOrderReq {
    1: string order_id (api.query="order_id")
}

// 获取订单详情响应
struct GetOrderResp {
    1: common.BaseResp base_resp
    2: optional OrderInfo order
}

// 获取订单列表请求
struct GetOrderListReq {
    1: optional string status (api.query="status")  // pending, confirmed, in_progress, completed, cancelled, refunded
    2: optional i64 kol_id (api.query="kol_id")  // 筛选指定KOL的订单
    3: optional i64 team_id (api.query="team_id")  // 筛选指定团队的订单
    4: optional i32 page (api.query="page")  // 默认1
    5: optional i32 page_size (api.query="page_size")  // 默认10
}

// 获取订单列表响应
struct GetOrderListResp {
    1: common.BaseResp base_resp
    2: list<OrderInfo> orders
    3: i64 total
}

// 更新订单状态请求
struct UpdateOrderStatusReq {
    1: string order_id (api.body="order_id")
    2: string status (api.body="status")  // confirmed, in_progress, completed, cancelled, refunded
    3: optional string reject_reason (api.body="reject_reason")  // 取消时需要提供原因
}

// 更新订单状态响应
struct UpdateOrderStatusResp {
    1: common.BaseResp base_resp
}

// 取消订单请求
struct CancelOrderReq {
    1: string order_id (api.body="order_id")
    2: string reason (api.body="reason")  // 取消原因
}

// 取消订单响应
struct CancelOrderResp {
    1: common.BaseResp base_resp
}

// 获取KOL收到的订单列表请求（KOL端使用）
struct GetKolOrdersReq {
    1: optional string status (api.query="status")
    2: optional i32 page (api.query="page")
    3: optional i32 page_size (api.query="page_size")
}

// 获取KOL收到的订单列表响应
struct GetKolOrdersResp {
    1: common.BaseResp base_resp
    2: list<OrderInfo> orders
    3: i64 total
}

// 订单服务
service OrderService {
    // 用户订单管理
    CreateOrderResp CreateOrder(1: CreateOrderReq req) (api.post="/api/v1/order/create")
    GetOrderResp GetOrder(1: GetOrderReq req) (api.post="/api/v1/order/detail")
    GetOrderListResp GetOrderList(1: GetOrderListReq req) (api.post="/api/v1/order/list")
    CancelOrderResp CancelOrder(1: CancelOrderReq req) (api.post="/api/v1/order/cancel")
    
    // KOL订单管理
    GetKolOrdersResp GetKolOrders(1: GetKolOrdersReq req) (api.post="/api/v1/order/kol/list")
    UpdateOrderStatusResp UpdateOrderStatus(1: UpdateOrderStatusReq req) (api.post="/api/v1/order/status/update")
}

