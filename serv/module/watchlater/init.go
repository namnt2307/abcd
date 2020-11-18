package watchlater

import (
	"net/http"

	. "cm-v5/serv/module"
	. "cm-v5/schema"
	"github.com/gin-gonic/gin"
	jsoniter "github.com/json-iterator/go"
)

var mRedisUSC RedisUSCModelStruct
var mRedisKV RedisKVModelStruct
var json = jsoniter.ConfigCompatibleWithStandardLibrary

func GetWatchlaterByUser(c *gin.Context) {
	cacheActive := StringToBool(c.DefaultQuery("cache", "true"))
	user_id := c.GetString("user_id")
	platform := c.DefaultQuery("platform", "web")
	page, err := StringToInt(c.DefaultQuery("page", "0"))
	if err != nil || page > PAGE_MAX {
		page = 0
	}
	limit, err := StringToInt(c.DefaultQuery("limit", "30"))
	if err != nil || limit > LIMIT_MAX {
		limit = 43
	}

	data, err := GetWatchlater(user_id, platform, page, limit, cacheActive)
	if err != nil {
		Sentry_log(err)
		c.JSON(http.StatusBadRequest, FormatResultAPI(http.StatusBadRequest, err.Error(), ""))
		return
	}
	c.JSON(http.StatusOK, FormatResultAPI(http.StatusOK, "", data))
}

func AddWatchlaterByUser(c *gin.Context) {
	content_id := c.PostForm("content_id")
	user_id := c.GetString("user_id")
	data, err := AddWatchlater(c, user_id, content_id)
	if err != nil {
		Sentry_log(err)
		c.JSON(http.StatusBadRequest, FormatResultAPI(http.StatusBadRequest, err.Error(), ""))
		return
	}
	c.JSON(http.StatusOK, FormatResultAPI(http.StatusOK, "", data))
}
