package content

import (
	"errors"
	"strings"
	"time"
	. "cm-v5/schema"
	. "cm-v5/serv/module"
	seo "cm-v5/serv/module/seo"
	tracking "cm-v5/serv/module/tracking"
	vod "cm-v5/serv/module/vod"
	watchlater "cm-v5/serv/module/watchlater"
)

var permissionDefaultVod = 0

func init() {
	permissionDefaultVod, _ = CommonConfig.GetInt("GLOBAL_CONFIG", "permission_default")
}

func GetContent(contentId string, platform string, cacheActive bool) (ContentObjOutputStruct, error) {
	platformInfo := Platform(platform)
	var ContentObjOutput ContentObjOutputStruct
	var ImagesMapping ImagesOutputObjectStruct
	var keyCache = KV_REDIS_CONTENT + "_" + platform + "_" + contentId

	if cacheActive {
		valueCache, err := mRedisKV.GetString(keyCache)
		if err == nil && valueCache != "" {
			err = json.Unmarshal([]byte(valueCache), &ContentObjOutput)
			if err == nil {
				return ContentObjOutput, nil
			}
		}
	}

	// Call func GetVodDetail get detail vod
	vodObj, err := vod.GetVodDetail(contentId, platformInfo.Id, cacheActive)
	if err != nil {
		mRedisKV.Del(keyCache)
		return ContentObjOutput, err
	}
	if vodObj.Id == "" {
		return ContentObjOutput, errors.New("GetVodDetail: Empty data - " + contentId)
	}

	dataByte, _ := json.Marshal(vodObj)
	err = json.Unmarshal(dataByte, &ContentObjOutput)
	if err != nil {
		return ContentObjOutput, err
	}

	ContentObjOutput.Min_rate = vodObj.Min_rate
	timeParse, _ := time.Parse("2006-01-02", vodObj.Release_date)
	ContentObjOutput.Created_at = timeParse.Unix()

	if ContentObjOutput.Type == VOD_TYPE_SEASON {
		// Call func GetVodByGroup get group VOD
		defaultEpisode := vod.GetDefaultEpisodeByGroupId(ContentObjOutput.Id, platformInfo.Id)

		dataByte, _ = json.Marshal(defaultEpisode)
		err = json.Unmarshal(dataByte, &ContentObjOutput.Default_episode)
		if err != nil {
			return ContentObjOutput, err
		}
		var ImagesMappingDefaultEps ImagesOutputObjectStruct
		// switch images follow platform
		switch platformInfo.Type {
		case "web":
			ImagesMappingDefaultEps.Vod_thumb_big = BuildImage(defaultEpisode.Images.Web.Vod_thumb_big)
			ImagesMappingDefaultEps.Vod_thumb = BuildImage(defaultEpisode.Images.Web.Vod_thumb)
			ImagesMappingDefaultEps.Thumbnail = BuildImage(defaultEpisode.Images.Web.Vod_thumb)
		case "smarttv":
			ImagesMappingDefaultEps.Banner = BuildImage(defaultEpisode.Images.Smarttv.Banner)
			ImagesMappingDefaultEps.Thumbnail = BuildImage(defaultEpisode.Images.Smarttv.Thumbnail)
		case "app":
			ImagesMappingDefaultEps.Thumbnail = BuildImage(defaultEpisode.Images.App.Thumbnail)
			ImagesMappingDefaultEps.Vod_thumb_big = BuildImage(defaultEpisode.Images.Web.Vod_thumb_big)
		}
		ImagesMappingDefaultEps = MappingImagesV4(platformInfo.Type, ImagesMappingDefaultEps, defaultEpisode.Images, false)
		ContentObjOutput.Default_episode.Images = ImagesMappingDefaultEps

		// Check is end epi
		ContentObjOutput.Default_episode.Is_end = CheckVodIsEndEps(ContentObjOutput.Default_episode.Group_id, ContentObjOutput.Default_episode.Id, platform)

		// get related season by season
		ContentObjOutput.Related_season = GetListRelatedSeasonBySeason(ContentObjOutput.Group_id, ContentObjOutput.Id, cacheActive)
	}

	ContentObjOutput.Is_end = true
	if ContentObjOutput.Type == VOD_TYPE_EPISODE {
		
		// Check is end epi
		ContentObjOutput.Is_end = CheckVodIsEndEps(ContentObjOutput.Group_id, ContentObjOutput.Id, platform)

		// Get more info seasion
		vodObjSeasion, _ := vod.GetVodDetail(ContentObjOutput.Group_id, platformInfo.Id, cacheActive)

		dataByte, _ = json.Marshal(vodObjSeasion)
		json.Unmarshal(dataByte, &ContentObjOutput.Movie)
	}

	//Check trailer exists
	if ContentObjOutput.Type == VOD_TYPE_SEASON || ContentObjOutput.Type == VOD_TYPE_MOVIE {
		totalTrailler, _ := GetListIDRelatedVideoByZrange(contentId, 0, 1, true)
		if len(totalTrailler) > 0 {
			ContentObjOutput.Have_trailer = 1
		}
	}

	// switch images follow platform
	switch platformInfo.Type {
	case "web":
		ImagesMapping.Vod_thumb_big = BuildImage(vodObj.Images.Web.Vod_thumb_big)
		ImagesMapping.Vod_thumb = BuildImage(vodObj.Images.Web.Vod_thumb)
		ImagesMapping.Thumbnail = BuildImage(vodObj.Images.Web.Vod_thumb)
		if ImagesMapping.Vod_thumb_big != "" {
			ImagesMapping.Vod_thumb = ImagesMapping.Vod_thumb_big
			ImagesMapping.Thumbnail = ImagesMapping.Vod_thumb_big
		}
	case "smarttv":
		ImagesMapping.Banner = BuildImage(vodObj.Images.Smarttv.Banner)
		ImagesMapping.Thumbnail = BuildImage(vodObj.Images.Smarttv.Thumbnail)
	case "app":
		ImagesMapping.Thumbnail = BuildImage(vodObj.Images.App.Thumbnail)
		ImagesMapping.Vod_thumb_big = BuildImage(vodObj.Images.App.Vod_thumb_big)
		if ImagesMapping.Vod_thumb_big != "" {
			ImagesMapping.Thumbnail = ImagesMapping.Vod_thumb_big
		}
	}
	ImagesMapping = MappingImagesV4(platformInfo.Type, ImagesMapping, vodObj.Images, false)

	ContentObjOutput.Images = ImagesMapping

	//Check permission
	dataContentObjOutput := GetPermissionDefaultVOD(ContentObjOutput)

	// Add code Ads Thư
	dataContentObjOutput.Ads = GetAds(dataContentObjOutput.Content_provider_id, dataContentObjOutput.Type, dataContentObjOutput.Group_id, cacheActive, platform)
	dataContentObjOutput.Default_episode.Ads = dataContentObjOutput.Ads

	// Handle show name
	if dataContentObjOutput.Show_name == "" {
		Show_nameTemp := strings.Split(dataContentObjOutput.Title, " - ")
		if len(Show_nameTemp) > 0 {
			dataContentObjOutput.Show_name = Show_nameTemp[0]
		}
	}

	// Hanle tag slug
	for k, val := range dataContentObjOutput.Tags {
		dataContentObjOutput.Tags[k].Seo = seo.FormatSeoByTag(val.Id, val.Slug, val.Name, 100, "", cacheActive)
	}

	// handle SEO info VOD
	dataContentObjOutput.Seo = seo.FormatSeoVODDetail(dataContentObjOutput, cacheActive)

	// Write Redis
	dataByte, _ = json.Marshal(dataContentObjOutput)
	mRedisKV.SetString(keyCache, string(dataByte), TTL_KVCACHE)

	return dataContentObjOutput, nil
}

