package cache

import (
	"net/http"

	. "cm-v5/serv/module"
	. "cm-v5/schema"
	content "cm-v5/serv/module/content"
	"github.com/gin-gonic/gin"
)

func UpdateAdsForCP(c *gin.Context) {
	providerId := c.Param("provider_id")
	var arrAds []AdsOutputStruct
	for _, val := range platforms {
		content.GetAdsMySQL(providerId, arrAds, false, val)
	}

	c.JSON(http.StatusOK, FormatResultAPI(http.StatusOK, "", "Update cache done"))
}
