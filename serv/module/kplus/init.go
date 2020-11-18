package kplus

import (
	"net/http"

	. "cm-v5/serv/module"
	// . "cm-v5/schema"
	"github.com/gin-gonic/gin"
	jsoniter "github.com/json-iterator/go"
)

var json = jsoniter.ConfigCompatibleWithStandardLibrary
var mRedis RedisModelStruct
var mRedisKV RedisKVModelStruct

func GetTypeBannerKPlusInit(c *gin.Context) {
	user_id := c.GetString("user_id")
	platform := c.DefaultQuery("platform", "web")
	bannerType := GetTypeBannerKplus(user_id, platform)
	c.JSON(http.StatusOK, FormatResultAPI(http.StatusOK, "", bannerType))
}
