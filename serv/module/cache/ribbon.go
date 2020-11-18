package cache

import (
	// . "cm-v5/schema"
	. "cm-v5/serv/module"
	page "cm-v5/serv/module/page"
	"net/http"

	"github.com/gin-gonic/gin"
)

func ClearCacheRibbonByMenuID(c *gin.Context) {
	var menu_id string = c.Param("menu_id")
	if menu_id == "" {
		menu_id = c.GetString("menu_id")
	}

	if menu_id == "" {
		c.JSON(http.StatusBadRequest, FormatResultAPI(http.StatusBadRequest, "", "menu_id is required"))
		return
	}

	dataRibbon, err := page.GetRibbonByMenuIDNoCache(menu_id)

	if err == nil {
		ClearCachePageRibbonV3(menu_id, dataRibbon.Type)
	}

	c.JSON(http.StatusOK, FormatResultAPI(http.StatusOK, "", "success"))
}

func ClearCacheRibbonByRibbonId(c *gin.Context) {
	var ribbon_id string = c.Param("ribbon_id")
	if ribbon_id == "" {
		ribbon_id = c.GetString("ribbon_id")
	}

	if ribbon_id == "" {
		c.JSON(http.StatusBadRequest, FormatResultAPI(http.StatusBadRequest, "", "ribbon_id is required"))
		return
	}

	ClearCacheRibbonV3(ribbon_id)

	c.JSON(http.StatusOK, FormatResultAPI(http.StatusOK, "", "success"))
}

func ClearCacheRibbonByRibbonItemId(c *gin.Context) {
	var ribbon_item_id string = c.Param("ribbon_item_id")
	if ribbon_item_id == "" {
		ribbon_item_id = c.GetString("ribbon_item_id")
	}

	if ribbon_item_id == "" {
		c.JSON(http.StatusBadRequest, FormatResultAPI(http.StatusBadRequest, "", "ribbon_item_id is required"))
		return
	}

	ribItemObj, err := page.GetRibbonItemByIdNoCache(ribbon_item_id)
	if err == nil {
		ClearCacheRibbonItemV3(ribbon_item_id, ribItemObj.Rib_ids, ribItemObj.Status)
	}

	c.JSON(http.StatusOK, FormatResultAPI(http.StatusOK, "", "success"))
}
