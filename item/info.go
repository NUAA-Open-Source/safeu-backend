package item

import (
	"net/http"

	"a2os/safeu-backend/common"

	"github.com/gin-gonic/gin"
)

type getItemInfoBody struct {
	UserToken string `json:"user_token"`
}

// GetItemInfo 获取文件信息
func GetItemInfo(c *gin.Context) {
	retrieveCode := c.Param("retrieveCode")
	var getItemInfoBody getItemInfoBody
	if common.FuncHandler(c, c.BindJSON(&getItemInfoBody), nil, http.StatusBadRequest, 20301) {
		return
	}
	tokenRedisClient := common.GetUserTokenRedisClient()
	if common.FuncHandler(c, KeyISExistInRedis(getItemInfoBody.UserToken, tokenRedisClient), true, http.StatusUnauthorized, 20201) {
		return
	}
	db := common.GetDB()
	var item Item
	if common.FuncHandler(c, db.Where("re_code = ?", retrieveCode).First(&item).RecordNotFound(), false, http.StatusUnauthorized, 20306) {
		return
	}
	c.JSON(http.StatusOK, gin.H{
		"downCount ": item.DownCount,
		"expireTime": item.ExpiredAt,
	})
	return
}
