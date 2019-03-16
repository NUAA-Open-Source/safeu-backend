package common

import (
	"fmt"
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
)

// FuncHandler 统一错误处理
// i 传入error,bool,int
// judge 触发正确值 非error环境下有效
// http status,errorCode
// 触发了错误 return True
// else return false
func FuncHandler(c *gin.Context, i interface{}, judge interface{}, option ...interface{}) bool {
	generalReturn := buildErrorMeta(option)
	errType := gin.ErrorTypePrivate
	// http返回码和错误码齐全则为公开错误
	if generalReturn.HTTPStatus != 0 || generalReturn.AppErrJSON.ErrCode != 0 {
		errType = gin.ErrorTypePublic
	}
	switch i.(type) {
	case nil:
		return false
	case error:
		c.Error(i.(error)).SetMeta(generalReturn).SetType(errType)
		return true
	case bool:
		if i.(bool) == judge.(bool) {
			return false
		}
		if generalReturn.CustomMessage != "" {
			c.Error(fmt.Errorf(generalReturn.CustomMessage)).SetMeta(generalReturn).SetType(errType)
		} else if generalReturn.AppErrJSON.Message != "" {
			c.Error(fmt.Errorf(generalReturn.AppErrJSON.Message)).SetMeta(generalReturn).SetType(errType)
		} else {
			c.Error(fmt.Errorf("no err")).SetMeta(generalReturn).SetType(errType)
		}
		return true
	}
	return true
}
func buildErrorMeta(option []interface{}) GeneralReturn {
	var generalReturn GeneralReturn
	for _, v := range option {
		switch v.(type) {
		case int:
			// RFC 2616 HTTP Status Code 是3位数字代码
			if v.(int) >= 1000 {
				generalReturn.AppErrJSON.ErrCode = v.(int)
				generalReturn.AppErrJSON.Message = Errors[v.(int)]
			} else {
				generalReturn.HTTPStatus = v.(int)
			}
			break
		case string:
			generalReturn.CustomMessage = v.(string)
			break
		}
	}
	return generalReturn
}

// GeneralReturn 通用码
type GeneralReturn struct {
	CustomMessage string
	HTTPStatus    int
	AppErrJSON    appErrJSON
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
		// 转义
		var metaData GeneralReturn
		switch err.Meta.(type) {
		case GeneralReturn:
			metaData = err.Meta.(GeneralReturn)
		default:
			return
		}
		switch err.Type {
		case gin.ErrorTypePublic:
			// 公开错误 返回对应Http状态码和错误码
			// 如果有自定义消息 写入日志
			if metaData.CustomMessage != "" {
				log.Println(metaData.CustomMessage)
			}
			c.JSON(metaData.HTTPStatus, metaData.AppErrJSON)
			return
		case gin.ErrorTypePrivate:
			// 如果有自定义消息 写入日志
			if metaData.CustomMessage != "" {
				log.Println(metaData.CustomMessage)
			}
			break
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
	20306: "The retrieve code mismatch auth",

	20501: "Incorrect password",
}
