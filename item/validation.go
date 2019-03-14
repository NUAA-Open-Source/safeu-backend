package item

import (
	"crypto/md5"
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"io"
	"log"
	"net/http"
	"strconv"
	"time"

	"a2os/safeu-backend/common"

	"github.com/gin-gonic/gin"
	"github.com/jinzhu/gorm"
)

type ResponseItem struct {
	Name         string    `json:"name"`
	OriginalName string    `json:"original_name"`
	DownCount    int       `json:"down_count"`
	Type         string    `json:"type"`
	Protocol     string    `json:"protocol"`
	Bucket       string    `json:"bucket"`
	Endpoint     string    `json:"endpoint"`
	Path         string    `json:"path"`
	ExpiredAt    time.Time `json:"expired_at"`
}

type ValiPass struct {
	Password string `form:"password" json:"password" binding:"required"`
}

func GenerateTokenByMd5() string {

	curTime := time.Now().Unix()
	h := md5.New()
	io.WriteString(h, strconv.FormatInt(curTime, 10))

	token := fmt.Sprintf("%x", h.Sum(nil))

	return token
}

func Validation(c *gin.Context) {
	retrieveCode := c.Param("retrieveCode")
	log.Println(c.ClientIP(), " Received validation request for ", retrieveCode, " resources")

	db := common.GetDB()

	// 是否存在文件
	var curItem Item
	if db.Where("re_code = ? AND (status = ? OR status = ?) AND archive_type = ?", retrieveCode, common.UPLOAD_FINISHED, common.FILE_ACTIVE, common.ARCHIVE_NULL).First(&curItem).RecordNotFound() {
		c.JSON(http.StatusNotFound, gin.H{
			"err_code": 10006,
			"message":  common.Errors[10006],
		})
		log.Println(c.ClientIP(), " Cannot find resource via the retrieve code ", retrieveCode)
		return
	}

	// 检查文件是否为公有
	if curItem.IsPublic {

		// 检查文件可下载次数和有效时间
		itemList, err := CheckDownCountAndExpiredTime(db, retrieveCode)
		if err != nil {
			log.Println(c.ClientIP(), " check down count and expired time failed")
			c.JSON(http.StatusInternalServerError, gin.H{
				"err_code": 10001,
				"message":  common.Errors[10001],
			})
			return
		}

		// 文件过期/无下载次数被清空，返回 404
		if len(itemList) == 0 {
			c.JSON(http.StatusNotFound, gin.H{
				"err_code": 10006,
				"message":  common.Errors[10006],
			})
			log.Println(c.ClientIP(), " Cannot find resource via the retrieve code ", retrieveCode)
			return
		}

		// 文件公有，生成 Token 并返回文件列表
		token := GenerateTokenByMd5()

		tokenRecord := Token{Token: token, RetrieveCode: retrieveCode, Valid: true, ExpiredAt: time.Now().Add((time.Duration)(common.TOKEN_VALID_MINUTES) * time.Minute)} // Token 有效期
		db.Create(&tokenRecord)
		log.Println(c.ClientIP(), " Generated token ", token, " for retrieve code ", retrieveCode)

		// 加工 itemList
		responseItemList := GetResponseItemList(itemList)
		c.JSON(http.StatusOK, gin.H{
			"token": token,
			"items": responseItemList,
		})
		return
	}

	// ------ 私有文件
	// 验证密码
	// 密码加密方式：SHA256( retrieveCode, SHA256(password) )
	var valiPass ValiPass
	if err := c.ShouldBindJSON(&valiPass); err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{
			"err_code": 10005,
			"message":  common.Errors[10005],
		})
		log.Println(c.ClientIP(), " Cannot get the password from client, return 401 Unauthorized")
		return
	}

	SHA256Password := valiPass.Password
	refPassword := curItem.Password
	hasher := sha256.New()
	hasher.Write([]byte(retrieveCode + SHA256Password))
	hasherSum := hex.EncodeToString(hasher.Sum(nil))

	//log.Println("[DEBUG] refPassword = ", refPassword)
	//log.Println("[DEBUG] password = ", SHA256Password)
	//log.Println("[DEBUG] hasher sum = ", hasherSum)

	// 密码不正确/密码缺失返回 401
	if hasherSum != refPassword {
		c.JSON(http.StatusUnauthorized, gin.H{
			"err_code": 20501,
			"message":  common.Errors[20501],
		})
		log.Println(c.ClientIP(), " The password is not correct, return 401 Unauthorized")
		return
	}

	// 检查文件可下载次数和有效时间
	itemList, err := CheckDownCountAndExpiredTime(db, retrieveCode)
	if err != nil {
		log.Println(c.ClientIP(), " check down count and expired time failed")
		c.JSON(http.StatusInternalServerError, gin.H{
			"err_code": 10001,
			"message":  common.Errors[10001],
		})
		return
	}

	// 文件过期/无下载次数被清空，返回 404
	if len(itemList) == 0 {
		c.JSON(http.StatusNotFound, gin.H{
			"err_code": 10006,
			"message":  common.Errors[10006],
		})
		log.Println(c.ClientIP(), " Cannot find resource via the retrieve code ", retrieveCode)
		return
	}

	// 密码正确生成 Token 并返回文件列表
	token := GenerateTokenByMd5()

	tokenRecord := Token{
		Token:        token,
		RetrieveCode: retrieveCode,
		Valid:        true,
		ExpiredAt:    time.Now().Add((time.Duration)(common.TOKEN_VALID_MINUTES) * time.Minute),
	} // Token 有效期

	db.Create(&tokenRecord)
	log.Println(c.ClientIP(), " Generated token ", token, " for retrieve code ", retrieveCode)

	// 加工 itemList
	responseItemList := GetResponseItemList(itemList)
	c.JSON(http.StatusOK, gin.H{
		"token": token,
		"items": responseItemList,
	})
	return
}

