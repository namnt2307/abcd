package rating

import (
	"net/http"

	. "cm-v5/serv/module"
	"github.com/gin-gonic/gin"
	jsoniter "github.com/json-iterator/go"
)

var mRedis RedisModelStruct
var mRedisUSC RedisUSCModelStruct
var json = jsoniter.ConfigCompatibleWithStandardLibrary

func ActionRatingByUser(c *gin.Context) {
	userId := c.GetString("user_id")
	contentId := c.Param("content_id")
	platform := c.DefaultQuery("platform", "web")
	point, err := StringToInt(c.Param("point"))
	if err != nil {
		c.JSON(http.StatusBadRequest, FormatResultAPI(http.StatusBadRequest, err.Error(), ""))
		return
	}

	data, err := RatingByUser(userId, contentId, platform , point)
	if err != nil {
		c.JSON(http.StatusBadRequest, FormatResultAPI(http.StatusBadRequest, err.Error(), ""))
		return
	}
	c.JSON(http.StatusOK, FormatResultAPI(http.StatusOK, "", data))
}
