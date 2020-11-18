package tag

import (
	"errors"
	"fmt"
	"strings"

	. "cm-v5/schema"
	. "cm-v5/serv/module"
	seo "cm-v5/serv/module/seo"
	vod "cm-v5/serv/module/vod"

	"gopkg.in/mgo.v2/bson"
)

func GetVODByTagIds(listTagID []string, sort, page, limit int, platforms string, cacheActive bool) (TagsOutputObjectStruct, error) {
	platform := Platform(platforms)
	var TagsOutputObjects TagsOutputObjectStruct
	var keyCache = TAGS_VOD + "_" + strings.Join(listTagID, "-") + "_" + fmt.Sprint(sort) + "_" + fmt.Sprint(page) + "_" + fmt.Sprint(limit) + "_" + platforms
	if cacheActive {
		valueCache, err := mRedis.GetString(keyCache)
		if err == nil && valueCache != "" {
			err = json.Unmarshal([]byte(valueCache), &TagsOutputObjects)
			if err == nil {
				return TagsOutputObjects, nil
			}
		}
	}

	//Lay danh sach id vod
	listVodID, err := GetListIDTagVOD(listTagID, sort, page, limit, platforms, cacheActive)
	if len(listVodID) > 0 && err == nil {

		var vodDataObjects []VODDataObjectStruct
		vodDataObjects, err := vod.GetVODByListID(listVodID, platform.Id, 1, true)
		if err != nil {
			return TagsOutputObjects, err
		}

		dataByte, _ := json.Marshal(vodDataObjects)
		err = json.Unmarshal(dataByte, &TagsOutputObjects.Items)
		if err != nil {
			return TagsOutputObjects, err
		}

		//Metadata tag
		dataTags, _ := GetTagSeoByID(listTagID, cacheActive)
		dataByte, _ = json.Marshal(dataTags)
		json.Unmarshal(dataByte, &TagsOutputObjects.Metadata.Tags)

		//Metadata seo
		dataTagsSeo := seo.FormatSeoByMutilTags(dataTags)
		TagsOutputObjects.Metadata.Seo = dataTagsSeo

		//Switch images platform
		platform := Platform(platforms)
		for k, vodData := range vodDataObjects {
			var ImagesMapping ImagesOutputObjectStruct
			// switch platform.Type {
			// case "web":
			// 	// ImagesMapping.Vod_thumb = BuildImage(vodData.Images.Web.Vod_thumb)
			// 	ImagesMapping.Thumbnail = BuildImage(vodData.Images.Web.Thumbnail)
			// case "smarttv":
			// 	ImagesMapping.Thumbnail = BuildImage(vodData.Images.Smarttv.Thumbnail)
			// case "app":
			// 	ImagesMapping.Thumbnail = BuildImage(vodData.Images.App.Thumbnail)
			// }
			switch platform.Type {
			case "web":
				ImagesMapping.Vod_thumb = BuildImage(vodData.Images.Web.Vod_thumb)
				ImagesMapping.Thumbnail = BuildImage(vodData.Images.Web.Vod_thumb)
			case "smarttv":
				ImagesMapping.Thumbnail = BuildImage(vodData.Images.Smarttv.Thumbnail)
			case "app":
				ImagesMapping.Thumbnail = BuildImage(vodData.Images.App.Thumbnail)
			}
			ImagesMapping = MappingImagesV4(platform.Type, ImagesMapping, vodData.Images, true)
			TagsOutputObjects.Items[k].Images = ImagesMapping
		}
	}

	total, _ := GetTotalTag(listTagID, platforms, cacheActive)
	//Pagination
	TagsOutputObjects.Metadata.Total = total
	TagsOutputObjects.Metadata.Limit = limit
	TagsOutputObjects.Metadata.Page = page

	// Handle SEO
	var contentArr []string
	for k , v := range TagsOutputObjects.Items {
		if k >= 5 {
			break
		}
		contentArr = append(contentArr, v.Title)
	}
	var contentStr string = strings.Join(contentArr , ", ")

	TagsOutputObjects.Metadata.Seo.Description = fmt.Sprintf(TAG_DESCRIPTION, TagsOutputObjects.Metadata.Seo.Title, contentStr)
	TagsOutputObjects.Metadata.Seo.Seo_text = TagsOutputObjects.Metadata.Seo.Description
	


	// Write Redis
	dataByte, _ := json.Marshal(TagsOutputObjects)
	mRedis.SetString(keyCache, string(dataByte), TTL_KVCACHE)

	return TagsOutputObjects, nil
}

