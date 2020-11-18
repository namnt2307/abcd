package cache

import (
	"net/http"

	. "cm-v5/serv/module"
	message "cm-v5/serv/module/message"
	"github.com/gin-gonic/gin"
)

func UpdateListMessage(c *gin.Context) {
	ClearListMessage()
	c.JSON(http.StatusOK, FormatResultAPI(http.StatusOK, "", "Update cache done"))
}

func ClearListMessage() {
	for _, val := range platforms {
		message.GetMessageListAll("", val, false)
	}
	// fmt.Println("Update cache done")
}
