package campaign

import (
	"net/http"

	. "cm-v5/serv/module"
	"github.com/gin-gonic/gin"
)

// var mRedisKV RedisKVModelStruct
// var KeyCacheFeedback = "FEEDBACK_USER_"


// Check Event Is valid
func Samsung_3008_3011(c *gin.Context) {
	// userId := c.GetString("user_id")


	c.JSON(http.StatusOK, FormatResultAPI(http.StatusOK, "", "OK"))
}
