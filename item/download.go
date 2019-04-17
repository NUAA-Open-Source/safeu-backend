package item

import (
	"bytes"
	"fmt"
	"log"
	"net/http"
	"time"

	"a2os/safeu-backend/common"

	"github.com/aliyun/aliyun-oss-go-sdk/oss"
	"github.com/gin-gonic/gin"
	"github.com/gin-gonic/gin/json"
	uuid "github.com/satori/go.uuid"
)

type PackRequest struct {
	ZipItems []ZipItem `json:"items"`
	Full     bool      `json:"full"`
}

type ZipItem struct {
	Protocol     string `json:"protocol"`
	Bucket       string `json:"bucket"`
	Endpoint     string `json:"endpoint"`
	Path         string `json:"path"`
	OriginalName string `json:"original_name"`
	//AccessKey       string
	//AccessKeySecret string
}

func DownloadItems(c *gin.Context) {
	retrieveCode := c.Param("retrieveCode")
	log.Println(c.ClientIP(), " Received download request for \"", retrieveCode, "\" resources")

	db := common.GetDB()

	// 检查并更新 Token 有效

	clientToken := c.Request.Header.Get("Token")
	if len(clientToken) == 0 { // if not get the token
		c.JSON(http.StatusUnauthorized, gin.H{
			"err_code": 20301,
			"message":  common.Errors[20301],
		})
		log.Println(c.ClientIP(), " Cannot get the client token from header")
		return
	}

	log.Println(c.ClientIP(), " Get client token ", clientToken)

	var tokenRecord Token
	if db.Where("token = ?", clientToken).First(&tokenRecord).RecordNotFound() {
		// 无法找到该 token
		c.JSON(http.StatusNotFound, gin.H{
			"err_code": 20304,
			"message":  common.Errors[20304],
		})
		log.Println(c.ClientIP(), " Invalid token ", clientToken)
		return
	}

	// 检查 token 是否失效
	if !tokenRecord.Valid {
		c.JSON(http.StatusUnauthorized, gin.H{
			"err_code": 20302,
			"message":  common.Errors[20302],
		})
		log.Println(c.ClientIP(), " Expired token ", clientToken)
		return
	}

	// 检查 token 是否过期
	tokenExpiredAt := tokenRecord.ExpiredAt
	if tokenExpiredAt.Before(time.Now()) {
		// Token 已过期，更新数据库并拒绝请求
		db.Model(&tokenRecord).Update("valid", false)
		c.JSON(http.StatusUnauthorized, gin.H{
			"err_code": 20303,
			"message":  common.Errors[20303],
		})
		log.Println(c.ClientIP(), " Expired token ", clientToken)
		return
	}

	// Token 没过期，核对提取码
	tokenRetrieveCode := tokenRecord.RetrieveCode
	if tokenRetrieveCode != retrieveCode {
		// 提取码不正确
		c.JSON(http.StatusUnauthorized, gin.H{
			"err_code": 20304,
			"message":  common.Errors[20304],
		})
		log.Println(c.ClientIP(), " Invalid token ", clientToken, " for resource ", retrieveCode)
		return
	}

	// ---------- 结束 Token 验证
	log.Println(c.ClientIP(), " The client token ", clientToken, " is valid")

	var itemList []Item
	db.Where("re_code = ? AND (status = ? OR status = ?)", retrieveCode, common.UPLOAD_FINISHED, common.FILE_ACTIVE).Find(&itemList)

	if len(itemList) == 0 {
		c.JSON(http.StatusNotFound, gin.H{
			"err_code": 10006,
			"message":  common.Errors[10006],
		})
		log.Println(c.ClientIP(), " resource ", retrieveCode, " not found")
		return
	}

	// 单文件下载
	if len(itemList) == 1 {
		var singleItem Item = itemList[0]

		// 检查文件有效时间
		if singleItem.ExpiredAt.Before(time.Now()) {
			// 文件已过期
			// 清除文件
			err := DeleteItem(singleItem.Bucket, singleItem.Path)
			if err != nil {
				log.Println("Cannot delete item in bucket ", singleItem.Bucket, ", path ", singleItem.Path)
				// TODO: 返回 500
			}

			// 删除数据库记录
			db.Delete(&singleItem)
			common.DeleteRedisRecodeFromRecode(singleItem.ReCode)
			// 返回 410 Gone
			c.JSON(http.StatusGone, gin.H{
				"err_code": 10006,
				"message":  common.Errors[10006],
			})
			log.Println(c.ClientIP(), " The retrieve code \"", retrieveCode, "\" resouce cannot be download due to the file duaration expired")
			return
		}

		// 检查剩余下载次数
		// 若为 0 次则返回 410 Gone 并删除文件
		if singleItem.DownCount <= 0 && singleItem.DownCount != common.INFINITE_DOWNLOAD {
			// 删除文件
			err := DeleteItem(singleItem.Bucket, singleItem.Path)
			if err != nil {
				log.Println("Cannot delete item in bucket ", singleItem.Bucket, ", path ", singleItem.Path)
				// TODO: 返回 500
			}

			// 删除数据库记录
			db.Delete(&singleItem)
			common.DeleteRedisRecodeFromRecode(singleItem.ReCode)
			c.JSON(http.StatusGone, gin.H{
				"err_code": 10006,
				"message":  common.Errors[10006],
			})
			log.Println(c.ClientIP(), " The retrieve code \"", retrieveCode, "\" resouce cannot be download due to downloadable counter = 0")
			return
		}

		// 剩余时间与下载次数合法，获取文件
		// 获取临时下载链接
		url, err := GetSignURL(singleItem.Bucket, singleItem.Path, common.GetAliyunOSSClient())
		if err != nil {
			log.Println("Cannot get the signed downloadable link for item \"", singleItem.Bucket, singleItem.Path, "\"")
			// TODO: 返回 500
		}
		log.Println(c.ClientIP(), " Get the zip file signed url: ", url)
		c.JSON(http.StatusOK, gin.H{
			"url": url,
		})
		return
	}

	// ----- 多文件打包下载

	var packRequest PackRequest
	if err := c.ShouldBindJSON(&packRequest); err != nil {
		// 缺少 ItemGroup
		log.Println(c.ClientIP(), " Cannot get the ItemGroup")
		c.JSON(http.StatusBadRequest, gin.H{
			"err_code": 10005,
			"message":  common.Errors[10005],
		})
		return
	}

	// 若为文件组中的单文件下载请求，直接返回签名链接
	if len(packRequest.ZipItems) == 1 {
		singleItem := packRequest.ZipItems[0]
		url, err := GetSignURL(singleItem.Bucket, singleItem.Path, common.GetAliyunOSSClient())
		if err != nil {
			log.Println("Cannot get the signed downloadable link for item \"", singleItem.Bucket, singleItem.Path, "\"")
			c.JSON(http.StatusInternalServerError, gin.H{
				"err_code": 20305,
				"message":  common.Errors[20305],
			})
			return
		}
		log.Println(c.ClientIP(), " Get the single file signed url: ", url)
		c.JSON(http.StatusOK, gin.H{
			"url": url,
		})
		return
	}

	zipEndpoint, err := GetZipEndpoint()
	// 若配置文件有误返回 503 Service Unavailable
	if err != nil {
		log.Println(err)
		c.JSON(http.StatusServiceUnavailable, gin.H{
			"err_code": 10002,
			"message":  common.Errors[10002],
		})
		log.Println("[ERROR] Cannot get the proper FaaS zip config from cloud config")
		return
	}

	// 全量打包下载
	var zipPack Item
	if packRequest.Full {
		// 全量打包下载
		if db.Where("re_code = ? AND (status = ? OR status = ?) AND archive_type = ?", retrieveCode, common.UPLOAD_FINISHED, common.FILE_ACTIVE, common.ARCHIVE_FULL).First(&zipPack).RecordNotFound() {

			// 没有全量打包，进行全量打包并将记录存储到数据库中
			resJson := ZipItemsFaaS(packRequest.ZipItems, retrieveCode, true, zipEndpoint)

			u := uuid.Must(uuid.NewV4())
			zipPack.Name = u.String()
			zipPack.Status = common.UPLOAD_FINISHED
			zipPack.ReCode = retrieveCode
			zipPack.Password = itemList[0].Password
			zipPack.IsPublic = itemList[0].IsPublic
			zipPack.Host = resJson["host"]
			zipPack.OriginalName = resJson["original_name"]
			zipPack.Protocol = resJson["protocol"]
			zipPack.Bucket = resJson["bucket"]
			zipPack.Endpoint = resJson["endpoint"]
			zipPack.Path = resJson["path"]
			zipPack.Type = resJson["type"]
			zipPack.ArchiveType = common.ARCHIVE_FULL
			zipPack.DownCount = common.INFINITE_DOWNLOAD
			zipPack.ExpiredAt = itemList[0].ExpiredAt // 过期时间跟随生成压缩包的文件有效时间

			db.Create(&zipPack)
			log.Println("Generated the full files zip package for retrieve code \"", retrieveCode, "\"")

			// 返回压缩包下载链接
			// 对压缩包签名
			url, err := GetSignURL(zipPack.Bucket, zipPack.Path, common.GetAliyunOSSClient())
			if err != nil {
				log.Println("Cannot get the signed downloadable link for item \"", zipPack.Bucket, zipPack.Path, "\"")
				// TODO: 返回 500
			}
			log.Println(c.ClientIP(), " Get the zip file signed url: ", url)
			c.JSON(http.StatusOK, gin.H{
				"url": url,
			})
			return

		}

		// 有全量打包，则直接发送打包文件
		// 对压缩包签名
		url, err := GetSignURL(zipPack.Bucket, zipPack.Path, common.GetAliyunOSSClient())
		if err != nil {
			log.Println("Cannot get the signed downloadable link for item \"", zipPack.Bucket, zipPack.Path, "\"")
			// TODO: 返回 500
		}
		c.JSON(http.StatusOK, gin.H{
			"url": url,
		})
		log.Println(c.ClientIP(), " Full zip pack has generated before, get the zip file signed url: ", url)
		return
	}

	// 自定义多文件打包下载
	resJson := ZipItemsFaaS(packRequest.ZipItems, retrieveCode, false, zipEndpoint)
	log.Println(c.ClientIP(), " Generated the custom zip file for retrieve code \"", retrieveCode, "\"")

	// 将自定义压缩包存入数据库记录
	zipPack = Item{
		Status:       common.UPLOAD_FINISHED,
		Name:         uuid.Must(uuid.NewV4()).String(),
		OriginalName: resJson["original_name"],
		Host:         resJson["host"],
		ReCode:       retrieveCode,
		Password:     itemList[0].Password,
		DownCount:    common.INFINITE_DOWNLOAD,
		Type:         resJson["type"],
		IsPublic:     itemList[0].IsPublic,
		ArchiveType:  common.ARCHIVE_CUSTOM,
		Protocol:     resJson["protocol"],
		Bucket:       resJson["bucket"],
		Endpoint:     resJson["endpoint"],
		Path:         resJson["path"],
		ExpiredAt:    itemList[0].ExpiredAt, // 过期时间跟随生成压缩包的文件有效时间
	}
	// 先清除数据库之前同提取码的自定义压缩包记录
	// [4.5.2019] 不需要，不同压缩包可以共存
	//var deleteZipPacks []Item
	//db.Where("re_code = ? AND (status = ? OR status = ?) AND archive_type = ?", retrieveCode, common.UPLOAD_FINISHED, common.FILE_ACTIVE, common.ARCHIVE_CUSTOM).Find(&deleteZipPacks)
	//for _, deleteZipPack := range deleteZipPacks {
	//	db.Delete(&deleteZipPack)
	//}

	db.Create(&zipPack)
	log.Println("Generated the custom files zip package for retrieve code \"", retrieveCode, "\"")

	// 返回压缩包路径
	// 对自定义压缩包签名
	url, err := GetSignURL(zipPack.Bucket, zipPack.Path, common.GetAliyunOSSClient())
	if err != nil {
		log.Println("Cannot get the signed downloadable link for item \"", zipPack.Bucket, zipPack.Path, "\"")
		// TODO: 返回 500
	}
	log.Println(c.ClientIP(), " Get the zip file signed url: ", url)
	c.JSON(http.StatusOK, gin.H{
		"url": url,
	})

	return
}

