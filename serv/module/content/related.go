package content

import (
	"errors"
	"strconv"
	"strings"
	. "cm-v5/serv/module"
	. "cm-v5/schema"
	seo "cm-v5/serv/module/seo"
	vod "cm-v5/serv/module/vod"
	recommendation "cm-v5/serv/module/recommendation"
	"gopkg.in/mgo.v2/bson"
)


func GetRelatedYusp(user_id, device_id, platforms, contentId string, page, limit int ) (RelatedObjOutputStruct, error) {
	platform := Platform(platforms)
	var RelatedObjOutput RelatedObjOutputStruct

	// get detail info by content id
	vodObj, err := vod.GetVodDetail(contentId, platform.Id, true)
	if err != nil {
		return RelatedObjOutput, err
	}

	var arrContentId []string
	val , _ := mRedis.GetString("CONFIG_YUSP_DEFAULT_CONTENT_DETAIL")
	if val != "" {
		arrContentId = strings.Split(val, ",")
	} else {
		scenario := recommendation.GetScenariosContentDetail()
		dataGetYusp, _ := recommendation.GetContentDetailYuspRecommendationAPI(user_id, device_id, scenario, platforms, contentId, page, limit)
		arrContentId = dataGetYusp.ItemIds
	}
	
	// Lay danh sach VOD Lien Quan
	if len(arrContentId) > 0 && err == nil {
		var VODDataObjects []VODDataObjectStruct
		VODDataObjects, err = vod.GetVODByListID(arrContentId, platform.Id, 1, true)
		if err != nil {
			return RelatedObjOutput, err
		}

		//nếu có tồn tại single movie là show thì loại bỏ
		var newVODDataObjects []VODDataObjectStruct
		for _, vodData := range VODDataObjects {
			if vodData.Type == 1 && vodData.Category == 2 {
				continue
			}
			newVODDataObjects = append(newVODDataObjects, vodData)
		}

		VODDataObjects = newVODDataObjects

		dataByte, _ := json.Marshal(VODDataObjects)
		err = json.Unmarshal(dataByte, &RelatedObjOutput.Items)
		if err != nil {
			return RelatedObjOutput, err
		}

		//Switch images platform
		for k, vodData := range VODDataObjects {

			var ImagesMapping ImagesOutputObjectStruct
			switch platform.Type {
			case "web":
				ImagesMapping.Thumbnail = BuildImage(vodData.Images.Web.Vod_thumb)
			case "smarttv":
				ImagesMapping.Thumbnail = BuildImage(vodData.Images.Smarttv.Thumbnail)
			case "app":
				ImagesMapping.Thumbnail = BuildImage(vodData.Images.App.Thumbnail)
			}
			ImagesMapping = MappingImagesV4(platform.Type, ImagesMapping, vodData.Images, true)
			RelatedObjOutput.Items[k].Images = ImagesMapping
		}

	}

	//Pagination
	RelatedObjOutput.Metadata.Total = len(RelatedObjOutput.Items)
	RelatedObjOutput.Metadata.Limit = limit
	RelatedObjOutput.Metadata.Page = page

	//get data by content id to create seo info
	// dataContent, _ := GetContent(contentId, platforms, true)
	RelatedObjOutput.Seo = seo.FormatSeoRelated(vodObj.Title, vodObj.Seo.Url)

	return RelatedObjOutput, nil
}


