package item

import (
	"crypto/sha256"
	"encoding/hex"
	"log"
	"net/http"
	"time"

	"a2os/safeu-backend/common"

	"github.com/gin-gonic/gin"
	"github.com/go-redis/redis"
	"github.com/jinzhu/gorm"
)

type ChangePassBody struct {
	Auth      string `json:"auth"`
	UserToken string `json:"user_token"`
}

type ChangeReCodeBody struct {
	NewReCode string `json:"new_re_code"`
	UserToken string `json:"user_token"`
	Auth      string `json:"auth"`
}

type ChangeDownCountBody struct {
	NewDownCount int    `json:"new_down_count"`
	UserToken    string `json:"user_token"`
}

type ChangeExpireTimeBody struct {
	NewExpireTime int    `json:"new_expire_time"`
	UserToken     string `json:"user_token"`
}

// 修改过期时间
func ChangeExpireTime(c *gin.Context) {
	retrieveCode := c.Param("retrieveCode")
	var changeExpireTimeBody ChangeExpireTimeBody
	err := c.BindJSON(&changeExpireTimeBody)
	if err != nil {
		log.Println(err)
		c.JSON(http.StatusInternalServerError, gin.H{
			"message": err,
		})
		return
	}
	// 时间长度检查
	if changeExpireTimeBody.NewExpireTime > common.FILE_MAX_EXIST_TIME || changeExpireTimeBody.NewExpireTime <= 0 {
		c.JSON(http.StatusBadRequest, gin.H{
			"message": "The length of the time is not within the right range",
		})
		return
	}
	tokenRedisClient := common.GetUserTokenRedisClient()
	files, err := tokenRedisClient.SMembers(changeExpireTimeBody.UserToken).Result()
	if len(files) == 0 {
		c.JSON(http.StatusUnauthorized, gin.H{
			"message": "Can`t Find User Token",
		})
		return
	}
	db := common.GetDB()
	// 获取创建时间
	var item Item
	db.Where("name = ? AND status = ? AND re_code = ?", files[0], common.UPLOAD_FINISHED, retrieveCode).First(&item)
	h, _ := time.ParseDuration("1h")
	newTime := item.CreatedAt.Add(time.Duration(changeExpireTimeBody.NewExpireTime) * h)
	for _, value := range files {
		db.Model(&Item{}).Where("name = ? AND status = ? AND re_code = ?", value, common.UPLOAD_FINISHED, retrieveCode).Update(map[string]interface{}{"expired_at": newTime})
	}
	c.JSON(http.StatusOK, gin.H{
		"message": newTime,
	})
	return
}

// 修改密码
func ChangePassword(c *gin.Context) {
	retrieveCode := c.Param("retrieveCode")
	var changePassBody ChangePassBody
	err := c.BindJSON(&changePassBody)
	if err != nil {
		log.Println(err)
		c.JSON(http.StatusInternalServerError, gin.H{
			"message": err,
		})
		return
	}

	tokenRedisClient := common.GetUserTokenRedisClient()
	files, err := tokenRedisClient.SMembers(changePassBody.UserToken).Result()
	// 无文件则未从redis成功读取用户Token 鉴权失败
	if len(files) == 0 {
		c.JSON(http.StatusUnauthorized, gin.H{
			"message": "Can`t Find User Token",
		})
		return
	}
	db := common.GetDB()

	// Auth 为空 重置密码为空
	var isPublic bool
	if changePassBody.Auth == "" {
		isPublic = true
		for _, value := range files {
			db.Model(&Item{}).Where("name = ? AND status = ? AND re_code = ?", value, common.UPLOAD_FINISHED, retrieveCode).Update(map[string]interface{}{"password": "", "is_public": isPublic})
		}
		c.JSON(http.StatusOK, gin.H{
			"message": "ok",
		})
		return
	}

	// Auth 不为空 加盐哈希密码
	hasher := sha256.New()
	hasher.Write([]byte(retrieveCode + changePassBody.Auth))
	hasherSum := hex.EncodeToString(hasher.Sum(nil))
	for _, value := range files {
		db.Model(&Item{}).Where("name = ? AND status = ? AND re_code = ?", value, common.UPLOAD_FINISHED, retrieveCode).Update(map[string]interface{}{"password": hasherSum, "is_public": isPublic})
	}
	c.JSON(http.StatusOK, gin.H{
		"message": "ok",
	})
}

