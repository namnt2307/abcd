package config

import (
	"net/http"

	. "cm-v5/serv/module"
	"github.com/gin-gonic/gin"
	jsoniter "github.com/json-iterator/go"
)

var mRedis RedisModelStruct
var mRedisKV RedisKVModelStruct
var json = jsoniter.ConfigCompatibleWithStandardLibrary

func init() {

}
func GetConfigByKeyInit(c *gin.Context) {
	key := c.DefaultQuery("key", "")
	platform := c.DefaultQuery("platform", "web")
	cacheActive := StringToBool(c.DefaultQuery("cache", "true"))
	if key == "" {
		c.JSON(http.StatusBadRequest, FormatResultAPI(http.StatusBadRequest, "key empty", ""))
		return
	}

	result, err := GetConfigByKey(key, platform, cacheActive)
	if err != nil {
		c.JSON(http.StatusBadRequest, FormatResultAPI(http.StatusBadRequest, err.Error(), ""))
		return
	}

	c.JSON(http.StatusOK, FormatResultAPI(1, "", result))
}
