package redirect

import (
	. "cm-v5/serv/module"
	"net/http"

	"github.com/gin-gonic/gin"
	jsoniter "github.com/json-iterator/go"
)

var mRedis RedisModelStruct
var json = jsoniter.ConfigCompatibleWithStandardLibrary

func RedirectInit(c *gin.Context) {
	from_url := c.PostForm("from_url")
	output, err := CheckURLDirect(from_url)
	if err != nil {
		c.JSON(http.StatusBadRequest, FormatResultAPI(http.StatusBadRequest, err.Error(), ""))
		return
	}
	c.JSON(http.StatusOK, FormatResultAPI(http.StatusOK, "", output))
}
