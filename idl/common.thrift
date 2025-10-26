namespace go common

// 通用响应结构
struct BaseResp {
    1: i32 code
    2: string message
}

// JWT Token 响应
struct TokenResp {
    1: string token
    2: i64 expires_in
    3: BaseResp base_resp
}

// 分页请求
struct PageReq {
    1: optional i32 page = 1 (api.query="page")
    2: optional i32 page_size = 10 (api.query="page_size")
}

// 分页响应
struct PageResp {
    1: i32 page
    2: i32 page_size
    3: i64 total
    4: i32 total_pages
}

// 通用搜索请求
struct SearchReq {
    1: optional string keyword (api.query="keyword") // 搜索关键字
    2: optional i32 page = 1 (api.query="page")
    3: optional i32 page_size = 10 (api.query="page_size")
}