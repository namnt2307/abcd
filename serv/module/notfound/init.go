package notfound

import (
	"net/http"
	. "cm-v5/serv/module"
	"github.com/gin-gonic/gin"
	jsoniter "github.com/json-iterator/go"
)

var (
	mRedis RedisModelStruct
	mRedisKV RedisKVModelStruct
	NOTFOUND_PREVIEW_VOD_COL_ID, NOTFOUND_PREVIEW_LIVETV_COL_ID     string
	json = jsoniter.ConfigCompatibleWithStandardLibrary
)

func init() {
	NOTFOUND_PREVIEW_VOD_COL_ID, _ = CommonConfig.GetString("NOTFOUND_PREVIEW", "vod_collection_id")
	NOTFOUND_PREVIEW_LIVETV_COL_ID, _ = CommonConfig.GetString("NOTFOUND_PREVIEW", "livtv_collection_id")
}

func VodPreviewInit(c *gin.Context) {

	platform := c.DefaultQuery("platform", "web")
	cacheActive := StringToBool(c.DefaultQuery("cache", "true"))

	limit, err := StringToInt(c.DefaultQuery("limit", "30"))
	if err != nil {
		limit = 30
	}
	page, err := StringToInt(c.DefaultQuery("page", "0"))
	if err != nil {
		page = 0
	}

	data, err := GetDataNotFoundPreview(NOTFOUND_PREVIEW_VOD_COL_ID, platform, page, limit, cacheActive)
	if err != nil {
		Sentry_log(err)
		c.JSON(http.StatusBadRequest, FormatResultAPI(http.StatusBadRequest, err.Error(), ""))
		return
	}

	c.JSON(http.StatusOK, FormatResultAPI(http.StatusOK, "", data))

}

func LiveTVPreviewInit(c *gin.Context) {

	platform := c.DefaultQuery("platform", "web")
	cacheActive := StringToBool(c.DefaultQuery("cache", "true"))

	limit, err := StringToInt(c.DefaultQuery("limit", "30"))
	if err != nil {
		limit = 30
	}
	page, err := StringToInt(c.DefaultQuery("page", "0"))
	if err != nil {
		page = 0
	}

	data, err := GetDataNotFoundPreview(NOTFOUND_PREVIEW_LIVETV_COL_ID, platform, page, limit, cacheActive)
	if err != nil {
		Sentry_log(err)
		c.JSON(http.StatusBadRequest, FormatResultAPI(http.StatusBadRequest, err.Error(), ""))
		return
	}

	c.JSON(http.StatusOK, FormatResultAPI(http.StatusOK, "", data))

}