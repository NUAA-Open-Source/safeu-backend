package common

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
