namespace go wallet

include "common.thrift"

// 钱包信息
struct WalletInfo {
    1: i64 id
    2: i64 user_id
    3: string balance // 可用余额
    4: string frozen_balance // 冻结余额
    5: string total_recharge // 累计充值
    6: string total_consume // 累计消费
    7: string created_at
    8: string updated_at
}

// 交易类型枚举
enum TransactionType {
    RECHARGE = 1 // 充值
    CONSUME = 2 // 消费
    REFUND = 3 // 退款
    FREEZE = 4 // 冻结
    UNFREEZE = 5 // 解冻
}

// 交易状态枚举
enum TransactionStatus {
    PENDING = 1 // 待处理
    PROCESSING = 2 // 处理中
    COMPLETED = 3 // 已完成
    FAILED = 4 // 失败
    CANCELLED = 5 // 已取消
}

// 支付方式枚举
enum PaymentMethod {
    CRYPTO = 1 // 加密货币
    ONLINE = 2 // 在线支付
}

// 加密货币类型
enum CryptoCurrency {
    USDT = 1
    USDC = 2
}

// 加密货币链
enum CryptoChain {
    ETH = 1 // 以太坊
    BSC = 2 // 币安智能链
    POLYGON = 3 // Polygon
    TRON = 4 // 波场
    ARBITRUM = 5 // Arbitrum
    OPTIMISM = 6 // Optimism
}

// 在线支付平台
enum OnlinePaymentPlatform {
    STRIPE = 1
    PAYPAL = 2
}

// 交易记录
struct Transaction {
    1: i64 id
    2: string transaction_id
    3: i64 user_id
    4: string type // TransactionType转字符串
    5: string amount
    6: string balance_before
    7: string balance_after
    8: string status // TransactionStatus转字符串
    9: optional string payment_method // PaymentMethod转字符串
    10: optional string crypto_currency
    11: optional string crypto_chain
    12: optional string crypto_address
    13: optional string crypto_tx_hash
    14: optional string online_payment_platform
    15: optional string online_payment_order_id
    16: optional string online_payment_url
    17: optional string related_order_id
    18: optional string remark
    19: optional string failed_reason
    20: optional string completed_at
    21: string created_at
    22: string updated_at
}

// 获取钱包信息请求
struct GetWalletInfoReq {
    // JWT中间件会自动解析用户ID，无需传参
}

// 获取钱包信息响应
struct GetWalletInfoResp {
    1: optional WalletInfo wallet
    2: common.BaseResp base_resp
}

// 充值请求（加密货币）
struct CryptoRechargeReq {
    1: required string amount (api.body="amount") // 充值金额（美元）
    2: required string crypto_currency (api.body="crypto_currency") // USDT 或 USDC
    3: required string crypto_chain (api.body="crypto_chain") // 链类型
    4: required string crypto_address (api.body="crypto_address") // 支付地址
}

// 充值请求（在线支付）
struct OnlineRechargeReq {
    1: required string amount (api.body="amount") // 充值金额（美元）
    2: required string platform (api.body="platform") // 支付平台：stripe, paypal
}

// 充值响应
struct RechargeResp {
    1: optional Transaction transaction
    2: optional string payment_url // 在线支付的URL
    3: common.BaseResp base_resp
}

// 获取交易记录列表请求
struct GetTransactionListReq {
    1: optional string type (api.body="type") // 交易类型筛选
    2: optional string status (api.body="status") // 状态筛选
    3: optional i32 page (api.body="page")
    4: optional i32 page_size (api.body="page_size")
}

// 获取交易记录列表响应
struct GetTransactionListResp {
    1: list<Transaction> transactions
    2: i64 total
    3: i32 page
    4: i32 page_size
    5: common.BaseResp base_resp
}

// 获取交易详情请求
struct GetTransactionDetailReq {
    1: required string transaction_id (api.path="transaction_id")
}

// 获取交易详情响应
struct GetTransactionDetailResp {
    1: optional Transaction transaction
    2: common.BaseResp base_resp
}

// 确认加密货币充值请求（用户提交交易哈希后调用）
struct ConfirmCryptoRechargeReq {
    1: required string transaction_id (api.body="transaction_id")
    2: required string crypto_tx_hash (api.body="crypto_tx_hash") // 加密货币交易哈希
}

// 确认加密货币充值响应
struct ConfirmCryptoRechargeResp {
    1: optional Transaction transaction
    2: common.BaseResp base_resp
}

// 钱包服务
service WalletService {
    GetWalletInfoResp GetWalletInfo(1: GetWalletInfoReq req) (api.post="/api/v1/wallet/info")
    RechargeResp CryptoRecharge(1: CryptoRechargeReq req) (api.post="/api/v1/wallet/recharge/crypto")
    RechargeResp OnlineRecharge(1: OnlineRechargeReq req) (api.post="/api/v1/wallet/recharge/online")
    ConfirmCryptoRechargeResp ConfirmCryptoRecharge(1: ConfirmCryptoRechargeReq req) (api.post="/api/v1/wallet/recharge/crypto/confirm")
    GetTransactionListResp GetTransactionList(1: GetTransactionListReq req) (api.post="/api/v1/wallet/transactions")
    GetTransactionDetailResp GetTransactionDetail(1: GetTransactionDetailReq req) (api.post="/api/v1/wallet/transaction/:transaction_id")
}



