package consts

// 错误码定义
const (
	SuccessCode = 0
	ErrorCode   = 10000

	// 参数错误 10001-10099
	ParamInvalidCode = 10001
	ParamMissingCode = 10002

	// 用户相关错误 10100-10199
	UserNotFoundCode     = 10100
	UserAlreadyExistCode = 10101
	UserCreateFailedCode = 10102
	UserUpdateFailedCode = 10103
	UserDeleteFailedCode = 10104

	// 数据库错误 10200-10299
	DBErrorCode       = 10200
	DBConnectFailCode = 10201
	DBQueryFailCode   = 10202
	DBInsertFailCode  = 10203
	DBUpdateFailCode  = 10204
	DBDeleteFailCode  = 10205

	// 业务错误 10300-10399
	BusinessErrorCode = 10300

	// 系统错误 10400-10499
	SystemErrorCode = 10400

	// 认证相关错误 10500-10599
	UnauthorizedCode = 10500
	ForbiddenCode    = 10501
)

// 错误消息
const (
	SuccessMsg = "success"
	ErrorMsg   = "internal error"

	ParamInvalidMsg = "参数无效"
	ParamMissingMsg = "参数缺失"

	UserNotFoundMsg     = "用户不存在"
	UserAlreadyExistMsg = "用户已存在"
	UserCreateFailedMsg = "创建用户失败"
	UserUpdateFailedMsg = "更新用户失败"
	UserDeleteFailedMsg = "删除用户失败"

	DBErrorMsg       = "数据库错误"
	DBConnectFailMsg = "数据库连接失败"
	DBQueryFailMsg   = "数据库查询失败"
	DBInsertFailMsg  = "数据库插入失败"
	DBUpdateFailMsg  = "数据库更新失败"
	DBDeleteFailMsg  = "数据库删除失败"

	BusinessErrorMsg = "业务错误"
	SystemErrorMsg   = "系统错误"

	UnauthorizedMsg = "未授权"
	ForbiddenMsg    = "无权访问"
)
