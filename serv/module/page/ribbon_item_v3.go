package page

import (
	. "cm-v5/schema"
	. "cm-v5/serv/module"
	seo "cm-v5/serv/module/seo"
	vod "cm-v5/serv/module/vod"
	"errors"
	"strings"

	"gopkg.in/mgo.v2/bson"
)

func GetListRibbonItemV3ByID(ribbonID string, platform string, page, limit int, cacheActive bool) ([]VODDataObjectStruct, error) {
	var RibItemObjects = make([]VODDataObjectStruct, 0)
	var listRibItemID []string
	var keyCache = LIST_RIB_ITEM_V3 + "_" + ribbonID + "_" + platform
	var start = page * limit
	var stop = (start + limit) - 1

	if cacheActive {
		vals, err := mRedis.ZRange(keyCache, int64(start), int64(stop))
		if err == nil && len(vals) > 0 {
			listRibItemID = vals
		}
	}

	if len(listRibItemID) <= 0 {
		var RibbonItemsV3 = make([]RibbonItemV3ObjectStruct, 0)
		// Connect MongoDB
		session, db, err := GetCollection()
		if err != nil {
			Sentry_log(err)
			return RibItemObjects, err
		}
		defer session.Close()

		var platformInfo = Platform(platform)
		var where = bson.M{
			"rib_ids":   ribbonID,
			"status":    1,
			"platforms": platformInfo.Id,
		}

		err = db.C(COLLECTION_RIB_ITEM_V3).Find(where).Sort("odr").All(&RibbonItemsV3)
		if err != nil {
			Sentry_log(err)
			return RibItemObjects, err
		}

		var listRibItemdAll []string

		//Write cache
		mRedis.Del(keyCache)
		for _, val := range RibbonItemsV3 {
			listRibItemdAll = append(listRibItemdAll, val.Rib_item_id)
			AddDataToCacheZRange(keyCache, val.Rib_item_id, float64(val.Odr))
		}

		for i := start; i <= stop; i++ {
			if i < len(listRibItemdAll) {
				listRibItemID = append(listRibItemID, listRibItemdAll[i])
			}
		}
	}

	var ResultRibItemObjects = make([]VODDataObjectStruct, 0)
	if len(listRibItemID) > 0 {
		ResultRibItemObjects, _ = GetListRibItemPushCache(listRibItemID, platform, cacheActive)
	}

	return ResultRibItemObjects, nil
}

func RemoveDataToCacheZRange(keyCache, RibItemId string) {
	mRedis.ZRem(keyCache, RibItemId)
}

func AddDataToCacheZRange(keyCache, RibItemId string, odr float64) {
	mRedis.ZAdd(keyCache, float64(odr), RibItemId)
}

