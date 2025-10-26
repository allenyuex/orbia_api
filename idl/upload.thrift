namespace go upload

include "common.thrift"

// 生成上传token请求
struct GenerateUploadTokenReq {
    1: required string file_extension (api.body="file_extension")  // 文件扩展名，如 .jpg, .png, .pdf（会自动转为小写）
    2: optional i64 file_size (api.body="file_size")              // 文件大小（字节）
}

// 生成上传token响应
struct GenerateUploadTokenResp {
    1: string upload_url                    // 预签名上传URL（使用PUT方法）
    2: string public_url                    // 上传成功后的公开访问URL
    3: i64 expires_in                       // 过期时间（秒）
    4: map<string, string> headers          // 上传时必需的HTTP请求头
    5: common.BaseResp base_resp
}

// 验证文件URL请求
struct ValidateFileURLReq {
    1: required string file_url (api.body="file_url")
}

// 验证文件URL响应
struct ValidateFileURLResp {
    1: bool is_valid                       // 是否有效
    2: optional string error_message       // 错误信息
    3: common.BaseResp base_resp
}

// 上传服务
service UploadService {
    // 生成上传token
    GenerateUploadTokenResp GenerateUploadToken(1: GenerateUploadTokenReq req) (api.post="/api/v1/upload/token")
    
    // 验证文件URL
    ValidateFileURLResp ValidateFileURL(1: ValidateFileURLReq req) (api.post="/api/v1/upload/validate")
}