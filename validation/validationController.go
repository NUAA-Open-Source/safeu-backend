package validation

import (
	"net/http"
	"time"
	"fmt"
	"strconv"
	"io"
	"log"
	"crypto/md5"
	"crypto/sha256"
	"encoding/hex"

	"a2os/safeu-backend/common"
	"a2os/safeu-backend/item"

	"github.com/gin-gonic/gin"
)

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
	log.Println(c.ClientIP() + " Received validation request for " + retrieveCode + " resources")

	db := common.GetDB()

	// 是否存在文件
	var curItem item.Item
	if db.Where("re_code = ? AND (status = 2 OR status = 3)", retrieveCode).First(&curItem).RecordNotFound() {
		c.JSON(http.StatusNotFound, gin.H{
			"error": "Cannot find resource via this retrieve code.",
		})
		log.Println(c.ClientIP() + " Cannot find resource via the retrieve code " + retrieveCode)
		return
	}

	// 检查文件是否为公有
	if curItem.IsPublic {
		// 文件公有，生成 Token 并返回文件列表
		token := GenerateTokenByMd5()
		var itemList []item.Item
		db.Where("re_code = ? AND (status = 2 OR status = 3)", retrieveCode).Find(&itemList)

		tokenRecord := Token{Token: token, RetrieveCode: retrieveCode, Valid: true, ExpiredAt: time.Now().Add(15 * time.Minute)} // Token 有效期为 15 min
		db.Create(&tokenRecord)
		log.Println(c.ClientIP() + " Generated token " + token + " for retrieve code " + retrieveCode)

		c.JSON(http.StatusOK, gin.H{
			"token": token,
			"items": itemList,
		})
		return
	}

	// 验证密码
	// 密码加密方式：SHA256( retrieveCode + SHA256(password) )
	var valiPass ValiPass
	if err := c.ShouldBindJSON(&valiPass); err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{
			"error": "Cannot get the password",
		})
		log.Println(c.ClientIP() + " Cannot get the password from client, return 401 Unauthorized")
		return
	}

	SHA256Password := valiPass.Password
	refPassword := curItem.Password
	hasher := sha256.New()
	hasher.Write([]byte(retrieveCode + SHA256Password))
	hasherSum := hex.EncodeToString(hasher.Sum(nil))

	//log.Println("[DEBUG] refPassword = " + refPassword)
	//log.Println("[DEBUG] password = " + SHA256Password)
	//log.Println("[DEBUG] hasher sum = " + hasherSum)

	// 密码不正确/密码缺失返回 401
	if hasherSum != refPassword {
		c.JSON(http.StatusUnauthorized, gin.H{
			"error": "The password is not correct",
		})
		log.Println(c.ClientIP() + " The password is not correct, return 401 Unauthorized")
		return
	}

	// 密码正确生成 Token 并返回文件列表
	token := GenerateTokenByMd5()
	var itemList []item.Item
	db.Where("re_code = ? AND (status = 2 OR status = 3)", retrieveCode).Find(&itemList)

	tokenRecord := Token{Token: token, RetrieveCode: retrieveCode, Valid: true, ExpiredAt: time.Now().Add(15 * time.Minute)} // Token 有效期为 15 min
	db.Create(&tokenRecord)
	log.Println(c.ClientIP() + " Generated token " + token + " for retrieve code " + retrieveCode)

	c.JSON(http.StatusOK, gin.H{
		"token": token,
		"items": itemList,
	})
	return

}
