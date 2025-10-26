namespace go payment_setting

include "common.thrift"

// 收款钱包设置信息
struct PaymentSetting {
    1: i64 id
    2: string network // 区块链网络
    3: string address // 钱包地址
    4: string label // 钱包标签
    5: i32 status // 状态：1-启用，0-禁用
    6: string created_at
    7: string updated_at
}

// 获取收款钱包设置列表请求
struct GetPaymentSettingListReq {
    1: optional string network (api.body="network") // 区块链网络筛选
    2: optional i32 status (api.body="status") // 状态筛选：1-启用，0-禁用
    3: optional i32 page (api.body="page")
    4: optional i32 page_size (api.body="page_size")
}

// 获取收款钱包设置列表响应
struct GetPaymentSettingListResp {
    1: list<PaymentSetting> list
    2: i64 total
    3: i32 page
    4: i32 page_size
    5: common.BaseResp base_resp
}

// 获取收款钱包设置详情请求
struct GetPaymentSettingDetailReq {
    1: required i64 id (api.path="id")
}

// 获取收款钱包设置详情响应
struct GetPaymentSettingDetailResp {
    1: optional PaymentSetting setting
    2: common.BaseResp base_resp
}

// 创建收款钱包设置请求
struct CreatePaymentSettingReq {
    1: required string network (api.body="network") // 区块链网络（如：TRC-20 - TRON Network (TRC-20)）
    2: required string address (api.body="address") // 钱包地址
    3: required string label (api.body="label") // 钱包标签
    4: optional i32 status (api.body="status") // 状态：1-启用，0-禁用，默认1
}

// 创建收款钱包设置响应
struct CreatePaymentSettingResp {
    1: optional PaymentSetting setting
    2: common.BaseResp base_resp
}

// 更新收款钱包设置请求
struct UpdatePaymentSettingReq {
    1: required i64 id (api.body="id") // 设置ID
    2: optional string network (api.body="network") // 区块链网络
    3: optional string address (api.body="address") // 钱包地址
    4: optional string label (api.body="label") // 钱包标签
    5: optional i32 status (api.body="status") // 状态：1-启用，0-禁用
}

// 更新收款钱包设置响应
struct UpdatePaymentSettingResp {
    1: optional PaymentSetting setting
    2: common.BaseResp base_resp
}

// 删除收款钱包设置请求
struct DeletePaymentSettingReq {
    1: required i64 id (api.body="id")
}

// 删除收款钱包设置响应
struct DeletePaymentSettingResp {
    1: common.BaseResp base_resp
}

// 获取启用的收款钱包设置列表（用户和管理员都可访问）
struct GetActivePaymentSettingsReq {
    1: optional string network (api.body="network") // 区块链网络筛选
}

// 获取启用的收款钱包设置列表响应
struct GetActivePaymentSettingsResp {
    1: list<PaymentSetting> list
    2: common.BaseResp base_resp
}

// 收款钱包设置服务
service PaymentSettingService {
    // 管理员接口
    GetPaymentSettingListResp GetPaymentSettingList(1: GetPaymentSettingListReq req) (api.post="/api/v1/admin/payment-settings/list")
    GetPaymentSettingDetailResp GetPaymentSettingDetail(1: GetPaymentSettingDetailReq req) (api.post="/api/v1/admin/payment-settings/:id")
    CreatePaymentSettingResp CreatePaymentSetting(1: CreatePaymentSettingReq req) (api.post="/api/v1/admin/payment-settings/create")
    UpdatePaymentSettingResp UpdatePaymentSetting(1: UpdatePaymentSettingReq req) (api.post="/api/v1/admin/payment-settings/update")
    DeletePaymentSettingResp DeletePaymentSetting(1: DeletePaymentSettingReq req) (api.post="/api/v1/admin/payment-settings/delete")
    
    // 用户接口
    GetActivePaymentSettingsResp GetActivePaymentSettings(1: GetActivePaymentSettingsReq req) (api.post="/api/v1/payment-settings/active")
}


