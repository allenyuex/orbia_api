namespace go common

// 通用响应结构
struct BaseResp {
    1: i32 code
    2: string message
    3: optional string data
}

// JWT Token 响应
struct TokenResp {
    1: string token
    2: i64 expires_in
    3: BaseResp base_resp
}

// 分页请求
struct PageReq {
    1: optional i32 page = 1
    2: optional i32 page_size = 10
}

// 分页响应
struct PageResp {
    1: i32 page
    2: i32 page_size
    3: i64 total
    4: i32 total_pages
}