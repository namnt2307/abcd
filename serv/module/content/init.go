package content

import (
	"fmt"
	"log"
	"math"
	"net/http"
	"strconv"
	"strings"
	"github.com/gin-gonic/gin"
	"github.com/gosimple/slug"

	. "cm-v5/schema"
	. "cm-v5/serv/module"

	seo "cm-v5/serv/module/seo"
	recommendation "cm-v5/serv/module/recommendation"
	vod "cm-v5/serv/module/vod"
	jsoniter "github.com/json-iterator/go"
	Subscription "cm-v5/serv/module/subscription"
	Packages_premium "cm-v5/serv/module/packages_premium"
	"gopkg.in/mgo.v2/bson"
)

var mRedis RedisModelStruct
var mRedisKV RedisKVModelStruct
var json = jsoniter.ConfigCompatibleWithStandardLibrary
var LOGO_LINK_4K string

func init() {
	LOGO_LINK_4K, _ = CommonConfig.GetString("SOFTLOGO_4K", "link_4k")
	if LOGO_LINK_4K == "" {
		LOGO_LINK_4K = "https://static.vieon.vn/vieplay-image/icon/2020/05/11/0wgsp99s_logo_vieon_4k_new_topright_new.png"
	}

	// Load top views content
	go InitContentTopViews()
}

func GetContentById(c *gin.Context) {
	userId := c.GetString("user_id")
	contentId := c.Param("content_id")
	if contentId == "" {
		contentId = c.GetString("content_id")
	}

	platform := c.DefaultQuery("platform", "web")
	cacheActive := StringToBool(c.DefaultQuery("cache", "true"))
	ipUser, _ := GetClientIPHelper(c.Request, c)

	recommendationId := c.DefaultQuery("recommendation_id", "")
	if recommendationId == "" {
		recommendationId = c.GetString("recommendation_id")
	}

	typeRec := c.DefaultQuery("type", "")
	if typeRec == "" {
		typeRec = c.GetString("type")
	}

	// Hit local data
	var keyCache = LOCAL_CONTENT + "_" + platform + "_" + contentId

	var SubtitlesEmpty []SubtitlesOuputStruct
	var AudiosEmpty []AudiosOutputStruct

	if cacheActive {
		valC, err := LocalCache.GetValue(keyCache)
		if err == nil {
			var dataContent ContentObjOutputStruct
			dataByte, _ := json.Marshal(valC)
			json.Unmarshal([]byte(dataByte), &dataContent)
			dataContent, err = GetUserInfoDataContent(userId, dataContent, ipUser)
			if err != nil {
				c.JSON(http.StatusBadRequest, FormatResultAPI(http.StatusBadRequest, err.Error(), ""))
				return
			}

			//Set empty content ouput subtitle + audio + ads
			dataContent.Subtitles = SubtitlesEmpty
			dataContent.Audios = AudiosEmpty

			dataContent , _ = CheckTagVipOutputInfo(dataContent , c)
			TrackingContent(c, contentId, recommendationId, typeRec)
			c.JSON(http.StatusOK, FormatResultAPI(http.StatusOK, "", dataContent))
			return
		}
	}

	dataContent, err := GetContent(contentId, platform, cacheActive)
	if err != nil {
		c.JSON(http.StatusBadRequest, FormatResultAPI(http.StatusBadRequest, err.Error(), ""))
		return
	}

	// Write local data
	LocalCache.SetValue(keyCache, dataContent, TTL_LOCALCACHE)
	dataContent, err = GetUserInfoDataContent(userId, dataContent, ipUser)
	if err != nil {
		c.JSON(http.StatusBadRequest, FormatResultAPI(http.StatusBadRequest, err.Error(), ""))
		return
	}

	//Set empty content ouput subtitle + audio + ads
	dataContent.Subtitles = SubtitlesEmpty
	dataContent.Audios = AudiosEmpty

	dataContent , _ = CheckTagVipOutputInfo(dataContent , c)
	TrackingContent(c, contentId, recommendationId, typeRec)
	c.JSON(http.StatusOK, FormatResultAPI(http.StatusOK, "", dataContent))
}