// 修改提取码
func ChangeRecode(c *gin.Context) {
	retrieveCode := c.Param("retrieveCode")
	var changeRecodeBody ChangeReCodeBody
	err := c.BindJSON(&changeRecodeBody)
	if err != nil {
		log.Println(err)
		c.JSON(http.StatusInternalServerError, gin.H{
			"message": err,
		})
		return
	}

	//1 Token 检测
	tokenRedisClient := common.GetUserTokenRedisClient()
	files, err := tokenRedisClient.SMembers(changeRecodeBody.UserToken).Result()
	if len(files) == 0 {
		log.Println("Can`t Find User Token In Redis", changeRecodeBody.UserToken)
		c.JSON(http.StatusUnauthorized, gin.H{
			"message": "Can`t Find User Token",
		})
		return
	}

	//2 合法性检测
	// 2.1 在Redis中判断提取码重复
	reCodeRedisClient := common.GetReCodeRedisClient()
	//if KeyISExistInRedis(changeRecodeBody.NewReCode, reCodeRedisClient) {
	//	//	log.Println("Find reCode Repeat In Redis", changeRecodeBody.NewReCode)
	//	//	c.JSON(http.StatusBadRequest, gin.H{
	//	//		"message": "reCode Repeat",
	//	//	})
	//	//	return
	//	//}
	// 2.2 在DB中判断判断提取码重复
	db := common.GetDB()
	if CheckReCodeRepeatInDB(changeRecodeBody.NewReCode, db) {
		log.Println("Find reCode Repeat In DB", changeRecodeBody.NewReCode)
		c.JSON(http.StatusBadRequest, gin.H{
			"message": "reCode Repeat",
		})
		return
	}

	// 2.3 检测密码
	// 2.4 Auth 检测
	// 2.4 修改数据库
	if !CheckItemHasPass(retrieveCode, db) {
		// 2.3.1 无密码
		// 3.1 直接修改数据库
		for _, value := range files {
			db.Model(&Item{}).Where("name = ? AND status = ? AND re_code = ?", value, common.UPLOAD_FINISHED, retrieveCode).Update(map[string]interface{}{"re_code": changeRecodeBody.NewReCode})
		}
		log.Println("Success Change ReCode", "Previous Recode", retrieveCode, "Now Recode", changeRecodeBody.NewReCode)
		c.JSON(http.StatusOK, gin.H{
			"message": "ok",
		})
		err := reCodeRedisClient.Rename(retrieveCode, changeRecodeBody.NewReCode).Err()
		if err != nil {
			log.Println("reCodeRedisClient Rename err", "old key", retrieveCode, "new key", changeRecodeBody.NewReCode)
		}
		return
	}

	// 2.3.2有密码
	// 2.4.1 Auth 未填充
	if changeRecodeBody.Auth == "" {
		log.Println("Item had password,but not Auth give  Previous Recode:", retrieveCode)
		c.JSON(http.StatusBadRequest, gin.H{
			"message": "require Auth",
		})
		return
	}
	// 2.4.2 Auth已填充
	// 3.2 加盐哈希密码存储
	hasher := sha256.New()
	hasher.Write([]byte(changeRecodeBody.NewReCode + changeRecodeBody.Auth))
	hasherSum := hex.EncodeToString(hasher.Sum(nil))
	for _, value := range files {
		db.Model(&Item{}).Where("name = ? AND status = ? AND re_code = ?", value, common.UPLOAD_FINISHED, retrieveCode).Update(map[string]interface{}{"re_code": changeRecodeBody.NewReCode, "password": hasherSum, "is_public": false})
	}
	log.Println("Success Change ReCode With New Password", "Previous Recode", retrieveCode, "Now Recode", changeRecodeBody.NewReCode)
	c.JSON(http.StatusOK, gin.H{
		"message": "ok",
	})
	err = reCodeRedisClient.Rename(retrieveCode, changeRecodeBody.NewReCode).Err()
	if err != nil {
		log.Println("reCodeRedisClient Rename err", "old key", retrieveCode, "new key", changeRecodeBody.NewReCode)
	}
	return
}

// 修改下载次数
func ChangeDownCount(c *gin.Context) {
	retrieveCode := c.Param("retrieveCode")
	var changeDownCount ChangeDownCountBody
	err := c.BindJSON(&changeDownCount)
	if err != nil {
		log.Println(err)
		c.JSON(http.StatusInternalServerError, gin.H{
			"message": err,
		})
		return
	}
	tokenRedisClient := common.GetUserTokenRedisClient()
	files, err := tokenRedisClient.SMembers(changeDownCount.UserToken).Result()
	if len(files) == 0 {
		c.JSON(http.StatusUnauthorized, gin.H{
			"message": "Can`t Find User Token",
		})
		return
	}
	db := common.GetDB()
	for _, value := range files {
		db.Model(&Item{}).Where("name = ? AND status = ? AND re_code = ?", value, common.UPLOAD_FINISHED, retrieveCode).Update(map[string]interface{}{"down_count": changeDownCount.NewDownCount})
	}
	c.JSON(http.StatusOK, gin.H{
		"message": "ok",
	})
}

func KeyISExistInRedis(str string, client *redis.Client) bool {
	if client.Exists(str).Val() == 0 {
		return false
	}
	return true
}

func CheckReCodeRepeatInDB(str string, db *gorm.DB) bool {
	if db.Where("re_code = ?", str).First(&Item{}).RecordNotFound() {
		return false
	}
	return true
}

func CheckItemHasPass(str string, db *gorm.DB) bool {
	var item Item
	db.Where("re_code = ?", str).First(&item)
	if item.Password != "" {
		return true
	}
	return false
}
