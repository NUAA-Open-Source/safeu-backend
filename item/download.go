package item

import (
	"log"
	"net/http"
	"time"

	"a2os/safeu-backend/common"

	"github.com/gin-gonic/gin"
)

func DownloadItems(c *gin.Context) {
	retrieveCode := c.Param("retrieveCode")
	log.Println(c.ClientIP(), " Received download request for ", retrieveCode, " resources")

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

	var itemList []Item
	db.Where("re_code = ? AND (status = ? OR status = ?)", retrieveCode, common.UPLOAD_FINISHED, common.FILE_ACTIVE).Find(&itemList)
	// 单文件下载
	if len(itemList) == 1 {
		url := itemList[0].Host
		c.JSON(http.StatusOK, gin.H{
			"url": url,
		})
		return
	}

	// FIXME: 暂时不支持多文件下载
	c.JSON(http.StatusNotAcceptable, gin.H{
		"error": "Sorry, we haven't support this type service yet",
	})
	return

	// TODO: 多文件 zip 打包
	// TODO: 多文件下载
}