func CheckDownCountAndExpiredTime(db *gorm.DB, retrieveCode string) ([]Item, error) {

	var itemList []Item
	db.Where("re_code = ? AND (status = ? OR status = ?) AND archive_type = ?", retrieveCode, common.UPLOAD_FINISHED, common.FILE_ACTIVE, common.ARCHIVE_NULL).Find(&itemList)

	itemListCount := len(itemList)
	isExpired := false
	deleted := 0
	for i := range itemList {
		j := i - deleted
		item := itemList[j]

		// 检查文件有效时间
		if item.ExpiredAt.Before(time.Now()) {
			// 文件已过期
			isExpired = true
			// 清除文件
			err := DeleteItem(item.Bucket, item.Path)
			if err != nil {
				log.Println("Cannot delete item in bucket ", item.Bucket, ", path ", item.Path)
				return nil, err
			}

			// 删除数据库记录
			db.Delete(&item)
			common.DeleteRedisRecodeFromRecode(item.ReCode)
			itemList = append(itemList[:j], itemList[j+1:]...)
			deleted++
			continue
		}

		// 检查剩余下载次数
		// 若为 0 次则返回 410 Gone 并删除文件
		if item.DownCount <= 0 && item.DownCount != common.INFINITE_DOWNLOAD {
			// 删除文件
			err := DeleteItem(item.Bucket, item.Path)
			if err != nil {
				log.Println("Cannot delete item in bucket ", item.Bucket, ", path ", item.Path)
				return nil, err
			}

			// 删除数据库记录
			db.Delete(&item)
			common.DeleteRedisRecodeFromRecode(item.ReCode)
			itemList = append(itemList[:j], itemList[j+1:]...)
			deleted++
			continue
		}

	}

	// 若为文件组，则需要删除压缩包
	if isExpired && itemListCount > 1 {
		var deleteZipList []Item
		db.Where("re_code = ? AND (status = ? OR status = ?) AND archive_type != ?", retrieveCode, common.UPLOAD_FINISHED, common.FILE_ACTIVE, common.ARCHIVE_NULL).Find(&deleteZipList)

		for _, deleteItem := range deleteZipList {
			// 删除压缩包
			err := DeleteItem(deleteItem.Bucket, deleteItem.Path)
			if err != nil {
				log.Println("Cannot delete zip in bucket ", deleteItem.Bucket, ", path ", deleteItem.Path)
				return nil, err
			}

			// 删除压缩包的数据库记录
			db.Delete(&deleteItem)
		}
	}

	return itemList, nil
}

func GetResponseItemList(itemList []Item) []ResponseItem {
	var responseItemList []ResponseItem
	for _, item := range itemList {
		responseItem := ResponseItem{
			Protocol:     item.Protocol,
			Bucket:       item.Bucket,
			Endpoint:     item.Endpoint,
			Path:         item.Path,
			OriginalName: item.OriginalName,
			Name:         item.Name,
			DownCount:    item.DownCount,
			Type:         item.Type,
			ExpiredAt:    item.ExpiredAt,
		}
		responseItemList = append(responseItemList, responseItem)
	}

	return responseItemList
}
