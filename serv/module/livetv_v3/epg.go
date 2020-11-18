package livetv_v3

import (
	// "encoding/json"

	"fmt"
	"time"

	. "cm-v5/serv/module"
	. "cm-v5/schema"
	"gopkg.in/mgo.v2/bson"
	// seo "cm-v5/serv/module/seo"
)

func GetEPGCurrent(livetv_id string) (EpgObjectOutputStruct, error) {
	var EpgObject EpgObjectOutputStruct
	start, end := GetStartEndTimeFromTimestamp(0)

	dataEpgTemp, err := GetEpgByLivetvID(livetv_id, start, end, true)
	if err != nil || len(dataEpgTemp) <= 0 {
		Sentry_log(err)
		return EpgObject, err
	}

	var currentTime = int(time.Now().Unix())
	for _, val := range dataEpgTemp {
		if val.Time_start <= currentTime && currentTime < val.Time_end {
			EpgObject = val
			break
		}
	}

	return EpgObject, nil
}

func GetEpgByLivetvID(livetv_id string, start, end int, cache bool) ([]EpgObjectOutputStruct, error) {
	var EpgObjects = make([]EpgObjectOutputStruct, 0)
	var keyCache = KV_REDIS_EPG_BY_LIVE_TV_ID + "_" + livetv_id + "_" + fmt.Sprint(start) + "_" + fmt.Sprint(end)

	if cache {
		valueCache, err := mRedisKV.GetString(keyCache)
		if err == nil && valueCache != "" {
			err = json.Unmarshal([]byte(valueCache), &EpgObjects)
			if err == nil {
				EpgObjects = ProcessIsCatchUpEpg(livetv_id, EpgObjects)
				return EpgObjects, nil
			}
		}
	}


	// Connect DB
	session, db, err := GetCollection()
	if err != nil {
		Sentry_log(err)
		return EpgObjects, err
	}
	defer session.Close()

	var where = bson.M{
		"group_id": livetv_id,
		"type":     7,
		"status":   1,
		// "time_start": bson.M{"$gte": start},
		// "time_end":   bson.M{"$lte": end},
		"$or": []bson.M{
			bson.M{"$and": []bson.M{
				bson.M{"time_start": bson.M{"$gte": start}},
				bson.M{"time_start": bson.M{"$lte": end}},
			}},
			bson.M{
				"time_start": bson.M{"$lte": start},
				"time_end":   bson.M{"$gte": start},
			},
		},
	}

	err = db.C(COLLECTION_LIVE_TV).Find(where).Sort("time_start").All(&EpgObjects)
	if err != nil {
		Sentry_log(err)
		return EpgObjects, err
	}

	// for k, val := range EpgObjects {
	// 	EpgObjects[k].Seo = seo.FormatSeoEPG("", val.Id, val.Title)
	// }

	// Write Redis
	dataByte, _ := json.Marshal(EpgObjects)
	mRedisKV.SetString(keyCache, string(dataByte), 60*60*2)

	EpgObjects = ProcessIsCatchUpEpg(livetv_id, EpgObjects)

	return EpgObjects, nil
}

func GetDetailEpgBySlug(slug_epg string, cache bool) (EpgObjectOutputStruct, error) {
	var EpgObject EpgObjectOutputStruct
	start, end := GetStartEndTimeFromTimestamp(0)
	start = start - 3600*24*3 //3day
	var keyCache = KV_REDIS_DETAIL_EPG + "_" + slug_epg

	if cache {
		valueCache, err := mRedisKV.GetString(keyCache)
		if err == nil && valueCache != "" {
			err = json.Unmarshal([]byte(valueCache), &EpgObject)
			if err == nil {
				return EpgObject, nil
			}
		}
	}

	// Connect DB
	session, db, err := GetCollection()
	if err != nil {
		Sentry_log(err)
		return EpgObject, err
	}
	defer session.Close()

	// Get Detail EPG
	var where = bson.M{
		"seo.url":    slug_epg,
		"type":       7,
		"status":     1,
		"time_start": bson.M{"$gte": start},
		"time_end":   bson.M{"$lte": end},
	}
	err = db.C(COLLECTION_LIVE_TV).Find(where).Sort("time_start").One(&EpgObject)
	if err != nil && err.Error() != "not found" {
		return EpgObject, err
	}

	// Write Redis
	dataByte, _ := json.Marshal(EpgObject)
	mRedisKV.SetString(keyCache, string(dataByte), 60*60*2)

	return EpgObject, nil
}

func ProcessIsCatchUpEpg(livetv_id string, listEpg []EpgObjectOutputStruct) []EpgObjectOutputStruct {
	if len(listEpg) == 0 {
		return listEpg
	}

	//get info livetv
	infoLivetv, err := GetLiveTVByListID([]string{livetv_id}, 0, true)

	if err == nil && len(infoLivetv) > 0 {
		livetvCatchUp := infoLivetv[0].Is_catch_up
		maxTimeCatchUp := infoLivetv[0].Max_time_catchup
		currentTime := int(time.Now().Unix())

		for i := 0; i < len(listEpg); i++ {
			rangeTime := currentTime - listEpg[i].Time_start

			//Nếu chương trình chưa được phát sóng thì không cần check
			//Vì thứ tự đã sắp xếp nên có thể break luôn
			if rangeTime < 0 {
				break
			}

			//nếu livetv set is_catch_up = true thì trả true cho tất cả epg nằm trong khoảng Max_time_catchup (hours) tính đến thời điểm hiện tại
			if livetvCatchUp && rangeTime/3600 < maxTimeCatchUp {
				listEpg[i].Is_catch_up = true
				continue
			}

			//nếu livetv set is_catch_up = false thì trả false cho tất cả epg
			listEpg[i].Is_catch_up = false
		}

	}
	return listEpg
}