func GetContentBySlug(c *gin.Context) {
	entitySlug := c.PostForm("entity_slug")
	platform := c.DefaultQuery("platform", "web")
	recommendationId := c.PostForm("recommendation_id")
	typeRec := c.PostForm("type")

	if entitySlug == "" {
		c.JSON(http.StatusBadRequest, FormatResultAPI(http.StatusBadRequest, "entity_slug not valid", ""))
		return
	}


	// Check slug exists
	vodIdSEO := seo.CheckExistsSEOBySlug(entitySlug)
	if vodIdSEO != "" {
		c.Set("content_id", vodIdSEO)
		c.Set("recommendation_id", recommendationId)
		c.Set("type", typeRec)
		GetContentById(c)
		return
	}


	platformInfo := Platform(platform)
	vodId, err := vod.GetVodIdBySlug(entitySlug, platformInfo.Id)
	if err != nil {
		c.JSON(http.StatusBadRequest, FormatResultAPI(http.StatusBadRequest, err.Error(), ""))
		return
	}

	if vodId == "" {
		c.JSON(http.StatusBadRequest, FormatResultAPI(http.StatusBadRequest, "content not found", ""))
		return
	}

	c.Set("content_id", vodId)
	c.Set("recommendation_id", recommendationId)
	c.Set("type", typeRec)
	GetContentById(c)
}

func TrackingContent(c *gin.Context, contentID, recommendationId, typeRec string) {
	var TrackingData TrackingDataStruct
	TrackingUserData, err := recommendation.GetTracking(c, contentID, TrackingData)
	TrackingUserData.Tracking_id = recommendationId
	TrackingUserData.Tracking_type = typeRec

	if err == nil {
		recommendation.PushDataToKafka(REC_CLICK, TrackingUserData)
	}
}
func GetContentTipById(c *gin.Context) {
	contentId := c.Param("content_id")
	platform := c.DefaultQuery("platform", "web")
	epsId := c.DefaultQuery("eps_id", "")
	cacheActive := StringToBool(c.DefaultQuery("cache", "true"))
	userId := c.GetString("user_id")

	dataContent, err := GetContentTip(contentId, epsId, userId, platform, cacheActive)
	if err != nil {
		Sentry_log(err)
		c.JSON(http.StatusBadRequest, FormatResultAPI(http.StatusBadRequest, err.Error(), ""))
		return
	}

	c.JSON(http.StatusOK, FormatResultAPI(http.StatusOK, "", dataContent))
}

/**
URL: /content_detail/:content_id?eps_id=:eps_id
*/
func GetContentDetailById(c *gin.Context) {
	userId := c.GetString("user_id")
	typeUser := c.GetInt("user_is_premium")
	contentId := c.Param("content_id")
	epsId := c.DefaultQuery("eps_id", "")
	platform := c.DefaultQuery("platform", "web")
	cacheActive := StringToBool(c.DefaultQuery("cache", "true"))
	ipUser, _ := GetClientIPHelper(c.Request, c)
	modelPlatform := c.DefaultQuery("model", "")

	dataOutput, err := MapContentInfoPersonal(contentId, epsId, userId, platform, modelPlatform, ipUser, fmt.Sprint(typeUser), cacheActive)
	if err != nil {
		Sentry_log(err)
		c.JSON(http.StatusBadRequest, FormatResultAPI(http.StatusBadRequest, err.Error(), ""))
		return
	}

	// Check user agent
	if CheckUserAgentAndroidAndIOSVersionLock(platform, c.Request.UserAgent()) {
		var LinkPlay LinkPlayStruct
		dataOutput.Link_play = LinkPlay
	}

	// DT-12265 : add views
	dataOutput.Views = GetRedisViews(contentId)
	dataOutput.UserInteractionCount = math.Round(float64(dataOutput.Views) * 2.3)
	c.JSON(http.StatusOK, FormatResultAPI(http.StatusOK, "", dataOutput))
}

// DT-12265 : add views
func GetRedisViews(contentId string) int64 {
	var views int64
	keyCache := "VIEW_" + contentId
	views = mRedisKV.Incr(keyCache)
	if views == 2 {
		// Get from DB
		//Clear cache ribbons v3
		session, db, err := GetCollection()
		if err != nil {
			Sentry_log(err)
			log.Println("AddView err", err)
		}
		defer session.Close()
		var where = bson.M{
			"contentid": contentId,
		}
		var viewObj ViewObjectStruct
		err = db.C(COLLECTION_VIEW).Find(where).One(&viewObj)
		valueFromDb := int64(0)
		if err == nil {
			valueFromDb = viewObj.View
		}

		views = mRedisKV.IncrBy(keyCache, valueFromDb)
	}
	if views%147 == 0 && views > 0 {
		var viewObj ViewObjectStruct
		viewObj.View = views
		viewObj.ContentId = contentId
		// Write DB
		go func(viewObj ViewObjectStruct) {
			//connect db
			session, db, err := GetCollection()
			if err != nil {
				Sentry_log(err)
				log.Println("AddView err", err)
				return
			}
			defer session.Close()
			var where = bson.M{
				"contentid": viewObj.ContentId,
			}
			_, err = db.C(COLLECTION_VIEW).Upsert(where, viewObj)

			if err != nil {
				log.Println("AddView err", err)
				return
			}

		}(viewObj)

	}

	return views

}

