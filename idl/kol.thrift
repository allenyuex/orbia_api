namespace go kol

include "common.thrift"

// KOL语言信息
struct KolLanguage {
    1: string language_code
    2: string language_name
}

// KOL标签
struct KolTag {
    1: string tag
}

// KOL数据统计
struct KolStats {
    1: i64 total_followers
    2: i64 tiktok_followers
    3: i64 youtube_subscribers
    4: i64 x_followers
    5: i64 discord_members
    6: i64 tiktok_avg_views
    7: double engagement_rate
}

// KOL报价Plan
struct KolPlan {
    1: i64 id
    2: string title
    3: string description
    4: double price
    5: string plan_type  // basic, standard, premium
    6: string created_at
    7: string updated_at
}

// KOL视频
struct KolVideo {
    1: i64 id
    2: string embed_code
    3: optional string cover_url
    4: string created_at
    5: string updated_at
}

// KOL详细信息
struct KolInfo {
    1: i64 id
    2: i64 user_id
    3: string avatar_url
    4: string display_name
    5: string description
    6: string country
    7: string tiktok_url
    8: string youtube_url
    9: string x_url
    10: string discord_url
    11: string status  // pending, approved, rejected
    12: string reject_reason
    13: string approved_at
    14: list<KolLanguage> languages
    15: list<KolTag> tags
    16: KolStats stats
    17: string created_at
    18: string updated_at
}

// 申请成为KOL请求
struct ApplyKolReq {
    1: string display_name (api.body="display_name")
    2: string description (api.body="description")
    3: string country (api.body="country")
    4: optional string avatar_url (api.body="avatar_url")
    5: optional string tiktok_url (api.body="tiktok_url")
    6: optional string youtube_url (api.body="youtube_url")
    7: optional string x_url (api.body="x_url")
    8: optional string discord_url (api.body="discord_url")
    9: list<string> language_codes (api.body="language_codes")  // ["en", "zh"]
    10: list<string> language_names (api.body="language_names")  // ["English", "中文"]
    11: list<string> tags (api.body="tags")  // ["Defi", "Web3"]
}

// 申请成为KOL响应
struct ApplyKolResp {
    1: common.BaseResp base_resp
    2: optional i64 kol_id
}

// 获取KOL信息请求
struct GetKolInfoReq {
    1: optional i64 kol_id (api.query="kol_id")  // 不传则获取当前登录用户的KOL信息
}

// 获取KOL信息响应
struct GetKolInfoResp {
    1: common.BaseResp base_resp
    2: optional KolInfo kol_info
}

// 更新KOL信息请求
struct UpdateKolInfoReq {
    1: optional string display_name (api.body="display_name")
    2: optional string description (api.body="description")
    3: optional string country (api.body="country")
    4: optional string avatar_url (api.body="avatar_url")
    5: optional string tiktok_url (api.body="tiktok_url")
    6: optional string youtube_url (api.body="youtube_url")
    7: optional string x_url (api.body="x_url")
    8: optional string discord_url (api.body="discord_url")
    9: optional list<string> language_codes (api.body="language_codes")
    10: optional list<string> language_names (api.body="language_names")
    11: optional list<string> tags (api.body="tags")
}

// 更新KOL信息响应
struct UpdateKolInfoResp {
    1: common.BaseResp base_resp
}

// 审核KOL请求（管理员使用）
struct ReviewKolReq {
    1: i64 kol_id (api.body="kol_id")
    2: string status (api.body="status")  // approved, rejected
    3: optional string reject_reason (api.body="reject_reason")
}

// 审核KOL响应
struct ReviewKolResp {
    1: common.BaseResp base_resp
}

// 更新KOL统计数据请求
struct UpdateKolStatsReq {
    1: optional i64 total_followers (api.body="total_followers")
    2: optional i64 tiktok_followers (api.body="tiktok_followers")
    3: optional i64 youtube_subscribers (api.body="youtube_subscribers")
    4: optional i64 x_followers (api.body="x_followers")
    5: optional i64 discord_members (api.body="discord_members")
    6: optional i64 tiktok_avg_views (api.body="tiktok_avg_views")
    7: optional double engagement_rate (api.body="engagement_rate")
}

// 更新KOL统计数据响应
struct UpdateKolStatsResp {
    1: common.BaseResp base_resp
}

// 创建/更新KOL报价Plan请求
struct SaveKolPlanReq {
    1: optional i64 id (api.body="id")  // 不传则创建，传了则更新
    2: string title (api.body="title")
    3: string description (api.body="description")
    4: double price (api.body="price")
    5: string plan_type (api.body="plan_type")  // basic, standard, premium
}

