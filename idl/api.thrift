namespace go api

// 通用响应结构
struct BaseResp {
    1: i32 code
    2: string message
}

// Hello 接口
struct HelloReq {
    1: required string name (api.query="name")
}

struct HelloResp {
    1: string message
    2: i64 timestamp
    3: BaseResp base_resp
}

// 用户相关
struct User {
    1: i64 id
    2: string name
    3: string email
    4: optional string phone
    5: string created_at
    6: string updated_at
}

struct CreateUserReq {
    1: required string name (api.body="name")
    2: required string email (api.body="email")
    3: optional string phone (api.body="phone")
}

struct CreateUserResp {
    1: i64 user_id
    2: BaseResp base_resp
}

struct GetUserReq {
    1: required i64 user_id (api.path="user_id")
}

struct GetUserResp {
    1: optional User user
    2: BaseResp base_resp
}

struct ListUsersReq {
    1: optional i32 page (api.query="page")
    2: optional i32 page_size (api.query="page_size")
}

struct ListUsersResp {
    1: list<User> users
    2: i32 total
    3: BaseResp base_resp
}

// API 服务定义
service ApiService {
    // Demo 接口
    HelloResp Hello(1: HelloReq req) (api.get="/api/v1/demo/hello")
    
    // 用户管理
    CreateUserResp CreateUser(1: CreateUserReq req) (api.post="/api/v1/users")
    GetUserResp GetUser(1: GetUserReq req) (api.get="/api/v1/users/:user_id")
    ListUsersResp ListUsers(1: ListUsersReq req) (api.get="/api/v1/users")
}

