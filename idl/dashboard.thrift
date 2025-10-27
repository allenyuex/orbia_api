namespace go dashboard

include "common.thrift"

// ==================== 优秀广告案例管理 ====================

// 创建优秀案例请求
struct CreateExcellentCaseReq {
    1: string video_url (api.body="video_url") // 视频URL
    2: string cover_url (api.body="cover_url") // 封面URL
    3: string title (api.body="title") // 案例标题
    4: optional string description (api.body="description") // 案例描述
    5: optional i32 sort_order (api.body="sort_order") // 排序序号
}

// 创建优秀案例响应
struct CreateExcellentCaseResp {
    1: common.BaseResp base_resp
    2: optional i64 id // 创建的案例ID
}

// 更新优秀案例请求
struct UpdateExcellentCaseReq {
    1: i64 id (api.body="id") // 案例ID
    2: optional string video_url (api.body="video_url") // 视频URL
    3: optional string cover_url (api.body="cover_url") // 封面URL
    4: optional string title (api.body="title") // 案例标题
    5: optional string description (api.body="description") // 案例描述
    6: optional i32 sort_order (api.body="sort_order") // 排序序号
    7: optional i32 status (api.body="status") // 状态：1-启用，0-禁用
}

// 更新优秀案例响应
struct UpdateExcellentCaseResp {
    1: common.BaseResp base_resp
}

// 删除优秀案例请求
struct DeleteExcellentCaseReq {
    1: i64 id (api.body="id") // 案例ID
}

// 删除优秀案例响应
struct DeleteExcellentCaseResp {
    1: common.BaseResp base_resp
}

// 获取优秀案例列表请求
struct GetExcellentCaseListReq {
    1: optional i32 status (api.query="status") // 状态筛选：1-启用，0-禁用
    2: optional i32 page = 1 (api.query="page")
    3: optional i32 page_size = 10 (api.query="page_size")
}

// 优秀案例项
struct ExcellentCaseItem {
    1: i64 id
    2: string video_url
    3: string cover_url
    4: string title
    5: optional string description
    6: i32 sort_order
    7: i32 status
    8: string created_at
    9: string updated_at
}

// 获取优秀案例列表响应
struct GetExcellentCaseListResp {
    1: common.BaseResp base_resp
    2: list<ExcellentCaseItem> cases
    3: common.PageResp page_info
}

// 获取优秀案例详情请求
struct GetExcellentCaseDetailReq {
    1: i64 id (api.path="id")
}

// 获取优秀案例详情响应
struct GetExcellentCaseDetailResp {
    1: common.BaseResp base_resp
    2: optional ExcellentCaseItem case_detail
}

// ==================== 内容趋势管理 ====================

// 创建内容趋势请求
struct CreateContentTrendReq {
    1: i32 ranking (api.body="ranking") // 排名
    2: string hot_keyword (api.body="hot_keyword") // 热点词
    3: string value_level (api.body="value_level") // 价值等级：low, medium, high
    4: i64 heat (api.body="heat") // 热度值
    5: double growth_rate (api.body="growth_rate") // 增长比例
}

// 创建内容趋势响应
struct CreateContentTrendResp {
    1: common.BaseResp base_resp
    2: optional i64 id // 创建的趋势ID
}

// 更新内容趋势请求
struct UpdateContentTrendReq {
    1: i64 id (api.body="id") // 趋势ID
    2: optional i32 ranking (api.body="ranking") // 排名
    3: optional string hot_keyword (api.body="hot_keyword") // 热点词
    4: optional string value_level (api.body="value_level") // 价值等级：low, medium, high
    5: optional i64 heat (api.body="heat") // 热度值
    6: optional double growth_rate (api.body="growth_rate") // 增长比例
    7: optional i32 status (api.body="status") // 状态：1-启用，0-禁用
}

// 更新内容趋势响应
struct UpdateContentTrendResp {
    1: common.BaseResp base_resp
}

// 删除内容趋势请求
struct DeleteContentTrendReq {
    1: i64 id (api.body="id") // 趋势ID
}

// 删除内容趋势响应
struct DeleteContentTrendResp {
    1: common.BaseResp base_resp
}

// 获取内容趋势列表请求
struct GetContentTrendListReq {
    1: optional i32 status (api.query="status") // 状态筛选：1-启用，0-禁用
    2: optional i32 page = 1 (api.query="page")
    3: optional i32 page_size = 10 (api.query="page_size")
}

// 内容趋势项
struct ContentTrendItem {
    1: i64 id
    2: i32 ranking
    3: string hot_keyword
    4: string value_level
    5: i64 heat
    6: double growth_rate
    7: i32 status
    8: string created_at
    9: string updated_at
}

