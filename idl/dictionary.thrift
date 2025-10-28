namespace go dictionary

include "common.thrift"

// ==================== 字典管理 ====================

// 创建字典请求
struct CreateDictionaryReq {
    1: string code (api.body="code") // 字典编码，只能大小写字母
    2: string name (api.body="name") // 字典名称
    3: optional string description (api.body="description") // 字典描述
}

// 创建字典响应
struct CreateDictionaryResp {
    1: common.BaseResp base_resp
    2: optional DictionaryInfo dictionary
}

// 更新字典请求（不能修改code）
struct UpdateDictionaryReq {
    1: i64 id (api.body="id") // 字典ID
    2: string name (api.body="name") // 字典名称
    3: optional string description (api.body="description") // 字典描述
    4: optional i32 status (api.body="status") // 状态：1-启用，0-禁用
}

// 更新字典响应
struct UpdateDictionaryResp {
    1: common.BaseResp base_resp
    2: optional DictionaryInfo dictionary
}

// 删除字典请求（软删除）
struct DeleteDictionaryReq {
    1: i64 id (api.body="id") // 字典ID
}

// 删除字典响应
struct DeleteDictionaryResp {
    1: common.BaseResp base_resp
}

// 获取字典列表请求
struct GetDictionariesReq {
    1: optional string keyword (api.query="keyword") // 搜索关键字（字典编码、名称）
    2: optional i32 status (api.query="status") // 状态筛选：1-启用，0-禁用
    3: optional i32 page = 1 (api.query="page")
    4: optional i32 page_size = 10 (api.query="page_size")
}

// 字典信息
struct DictionaryInfo {
    1: i64 id
    2: string code
    3: string name
    4: optional string description
    5: i32 status // 1-启用，0-禁用
    6: string created_at
    7: string updated_at
}

// 获取字典列表响应
struct GetDictionariesResp {
    1: common.BaseResp base_resp
    2: list<DictionaryInfo> dictionaries
    3: common.PageResp page_info
}

// 获取字典详情请求
struct GetDictionaryDetailReq {
    1: i64 id (api.path="id") // 字典ID
}

// 获取字典详情响应
struct GetDictionaryDetailResp {
    1: common.BaseResp base_resp
    2: optional DictionaryInfo dictionary
}

// 根据编码获取字典请求
struct GetDictionaryByCodeReq {
    1: string code (api.query="code") // 字典编码
}

// 根据编码获取字典响应
struct GetDictionaryByCodeResp {
    1: common.BaseResp base_resp
    2: optional DictionaryInfo dictionary
}

// ==================== 字典项管理 ====================

// 创建字典项请求
struct CreateDictionaryItemReq {
    1: i64 dictionary_id (api.body="dictionary_id") // 字典ID
    2: i64 parent_id (api.body="parent_id") // 父级ID（0表示根节点）
    3: string code (api.body="code") // 字典项编码
    4: string name (api.body="name") // 字典项名称
    5: optional string description (api.body="description") // 字典项描述
    6: optional string icon_url (api.body="icon_url") // 图标URL
    7: optional i32 sort_order (api.body="sort_order") // 排序序号
}

// 创建字典项响应
struct CreateDictionaryItemResp {
    1: common.BaseResp base_resp
    2: optional DictionaryItemInfo item
}

// 更新字典项请求
struct UpdateDictionaryItemReq {
    1: i64 id (api.body="id") // 字典项ID
    2: optional string name (api.body="name") // 字典项名称
    3: optional string description (api.body="description") // 字典项描述
    4: optional string icon_url (api.body="icon_url") // 图标URL
    5: optional i32 sort_order (api.body="sort_order") // 排序序号
    6: optional i32 status (api.body="status") // 状态：1-启用，0-禁用
}

// 更新字典项响应
struct UpdateDictionaryItemResp {
    1: common.BaseResp base_resp
    2: optional DictionaryItemInfo item
}

// 删除字典项请求（软删除）
struct DeleteDictionaryItemReq {
    1: i64 id (api.body="id") // 字典项ID
}

