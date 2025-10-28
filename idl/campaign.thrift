namespace go campaign

include "common.thrift"

// Campaign附件信息
struct CampaignAttachment {
    1: i64 id
    2: string file_url
    3: string file_name
    4: string file_type
    5: i64 file_size
    6: string created_at
}

// Campaign详细信息
struct CampaignInfo {
    1: i64 id
    2: string campaign_id
    3: i64 user_id
    4: i64 team_id
    5: string campaign_name
    6: string promotion_objective  // awareness, consideration, conversion
    7: string optimization_goal
    8: optional list<i64> location  // JSON数组的数据字典ID列表
    9: optional i64 age  // 引用数据字典ID
    10: optional i64 gender  // 引用数据字典ID
    11: optional list<i64> languages  // JSON数组的数据字典ID列表
    12: optional i64 spending_power  // 引用数据字典ID
    13: optional i64 operating_system  // 引用数据字典ID
    14: optional list<i64> os_versions  // JSON数组的数据字典ID列表
    15: optional list<i64> device_models  // JSON数组的数据字典ID列表
    16: optional list<i64> connection_types  // JSON数组的数据字典ID列表
    17: i32 device_price_type  // 0-any, 1-specific range
    18: optional double device_price_min
    19: optional double device_price_max
    20: string planned_start_time
    21: string planned_end_time
    22: optional i64 time_zone  // 引用数据字典ID
    23: i32 dayparting_type  // 0-全天, 1-特定时段
    24: optional string dayparting_schedule  // JSON格式
    25: i32 frequency_cap_type  // 0-每七天不超过三次, 1-每天不超过一次, 2-自定义
    26: optional i32 frequency_cap_times
    27: optional i32 frequency_cap_days
    28: i32 budget_type  // 0-每日预算, 1-总预算
    29: double budget_amount
    30: optional string website
    31: optional string ios_download_url
    32: optional string android_download_url
    33: string status  // pending, active, paused, ended
    34: list<CampaignAttachment> attachments
    35: string created_at
    36: string updated_at
}

// 创建Campaign请求
struct CreateCampaignReq {
    1: string campaign_name (api.body="campaign_name")
    2: string promotion_objective (api.body="promotion_objective")
    3: string optimization_goal (api.body="optimization_goal")
    4: optional list<i64> location (api.body="location")
    5: optional i64 age (api.body="age")
    6: optional i64 gender (api.body="gender")
    7: optional list<i64> languages (api.body="languages")
    8: optional i64 spending_power (api.body="spending_power")
    9: optional i64 operating_system (api.body="operating_system")
    10: optional list<i64> os_versions (api.body="os_versions")
    11: optional list<i64> device_models (api.body="device_models")
    12: optional list<i64> connection_types (api.body="connection_types")
    13: i32 device_price_type (api.body="device_price_type")
    14: optional double device_price_min (api.body="device_price_min")
    15: optional double device_price_max (api.body="device_price_max")
    16: string planned_start_time (api.body="planned_start_time")
    17: string planned_end_time (api.body="planned_end_time")
    18: optional i64 time_zone (api.body="time_zone")
    19: i32 dayparting_type (api.body="dayparting_type")
    20: optional string dayparting_schedule (api.body="dayparting_schedule")
    21: i32 frequency_cap_type (api.body="frequency_cap_type")
    22: optional i32 frequency_cap_times (api.body="frequency_cap_times")
    23: optional i32 frequency_cap_days (api.body="frequency_cap_days")
    24: i32 budget_type (api.body="budget_type")
    25: double budget_amount (api.body="budget_amount")
    26: optional string website (api.body="website")
    27: optional string ios_download_url (api.body="ios_download_url")
    28: optional string android_download_url (api.body="android_download_url")
    29: optional list<string> attachment_urls (api.body="attachment_urls")  // 附件URL列表
} (api.post="/api/v1/campaign/create")

// 创建Campaign响应
struct CreateCampaignResp {
    1: CampaignInfo campaign
    2: common.BaseResp base_resp
}

