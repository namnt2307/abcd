package artist

import (
	"fmt"
	"strings"
	. "cm-v5/serv/module"
	. "cm-v5/schema"
	vod "cm-v5/serv/module/vod"
	seo "cm-v5/serv/module/seo"
)

func GetVodArtist(peopleID, platform string, page, limit, sort int, cacheActive bool) (VodArtistOutputObjectStruct, error) {
	var VodArtistOutputObject VodArtistOutputObjectStruct
	VodArtistOutputObject.Items = make([]ItemsVodArtistOutputStruct, 0)
	platformInfo := Platform(platform)

	var keyCache = ARTIST_VOD + "_" + fmt.Sprint(page) + "_" + fmt.Sprint(limit) + "_" + fmt.Sprint(sort) + "_" + peopleID + "_" + platformInfo.Type

	if cacheActive {
		dataCache, err := mRedis.GetString(keyCache)
		if err == nil && dataCache != "" {
			err = json.Unmarshal([]byte(dataCache), &VodArtistOutputObject)
			if err == nil {
				return VodArtistOutputObject, nil
			}
		}
	}

	//Kiem tra thong tin nghe si va lay thong tin SEO nghe si
	dataArtist, err := GetInfoArtist(peopleID, cacheActive)
	if err != nil || dataArtist.Id == "" {
		fmt.Println(err)
		return VodArtistOutputObject, err
	}

	//Page limit
	VodArtistOutputObject.Metadata.Page = page
	VodArtistOutputObject.Metadata.Limit = limit

	//Lay danh sach ID Vod cua nghe si
	listVodID, err := GetVodIdsArtist(peopleID, platform, page, limit, sort, cacheActive)
	if err != nil || len(listVodID) <= 0 {
		return VodArtistOutputObject, err
	}

	//Lay total vod
	total, _ := GetTotalVodArtist(peopleID, platformInfo.Type, cacheActive)
	VodArtistOutputObject.Metadata.Total = total

	//Lay danh sach Vod thuoc nghe si
	vodObjects, err := vod.GetVODByListID(listVodID, platformInfo.Id, 1, cacheActive)
	if err != nil || len(vodObjects) <= 0 {
		return VodArtistOutputObject, nil
	}

	dataByte, _ := json.Marshal(vodObjects)
	err = json.Unmarshal(dataByte, &VodArtistOutputObject.Items)
	if err != nil {
		return VodArtistOutputObject, err
	}

	//Switch images platform
	for k, vodData := range vodObjects {
		var ImagesMapping ImagesOutputObjectStruct
		// switch platformInfo.Type {
		// case "web":
		// 	// ImagesMapping.Vod_thumb = BuildImage(vodData.Images.Web.Vod_thumb)
		// 	ImagesMapping.Thumbnail = BuildImage(vodData.Images.Web.Vod_thumb)
		// case "smarttv":
		// 	// ImagesMapping.Vod_thumb_big = BuildImage(vodData.Images.Smarttv.Vod_thumb_big)
		// 	// ImagesMapping.Home_vod_hot = BuildImage(vodData.Images.Smarttv.Home_vod_hot)
		// 	ImagesMapping.Thumbnail = BuildImage(vodData.Images.Smarttv.Thumbnail)
		// case "app":
		// 	// ImagesMapping.Vod_thumb_big = BuildImage(vodData.Images.App.Vod_thumb_big)
		// 	// ImagesMapping.Home_vod_hot = BuildImage(vodData.Images.App.Home_vod_hot)
		// 	ImagesMapping.Thumbnail = BuildImage(vodData.Images.App.Thumbnail)
		// }
		switch platformInfo.Type {
		case "web":
			ImagesMapping.Vod_thumb = BuildImage(vodData.Images.Web.Vod_thumb)
			ImagesMapping.Thumbnail = BuildImage(vodData.Images.Web.Vod_thumb)
		case "smarttv":
			ImagesMapping.Vod_thumb_big = BuildImage(vodData.Images.Smarttv.Vod_thumb_big)
			ImagesMapping.Home_vod_hot = BuildImage(vodData.Images.Smarttv.Home_vod_hot)
			ImagesMapping.Thumbnail = BuildImage(vodData.Images.Smarttv.Thumbnail)
		case "app":
			ImagesMapping.Vod_thumb_big = BuildImage(vodData.Images.App.Vod_thumb_big)
			ImagesMapping.Home_vod_hot = BuildImage(vodData.Images.App.Home_vod_hot)
			ImagesMapping.Thumbnail = BuildImage(vodData.Images.App.Thumbnail)
		}

		ImagesMapping = MappingImagesV4(platformInfo.Type, ImagesMapping, vodData.Images, true)
		VodArtistOutputObject.Items[k].Images = ImagesMapping
	}


	var ContentArtistArr []string
	for k , val := range VodArtistOutputObject.Items {
		if k >= 5 {
			break
		}
		ContentArtistArr = append(ContentArtistArr , val.Title)
	}
	var ContentArtistStr string = strings.Join(ContentArtistArr , ", ")
	VodArtistOutputObject.Seo = seo.FormatSeoArtist(dataArtist, total , ContentArtistStr)


	// Write cache
	dataByte, _ = json.Marshal(VodArtistOutputObject)
	mRedis.SetString(keyCache, string(dataByte), TTL_REDIS_LV1)

	return VodArtistOutputObject, nil
}

func GetTotalVodArtist(peopleID string, platform string, cacheActive bool) (int, error) {
	var keyCache = ARTIST_VOD_TOTAL + "_" + peopleID + "_" + platform
	var total int
	if cacheActive {
		total, err := mRedis.GetInt(keyCache)
		if err == nil {
			return total, nil
		}
	}

	//Connect mysql
	db_mysql, err := ConnectMySQL()
	if err != nil {
		return total, err
	}
	defer db_mysql.Close()

	platformInfo := Platform(platform)
	// sqlRaw := fmt.Sprintf(`
	// 	SELECT count(DISTINCT (a.entity_id)) FROM entity_people as a
	// 	LEFT JOIN entity_vod as b ON b.id = a.entity_id
	// 	LEFT JOIN entity_vod_platform as c ON c.entity_id = b.id
	// 	WHERE a.people_id = '%s' AND b.status IN (3, 5) AND b.type IN (1,3,4,5)
	// 	AND c.platform_id = %s`, peopleID, fmt.Sprint(platformInfo.Id))

	totalVodArtist, err := db_mysql.Query(`
	SELECT count(DISTINCT (a.entity_id)) FROM entity_people as a 
	LEFT JOIN entity_vod as b ON b.id = a.entity_id
	LEFT JOIN entity_vod_platform as c ON c.entity_id = b.id
	WHERE a.people_id = ? AND b.status IN (3, 5) AND b.type IN (1,3,4,5) 
	AND c.platform_id = ?`, peopleID, fmt.Sprint(platformInfo.Id))
	if err != nil {
		return total, err
	}

	for totalVodArtist.Next() {
		err := totalVodArtist.Scan(&total)
		if err != nil {
			return total, err
		}
	}

	mRedis.SetInt(keyCache, total, TTL_REDIS_LV1)
	return total, nil
}
