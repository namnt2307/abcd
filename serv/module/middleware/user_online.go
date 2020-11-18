package middleware

import (
	// . "cm-v5/serv/module"
	// . "cm-v5/schema"
	// "fmt"
	// "log"
	// "time"
	// "gopkg.in/mgo.v2/bson"
)

// var UscOnline_NumToWriteDB int64 = 1

// func TrackUserOnline(user_id string) {
// 	return
// 	log.Println("TrackUserOnline: " , user_id)


// 	// Get current time
// 	t := time.Now()

// 	t = t.Add(time.Second * 3600 * (26))

// 	var currentDate string = t.Format("2006-01-02")
// 	var currentDayStr string = t.Format("Monday")
// 	var currentHour string = t.Format("03PM")
// 	var currentMinute string = t.Format("04")

// 	keyCacheCurrentHour := fmt.Sprintf(PREFIX_LOG_USC_ONLINE_PER_HOUS , currentHour ,  user_id) 

// 	log.Println(currentDate)
// 	log.Println(currentHour)
// 	log.Println(currentMinute)
// 	log.Println(keyCacheCurrentHour)

// 	// Khi user có activity (gọi API) ở khung giờ nào, sẽ tạo key tương ứng khung giờ đó (check exist trước khi tạo)
// 	if mRedisLog.Exists(keyCacheCurrentHour) == 0 {
// 		mRedisLog.SetInt(keyCacheCurrentHour , 1 , 3600)

// 		// Tang total online per hour
// 		keyCacheTotalCurrentHour := fmt.Sprintf(PREFIX_LOG_TOTAL_ONLINE_PER_HOUS , currentHour) 
// 		totalOnlineCurrentHour := mRedisLog.Incr(keyCacheTotalCurrentHour)

// 		if totalOnlineCurrentHour % UscOnline_NumToWriteDB == 0 {
// 			// Ghi vào DB 
// 			DBIncrUSCOnlinePerHour(currentDate , currentDayStr , currentHour , totalOnlineCurrentHour)
// 		}
		
// 	}

// 	// Nếu thời gian phát sinh activity (gọi API) >= "Next hour" - "Activity Expired", tạo key cho khung giờ tiếp theo


// 	// USCOnlineLogObjectStruct
	
// }


// func DBIncrUSCOnlinePerHour(date, day_str, hour string , total int64) error {
// 	// Get data in DB
// 	session, db, err := GetCollection()
// 	if err != nil {
// 		return err
// 	}
// 	defer session.Close()
// 	var where = bson.M{
// 		"date":         date,
// 	}

// 	var USCOnlineLogObject	USCOnlineLogObjectStruct
// 	cOnlineLog := db.C(COLLECTION_USC_ONLINE_LOG)
// 	// Check exist
// 	err = cOnlineLog.Find(where).One(&USCOnlineLogObject)
// 	if err == nil {
// 		var updateData = bson.M{"$set": bson.M{"detail_hour." + hour: total }}
// 		err = cOnlineLog.Update(where, updateData)	
// 		if err != nil {
// 			Sentry_log(err)
// 			return err
// 		}
// 	} else {
// 		var detailHour = make(map[string]int64 , 0)
// 		detailHour[hour] = total

// 		USCOnlineLogObject.Date = date
// 		USCOnlineLogObject.Day_str = day_str
// 		USCOnlineLogObject.Detail_hour = detailHour
// 		_ , err = cOnlineLog.Upsert(where , USCOnlineLogObject)	
// 		if err != nil {
// 			Sentry_log(err)
// 			return err
// 		}
// 	}
// 	return nil
// }

