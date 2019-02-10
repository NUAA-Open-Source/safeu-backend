package item

import (
	"a2os/safeu-backend/common"
	"crypto/sha256"
	"encoding/hex"
	"github.com/gin-gonic/gin"
	"log"
	"net/http"
)

type ChangePassBody struct {
	Auth      string `json:"auth"`
	UserToken string `json:"user_token"`
}

type ChangeReCodeBody struct {
	NewReCode string `json:"new_re_code"`
	UserToken string `json:"user_token"`
}

type ChangeDownCountBody struct {
	NewDownCount int    `json:"new_down_count"`
	UserToken    string `json:"user_token"`
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
	var isPublic bool
	if changePassBody.Auth == "" {
		isPublic = true
	}
	hasher := sha256.New()
	hasher.Write([]byte(retrieveCode + changePassBody.Auth))
	hasherSum := hex.EncodeToString(hasher.Sum(nil))

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
	tokenRedisClient := common.GetUserTokenRedisClient()
	files, err := tokenRedisClient.SMembers(changeRecodeBody.UserToken).Result()
	if len(files) == 0 {
		c.JSON(http.StatusUnauthorized, gin.H{
			"message": "Can`t Find User Token",
		})
		return
	}
	db := common.GetDB()
	for _, value := range files {
		db.Model(&Item{}).Where("name = ? AND status = ? AND re_code = ?", value, common.UPLOAD_FINISHED, retrieveCode).Update(map[string]interface{}{"re_code": changeRecodeBody.NewReCode})
	}
	c.JSON(http.StatusOK, gin.H{
		"message": "ok",
	})
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