func GetTotalTag(listTagID []string, platforms string, cacheActive bool) (int, error) {
	platform := Platform(platforms)
	var total int = 0
	var keyCache = TAGS_VOD_TOTAL + "_" + strings.Join(listTagID, "-") + "_" + platforms

	if cacheActive {
		total, err := mRedis.GetInt(keyCache)
		if err == nil {
			return total, nil
		}
	}

	// Connect DB
	session, db, err := GetCollection()
	if err != nil {
		return total, err
	}
	defer session.Close()

	var where = bson.M{
		"tags.id":   bson.M{"$in": listTagID},
		"platforms": platform.Id,
	}

	//Write cache
	total, _ = db.C(COLLECTION_VOD).Find(where).Count()
	mRedis.SetInt(keyCache, total, TTL_REDIS_LV1)
	return total, nil
}

func GetTagSeoByID(listTagIds []string, cacheActive bool) ([]TagObjectStruct, error) {
	var TagObjects []TagObjectStruct
	if len(listTagIds) < 0 {
		return TagObjects, errors.New("GetTagSeoByID: Empty data")
	}

	if cacheActive {
		dataRedis, err := mRedis.HMGet(REFIX_REDIS_TAGS_SEO_HASH, listTagIds)
		if err == nil && len(dataRedis) > 0 {
			for _, val := range dataRedis {
				if str, ok := val.(string); ok {
					var TagSeoObject TagObjectStruct
					err = json.Unmarshal([]byte(str), &TagSeoObject)
					if err != nil {
						continue
					}
					TagObjects = append(TagObjects, TagSeoObject)
				}
			}
		}
	}

	//Connect mysql
	db_mysql, err := ConnectMySQL()
	if err != nil {
		return TagObjects, err
	}
	defer db_mysql.Close()

	queryTagIds := strings.Join(listTagIds, "','")
	sqlRaw := fmt.Sprintf(`SELECT id, name, type, slug FROM tag WHERE id IN ('%s') ORDER BY created_at DESC`, queryTagIds)
	entityTags, err := db_mysql.Query(sqlRaw)
	if err != nil {
		return TagObjects, err
	}

	for entityTags.Next() {
		var TagOne TagObjectStruct
		var id, name, slug string
		var typeTag int

		err := entityTags.Scan(&id, &name, &typeTag, &slug)
		if err != nil {
			return TagObjects, err
		}

		TagOne.Id = id
		TagOne.Name = name
		TagOne.Type = GetTagType(typeTag)
		TagOne.Slug = slug
		TagObjects = append(TagObjects, TagOne)
	}

	if len(TagObjects) > 0 {
		for k, val := range TagObjects {
			dataTagsSeo := seo.FormatSeoByTag(val.Id, val.Slug, val.Name, 100 , "", cacheActive)
			TagObjects[k].Seo = dataTagsSeo

			// Write cache follow ID
			dataByte, _ := json.Marshal(TagObjects)
			mRedis.HSet(REFIX_REDIS_TAGS_SEO_HASH, val.Id, string(dataByte))
		}
	}
	return TagObjects, nil
}

func GetTagsCategorys(cacheActive bool) ([]TagInfoOutputObjectStruct, error) {
	var TagInfoOutputObjects []TagInfoOutputObjectStruct
	var keyCacheKV = KV_REDIS_CATEGORY_TAGS

	if cacheActive {
		valueCache, err := mRedisKV.GetString(keyCacheKV)
		if err == nil && valueCache != "" {
			err = json.Unmarshal([]byte(valueCache), &TagInfoOutputObjects)
			if err == nil {
				return TagInfoOutputObjects, nil
			}
		}
	}

	//Connect mysql
	db_mysql, err := ConnectMySQL()
	if err != nil {
		return TagInfoOutputObjects, err
	}
	defer db_mysql.Close()

	sqlRaw := fmt.Sprintf(`SELECT id, name, type, slug FROM tag WHERE type = 4 and status = 1 and cloud = 1`)
	entityTags, err := db_mysql.Query(sqlRaw)
	if err != nil {
		return TagInfoOutputObjects, err
	}

	for entityTags.Next() {
		var TagOne TagInfoOutputObjectStruct
		var id, name, slug string
		var typeTag int

		err := entityTags.Scan(&id, &name, &typeTag, &slug)
		if err != nil {
			return TagInfoOutputObjects, err
		}
		TagOne.Id = id
		TagOne.Name = name
		TagOne.Type = GetTagType(typeTag)
		TagOne.Slug = slug
		TagInfoOutputObjects = append(TagInfoOutputObjects, TagOne)
	}

	var TagOne TagInfoOutputObjectStruct
	TagOne.Id = "99999999-9999-9999-9999-999999999999"
	TagOne.Name = "Nghệ sĩ"
	TagOne.Type = "category"
	TagOne.Slug = "nghe-si"
	TagInfoOutputObjects = append(TagInfoOutputObjects, TagOne)

	// Write Redis
	dataByte, _ := json.Marshal(TagInfoOutputObjects)
	mRedisKV.SetString(keyCacheKV, string(dataByte), TTL_KVCACHE)

	return TagInfoOutputObjects, nil
}
