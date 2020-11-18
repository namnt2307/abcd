package comment

import (
	"log"
	"net/http"

	. "cm-v5/serv/module"
	. "cm-v5/schema"
	"github.com/gin-gonic/gin"
	jsoniter "github.com/json-iterator/go"
)

var mRedis RedisModelStruct
var mRedisKV RedisKVModelStruct
var json = jsoniter.ConfigCompatibleWithStandardLibrary

func init() {

}

func PostCommentInit(c *gin.Context) {
	userId := c.GetString("user_id")
	platform := c.DefaultQuery("platform", "web")
	contentId := c.Param("content_id")
	parentId := c.PostForm("parent_id")
	message := c.PostForm("message")

	message, err := CheckRuleMessage(message)
	if err != nil {
		log.Println("CheckRuleMessage fail userID "+userId+": ", err.Error())
		c.JSON(http.StatusBadRequest, FormatResultAPI(http.StatusBadRequest, err.Error(), ""))
		return
	}

	// Thu add DT-9922 filter blackwords comment
	old_message := message
	message = FilterMessage(message, userId)
	message = FilterMessageLike(message)
	log.Println("Data sau khi loc ==== " + message)
	commentObj, err := PostComment(userId, platform, contentId, parentId, message, old_message, c)
	if err != nil {
		c.JSON(http.StatusBadRequest, FormatResultAPI(http.StatusBadRequest, err.Error(), ""))
		return
	}

	c.JSON(http.StatusOK, FormatResultAPI(1, "", commentObj))
}

func GetCommentsByContentIdInit(c *gin.Context) {

	contentID := c.Param("content_id")
	if contentID == "" {
		contentID = c.GetString("content_id")
	}
	page, err := StringToInt(c.DefaultQuery("page", "0"))
	if err != nil || page > PAGE_MAX || page < 0 {
		page = 0
	}
	limit, err := StringToInt(c.DefaultQuery("limit", "10"))
	if err != nil || limit > LIMIT_MAX || limit < 0 {
		limit = 10
	}
	cacheActive := StringToBool(c.DefaultQuery("cache", "true"))
	if contentID == "" {
		c.JSON(http.StatusBadRequest, FormatResultAPI(http.StatusBadRequest, "Content not found", ""))
		return
	}
	dataComment := GetCommentsByContentId(contentID, page, limit, cacheActive)
	// DT-9922
	userId := c.GetString("user_id")
	dataComment = FilterByUserID(dataComment, page, limit, userId)
	c.JSON(http.StatusOK, FormatResultAPI(http.StatusOK, "", dataComment))

}

func DelCommentByUser(c *gin.Context) {
	userId := c.GetString("user_id")
	contentId := c.Param("content_id")
	parentId := c.PostForm("parent_id")
	commentId := c.PostForm("comment_id")

	if contentId == "" {
		c.JSON(http.StatusBadRequest, FormatResultAPI(http.StatusBadRequest, "contentId is required", ""))
		return
	}

	if commentId == "" {
		c.JSON(http.StatusBadRequest, FormatResultAPI(http.StatusBadRequest, "commentId is required", ""))
		return
	}
	// Thu DT-9922 truyền thêm status khi xóa và duyệt
	err := DelComment(userId, contentId, commentId, parentId, -1)
	if err != nil {
		c.JSON(http.StatusBadRequest, FormatResultAPI(http.StatusBadRequest, err.Error(), ""))
		return
	}
	c.JSON(http.StatusOK, FormatResultAPI(http.StatusOK, "", "success"))
}
