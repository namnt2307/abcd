package tracking

import (
	"errors"
	"fmt"
	"strconv"
	"strings"
	"time"

	. "cm-v5/schema"
	. "cm-v5/serv/module"
	fingering "cm-v5/serv/module/fingering"
	Kplus "cm-v5/serv/module/kplus"
	livetv_v3 "cm-v5/serv/module/livetv_v3"
	recommendation "cm-v5/serv/module/recommendation"
	vod "cm-v5/serv/module/vod"
	watchlater "cm-v5/serv/module/watchlater"

	"github.com/gin-gonic/gin"
	"github.com/gosimple/slug"
)

func TrackingWatch(c *gin.Context, user_id string, trackingWatchRequestData TrackingWatchRequestDataStruct) (TrackingWatchResultStruct, error) {
	var TrackingWatchResult TrackingWatchResultStruct
	TrackingWatchResult.Data.Next_log = TRACKING_WATCH_NEXT_LOG
	TrackingWatchResult.Data.Record_time = TRACKING_WATCH_RECORD_TIME
	TrackingWatchResult.Success = true

	platform := c.DefaultQuery("platform", "web")

	if len(trackingWatchRequestData.Data) <= 0 {
		return TrackingWatchResult, nil
	}

	var content_id string = trackingWatchRequestData.Content_id
	switch trackingWatchRequestData.Content_type {
	case VOD_TYPE_LIVESTREAM, VOD_TYPE_EPG:
		return TrackingWatchResult, nil
	case VOD_TYPE_LIVETV:
		liveObj, err := livetv_v3.GetLiveTVByListID([]string{content_id}, 0, true)
		if err != nil || len(liveObj) <= 0 {
			TrackingWatchResult.Message = "not found content"
			TrackingWatchResult.Success = false
			return TrackingWatchResult, nil
		}

		if trackingWatchRequestData.Usi != "" && user_id != "" {
			// Gia hạn key stream từng USI
			fingering.KeepUniqueStreamingID(trackingWatchRequestData.Usi, user_id)
		}

		if Kplus.CheckIdLivetvIsKplus(trackingWatchRequestData.Content_id) {
			// Gia hạn key active stream của user
			fingering.KeepStreamingByUserID(user_id)
		}

		dataByte, _ := json.Marshal(liveObj[0])
		err = json.Unmarshal(dataByte, &TrackingWatchResult.Content_info)
		if err != nil {
			return TrackingWatchResult, err
		}

		return TrackingWatchResult, nil
	}

	var duration, progress, timestamp, percentage int64
	var endData = trackingWatchRequestData.Data[len(trackingWatchRequestData.Data)-1]

	// Check content validate
	vodObj, err := vod.GetVodDetail(content_id, 0, true)
	if err != nil {
		TrackingWatchResult.Message = "not found content"
		TrackingWatchResult.Success = false
		return TrackingWatchResult, nil
	}

	//sync info
	dataByte, _ := json.Marshal(vodObj)
	err = json.Unmarshal(dataByte, &TrackingWatchResult.Content_info)
	if err != nil {
		return TrackingWatchResult, err
	}

	duration = endData.Duration
	progress = endData.Progress
	timestamp = time.Now().Unix()

	if duration == 0 || progress < 0 {
		TrackingWatchResult.Message = "duration / progress not valalidate"
		TrackingWatchResult.Success = false
		return TrackingWatchResult, nil
	}

	if user_id == "" {
		return TrackingWatchResult, nil
	}
	if strings.HasPrefix(user_id, "anonymous_") == true {
		return TrackingWatchResult, nil
	}

	if user_id == "" {
		return TrackingWatchResult, nil
	}
	if strings.HasPrefix(user_id, "anonymous_") == true {
		return TrackingWatchResult, nil
	}

	if trackingWatchRequestData.Usi != "" && user_id != "" {
		// Gia hạn key stream
		fingering.KeepUniqueStreamingID(trackingWatchRequestData.Usi, user_id)
	}

	percentage = (progress * 100) / duration
	if percentage >= 95 {
		// Remove cache
		mRedisUSC.ZRem(PREFIX_REDIS_ZRANGE_USC_WATCHING+user_id, content_id)
		mRedisUSC.HDel(PREFIX_REDIS_HASH_USC_WATCHING+user_id, content_id)

		//check episode is end clear history
		if vodObj.Type == VOD_TYPE_EPISODE {

			//async
			go func(groupID, platform string, episode int) {
				seasonObj, _ := vod.GetVodDetail(groupID, 0, true)
				currentEpisode, _ := StringToInt(seasonObj.Current_episode)

				//nếu phim này đã ra tập cuối
				if episode == seasonObj.Episode && seasonObj.Current_episode != "" && currentEpisode == seasonObj.Episode {
					platformInfo := Platform(platform)
					//get list episode
					//delete all history in season
					listIdEpisode := vod.GetListIdEpisodeByCache(groupID, platformInfo.Id, 0, 100)
					for _, ID := range listIdEpisode {
						mRedisUSC.ZRem(PREFIX_REDIS_ZRANGE_USC_WATCHING+user_id, ID)
						mRedisUSC.HDel(PREFIX_REDIS_HASH_USC_WATCHING+user_id, ID)
					}
					mRedisUSC.HDel(PREFIX_REDIS_HASH_USC_WATCHING_SESSION+user_id, groupID)

				}
			}(vodObj.Group_id, platform, vodObj.Episode)

		}
		return TrackingWatchResult, nil
	}

	if strings.ToUpper(endData.Action) != "PLAY" {
		// Detail progress per connent by user
		var valRedis string = content_id + "_" + fmt.Sprint(progress)
		mRedisUSC.HSet(PREFIX_REDIS_HASH_USC_WATCHING+user_id, content_id, valRedis)

		// Detail progress per session by user
		if vodObj.Type == VOD_TYPE_EPISODE {

			// Remove eps old in list watching by user
			dataRedis, err := mRedisUSC.HMGet(PREFIX_REDIS_HASH_USC_WATCHING_SESSION+user_id, []string{vodObj.Group_id})
			if err == nil && len(dataRedis) > 0 {
				if str, ok := dataRedis[0].(string); ok {
					mRedisUSC.ZRem(PREFIX_REDIS_ZRANGE_USC_WATCHING+user_id, str)
				}
			}

			mRedisUSC.HSet(PREFIX_REDIS_HASH_USC_WATCHING_SESSION+user_id, vodObj.Group_id, content_id)
		}

		// Danh sách VOD dang xem theo thu tu (moi -> cu)
		mRedisUSC.ZAdd(PREFIX_REDIS_ZRANGE_USC_WATCHING+user_id, float64(timestamp), content_id)
	}

	//Push Kafka
	var TrackingData TrackingDataStruct
	TrackingUserData, _ := recommendation.GetTracking(c, content_id, TrackingData)
	TrackingUserData.Duration = duration
	TrackingUserData.Percentage = percentage
	recommendation.PushDataToKafka(WATCH, TrackingUserData)

	return TrackingWatchResult, nil
}

