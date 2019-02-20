package item

import (
	"fmt"
	"log"
	"net/http"
	"sync"

	"a2os/safeu-backend/common"
	"github.com/gin-gonic/gin"

	"github.com/aliyun/aliyun-oss-go-sdk/oss"
)

type DelateItemBody struct {
	UserToken string `json:"user_token"`
}

func DeleteManual(c *gin.Context) {
	retrieveCode := c.Param("retrieveCode")
	var deleteItemBody DelateItemBody
	err := c.BindJSON(&deleteItemBody)
	if err != nil {
		log.Println(err)
		c.JSON(http.StatusInternalServerError, gin.H{
			"message": err,
		})
		return
	}
	// 判断提供的提取码和提供的Auth是否相符
	reCodeRedisClient := common.GetReCodeRedisClient()
	if deleteItemBody.UserToken != reCodeRedisClient.Get(retrieveCode).Val() {
		c.JSON(http.StatusUnauthorized, gin.H{
			"message": "this reCode mismatch auth",
		})
		return
	}
	db := common.GetDB()
	var itemList []Item
	db.Where("re_code = ? AND (status = ? OR status = ?)", retrieveCode, common.UPLOAD_FINISHED, common.FILE_ACTIVE).Find(&itemList)
	for _, item := range itemList {
		db.Delete(item)
	}
	DeleteItems(itemList)
	c.JSON(http.StatusOK, gin.H{
		"message": "ok",
	})
	return
}

// 多线程删除文件
func DeleteItems(items []Item) {
	threads := len(items)
	errChan := make(chan error, threads)
	var waitGroup sync.WaitGroup
	for i := 0; i < threads; i++ {
		waitGroup.Add(1)
		go DeleteItemWaitGroup(&waitGroup, items[i], errChan)
	}
	waitGroup.Wait()
	var errs []error
	for i := 0; i < threads; i++ {
		errs = append(errs, <-errChan)
	}
	for i := 0; i < len(errs); i++ {
		if errs[i] != nil {
			log.Println(errs[i])
		}
	}

}

func DeleteItemWaitGroup(group *sync.WaitGroup, item Item, errChan chan error) {
	defer func() {
		group.Done()
		if err := recover(); err != nil {
			fmt.Println("work thread error:", err)
		}
	}()
	err := DeleteItem(item.Bucket, item.Path)
	errChan <- err
}

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
			log.Println(bucketName, "Aliyun OSS get bucket error: ", err, ", retrying...")
			continue
		}
		break
	}
	if err != nil {
		log.Println(bucketName, "Aliyun OSS get bucket error: ", err, ", retries out")
		return err
	}

	for i := retryCount; i > 0; i-- {
		// 删除单个文件。
		err = bucket.DeleteObject(objectName)
		if err != nil {
			log.Println("bucket: ", bucketName, "object:", objectName, "Aliyun OSS delete item error: ", err, ", retrying...")
			continue
		}
		break
	}
	if err != nil {
		log.Println("bucket: ", bucketName, "object:", objectName, "Aliyun OSS delete item error: ", err, ", retries out")
		return err
	}

	return nil
}
