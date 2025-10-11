namespace go auth

include "common.thrift"

// 钱包登录请求
struct WalletLoginReq {
    1: required string wallet_address (api.body="wallet_address")
    2: required string signature (api.body="signature")
    3: optional string message (api.body="message")
}

// 钱包登录响应
struct WalletLoginResp {
    1: string token
    2: i64 expires_in
    3: common.BaseResp base_resp
}

// 邮箱登录请求（预留）
struct EmailLoginReq {
    1: required string email (api.body="email")
    2: required string password (api.body="password")
}

// 邮箱登录响应（预留）
struct EmailLoginResp {
    1: string token
    2: i64 expires_in
    3: common.BaseResp base_resp
}

// 发送验证码请求（预留）
struct SendCodeReq {
    1: required string email (api.body="email")
}

// 发送验证码响应（预留）
struct SendCodeResp {
    1: common.BaseResp base_resp
}

// 认证服务
service AuthService {
    WalletLoginResp WalletLogin(1: WalletLoginReq req) (api.post="/api/v1/auth/wallet-login")
    EmailLoginResp EmailLogin(1: EmailLoginReq req) (api.post="/api/v1/auth/email-login")
    SendCodeResp SendCode(1: SendCodeReq req) (api.post="/api/v1/auth/send-code")
}