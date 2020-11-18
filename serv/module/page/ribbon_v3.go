package page

import (
	// "errors"
	. "cm-v5/schema"
	. "cm-v5/serv/module"
	recommendation "cm-v5/serv/module/recommendation"
	seo "cm-v5/serv/module/seo"
	vod "cm-v5/serv/module/vod"
	"fmt"
	"strings"

	"gopkg.in/mgo.v2/bson"
)

func GetRibbonIdBySlugV3(ribbonSlug string) (string, error) {
	var keyCacheSlug = "DETAIL_RIB_SLUG_V3_" + ribbonSlug
	value, err := mRedis.GetString(keyCacheSlug)
	if err == nil && value != "" {
		return value, nil
	}

	var ribbonId string = ""

	if len(ribbonSlug) > 10 {
		ribbonSlug = ribbonSlug[12 : len(ribbonSlug)-1]
	}

	// Check slug validate
	// slugSplits := strings.Split(ribbonSlug, fmt.Sprintf(SEO_RIBBON_URL, ""))
	// if len(slugSplits) != 2 {
	// 	return ribbonId, errors.New("Slug not validated")
	// }
	// ribbonSlug = slugSplits[1]
	var RibbonsV3 RibbonsV3ObjectStruct

	// Get data in DB
	session, db, err := GetCollection()
	if err != nil {
		Sentry_log(err)
		return ribbonId, err
	}
	defer session.Close()
	var where = bson.M{
		"status": 1,
		"$or": []bson.M{
			bson.M{"slug": ribbonSlug},
			bson.M{"slug_filter": ribbonSlug},
		},
	}

	err = db.C(COLLECTION_RIB_V3).Find(where).One(&RibbonsV3)
	if err != nil && err.Error() != "not found" {
		Sentry_log(err)
		return ribbonId, err
	}

	ribbonId = RibbonsV3.Id

	keyCacheSlug = "DETAIL_RIB_SLUG_V3_" + RibbonsV3.Slug
	mRedis.SetString(keyCacheSlug, ribbonId, TTL_REDIS_LV1)

	if RibbonsV3.Slug_filter != "" {
		var keyCacheSlugFilter = "DETAIL_RIB_SLUG_V3_" + RibbonsV3.Slug_filter
		mRedis.SetString(keyCacheSlugFilter, ribbonId, TTL_REDIS_LV1)
	}
	return ribbonId, err
}