// 获取内容趋势列表响应
struct GetContentTrendListResp {
    1: common.BaseResp base_resp
    2: list<ContentTrendItem> trends
    3: common.PageResp page_info
}

// 获取内容趋势详情请求
struct GetContentTrendDetailReq {
    1: i64 id (api.path="id")
}

// 获取内容趋势详情响应
struct GetContentTrendDetailResp {
    1: common.BaseResp base_resp
    2: optional ContentTrendItem trend_detail
}

// ==================== 平台数据统计管理 ====================

// 更新平台数据请求
struct UpdatePlatformStatsReq {
    1: optional i64 active_kols (api.body="active_kols") // 活跃的KOLs数量
    2: optional i64 total_coverage (api.body="total_coverage") // 总覆盖用户数
    3: optional i64 total_ad_impressions (api.body="total_ad_impressions") // 累计广告曝光次数
    4: optional double total_transaction_amount (api.body="total_transaction_amount") // 平台总交易额
    5: optional double average_roi (api.body="average_roi") // 平均ROI
    6: optional double average_cpm (api.body="average_cpm") // 平均CPM
    7: optional i64 web3_brand_count (api.body="web3_brand_count") // 合作Web3品牌数
}

// 更新平台数据响应
struct UpdatePlatformStatsResp {
    1: common.BaseResp base_resp
}

// 获取平台数据请求
struct GetPlatformStatsReq {
}

// 平台数据统计
struct PlatformStatsData {
    1: i64 id
    2: i64 active_kols
    3: i64 total_coverage
    4: i64 total_ad_impressions
    5: double total_transaction_amount
    6: double average_roi
    7: double average_cpm
    8: i64 web3_brand_count
    9: string created_at
    10: string updated_at
}

// 获取平台数据响应
struct GetPlatformStatsResp {
    1: common.BaseResp base_resp
    2: optional PlatformStatsData stats
}

// ==================== Dashboard 数据（普通用户接口） ====================

// 获取 Dashboard 数据请求
struct GetDashboardDataReq {
}

// Dashboard 数据响应
struct GetDashboardDataResp {
    1: common.BaseResp base_resp
    2: list<ExcellentCaseItem> excellent_cases // 优秀广告案例列表（只返回启用的）
    3: list<ContentTrendItem> content_trends // 内容趋势列表（只返回启用的，按排名排序）
    4: optional PlatformStatsData platform_stats // 平台数据统计
}

// Dashboard 服务
service DashboardService {
    // 优秀广告案例管理（Admin接口）
    CreateExcellentCaseResp CreateExcellentCase(1: CreateExcellentCaseReq req) (api.post="/api/v1/admin/dashboard/excellent-case/create")
    UpdateExcellentCaseResp UpdateExcellentCase(1: UpdateExcellentCaseReq req) (api.post="/api/v1/admin/dashboard/excellent-case/update")
    DeleteExcellentCaseResp DeleteExcellentCase(1: DeleteExcellentCaseReq req) (api.post="/api/v1/admin/dashboard/excellent-case/delete")
    GetExcellentCaseListResp GetExcellentCaseList(1: GetExcellentCaseListReq req) (api.post="/api/v1/admin/dashboard/excellent-case/list")
    GetExcellentCaseDetailResp GetExcellentCaseDetail(1: GetExcellentCaseDetailReq req) (api.post="/api/v1/admin/dashboard/excellent-case/:id")
    
    // 内容趋势管理（Admin接口）
    CreateContentTrendResp CreateContentTrend(1: CreateContentTrendReq req) (api.post="/api/v1/admin/dashboard/content-trend/create")
    UpdateContentTrendResp UpdateContentTrend(1: UpdateContentTrendReq req) (api.post="/api/v1/admin/dashboard/content-trend/update")
    DeleteContentTrendResp DeleteContentTrend(1: DeleteContentTrendReq req) (api.post="/api/v1/admin/dashboard/content-trend/delete")
    GetContentTrendListResp GetContentTrendList(1: GetContentTrendListReq req) (api.post="/api/v1/admin/dashboard/content-trend/list")
    GetContentTrendDetailResp GetContentTrendDetail(1: GetContentTrendDetailReq req) (api.post="/api/v1/admin/dashboard/content-trend/:id")
    
    // 平台数据统计管理（Admin接口）
    UpdatePlatformStatsResp UpdatePlatformStats(1: UpdatePlatformStatsReq req) (api.post="/api/v1/admin/dashboard/platform-stats/update")
    GetPlatformStatsResp GetPlatformStats(1: GetPlatformStatsReq req) (api.post="/api/v1/admin/dashboard/platform-stats")
    
    // 普通用户接口
    GetDashboardDataResp GetDashboardData(1: GetDashboardDataReq req) (api.post="/api/v1/dashboard/data")
}

