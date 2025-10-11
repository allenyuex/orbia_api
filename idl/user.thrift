namespace go user

include "common.thrift"

// 用户信息
struct UserInfo {
    1: i64 id
    2: optional string wallet_address
    3: optional string email
    4: optional string nickname
    5: optional string avatar_url
    6: string created_at
    7: string updated_at
}

// 获取用户信息请求
struct GetProfileReq {
    // JWT中间件会自动解析用户ID，无需传参
}

// 获取用户信息响应
struct GetProfileResp {
    1: optional UserInfo user
    2: common.BaseResp base_resp
}

// 更新用户信息请求
struct UpdateProfileReq {
    1: optional string nickname (api.body="nickname")
    2: optional string avatar_url (api.body="avatar_url")
}

// 更新用户信息响应
struct UpdateProfileResp {
    1: common.BaseResp base_resp
}

// 根据ID获取用户请求
struct GetUserByIdReq {
    1: required i64 user_id (api.path="user_id")
}

// 根据ID获取用户响应
struct GetUserByIdResp {
    1: optional UserInfo user
    2: common.BaseResp base_resp
}

// 用户服务
service UserService {
    GetProfileResp GetProfile(1: GetProfileReq req) (api.post="/api/v1/user/profile")
    UpdateProfileResp UpdateProfile(1: UpdateProfileReq req) (api.post="/api/v1/user/update-profile")
    GetUserByIdResp GetUserById(1: GetUserByIdReq req) (api.post="/api/v1/user/:user_id")
}