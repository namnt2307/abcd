package cache

import (
	"net/http"

	. "cm-v5/serv/module"
	tag "cm-v5/serv/module/tag"
	"github.com/gin-gonic/gin"
)

func UpdateCacheTagListID(c *gin.Context) {
	tags := c.PostForm("tags")
	//Parse string to array
	var arrTagIds []string
	arrTagIds, _ = ParseStringToArray(tags)

	for _, val := range platforms {
		tag.GetVODByTagIds(arrTagIds, 1, 0, 30, val, false)
		tag.GetVODByTagIds(arrTagIds, 2, 0, 30, val, false)
		tag.GetVODByTagIds(arrTagIds, 3, 0, 30, val, false)
	}
	c.JSON(http.StatusOK, FormatResultAPI(http.StatusOK, "", "Update cache done"))
}
