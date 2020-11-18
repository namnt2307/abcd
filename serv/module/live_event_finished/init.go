package live_event_finished

import (
	"fmt"
	"net/http"

	. "cm-v5/serv/module"
	. "cm-v5/schema"
	"github.com/gin-gonic/gin"
	jsoniter "github.com/json-iterator/go"
)

var json = jsoniter.ConfigCompatibleWithStandardLibrary
var mRedis RedisModelStruct
var mRedisKV RedisKVModelStruct
var mLocal LocalModelStruct

func GetListEventFinished(c *gin.Context) {
	page, err := StringToInt(c.DefaultQuery("page", "0"))
	if err != nil {
		page = 0
	}
	limit, err := StringToInt(c.DefaultQuery("limit", "30"))
	if err != nil {
		limit = 43
	}
	platform := c.DefaultQuery("platform", "web")
	cacheActive := StringToBool(c.DefaultQuery("cache", "true"))
	// userId := c.GetString("user_id")
	// ipUser := c.ClientIP()

	var keyCache = LOCAL_LIVE_EVENT_FINISHED + "_" + fmt.Sprint(page) + "_" + fmt.Sprint(limit) + "_" + platform
	if cacheActive {
		valC, err := mLocal.GetValue(keyCache)
		if err == nil {
			var dataLiveEvents LiveEventFinishedOutputObjectStruct
			dataByte, _ := json.Marshal(valC)
			json.Unmarshal([]byte(dataByte), &dataLiveEvents)
			c.JSON(http.StatusOK, FormatResultAPI(http.StatusOK, "", dataLiveEvents))
			return
		}
	}

	// dataLiveEvents, err := GetLiveEventFinishedByMongoDB(platform, page, limit, cacheActive)
	var dataLiveEvents LiveEventFinishedOutputObjectStruct
	if err != nil {
		c.JSON(http.StatusBadRequest, FormatResultAPI(http.StatusBadRequest, err.Error(), "Content not found"))
		return
	}
	// Write local data
	mLocal.SetValue(keyCache, dataLiveEvents, TTL_LOCALCACHE)
	// dataLiveEvents.Items = GetPermissionListLiveEvent(dataLiveEvents.Items, userId, ipUser)
	c.JSON(http.StatusOK, FormatResultAPI(http.StatusOK, "", dataLiveEvents))
}
