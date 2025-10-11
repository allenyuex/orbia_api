namespace go demo

struct HelloReq {
    1: required string name
}

struct HelloResp {
    1: required string message
    2: required i64 timestamp
}

struct User {
    1: required i64 id
    2: required string name
    3: required string email
    4: optional string phone
}

struct CreateUserReq {
    1: required string name
    2: required string email
    3: optional string phone
}

struct CreateUserResp {
    1: required i64 user_id
    2: required string message
}

struct GetUserReq {
    1: required i64 user_id
}

struct GetUserResp {
    1: optional User user
    2: required string message
}

service DemoService {
    HelloResp Hello(1: HelloReq req)
    CreateUserResp CreateUser(1: CreateUserReq req)
    GetUserResp GetUser(1: GetUserReq req)
}

