package cache

import (
	"net/http"

	. "cm-v5/serv/module"
	artist "cm-v5/serv/module/artist"
	"github.com/gin-gonic/gin"
)

func UpdateCacheContentsByID(c *gin.Context) {
	peopleID := c.Param("people_id")

	for _, val := range platforms {
		artist.GetVodArtist(peopleID, val, 0, 30, 0, false)
	}
	c.JSON(http.StatusOK, FormatResultAPI(http.StatusOK, "", "Update cache done"))
}

func UpdateCacheArtistById(c *gin.Context) {
	peopleID := c.Param("people_id")
	artist.GetInfoArtist(peopleID, false)

	c.JSON(http.StatusOK, FormatResultAPI(http.StatusOK, "", "Update cache done"))
}

func UpdateCacheArtistRelatedById(c *gin.Context) {
	peopleID := c.Param("people_id")
	for _, val := range platforms {
		artist.GetArtistRelated(peopleID, val, 0, 30, false)
	}
	c.JSON(http.StatusOK, FormatResultAPI(http.StatusOK, "", "Update cache done"))
}
