package shp

import (
	"net/http"

	. "cm-v5/serv/module"
	"github.com/gin-gonic/gin"
	jsoniter "github.com/json-iterator/go"
)

var mRedis RedisModelStruct
var mRedisKV RedisKVModelStruct
var json = jsoniter.ConfigCompatibleWithStandardLibrary

func SmartHubPreviewInit(c *gin.Context) {

	userId := c.GetString("user_id")
	cacheActive := StringToBool(c.DefaultQuery("cache", "true"))

	dataSHP := GetDataSamSungPreview(userId, cacheActive)

	c.JSON(http.StatusOK, FormatResultAPI(http.StatusOK, "", dataSHP))

}
