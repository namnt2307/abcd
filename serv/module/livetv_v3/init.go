package livetv_v3

import (
	"fmt"
	"net/http"
	"strings"

	. "cm-v5/serv/module"

	// . "ott-backend-go/serv/module_v4"
	. "cm-v5/schema"
	Packages_premium "cm-v5/serv/module/packages_premium"
	seo "cm-v5/serv/module/seo"
	Subscription "cm-v5/serv/module/subscription"

	"github.com/gin-gonic/gin"
	jsoniter "github.com/json-iterator/go"
)

var mRedisUSC RedisUSCModelStruct
var mRedis RedisModelStruct

var mRedisKV RedisKVModelStruct
var json = jsoniter.ConfigCompatibleWithStandardLibrary

func GetLiveTVCategory(c *gin.Context) {
	cache := StringToBool(c.DefaultQuery("cache", "true"))
	platform := c.DefaultQuery("platform", "android")
	userId := c.GetString("user_id")

	keyCache := LOCAL_LIVE_TV_GROUP
	if cache {
		valC, err := LocalCache.GetValue(keyCache)
		if err == nil {
			//Khanh DT-10720
			listLiveTVGroup := make([]LiveTVGroupOutputObject, 0)
			dataByte, _ := json.Marshal(valC)
			json.Unmarshal(dataByte, &listLiveTVGroup)

			//check group super premium
			listLiveTVGroup = CheckGroupSuperPremium(listLiveTVGroup, userId)
			c.JSON(http.StatusOK, FormatResultAPI(http.StatusOK, "", listLiveTVGroup))
			return
		}
	}

	dataLiveTVGroup, err := GetLiveTVGroup(platform, cache)
	if err != nil {
		c.JSON(http.StatusBadRequest, FormatResultAPI(http.StatusBadRequest, err.Error(), ""))
		return
	}

	// Write local data
	LocalCache.SetValue(keyCache, dataLiveTVGroup, TTL_LOCALCACHE)

	//Khanh DT-10720
	//check group super premium
	dataLiveTVGroup = CheckGroupSuperPremium(dataLiveTVGroup, userId)

	c.JSON(http.StatusOK, FormatResultAPI(http.StatusOK, "", dataLiveTVGroup))
}

func GetDetailLiveTVByID(c *gin.Context) {
	livetv_id := c.Param("livetv_id")
	epg_slug := c.Param("epg_slug")
	if livetv_id == "" {
		livetv_id = c.GetString("livetv_id")
	}
	if epg_slug == "" {
		epg_slug = c.GetString("epg_slug")
	}
	platform := c.DefaultQuery("platform", "android")
	cache := StringToBool(c.DefaultQuery("cache", "true"))
	user_id := c.GetString("user_id")
	statusUserIsPremium := c.GetInt("user_is_premium")
	ipUser, _ := GetClientIPHelper(c.Request, c)
	tokenUser := c.GetHeader("Authorization")

	// VBE-17
	if user_id != "" {
		keyCacheLock := "lock_detail_" + user_id
		valIncr := mRedis.Incr(keyCacheLock)
		if valIncr > int64(LIVETV_LIMIT_REQUEST) {
			c.JSON(http.StatusLocked, FormatResultAPI(http.StatusLocked, "", "Vượt quá giới hạn. Vui lòng thử lại sau 60s"))
			return
		}
		mRedis.Expire(keyCacheLock, 60)
	}
	// VBE-17

	keyCache := LOCAL_DETAIL_LIVE_TV + "_" + livetv_id + "_" + epg_slug + "_" + platform
	if cache {
		valC, err := LocalCache.GetValue(keyCache)
		if err != nil && valC != "" {
			var DetailLiveTVObject DetailLiveTVObjectOutputStruct
			dataByte, _ := json.Marshal(valC)
			json.Unmarshal([]byte(dataByte), &DetailLiveTVObject)
			RouterGetPermission(&DetailLiveTVObject, user_id, ipUser, tokenUser, statusUserIsPremium)
			c.JSON(http.StatusOK, FormatResultAPI(http.StatusOK, "", DetailLiveTVObject))
			return
		}
	}

	dataDetailLiveTV, err := GetDetailLiveTV(livetv_id, epg_slug, platform, cache)
	if err != nil {
		c.JSON(http.StatusBadRequest, FormatResultAPI(http.StatusBadRequest, err.Error(), ""))
		return
	}

	// Write local data
	LocalCache.SetValue(keyCache, dataDetailLiveTV, TTL_LOCALCACHE)
	RouterGetPermission(&dataDetailLiveTV, user_id, ipUser, tokenUser, statusUserIsPremium)
	c.JSON(http.StatusOK, FormatResultAPI(http.StatusOK, "", dataDetailLiveTV))
}

