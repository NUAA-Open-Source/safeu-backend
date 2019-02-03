package common

var CloudConfig *CloudConfiguration

type StatusCode int32

const (
	DBConfigFile          = "conf/db.json"
	CloudConfigFile       = "conf/cloud.yml"
	CloudConfigDBName     = "cloud"
	AliyunOSSCallbackBody = `"bucket":${bucket},"object":${object},"etag":${etag},"size":${size},"mimeType":${mimeType},"height":${imageInfo.height},"width":${imageInfo.width},"format":${imageInfo.format}`
	ReCodeLength          = 4
)

// 文件状态码
const (
	UPLOAD_BEGIN    StatusCode = 0
	CANCEL_UPLOAD   StatusCode = 1
	UPLOAD_FINISHED StatusCode = 2
	FILE_ACTIVE     StatusCode = 3
	FILE_DELETE     StatusCode = 4
)
