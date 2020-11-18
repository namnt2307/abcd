package cache

import (
	"net/http"

	. "cm-v5/serv/module"
	"github.com/gin-gonic/gin"
)

func UpdateCacheSetting(c *gin.Context) {
	c.JSON(http.StatusOK, FormatResultAPI(http.StatusOK, "", "Update cache done"))
}
