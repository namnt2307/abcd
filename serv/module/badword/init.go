package badword

import (
	// "log"
	// "net/http"

	. "cm-v5/serv/module"
	// "github.com/gin-gonic/gin"
	jsoniter "github.com/json-iterator/go"
)

var mRedis RedisModelStruct
var mRedisKV RedisKVModelStruct
var json = jsoniter.ConfigCompatibleWithStandardLibrary

func init() {

}

// func PostBadwordInit(c *gin.Context) {
// 	userId := c.GetString("user_id")
// 	platform := c.DefaultQuery("platform", "web")
// 	contentId := c.Param("content_id")
// 	parentId := c.PostForm("parent_id")
// 	message := c.PostForm("message")

// 	message, err := CheckRuleMessage(message)
// 	if err != nil {
// 		log.Println("CheckRuleMessage fail userID "+userId+": ", err.Error())
// 		c.JSON(http.StatusBadRequest, FormatResultAPI(http.StatusBadRequest, err.Error(), ""))
// 		return
// 	}

// 	badwordObj, err := PostBadword(userId, platform, contentId, parentId, message, c)
// 	if err != nil {
// 		c.JSON(http.StatusBadRequest, FormatResultAPI(http.StatusBadRequest, err.Error(), ""))
// 		return
// 	}

// 	c.JSON(http.StatusOK, FormatResultAPI(1, "", badwordObj))
// }

// func GetBadwordsByContentIdInit(c *gin.Context) {
// 	contentID := c.Param("content_id")
// 	if contentID == "" {
// 		contentID = c.GetString("content_id")
// 	}
// 	page, err := StringToInt(c.DefaultQuery("page", "0"))
// 	if err != nil {
// 		page = 0
// 	}
// 	limit, err := StringToInt(c.DefaultQuery("limit", "10"))
// 	if err != nil {
// 		limit = 10
// 	}
// 	cacheActive := StringToBool(c.DefaultQuery("cache", "true"))
// 	if contentID == "" {
// 		c.JSON(http.StatusBadRequest, FormatResultAPI(http.StatusBadRequest, "Content not found", ""))
// 		return
// 	}

// 	dataBadword := GetBadwordsByContentId(page, limit, cacheActive)
// 	c.JSON(http.StatusOK, FormatResultAPI(http.StatusOK, "", dataBadword))

// }

// func DelBadwordByUser(c *gin.Context) {
// 	userId := c.GetString("user_id")
// 	contentId := c.Param("content_id")
// 	parentId := c.PostForm("parent_id")
// 	badwordId := c.PostForm("badword_id")

// 	if contentId == "" {
// 		c.JSON(http.StatusBadRequest, FormatResultAPI(http.StatusBadRequest, "contentId is required", ""))
// 		return
// 	}

// 	if badwordId == "" {
// 		c.JSON(http.StatusBadRequest, FormatResultAPI(http.StatusBadRequest, "BadwordId is required", ""))
// 		return
// 	}
// 	// Thu DT-9922 truyền thêm status khi xóa và duyệt 
// 	// err := DelBadword(userId, contentId, badwordId, parentId, "-1")
// 	// if err != nil {
// 	// 	c.JSON(http.StatusBadRequest, FormatResultAPI(http.StatusBadRequest, err.Error(), ""))
// 	// 	return
// 	// }
// 	c.JSON(http.StatusOK, FormatResultAPI(http.StatusOK, "", "success"))
// }