func GetEpisodeById(c *gin.Context) {
	contentId := c.Param("content_id")
	userId := c.GetString("user_id")
	page, err := StringToInt(c.DefaultQuery("page", "0"))
	if err != nil || page > PAGE_MAX {
		page = 0
	}
	limit, err := StringToInt(c.DefaultQuery("limit", "30"))
	if err != nil || limit > LIMIT_MAX {
		limit = 30
	}
	platform := c.DefaultQuery("platform", "web")
	cacheActive := StringToBool(c.DefaultQuery("cache", "true"))

	// Hit local data
	var keyCache = LOCAL_LIST_EPISODE + "_" + platform + "_" + contentId + "_" + strconv.Itoa(page) + "_" + strconv.Itoa(limit)
	if cacheActive {
		valC, err := LocalCache.GetValue(keyCache)
		if err == nil {
			var dataEpisode EpisodeObjOutputStruct
			dataByte, _ := json.Marshal(valC)
			json.Unmarshal([]byte(dataByte), &dataEpisode)
			dataEpisode = GetTrackingDataEpisode(platform, dataEpisode)
			dataEpisode = GetUserInfoDataEpisode(userId, dataEpisode)
			dataEpisode , _ = CheckTagVipOutputEpisode(dataEpisode , c)
			
			c.JSON(http.StatusOK, FormatResultAPI(1, "", dataEpisode))
			return
		}
	}

	dataEpisode, err := GetEpisode(contentId, "", page, limit, platform, cacheActive)
	if err != nil {
		Sentry_log(err)
		c.JSON(http.StatusBadRequest, FormatResultAPI(0, err.Error(), ""))
		return
	}

	// Write local data
	LocalCache.SetValue(keyCache, dataEpisode, TTL_LOCALCACHE)
	dataEpisode = GetTrackingDataEpisode(platform, dataEpisode)
	dataEpisode = GetUserInfoDataEpisode(userId, dataEpisode)
	dataEpisode , _ = CheckTagVipOutputEpisode(dataEpisode , c)

	c.JSON(http.StatusOK, FormatResultAPI(1, "", dataEpisode))
}

func GetEpisodeBySlug(c *gin.Context) {
	entitySlug := c.PostForm("entity_slug")
	if entitySlug == "" {
		c.JSON(http.StatusBadRequest, FormatResultAPI(http.StatusBadRequest, "entity_slug not valid", ""))
		return
	}

	userId := c.GetString("user_id")
	page, err := StringToInt(c.DefaultQuery("page", "0"))
	if err != nil || page > PAGE_MAX {
		page = 0
	}
	limit, err := StringToInt(c.DefaultQuery("limit", "30"))
	if err != nil || limit > LIMIT_MAX {
		limit = 30
	}
	platform := c.DefaultQuery("platform", "web")
	cacheActive := StringToBool(c.DefaultQuery("cache", "true"))

	// Hit local data
	var keyCache = LOCAL_LIST_EPISODE + "_" + platform + "_" + entitySlug + "_" + strconv.Itoa(page) + "_" + strconv.Itoa(limit)
	if cacheActive {
		valC, err := LocalCache.GetValue(keyCache)
		if err == nil {
			var dataEpisode EpisodeObjOutputStruct
			dataByte, _ := json.Marshal(valC)
			json.Unmarshal([]byte(dataByte), &dataEpisode)
			dataEpisode = GetTrackingDataEpisode(platform, dataEpisode)
			dataEpisode = GetUserInfoDataEpisode(userId, dataEpisode)
			dataEpisode , _ = CheckTagVipOutputEpisode(dataEpisode , c)

			c.JSON(http.StatusOK, FormatResultAPI(1, "", dataEpisode))
			return
		}
	}

	dataEpisode, err := GetEpisode("", entitySlug, page, limit, platform, cacheActive)
	if err != nil {
		Sentry_log(err)
		c.JSON(http.StatusBadRequest, FormatResultAPI(http.StatusBadRequest, err.Error(), ""))
		return
	}

	// Write local data
	LocalCache.SetValue(keyCache, dataEpisode, TTL_LOCALCACHE)
	dataEpisode = GetTrackingDataEpisode(platform, dataEpisode)
	dataEpisode = GetUserInfoDataEpisode(userId, dataEpisode)
	dataEpisode , _ = CheckTagVipOutputEpisode(dataEpisode , c)

	c.JSON(http.StatusOK, FormatResultAPI(1, "", dataEpisode))
}

