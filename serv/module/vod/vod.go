package vod

import (
	. "cm-v5/schema"
	. "cm-v5/serv/module"
	seo "cm-v5/serv/module/seo"
	topviews "cm-v5/serv/module/topviews"
	"encoding/json"
	"errors"
	"fmt"
	"math"

	jsoniter "github.com/json-iterator/go"
	"gopkg.in/mgo.v2/bson"
	// live_event_finished "cm-v5/serv/module/live_event_finished"
)

var mRedis RedisModelStruct

/**
	Status
		- 0: khong tam status vod
		- con lại: lấy vod theo status truyền vào
**/

func GetVODByListID(listId []string, platform, status int, cache bool) ([]VODDataObjectStruct, error) {

	var vodDataObjects = make([]VODDataObjectStruct, 0)
	var json = jsoniter.ConfigCompatibleWithStandardLibrary

	if len(listId) <= 0 {
		return vodDataObjects, errors.New("GetVODByListID: Empty data ")
	}
	if cache {
		dataRedis, err := mRedis.HMGet(PREFIX_REDIS_HASH_CONTENT, listId)
		if err != nil {
			Sentry_log(err)
			return vodDataObjects, err
		}

		if err == nil && len(dataRedis) > 0 {
			for _, val := range dataRedis {
				if str, ok := val.(string); ok {
					var vodDataObject VODDataObjectStruct
					err = json.Unmarshal([]byte(str), &vodDataObject)
					if err != nil {
						continue
					}
					vodDataObjects = append(vodDataObjects, vodDataObject)
				}
			}
		}
	}

	var DataSeasion = make(map[string][]VODDataObjectStruct)
	if len(vodDataObjects) != len(listId) {

		//reset vodDataObjects = []
		vodDataObjects = make([]VODDataObjectStruct, 0)

		// Connect DB
		session, db, err := GetCollection()
		if err != nil {
			Sentry_log(err)
			return vodDataObjects, err
		}
		defer session.Close()

		var where = bson.M{
			"id": bson.M{"$in": listId},
		}

		err = db.C(COLLECTION_VOD).Find(where).All(&vodDataObjects)
		if err != nil {
			Sentry_log(err)
			return vodDataObjects, err
		}

		if len(vodDataObjects) <= 0 {
			// Remove cache
			for _, id := range listId {
				mRedis.HDel(PREFIX_REDIS_HASH_CONTENT, id)
			}
			return vodDataObjects, errors.New("GetVODByListID: Empty data")
		}

		// Write cache
		for i, val := range vodDataObjects {

			//lấy thông tin từ season nếu type là episode hoặc trailer
			if (val.Type == VOD_TYPE_EPISODE || val.Type == VOD_TYPE_TRAILER) && val.Id != val.Group_id {
				if len(DataSeasion[val.Group_id]) <= 0 {
					ArrIds := []string{val.Group_id}
					DataSeasion[val.Group_id], err = GetVODByListID(ArrIds, platform, 0, true)
					if err != nil {
						continue
					}
				}

				//KHANH DT-11739
				if len(DataSeasion[val.Group_id]) > 0 {
					vodDataObjects[i].Avg_rate = DataSeasion[val.Group_id][0].Avg_rate
					vodDataObjects[i].Total_rate = DataSeasion[val.Group_id][0].Total_rate
					vodDataObjects[i].Resolution = DataSeasion[val.Group_id][0].Resolution
					vodDataObjects[i].Tags = DataSeasion[val.Group_id][0].Tags
					vodDataObjects[i].Release_year = DataSeasion[val.Group_id][0].Release_year
				}
			}
			if val.Type == VOD_TYPE_MOVIE || val.Type == VOD_TYPE_SEASON {
				//get rating in mongo
				dataRatingContent := GetDataRatingContentInMongo(val.Id)
				if dataRatingContent.Total_rate > 0 {
					vodDataObjects[i].Avg_rate = dataRatingContent.Avg_rate
					vodDataObjects[i].Total_rate = dataRatingContent.Total_rate
				}
				vodDataObjects[i].Avg_rate = ReCalculatorRating(vodDataObjects[i].Min_rate, vodDataObjects[i].Avg_rate)
			}

			dataByte, _ := json.Marshal(vodDataObjects[i])
			// Write cache follow ID
			mRedis.HSet(PREFIX_REDIS_HASH_CONTENT, val.Id, string(dataByte))
			// Write cache follow Slug
			mRedis.HSet(PREFIX_REDIS_HASH_CONTENT_SLUG, val.Slug_seo, string(val.Id))

			//add cache zrange VOD_TYPE_LIVE_EVENT_FINISHED
			if val.Type == VOD_TYPE_LIVE_EVENT_FINISHED {
				mRedis.ZRem(PREFIX_REDIS_ZRANGE_LIVE_EVENT_FINISHED, val.Id)
				mRedis.ZAdd(PREFIX_REDIS_ZRANGE_LIVE_EVENT_FINISHED, float64(val.Created_at), val.Id)
			}
		}

		//Clear key cache hash vod slug sau 1 ngay
		mRedis.Expire(PREFIX_REDIS_HASH_CONTENT_SLUG, TTL_REDIS_LV1)
	}

	// Check platform
	// var VODDataObjectNew []VODDataObjectStruct
	var DataVods = make(map[string]VODDataObjectStruct)
	for _, valVod := range vodDataObjects {
		exists, _ := In_array(platform, valVod.Platforms)
		if (exists || platform == 0) && (0 == status || (0 != status && 9 != valVod.Status)) {
			//GetTopViewsRankingByContentID
			valVod.Ranking = topviews.GetTopViewsRankingByContentID(valVod.Id, true)

			valVod.Seo = seo.FormatSeoVOD(valVod.Slug_seo_v5, valVod.Seo)
			valVod.Share_url = HandleSeoShareUrlVieON(valVod.Seo.Url)
			valVod.Share_url_seo = HandleSeoShareUrlVieON(valVod.Seo.Url)
			// VODDataObjectNew = append(VODDataObjectNew, valVod)
			DataVods[valVod.Id] = valVod
		}
	}

	// Sap xep mang theo thu tu id truyen vao
	var VODDataObjectSort []VODDataObjectStruct
	for _, valId := range listId {

		if _, ok := DataVods[valId]; ok {
			if DataVods[valId].Id != "" {
				VODDataObjectSort = append(VODDataObjectSort, DataVods[valId])
				delete(DataVods, valId)
			}
		}

	}

	return VODDataObjectSort, nil
}

