package common

import (
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
			c.Error(nil).SetType(errorType).SetMeta(buildErrorMeta(errorCode))
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
	//0:  "General error.", // 20000
	//1:  "Parameter error.", // 10003
	//2:  "The length of the time is not within the right range.", // 10004
	//3:  "Can't Find User Token.", // 20201
	//4:  "Repeat retrieve code.", //
	//5:  "Auth field required.", // 10005
	//6:  "The retrieve code mismatch auth.",
	//7:  "Cannot find resource via this retrieve code.", // 10006
	//8:  "There is a problem for this resource, please contact the maintainer.", // 10001
	//9:  "Cannot get the password.", // 10005
	//10: "The password is not correct.", // 20501
	//11: "Cannot get the token.", // 20301
	//12: "Token invalid.", // 203[02:04]
	//13: "Over the expired time.", // 10006
	//14: "Out of downloadable count.", // 10006
	//15: "Cannot get the items field.", //10005
	//16: "Cannot get the download link.", //20305
	//17: "Service Unavailable, please contact the maintainer.", // 10002

	0: "OK",

	// 系统级错误
	10001: "System error",
	10002: "Service unavailable",
	10003: "Parameter error",
	10004: "Parameter value invalid",
	10005: "Missing required parameter",
	10006: "Resource unavailable",

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