func GetEpisodeRangeById(c *gin.Context) {
	contentId := c.Param("content_id")
	platform := c.DefaultQuery("platform", "web")
	cacheActive := StringToBool(c.DefaultQuery("cache", "true"))

	// Hit local data
	var keyCache = LOCAL_LIST_EPISODE_RANGE + "_" + platform + "_" + contentId

	if cacheActive {
		valC, err := LocalCache.GetValue(keyCache)
		if err == nil {
			c.JSON(http.StatusOK, FormatResultAPI(1, "", valC))
			return
		}
	}

	dataEpisodeRange := GetRangeEpisode(contentId, platform, cacheActive)

	// Write local data
	LocalCache.SetValue(keyCache, dataEpisodeRange, TTL_LOCALCACHE)
	c.JSON(http.StatusOK, FormatResultAPI(1, "", dataEpisodeRange))
}

func GetRelatedById(c *gin.Context) {
	contentId := c.Param("content_id")
	if contentId == "" {
		contentId = c.GetString("content_id")
	}
	page, err := StringToInt(c.DefaultQuery("page", "0"))
	if err != nil || page > PAGE_MAX {
		page = 0
	}
	limit, err := StringToInt(c.DefaultQuery("limit", "30"))
	if err != nil || limit > LIMIT_MAX {
		limit = 43
	}
	platform := c.DefaultQuery("platform", "web")
	// cacheActive := StringToBool(c.DefaultQuery("cache", "true"))

	user_id := c.GetString("user_id")
	device_id := c.GetString("device_id")

	dataRelated, err := GetRelatedYusp(user_id, device_id, platform, contentId, page, limit)
	if err != nil {
		Sentry_log(err)
		c.JSON(http.StatusBadRequest, FormatResultAPI(http.StatusBadRequest, err.Error(), ""))
		return
	}

	/**************************************************/
	// Hit local data
	// var keyCache = LOCAL_RELATED + "_" + fmt.Sprint(page) + "_" + fmt.Sprint(limit) + "_" + platform + "_" + contentId
	// if cacheActive {
	// 	valC, err := LocalCache.GetValue(keyCache)
	// 	if err == nil {
	// 		var RelatedObjOutput RelatedObjOutputStruct
	// 		json.Unmarshal([]byte(valC.(string)), &RelatedObjOutput)
	// 		dataRelated := CheckLocationReleated(RelatedObjOutput, c)
	// 		c.JSON(http.StatusOK, FormatResultAPI(http.StatusOK, "", dataRelated))
	// 		return
	// 	}
	// }

	// dataRelated, err := GetRelated(contentId, page, limit, platform, cacheActive)
	// if err != nil {
	// 	Sentry_log(err)
	// 	c.JSON(http.StatusBadRequest, FormatResultAPI(http.StatusBadRequest, err.Error(), ""))
	// 	return
	// }

	// // Write local data
	// dataByte, _ := json.Marshal(dataRelated)
	// LocalCache.SetValue(keyCache, string(dataByte), TTL_LOCALCACHE)
	/**************************************************/

	dataRelated = CheckLocationReleated(dataRelated, c)
	dataRelated, _ = CheckTagVipOutputReleated(dataRelated, c)

	c.JSON(http.StatusOK, FormatResultAPI(http.StatusOK, "", dataRelated))
}

func GetRelatedBySlug(c *gin.Context) {
	entitySlug := c.PostForm("entity_slug")
	if entitySlug == "" {
		c.JSON(http.StatusBadRequest, FormatResultAPI(http.StatusBadRequest, "Bad entity_slug", ""))
		return
	}
	platform := c.DefaultQuery("platform", "web")
	platformInfo := Platform(platform)
	vodId, err := vod.GetVodIdBySlug(entitySlug, platformInfo.Id)

	if err != nil {
		c.JSON(http.StatusBadRequest, FormatResultAPI(http.StatusBadRequest, err.Error(), ""))
		return
	}

	c.Set("content_id", vodId)
	GetRelatedById(c)
}

