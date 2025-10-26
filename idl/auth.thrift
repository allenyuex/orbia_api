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

// 发送验证码请求
struct SendVerificationCodeReq {
    1: required string email (api.body="email")
    2: optional string code_type (api.body="code_type")  // 验证码类型：login, register, reset_password，默认为login
}

// 发送验证码响应
struct SendVerificationCodeResp {
    1: common.BaseResp base_resp
}

// 邮箱验证码登录请求
struct EmailLoginReq {
    1: required string email (api.body="email")
    2: required string code (api.body="code")
}

// 邮箱验证码登录响应
struct EmailLoginResp {
    1: string token
    2: i64 expires_in
    3: common.BaseResp base_resp
}

// 认证服务
service AuthService {
    WalletLoginResp WalletLogin(1: WalletLoginReq req) (api.post="/api/v1/auth/wallet-login")
    SendVerificationCodeResp SendVerificationCode(1: SendVerificationCodeReq req) (api.post="/api/v1/auth/send-verification-code")
    EmailLoginResp EmailLogin(1: EmailLoginReq req) (api.post="/api/v1/auth/email-login")
}