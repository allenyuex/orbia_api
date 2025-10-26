namespace go recharge_order

include "common.thrift"

// 支付类型枚举
enum PaymentType {
    CRYPTO = 1 // 加密货币
    ONLINE = 2 // 在线支付
}

// 充值订单状态枚举
enum RechargeOrderStatus {
    PENDING = 1 // 待确认
    CONFIRMED = 2 // 已确认
    FAILED = 3 // 失败
    CANCELLED = 4 // 已取消
}

// 充值订单信息
struct RechargeOrder {
    1: i64 id
    2: string order_id
    3: i64 user_id
    4: string amount // 充值金额（美元）
    5: string payment_type // PaymentType转字符串：crypto, online
    6: optional i64 payment_setting_id
    7: optional string payment_network // 快照-区块链网络
    8: optional string payment_address // 快照-钱包地址
    9: optional string payment_label // 快照-钱包标签
    10: optional string user_crypto_address // 用户的转账钱包地址
    11: optional string crypto_tx_hash // 加密货币交易哈希
    12: optional string online_payment_platform // 在线支付平台：stripe, paypal
    13: optional string online_payment_order_id // 在线支付平台订单ID
    14: optional string online_payment_url // 在线支付URL
    15: string status // RechargeOrderStatus转字符串：pending, confirmed, failed, cancelled
    16: optional i64 confirmed_by // 确认人ID（管理员）
    17: optional string confirmed_at // 确认时间
    18: optional string failed_reason // 失败原因
    19: optional string remark // 备注
    20: string created_at
    21: string updated_at
}

// 创建充值订单请求（加密货币）
struct CreateCryptoRechargeOrderReq {
    1: required string amount (api.body="amount") // 充值金额（美元）
    2: required i64 payment_setting_id (api.body="payment_setting_id") // 选择的收款钱包ID
    3: required string user_crypto_address (api.body="user_crypto_address") // 用户的转账钱包地址
    4: optional string crypto_tx_hash (api.body="crypto_tx_hash") // 加密货币交易哈希（可选）
    5: optional string remark (api.body="remark") // 备注
}

// 创建充值订单响应
struct CreateRechargeOrderResp {
    1: optional RechargeOrder order
    2: common.BaseResp base_resp
}

// 创建充值订单请求（在线支付）
struct CreateOnlineRechargeOrderReq {
    1: required string amount (api.body="amount") // 充值金额（美元）
    2: required string platform (api.body="platform") // 支付平台：stripe, paypal
}

// 查询充值订单列表请求（normal用户）
struct GetMyRechargeOrdersReq {
    1: optional string status (api.body="status") // 状态筛选：pending, confirmed, failed, cancelled
    2: optional i32 page (api.body="page")
    3: optional i32 page_size (api.body="page_size")
}

// 查询充值订单列表响应
struct GetRechargeOrdersResp {
    1: list<RechargeOrder> orders
    2: i64 total
    3: i32 page
    4: i32 page_size
    5: common.BaseResp base_resp
}

// 查询所有充值订单列表请求（admin用户）
struct GetAllRechargeOrdersReq {
    1: optional string status (api.body="status") // 状态筛选
    2: optional string payment_type (api.body="payment_type") // 支付类型筛选：crypto, online
    3: optional i64 user_id (api.body="user_id") // 用户ID筛选
    4: optional i32 page (api.body="page")
    5: optional i32 page_size (api.body="page_size")
}

// 查询充值订单详情请求
struct GetRechargeOrderDetailReq {
    1: required string order_id (api.path="order_id")
}

// 查询充值订单详情响应
struct GetRechargeOrderDetailResp {
    1: optional RechargeOrder order
    2: common.BaseResp base_resp
}

// admin确认充值订单请求
struct ConfirmRechargeOrderReq {
    1: required string order_id (api.body="order_id")
    2: optional string crypto_tx_hash (api.body="crypto_tx_hash") // 管理员可以填写或更新交易哈希
    3: optional string remark (api.body="remark") // 备注
}

// admin确认充值订单响应
struct ConfirmRechargeOrderResp {
    1: optional RechargeOrder order
    2: common.BaseResp base_resp
}

// admin拒绝充值订单请求
struct RejectRechargeOrderReq {
    1: required string order_id (api.body="order_id")
    2: required string failed_reason (api.body="failed_reason") // 失败原因
}

// admin拒绝充值订单响应
struct RejectRechargeOrderResp {
    1: optional RechargeOrder order
    2: common.BaseResp base_resp
}

// 充值订单服务
service RechargeOrderService {
    // normal用户创建充值订单
    CreateRechargeOrderResp CreateCryptoRechargeOrder(1: CreateCryptoRechargeOrderReq req) (api.post="/api/v1/recharge/create/crypto")
    CreateRechargeOrderResp CreateOnlineRechargeOrder(1: CreateOnlineRechargeOrderReq req) (api.post="/api/v1/recharge/create/online")
    
    // normal用户查询自己的充值订单
    GetRechargeOrdersResp GetMyRechargeOrders(1: GetMyRechargeOrdersReq req) (api.post="/api/v1/recharge/my/list")
    GetRechargeOrderDetailResp GetRechargeOrderDetail(1: GetRechargeOrderDetailReq req) (api.post="/api/v1/recharge/detail/:order_id")
    
    // admin查询所有充值订单
    GetRechargeOrdersResp GetAllRechargeOrders(1: GetAllRechargeOrdersReq req) (api.post="/api/v1/admin/recharge/list")
    
    // admin确认/拒绝充值订单
    ConfirmRechargeOrderResp ConfirmRechargeOrder(1: ConfirmRechargeOrderReq req) (api.post="/api/v1/admin/recharge/confirm")
    RejectRechargeOrderResp RejectRechargeOrder(1: RejectRechargeOrderReq req) (api.post="/api/v1/admin/recharge/reject")
}

