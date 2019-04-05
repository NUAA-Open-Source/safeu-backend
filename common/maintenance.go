package common

import (
	"log"
	"net/http"

	"a2os/behavior/common"

	"github.com/gin-gonic/gin"
)

func MaintenanceHandling() gin.HandlerFunc {
	return func(c *gin.Context) {
		if MAINTENANCE {
			c.JSON(http.StatusServiceUnavailable, gin.H{
				"err_code": 10008,
				"message":  common.Errors[10008],
			})
			log.Println(c.ClientIP(), "Maintenance mode is open")
			c.Abort()
		}

		c.Next()
	}
}
