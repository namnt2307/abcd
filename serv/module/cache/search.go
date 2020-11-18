package cache

import (
	"net/http"

	. "cm-v5/serv/module"
	search "cm-v5/serv/module/search"

	// "fmt"
	"github.com/gin-gonic/gin"
)

func UpdateCacheSearchKeyword(c *gin.Context) {
	for _, val := range platforms {
		platform := Platform(val)
		search.GetSearchKeywordByMysql(platform.Type, 0, 10, false)
	}
	c.JSON(http.StatusOK, FormatResultAPI(http.StatusOK, "", "Update cache done"))
}
