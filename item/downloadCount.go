package item

import (
	"a2os/safeu-backend/common"

	"github.com/gin-gonic/gin"
	"net/http"
	"log"
)

func DownloadCount(c *gin.Context) {
	retrieveCode := c.Param("retrieveCode")
	bucket := c.Query("bucket")
	path := c.Query("path")

	db := common.GetDB()

	var itemList []Item
	db.Where("re_code = ? AND (status = ? OR status = ?) AND archive_type = ?", retrieveCode, common.UPLOAD_FINISHED, common.FILE_ACTIVE, common.ARCHIVE_NULL).Find(&itemList)

	// 提取码错误
	if len(itemList) == 0 {
		c.JSON(http.StatusNotFound, gin.H{
			"error": "Cannot find the resource.",
		})
		log.Println(c.ClientIP(), " resource ", retrieveCode, " not found, cannot change download count")
		return
	}

	// 单文件
	if len(itemList) == 1 {
		singleItem := itemList[0]

		// 无限下载
		if singleItem.DownCount == common.INFINITE_DOWNLOAD {
			c.String(http.StatusOK, "OK")
			return
		}

		// 下载次数--
		singleItem.DownCount -= 1
		if singleItem.DownCount <= 0 {
			// 删除文件
			err := DeleteItem(singleItem.Bucket, singleItem.Path)
			if err != nil {
				log.Println("Cannot delete item in bucket ", singleItem.Bucket, ", path ", singleItem.Path)
			}

			// 删除数据库记录
			db.Delete(&singleItem)

			c.String(http.StatusOK, "DELETED")
			return
		}

		// 更新下载次数
		db.Update(&singleItem)
		c.String(http.StatusOK, "MINUS")
		return
	}

	// TODO: 多文件

}
