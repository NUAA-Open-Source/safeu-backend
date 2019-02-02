package common

var CloudConfig *CloudConfiguration

const (
	DBConfigFile          = "conf/db.json"
	CloudConfigFile       = "conf/cloud.yml"
	CloudConfigDBName     = "cloud"
	AliyunOSSCallbackBody = `"bucket":${bucket},"object":${object},"etag":${etag},"size":${size},"mimeType":${mimeType},"height":${imageInfo.height},"width":${imageInfo.width},"format":${imageInfo.format}`
)
