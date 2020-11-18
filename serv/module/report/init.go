package report

import (
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
func GetReportTypeInit(c *gin.Context) {
	platform := c.DefaultQuery("platform", "web")
	cacheActive := StringToBool(c.DefaultQuery("cache", "true"))
	limit, err := StringToInt(c.DefaultQuery("limit", "10"))
	if err != nil || limit > LIMIT_MAX {
		limit = 10
	}

	page, err := StringToInt(c.DefaultQuery("limit", "0"))
	if err != nil || page > PAGE_MAX {
		page = 0
	}

	var keyCache = LOCAL_REPORT_TYPE + "_" + platform + "_" + string(limit) + "_" + string(page)
	 
	if cacheActive {
		valC, err := LocalCache.GetValue(keyCache)
		if err == nil {

			c.JSON(http.StatusOK, FormatResultAPI(http.StatusOK, "", valC))
			return
		}
	}

	dataReportType, err := GetReportType(platform, page, limit, cacheActive)
	if err != nil {
		Sentry_log(err)
		c.JSON(http.StatusBadRequest, FormatResultAPI(http.StatusBadRequest, err.Error(), ""))
		return
	}

	LocalCache.SetValue(keyCache, dataReportType, TTL_LOCALCACHE)
	c.JSON(http.StatusOK, FormatResultAPI(http.StatusOK, "", dataReportType))
}

func ReportContentInit(c *gin.Context) {
	// platform := c.DefaultQuery("platform", "web")
	contentId := c.Param("content_id")
	userId := c.GetString("user_id")
	userAgent := c.GetHeader("User-Agent")
	token := c.GetHeader("Authorization")
	reportId := c.PostForm("report_id")
	message := c.PostForm("message")
	videoProfile := c.PostForm("video_profile")
	audioProfile := c.PostForm("audio_profile")
	subtitle := c.PostForm("subtitle")

	if reportId == "" {
		c.JSON(http.StatusBadRequest, FormatResultAPI(http.StatusBadRequest, "reportId required", ""))
		return
	}
	if contentId == "" {
		c.JSON(http.StatusBadRequest, FormatResultAPI(http.StatusBadRequest, "contentId required", ""))
		return
	}
	var ReportContent UserReportContentStruct
	ReportContent.Entity_id = contentId
	ReportContent.Report_id = reportId
	ReportContent.User_id = userId
	ReportContent.User_agent = userAgent
	ReportContent.Access_token = token
	ReportContent.Video_profile = videoProfile
	ReportContent.Audio_profile = audioProfile
	ReportContent.Subtitle = subtitle
	ReportContent.Message = message

	err := UserReportContent(ReportContent)
	if err != nil {
		c.JSON(http.StatusBadRequest, FormatResultAPI(http.StatusBadRequest, err.Error(), ""))
		return
	}
	c.JSON(http.StatusOK, FormatResultAPI(http.StatusOK, "success", "success"))

}