func GetRelatedVideosBySlug(c *gin.Context) {
	entitySlug := c.PostForm("entity_slug")
	platform := c.DefaultQuery("platform", "web")

	platformInfo := Platform(platform)
	vodId, err := vod.GetVodIdBySlug(entitySlug, platformInfo.Id)
	if err != nil {
		c.JSON(http.StatusBadRequest, FormatResultAPI(http.StatusBadRequest, err.Error(), ""))
		return
	}

	c.Set("group_id", vodId)
	GetRelatedVideosById(c)
}

func GetRelatedVideosById(c *gin.Context) {
	groupId := c.Param("group_id")
	if groupId == "" {
		groupId = c.GetString("group_id")
	}
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

	// Hit local data
	var keyCache = "related_videos_" + groupId + "_" + platform + "_" + strconv.Itoa(page) + "_" + strconv.Itoa(limit)

	if cacheActive {
		valC, err := LocalCache.GetValue(keyCache)
		if err == nil {
			var RelatedVideosOutput RelatedVideosObjOutputStruct
			json.Unmarshal([]byte(valC.(string)), &RelatedVideosOutput)
			dataRelatedVideos := CheckLocationReleatedVideos(RelatedVideosOutput, c)
			dataRelatedVideos, _ = CheckTagVipOutputReleatedVideos(dataRelatedVideos, c)
			c.JSON(http.StatusOK, FormatResultAPI(http.StatusOK, "", dataRelatedVideos))
			return
		}
	}

	dataRelatedVideos, err := GetRelatedVideos(groupId, page, limit, platform, cacheActive)
	if err != nil {
		Sentry_log(err)
		c.JSON(http.StatusBadRequest, FormatResultAPI(http.StatusBadRequest, err.Error(), ""))
		return
	}

	// Write local data
	dataByte, _ := json.Marshal(dataRelatedVideos)
	LocalCache.SetValue(keyCache, string(dataByte), TTL_LOCALCACHE)

	dataRelatedVideos = CheckLocationReleatedVideos(dataRelatedVideos, c)
	dataRelatedVideos, _ = CheckTagVipOutputReleatedVideos(dataRelatedVideos, c)
	
	c.JSON(http.StatusOK, FormatResultAPI(http.StatusOK, "", dataRelatedVideos))
}

func GetTrackingDataEpisode(platform string, dataEpisode EpisodeObjOutputStruct) EpisodeObjOutputStruct {
	trackingType := platform + "_" + RECOMMENDATION_CONTENT_LISTING
	trackingType = slug.Make(trackingType)
	dataEpisode.Tracking_data = recommendation.GetRandomDefaultTrackingData(strings.ToUpper(trackingType))
	return dataEpisode
}

func CheckLocationReleated(RelatedObjOutput RelatedObjOutputStruct, c *gin.Context) RelatedObjOutputStruct {
	ipUser, _ := GetClientIPHelper(c.Request, c)
	// Check localtion item
	valCheckIpVN := 2
	total := 0

	var ReleatedItemTemps []ItemOutputObjStruct
	for i, s := 0, len(RelatedObjOutput.Items); i < s; i++ {
		if RelatedObjOutput.Items[i].Geo_check == 1 {
			// Content nay  Chi duoc xem o VN
			if valCheckIpVN == 2 {
				valCheckIpVN = 1 // ip vn
				// Check ip user is VN
				if CheckIpIsVN(ipUser) == false {
					valCheckIpVN = 0 // ip nuoc ngoai
				}
			}
			if valCheckIpVN == 0 {
				total = total + 1
			} else {
				ReleatedItemTemps = append(ReleatedItemTemps, RelatedObjOutput.Items[i])
			}
		} else {
			ReleatedItemTemps = append(ReleatedItemTemps, RelatedObjOutput.Items[i])
		}
	}

	RelatedObjOutput.Items = ReleatedItemTemps
	RelatedObjOutput.Metadata.Total = RelatedObjOutput.Metadata.Total - total

	return RelatedObjOutput
}