// 创建/更新KOL报价Plan响应
struct SaveKolPlanResp {
    1: common.BaseResp base_resp
    2: optional i64 plan_id
}

// 删除KOL报价Plan请求
struct DeleteKolPlanReq {
    1: i64 plan_id (api.body="plan_id")
}

// 删除KOL报价Plan响应
struct DeleteKolPlanResp {
    1: common.BaseResp base_resp
}

// 获取KOL报价Plans请求
struct GetKolPlansReq {
    1: optional i64 kol_id (api.query="kol_id")  // 不传则获取当前登录用户的Plans
}

// 获取KOL报价Plans响应
struct GetKolPlansResp {
    1: common.BaseResp base_resp
    2: list<KolPlan> plans
}

// 创建KOL视频请求
struct CreateKolVideoReq {
    1: string embed_code (api.body="embed_code")
    2: optional string cover_url (api.body="cover_url")
}

// 创建KOL视频响应
struct CreateKolVideoResp {
    1: common.BaseResp base_resp
    2: optional i64 video_id
}

// 更新KOL视频请求
struct UpdateKolVideoReq {
    1: i64 video_id (api.body="video_id")
    2: string embed_code (api.body="embed_code")
    3: optional string cover_url (api.body="cover_url")
}

// 更新KOL视频响应
struct UpdateKolVideoResp {
    1: common.BaseResp base_resp
}

// 删除KOL视频请求
struct DeleteKolVideoReq {
    1: i64 video_id (api.body="video_id")
}

// 删除KOL视频响应
struct DeleteKolVideoResp {
    1: common.BaseResp base_resp
}

// 获取KOL视频列表请求
struct GetKolVideosReq {
    1: optional i64 kol_id (api.query="kol_id")  // 不传则获取当前登录用户的视频
    2: optional i32 page (api.query="page")  // 默认1
    3: optional i32 page_size (api.query="page_size")  // 默认10
}

// 获取KOL视频列表响应
struct GetKolVideosResp {
    1: common.BaseResp base_resp
    2: list<KolVideo> videos
    3: i64 total
}

// KOL列表请求
struct GetKolListReq {
    1: optional string status (api.query="status")  // pending, approved, rejected，不传则查所有approved
    2: optional string country (api.query="country")
    3: optional string tag (api.query="tag")
    4: optional i32 page (api.query="page")  // 默认1
    5: optional i32 page_size (api.query="page_size")  // 默认10
}

// KOL列表响应
struct GetKolListResp {
    1: common.BaseResp base_resp
    2: list<KolInfo> kol_list
    3: i64 total
}

// KOL服务
service KolService {
    // KOL申请和基本信息管理
    ApplyKolResp ApplyKol(1: ApplyKolReq req) (api.post="/api/v1/kol/apply")
    GetKolInfoResp GetKolInfo(1: GetKolInfoReq req) (api.post="/api/v1/kol/info")
    UpdateKolInfoResp UpdateKolInfo(1: UpdateKolInfoReq req) (api.post="/api/v1/kol/update")
    ReviewKolResp ReviewKol(1: ReviewKolReq req) (api.post="/api/v1/kol/review")
    GetKolListResp GetKolList(1: GetKolListReq req) (api.post="/api/v1/kol/list")
    
    // KOL统计数据管理
    UpdateKolStatsResp UpdateKolStats(1: UpdateKolStatsReq req) (api.post="/api/v1/kol/stats/update")
    
    // KOL报价Plans管理
    SaveKolPlanResp SaveKolPlan(1: SaveKolPlanReq req) (api.post="/api/v1/kol/plan/save")
    DeleteKolPlanResp DeleteKolPlan(1: DeleteKolPlanReq req) (api.post="/api/v1/kol/plan/delete")
    GetKolPlansResp GetKolPlans(1: GetKolPlansReq req) (api.post="/api/v1/kol/plans")
    
    // KOL视频管理
    CreateKolVideoResp CreateKolVideo(1: CreateKolVideoReq req) (api.post="/api/v1/kol/video/create")
    UpdateKolVideoResp UpdateKolVideo(1: UpdateKolVideoReq req) (api.post="/api/v1/kol/video/update")
    DeleteKolVideoResp DeleteKolVideo(1: DeleteKolVideoReq req) (api.post="/api/v1/kol/video/delete")
    GetKolVideosResp GetKolVideos(1: GetKolVideosReq req) (api.post="/api/v1/kol/videos")
}

