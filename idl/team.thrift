namespace go team

include "common.thrift"

// 团队角色枚举
enum TeamRole {
    CREATOR = 1,
    OWNER = 2,
    MEMBER = 3
}

// 邀请状态枚举
enum InvitationStatus {
    PENDING = 1,
    ACCEPTED = 2,
    REJECTED = 3,
    EXPIRED = 4
}

// 团队信息
struct Team {
    1: i64 id
    2: string name
    3: optional string icon_url
    4: i64 creator_id
    5: string created_at
    6: string updated_at
}

// 团队成员信息
struct TeamMember {
    1: i64 id
    2: i64 team_id
    3: i64 user_id
    4: TeamRole role
    5: string joined_at
    6: optional string user_nickname
    7: optional string user_avatar_url
    8: optional string user_email
    9: optional string user_wallet_address
}

// 团队邀请信息
struct TeamInvitation {
    1: i64 id
    2: i64 team_id
    3: i64 inviter_id
    4: optional string invitee_email
    5: optional string invitee_wallet
    6: TeamRole role
    7: InvitationStatus status
    8: string invitation_code
    9: string expires_at
    10: string created_at
    11: optional string team_name
    12: optional string inviter_nickname
}

// 创建团队请求
struct CreateTeamReq {
    1: string name (api.body="name")
    2: optional string icon_url (api.body="icon_url")
}

// 创建团队响应
struct CreateTeamResp {
    1: Team team
    2: common.BaseResp base_resp
}

// 获取团队详情请求
struct GetTeamReq {
    1: string team_id (api.body="team_id")
}

// 获取团队详情响应
struct GetTeamResp {
    1: Team team
    2: common.BaseResp base_resp
}

// 更新团队请求
struct UpdateTeamReq {
    1: string team_id (api.body="team_id")
    2: optional string name (api.body="name")
    3: optional string icon_url (api.body="icon_url")
}

// 更新团队响应
struct UpdateTeamResp {
    1: Team team
    2: common.BaseResp base_resp
}

// 获取用户团队列表请求
struct GetUserTeamsReq {
    1: optional common.PageReq page_req
}

// 获取用户团队列表响应
struct GetUserTeamsResp {
    1: list<Team> teams
    2: optional common.PageResp page_resp
    3: common.BaseResp base_resp
}

// 邀请用户加入团队请求
struct InviteUserReq {
    1: string team_id (api.body="team_id")
    2: optional string email (api.body="email")
    3: optional string wallet_address (api.body="wallet_address")
    4: TeamRole role (api.body="role")
}

// 邀请用户加入团队响应
struct InviteUserResp {
    1: TeamInvitation invitation
    2: common.BaseResp base_resp
}

// 获取团队成员列表请求
struct GetTeamMembersReq {
    1: string team_id (api.body="team_id")
    2: optional common.PageReq page_req
}

// 获取团队成员列表响应
struct GetTeamMembersResp {
    1: list<TeamMember> members
    2: optional common.PageResp page_resp
    3: common.BaseResp base_resp
}

// 移除团队成员请求
struct RemoveTeamMemberReq {
    1: string team_id (api.body="team_id")
    2: string user_id (api.body="user_id")
}

// 移除团队成员响应
struct RemoveTeamMemberResp {
    1: common.BaseResp base_resp
}

// 接受邀请请求
struct AcceptInvitationReq {
    1: string invitation_code (api.body="invitation_code")
}

// 接受邀请响应
struct AcceptInvitationResp {
    1: TeamMember member
    2: common.BaseResp base_resp
}

// 拒绝邀请请求
struct RejectInvitationReq {
    1: string invitation_code (api.body="invitation_code")
}

// 拒绝邀请响应
struct RejectInvitationResp {
    1: common.BaseResp base_resp
}

// 获取用户邀请列表请求
struct GetUserInvitationsReq {
    1: optional common.PageReq page_req
}

// 获取用户邀请列表响应
struct GetUserInvitationsResp {
    1: list<TeamInvitation> invitations
    2: optional common.PageResp page_resp
    3: common.BaseResp base_resp
}

// 团队服务接口
service TeamService {
    // 创建团队
    CreateTeamResp CreateTeam(1: CreateTeamReq req) (api.post="/api/v1/team/create")
    
    // 获取团队详情
    GetTeamResp GetTeam(1: GetTeamReq req) (api.post="/api/v1/team/get")
    
    // 更新团队
    UpdateTeamResp UpdateTeam(1: UpdateTeamReq req) (api.post="/api/v1/team/update")
    
    // 获取用户团队列表
    GetUserTeamsResp GetUserTeams(1: GetUserTeamsReq req) (api.post="/api/v1/team/user-teams")
    
    // 邀请用户加入团队
    InviteUserResp InviteUser(1: InviteUserReq req) (api.post="/api/v1/team/invite")
    
    // 获取团队成员列表
    GetTeamMembersResp GetTeamMembers(1: GetTeamMembersReq req) (api.post="/api/v1/team/members")
    
    // 移除团队成员
    RemoveTeamMemberResp RemoveTeamMember(1: RemoveTeamMemberReq req) (api.post="/api/v1/team/remove-member")
    
    // 接受邀请
    AcceptInvitationResp AcceptInvitation(1: AcceptInvitationReq req) (api.post="/api/v1/team/invitation/accept")
    
    // 拒绝邀请
    RejectInvitationResp RejectInvitation(1: RejectInvitationReq req) (api.post="/api/v1/team/invitation/reject")
    
    // 获取用户邀请列表
    GetUserInvitationsResp GetUserInvitations(1: GetUserInvitationsReq req) (api.post="/api/v1/team/invitations")
}