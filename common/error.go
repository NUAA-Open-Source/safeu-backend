package common

import (
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
)

// FuncHandler 统一错误处理
// i 传入error,bool,int
// judege 触发正确值 非error环境下有效
// 当 errorType 为 ErrorTypePublic 时可用 errorCode[0]是自定义错误码 errorCode[1]是Http状态码
// 触发了错误 return True
// else return false
func FuncHandler(c *gin.Context, i interface{}, judge interface{}, errorType gin.ErrorType, errorCode ...int) bool {
	switch i.(type) {
	case nil:
		return false
	case error:
		if len(errorCode) == 2 && errorType == gin.ErrorTypePublic {
			c.Error(i.(error)).SetType(errorType).SetMeta(buildErrorMeta(errorCode))
			return true
		}
		c.Error(i.(error)).SetType(errorType)
		return true
	case bool:
		if i.(bool) == judge.(bool) {
			return false
		}
		if len(errorCode) == 2 && errorType == gin.ErrorTypePublic {
			c.Error(fmt.Errorf("no err")).SetType(errorType).SetMeta(buildErrorMeta(errorCode))
			return true
		}
		return true
	}
	return true
}

func buildErrorMeta(errorCode []int) (generalErr generalErr) {
	generalErr.AppErrJSON.ErrCode = errorCode[0]
	generalErr.AppErrJSON.Message = Errors[errorCode[0]]
	generalErr.HTTPStatus = errorCode[1]
	return generalErr
}

type generalErr struct {
	HTTPStatus int
	AppErrJSON appErrJSON
}
type appErrJSON struct {
	ErrCode int    `json:"err_code"`
	Message string `json:"message"`
}

// ErrorHandling 错误处理中间件
func ErrorHandling() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Next()
		err := c.Errors.Last()
		if err == nil {
			return
		}
		var metaData generalErr
		switch err.Meta.(type) {
		case generalErr:
			metaData = err.Meta.(generalErr)
		default:
			return
		}
		switch err.Type {
		case gin.ErrorTypePublic:
			// 公开错误 返回对应Http状态码和错误码
			c.JSON(metaData.HTTPStatus, metaData.AppErrJSON)
			return
		default:
			// 非公开错误 统一返回 500
			// 日志处理
			c.JSON(http.StatusInternalServerError, gin.H{"message": "Something went wrong."})
			return
		}

	}
}

// Errors 错误码
var Errors = map[int]string{

	0: "OK",

	// 系统级错误
	10001: "System error",
	10002: "Service unavailable",
	10003: "Parameter error",
	10004: "Parameter value invalid",
	10005: "Missing required parameter",
	10006: "Resource unavailable",
	10007: "CSRF token mismatch",

	// 应用级错误
	20000: "Application error",

	20201: "Can't find user token",

	20301: "Missing token in header",
	20302: "Token used",
	20303: "Token expired",
	20304: "Token revoked",
	20305: "Can't get the download link",

	20501: "Incorrect password",
}
