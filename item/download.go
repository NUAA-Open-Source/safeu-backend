package item

import (
	"bytes"
	"log"
	"net/http"
	"time"

	"a2os/safeu-backend/common"

	"github.com/gin-gonic/gin"
	"github.com/gin-gonic/gin/json"
	"github.com/satori/go.uuid"
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
			"error": "Cannot get the token",
		})
		log.Println(c.ClientIP(), " Cannot get the client token from header")
		return
	}

	log.Println(c.ClientIP(), " Get client token ", clientToken)

	var tokenRecord Token
	if db.Where("token = ?", clientToken).First(&tokenRecord).RecordNotFound() {
		// 无法找到该 token
		c.JSON(http.StatusNotFound, gin.H{
			"error": "Token invalid",
		})
		log.Println(c.ClientIP(), " Invalid token ", clientToken)
		return
	}

	// 检查 token 是否失效
	if !tokenRecord.Valid {
		c.JSON(http.StatusUnauthorized, gin.H{
			"error": "Token invalid",
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
			"error": "Token invalid",
		})
		log.Println(c.ClientIP(), " Expired token ", clientToken)
		return
	}

	// Token 没过期，核对提取码
	tokenRetrieveCode := tokenRecord.RetrieveCode
	if tokenRetrieveCode != retrieveCode {
		// 提取码不正确
		c.JSON(http.StatusUnauthorized, gin.H{
			"error": "Invalid token for this item",
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
			"error": "Cannot find the resource by retrieve code " + retrieveCode,
		})
		log.Println(c.ClientIP(), " resource ", retrieveCode, " not found")
		return
	}

	// 单文件下载
	if len(itemList) == 1 {
		url := itemList[0].Host
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
			"error": "Cannot get the items.",
		})
		return
	}
	// 全量打包下载
	var zipPack Item
	if packRequest.Full {
		// 全量打包下载
		if db.Where("re_code = ? AND (status = ? OR status = ?) AND is_archive = ?", retrieveCode, common.UPLOAD_FINISHED, common.FILE_ACTIVE, true).First(&zipPack).RecordNotFound() {
			// 没有全量打包，进行全量打包并将记录存储到数据库中

			// FIXME: 此处最好通过 key-value 来进行检索 (key = name, value = zip)
			resJson := ZipItemsFaaS(packRequest.ZipItems, retrieveCode, true, common.CloudConfig.FaaS[0].Endpoint)

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
			zipPack.IsArchive = true

			db.Create(&zipPack)
			log.Println("Generated the full files zip package for retrieve code \"", retrieveCode, "\"")

			downloadLink := resJson["host"]

			// 返回压缩包下载链接
			c.JSON(http.StatusOK, gin.H{
				"url": downloadLink,
			})
			log.Println(c.ClientIP(), " Get the zip file url: ", downloadLink)
			return

		}

		// 有全量打包，则直接发送打包文件
		c.JSON(http.StatusOK, gin.H{
			"url": zipPack.Host,
		})
		log.Println(c.ClientIP(), " Full zip pack has generated before, get the zip file url: ", zipPack.Host)
		return
	}

	// 自定义多文件打包下载

	// FIXME: 此处最好通过 key-value 来进行检索 (key = name, value = zip)
	resJson := ZipItemsFaaS(packRequest.ZipItems, retrieveCode, false, common.CloudConfig.FaaS[0].Endpoint)
	log.Println(c.ClientIP(), " Generated the custom zip file for retrieve code \"", retrieveCode, "\"")

	downloadLink := resJson["host"]
	log.Println(c.ClientIP(), " Get the zip file url: ", downloadLink)

	// 返回压缩包路径
	c.JSON(http.StatusOK, gin.H{
		"url": downloadLink,
	})
	return
}

func ZipItemsFaaS(zipItems []ZipItem, retrieveCode string, isFull bool, endpoint string) map[string]string {
	reqJson := map[string]interface{}{
		"re_code": retrieveCode,
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
		log.Println(err)
	}

	var resJson map[string]string
	json.NewDecoder(res.Body).Decode(&resJson)

	return resJson
}
