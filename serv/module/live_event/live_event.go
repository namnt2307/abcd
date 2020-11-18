package live_event

import (
	"errors"

	. "cm-v5/serv/module"
	. "cm-v5/schema"
	"gopkg.in/mgo.v2/bson"
)

func GetDetailByID(liveEventID string, platform string, cacheActive bool) (LiveEventOutputObjectStruct, error) {
	var liveEventIds = []string{liveEventID}
	var LiveEventOutput LiveEventOutputObjectStruct
	var keyCache = KV_REDIS_LIVE_EVENT_ID + "_" + platform + "_" + liveEventID
	if cacheActive {
		valueCache, err := mRedisKV.GetString(keyCache)
		if err == nil && valueCache != "" {
			err = json.Unmarshal([]byte(valueCache), &LiveEventOutput)
			if err == nil {
				return LiveEventOutput, nil
			}
		}
	}

	dataLiveEvents, err := GetLiveEventByListID(liveEventIds, platform, 0, 30, cacheActive)
	if len(dataLiveEvents) <= 0 || err != nil {
		return LiveEventOutput, errors.New("Live Event Empty")
	}

	// Write Redis
	dataByte, _ := json.Marshal(dataLiveEvents[0])
	mRedisKV.SetString(keyCache, string(dataByte), TTL_KVCACHE)
	return dataLiveEvents[0], nil

}

func GetIdBySLug(liveEventSlug string, cacheActive bool) (string, error) {
	var liveEventID string
	var keyCache = "live_event_slug_" + liveEventSlug
	if cacheActive {
		dataCache, err := mRedis.GetString(keyCache)
		if err == nil && dataCache != "" {
			err = json.Unmarshal([]byte(dataCache), &liveEventID)
			if err == nil {
				return liveEventID, nil
			}
		}
	}

	// Connect DB
	session, db, err := GetCollection()
	if err != nil {
		Sentry_log(err)
		return liveEventID, err
	}
	defer session.Close()

	var where = bson.M{
		"slug": liveEventSlug,
	}

	var LiveEventObject LiveEventObjectStruct
	err = db.C(COLLECTION_LIVE_EVENT).Find(where).One(&LiveEventObject)
	if err != nil {
		Sentry_log(err)
		return liveEventID, err
	}

	// Write cache
	dataByte, _ := json.Marshal(LiveEventObject.Id)
	mRedis.SetString(keyCache, string(dataByte), TTL_REDIS_LV1)

	return LiveEventObject.Id, nil
}

func GetAllLiveEventByCache(platform string, page, limit int, cacheActive bool) ([]LiveEventOutputObjectStruct, error) {
	var LiveEventOutputs = make([]LiveEventOutputObjectStruct, 0)

	// Check Exists Data in cache
	if mRedis.Exists(PREFIX_REDIS_HASH_LIVE_EVENT) == 0 || cacheActive == false {

		// No have Data in cache
		// Push data from DB
		err := GetAllLiveEventByMongoDB()
		if err != nil {
			Sentry_log(err)
			return LiveEventOutputs, err
		}
	}

	start := page * limit
	stop := (start + limit) - 1
	listLiveEventIds, err := mRedis.ZRange(PREFIX_REDIS_ZRANGE_LIVE_EVENT, int64(start), int64(stop))
	if err != nil || len(listLiveEventIds) <= 0 {
		Sentry_log(err)
		return LiveEventOutputs, err
	}

	dataLiveEvent, err := GetLiveEventByListID(listLiveEventIds, platform, page, limit, true)
	if err != nil || len(dataLiveEvent) <= 0 {
		Sentry_log(err)
		return LiveEventOutputs, err
	}

	return dataLiveEvent, nil
}

