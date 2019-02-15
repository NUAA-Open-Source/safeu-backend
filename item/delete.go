package item

import (
	"log"

	"a2os/safeu-backend/common"

	"github.com/aliyun/aliyun-oss-go-sdk/oss"
)

func DeleteItem(bucketName string, objectName string) error {

	client := common.GetAliyunOSSClient()

	retryCount := common.CloudConfig.AliyunConfig.Retry
	var (
		bucket *oss.Bucket
		err    error
	)

	// 阿里云操作重试机制
	for i := retryCount; i > 0; i-- {
		// 获取存储空间
		bucket, err = client.Bucket(bucketName)
		if err != nil {
			log.Println("Aliyun OSS get bucket error: ", err, ", retrying...")
			continue
		}
		break
	}
	if err != nil {
		log.Fatalln("Aliyun OSS get bucket error: ", err, ", retries out")
		return err
	}

	for i := retryCount; i > 0; i-- {
		// 删除单个文件。
		err = bucket.DeleteObject(objectName)
		if err != nil {
			log.Println("Aliyun OSS delete item error: ", err, ", retrying...")
			continue
		}
		break
	}
	if err != nil {
		log.Println("Aliyun OSS delete item error: ", err, ", retries out")
		return err
	}

	return nil
}
