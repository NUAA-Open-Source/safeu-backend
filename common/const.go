package common

var CloudConfig *CloudConfiguration

const (
	DEBUG       = true
	MAINTENANCE = false
	PORT        = "8080"
)

const (
	DBConfigFile          = "conf/db.json"
	CloudConfigFile       = "conf/cloud.yml"
	CloudConfigDBName     = "cloud"
	AliyunOSSCallbackBody = `"bucket":${bucket},"object":${object},"etag":${etag},"size":${size},"mimeType":${mimeType}`
	ReCodeLength          = 4
	UserTokenLength       = 32
	MYSQLTIMEZONE         = "Asia%2FShanghai"
	SHADOWKEYPREFIX       = "shadowKey:"
)

// 文件状态码
const (
	UPLOAD_BEGIN    = iota // 0
	CANCEL_UPLOAD          // 1
	UPLOAD_FINISHED        // 2
	FILE_ACTIVE            // 3
	FILE_DELETED           // 4
)

// 文件时长
const FILE_DEFAULT_EXIST_TIME = "8h" //文件创建默认存在时长
const FILE_MAX_EXIST_TIME = 24       //文件最大存在时长

// 文件默认可下载次数
const FILE_DEFAULT_DOWNCOUNT = 10

// 文件签名下载签名过期时长 单位:秒
const FILE_DOWNLOAD_SIGNURL_TIME = 60

// 压缩包/归档类型
const (
	ARCHIVE_NULL   = iota // 0，非压缩包
	ARCHIVE_CUSTOM        // 1，临时压缩包
	ARCHIVE_FULL          // 2，全量压缩包
)

// 初始化参数
const (
	INFINITE_DOWNLOAD = -100
	DEFAULT_PROTOCOL  = "https"
)

const (
	TOKEN_VALID_MINUTES int32 = 15 // Token 有效时长
)

// 数据库连接失败重试间隔
const DB_CONNECT_FAIL_RETRY_INTERVAL = 20

// RedisDB
const (
	USER_TOKEN = iota // 0
	RECODE
	TASK_QUEUE // 1
)

// Cross-sites resource sharing settings
var CORS_ALLOW_ORIGINS = []string{
	"https://safeu.a2os.club",
	"https://test.safeu.a2os.club",
	"http://safeu.a2os.club",
	"http://test.safeu.a2os.club",
}

var CORS_ALLOW_DEBUG_ORIGINS = []string{
	"http://*",
	"https://*",
}

var CORS_ALLOW_HEADERS = []string{
	"Origin",
	"Content-Length",
	"Content-Type",
	"Token",
	"X-CSRF-TOKEN",
	"withCredentials",
}

var CORS_ALLOW_METHODS = []string{
	"GET",
	"POST",
	"PUT",
	"PATCH",
	"DELETE",
	"HEAD",
}

var CORS_EXPOSE_HEADERS = []string{
	"X-CSRF-TOKEN",
	"Token",
}

var CSRF_COOKIE_SECRET = []byte("csrf-secret")

const (
	CSRF_SESSION_NAME string = "safeu-session"
	CSRF_SECRET       string = "safeu-secret"
)

// OSS 下载请求所包含的 Content-Type 值，用于 URL 签名
var OSS_DOWNLOAD_CONTENT_TYPE = ""