func GetLiveEventByListID(liveEventIds []string, platform string, page, limit int, cacheActive bool) ([]LiveEventOutputObjectStruct, error) {
	var LiveEventOutputs = make([]LiveEventOutputObjectStruct, 0)
	if len(liveEventIds) < 0 {
		return LiveEventOutputs, errors.New("GetLiveEventByListID: Empty data ")
	}

	var LiveEventObjects = make([]LiveEventObjectStruct, 0)

	// Check Exists Data in cache
	if cacheActive {
		dataRedis, err := mRedis.HMGet(PREFIX_REDIS_HASH_LIVE_EVENT, liveEventIds)
		if err == nil && len(dataRedis) > 0 {
			for _, val := range dataRedis {
				if str, ok := val.(string); ok {
					var LiveEventObject LiveEventObjectStruct
					err = json.Unmarshal([]byte(str), &LiveEventObject)
					if err != nil {
						continue
					}
					LiveEventObjects = append(LiveEventObjects, LiveEventObject)
				}
			}
		}
	}

	platformInfo := Platform(platform)
	if len(LiveEventObjects) != len(liveEventIds) {

		// Connect DB
		session, db, err := GetCollection()
		if err != nil {
			Sentry_log(err)
			return LiveEventOutputs, err
		}
		defer session.Close()

		var where = bson.M{
			"id": bson.M{"$in": liveEventIds},
		}
		err = db.C(COLLECTION_LIVE_EVENT).Find(where).Sort("-str_to_time").Skip(page * limit).Limit(limit).All(&LiveEventObjects)
		if err != nil {
			Sentry_log(err)
			return LiveEventOutputs, err
		}

		if len(LiveEventObjects) <= 0 {
			for _, id := range liveEventIds {
				mRedis.HDel(PREFIX_REDIS_HASH_LIVE_EVENT, id)
			}
			return LiveEventOutputs, err
		}

		dataByte, _ := json.Marshal(LiveEventObjects)
		err = json.Unmarshal(dataByte, &LiveEventOutputs)
		if err != nil {
			Sentry_log(err)
			return LiveEventOutputs, err
		}

		//Write cache
		for _, liveEvent := range LiveEventObjects {
			dataByte, _ := json.Marshal(liveEvent)
			mRedis.HSet(PREFIX_REDIS_HASH_LIVE_EVENT, liveEvent.Id, string(dataByte))
			mRedis.ZAdd(PREFIX_REDIS_ZRANGE_LIVE_EVENT, float64(liveEvent.Str_to_time), liveEvent.Id)
		}
	}

	dataByte, _ := json.Marshal(LiveEventObjects)
	err := json.Unmarshal(dataByte, &LiveEventOutputs)
	if err != nil {
		Sentry_log(err)
		return LiveEventOutputs, err
	}

	for k, val := range LiveEventObjects {
		var ImagesMapping ImagesOutputObjectStruct

		// switch images follow platform
		switch platformInfo.Type {
		case "web":
			ImagesMapping.Vod_thumb_big = BuildImage(val.Images.Web.Vod_thumb_big)
			ImagesMapping.Thumbnail = BuildImage(val.Images.Web.Thumbnail)
		case "smarttv":
			ImagesMapping.Vod_thumb_big = BuildImage(val.Images.Smarttv.Vod_thumb_big)
			ImagesMapping.Thumbnail = BuildImage(val.Images.Smarttv.Thumbnail)
		case "app":
			ImagesMapping.Vod_thumb_big = BuildImage(val.Images.App.Vod_thumb_big)
			ImagesMapping.Thumbnail = BuildImage(val.Images.App.Thumbnail)
		}
		ImagesMapping = MappingImagesV4(platformInfo.Type, ImagesMapping, val.Images, true)
		LiveEventOutputs[k].Images = ImagesMapping
		LiveEventOutputs[k].Share_url = HandleShareUrlVieON(LiveEventOutputs[k].Share_url)
		LiveEventOutputs[k].Share_url_seo = HandleShareUrlVieON(LiveEventOutputs[k].Share_url_seo)

		//Lay thong tin goi
		LiveEventOutputs[k] = GetPackageLiveEvent(LiveEventOutputs[k])
	}

	return LiveEventOutputs, nil
}

func GetAllLiveEventByMongoDB() error {

	//Clear cache
	mRedis.Del(PREFIX_REDIS_HASH_LIVE_EVENT)
	mRedis.Del(PREFIX_REDIS_ZRANGE_LIVE_EVENT)

	// Connect DB
	session, db, err := GetCollection()
	if err != nil {
		Sentry_log(err)
		return err
	}
	defer session.Close()

	var LiveEventObjects []LiveEventObjectStruct
	where := bson.M{"is_live": 1}
	err = db.C(COLLECTION_LIVE_EVENT).Find(where).All(&LiveEventObjects)
	if err != nil {
		Sentry_log(err)
		return err
	}

	if len(LiveEventObjects) <= 0 {
		return err
	}

	//Write cache
	for _, liveEvent := range LiveEventObjects {
		dataByte, _ := json.Marshal(liveEvent)
		mRedis.HSet(PREFIX_REDIS_HASH_LIVE_EVENT, liveEvent.Id, string(dataByte))
		mRedis.ZAdd(PREFIX_REDIS_ZRANGE_LIVE_EVENT, float64(liveEvent.Str_to_time), liveEvent.Id)
	}
	return nil
}