func GetAds(idProvider string, typeContent int, groupId string, cache bool, platform string) []AdsOutputStruct {
	var arrAds []AdsOutputStruct
	var arrIdProvider []string
	if typeContent == VOD_TYPE_EPISODE || typeContent == VOD_TYPE_TRAILER {
		// GET mysql
		idProvider = GetContentProviderIdMySQL(groupId)
	}
	// Convert string idProvider to array
	arrIdProvider = strings.Split(idProvider, ",")
	if len(arrIdProvider) > 0 {
		arrAds = GetAdsMySQL(arrIdProvider[0], arrAds, cache, platform)
	}

	return arrAds
}

func GetContentProviderIdMySQL(ContentId string) string {
	var arrIdProvider string
	vodObj, err := vod.GetVodDetail(ContentId, 0, true)
	if err == nil {
		arrIdProvider = vodObj.Content_provider_id
	}
	return arrIdProvider
}

func GetUserInfoDataContent(userId string, dataContent ContentObjOutputStruct, ipUser string) (ContentObjOutputStruct, error) {
	// Check localtion
	if dataContent.Geo_check == 1 {
		// Content nay  Chi duoc xem o VN
		if CheckIpIsVN(ipUser) == false {
			// IP khong phai o VN
			return dataContent, errors.New("Content not validate your country")
		}
	}

	// if STATUS_LISTEPISODE_GETUSERINFO == "true" {
	// 	// Rating info
	// 	RatingDetailContent, err := rating.RatingInfoByContentId(dataContent.Id, "web")
	// 	if err == nil {
	// 		dataContent.Avg_rate = RatingDetailContent.Avg_rate
	// 		dataContent.Total_rate = RatingDetailContent.Total_rate
	// 	}
	// }

	// handle rating
	dataContent.Avg_rate = ReCalculatorRating(dataContent.Min_rate, dataContent.Avg_rate)
	dataContent.Min_rate = 0

	// get Permission
	dataContent = GetPermissionVOD(dataContent, userId)

	// Check index of range eps
	if dataContent.Type == VOD_TYPE_SEASON {
		// Check index of range default_episode
		dataContent.Range_page_index = vod.GetPageIndexZRangeByGroupIdAndEpsID(dataContent.Id, dataContent.Default_episode.Id)
	} else if dataContent.Type == VOD_TYPE_EPISODE {
		// Check index of range eps
		dataContent.Range_page_index = vod.GetPageIndexZRangeByGroupIdAndEpsID(dataContent.Group_id, dataContent.Id)
	}

	//Replace hinh nghe si
	for k, Director := range dataContent.People.Director {
		dataContent.People.Director[k].Images.Avatar = HandleLinkPlayVieON(Director.Images.Avatar)
	}
	for k, Actor := range dataContent.People.Actor {
		dataContent.People.Actor[k].Images.Avatar = HandleLinkPlayVieON(Actor.Images.Avatar)
	}

	if strings.HasPrefix(userId, "anonymous_") == true || userId == "" {
		return dataContent, nil
	}

	// Check watchlater
	// dataContent.Is_watchlater = watchlater.CheckContentIsWatchLater(dataContent.Id, userId)

	// Check progress
	if dataContent.Type == VOD_TYPE_SEASON {
		if STATUS_LISTEPISODE_GETUSERINFO == "true" {
			// check default_episode exist
			var contentIdDefault = tracking.GetIdDefaultEpsBySessionID(dataContent.Id, userId)
			if contentIdDefault != "" {
				// Get info Default_episode
				VODDataObject, err := vod.GetVodDetail(contentIdDefault, 0, true)
				if err == nil {
					dataByte, _ := json.Marshal(VODDataObject)
					json.Unmarshal(dataByte, &dataContent.Default_episode)
					dataContent.Default_episode.Is_watchlater = watchlater.CheckContentIsWatchLater(contentIdDefault, userId)
				}
			}
		}
	}

	// Check rating
	// dataContent.User_rating = rating.GetRatingContentByUser(userId, dataContent.Id)

	return dataContent, nil
}
