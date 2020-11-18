package cache

import (
	"net/http"

	. "cm-v5/serv/module"
	config "cm-v5/serv/module/config"
	"github.com/gin-gonic/gin"
)

func UpdateCacheConfigByKey(c *gin.Context) {
	key := c.DefaultQuery("key", "")

	for _, val := range platforms {
		config.GetConfigByKey(key, val, false)
	}
	c.JSON(http.StatusOK, FormatResultAPI(http.StatusOK, "", "Update cache done"))
}