func GetDetailLiveTVBySlug(c *gin.Context) {
	var livetvSlug = c.PostForm("livetv_slug")
	var epgSlug string

	// Check slug exists
	livetv_id := seo.CheckExistsSEOBySlug(livetvSlug)
	if livetv_id != "" {
		c.Set("livetv_id", livetv_id)
		GetDetailLiveTVByID(c)
		return
	}

	// Check slug validate
	slugSplits := strings.Split(livetvSlug, fmt.Sprintf(SEO_LIVETV_URL, ""))
	if len(slugSplits) != 2 {
		c.JSON(http.StatusBadRequest, FormatResultAPI(http.StatusBadRequest, "Slug not validated", ""))
		return
	}
	livetvSlug = slugSplits[1]
	// Check slug LiveTV / EPG
	slugSplits = strings.Split(livetvSlug, "/")
	if len(slugSplits) == 2 {
		// Lay EPG
		livetvSlug = slugSplits[0]
		epgSlug = slugSplits[1]
	}

	livetvSlugSeo := fmt.Sprintf(SEO_LIVETV_URL, livetvSlug)

	livetv_id, err := GetIdBySlug(livetvSlugSeo)
	if err != nil {
		c.JSON(http.StatusBadRequest, FormatResultAPI(http.StatusBadRequest, err.Error(), ""))
		return
	}
	if epgSlug != "" {
		livetvEpgSlugSeo := fmt.Sprintf(SEO_LIVETV_EPG_URL, livetvSlug, epgSlug)
		c.Set("epg_slug", livetvEpgSlugSeo)
	}
	c.Set("livetv_id", livetv_id)
	GetDetailLiveTVByID(c)
}

func GetLiveTVByGroup(c *gin.Context) {
	platform := c.DefaultQuery("platform", "android")
	cache := StringToBool(c.DefaultQuery("cache", "true"))
	livetv_group_id := c.Param("livetv_group_id")
	user_id := c.GetString("user_id")
	page, err := StringToInt(c.DefaultQuery("page", "0"))
	if err != nil || page > PAGE_MAX {
		page = 0
	}

	limit, err := StringToInt(c.DefaultQuery("limit", "60"))
	if err != nil || limit > LIMIT_MAX {
		limit = 43
	}

	keyCache := LOCAL_LIVE_TV_BY_GROUP + "_" + livetv_group_id + "_" + fmt.Sprint(page) + "_" + fmt.Sprint(limit)
	if cache {
		valC, err := LocalCache.GetValue(keyCache)
		if err == nil {
			var LiveTVObjectOutput LiveTVObjectOutputStruct
			dataByte, _ := json.Marshal(valC)
			json.Unmarshal([]byte(dataByte), &LiveTVObjectOutput)

			//Lay danh sach kenh yeu thich cua user
			LiveTVObjectOutput.Items = GetFavoriteListLiveTV(LiveTVObjectOutput.Items, user_id)

			LiveTVObjectOutput, _ = CheckTagVipOutputLiveTV(LiveTVObjectOutput, c)
			c.JSON(http.StatusOK, FormatResultAPI(http.StatusOK, "", LiveTVObjectOutput))
			return
		}
	}

	dataLiveTV, err := GetLiveTVByGroupId(livetv_group_id, platform, page, limit, cache)
	if err != nil {
		c.JSON(http.StatusBadRequest, FormatResultAPI(http.StatusBadRequest, err.Error(), ""))
		return
	}

	// Write local data
	LocalCache.SetValue(keyCache, dataLiveTV, TTL_LOCALCACHE)

	//Lay danh sach kenh yeu thich cua user
	dataLiveTV.Items = GetFavoriteListLiveTV(dataLiveTV.Items, user_id)

	dataLiveTV, _ = CheckTagVipOutputLiveTV(dataLiveTV, c)
	c.JSON(http.StatusOK, FormatResultAPI(http.StatusOK, "", dataLiveTV))
}

func GetLiveTVEPG(c *gin.Context) {
	livetv_id := c.DefaultQuery("livetv_id", "")
	start, _ := StringToInt(c.DefaultQuery("start", ""))
	end, _ := StringToInt(c.DefaultQuery("end", ""))
	str_date := c.DefaultQuery("str_date", "")
	cache := StringToBool(c.DefaultQuery("cache", "true"))

	//get date from timestamp type d/m/Y
	// strDateFromStart := GetDateFromTimestamp(int64(start))
	// keyCache := LOCAL_LIVE_TV_EPG + "_" + livetv_id + "_" + strDateFromStart

	if str_date != "" {
		start, end = GetTimeStampStartEndFromString(str_date)
		if start <= 0 || end <= 0 {
			c.JSON(http.StatusBadRequest, FormatResultAPI(http.StatusBadRequest, "Input str_date not correct", ""))
			return
		}

	} else {
		//sync start to 0h of start and end to 23h59 of end
		start, _ = GetStartEndTimeFromTimestamp(int64(start))
		_, end = GetStartEndTimeFromTimestamp(int64(end))
	}

	var keyCache = LOCAL_LIVE_TV_EPG + "_" + livetv_id + "_" + fmt.Sprint(start) + "_" + fmt.Sprint(end)

	if cache {
		valC, err := LocalCache.GetValue(keyCache)
		if err == nil {
			c.JSON(http.StatusOK, FormatResultAPI(http.StatusOK, "", valC))
			return
		}
	}

	dataLiveEpg, err := GetEpgByLivetvID(livetv_id, start, end, cache)
	if err != nil {
		c.JSON(http.StatusBadRequest, FormatResultAPI(http.StatusBadRequest, err.Error(), ""))
		return
	}
	// Write local data
	LocalCache.SetValue(keyCache, dataLiveEpg, TTL_LOCALCACHE)
	c.JSON(http.StatusOK, FormatResultAPI(http.StatusOK, "", dataLiveEpg))
}