func ListWatching(userId string, platform string, page int, limit int) (WatchingOutputObjectStruct, error) {
	var platformDetail = Platform(platform)
	var WatchingOutputObj WatchingOutputObjectStruct
	WatchingOutputObj.Metadata.Limit = limit
	WatchingOutputObj.Metadata.Page = page
	WatchingOutputObj.Items = make([]ItemWatchingOutputObjectStruct, 0)

	//Config off push data kafka
	if ENABLE_WATCHMORE == "disable" {
		return WatchingOutputObj, nil
	}

	contentIds, err := GetListIDWatchingByUserID(userId, page, limit)
	if len(contentIds) > 0 && err == nil {
		// Get List Progress
		listProgress, err := GetProgressByListID(contentIds, userId)
		if len(listProgress) <= 0 || err != nil {
			return WatchingOutputObj, nil
		}

		// Get List VOD Detail
		var VODDataObjects []VODDataObjectStruct
		VODDataObjects, err = vod.GetVODByListID(contentIds, platformDetail.Id, 1, true)
		if err != nil {
			return WatchingOutputObj, nil
		}

		// Get List Group Id of Episode
		var ListGroupID []string
		for _, vodData := range VODDataObjects {
			if vodData.Type == VOD_TYPE_EPISODE {
				ListGroupID = append(ListGroupID, vodData.Group_id)
			}
		}

		// get data season from groupID
		var VODDataObjectsGetByGroupIDs []VODDataObjectStruct
		if len(ListGroupID) > 0 {
			VODDataObjectsGetByGroupIDs, _ = vod.GetVODByListID(ListGroupID, platformDetail.Id, 1, true)
		}

		// Format data output
		dataByte, _ := json.Marshal(VODDataObjects)
		err = json.Unmarshal(dataByte, &WatchingOutputObj.Items)
		if err != nil {
			return WatchingOutputObj, err
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
			WatchingOutputObj.Items[k].Images = ImagesMapping
			WatchingOutputObj.Items[k].Is_watchlater = watchlater.CheckContentIsWatchLater(vodData.Id, userId)

			if progress, ok := listProgress[vodData.Id]; ok {
				WatchingOutputObj.Items[k].Progress = progress
			}
		}

		// Get total
		WatchingOutputObj.Metadata.Total = GetWatchingTotalInCache(userId)
		WatchingOutputObj.Tracking_data = GetTrackingData(platform)
	}

	return WatchingOutputObj, nil
}