func CheckLocationReleatedVideos(RelatedVideosObjOutput RelatedVideosObjOutputStruct, c *gin.Context) RelatedVideosObjOutputStruct {
	ipUser, _ := GetClientIPHelper(c.Request, c)
	// Check localtion item
	valCheckIpVN := 2
	total := 0

	var ReleatedItemVideoTemps []ItemOutputObjStruct
	for i, s := 0, len(RelatedVideosObjOutput.Items); i < s; i++ {
		if RelatedVideosObjOutput.Items[i].Geo_check == 1 {
			// Content nay  Chi duoc xem o VN
			if valCheckIpVN == 2 {
				valCheckIpVN = 1 // ip vn
				// Check ip user is VN
				if CheckIpIsVN(ipUser) == false {
					valCheckIpVN = 0 // ip nuoc ngoai
				}
			}
			if valCheckIpVN == 0 {
				total = total + 1
			} else {
				ReleatedItemVideoTemps = append(ReleatedItemVideoTemps, RelatedVideosObjOutput.Items[i])
			}
		} else {
			ReleatedItemVideoTemps = append(ReleatedItemVideoTemps, RelatedVideosObjOutput.Items[i])
		}
	}

	RelatedVideosObjOutput.Items = ReleatedItemVideoTemps
	RelatedVideosObjOutput.Metadata.Total = RelatedVideosObjOutput.Metadata.Total - total

	return RelatedVideosObjOutput
}

func GetVodById(c *gin.Context) {
	// userId := c.GetString("user_id")
	contentId := c.Param("content_id")
	epsId := c.DefaultQuery("eps_id", "")
	contentTracking := contentId
	if contentId == "" {
		contentId = c.GetString("content_id")
	}
	if epsId == "" {
		epsId = c.GetString("eps_id")
		contentTracking = epsId
	}

	platform := c.DefaultQuery("platform", "web")
	cacheActive := StringToBool(c.DefaultQuery("cache", "true"))
	// ipUser, _ := GetClientIPHelper(c.Request, c)

	recommendationId := c.DefaultQuery("recommendation_id", "")
	if recommendationId == "" {
		recommendationId = c.GetString("recommendation_id")
	}

	typeRec := c.DefaultQuery("type", "")
	if typeRec == "" {
		typeRec = c.GetString("type")
	}

	// Hit local data
	var keyCache = LOCAL_VOD_V3 + "_" + platform + "_" + contentId + "_" + epsId

	if cacheActive {
		valC, err := LocalCache.GetValue(keyCache)
		if err == nil {
			TrackingContent(c, contentTracking, recommendationId, typeRec)
			c.JSON(http.StatusOK, FormatResultAPI(http.StatusOK, "", valC))
			return
		}
	}

	dataContent, err := GetInfoVod(contentId, epsId, platform, cacheActive)
	if err != nil {
		c.JSON(http.StatusBadRequest, FormatResultAPI(http.StatusBadRequest, err.Error(), ""))
		return
	}

	// Write local data
	LocalCache.SetValue(keyCache, dataContent, TTL_LOCALCACHE)
	// dataContent, err = GetUserInfoDataContent(userId, dataContent, ipUser)
	// if err != nil {
	// 	c.JSON(http.StatusBadRequest, FormatResultAPI(http.StatusBadRequest, err.Error(), ""))
	// 	return
	// }

	TrackingContent(c, contentTracking, recommendationId, typeRec)
	c.JSON(http.StatusOK, FormatResultAPI(http.StatusOK, "", dataContent))
}

func GetVodDetailById(c *gin.Context) {
	userId := c.GetString("user_id")
	contentId := c.Param("content_id")
	epsId := c.DefaultQuery("eps_id", "")

	platform := c.DefaultQuery("platform", "web")
	cacheActive := StringToBool(c.DefaultQuery("cache", "true"))
	ipUser, _ := GetClientIPHelper(c.Request, c)
	modelPlatform := c.DefaultQuery("model", "")

	dataContent, err := GetVodInfoPersonal(contentId, epsId, userId, platform, modelPlatform, ipUser, cacheActive)
	if err != nil {
		Sentry_log(err)
		c.JSON(http.StatusBadRequest, FormatResultAPI(http.StatusBadRequest, err.Error(), ""))
		return
	}

	c.JSON(http.StatusOK, FormatResultAPI(http.StatusOK, "", dataContent))
}