// 删除字典项响应
struct DeleteDictionaryItemResp {
    1: common.BaseResp base_resp
}

// 获取字典项列表请求
struct GetDictionaryItemsReq {
    1: i64 dictionary_id (api.query="dictionary_id") // 字典ID
    2: optional i64 parent_id (api.query="parent_id") // 父级ID（不传则获取所有）
    3: optional i32 status (api.query="status") // 状态筛选：1-启用，0-禁用
    4: optional i32 page = 1 (api.query="page")
    5: optional i32 page_size = 100 (api.query="page_size")
}

// 字典项信息（平铺结构）
struct DictionaryItemInfo {
    1: i64 id
    2: i64 dictionary_id
    3: i64 parent_id
    4: string code
    5: string name
    6: optional string description
    7: optional string icon_url
    8: i32 sort_order
    9: i32 level
    10: string path
    11: i32 status // 1-启用，0-禁用
    12: string created_at
    13: string updated_at
}

// 获取字典项列表响应
struct GetDictionaryItemsResp {
    1: common.BaseResp base_resp
    2: list<DictionaryItemInfo> items
    3: common.PageResp page_info
}

// 字典项树形节点
struct DictionaryItemTreeNode {
    1: i64 id
    2: string code
    3: string name
    4: optional string description
    5: optional string icon_url
    6: i32 sort_order
    7: i32 level
    8: i32 status
    9: list<DictionaryItemTreeNode> children // 子节点
}

// 字典及其字典项树形结构
struct DictionaryWithTree {
    1: DictionaryInfo dictionary // 字典基本信息
    2: list<DictionaryItemTreeNode> tree // 字典项树形结构
}

// 批量获取字典和字典项请求（用于前端冷启动）
struct GetAllDictionariesWithItemsReq {
    1: optional i32 page = 1 (api.body="page") // 页码，默认1
    2: optional i32 page_size = 20 (api.body="page_size") // 每页数量，默认20，最大20
}

// 批量获取字典和字典项响应
struct GetAllDictionariesWithItemsResp {
    1: common.BaseResp base_resp
    2: list<DictionaryWithTree> dictionaries // 字典列表（包含树形字典项）
    3: common.PageResp page_info // 分页信息
}

// 字典服务（仅管理员可用）
service DictionaryService {
    // 字典管理
    CreateDictionaryResp CreateDictionary(1: CreateDictionaryReq req) (api.post="/api/v1/admin/dictionary/create")
    UpdateDictionaryResp UpdateDictionary(1: UpdateDictionaryReq req) (api.post="/api/v1/admin/dictionary/update")
    DeleteDictionaryResp DeleteDictionary(1: DeleteDictionaryReq req) (api.post="/api/v1/admin/dictionary/delete")
    GetDictionariesResp GetDictionaries(1: GetDictionariesReq req) (api.post="/api/v1/admin/dictionary/list")
    GetDictionaryDetailResp GetDictionaryDetail(1: GetDictionaryDetailReq req) (api.post="/api/v1/admin/dictionary/:id")
    GetDictionaryByCodeResp GetDictionaryByCode(1: GetDictionaryByCodeReq req) (api.post="/api/v1/admin/dictionary/by-code")
    
    // 字典项管理
    CreateDictionaryItemResp CreateDictionaryItem(1: CreateDictionaryItemReq req) (api.post="/api/v1/admin/dictionary/item/create")
    UpdateDictionaryItemResp UpdateDictionaryItem(1: UpdateDictionaryItemReq req) (api.post="/api/v1/admin/dictionary/item/update")
    DeleteDictionaryItemResp DeleteDictionaryItem(1: DeleteDictionaryItemReq req) (api.post="/api/v1/admin/dictionary/item/delete")
    GetDictionaryItemsResp GetDictionaryItems(1: GetDictionaryItemsReq req) (api.post="/api/v1/admin/dictionary/item/list")
    
    // 批量获取所有字典和字典项（公开接口，用于前端冷启动）
    GetAllDictionariesWithItemsResp GetAllDictionariesWithItems(1: GetAllDictionariesWithItemsReq req) (api.post="/api/v1/dictionary/all")
}

