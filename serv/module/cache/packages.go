package cache

import (
	"net/http"

	. "cm-v5/serv/module"
	packages_module "cm-v5/serv/module/packages"

	// . "cm-v5/schema"
	"github.com/gin-gonic/gin"
)

func UpdateCachePackagesByContent(c *gin.Context) {
	contentId := c.Param("content_id")

	if contentId != "" {
		var Pack packages_module.PackagesObjectStruct
		Pack.GetListPackageByDenyContent(contentId, false)
	}

	c.JSON(http.StatusOK, FormatResultAPI(http.StatusOK, "", "Update cache done"))
}
