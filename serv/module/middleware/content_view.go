package middleware

import (
	// . "cm-v5/serv/module"
	// . "cm-v5/schema"
	// vod "cm-v5/serv/module/vod"
	// tracking "cm-v5/serv/module/tracking"
	// live_event "cm-v5/serv/module/live_event"
	// "fmt"
	// "time"
	// "gopkg.in/mgo.v2/bson"
)

// var mRedisLog RedisLogModelStruct

func init() {
	
}

// func TrackContentViewInDay(dataInputProcessRequest  DataInputLogStruct , platform string) {
// 	// Get current day
// 	t := time.Now()
// 	var currentDate string = t.Format("2006-01-02")

// 	// Check content validate
// 	var TrackingWatchRequestData tracking.TrackingWatchRequestDataStruct
// 	err := json.Unmarshal([]byte(dataInputProcessRequest.Data), &TrackingWatchRequestData)
// 	if err == nil {
// 		switch TrackingWatchRequestData.Content_type {
// 		case VOD_TYPE_MOVIE , VOD_TYPE_EPISODE , VOD_TYPE_TRAILER:
// 			keyCache := fmt.Sprintf(PREFIX_LOG_VOD_VIEW , TrackingWatchRequestData.Content_id , GetMd5(dataInputProcessRequest.User_agent + dataInputProcessRequest.Ip + TrackingWatchRequestData.Content_id))
// 			if mRedisLog.Exists(keyCache) == 0 {
// 				mRedisLog.SetInt(keyCache , 1 , 2 * 3600)
// 				// Ghi vào DB 
// 				IncrContentViewInDate(currentDate , TrackingWatchRequestData , platform)
// 			}
// 		case VOD_TYPE_LIVESTREAM:
// 			keyCache := fmt.Sprintf(PREFIX_LOG_LIVE_STREAM_VIEW , TrackingWatchRequestData.Content_id , GetMd5(dataInputProcessRequest.User_agent + dataInputProcessRequest.Ip + TrackingWatchRequestData.Content_id))
// 			if mRedisLog.Exists(keyCache) == 0 {
// 				mRedisLog.SetInt(keyCache , 1 , 3600)
// 				// Ghi vào DB 
// 				IncrContentViewInDate(currentDate , TrackingWatchRequestData , platform)
// 			}
// 		case VOD_TYPE_LIVETV , VOD_TYPE_EPG:
// 			keyCache := fmt.Sprintf(PREFIX_LOG_LIVE_TV_VIEW , TrackingWatchRequestData.Content_id , TrackingWatchRequestData.Content_name , GetMd5(dataInputProcessRequest.User_agent + dataInputProcessRequest.Ip + TrackingWatchRequestData.Content_id))
// 			if mRedisLog.Exists(keyCache) == 0 {
// 				mRedisLog.SetInt(keyCache , 1 , 2 * 3600)
// 				// Ghi vào DB 
// 				IncrContentViewInDate(currentDate , TrackingWatchRequestData , platform)
// 			}
// 		}
// 	}
// }


// func IncrContentViewInDate(currentDate string , TrackingWatchRequestData tracking.TrackingWatchRequestDataStruct , platform string) error {
// 	// Get data in DB
// 	session, db, err := GetCollection()
// 	if err != nil {
// 		return err
// 	}
// 	defer session.Close()

// 	var timeView int64 = int64(len(TrackingWatchRequestData.Data)) * int64(tracking.TRACKING_WATCH_RECORD_TIME)
// 	var platformInfo = Platform(platform)
// 	var where = bson.M{}

// 	switch TrackingWatchRequestData.Content_type {
// 	case VOD_TYPE_LIVETV , VOD_TYPE_EPG:
// 		where = bson.M{
// 			"date":         currentDate,
// 			"group_id": 	TrackingWatchRequestData.Content_id,
// 			"entity_name": 	TrackingWatchRequestData.Content_name,
// 		}
// 	default:
// 		where = bson.M{
// 			"date":         currentDate,
// 			"entity_id": 	TrackingWatchRequestData.Content_id,
// 		}
// 	}

// 	var EntityLogObject	EntityLogObjectStruct
// 	cEntityLog := db.C(COLLECTION_ENTITY_LOG)
// 	// Check exist
// 	err = cEntityLog.Find(where).One(&EntityLogObject)
// 	if err == nil {
// 		var updateData = bson.M{"$inc": bson.M{"num_view": 1 , "web_view": 1 , "time_view": timeView}}
// 		switch platformInfo.Type {
// 		case "smarttv":
// 			updateData = bson.M{"$inc": bson.M{"num_view": 1 , "tv_view": 1 , "time_view": timeView}}
// 		case "app":
// 			updateData = bson.M{"$inc": bson.M{"num_view": 1 , "app_view": 1 , "time_view": timeView}}
// 		}
// 		err = cEntityLog.Update(where, updateData)	
// 		if err != nil {
// 			Sentry_log(err)
// 			return err
// 		}
// 	} else {
// 		switch TrackingWatchRequestData.Content_type {
// 		case VOD_TYPE_MOVIE , VOD_TYPE_EPISODE , VOD_TYPE_TRAILER:
// 			// Get VOD Detail
// 			dataVod, err := vod.GetVodDetail(TrackingWatchRequestData.Content_id, "", 1, true)
// 			if err != nil {
// 				return err
// 			}
// 			EntityLogObject.Entity_id 		= TrackingWatchRequestData.Content_id
// 			EntityLogObject.Entity_name 	= dataVod.Title
// 			if TrackingWatchRequestData.Content_type == VOD_TYPE_EPISODE {
// 				EntityLogObject.Group_id 	= dataVod.Group_id
// 				EntityLogObject.Group_name 	= dataVod.Movie.Title
// 			}

// 		case VOD_TYPE_LIVESTREAM:
// 			// Get Event Detail
// 			dataLiveEvent, err := live_event.GetDetailByID(TrackingWatchRequestData.Content_id, platform, true)
// 			if err != nil {
// 				return err
// 			}
// 			EntityLogObject.Entity_id 	= TrackingWatchRequestData.Content_id
// 			EntityLogObject.Entity_name = dataLiveEvent.Title

// 		case VOD_TYPE_LIVETV , VOD_TYPE_EPG:
// 			EntityLogObject.Group_id 	= TrackingWatchRequestData.Content_id
// 			EntityLogObject.Entity_name = TrackingWatchRequestData.Content_name

// 		}

// 		EntityLogObject.Entity_type = TrackingWatchRequestData.Content_type
// 		EntityLogObject.Num_view = 1
// 		EntityLogObject.Time_view = timeView
// 		EntityLogObject.Date = currentDate
// 		switch platformInfo.Type {
// 		case "smarttv":
// 			EntityLogObject.Tv_view = 1
// 		case "app":
// 			EntityLogObject.App_view = 1
// 		default:
// 			EntityLogObject.Web_view = 1
// 		}

// 		_ , err = cEntityLog.Upsert(where , EntityLogObject)	
// 		if err != nil {
// 			Sentry_log(err)
// 			return err
// 		}
// 	}
// 	return nil
// }

