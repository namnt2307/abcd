package live_event

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

func GetListEvent(c *gin.Context) {
	page, err := StringToInt(c.DefaultQuery("page", "0"))
	if err != nil || page > PAGE_MAX {
		page = 0
	}
	limit, err := StringToInt(c.DefaultQuery("limit", "30"))
	if err != nil || limit > LIMIT_MAX {
		limit = 43
	}
	platform := c.DefaultQuery("platform", "web")
	cacheActive := StringToBool(c.DefaultQuery("cache", "true"))
	// userId := c.GetString("user_id")
	// ipUser, _ := GetClientIPHelper(c.Request, c)

	var keyCache = LOCAL_LIVE_EVENT + "_" + fmt.Sprint(page) + "_" + fmt.Sprint(limit) + "_" + platform
	if cacheActive {
		valC, err := LocalCache.GetValue(keyCache)
		if err == nil {
			var dataLiveEvents []LiveEventOutputObjectStruct
			dataByte, _ := json.Marshal(valC)
			json.Unmarshal([]byte(dataByte), &dataLiveEvents)
			// dataLiveEventsPer := GetPermissionListLiveEvent(dataLiveEvents, userId, ipUser)
			c.JSON(http.StatusOK, FormatResultAPI(http.StatusOK, "", dataLiveEvents))
			return
		}
	}

	dataLiveEvents, err := GetAllLiveEventByCache(platform, page, limit, cacheActive)
	if err != nil {
		c.JSON(http.StatusBadRequest, FormatResultAPI(http.StatusBadRequest, err.Error(), "Content not found"))
		return
	}
	// Write local data
	LocalCache.SetValue(keyCache, dataLiveEvents, TTL_LOCALCACHE)
	// dataLiveEventsPer := GetPermissionListLiveEvent(dataLiveEvents, userId, ipUser)
	c.JSON(http.StatusOK, FormatResultAPI(http.StatusOK, "", dataLiveEvents))
}

func GetListEventByID(c *gin.Context) {
	eventID := c.Param("event_id")
	if eventID == "" {
		eventID = c.GetString("event_id")
	}
	platform := c.DefaultQuery("platform", "web")
	cacheActive := StringToBool(c.DefaultQuery("cache", "true"))
	userId := c.GetString("user_id")
	ipUser, _ := GetClientIPHelper(c.Request, c)
	statusUserIsPremium := c.GetInt("user_is_premium")

	var keyCache = LOCAL_LIVE_EVENT_ID + "_" + eventID + "_" + platform
	if cacheActive {
		valC, err := LocalCache.GetValue(keyCache)
		if err == nil {
			var dataLiveEvent LiveEventOutputObjectStruct
			dataByte, _ := json.Marshal(valC)
			json.Unmarshal([]byte(dataByte), &dataLiveEvent)
			GetPermissionLiveEvent(&dataLiveEvent, userId, ipUser, statusUserIsPremium)
			c.JSON(http.StatusOK, FormatResultAPI(http.StatusOK, "", dataLiveEvent))
			return
		}
	}

	dataLiveEvent, err := GetDetailByID(eventID, platform, cacheActive)
	if err != nil {
		c.JSON(http.StatusBadRequest, FormatResultAPI(http.StatusBadRequest, err.Error(), "Content not found"))
		return
	}

	// Write local data
	LocalCache.SetValue(keyCache, dataLiveEvent, TTL_LOCALCACHE)
	GetPermissionLiveEvent(&dataLiveEvent, userId, ipUser, statusUserIsPremium)
	c.JSON(http.StatusOK, FormatResultAPI(http.StatusOK, "", dataLiveEvent))
}

func GetListEventBySlug(c *gin.Context) {
	liveEventSlug := c.PostForm("live_event_slug")
	cacheActive := StringToBool(c.DefaultQuery("cache", "true"))
	liveEventID, err := GetIdBySLug(liveEventSlug, cacheActive)
	if err != nil {
		c.JSON(http.StatusBadRequest, FormatResultAPI(http.StatusBadRequest, err.Error(), "Content not found"))
		return
	}

	c.Set("event_id", liveEventID)
	GetListEventByID(c)
}

// func GetListEventFinished(c *gin.Context) {
// 	page, err := StringToInt(c.DefaultQuery("page", "0"))
// 	if err != nil || page > PAGE_MAX {
// 		page = 0
// 	}
// 	limit, err := StringToInt(c.DefaultQuery("limit", "30"))
// 	if err != nil || limit > LIMIT_MAX {
// 		limit = 43
// 	}
// 	platform := c.DefaultQuery("platform", "web")
// 	cacheActive := StringToBool(c.DefaultQuery("cache", "true"))
// 	// userId := c.GetString("user_id")
// 	// ipUser := c.ClientIP()

// 	var keyCache = LOCAL_LIVE_EVENT_FINISHED + "_" + fmt.Sprint(page) + "_" + fmt.Sprint(limit) + "_" + platform
// 	if cacheActive {
// 		valC, err := LocalCache.GetValue(keyCache)
// 		if err == nil {
// 			var dataLiveEvents LiveEventFinishedOutputObjectStruct
// 			dataByte, _ := json.Marshal(valC)
// 			json.Unmarshal([]byte(dataByte), &dataLiveEvents)
// 			// dataLiveEvents.Items = GetPermissionListLiveEvent(dataLiveEvents.Items, userId, ipUser)
// 			c.JSON(http.StatusOK, FormatResultAPI(http.StatusOK, "", dataLiveEvents))
// 			return
// 		}
// 	}

// 	dataLiveEvents, err := GetLiveEventFinishedByCache(platform, page, limit, cacheActive)
// 	if err != nil {
// 		c.JSON(http.StatusBadRequest, FormatResultAPI(http.StatusBadRequest, err.Error(), "Content not found"))
// 		return
// 	}
// 	// Write local data
// 	LocalCache.SetValue(keyCache, dataLiveEvents, TTL_LOCALCACHE)
// 	// dataLiveEvents.Items = GetPermissionListLiveEvent(dataLiveEvents.Items, userId, ipUser)
// 	c.JSON(http.StatusOK, FormatResultAPI(http.StatusOK, "", dataLiveEvents))
// }