func GetVodDetailBySlug(entitySlug string, platform int) (VODDataObjectStruct, error) {
	var VODDataObject VODDataObjectStruct
	var keyCache = "vod_detail_slug_temp_" + entitySlug
	valueCache, err := mRedis.GetString(keyCache)
	if err == nil {
		err = json.Unmarshal([]byte(valueCache), &VODDataObject)
		if err == nil {
			return VODDataObject, nil
		}
	}

	// Connect DB
	session, db, err := GetCollection()
	if err != nil {
		Sentry_log(err)
		return VODDataObject, err
	}
	defer session.Close()

	var where = bson.M{
		"slug_seo_v5": entitySlug,
	}

	err = db.C(COLLECTION_VOD).Find(where).One(&VODDataObject)
	if err != nil && err.Error() != "not found" {
		return VODDataObject, err
	}

	// Write cache
	mRedis.HSet(PREFIX_REDIS_HASH_CONTENT_SLUG, VODDataObject.Slug_seo_v5, string(VODDataObject.Id))

	// Write cache
	dataByte, _ := json.Marshal(VODDataObject)
	mRedis.SetString(keyCache, string(dataByte), TTL_REDIS_5_MINUTE)

	return VODDataObject, nil
}

func GetVodIdBySlug(entitySlug string, platform int) (string, error) {
	dataRedisBySlugId, _ := mRedis.HMGet(PREFIX_REDIS_HASH_CONTENT_SLUG, []string{entitySlug})
	for _, val := range dataRedisBySlugId {
		if str, ok := val.(string); ok {
			return str, nil
		}
	}

	// Check slug Movie / Seasion
	dataVodBySlug, err := GetVodDetailBySlug(entitySlug, platform)
	return dataVodBySlug.Id, err
}

func GetVodDetail(vodId string, platform int, cache bool) (VODDataObjectStruct, error) {
	var listID []string
	var VODDataObject VODDataObjectStruct

	if vodId == "" {
		return VODDataObject, errors.New("GetVodDetail: Empty data - " + vodId)
	}

	var keyCache = "vod_detail_id_temp_" + vodId + "_" + fmt.Sprint(platform)

	if cache {
		valueCache, err := mRedis.GetString(keyCache)
		if err == nil {
			err = json.Unmarshal([]byte(valueCache), &VODDataObject)
			if err == nil {
				return VODDataObject, nil
			}
		}

	}

	// Get VodDetail by ID
	listID = append(listID, vodId)
	dataVodDetail, err := GetVODByListID(listID, platform, 0, cache)

	if err != nil && err.Error() != "GetVODByListID: Empty data" {
		return VODDataObject, err
	}

	if len(dataVodDetail) > 0 {
		dataByte, err := json.Marshal(dataVodDetail[0])
		err = json.Unmarshal(dataByte, &VODDataObject)
		if err != nil {
			Sentry_log(err)
			return VODDataObject, err
		}
		// Write cache
		mRedis.SetString(keyCache, string(dataByte), TTL_REDIS_5_MINUTE)
		return VODDataObject, nil
	}

	// Write cache
	dataByte, _ := json.Marshal(VODDataObject)
	mRedis.SetString(keyCache, string(dataByte), TTL_REDIS_5_MINUTE)

	return VODDataObject, errors.New("GetVodDetail: Empty data - " + vodId)
}

func GetDefaultEpisodeByGroupId(groupId string, platformId int) VODDataObjectStruct {
	var VODDataObject VODDataObjectStruct
	VODDataObjects, _ := GetVodByListGroup(groupId, platformId, 0, 10, true)
	if len(VODDataObjects) > 0 {

		if VODDataObjects[0].Episode < VODDataObjects[len(VODDataObjects)-1].Episode {
			VODDataObject = VODDataObjects[0]
		} else {
			VODDataObjectZRangeByScore, _ := GetDefaulEpsByZRevRange(groupId, platformId, 0, 1)
			if len(VODDataObjectZRangeByScore) > 0 {
				VODDataObject = VODDataObjectZRangeByScore[0]
			}
		}
	}

	return VODDataObject
}

