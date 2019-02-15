package common

import (
	"log"

	"github.com/aliyun/aliyun-oss-go-sdk/oss"
)

var ossClient *oss.Client

func InitAliyunOSSClient() *oss.Client {

	var (
		ossEndpoint     string
		accessKeyID     string
		accessKeySecret string
	)

	// 从云配置中获取相关配置
	for _, account := range CloudConfig.Aliyun {
		accessKeyID = account.AccessKey
		accessKeySecret = account.AccessKeySecret
		for _, endpoint := range account.EndPoint {
			ossEndpoint = endpoint.Base
		}
	}

	// FIXME: 每个用户的每个地区都需要初始化一个 client 实例，若有多地域使用 map 来映射存储
	client, err := oss.New(
		ossEndpoint,
		accessKeyID,
		accessKeySecret)

	if err != nil {
		log.Fatalln("Cannot init Aliyun OSS Client: ", err)
	}

	ossClient = client

	return ossClient
}

func GetAliyunOSSClient() *oss.Client {
	return ossClient
}
