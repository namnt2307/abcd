package watchlater

import (
	"errors"
	"strings"
	"time"
	"fmt"
	. "cm-v5/serv/module"
	. "cm-v5/schema"
	recommendation "cm-v5/serv/module/recommendation"
	vod "cm-v5/serv/module/vod"
	"github.com/gin-gonic/gin"
	"gopkg.in/mgo.v2/bson"
)

type DataResponseStruct struct {
	Status string `json:"status" `
}

func GetWatchlater(userId string, platform string, page int, limit int, cacheActive bool) (WatchlaterContentObjStruct, error) {
	var platformDetail = Platform(platform)
	var WatchlaterContentObj WatchlaterContentObjStruct
	WatchlaterContentObj.Metadata.Limit = limit
	WatchlaterContentObj.Metadata.Page = page
	WatchlaterContentObj.Items = make([]WatchlaterContentDetailObjStruct, 0)

	//Config off ENABLE_WATCHLATER
	if ENABLE_WATCHLATER == "disable" {
		return WatchlaterContentObj, nil
	}

	contentIds, err := GetListIDWatchlaterByUserID(userId, page, limit, cacheActive)
	if len(contentIds) > 0 && err == nil {
		// Get List VOD Detail
		var VODDataObjects []VODDataObjectStruct
		VODDataObjects, err = vod.GetVODByListID(contentIds, platformDetail.Id, 0, cacheActive)
		if err != nil {
			return WatchlaterContentObj, nil
		}

		// Get List Group Id of Episode
		var ListGroupID []string
		for _, vodData := range VODDataObjects {
			if vodData.Type == VOD_TYPE_EPISODE {
				ListGroupID = append(ListGroupID, vodData.Group_id)
			}
		}

		//get data season from groupID
		var VODDataObjectsGetByGroupIDs []VODDataObjectStruct
		if len(ListGroupID) > 0 {
			VODDataObjectsGetByGroupIDs, _ = vod.GetVODByListID(ListGroupID, platformDetail.Id, 0, true)
		}

		// Format data output
		dataByte, _ := json.Marshal(VODDataObjects)
		err = json.Unmarshal(dataByte, &WatchlaterContentObj.Items)
		if err != nil {
			return WatchlaterContentObj, err
		}

		//Switch images platform
		for k, vodData := range VODDataObjects {
			var ImagesData = vodData.Images
			var ImagesMapping ImagesOutputObjectStruct

			//sync images from season to episode
			if vodData.Type == VOD_TYPE_EPISODE && len(VODDataObjectsGetByGroupIDs) > 0 {
				for _, seasonData := range VODDataObjectsGetByGroupIDs {
					if seasonData.Id == vodData.Group_id {
						ImagesData = seasonData.Images
						break
					}
				}
			}

			switch platformDetail.Type {
			case "web":
				ImagesMapping.Vod_thumb = BuildImage(ImagesData.Web.Vod_thumb)
				ImagesMapping.Thumbnail = BuildImage(ImagesData.Web.Vod_thumb)
			case "smarttv":
				ImagesMapping.Thumbnail = BuildImage(ImagesData.Smarttv.Thumbnail)
			case "app":
				ImagesMapping.Thumbnail = BuildImage(ImagesData.App.Thumbnail)
			default:
				ImagesMapping.Vod_thumb = BuildImage(ImagesData.Web.Vod_thumb)
			}

			ImagesMapping = MappingImagesV4(platformDetail.Type, ImagesMapping, ImagesData, true)
			WatchlaterContentObj.Items[k].Images = ImagesMapping
			WatchlaterContentObj.Items[k].Is_watchlater = true
		}

		// Get total
		WatchlaterContentObj.Metadata.Total = GetWatchlaterTotalInCache(userId)
	}

	return WatchlaterContentObj, nil
}

func GetListIDWatchlaterByUserID(userId string, page int, limit int, cacheActive bool) ([]string, error) {
	// Check Exists Data in cache
	if mRedisUSC.Exists(PREFIX_REDIS_ZRANGE_USC_WATCHLATER+userId) == 0 {
		// No have Data in cache
		// Push data from DB
		// PushWatchlaterToCache(userId)
		return make([]string, 0), nil
	}

	var start = page * limit
	var stop = (start + limit) - 1
	contentIds := mRedisUSC.ZRevRange(PREFIX_REDIS_ZRANGE_USC_WATCHLATER+userId, int64(start), int64(stop))

	return contentIds, nil
}

