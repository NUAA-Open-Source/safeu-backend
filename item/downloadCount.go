package item

import (
	"log"
	"net/http"

	"a2os/safeu-backend/common"

	"github.com/gin-gonic/gin"
)

func DownloadCount(c *gin.Context) {
	retrieveCode := c.Param("retrieveCode")
	// 为文件组生命周期准备
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

	// ------- 文件组生命周期

	isDelete := false
	for _, item := range itemList {
		// 检查是否为当前文件，次数无限直接返回
		if item.DownCount == common.INFINITE_DOWNLOAD && item.Bucket == bucket && item.Path == path {
			c.String(http.StatusOK, "OK")
			return
		}

		// 当前文件下载次数--
		if item.DownCount != common.INFINITE_DOWNLOAD && item.Bucket == bucket && item.Path == path {
			item.DownCount -= 1
			db.Model(&item).Update("down_count", item.DownCount)
			log.Println(c.ClientIP(), " Item ", item.ID, " remain downloadable count: ", item.DownCount)
		}

		if item.DownCount == common.INFINITE_DOWNLOAD {
			continue
		} else if item.DownCount <= 0 {
			// 删除文件组中的过期文件
			err := DeleteItem(item.Bucket, item.Path)
			if err != nil {
				log.Println("Cannot delete item in bucket ", item.Bucket, ", path ", item.Path)
			}

			// 删除数据库记录
			db.Delete(&item)
			log.Println(c.ClientIP(), " Deleted item id ", item.ID, " by retrieve code \"", retrieveCode, "\"")
			isDelete = true
		}
	}

	// 重新检索，查看删除后的记录
	if isDelete {
		db.Where("re_code = ? AND (status = ? OR status = ?) AND archive_type = ?", retrieveCode, common.UPLOAD_FINISHED, common.FILE_ACTIVE, common.ARCHIVE_NULL).Find(&itemList)
	}

	// 有效文件数 <= 1，清除所有压缩包
	if len(itemList) <= 1 {
		var deleteItems []Item
		db.Where("re_code = ? AND (status = ? OR status = ?) AND archive_type != ?", retrieveCode, common.UPLOAD_FINISHED, common.FILE_ACTIVE, common.ARCHIVE_NULL).Find(&deleteItems)

		for _, deleteItem := range deleteItems {
			// 删除压缩包
			err := DeleteItem(deleteItem.Bucket, deleteItem.Path)
			if err != nil {
				log.Println("Cannot delete item in bucket ", deleteItem.Bucket, ", path ", deleteItem.Path)
			}

			// 删除压缩包的数据库记录
			db.Delete(&deleteItem)
		}

		log.Println(c.ClientIP(), " Clean all the zip package for retrieve code \"", retrieveCode, "\"")
	}

	if isDelete {
		c.String(http.StatusOK, "DELETED")
	} else {
		c.String(http.StatusOK, "MINUS")
	}

	return
}