func GetDefaulEpsByZRevRange(groupId string, platformId, page, limit int) ([]VODDataObjectStruct, error) {
	var keyCache = CONTENT_LIST_EPS + groupId
	var start = page * limit
	var stop = (start + limit) - 1
	var listVodID []string

	vals := mRedis.ZRevRange(keyCache, int64(start), int64(stop))
	if len(vals) > 0 {
		listVodID = vals
	}

	var ResultVodData = make([]VODDataObjectStruct, 0)
	if len(listVodID) > 0 {
		ResultVodData, _ = GetVODByListID(listVodID, platformId, 1, true)
	}

	return ResultVodData, nil

}

func GetListIdEpisodeByCache(groupId string, platformId, page, limit int) []string {
	var keyCache = CONTENT_LIST_EPS + groupId
	var start = page * limit
	var stop = (start + limit) - 1
	var listVodID []string

	vals, err := mRedis.ZRange(keyCache, int64(start), int64(stop))
	if err == nil && len(vals) > 0 {
		listVodID = vals
	}
	return listVodID
}

func GetVodByListGroup(groupId string, platformId, page, limit int, cache bool) ([]VODDataObjectStruct, error) {
	var ListVodData []VODDataObjectStruct
	var listVodID []string
	var keyCache = CONTENT_LIST_EPS + groupId
	var start = page * limit
	var stop = (start + limit) - 1

	if cache {
		listVodID = GetListIdEpisodeByCache(groupId, platformId, page, limit)
	}

	if len(listVodID) <= 0 {
		// Connect DB
		session, db, err := GetCollection()
		if err != nil {
			Sentry_log(err)
			return ListVodData, err
		}
		defer session.Close()

		var where = bson.M{
			"group_id": groupId,
			"type":     4,
			"status":   bson.M{"$ne": 9},
		}

		// Lây tất cả epi
		err = db.C(COLLECTION_VOD).Find(where).Sort("odr").All(&ListVodData)
		if err != nil {
			Sentry_log(err)
			return ListVodData, err
		}

		var listVodIdAll []string
		//Write cache
		mRedis.Del(keyCache)
		for _, val := range ListVodData {
			listVodIdAll = append(listVodIdAll, val.Id)
			mRedis.ZAdd(keyCache, float64(val.Odr), val.Id)
		}
		for i := start; i <= stop; i++ {

			if i < len(listVodIdAll) {
				listVodID = append(listVodID, listVodIdAll[i])
			}

			// if Isset(listVodIdAll, i) {
			// 	listVodID = append(listVodID, listVodIdAll[i])
			// }
		}
	}

	if len(listVodID) > 0 {
		VodsData, _ := GetVODByListID(listVodID, 0, 1, cache)

		//check miss cache
		if len(VodsData) < len(listVodID) {
			VodsData, _ = GetVODByListID(listVodID, 0, 1, false)
			for _, ID := range listVodID {
				isPublic := false
				for _, val := range VodsData {
					if val.Id == ID {
						isPublic = true
						break
					}

				}
				if isPublic == false {
					mRedis.ZRem(keyCache, ID)
				}
			}
		}

		ListVodData = make([]VODDataObjectStruct, 0)
		//check platform
		for _, val := range VodsData {
			exists, _ := In_array(platformId, val.Platforms)
			if exists {
				ListVodData = append(ListVodData, val)
			}
		}
	}

	return ListVodData, nil
}

func WriteCacheZRangeByListGroup(keyCache string, score float64, value string) error {
	err := mRedis.ZAdd(keyCache, score, value)
	return err
}

func GetPageIndexZRangeByGroupIdAndEpsID(groupId, epsId string) int64 {
	var keyCache = CONTENT_LIST_EPS + groupId
	var indexOfRank int64 = mRedis.ZRank(keyCache, epsId)
	// get indexPage
	var indexPage = math.Floor(float64(indexOfRank) / LIMIT_EPS_PER_RANGE)
	return int64(indexPage)
}

func GetDataRatingContentInMongo(contentId string) RatingDetailContentStruct {
	var RatingDetailContent RatingDetailContentStruct
	session, db, err := GetCollection()
	if err != nil {
		Sentry_log(err)
	}
	defer session.Close()
	var where = bson.M{
		"contentid": contentId,
	}
	var rateObj RateObjectStruct
	err = db.C(COLLECTION_RATE).Find(where).One(&rateObj)
	// valueFromDb := int64(0)
	if err == nil && rateObj.Total_rate > 0 {
		RatingDetailContent.Avg_rate = rateObj.Avg_rate
		RatingDetailContent.Total_rate = rateObj.Total_rate
		RatingDetailContent.Total_point = rateObj.Total_point
		return RatingDetailContent
	}

	return RatingDetailContent
}