func GetRibbonInfoV3(rib_id string, platform string, page, limit, order int, cacheActive bool) (RibbonDetailOutputObjectStruct, error) {
	var RibbonDetailOutputV3Object RibbonDetailOutputObjectStruct
	var keyCache = KV_DETAIL_RIBBON_INFO_V3 + "_" + platform + "_" + rib_id + "_" + fmt.Sprint(page) + "_" + fmt.Sprint(limit) + "_" + fmt.Sprint(order)

	if cacheActive {
		valueCache, err := mRedisKV.GetString(keyCache)
		if err == nil && valueCache != "" {
			err = json.Unmarshal([]byte(valueCache), &RibbonDetailOutputV3Object)
			if err == nil {
				return RibbonDetailOutputV3Object, nil
			}
		}
	}

	var platformInfo = Platform(platform)

	//get info ribbon v3
	ribbonInfo, err := GetRibbonInfoByRibIDV3(rib_id, platform, cacheActive)
	if err != nil {
		return RibbonDetailOutputV3Object, err
	}
	dataByte, _ := json.Marshal(ribbonInfo)
	json.Unmarshal(dataByte, &RibbonDetailOutputV3Object)

	//mapping images ribbon
	var ImagesMapping ImagesOutputObjectStruct
	ImagesMapping = MappingImagesV4(platformInfo.Type, ImagesMapping, ribbonInfo.Images, false)
	RibbonDetailOutputV3Object.Images = ImagesMapping

	//get ribbon item
	var shortcutMapping = make([]VODDataObjectStruct, 0)
	shortcutMapping, err = GetListRibbonItemV3ByID(RibbonDetailOutputV3Object.Id, platform, page, limit, cacheActive)
	if err != nil {
		Sentry_log(err)
		return RibbonDetailOutputV3Object, err
	}

	dataByte, _ = json.Marshal(shortcutMapping)

	err = json.Unmarshal(dataByte, &RibbonDetailOutputV3Object.Ribbon_items)
	if err != nil {
		Sentry_log(err)
		return RibbonDetailOutputV3Object, err
	}
	//Format field image follow platform
	for k, vodData := range shortcutMapping {
		var ImagesMapping ImagesOutputObjectStruct
		ImagesMapping = MappingImagesV4(platformInfo.Type, ImagesMapping, vodData.Images, k != 0)

		RibbonDetailOutputV3Object.Ribbon_items[k].Images = ImagesMapping
	}

	// Mapping properties
	switch platformInfo.Type {
	case "web":
		RibbonDetailOutputV3Object.Properties.Line = ribbonInfo.Properties.Web.Line
		RibbonDetailOutputV3Object.Properties.Is_title = ribbonInfo.Properties.Web.Is_title
		RibbonDetailOutputV3Object.Properties.Is_slide = ribbonInfo.Properties.Web.Is_slide
		RibbonDetailOutputV3Object.Properties.Is_refresh = ribbonInfo.Properties.Web.Is_refresh
		RibbonDetailOutputV3Object.Properties.Is_view_all = ribbonInfo.Properties.Web.Is_view_all
		RibbonDetailOutputV3Object.Properties.Thumb = ribbonInfo.Properties.Web.Thumb
	case "smarttv":
		RibbonDetailOutputV3Object.Properties.Line = ribbonInfo.Properties.Smarttv.Line
		RibbonDetailOutputV3Object.Properties.Is_title = ribbonInfo.Properties.Smarttv.Is_title
		RibbonDetailOutputV3Object.Properties.Is_slide = ribbonInfo.Properties.Smarttv.Is_slide
		RibbonDetailOutputV3Object.Properties.Is_refresh = ribbonInfo.Properties.Smarttv.Is_refresh
		RibbonDetailOutputV3Object.Properties.Is_view_all = ribbonInfo.Properties.Smarttv.Is_view_all
		RibbonDetailOutputV3Object.Properties.Thumb = ribbonInfo.Properties.Smarttv.Thumb
	case "app":
		RibbonDetailOutputV3Object.Properties.Line = ribbonInfo.Properties.App.Line
		RibbonDetailOutputV3Object.Properties.Is_title = ribbonInfo.Properties.App.Is_title
		RibbonDetailOutputV3Object.Properties.Is_slide = ribbonInfo.Properties.App.Is_slide
		RibbonDetailOutputV3Object.Properties.Is_refresh = ribbonInfo.Properties.App.Is_refresh
		RibbonDetailOutputV3Object.Properties.Is_view_all = ribbonInfo.Properties.App.Is_view_all
		RibbonDetailOutputV3Object.Properties.Thumb = ribbonInfo.Properties.App.Thumb
	}

	//Get total
	total := mRedis.ZCountAll(LIST_RIB_ITEM_V3 + "_" + rib_id + "_" + platform)
	RibbonDetailOutputV3Object.Metadata.Total = int(total)
	RibbonDetailOutputV3Object.Metadata.Limit = limit
	RibbonDetailOutputV3Object.Metadata.Page = page

	// Handle new SEO
	var contentArr []string
	for k, v := range RibbonDetailOutputV3Object.Ribbon_items {
		if k >= 5 {
			break
		}
		contentArr = append(contentArr, v.Title)
	}

	var contentStr string = strings.Join(contentArr, ", ")
	RibbonDetailOutputV3Object.Seo = seo.FormatSeoRibbon(RibbonDetailOutputV3Object.Id, RibbonDetailOutputV3Object.Seo.Slug, RibbonDetailOutputV3Object.Name, int(total), contentStr, cacheActive)

	// Write Redis
	dataByte, _ = json.Marshal(RibbonDetailOutputV3Object)
	mRedisKV.SetString(keyCache, string(dataByte), TTL_KVCACHE)
	return RibbonDetailOutputV3Object, nil
}

