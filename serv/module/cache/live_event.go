package cache

import (
	"net/http"

	. "cm-v5/serv/module"
	live_event "cm-v5/serv/module/live_event"
	"github.com/gin-gonic/gin"
)

func UpdateCacheLiveEvent(c *gin.Context) {
	var EventID string = c.Param("event_id")
	for _, val := range platforms {
		live_event.GetDetailByID(EventID, val, false)
	}
	c.JSON(http.StatusOK, FormatResultAPI(http.StatusOK, "", "Update cache done"))
}

func UpdateCacheAllLiveEvent(c *gin.Context) {
	live_event.GetAllLiveEventByMongoDB()
	c.JSON(http.StatusOK, FormatResultAPI(http.StatusOK, "", "Update cache done"))
}