func GetListRibItemPushCache(listRibItemID []string, platform string, cacheActive bool) ([]VODDataObjectStruct, error) {
	var shortcutItemsObjects = make([]VODDataObjectStruct, 0)
	if len(listRibItemID) <= 0 {
		return shortcutItemsObjects, errors.New("GetVODByListID: Empty data ")
	}

	if cacheActive {
		dataRedis, err := mRedis.HMGet(PREFIX_REDIS_HASH_RIB_ITEM_V3, listRibItemID)
		if err == nil && len(dataRedis) > 0 {
			for _, val := range dataRedis {
				if str, ok := val.(string); ok {
					var shortcutDataObject VODDataObjectStruct
					err = json.Unmarshal([]byte(str), &shortcutDataObject)
					if err != nil {
						continue
					}
					shortcutItemsObjects = append(shortcutItemsObjects, shortcutDataObject)
				}
			}
		}
	}

	if len(shortcutItemsObjects) != len(listRibItemID) {
		// Connect DB
		session, db, err := GetCollection()
		if err != nil {
			Sentry_log(err)
			return shortcutItemsObjects, err
		}

		defer session.Close()

		var where = bson.M{
			"rib_item_id": bson.M{"$in": listRibItemID},
			"status":      1,
		}

		var RibbonItemsV3Objects = make([]RibbonItemV3ObjectStruct, 0)
		err = db.C(COLLECTION_RIB_ITEM_V3).Find(where).All(&RibbonItemsV3Objects)
		if err != nil {
			Sentry_log(err)
			return shortcutItemsObjects, err
		}

		if len(RibbonItemsV3Objects) <= 0 {
			// Remove cache
			for _, id := range listRibItemID {
				mRedis.HDel(PREFIX_REDIS_HASH_RIB_ITEM_V3, id)
			}
			return shortcutItemsObjects, errors.New("GetVODByListID: Empty data")
		}

		//Lay thong tin cua VOD
		var listVodId []string
		for _, val := range RibbonItemsV3Objects {
			listVodId = append(listVodId, val.Ref_id)
		}

		var vodDataObjects = make([]VODDataObjectStruct, 0)
		if len(listVodId) > 0 {
			var platformInfo = Platform(platform)
			vodDataObjects, _ = vod.GetVODByListID(listVodId, platformInfo.Id, 1, true)
		}

		dataByte, _ := json.Marshal(RibbonItemsV3Objects)
		err = json.Unmarshal([]byte(dataByte), &shortcutItemsObjects)
		if err != nil {
			return shortcutItemsObjects, err
		}

		var shortcutItemsMapping = make([]VODDataObjectStruct, 0)
		for k, valVod := range shortcutItemsObjects {

			//những thông tin rib_item được giữ lại
			valVod.Rib_item_id = RibbonItemsV3Objects[k].Rib_item_id
			valVod.Id = RibbonItemsV3Objects[k].Ref_id
			valVod.Title = RibbonItemsV3Objects[k].Title
			valVod.Rib_ids = RibbonItemsV3Objects[k].Rib_ids
			valVod.Seo.Title = RibbonItemsV3Objects[k].Title
			valVod.Seo.Description = RibbonItemsV3Objects[k].Short_description
			valVod.Image_soucre = RibbonItemsV3Objects[k].Image_soucre
			valVod.Images = RibbonItemsV3Objects[k].Images
			valVod.Type_name = RibbonItemsV3Objects[k].Type_name


			if valVod.Type == 10 {
				// item là ribbon
				slugSplits := strings.Split(valVod.Seo.Url, "/tuyen-tap/")
				if len(slugSplits) == 2 {
					valVod.Seo = seo.FormatSeoRibbon(valVod.Id, slugSplits[1], valVod.Title, 10, "", cacheActive)
				}
			}

			for _, val := range vodDataObjects {
				if valVod.Id == val.Id {
					MappingDataFromVodToRibItem(&valVod, val)
				}
			}
			shortcutItemsMapping = append(shortcutItemsMapping, valVod)
		}

		shortcutItemsObjects = shortcutItemsMapping
		// Write cache
		for _, val := range shortcutItemsObjects {
			dataByte, _ := json.Marshal(val)
			// Write cache follow id
			mRedis.HSet(PREFIX_REDIS_HASH_RIB_ITEM_V3, val.Rib_item_id, string(dataByte))
		}
	}

	var TypeName = []string{"MOVIE", "SHOW", "SEASON", "EPISODE", "VOD"}
	//Sap xep mang theo thu tu id truyen vao
	var shortcutItemsObjectsSort = make([]VODDataObjectStruct, 0)
	for _, id := range listRibItemID {
		for _, valShort := range shortcutItemsObjects {
			if id == valShort.Rib_item_id {
				exists, _ := In_array(valShort.Type_name, TypeName)
				if exists {
					var vodDataObjects = make([]VODDataObjectStruct, 0)
					var platformInfo = Platform(platform)
					var listVodId = []string{valShort.Id}
					vodDataObjects, _ = vod.GetVODByListID(listVodId, platformInfo.Id, 1, true)
					if len(vodDataObjects) <= 0 {
						continue
					}
					//mapping data from VOD
					MappingDataFromVodToRibItem(&valShort, vodDataObjects[0])

					//Khanh DT-15025
					valShort.Avg_rate = ReCalculatorRating(vodDataObjects[0].Min_rate, vodDataObjects[0].Avg_rate)

				}

				shortcutItemsObjectsSort = append(shortcutItemsObjectsSort, valShort)
			}
		}
	}

	return shortcutItemsObjectsSort, nil
}

func MappingDataFromVodToRibItem(dataRibitem *VODDataObjectStruct, dataVod VODDataObjectStruct) {
	dataRibitem.Seo.Url = dataVod.Seo.Url
	dataRibitem.Group_id = dataVod.Group_id
	dataRibitem.Resolution = dataVod.Resolution
	dataRibitem.Avg_rate = dataVod.Avg_rate
	dataRibitem.Min_rate = dataVod.Min_rate
	dataRibitem.Total_rate = dataVod.Total_rate
	dataRibitem.Is_watchlater = dataVod.Is_watchlater
	dataRibitem.Episode = dataVod.Episode
	dataRibitem.Current_episode = dataVod.Current_episode
	dataRibitem.Slug = dataVod.Slug
	dataRibitem.R_id = dataVod.R_id
	dataRibitem.Share_url = dataVod.Share_url
	dataRibitem.Share_url_seo = dataVod.Share_url_seo
	dataRibitem.Is_premium = dataVod.Is_premium
	dataRibitem.Label_subtitle_audio = dataVod.Label_subtitle_audio
	dataRibitem.Label_public_day = dataVod.Label_public_day
	dataRibitem.Is_coming_soon = dataVod.Is_coming_soon
	dataRibitem.Tags = dataVod.Tags
	dataRibitem.Release_year = dataVod.Release_year
	dataRibitem.Link_play.Hls_link_play = dataVod.Trailer_link_play.Hls_link_play
	dataRibitem.Link_play.Dash_link_play = dataVod.Trailer_link_play.Dash_link_play
	dataRibitem.Tags_display = dataVod.Tags_display
	dataRibitem.Ranking = dataVod.Ranking
	dataRibitem.People = dataVod.People
	
}

//no cache for api update cache
func GetRibbonItemByIdNoCache(rib_item_id string) (RibbonItemV3ObjectStruct, error) {
	var RibbonItem RibbonItemV3ObjectStruct

	// Connect MongoDB
	session, db, err := GetCollection()
	if err != nil {
		return RibbonItem, err
	}

	defer session.Close()

	var where = bson.M{
		"rib_item_id": rib_item_id,
	}

	err = db.C(COLLECTION_RIB_ITEM_V3).Find(where).One(&RibbonItem)
	if err != nil && err.Error() != "not found" {
		return RibbonItem, err
	}

	return RibbonItem, nil
}
