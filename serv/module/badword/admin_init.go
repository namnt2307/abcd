package badword

import (
	"net/http"
	// "log"
	. "cm-v5/serv/module"
	. "cm-v5/schema"
	"github.com/gin-gonic/gin"
)

func GetBadwordsInit(c *gin.Context) {
	status := c.DefaultQuery("status", "")
	sort := c.DefaultQuery("sort", "asc")
	textSearch := c.DefaultQuery("text_search", "")

	page, err := StringToInt(c.DefaultQuery("page", "0"))
	if err != nil {
		page = 0
	}
	limit, err := StringToInt(c.DefaultQuery("limit", "30"))
	if err != nil {
		limit = 30
	}
	var dataPage BadwordPageStruct
	dataPage.Items = GetBadwords(status, textSearch, sort, page, limit)
	dataPage.Metadata.Page = page
	dataPage.Metadata.Limit = limit
	dataPage.Metadata.Total = GetTotalBadword(status, textSearch)
	c.JSON(http.StatusOK, FormatResultAPI(http.StatusOK, "", dataPage))

}


func UpdateWordInit(c *gin.Context) {

	id := c.PostForm("id")
	content := c.PostForm("content")
 
	if id == "" {
		c.JSON(http.StatusBadRequest, FormatResultAPI(http.StatusBadRequest, "ID is required", ""))
		return
	}
	if content == "" {
		c.JSON(http.StatusBadRequest, FormatResultAPI(http.StatusBadRequest, "content is required", ""))
		return
	}

	err := UpdateWord(id, content)
	if err != nil {
		c.JSON(http.StatusBadRequest, FormatResultAPI(http.StatusBadRequest, err.Error(), ""))
		return
	}
	c.JSON(http.StatusOK, FormatResultAPI(http.StatusOK, "", "success"))
}

func AddWordInit(c *gin.Context) {

	content := c.PostForm("content")
 
	if content == "" {
		c.JSON(http.StatusBadRequest, FormatResultAPI(http.StatusBadRequest, "content is required", ""))
		return
	}

	err := AddWord(content)
	if err != nil {
		c.JSON(http.StatusBadRequest, FormatResultAPI(http.StatusBadRequest, err.Error(), ""))
		return
	}
	c.JSON(http.StatusOK, FormatResultAPI(http.StatusOK, "", "success"))
}


func DelBadwordByAdminInit(c *gin.Context) {

	id := c.PostForm("id")
	status := c.PostForm("status")
 
	if id == "" {
		c.JSON(http.StatusBadRequest, FormatResultAPI(http.StatusBadRequest, "ID is required", ""))
		return
	}
	if status == "" {
		c.JSON(http.StatusBadRequest, FormatResultAPI(http.StatusBadRequest, "status is required", ""))
		return
	}

	err := DelBadword(id, status)
	if err != nil {
		c.JSON(http.StatusBadRequest, FormatResultAPI(http.StatusBadRequest, err.Error(), ""))
		return
	}
	c.JSON(http.StatusOK, FormatResultAPI(http.StatusOK, "", "success"))
}
