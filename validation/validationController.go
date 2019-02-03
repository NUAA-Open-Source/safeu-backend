package validation

import (
	"net/http"
	"time"
	"fmt"
	"crypto/md5"
	"strconv"
	"io"

	"a2os/safeu-backend/common"
	"a2os/safeu-backend/item"

	"github.com/gin-gonic/gin"
)

type ValiPass struct {
	Password	string	`form:"password" json:"password"`
}

func GenerateTokenByMd5() string {

	curTime := time.Now().Unix()
	h := md5.New()
	io.WriteString(h, strconv.FormatInt(curTime, 10))

	token := fmt.Sprintf("%x", h.Sum(nil))

	return token
}

func Validation (c *gin.Context) {
	retrieveCode := c.Param("retrieveCode")

	db := common.GetDB()

	// 是否存在文件
	var curItem item.Item
	if db.Where("re_code = ? AND (status = 2 OR status = 3)", retrieveCode).First(&curItem).RecordNotFound() {
		c.JSON(http.StatusNotFound, gin.H{
			"error": "Cannot find resource via this retrieve code.",
		})
		return
	}

	// 检查文件是否为公有
	if curItem.IsPublic {
		// 文件公有，生成 Token 并返回文件列表
		token := GenerateTokenByMd5()
		var itemList []item.Item
		db.Where("re_code = ? AND (status = 2 OR status = 3)", retrieveCode).Find(&itemList)

		tokenRecord := Token{Token: token, RetrieveCode: retrieveCode, Valid: true, ExpiredAt: time.Now().Add(15*time.Minute)}  // Token 有效期为 15 min
		db.Create(&tokenRecord)

		c.JSON(http.StatusOK, gin.H{
			"token": token,
			"items": itemList,
		})
		return
	}

	// TODO: 验证密码

	// TODO: 密码不正确/密码缺失返回 401

	// TODO: 密码正确生成 Token 并返回文件列表

}