func AddLivetvFavoriteByUser(c *gin.Context) {
	livetv_ids := c.PostForm("livetv_ids")
	user_id := c.GetString("user_id")
	err := AddLivetvFavorite(user_id, livetv_ids)
	if err != nil {
		c.JSON(http.StatusBadRequest, FormatResultAPI(http.StatusBadRequest, err.Error(), ""))
		return
	}
	c.JSON(http.StatusOK, FormatResultAPI(http.StatusOK, "", "success"))
}

func GetLivetvFavoriteByUser(c *gin.Context) {
	platform := c.DefaultQuery("platform", "android")
	cacheActive := StringToBool(c.DefaultQuery("cache", "true"))
	user_id := c.GetString("user_id")
	dataLiveTV, err := GetLiveFavorite(user_id, platform, cacheActive)
	if err != nil {
		c.JSON(http.StatusBadRequest, FormatResultAPI(http.StatusBadRequest, err.Error(), ""))
		return
	}

	dataLiveTV = GetFavoriteListLiveTV(dataLiveTV, user_id)
	c.JSON(http.StatusOK, FormatResultAPI(http.StatusOK, "", dataLiveTV))
}

func AddLivetvWatchedByUser(c *gin.Context) {
	livetv_ids := c.PostForm("livetv_ids")
	user_id := c.GetString("user_id")
	err := AddLivetvWatched(user_id, livetv_ids)
	if err != nil {
		c.JSON(http.StatusBadRequest, FormatResultAPI(http.StatusBadRequest, err.Error(), ""))
		return
	}
	c.JSON(http.StatusOK, FormatResultAPI(http.StatusOK, "", "success"))
}

func GetLivetvWatchedByUser(c *gin.Context) {
	platform := c.DefaultQuery("platform", "android")
	cache := StringToBool(c.DefaultQuery("cache", "true"))
	user_id := c.GetString("user_id")
	dataLiveTV, err := GetLivetvWatched(user_id, platform, cache)
	if err != nil {
		c.JSON(http.StatusBadRequest, FormatResultAPI(http.StatusBadRequest, err.Error(), ""))
		return
	}

	dataLiveTV = GetFavoriteListLiveTV(dataLiveTV, user_id)
	c.JSON(http.StatusOK, FormatResultAPI(http.StatusOK, "", dataLiveTV))
}

func CheckTagVipOutputLiveTV(LiveTVObjectOutput LiveTVObjectOutputStruct, c *gin.Context) (LiveTVObjectOutputTemp LiveTVObjectOutputStruct, err error) {
	LiveTVObjectOutputTemp = LiveTVObjectOutput

	userId := c.GetString("user_id")

	// Kiểm tra User chưa login => Return
	if strings.HasPrefix(userId, "anonymous_") == true || userId == "" {
		return LiveTVObjectOutputTemp, nil
	}

	var userType bool = false // chua mua goi
	var Sub Subscription.SubcriptionObjectStruct
	SubcriptionsOfUser, err := Sub.GetListByUserId(userId)
	if err == nil && len(SubcriptionsOfUser) > 0 {
		userType = true // Co gói đang active
	}

	// Kiểm tra User chưa mua gói => Return
	if userType == false {
		return LiveTVObjectOutputTemp, nil
	}

	var userIsVipOrPremium bool = false // False: Premium , True: VIP
	for _, val := range SubcriptionsOfUser {
		if val.Type == 1 {
			userIsVipOrPremium = true
			break
		}
	}

	// Kiểm tra User có gói VIP
	// chuyển is_premium = 1 => 0
	if userIsVipOrPremium {
		Items := []LiveTVObjectStruct{}
		for i := 0; i < len(LiveTVObjectOutput.Items); i++ {
			var Item = LiveTVObjectOutput.Items[i]
			if Item.Is_premium == 1 {
				// La content có tag VIP
				Item.Is_premium = 0
			}
			Items = append(Items, Item)
		}
		LiveTVObjectOutputTemp.Items = Items
		return LiveTVObjectOutputTemp, nil
	}

	// User Premium
	// Kiểm tra từng item is_premium = 1, kiểm tra gói của content và gói của user
	listContentPremium, _ := Packages_premium.GetListContentIdInPackagePremium(true)

	Items := []LiveTVObjectStruct{}
	for i := 0; i < len(LiveTVObjectOutput.Items); i++ {
		var Item = LiveTVObjectOutput.Items[i]
		if Item.Is_premium == 1 {
			if ok, _ := In_array(Item.Id, listContentPremium); ok {
				Item.Is_premium = 0
			}
		}
		Items = append(Items, Item)
	}

	LiveTVObjectOutputTemp.Items = Items
	return LiveTVObjectOutputTemp, nil
}