func GetRelated(contentId string, page int, limit int, platforms string, cacheActive bool) (RelatedObjOutputStruct, error) {
	platform := Platform(platforms)
	var RelatedObjOutput RelatedObjOutputStruct
	var keyCache = KV_RELATED + "_" + strconv.Itoa(page) + "_" + strconv.Itoa(limit) + "_" + contentId + "_" + platforms

	if cacheActive {
		valueCache, err := mRedisKV.GetString(keyCache)
		if err == nil && valueCache != "" {
			err = json.Unmarshal([]byte(valueCache), &RelatedObjOutput)
			if err == nil {
				return RelatedObjOutput, nil
			}
		}
	}

	// get detail info by content id
	vodObj, err := vod.GetVodDetail(contentId, platform.Id, cacheActive)
	if err != nil {
		return RelatedObjOutput, err
	}

	// Get list tag id by content id
	var listTagId []string
	for _, tag := range vodObj.Tags {
		listTagId = append(listTagId, tag.Id)
	}

	var total int = 0

	// Lay danh sach VOD Lien Quan
	listVodID, err := GetListIDRelatedVOD(vodObj.Id, listTagId, page, limit, platforms, cacheActive)
	if len(listVodID) > 0 && err == nil {

		var VODDataObjects []VODDataObjectStruct
		VODDataObjects, err = vod.GetVODByListID(listVodID, platform.Id, 1, cacheActive)
		if err != nil {
			return RelatedObjOutput, err
		}

		//nếu có tồn tại single movie là show thì loại bỏ
		// var newVODDataObjects []VODDataObjectStruct
		// for _, vodData := range VODDataObjects {
		// 	if vodData.Type == 1 && vodData.Category == 2 {
		// 		continue
		// 	}
		// 	newVODDataObjects = append(newVODDataObjects, vodData)
		// }

		dataByte, _ := json.Marshal(VODDataObjects)
		err = json.Unmarshal(dataByte, &RelatedObjOutput.Items)
		if err != nil {
			return RelatedObjOutput, err
		}

		//Switch images platform
		for k, vodData := range VODDataObjects {

			var ImagesMapping ImagesOutputObjectStruct
			switch platform.Type {
			case "web":
				ImagesMapping.Thumbnail = BuildImage(vodData.Images.Web.Vod_thumb)
			case "smarttv":
				ImagesMapping.Thumbnail = BuildImage(vodData.Images.Smarttv.Thumbnail)
			case "app":
				ImagesMapping.Thumbnail = BuildImage(vodData.Images.App.Thumbnail)
			}
			ImagesMapping = MappingImagesV4(platform.Type, ImagesMapping, vodData.Images, true)
			RelatedObjOutput.Items[k].Images = ImagesMapping
		}

	}
	total, _ = GetTotalRelatedVOD(vodObj.Id, listTagId, platforms, cacheActive)
	//Pagination
	RelatedObjOutput.Metadata.Total = total
	RelatedObjOutput.Metadata.Limit = limit
	RelatedObjOutput.Metadata.Page = page

	//get data by content id to create seo info
	// dataContent, _ := GetContent(contentId, platforms, true)
	RelatedObjOutput.Seo = seo.FormatSeoRelated(vodObj.Title, vodObj.Seo.Url)

	// Write Redis
	dataByte, _ := json.Marshal(RelatedObjOutput)
	mRedisKV.SetString(keyCache, string(dataByte), TTL_KVCACHE)

	return RelatedObjOutput, nil
}

func GetListIDRelatedVOD(contentID string, listTagId []string, page int, limit int, platforms string, cacheActive bool) ([]string, error) {
	platform := Platform(platforms)
	var keyCache = RELATED_VOD + "_" + contentID + "_" + strconv.Itoa(page) + "_" + strconv.Itoa(limit) + "_" + platforms
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
		"tags.id":   bson.M{"$in": listTagId},
		"id":        bson.M{"$ne": contentID},
		"platforms": platform.Id,
		"type":      bson.M{"$ne": 1}, //nếu có tồn tại single movie là show thì loại bỏ
		"category":  bson.M{"$ne": 2}, //nếu có tồn tại single movie là show thì loại bỏ
	}

	err = db.C(COLLECTION_VOD).Find(where).Sort("-release_date").Skip(page * limit).Limit(limit).All(&VODDataObjects)
	if err != nil || len(VODDataObjects) <= 0 {
		return listVodID, err
	}

	if len(VODDataObjects) <= 0 {
		return listVodID, errors.New("Empty data")
	}

	for _, val := range VODDataObjects {
		listVodID = append(listVodID, val.Id)
	}

	// Write cache
	dataByte, _ := json.Marshal(listVodID)
	mRedis.SetString(keyCache, string(dataByte), TTL_REDIS_LV1)

	return listVodID, nil
}

func GetTotalRelatedVOD(contentID string, listTagId []string, platforms string, cacheActive bool) (int, error) {
	platform := Platform(platforms)
	var keyCache = RELATED_VOD_TOTAL + "_" + contentID + "_" + platforms
	if cacheActive {
		total, err := mRedis.GetInt(keyCache)
		if err == nil {
			return total, nil
		}
	}

	// Connect DB
	session, db, err := GetCollection()
	if err != nil {
		return 0, err
	}
	defer session.Close()

	var where = bson.M{
		"tags.id":   bson.M{"$in": listTagId},
		"id":        bson.M{"$ne": contentID},
		"platforms": platform.Id,
		"type":      bson.M{"$ne": 1}, //nếu có tồn tại single movie là show thì loại bỏ
		"category":  bson.M{"$ne": 2}, //nếu có tồn tại single movie là show thì loại bỏ
	}

	//Get total related_videos
	total, _ := db.C(COLLECTION_VOD).Find(where).Count()
	mRedis.SetInt(keyCache, total, TTL_REDIS_LV1)

	return total, nil
}