func GetVodBySlug(c *gin.Context) {
	entitySlug := c.PostForm("entity_slug")
	platform := c.DefaultQuery("platform", "web")
	recommendationId := c.PostForm("recommendation_id")
	typeRec := c.PostForm("type")

	// Check slug validate
	slugSplits := strings.Split(entitySlug, "/")
	if len(slugSplits) < 3 {
		c.JSON(http.StatusBadRequest, FormatResultAPI(http.StatusBadRequest, "Slug not validated", ""))
		return
	}
	arr := strings.Split(entitySlug, "/")
	//len = 5 slug of episode
	if len(arr) == 5 {
		arr = arr[:len(arr)-1]
	}
	var entitySlugSeason = strings.Join(arr, "/")

	platformInfo := Platform(platform)
	vodId, err := vod.GetVodIdBySlug(entitySlug, platformInfo.Id)
	if err != nil {
		c.JSON(http.StatusBadRequest, FormatResultAPI(http.StatusBadRequest, err.Error(), ""))
		return
	}
	//set recommendation_id and type
	c.Set("recommendation_id", recommendationId)
	c.Set("type", typeRec)

	//check if request for slug episode
	if entitySlugSeason != entitySlug {

		seasonId, err := vod.GetVodIdBySlug(entitySlugSeason, platformInfo.Id)
		if err != nil {
			c.JSON(http.StatusBadRequest, FormatResultAPI(http.StatusBadRequest, err.Error(), ""))
			return
		}
		c.Set("content_id", seasonId)
		c.Set("eps_id", vodId)

		GetVodById(c)
		return
	}
	c.Set("content_id", vodId)
	GetVodById(c)
}

func GetViewsContents(c *gin.Context) {
	var listViews []ViewObjectStruct
	contentIds := c.PostForm("content_ids")
	if contentIds != "" {
		arrIDContent, _ := ParseStringToArray(contentIds)
		for _, val := range arrIDContent {
			var viewObj ViewObjectStruct
			keyCache := "VIEW_" + val
			views, _ := mRedis.GetInt(keyCache)
			viewObj.View = int64(views)
			viewObj.ContentId = val

			listViews = append(listViews, viewObj)

		}
	}

	c.JSON(http.StatusOK, FormatResultAPI(http.StatusOK, "", listViews))
}

func CheckTagVipOutputReleated(RelatedObjOutput RelatedObjOutputStruct, c *gin.Context) (RelatedObjOutputTemp RelatedObjOutputStruct, err error) {
	RelatedObjOutputTemp = RelatedObjOutput

	userId := c.GetString("user_id")

	// Kiểm tra User chưa login => Return
	if strings.HasPrefix(userId, "anonymous_") == true || userId == "" {
		return RelatedObjOutputTemp, nil
	}

	var userType bool = false // chua mua goi
	var Sub Subscription.SubcriptionObjectStruct
	SubcriptionsOfUser, err := Sub.GetListByUserId(userId)
	if err == nil && len(SubcriptionsOfUser) > 0 {
		userType = true // Co gói đang active
	}

	// Kiểm tra User chưa mua gói => Return
	if userType == false {
		return RelatedObjOutputTemp, nil
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
		Items := []ItemOutputObjStruct{}
		for i := 0; i < len(RelatedObjOutput.Items); i++ {
			var Item = RelatedObjOutput.Items[i]
			if Item.Is_premium == 1 {
				// La content có tag VIP
				Item.Is_premium = 0 
			} 
			Items = append(Items, Item)
		}
		RelatedObjOutputTemp.Items = Items
		return RelatedObjOutputTemp, nil
	}

	// User Premium
	// Kiểm tra từng item is_premium = 1, kiểm tra gói của content và gói của user
	listContentPremium , _ := Packages_premium.GetListContentIdInPackagePremium(true)

	Items := []ItemOutputObjStruct{}
	for i := 0; i < len(RelatedObjOutput.Items); i++ {
		var Item = RelatedObjOutput.Items[i]
		if Item.Is_premium == 1 {
			if ok, _ := In_array(Item.Id, listContentPremium); ok {
				Item.Is_premium = 0
			}
		} 
		Items = append(Items, Item)
	}

	RelatedObjOutputTemp.Items = Items
	return RelatedObjOutputTemp, nil
}

