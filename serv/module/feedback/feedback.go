package feedback

import (
	"net/http"

	. "cm-v5/serv/module"
	"github.com/gin-gonic/gin"
)

var mRedisKV RedisKVModelStruct
var KeyCacheFeedback = "FEEDBACK_USER_"

func Feedback_GetList(c *gin.Context) {
	userId := c.GetString("user_id")
	key := KeyCacheFeedback + userId
	dataCache := mRedisKV.Incr(key)
	if dataCache > 1 {
		var dataEmpty struct {}
		c.JSON(http.StatusOK, FormatResultAPI(http.StatusOK, "", dataEmpty))
		return
	}


	var FeedbackData struct {
		Like []string
		Dislike []string
	}

	FeedbackData.Like = append(FeedbackData.Like , "Có nhiều nội dung mà tôi thích")
	FeedbackData.Like = append(FeedbackData.Like , "Tải video rất nhanh")
	FeedbackData.Like = append(FeedbackData.Like , "Ứng dụng dễ sử dụng")

	FeedbackData.Dislike = append(FeedbackData.Dislike , "Không có nhiều nội dung")
	FeedbackData.Dislike = append(FeedbackData.Dislike , "Tải video rất chậm")
	FeedbackData.Dislike = append(FeedbackData.Dislike , "Ứng dụng khó sử dụng")

	c.JSON(http.StatusOK, FormatResultAPI(http.StatusOK, "", FeedbackData))
}

func Feedback_Open(c *gin.Context) {
	var screen = c.PostForm("screen") // type: QuestionScreen / LikeScreen / DislikeScreen
	if screen != "QuestionScreen" && screen != "LikeScreen" && screen != "DislikeScreen" {
		c.JSON(http.StatusBadRequest, FormatResultAPI(http.StatusBadRequest, "", "QuestionScreen/LikeScreen/DislikeScreen"))
		return
	}

	c.JSON(http.StatusOK, FormatResultAPI(http.StatusOK, "", "opened"))
}

func Feedback_Close(c *gin.Context) {
	var screen = c.PostForm("screen") // type: QuestionScreen / LikeScreen / DislikeScreen
	if screen != "QuestionScreen" && screen != "LikeScreen" && screen != "DislikeScreen" {
		c.JSON(http.StatusBadRequest, FormatResultAPI(http.StatusBadRequest, "", "QuestionScreen/LikeScreen/DislikeScreen"))
		return
	}

	userId := c.GetString("user_id")
	key := KeyCacheFeedback + userId

	switch screen {
	case "QuestionScreen":
		mRedisKV.Expire(key, 86400 * 7 * 2)
	case "LikeScreen":
		mRedisKV.Expire(key, 86400 * 30 * 2)
	case "DislikeScreen":
		mRedisKV.Expire(key, 86400 * 30 * 2)
	}

	c.JSON(http.StatusOK, FormatResultAPI(http.StatusOK, "", "closed"))
}


func Feedback_Clear(c *gin.Context) {
	userId := c.GetString("user_id")
	key := KeyCacheFeedback + userId
	mRedisKV.Del(key)
	
	c.JSON(http.StatusOK, FormatResultAPI(http.StatusOK, "", "ok"))
}


func Feedback_Save(c *gin.Context) {
	var data struct {
		Status string
		Message string
	}
	data.Status = "Success"
	data.Message = "Cám ơn ý kiến đóng góp của bạn!"

	var dataF = c.PostForm("data") 
	if dataF == "" {
		data.Status = "Error"
		data.Message = "Vui lòng chọn ít nhất một góp ý!"
		c.JSON(http.StatusBadRequest, FormatResultAPI(http.StatusBadRequest, "", data))
	}

	c.JSON(http.StatusOK, FormatResultAPI(http.StatusOK, "", data))
}