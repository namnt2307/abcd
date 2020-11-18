package page

import (
	"fmt"

	. "cm-v5/serv/module"
	. "cm-v5/schema"
	"gopkg.in/mgo.v2/bson"
)

func GetInfoPageBannersV3(page_id string, platform string, limit int, cacheActive bool) ([]BannerRibbonsOutputObjectStruct, error) {
	var BannerRibbonsOutputV3 = make([]BannerRibbonsOutputObjectStruct, 0)
	var keyCache = KV_DETAIL_RIBBONS_BANNERS_V3 + "_" + platform + "_" + page_id + fmt.Sprint(limit)

	if cacheActive {
		valueCache, err := mRedisKV.GetString(keyCache)
		if err == nil && valueCache != "" {
			err = json.Unmarshal([]byte(valueCache), &BannerRibbonsOutputV3)
			if err == nil {
				return BannerRibbonsOutputV3, nil
			}
		}
	}

	RibbonBannerV3, err := GetRibbonBannerV3(page_id, platform, cacheActive)
	if err != nil {
		Sentry_log(err)
		return BannerRibbonsOutputV3, err
	}

	var VodDataMapping = make([]VODDataObjectStruct, 0)
	VodDataMapping, err = GetListRibbonItemV3ByID(RibbonBannerV3.Id, platform, 0, limit, true)

	dataByte, _ := json.Marshal(VodDataMapping)
	err = json.Unmarshal(dataByte, &BannerRibbonsOutputV3)
	if err != nil {
		Sentry_log(err)
		return BannerRibbonsOutputV3, err
	}

	var platformInfo = Platform(platform)

	//Format field image follow platform
	for k, vodData := range VodDataMapping {
		var ImagesMapping ImagesOutputObjectStruct

		switch platformInfo.Type {
		case "web":
			ImagesMapping.Home_carousel_web = BuildImage(vodData.Images.Web.Home_carousel_web)
			ImagesMapping.Vod_thumb = BuildImage(vodData.Images.Web.Vod_thumb)
			ImagesMapping.Thumbnail = BuildImage(vodData.Images.Web.Vod_thumb)
			ImagesMapping.Vod_thumb_big = BuildImage(vodData.Images.Web.Home_vod_hot)
			ImagesMapping.Home_carousel = BuildImage(vodData.Images.Web.Home_carousel_web)
		case "smarttv":
			ImagesMapping.Home_carousel_tv = BuildImage(vodData.Images.Smarttv.Home_carousel_tv)
			ImagesMapping.Home_carousel = BuildImage(vodData.Images.Smarttv.Home_carousel_tv)
			ImagesMapping.Thumbnail = BuildImage(vodData.Images.Smarttv.Vod_thumb)
		case "app":
			ImagesMapping.Banner = BuildImage(vodData.Images.App.Banner)
			ImagesMapping.Thumbnail = BuildImage(vodData.Images.App.Banner)
		}
		ImagesMapping = MappingImagesV4(platformInfo.Type, ImagesMapping, vodData.Images, false)
		BannerRibbonsOutputV3[k].Images = ImagesMapping
	}

	// Write Redis
	dataByte, _ = json.Marshal(BannerRibbonsOutputV3)
	mRedisKV.SetString(keyCache, string(dataByte), TTL_KVCACHE_1_MINUTE)

	return BannerRibbonsOutputV3, nil
}

func GetRibbonBannerV3(page_id, platform string, cacheActive bool) (RibbonsV3ObjectStruct, error) {
	var RibbonsV3 RibbonsV3ObjectStruct
	var keyCache = LIST_RIBBONS_BANNER_V3 + "_" + platform + "_" + page_id

	if cacheActive {
		value, err := mRedis.GetString(keyCache)
		if err == nil && value != "" {
			err = json.Unmarshal([]byte(value), &RibbonsV3)
			if err == nil {
				return RibbonsV3, nil
			}
		}
	}

	var platformInfo = Platform(platform)

	// Connect MongoDB
	session, db, err := GetCollection()
	if err != nil {
		Sentry_log(err)
		return RibbonsV3, err
	}
	defer session.Close()

	var where = bson.M{
		"menus.id":        page_id,
		"menus.platforms": platformInfo.Id,
		"platforms":       platformInfo.Id,
		"status":          1,
		"type":            0,
	}

	err = db.C(COLLECTION_RIB_V3).Find(where).Sort("odr").One(&RibbonsV3)
	if err != nil && err.Error() != "not found" {
		mRedis.Del(keyCache)
		Sentry_log(err)
		return RibbonsV3, err
	}

	// Write cache
	dataByte, _ := json.Marshal(RibbonsV3)
	mRedis.SetString(keyCache, string(dataByte), TTL_REDIS_LV1)
	return RibbonsV3, nil
}
