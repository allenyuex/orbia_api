namespace go user

include "common.thrift"
include "team.thrift"

// 用户信息
struct UserInfo {
    1: i64 id
    2: optional string wallet_address
    3: optional string email
    4: optional string nickname
    5: optional string avatar_url
    6: string role // 用户角色：user-普通用户，admin-管理员
    7: string created_at
    8: string updated_at
    9: optional team.Team current_team
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

// 切换当前团队请求
struct SwitchCurrentTeamReq {
    1: required i64 team_id (api.body="team_id")
}

// 切换当前团队响应
struct SwitchCurrentTeamResp {
    1: common.BaseResp base_resp
    2: optional team.Team current_team
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
    SwitchCurrentTeamResp SwitchCurrentTeam(1: SwitchCurrentTeamReq req) (api.post="/api/v1/user/switch-team")
    GetUserByIdResp GetUserById(1: GetUserByIdReq req) (api.post="/api/v1/user/:user_id")
}