func AddWatchlater(c *gin.Context, userId string, contentId string) (DataResponseStruct, error) {
	var dataResponse DataResponseStruct
	if contentId == "" || userId == "" {
		dataResponse.Status = "fail"
		return dataResponse, errors.New("Data input not empty")
	}

	// Update cache
	status, _ := UpdateWatchlaterCache(userId, contentId)
	var TrackingData TrackingDataStruct
	// //Push data Kafka
	TrackingUserData, err := recommendation.GetTracking(c, contentId, TrackingData)
	if err != nil {
		var eventType string
		if status%2 == 0 {
			//Remove
			eventType = REMOVE_FROM_PLAYLIST
		} else {
			//add
			eventType = ADD_TO_PLAYLIST
		}
		recommendation.PushDataToKafka(eventType, TrackingUserData)
	}

	// Update Mongo
	UpdateWatchlaterDB(userId, contentId, status, 0)

	dataResponse.Status = "success"
	return dataResponse, nil
}

func UpdateWatchlaterCache(userId string, contentId string) (int64, error) {
	var curTime = time.Now().Unix()
	var status int64 = mRedisUSC.Incr(PREFIX_REDIS_USC_WATCHLATER + userId + "_" + contentId)
	fmt.Println(PREFIX_REDIS_USC_WATCHLATER + userId + "_" + contentId)
	if status%2 == 0 {
		// UnFavorite -> Remove cache
		mRedisUSC.ZRem(PREFIX_REDIS_ZRANGE_USC_WATCHLATER+userId, contentId)
	} else {
		// Favorite -> add cache
		mRedisUSC.ZAdd(PREFIX_REDIS_ZRANGE_USC_WATCHLATER+userId, float64(curTime), contentId)
	}
	return status, nil
}

func PushWatchlaterToCache(userId string) error {
	/*
	   GetAll se bi timeout khi data qua lon
	   Can optimize them
	*/
	// Get data in DB
	session, db, err := GetCollection()
	if err != nil {
		return err
	}
	defer session.Close()

	var UscUserWatchlaterContentObjs []UscUserWatchlaterContentObjStruct
	var where = bson.M{
		"user_id": userId,
		"status":  1,
	}

	err = db.C(COLLECTION_USC_USER_WATCHLATER).Find(where).Sort("-updated_date").All(&UscUserWatchlaterContentObjs)
	if err != nil {
		Sentry_log(err)
		return err
	}

	for _, val := range UscUserWatchlaterContentObjs {
		mRedisUSC.ZAdd(PREFIX_REDIS_ZRANGE_USC_WATCHLATER+userId, float64(val.Updated_date), val.Content_id)
	}
	return nil
}

func UpdateWatchlaterDB(userId string, contentId string, status int64, entityType int) error {
	// Get data in DB
	session, db, err := GetCollection()
	if err != nil {
		return err
	}
	defer session.Close()

	status = status % 2

	var where = bson.M{
		"user_id":    userId,
		"content_id": contentId,
	}
	var UscUserWatchlaterContentObj UscUserWatchlaterContentObjStruct
	UscUserWatchlaterContentObj.User_id = userId
	UscUserWatchlaterContentObj.Content_id = contentId
	UscUserWatchlaterContentObj.Entity_type = entityType
	UscUserWatchlaterContentObj.Updated_date = time.Now().Unix()
	UscUserWatchlaterContentObj.Status = int(status)

	_, err = db.C(COLLECTION_USC_USER_WATCHLATER).Upsert(where, UscUserWatchlaterContentObj)
	if err != nil {
		return err
	}
	return nil
}

func GetWatchlaterTotalInCache(userId string) int64 {
	total := mRedisUSC.ZCountAll(PREFIX_REDIS_ZRANGE_USC_WATCHLATER + userId)
	return total
}

func CheckContentIsWatchLater(contentId string, userId string) bool {
	if userId == "" || strings.HasPrefix(userId, "anonymous_") == true {
		return false
	}

	// Check watch later
	zIndex := mRedisUSC.ZScore(PREFIX_REDIS_ZRANGE_USC_WATCHLATER+userId, contentId)
	if zIndex == 0 {
		return false
	}
	return true
}
