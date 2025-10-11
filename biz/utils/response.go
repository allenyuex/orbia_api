package utils

import (
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

// Error 错误响应
func Error(c *app.RequestContext, code int, message string) {
	c.JSON(consts.StatusOK, Response{
		Code:    code,
		Message: message,
	})
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
