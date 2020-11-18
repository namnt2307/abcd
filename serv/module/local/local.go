package local

import (
	"net/http"
	. "cm-v5/serv/module"
	"github.com/gin-gonic/gin"
)

func GetLocalInfo(c *gin.Context) {
	//  
	// val := LocalCache.GetLocalInfo()
	c.JSON(http.StatusOK, FormatResultAPI(http.StatusOK, "", ""))
}
