package item

import (
	"net/http"

	"a2os/safeu-backend/common"

	"github.com/gin-gonic/gin"
)

type GetItemInfoBody struct {
	UserToken string `json:"user_token"`
}

func GetItemInfo(c *gin.Context) {
	retrieveCode := c.Param("retrieveCode")
	var getItemInfoBody GetItemInfoBody
	if err := c.BindJSON(&getItemInfoBody); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"err_code": 1,
			"message":  common.Errors[1],
		})
		return
	}
	tokenRedisClient := common.GetUserTokenRedisClient()
	if !KeyISExistInRedis(getItemInfoBody.UserToken, tokenRedisClient) {
		c.JSON(http.StatusBadRequest, gin.H{
			"err_code": 3,
			"message":  common.Errors[3],
		})
		return
	}
	db := common.GetDB()
	var item Item
	if db.Where("re_code = ?", retrieveCode).First(&item).RecordNotFound() {
		c.JSON(http.StatusBadRequest, gin.H{
			"err_code": 6,
			"message":  common.Errors[6],
		})
		return
	}
	c.JSON(http.StatusOK, gin.H{
		"downCount ": item.DownCount,
		"expireTime": item.ExpiredAt,
	})
	return
}