// 更新Campaign请求
struct UpdateCampaignReq {
    1: string campaign_id (api.body="campaign_id")
    2: optional string campaign_name (api.body="campaign_name")
    3: optional string promotion_objective (api.body="promotion_objective")
    4: optional string optimization_goal (api.body="optimization_goal")
    5: optional list<i64> location (api.body="location")
    6: optional i64 age (api.body="age")
    7: optional i64 gender (api.body="gender")
    8: optional list<i64> languages (api.body="languages")
    9: optional i64 spending_power (api.body="spending_power")
    10: optional i64 operating_system (api.body="operating_system")
    11: optional list<i64> os_versions (api.body="os_versions")
    12: optional list<i64> device_models (api.body="device_models")
    13: optional list<i64> connection_types (api.body="connection_types")
    14: optional i32 device_price_type (api.body="device_price_type")
    15: optional double device_price_min (api.body="device_price_min")
    16: optional double device_price_max (api.body="device_price_max")
    17: optional string planned_start_time (api.body="planned_start_time")
    18: optional string planned_end_time (api.body="planned_end_time")
    19: optional i64 time_zone (api.body="time_zone")
    20: optional i32 dayparting_type (api.body="dayparting_type")
    21: optional string dayparting_schedule (api.body="dayparting_schedule")
    22: optional i32 frequency_cap_type (api.body="frequency_cap_type")
    23: optional i32 frequency_cap_times (api.body="frequency_cap_times")
    24: optional i32 frequency_cap_days (api.body="frequency_cap_days")
    25: optional i32 budget_type (api.body="budget_type")
    26: optional double budget_amount (api.body="budget_amount")
    27: optional string website (api.body="website")
    28: optional string ios_download_url (api.body="ios_download_url")
    29: optional string android_download_url (api.body="android_download_url")
    30: optional list<string> attachment_urls (api.body="attachment_urls")  // 附件URL列表
} (api.post="/api/v1/campaign/update")

// 更新Campaign响应
struct UpdateCampaignResp {
    1: CampaignInfo campaign
    2: common.BaseResp base_resp
}

// 暂停/重启Campaign请求
struct UpdateCampaignStatusReq {
    1: string campaign_id (api.body="campaign_id")
    2: string status (api.body="status")  // active, paused
} (api.post="/api/v1/campaign/status")

// 暂停/重启Campaign响应
struct UpdateCampaignStatusResp {
    1: common.BaseResp base_resp
}

// 获取Campaign列表请求
struct ListCampaignsReq {
    1: optional string keyword (api.body="keyword")  // 搜索关键字
    2: optional string status (api.body="status")  // 状态筛选
    3: optional string promotion_objective (api.body="promotion_objective")  // 推广目标筛选
    4: i32 page = 1 (api.body="page")
    5: i32 page_size = 10 (api.body="page_size")
} (api.post="/api/v1/campaign/list")

// 获取Campaign列表响应
struct ListCampaignsResp {
    1: list<CampaignInfo> campaigns
    2: common.PageResp page_info
    3: common.BaseResp base_resp
}

// 获取Campaign详情请求
struct GetCampaignReq {
    1: string campaign_id (api.body="campaign_id")
} (api.post="/api/v1/campaign/detail")

// 获取Campaign详情响应
struct GetCampaignResp {
    1: CampaignInfo campaign
    2: common.BaseResp base_resp
}

// Admin - 获取所有Campaign列表请求
struct AdminListCampaignsReq {
    1: optional string keyword (api.body="keyword")
    2: optional string status (api.body="status")
    3: optional string promotion_objective (api.body="promotion_objective")
    4: optional i64 user_id (api.body="user_id")  // 按用户筛选
    5: optional i64 team_id (api.body="team_id")  // 按团队筛选
    6: i32 page = 1 (api.body="page")
    7: i32 page_size = 10 (api.body="page_size")
} (api.post="/api/v1/admin/campaign/list")

// Admin - 获取所有Campaign列表响应
struct AdminListCampaignsResp {
    1: list<CampaignInfo> campaigns
    2: common.PageResp page_info
    3: common.BaseResp base_resp
}

// Admin - 更新Campaign状态请求
struct AdminUpdateCampaignStatusReq {
    1: string campaign_id (api.body="campaign_id")
    2: string status (api.body="status")  // active, paused, ended
} (api.post="/api/v1/admin/campaign/status")

// Admin - 更新Campaign状态响应
struct AdminUpdateCampaignStatusResp {
    1: common.BaseResp base_resp
}

// Campaign服务
service CampaignService {
    // 普通用户接口
    CreateCampaignResp CreateCampaign(1: CreateCampaignReq req) (api.post="/api/v1/campaign/create")
    UpdateCampaignResp UpdateCampaign(1: UpdateCampaignReq req) (api.post="/api/v1/campaign/update")
    UpdateCampaignStatusResp UpdateCampaignStatus(1: UpdateCampaignStatusReq req) (api.post="/api/v1/campaign/status")
    ListCampaignsResp ListCampaigns(1: ListCampaignsReq req) (api.post="/api/v1/campaign/list")
    GetCampaignResp GetCampaign(1: GetCampaignReq req) (api.post="/api/v1/campaign/detail")
    
    // 管理员接口
    AdminListCampaignsResp AdminListCampaigns(1: AdminListCampaignsReq req) (api.post="/api/v1/admin/campaign/list")
    AdminUpdateCampaignStatusResp AdminUpdateCampaignStatus(1: AdminUpdateCampaignStatusReq req) (api.post="/api/v1/admin/campaign/status")
}