func CheckTagVipOutputReleatedVideos(RelatedVideosObjOutput RelatedVideosObjOutputStruct, c *gin.Context) (RelatedVideosObjOutputTemp RelatedVideosObjOutputStruct, err error) {
	RelatedVideosObjOutputTemp = RelatedVideosObjOutput

	userId := c.GetString("user_id")

	// Kiểm tra User chưa login => Return
	if strings.HasPrefix(userId, "anonymous_") == true || userId == "" {
		return RelatedVideosObjOutputTemp, nil
	}

	var userType bool = false // chua mua goi
	var Sub Subscription.SubcriptionObjectStruct
	SubcriptionsOfUser, err := Sub.GetListByUserId(userId)
	if err == nil && len(SubcriptionsOfUser) > 0 {
		userType = true // Co gói đang active
	}

	// Kiểm tra User chưa mua gói => Return
	if userType == false {
		return RelatedVideosObjOutputTemp, nil
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
		Items := []ItemOutputObjStruct{}
		for i := 0; i < len(RelatedVideosObjOutput.Items); i++ {
			var Item = RelatedVideosObjOutput.Items[i]
			if Item.Is_premium == 1 {
				// La content có tag VIP
				Item.Is_premium = 0 
			} 
			Items = append(Items, Item)
		}
		RelatedVideosObjOutputTemp.Items = Items
		return RelatedVideosObjOutputTemp, nil
	}

	// User Premium
	// Kiểm tra từng item is_premium = 1, kiểm tra gói của content và gói của user
	listContentPremium , _ := Packages_premium.GetListContentIdInPackagePremium(true)

	Items := []ItemOutputObjStruct{}
	for i := 0; i < len(RelatedVideosObjOutput.Items); i++ {
		var Item = RelatedVideosObjOutput.Items[i]
		if Item.Is_premium == 1 {
			if ok, _ := In_array(Item.Id, listContentPremium); ok {
				Item.Is_premium = 0
			}
		} 
		Items = append(Items, Item)
	}

	RelatedVideosObjOutputTemp.Items = Items
	return RelatedVideosObjOutputTemp, nil
}

func CheckTagVipOutputEpisode(DataEpisode EpisodeObjOutputStruct, c *gin.Context) (DataEpisodeTemp EpisodeObjOutputStruct, err error) {
	DataEpisodeTemp = DataEpisode

	userId := c.GetString("user_id")

	// Kiểm tra User chưa login => Return
	if strings.HasPrefix(userId, "anonymous_") == true || userId == "" {
		return DataEpisodeTemp, nil
	}

	var userType bool = false // chua mua goi
	var Sub Subscription.SubcriptionObjectStruct
	SubcriptionsOfUser, err := Sub.GetListByUserId(userId)
	if err == nil && len(SubcriptionsOfUser) > 0 {
		userType = true // Co gói đang active
	}

	// Kiểm tra User chưa mua gói => Return
	if userType == false {
		return DataEpisodeTemp, nil
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
		Items := []ItemsEpisodeObjOutputStruct{}
		for i := 0; i < len(DataEpisode.Items); i++ {
			var Item = DataEpisode.Items[i]
			if Item.Is_premium == 1 {
				// La content có tag VIP
				Item.Is_premium = 0 
			} 
			Items = append(Items, Item)
		}
		DataEpisodeTemp.Items = Items
		return DataEpisodeTemp, nil
	}

	// User Premium
	// Kiểm tra từng item is_premium = 1, kiểm tra gói của content và gói của user
	listContentPremium , _ := Packages_premium.GetListContentIdInPackagePremium(true)

	Items := []ItemsEpisodeObjOutputStruct{}
	for i := 0; i < len(DataEpisode.Items); i++ {
		var Item = DataEpisode.Items[i]
		if Item.Is_premium == 1 {
			if ok, _ := In_array(Item.Id, listContentPremium); ok {
				Item.Is_premium = 0
			}
		} 
		Items = append(Items, Item)
	}

	DataEpisodeTemp.Items = Items
	return DataEpisodeTemp, nil
}

func CheckTagVipOutputInfo(DataContent ContentObjOutputStruct, c *gin.Context) (DataContentTemp ContentObjOutputStruct, err error) {
	DataContentTemp = DataContent

	userId := c.GetString("user_id")

	// Kiểm tra User chưa login => Return
	if strings.HasPrefix(userId, "anonymous_") == true || userId == "" {
		return DataContentTemp, nil
	}

	var userType bool = false // chua mua goi
	var Sub Subscription.SubcriptionObjectStruct
	SubcriptionsOfUser, err := Sub.GetListByUserId(userId)
	if err == nil && len(SubcriptionsOfUser) > 0 {
		userType = true // Co gói đang active
	}

	// Kiểm tra User chưa mua gói => Return
	if userType == false {
		return DataContentTemp, nil
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
		if DataContentTemp.Is_premium == 1 {
			// La content có tag VIP
			DataContentTemp.Is_premium = 0 
		}
		return DataContentTemp, nil
	}

	// User Premium
	// Kiểm tra từng item is_premium = 1, kiểm tra gói của content và gói của user
	listContentPremium , _ := Packages_premium.GetListContentIdInPackagePremium(true)

	if ok, _ := In_array(DataContentTemp.Id, listContentPremium); ok {
		DataContentTemp.Is_premium = 0
	}
	return DataContentTemp, nil
}