func GetRibbonInfoByRibIDV3(rib_id, platform string, cacheActive bool) (RibbonsV3ObjectStruct, error) {
	var RibbonDetailOutputV3Object RibbonsV3ObjectStruct

	var keyCache = "DETAIL_RIB_" + platform + "_" + rib_id
	if cacheActive {
		value, err := mRedis.GetString(keyCache)
		if err == nil && value != "" {
			err = json.Unmarshal([]byte(value), &RibbonDetailOutputV3Object)
			if err == nil {
				return RibbonDetailOutputV3Object, nil
			}
		}
	}

	// Connect MongoDB
	session, db, err := GetCollection()
	if err != nil {
		return RibbonDetailOutputV3Object, err
	}

	defer session.Close()

	var platformInfo = Platform(platform)
	var where = bson.M{
		"id":        rib_id,
		"platforms": platformInfo.Id,
		"status":    1,
		"type":      bson.M{"$in": []int{0, 1, 2, 3, 4, 5, 6, 7, 8, 9, 10, 11, 12}},
	}

	err = db.C(COLLECTION_RIB_V3).Find(where).One(&RibbonDetailOutputV3Object)
	if err != nil && err.Error() != "not found" {
		return RibbonDetailOutputV3Object, err
	}

	RibbonDetailOutputV3Object.Seo = seo.FormatSeoRibbon(RibbonDetailOutputV3Object.Id, RibbonDetailOutputV3Object.Slug, RibbonDetailOutputV3Object.Name, 10, "", cacheActive)

	dataByte, _ := json.Marshal(RibbonDetailOutputV3Object)
	// Write cache
	mRedis.SetString(keyCache, string(dataByte), 0)
	return RibbonDetailOutputV3Object, nil
}

