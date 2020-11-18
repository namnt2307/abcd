package question

import (
	"encoding/hex"
	"net/http"

	. "cm-v5/serv/module"
	"github.com/gin-gonic/gin"
)

func ExportRankingChannel(c *gin.Context) {

	channel_id := c.DefaultQuery("channel_id", "")
	if channel_id == "" {
		c.JSON(http.StatusBadRequest, FormatResultAPI(http.StatusBadRequest, "channel_id required", "Fail"))
		return
	}

	resultData, err := GetRankingByChannel(channel_id)

	if err != nil {
		c.JSON(http.StatusBadRequest, FormatResultAPI(http.StatusBadRequest, err.Error(), "Fail"))
		return
	}
	if resultData.Channel_id == "" {
		c.JSON(http.StatusBadRequest, FormatResultAPI(http.StatusBadRequest, "Channel không hợp lệ hoặc chưa tồn tại BXH", "Fail"))
		return
	}

	c.JSON(http.StatusOK, FormatResultAPI(http.StatusOK, "", resultData))
}
func StatisticQuestion(c *gin.Context) {

	question_id := c.DefaultQuery("question_id", "")
	if question_id == "" {
		c.JSON(http.StatusBadRequest, FormatResultAPI(http.StatusBadRequest, "question Id required", "Fail"))
		return
	}
	d, _ := hex.DecodeString(question_id)
	if len(d) != 12 {
		c.JSON(http.StatusBadRequest, FormatResultAPI(http.StatusBadRequest, "question Id not correct", "Fail"))
		return
	}

	resultData, err := QuestionStatistics(question_id)

	if err != nil {
		c.JSON(http.StatusBadRequest, FormatResultAPI(http.StatusBadRequest, err.Error(), "Fail"))
		return
	}

	c.JSON(http.StatusOK, FormatResultAPI(http.StatusOK, "", resultData))
}
