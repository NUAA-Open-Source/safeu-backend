package event

import (
	"log"
	"net/http"

	"a2os/safeu-backend/common"

	"github.com/gin-gonic/gin"
)

func ReportEvent(c *gin.Context) {
	eventName := c.Query("event")
	from := c.Query("from")

	if eventName == "" {
		// return 400 Bad Request
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Cannot get the event name from request.",
		})
		log.Println(c.ClientIP(), "Cannot get the event name from request")
		return
	}

	newEvent := Event{
		Name: eventName,
		From: from,
	}

	db := common.GetDB()
	db.Create(&newEvent)
	log.Println(c.ClientIP(), "Get front-end event", eventName, "from", from)

	c.JSON(http.StatusOK, gin.H{
		"msg": "success",
	})
	return
}
