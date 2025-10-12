namespace go upload

include "common.thrift"

// 图片上传类型枚举
enum ImageType {
    AVATAR = 1,     // 头像
    TEAM_ICON = 2   // 团队图标
}

// 生成上传token请求
struct GenerateUploadTokenReq {
    1: required ImageType image_type (api.body="image_type")
    2: required string file_extension (api.body="file_extension")  // 文件扩展名，如 .jpg, .png
    3: optional i64 file_size (api.body="file_size")              // 文件大小（字节）
}

// 生成上传token响应
struct GenerateUploadTokenResp {
    1: string upload_url                    // 上传URL
    2: string access_key_id                 // 访问密钥ID
    3: string secret_access_key             // 访问密钥
    4: string session_token                 // 会话令牌
    5: string bucket                        // 存储桶名称
    6: string key                          // 对象键（文件路径）
    7: string public_url                   // 公开访问URL
    8: i64 expires_in                      // 过期时间（秒）
    9: common.BaseResp base_resp
}

// 验证图片URL请求
struct ValidateImageURLReq {
    1: required string image_url (api.body="image_url")
    2: required ImageType image_type (api.body="image_type")
}

// 验证图片URL响应
struct ValidateImageURLResp {
    1: bool is_valid                       // 是否有效
    2: optional string error_message       // 错误信息
    3: common.BaseResp base_resp
}

// 上传服务
service UploadService {
    // 生成上传token
    GenerateUploadTokenResp GenerateUploadToken(1: GenerateUploadTokenReq req) (api.post="/api/v1/upload/token")
    
    // 验证图片URL
    ValidateImageURLResp ValidateImageURL(1: ValidateImageURLReq req) (api.post="/api/v1/upload/validate")
}