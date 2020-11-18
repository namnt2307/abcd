package cache

import (
	"net/http"
	"strconv"

	. "cm-v5/serv/module"
	content "cm-v5/serv/module/content"
	vod "cm-v5/serv/module/vod"

	. "cm-v5/schema"
	"github.com/gin-gonic/gin"
)

var platforms []string = []string{"app", "smarttv", "web", "ios", "android", "samsung_tv", "sony_androidtv", "lg_tv", "androidtv", "mobile_web"}

//Detail
func UpdateCacheDetail(c *gin.Context) {
	contentId := c.Query("content_id")
	epsId := c.Query("eps_id")
	eps := c.Query("eps")

	if contentId != "" {
		UpdateCacheContent(contentId)
		UpdateCacheRelated(contentId)
	}

	if contentId != "" && epsId != "" && eps != "" {
		UpdateCacheEpisode(contentId, eps, epsId)
	}

	c.JSON(http.StatusOK, FormatResultAPI(http.StatusOK, "", "Update done"))
}

func UpdateCacheContent(contentId string) {
	// fmt.Println("UpdateCacheContent")
	// fmt.Println("----------------------------------------------------------")
	for _, val := range platforms {
		_, err := content.GetContent(contentId, val, false)
		if err != nil {
			// fmt.Println("Clear cache content not fail")
			return
		}

		//Khanh DT-11009
		content.GetInfoVod(contentId, "", val, false)
		// fmt.Println("UpdateCacheContent with ID : " + contentId + " platform: " + val)
	}

	// fmt.Println("----------------------------------------------------------")

}

func UpdateCacheEpisode(groupId string, eps string, epsId string) {
	// fmt.Println("UpdateCacheEpisode")
	// fmt.Println("----------------------------------------------------------")
	score, _ := strconv.ParseFloat(eps, 64)
	keyCache := CONTENT_LIST_EPS + "_" + groupId

	err := vod.WriteCacheZRangeByListGroup(keyCache, score, epsId)
	if err != nil {
		// fmt.Println("Clear cache episode not fail")
		return
	}

	for _, val := range platforms {
		_, err := content.GetTotalEpisodeVOD(groupId, val, false)
		if err != nil {
			// fmt.Println("Clear cache episode total not fail")
			return
		}
		// fmt.Println("Update cache episode total with ID : " + groupId + " platform: " + val)
	}

	// fmt.Println("----------------------------------------------------------")
}

func UpdateCacheRelated(contentId string) {
	// fmt.Println("UpdateCacheRelated")
	// fmt.Println("----------------------------------------------------------")

	for _, val := range platforms {
		_, err := content.GetRelated(contentId, 0, 30, val, false)
		if err != nil {
			// fmt.Println("Clear cache related not fail")
			return
		}
		// fmt.Println("UpdateCacheRelated with ID : " + contentId + " platform: " + val)
	}

	// fmt.Println("----------------------------------------------------------")

}