func GetListIDWatchingByUserID(userId string, page int, limit int) ([]string, error) {
	// Check Exists Data in cache
	if mRedisUSC.Exists(PREFIX_REDIS_ZRANGE_USC_WATCHING+userId) == 0 {
		return make([]string, 0), nil
	}
	var start = page * limit
	var stop = (start + limit) - 1
	contentIds := mRedisUSC.ZRevRange(PREFIX_REDIS_ZRANGE_USC_WATCHING+userId, int64(start), int64(stop))
	return contentIds, nil
}

func GetWatchingTotalInCache(userId string) int {
	total := mRedisUSC.ZCountAll(PREFIX_REDIS_ZRANGE_USC_WATCHING + userId)
	return int(total)
}

func GetTrackingData(platform string) TrackingDataStruct {
	trackingType := platform + "_VIDEO_DETAIL"
	trackingType = slug.Make(trackingType)
	return recommendation.GetRandomDefaultTrackingData(strings.ToUpper(trackingType))
}

func GetProgressByListID(listId []string, userId string) (map[string]int64, error) {
	var listProgress = make(map[string]int64, 0)
	if len(listId) <= 0 {
		return listProgress, errors.New("GetProgressByListID: Empty data")
	}

	if userId == "" || strings.HasPrefix(userId, "anonymous_") == true {
		return listProgress, nil
	}

	dataRedis, err := mRedisUSC.HMGet(PREFIX_REDIS_HASH_USC_WATCHING+userId, listId)
	if err != nil || len(dataRedis) <= 0 {
		return listProgress, err
	}

	for _, val := range dataRedis {
		if str, ok := val.(string); ok {
			strs := strings.Split(str, "_")
			if len(strs) > 1 {
				progress, err := strconv.ParseInt(strs[1], 10, 64)
				if err != nil {
					continue
				}
				listProgress[strs[0]] = progress
			}
		}
	}

	return listProgress, nil
}

func GetIdDefaultEpsBySessionID(ssId string, userId string) string {
	dataRedis, err := mRedisUSC.HGet(PREFIX_REDIS_HASH_USC_WATCHING_SESSION+userId, ssId)
	if err == nil {
		return dataRedis
	}
	return ""
}

func RemoveWatching(userId string) (DataResponseStruct, error) {
	var DataResponse DataResponseStruct
	DataResponse.Success = true

	mRedisUSC.Del(PREFIX_REDIS_HASH_USC_WATCHING + userId)
	mRedisUSC.Del(PREFIX_REDIS_ZRANGE_USC_WATCHING + userId)
	mRedisUSC.Del(PREFIX_REDIS_HASH_USC_WATCHING_SESSION + userId)
	return DataResponse, nil
}
