package item

import (
	"log"

	"a2os/safeu-backend/common"
)

func DeleteItem(bucketName string, objectName string) (error) {

	client := common.GetAliyunOSSClient()

	// 获取存储空间。
	bucket, err := client.Bucket(bucketName)
	if err != nil {
		log.Fatalln("Aliyun OSS get bucket error: ", err)
		return err
	}

	// 删除单个文件。
	err = bucket.DeleteObject(objectName)
	if err != nil {
		log.Fatalln("Aliyun OSS delete item error: ", err)
		return err
	}

	return nil
}
