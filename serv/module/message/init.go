package message

import (
	"net/http"

	. "cm-v5/serv/module"
	. "cm-v5/schema"
	"github.com/gin-gonic/gin"
	jsoniter "github.com/json-iterator/go"
)

var mRedisUSC RedisUSCModelStruct
var mRedis RedisModelStruct
var json = jsoniter.ConfigCompatibleWithStandardLibrary

func CountMessageByUser(c *gin.Context) {
	cacheActive := StringToBool(c.DefaultQuery("cache", "true"))
	platform := c.DefaultQuery("platform", "web")
	user_id := c.GetString("user_id")

	data, err := CountMessage(user_id, platform, cacheActive)
	if err != nil {
		Sentry_log(err)
		c.JSON(http.StatusBadRequest, FormatResultAPI(http.StatusBadRequest, err.Error(), ""))
		return
	}
	c.JSON(http.StatusOK, FormatResultAPI(http.StatusOK, "", data))
}

func ListMessageByUser(c *gin.Context) {
	cacheActive := StringToBool(c.DefaultQuery("cache", "true"))
	platform := c.DefaultQuery("platform", "web")
	user_id := c.GetString("user_id")
	page, err := StringToInt(c.DefaultQuery("page", "0"))
	if err != nil || page > PAGE_MAX {
		page = 0
	}
	limit, err := StringToInt(c.DefaultQuery("limit", "10"))
	if err != nil || limit > LIMIT_MAX {
		limit = 10
	}

	data, err := GetMessageByUser(user_id, platform, page, limit, cacheActive)
	if err != nil {
		Sentry_log(err)
		c.JSON(http.StatusBadRequest, FormatResultAPI(http.StatusBadRequest, err.Error(), ""))
		return
	}
	c.JSON(http.StatusOK, FormatResultAPI(http.StatusOK, "", data))
}

func ActionMessageByUser(c *gin.Context) {
	message_id := c.PostForm("message_id")
	user_id := c.GetString("user_id")
	action := c.PostForm("action")
	listAction := map[string]bool{
		"mark_all":   true,
		"unmark_all": true,
		"mark":       true,
		"unmark":     true,
		"delete":     true,
	}
	if _, ok := listAction[action]; !ok {
		c.JSON(http.StatusBadRequest, FormatResultAPI(http.StatusBadRequest, "Action: mark_all/unmark_all/mark/unmark/delete", ""))
		return
	}

	data, err := ActionMessage(user_id, message_id, action)
	if err != nil {
		Sentry_log(err)
		c.JSON(http.StatusBadRequest, FormatResultAPI(http.StatusBadRequest, err.Error(), ""))
		return
	}
	c.JSON(http.StatusOK, FormatResultAPI(http.StatusOK, "", data))
}

func SaveDataPushTokenByUser(c *gin.Context) {
	push_token := c.PostForm("push_token")
	platform := c.DefaultQuery("platform", "web")
	access_token := c.GetString("access_token")
	user_id := c.GetString("user_id")

	if push_token == "" {
		c.JSON(http.StatusBadRequest, FormatResultAPI(http.StatusBadRequest, "Required: push_token", ""))
		return
	}

	err := SaveDataPushToken(push_token, access_token, user_id, platform)
	if err != nil {
		c.JSON(http.StatusBadRequest, FormatResultAPI(http.StatusBadRequest, err.Error(), ""))
		return
	}

	c.JSON(http.StatusOK, FormatResultAPI(http.StatusOK, "", "Success"))
}
