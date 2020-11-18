package content

import (
	"strconv"

	. "cm-v5/serv/module"
	. "cm-v5/schema"
	vod "cm-v5/serv/module/vod"
	"gopkg.in/mgo.v2/bson"
)

func GetRelatedVideos(groupId string, page int, limit int, platforms string, cacheActive bool) (RelatedVideosObjOutputStruct, error) {
	platform := Platform(platforms)
	var RelatedVideosObjOutput RelatedVideosObjOutputStruct
	var keyCache = "KV_related_videos_" + groupId + "_" + strconv.Itoa(platform.Id) + "_" + strconv.Itoa(page) + "_" + strconv.Itoa(limit)

	if cacheActive {
		valueCache, err := mRedisKV.GetString(keyCache)
		if err == nil && valueCache != "" {
			err = json.Unmarshal([]byte(valueCache), &RelatedVideosObjOutput)
			if err == nil {
				return RelatedVideosObjOutput, nil
			}
		}
	}

	var total int = 0

	// Lay danh sach Video Lien Quan
	listVodID, err := GetListIDRelatedVideoByZrange(groupId, page, limit, cacheActive)
	if len(listVodID) > 0 && err == nil {

		var VODDataObjects []VODDataObjectStruct
		VODDataObjects, err = vod.GetVODByListID(listVodID, platform.Id, 1, cacheActive)
		if err != nil {
			return RelatedVideosObjOutput, err
		}

		dataByte, _ := json.Marshal(VODDataObjects)
		err = json.Unmarshal(dataByte, &RelatedVideosObjOutput.Items)
		if err != nil {
			return RelatedVideosObjOutput, err
		}

		//Switch images platform
		for k, vodData := range VODDataObjects {
			var ImagesMapping ImagesOutputObjectStruct
			switch platform.Type {
			case "web":
				ImagesMapping.Vod_thumb_big = BuildImage(vodData.Images.Web.Vod_thumb)
				ImagesMapping.Home_vod_hot = BuildImage(vodData.Images.Web.Home_vod_hot)
				ImagesMapping.Vod_thumb = BuildImage(vodData.Images.Web.Vod_thumb)
				ImagesMapping.Thumbnail = BuildImage(vodData.Images.Web.Vod_thumb)
			case "smarttv":
				ImagesMapping.Thumbnail = BuildImage(vodData.Images.Smarttv.Thumbnail)
			case "app":
				ImagesMapping.Thumbnail = BuildImage(vodData.Images.App.Thumbnail)
				ImagesMapping.Vod_thumb = BuildImage(vodData.Images.App.Vod_thumb)
			}

			ImagesMapping = MappingImagesV4(platform.Type, ImagesMapping, vodData.Images, true)
			RelatedVideosObjOutput.Items[k].Images = ImagesMapping

			RelatedVideosObjOutput.Items[k].Seo = vodData.Seo
		}

		total, _ = GetTotalRelatedVideos(groupId, platforms, cacheActive)
	}

	//Pagination
	RelatedVideosObjOutput.Metadata.Total = total
	RelatedVideosObjOutput.Metadata.Limit = limit
	RelatedVideosObjOutput.Metadata.Page = page

	// Write Redis
	dataByte, _ := json.Marshal(RelatedVideosObjOutput)
	mRedisKV.SetString(keyCache, string(dataByte), TTL_KVCACHE)

	return RelatedVideosObjOutput, err
}

func GetListIDRelatedVideoByZrange(groupId string, page, limit int, cache bool) ([]string, error) {
	var listVodID []string
	var keyCache = PREFIX_REDIS_ZRANGE_RELATED_VIDEOS + "_" + groupId
	var start = page * limit
	var stop = (start + limit) - 1

	if cache {
		vals, err := mRedis.ZRange(keyCache, int64(start), int64(stop))
		if err == nil && len(vals) > 0 {
			listVodID = vals
		}
	}

	if len(listVodID) <= 0 {
		// Connect DB
		session, db, err := GetCollection()
		if err != nil {
			Sentry_log(err)
			return listVodID, err
		}
		defer session.Close()

		var where = bson.M{
			"group_id": groupId,
			"type":     bson.M{"$in": []int{6, 8}},
		}

		// Lây tất cả releated
		var VODDataObjects []VODDataObjectStruct
		err = db.C(COLLECTION_VOD).Find(where).Sort("odr").All(&VODDataObjects)
		if err != nil {
			Sentry_log(err)
			return listVodID, err
		}

		var listVodIdAll []string
		//Write cache
		mRedis.Del(keyCache)
		for _, val := range VODDataObjects {
			listVodIdAll = append(listVodIdAll, val.Id)
			mRedis.ZAdd(keyCache, float64(val.Odr), val.Id)
		}

		for i := start; i <= stop; i++ {
			if i < len(listVodIdAll) {
				listVodID = append(listVodID, listVodIdAll[i])
			}
		}
	}
	return listVodID, nil
}

func GetListIDRelatedVideos(groupId string, page int, limit int, platforms string, cacheActive bool) ([]string, error) {
	platform := Platform(platforms)
	var keyCache = RELATED_VIDEOS + "_" + groupId + "_" + strconv.Itoa(platform.Id) + "_" + strconv.Itoa(page) + "_" + strconv.Itoa(limit)
	var listVodID []string

	if cacheActive {
		dataCache, err := mRedis.GetString(keyCache)
		if err == nil && dataCache != "" {
			err = json.Unmarshal([]byte(dataCache), &listVodID)
			if err == nil {
				return listVodID, nil
			}
		}
	}

	// Connect DB
	var VODDataObjects []VODDataObjectStruct
	session, db, err := GetCollection()
	if err != nil {
		return listVodID, err
	}
	defer session.Close()

	var where = bson.M{
		"group_id": groupId,
		"type":     bson.M{"$in": []int{6, 8}},
	}

	err = db.C(COLLECTION_VOD).Find(where).Sort("-release_date").Skip(page * limit).Limit(limit).All(&VODDataObjects)
	if err != nil || len(VODDataObjects) <= 0 {
		return listVodID, err
	}

	for _, val := range VODDataObjects {
		listVodID = append(listVodID, val.Id)
	}

	// Write cache
	dataByte, _ := json.Marshal(listVodID)
	mRedis.SetString(keyCache, string(dataByte), TTL_REDIS_LV1)

	return listVodID, nil
}

func GetTotalRelatedVideos(groupId string, platforms string, cacheActive bool) (int, error) {
	var keyCache = PREFIX_REDIS_ZRANGE_RELATED_VIDEOS + groupId
	total := mRedis.ZCountAll(keyCache)
	return int(total), nil

	// platform := Platform(platforms)
	// var keyCache = RELATED_VIDEOS_TOTAL + "_" + groupId + "_" + strconv.Itoa(platform.Id)
	// if cacheActive {
	// 	total, err := mRedis.GetInt(keyCache)
	// 	if err == nil {
	// 		return total, nil
	// 	}
	// }

	// // Connect DB
	// session, db, err := GetCollection()
	// if err != nil {
	// 	return 0, err
	// }
	// defer session.Close()

	// var where = bson.M{
	// 	"group_id":  groupId,
	// 	"type":      bson.M{"$in": []int{6, 8}},
	// 	"platforms": platform.Id,
	// }

	// //Get total related_videos
	// total, _ := db.C(COLLECTION_VOD).Find(where).Count()
	// mRedis.SetInt(keyCache, total, TTL_REDIS_LV1)

	// return total, nil
}
