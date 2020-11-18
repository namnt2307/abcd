package live_event

import (
	. "cm-v5/serv/module"
	. "cm-v5/schema"
	"gopkg.in/mgo.v2/bson"
)

func GetLiveEventFinishedByCache(platform string, page, limit int, cacheActive bool) (LiveEventFinishedOutputObjectStruct, error) {
	var keyCache = KV_REDIS_LIVE_EVENT_FINISHED + "_" + platform + "_" + string(page) + "_" + string(limit)

	var LiveEventOutputs LiveEventFinishedOutputObjectStruct

	// Check Exists Data in cache
	if cacheActive == true {

		valueCache, err := mRedisKV.GetString(keyCache)
		if err == nil && valueCache != "" {
			err = json.Unmarshal([]byte(valueCache), &LiveEventOutputs)
			if err == nil {
				return LiveEventOutputs, nil
			}
		}
	}

	err, dataLiveEvent := GetLiveEventFinishedByMongoDB(platform, page, limit)
	if err != nil {
		Sentry_log(err)
		return LiveEventOutputs, err
	}

	dataByte, _ := json.Marshal(dataLiveEvent)
	err = json.Unmarshal(dataByte, &LiveEventOutputs.Items)
	if err != nil {
		Sentry_log(err)
		return LiveEventOutputs, err
	}

	platformInfo := Platform(platform)
	for k, val := range dataLiveEvent {
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
		LiveEventOutputs.Items[k].Images = ImagesMapping
		LiveEventOutputs.Items[k].Share_url = HandleShareUrlVieON(LiveEventOutputs.Items[k].Share_url)
		LiveEventOutputs.Items[k].Share_url_seo = HandleShareUrlVieON(LiveEventOutputs.Items[k].Share_url_seo)
	}

	LiveEventOutputs.Metadata.Page = page
	LiveEventOutputs.Metadata.Limit = limit
	LiveEventOutputs.Metadata.Total = GetTotalLiveEventFinished(platform, limit, cacheActive)

	dataByte, _ = json.Marshal(LiveEventOutputs)
	mRedisKV.SetString(keyCache, string(dataByte), TTL_KVCACHE)

	return LiveEventOutputs, nil
}
func GetTotalLiveEventFinished(platform string, limit int, cacheActive bool) int {
	var keyCache = PREFIX_REDIS_LIVE_EVENT_FINISHED_TOTAL + "_" + string(limit) + "_" + platform
	var total int

	if cacheActive {
		total, err := mRedis.GetInt(keyCache)
		if err == nil {
			return total
		}
	}

	// Connect DB
	session, db, err := GetCollection()
	if err != nil {
		Sentry_log(err)
		return 0
	}
	defer session.Close()

	where := bson.M{"is_live": 0}
	total, _ = db.C(COLLECTION_LIVE_EVENT).Find(where).Count()
	mRedis.SetInt(keyCache, total, TTL_KVCACHE)

	return total
}

func GetLiveEventFinishedByMongoDB(platform string, page, limit int) (error, []LiveEventObjectStruct) {

	var LiveEventObjects = make([]LiveEventObjectStruct, 0)
	// Connect DB
	session, db, err := GetCollection()
	if err != nil {
		Sentry_log(err)
		return err, LiveEventObjects
	}
	defer session.Close()

	where := bson.M{"is_live": 0}

	err = db.C(COLLECTION_LIVE_EVENT).Find(where).Sort("-str_to_time").Skip(page * limit).Limit(limit).All(&LiveEventObjects)
	if err != nil {
		Sentry_log(err)
		return err, LiveEventObjects
	}

	if len(LiveEventObjects) <= 0 {
		return err, LiveEventObjects
	}

	//Write cache
	// for _, liveEvent := range LiveEventObjects {
	// 	var keyCache = KV_REDIS_LIVE_EVENT_FINISHED_ID + "_" + fmt.Sprint(platform) + "_" + liveEvent.Id
	// 	dataByte, _ := json.Marshal(liveEvent)
	// 	mRedisKV.SetString(keyCache, string(dataByte), TTL_KVCACHE)
	// }
	return nil, LiveEventObjects
}
