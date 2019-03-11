package common

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
