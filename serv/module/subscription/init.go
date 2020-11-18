package subscription

import (
	"net/http"

	. "cm-v5/serv/module"
	"github.com/gin-gonic/gin"
	jsoniter "github.com/json-iterator/go"
)

var json = jsoniter.ConfigCompatibleWithStandardLibrary
var mRedisUSC RedisUSCModelStruct

func ClearCacheSubscription(c *gin.Context) {
	userId := c.PostForm("user_id")
	if userId == "" {
		c.JSON(http.StatusBadRequest, FormatResultAPI(http.StatusBadRequest, "", ""))
		return
	}
	var Sub SubcriptionObjectStruct
	Sub.GetListByUserIdWithDB(userId)
	c.JSON(http.StatusOK, FormatResultAPI(http.StatusOK, "", "OK"))
}
