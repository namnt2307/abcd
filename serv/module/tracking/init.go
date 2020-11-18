package tracking

import (
	"net/http"

	. "cm-v5/serv/module"
	. "cm-v5/schema"
	"github.com/gin-gonic/gin"
	jsoniter "github.com/json-iterator/go"
)

// var mRedisUSC RedisUSCModelStruct
// var mRedisKV RedisKVModelStruct
var json = jsoniter.ConfigCompatibleWithStandardLibrary
var mRedisUSC RedisUSCModelStruct

func TrackingWatchByUser(c *gin.Context) {
	user_id := c.GetString("user_id")
	var trackingWatchRequestData TrackingWatchRequestDataStruct
	if err := c.ShouldBindJSON(&trackingWatchRequestData); err != nil {
		c.JSON(http.StatusBadRequest, FormatResultAPI(http.StatusBadRequest, err.Error(), ""))
		return
	}

	data, err := TrackingWatch(c, user_id, trackingWatchRequestData)
	if err != nil {
		Sentry_log(err)
		c.JSON(http.StatusBadRequest, FormatResultAPI(http.StatusBadRequest, err.Error(), ""))
		return
	}
	c.JSON(http.StatusOK, FormatResultAPI(http.StatusOK, "", data))
}

func ListWatchingByUser(c *gin.Context) {
	user_id := c.GetString("user_id")
	platform := c.DefaultQuery("platform", "web")
	page, err := StringToInt(c.DefaultQuery("page", "0"))
	if err != nil || page > PAGE_MAX {
		page = 0
	}
	limit, err := StringToInt(c.DefaultQuery("limit", "15"))
	if err != nil || limit > LIMIT_MAX {
		limit = 15
	}

	data, err := ListWatching(user_id, platform, page, limit)
	if err != nil {
		Sentry_log(err)
		c.JSON(http.StatusBadRequest, FormatResultAPI(http.StatusBadRequest, err.Error(), ""))
		return
	}
	c.JSON(http.StatusOK, FormatResultAPI(http.StatusOK, "", data))
}

func RemoveWatchingByUser(c *gin.Context) {
	user_id := c.GetString("user_id")
	// isAll := c.DefaultQuery("is_all", "true")

	data, err := RemoveWatching(user_id)
	if err != nil {
		Sentry_log(err)
		c.JSON(http.StatusBadRequest, FormatResultAPI(http.StatusBadRequest, err.Error(), ""))
		return
	}
	c.JSON(http.StatusOK, FormatResultAPI(http.StatusOK, "", data))
}

func TrackingNotification(c *gin.Context) {
	c.JSON(http.StatusOK, FormatResultAPI(http.StatusOK, "", "Success"))
}
