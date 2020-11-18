package comment

import (
	"net/http"

	"log"

	"time"

	. "cm-v5/serv/module"
	. "cm-v5/schema"
	"github.com/gin-gonic/gin"
)

func GetCommentsInit(c *gin.Context) {
	contentId := c.DefaultQuery("content_id", "")
	parentId := c.DefaultQuery("parent_id", "")
	userId := c.DefaultQuery("user_id", "")
	status := c.DefaultQuery("status", "")
	sort := c.DefaultQuery("sort", "asc")
	textSearch := c.DefaultQuery("text_search", "")
	dateFilter := c.DefaultQuery("date_filter", "")

	page, err := StringToInt(c.DefaultQuery("page", "0"))
	if err != nil {
		page = 0
	}
	limit, err := StringToInt(c.DefaultQuery("limit", "30"))
	if err != nil {
		limit = 30
	}
	var dataPage CommentPageStruct
	dataPage.Items = GetComments(contentId, parentId, userId, dateFilter, status, textSearch, sort, page, limit)
	dataPage.Metadata.Page = page
	dataPage.Metadata.Limit = limit
	dataPage.Metadata.Total = GetTotalComment(contentId, parentId, userId, dateFilter, status, textSearch)
	c.JSON(http.StatusOK, FormatResultAPI(http.StatusOK, "", dataPage))

}

func DelCommentByAdminInit(c *gin.Context) {

	userId := c.PostForm("user_id")
	contentId := c.PostForm("content_id")
	parentId := c.PostForm("parent_id")
	commentId := c.PostForm("comment_id")
	status, _ := StringToInt(c.PostForm("status"))

	if contentId == "" {
		c.JSON(http.StatusBadRequest, FormatResultAPI(http.StatusBadRequest, "contentId is required", ""))
		return
	}
	if userId == "" {
		c.JSON(http.StatusBadRequest, FormatResultAPI(http.StatusBadRequest, "userId is required", ""))
		return
	}

	if commentId == "" {
		c.JSON(http.StatusBadRequest, FormatResultAPI(http.StatusBadRequest, "commentId is required", ""))
		return
	}

	err := DelComment(userId, contentId, commentId, parentId, status)
	if err != nil {

		c.JSON(http.StatusBadRequest, FormatResultAPI(http.StatusBadRequest, err.Error(), ""))
		return
	}
	c.JSON(http.StatusOK, FormatResultAPI(http.StatusOK, "", "success"))
}

func PinCommentByAdminInit(c *gin.Context) {

	userId := c.PostForm("user_id")
	contentId := c.PostForm("content_id")
	commentId := c.PostForm("comment_id")
	status, _ := StringToInt(c.PostForm("status"))

	if contentId == "" {
		c.JSON(http.StatusBadRequest, FormatResultAPI(http.StatusBadRequest, "contentId is required", ""))
		return
	}
	if userId == "" {
		c.JSON(http.StatusBadRequest, FormatResultAPI(http.StatusBadRequest, "userId is required", ""))
		return
	}

	if commentId == "" {
		c.JSON(http.StatusBadRequest, FormatResultAPI(http.StatusBadRequest, "commentId is required", ""))
		return
	}

	err := PinComment(contentId, commentId, status)
	if err != nil {
		c.JSON(http.StatusBadRequest, FormatResultAPI(http.StatusBadRequest, err.Error(), ""))
		return
	}
	log.Println(err)
	c.JSON(http.StatusOK, FormatResultAPI(http.StatusOK, "", "success"))
}

func AddCommentByAdminInit(c *gin.Context) {

	contentId := c.PostForm("content_id")
	message := c.PostForm("msg")

	if contentId == "" {
		c.JSON(http.StatusBadRequest, FormatResultAPI(http.StatusBadRequest, "contentId is required", ""))
		return
	}

	if message == "" {
		c.JSON(http.StatusBadRequest, FormatResultAPI(http.StatusBadRequest, "msg is required", ""))
		return
	}

	t := time.Now()
	var commentObj CommentObjectStruct
	commentObj.Id = UUIDV4Generate()
	commentObj.Content_id = contentId
	commentObj.Message = message
	commentObj.Status = 1
	commentObj.Created_at = t.Unix()
	commentObj.Updated_at = t.Unix()
	commentObj.User.Id = "58a5d91e-cfc7-ff49-2943-3ca9f705bcf7"
	commentObj.User.Avatar = "https://static.vieon.vn/vieplay-image/artis_avatar/2020/06/15/qd76tzsq_logo-ad746d92e400f7b9ba07bd9456d5f48b8b.png"
	commentObj.User.Name = "Admin VieON"
	commentObj.User.Gender = 0
	platforms := Platform("web")
	commentObj.Platforms = platforms

	err := AddComment(commentObj)
	if err != nil {
		c.JSON(http.StatusBadRequest, FormatResultAPI(http.StatusBadRequest, err.Error(), ""))
		return
	}
	c.JSON(http.StatusOK, FormatResultAPI(http.StatusOK, "", "success"))
}