func ZipItemsFaaS(zipItems []ZipItem, retrieveCode string, isFull bool, endpoint string) map[string]string {
	reqJson := map[string]interface{}{
		"re_code": retrieveCode,
		"uuid":    uuid.Must(uuid.NewV4()).String(),
		"items":   zipItems,
		"full":    isFull,
	}

	bytesRepresentation, err := json.Marshal(reqJson)
	if err != nil {
		log.Println(err)
	}

	// TODO: 加入阿里云认证
	// 请求函数计算
	res, err := http.Post(endpoint, "application/json", bytes.NewBuffer(bytesRepresentation))
	if err != nil {
		// TODO: 加入重试机制
		log.Println(err)
	}

	var resJson map[string]string
	json.NewDecoder(res.Body).Decode(&resJson)

	return resJson
}

func GetZipEndpoint() (string, error) {
	var zipEndpoint string
	for _, faasConfig := range common.CloudConfig.FaaS {
		if faasConfig.Name == "zip" {
			zipEndpoint = faasConfig.Endpoint
		}
	}

	if len(zipEndpoint) == 0 {
		return zipEndpoint, fmt.Errorf("cannot get the zip FaaS endpoint from cloud config")
	}

	return zipEndpoint, nil
}

// 获取签名URL
func GetSignURL(itemBucket string, itemPath string, client *oss.Client) (string, error) {

	// TODO: 阿里云重试机制
	bucket, err := client.Bucket(itemBucket)
	if err != nil {
		log.Println(fmt.Sprintf("Func: GetSignURL Get Client %v Bucket %s Failed %s", client, itemBucket, err.Error()))
		return "", err
	}

	// 请求头信息进行签名
	options := []oss.Option{
		oss.ContentType(common.OSS_DOWNLOAD_CONTENT_TYPE),
	}

	signedURL, err := bucket.SignURL(itemPath, oss.HTTPGet, common.FILE_DOWNLOAD_SIGNURL_TIME, options...)
	if err != nil {
		log.Println(fmt.Sprintf("Func: GetSignURL Get Bucket %s Object %s Failed %s", itemBucket, itemPath, err.Error()))
		return "", err
	}

	// TODO: 优雅一点……这个太暴力了
	signedHttpsURL := "https" + signedURL[4:]
	log.Println("signed url: ", signedHttpsURL)
	return signedHttpsURL, nil
}
