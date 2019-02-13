package common

var CloudConfig *CloudConfiguration

const (
	DEBUG bool = true
)

const (
	DBConfigFile          = "conf/db.json"
	CloudConfigFile       = "conf/cloud.yml"
	CloudConfigDBName     = "cloud"
	AliyunOSSCallbackBody = `"bucket":${bucket},"object":${object},"etag":${etag},"size":${size},"mimeType":${mimeType},"height":${imageInfo.height},"width":${imageInfo.width},"format":${imageInfo.format}`
	ReCodeLength          = 4
	UserTokenLength       = 32
	MYSQLTIMEZONE         = "Asia%2FShanghai"
)

// 文件状态码
const (
	UPLOAD_BEGIN    = iota // 0
	CANCEL_UPLOAD          // 1
	UPLOAD_FINISHED        // 2
	FILE_ACTIVE            // 3
	FILE_DELETE            // 4
)

const (
	TOKEN_VALID_MINUTES int32 = 15 // Token 有效时长
)

// RedisDB
const (
	USER_TOKEN = iota // 0
	RECODE
	TASK_QUEUE // 1
)
