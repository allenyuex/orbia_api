namespace go admin

include "common.thrift"
include "user.thrift"
include "kol.thrift"
include "team.thrift"

// ==================== 用户管理 ====================

// 管理员获取所有用户列表请求
struct GetAllUsersReq {
    1: optional string keyword (api.query="keyword") // 搜索关键字（用户名、邮箱、钱包地址）
    2: optional string role (api.query="role") // 角色筛选：user, admin
    3: optional string status (api.query="status") // 状态筛选：normal, disabled, deleted
    4: optional i32 page = 1 (api.query="page")
    5: optional i32 page_size = 10 (api.query="page_size")
}

// 用户列表项
struct UserListItem {
    1: i64 id
    2: optional string wallet_address
    3: optional string email
    4: optional string nickname
    5: optional string avatar_url
    6: string role // user, admin
    7: string status // normal, disabled, deleted
    8: optional i64 kol_id
    9: string created_at
    10: string updated_at
}

// 管理员获取所有用户列表响应
struct GetAllUsersResp {
    1: common.BaseResp base_resp
    2: list<UserListItem> users
    3: common.PageResp page_info
}

// 设置用户状态请求
struct SetUserStatusReq {
    1: i64 user_id (api.body="user_id")
    2: string status (api.body="status") // normal, disabled, deleted
}

// 设置用户状态响应
struct SetUserStatusResp {
    1: common.BaseResp base_resp
}

// ==================== KOL管理 ====================

// 管理员获取所有KOL列表请求
struct GetAllKolsReq {
    1: optional string keyword (api.query="keyword") // 搜索关键字（显示名称、国家）
    2: optional string status (api.query="status") // 状态筛选：pending, approved, rejected
    3: optional string country (api.query="country") // 国家筛选
    4: optional string tag (api.query="tag") // 标签筛选
    5: optional i32 page = 1 (api.query="page")
    6: optional i32 page_size = 10 (api.query="page_size")
}

// KOL列表项（简化版，用于列表展示）
struct KolListItem {
    1: i64 id
    2: i64 user_id
    3: string display_name
    4: string avatar_url
    5: string country
    6: string status // pending, approved, rejected
    7: optional i64 total_followers
    8: string created_at
    9: string updated_at
}

// 管理员获取所有KOL列表响应
struct GetAllKolsResp {
    1: common.BaseResp base_resp
    2: list<KolListItem> kols
    3: common.PageResp page_info
}

// 管理员审核KOL请求
struct AdminReviewKolReq {
    1: i64 kol_id (api.body="kol_id")
    2: string status (api.body="status") // approved, rejected
    3: optional string reject_reason (api.body="reject_reason")
}

// 管理员审核KOL响应
struct AdminReviewKolResp {
    1: common.BaseResp base_resp
}

// ==================== 团队管理 ====================

// 管理员获取所有团队列表请求
struct GetAllTeamsReq {
    1: optional string keyword (api.query="keyword") // 搜索关键字（团队名称）
    2: optional i32 page = 1 (api.query="page")
    3: optional i32 page_size = 10 (api.query="page_size")
}

// 团队列表项
struct TeamListItem {
    1: i64 id
    2: string name
    3: optional string icon_url
    4: i64 creator_id
    5: optional string creator_name
    6: i64 member_count
    7: string created_at
}

// 管理员获取所有团队列表响应
struct GetAllTeamsResp {
    1: common.BaseResp base_resp
    2: list<TeamListItem> teams
    3: common.PageResp page_info
}

// 管理员获取特定团队的所有用户请求
struct GetTeamMembersReq {
    1: i64 team_id (api.path="team_id")
    2: optional i32 page = 1 (api.query="page")
    3: optional i32 page_size = 10 (api.query="page_size")
}

// 团队成员项
struct TeamMemberItem {
    1: i64 user_id
    2: optional string nickname
    3: optional string email
    4: optional string avatar_url
    5: string role // creator, owner, member
    6: string joined_at
}

// 管理员获取特定团队的所有用户响应
struct GetTeamMembersResp {
    1: common.BaseResp base_resp
    2: list<TeamMemberItem> members
    3: common.PageResp page_info
}

// ==================== 订单管理 ====================

// 管理员获取所有订单列表请求
struct GetAllOrdersReq {
    1: optional string keyword (api.query="keyword") // 搜索关键字（订单ID、用户名、邮箱、钱包地址）
    2: optional string status (api.query="status") // 状态筛选
    3: optional i32 page = 1 (api.query="page")
    4: optional i32 page_size = 10 (api.query="page_size")
}

// 订单列表项
struct OrderListItem {
    1: string order_id
    2: i64 user_id
    3: optional string user_name
    4: optional string user_email
    5: i64 kol_id
    6: optional string kol_name
    7: string plan_title
    8: double plan_price
    9: string status
    10: string created_at
    11: optional string completed_at
}

// 管理员获取所有订单列表响应
struct GetAllOrdersResp {
    1: common.BaseResp base_resp
    2: list<OrderListItem> orders
    3: common.PageResp page_info
}

// ==================== 钱包管理 ====================

// 管理员获取特定用户钱包信息请求
struct GetUserWalletReq {
    1: i64 user_id (api.path="user_id")
}

// 用户钱包信息
struct UserWalletInfo {
    1: i64 user_id
    2: optional string user_name
    3: optional string user_email
    4: double balance
    5: double frozen_balance
    6: double total_recharge
    7: double total_consume
    8: string created_at
    9: string updated_at
}

// 管理员获取特定用户钱包信息响应
struct GetUserWalletResp {
    1: common.BaseResp base_resp
    2: optional UserWalletInfo wallet
}

// 管理员服务
service AdminService {
    // 用户管理
    GetAllUsersResp GetAllUsers(1: GetAllUsersReq req) (api.post="/api/v1/admin/users")
    SetUserStatusResp SetUserStatus(1: SetUserStatusReq req) (api.post="/api/v1/admin/user/status")
    
    // KOL管理
    GetAllKolsResp GetAllKols(1: GetAllKolsReq req) (api.post="/api/v1/admin/kols")
    AdminReviewKolResp AdminReviewKol(1: AdminReviewKolReq req) (api.post="/api/v1/admin/kol/review")
    
    // 团队管理
    GetAllTeamsResp GetAllTeams(1: GetAllTeamsReq req) (api.post="/api/v1/admin/teams")
    GetTeamMembersResp GetTeamMembers(1: GetTeamMembersReq req) (api.post="/api/v1/admin/team/:team_id/members")
    
    // 订单管理
    GetAllOrdersResp GetAllOrders(1: GetAllOrdersReq req) (api.post="/api/v1/admin/orders")
    
    // 钱包管理
    GetUserWalletResp GetUserWallet(1: GetUserWalletReq req) (api.post="/api/v1/admin/user/:user_id/wallet")
}