func GetRibbonByYUSP(user_id, device_id, scenario, rib_id, platform string, page, limit int, cacheActive bool) (RibbonDetailOutputObjectStruct, error) {
	var RibbonDetailOutputObject RibbonDetailOutputObjectStruct

	var keyCache = "ribbon_yusp_" + device_id + "_" + user_id + "_" + scenario + "_" + platform + "_" + rib_id + "_" + fmt.Sprint(page) + "_" + fmt.Sprint(limit)
	if cacheActive {
		// valueCache, err := mRedis.GetString(keyCache)
		// if err == nil && valueCache != "" {
		// 	err = json.Unmarshal([]byte(valueCache), &RibbonDetailOutputObject)
		// 	if err == nil {
		// 		return RibbonDetailOutputObject, nil
		// 	}
		// }
	}

	platformInfo := Platform(platform)

	//get info ribbon v3
	ribbonInfo, err := GetRibbonInfoByRibIDV3(rib_id, platform, cacheActive)
	if err != nil {
		return RibbonDetailOutputObject, err
	}
	dataByte, _ := json.Marshal(ribbonInfo)
	json.Unmarshal(dataByte, &RibbonDetailOutputObject)

	//mapping images ribbon
	var ImagesMapping ImagesOutputObjectStruct
	ImagesMapping = MappingImagesV4(platformInfo.Type, ImagesMapping, ribbonInfo.Images, true)
	RibbonDetailOutputObject.Images = ImagesMapping

	dataGetYusp, err := recommendation.GetCommonYuspRecommendationAPI(user_id, device_id, scenario, platform, page, limit)
	if err != nil {
		Sentry_log(err)
		return RibbonDetailOutputObject, err
	}

	if scenario == "BYWATCHED" && len(dataGetYusp.Items) > 0 {
		ribName := RibbonDetailOutputObject.Name
		prefixReplace := "{movie_name}"
		firstItemFromYusp := dataGetYusp.Items[0]

		if strings.Index(ribName, prefixReplace) > -1 {
			RibbonDetailOutputObject.Name = strings.Replace(ribName, prefixReplace, strings.ToUpper(firstItemFromYusp.Title), 1)
		} else {
			RibbonDetailOutputObject.Name = ribName + " " + strings.ToUpper(firstItemFromYusp.Title)
		}

		RibbonDetailOutputObject.Seo = seo.FormatSeoRibbon(ribbonInfo.Id, ribbonInfo.Slug, ribbonInfo.Name, 10, "", cacheActive)

		//remove first item from yusp
		dataGetYusp.ItemIds = dataGetYusp.ItemIds[1:]
		dataGetYusp.TotalResults = len(dataGetYusp.ItemIds)
	}

	// Get List VOD Info by List VOD ID
	var VODDataObjectTemp []VODDataObjectStruct
	if len(dataGetYusp.ItemIds) > 0 {
		VODDataObjectTemp, err = vod.GetVODByListID(dataGetYusp.ItemIds, platformInfo.Id, 1, true)
		if err != nil || len(VODDataObjectTemp) <= 0 {
			Sentry_log(err)
			return RibbonDetailOutputObject, err
		}
	}

	dataByte, err = json.Marshal(VODDataObjectTemp)
	if err != nil {
		Sentry_log(err)
		return RibbonDetailOutputObject, err
	}

	err = json.Unmarshal(dataByte, &RibbonDetailOutputObject.Ribbon_items)
	if err != nil {
		Sentry_log(err)
		return RibbonDetailOutputObject, err
	}

	//Format field image follow platform
	for k, vodData := range VODDataObjectTemp {
		var ImagesMapping ImagesOutputObjectStruct

		// switch images follow platform
		var urlPoster string
		for _, val := range vodData.Image_soucre {
			if val.Image_type == "poster" {
				urlPoster = val.Url
				break
			}
		}

		// switch images follow platform
		switch platformInfo.Type {
		case "web":
			ImagesMapping.Home_carousel_web = BuildImage(vodData.Images.Web.Home_carousel_web)
			ImagesMapping.Vod_thumb = BuildImage(vodData.Images.Web.Vod_thumb)
			ImagesMapping.Thumbnail = BuildImage(vodData.Images.Web.Vod_thumb)
			ImagesMapping.Home_vod_hot = BuildImage(vodData.Images.Web.Home_vod_hot)
			ImagesMapping.Poster = BuildImage(vodData.Images.Web.Poster)
			if ImagesMapping.Poster == "" {
				ImagesMapping.Poster = BuildImage(SupportUrlImg(urlPoster, "350_466"))
			}

		case "smarttv":
			ImagesMapping.Home_carousel_tv = BuildImage(vodData.Images.Smarttv.Home_carousel_tv)
			ImagesMapping.Thumbnail = BuildImage(vodData.Images.Smarttv.Vod_thumb)
			ImagesMapping.Poster = BuildImage(vodData.Images.Smarttv.Poster)
			if ImagesMapping.Poster == "" {
				ImagesMapping.Poster = BuildImage(SupportUrlImg(urlPoster, "350_466"))
			}
		case "app":
			ImagesMapping.Banner = BuildImage(vodData.Images.App.Banner)
			ImagesMapping.Thumbnail = BuildImage(vodData.Images.App.Vod_thumb)
			ImagesMapping.Poster = BuildImage(vodData.Images.App.Poster)
			if ImagesMapping.Poster == "" {
				ImagesMapping.Poster = BuildImage(SupportUrlImg(urlPoster, "140_187"))
			}
		}

		ImagesMapping = MappingImagesV4(platformInfo.Type, ImagesMapping, vodData.Images, true)

		// Output
		RibbonDetailOutputObject.Ribbon_items[k].Images = ImagesMapping
	}

	RibbonDetailOutputObject.Tracking_data = recommendation.GetTrackingData(dataGetYusp.RecommendationId, RECOMMENDATION_NAME_YUSP)
	RibbonDetailOutputObject.Metadata.Total = dataGetYusp.TotalResults
	RibbonDetailOutputObject.Metadata.Limit = limit
	RibbonDetailOutputObject.Metadata.Page = page

	//Set cache data YUSP
	dataByte, _ = json.Marshal(RibbonDetailOutputObject)
	mRedis.SetString(keyCache, string(dataByte), TTL_REDIS_6_HOURS)

	return RibbonDetailOutputObject, nil
}

//no cache
func GetRibbonByMenuIDNoCache(menu_id string) (RibbonsV3ObjectStruct, error) {
	var RibbonDetailOutputV3Object RibbonsV3ObjectStruct

	// Connect MongoDB
	session, db, err := GetCollection()
	if err != nil {
		return RibbonDetailOutputV3Object, err
	}

	defer session.Close()

	var where = bson.M{
		"menus.id": menu_id,
	}

	err = db.C(COLLECTION_RIB_V3).Find(where).One(&RibbonDetailOutputV3Object)
	if err != nil && err.Error() != "not found" {
		return RibbonDetailOutputV3Object, err
	}

	return RibbonDetailOutputV3Object, nil
}
