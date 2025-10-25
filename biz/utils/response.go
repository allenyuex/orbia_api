package utils

import (
	"time"

	"github.com/cloudwego/hertz/pkg/app"
	"github.com/cloudwego/hertz/pkg/protocol/consts"

	apiconsts "orbia_api/biz/consts"
)

// Response 统一响应结构
type Response struct {
	Code    int         `json:"code"`
	Message string      `json:"message"`
	Data    interface{} `json:"data,omitempty"`
}

// Success 成功响应
func Success(c *app.RequestContext, data interface{}) {
	c.JSON(consts.StatusOK, Response{
		Code:    apiconsts.SuccessCode,
		Message: apiconsts.SuccessMsg,
		Data:    data,
	})
}

// SuccessResponse 成功响应（另一种写法）
func SuccessResponse(c *app.RequestContext, data interface{}) {
	Success(c, data)
}

// Error 错误响应
func Error(c *app.RequestContext, code int, message string) {
	c.JSON(consts.StatusOK, Response{
		Code:    code,
		Message: message,
	})
}

// ErrorResponse 错误响应（另一种写法）
func ErrorResponse(c *app.RequestContext, code int, message string) {
	Error(c, code, message)
}

// ParamError 参数错误响应
func ParamError(c *app.RequestContext, message string) {
	if message == "" {
		message = apiconsts.ParamInvalidMsg
	}
	Error(c, apiconsts.ParamInvalidCode, message)
}

// SystemError 系统错误响应
func SystemError(c *app.RequestContext) {
	Error(c, apiconsts.SystemErrorCode, apiconsts.SystemErrorMsg)
}

// FormatTime 格式化时间
func FormatTime(t *time.Time) string {
	if t == nil {
		return ""
	}
	return t.Format(apiconsts.DateTimeFormat)
}

// BuildErrorResp 构建错误响应
func BuildErrorResp(code int, message string) Response {
	return Response{
		Code:    code,
		Message: message,
	}
}
