package cache

import (
	"net/http"

	. "cm-v5/serv/module"
	menu "cm-v5/serv/module/menu"
	"github.com/gin-gonic/gin"
)

func UpdateCacheMenu(c *gin.Context) {
	for _, val := range platforms {
		platform := Platform(val)
		menu.GetListMenuInfoMySQL(platform.Type, platform.Id, false)
	}
	c.JSON(http.StatusOK, FormatResultAPI(http.StatusOK, "", "Update cache done"))
}
