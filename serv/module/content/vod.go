package content

import (
	// "errors"
	// "fmt"
	// "strings"
	// "time"

	"strings"

	. "cm-v5/schema"
	. "cm-v5/serv/module"
	Package "cm-v5/serv/module/packages"
	rating "cm-v5/serv/module/rating"
	Subscription "cm-v5/serv/module/subscription"
	tracking "cm-v5/serv/module/tracking"
	vod "cm-v5/serv/module/vod"
	watchlater "cm-v5/serv/module/watchlater"
)

func GetInfoVod(contentId string, epsId string, platform string, cacheActive bool) (VodOptimizeOutputStruct, error) {
	platformInfo := Platform(platform)
	var VodObjOutput VodOptimizeOutputStruct

	var keyCache = KV_REDIS_VOD_V3 + "_" + platformInfo.Type + "_" + contentId + "_" + epsId

	if cacheActive {
		valueCache, err := mRedisKV.GetString(keyCache)
		if err == nil && valueCache != "" {
			err = json.Unmarshal([]byte(valueCache), &VodObjOutput)
			if err == nil {
				return VodObjOutput, nil
			}
		}
	}
	var listId []string = []string{contentId}
	if epsId != "" {
		listId = append(listId, epsId)

	}
	listDataVod, err := vod.GetVODByListID(listId, platformInfo.Id, 0, true)
	if err != nil {
		return VodObjOutput, err
	}
	var EpisodeObj EpisodeOptimizeOutputStruct
	for _, val := range listDataVod {
		var ImagesMapping ImagesOutputObjectStruct
		// switch images follow platform
		switch platformInfo.Type {
		case "web":
			ImagesMapping.Vod_thumb_big = BuildImage(val.Images.Web.Vod_thumb_big)
			ImagesMapping.Thumbnail = BuildImage(val.Images.Web.Vod_thumb)
		case "smarttv":
			ImagesMapping.Vod_thumb_big = BuildImage(val.Images.Smarttv.Banner)
			ImagesMapping.Thumbnail = BuildImage(val.Images.Smarttv.Thumbnail)
		case "app":
			ImagesMapping.Thumbnail = BuildImage(val.Images.App.Thumbnail)
			ImagesMapping.Vod_thumb_big = BuildImage(val.Images.App.Vod_thumb_big)
		}

		ImagesMapping = MappingImagesV4(platformInfo.Type, ImagesMapping, val.Images, false)
		if val.Id == contentId {

			VodObjOutput.Images = ImagesMapping
			//sync data
			dataByte, _ := json.Marshal(val)
			err = json.Unmarshal(dataByte, &VodObjOutput)
			if err != nil {
				return VodObjOutput, err
			}

		} else if val.Id == epsId && val.Type == VOD_TYPE_EPISODE {

			//sync data episode
			dataByte, _ := json.Marshal(val)
			err = json.Unmarshal(dataByte, &EpisodeObj)
			if err != nil {
				return VodObjOutput, err
			}
			EpisodeObj.Images = ImagesMapping
			// Check is end epi
			EpisodeObj.Is_end = CheckVodIsEndEps(EpisodeObj.Group_id, EpisodeObj.Id, platform)

		}
	}
	if EpisodeObj.Id != "" {
		VodObjOutput.Episode_info = EpisodeObj
	} else {
		VodObjOutput.Episode_info = ""
	}

	//get list packages
	VodObjOutput.PackageGroup = GetListPackageGroupByVOD(VodObjOutput)

	VodObjOutput.Is_end = true

	//Check trailer exists
	if VodObjOutput.Type == VOD_TYPE_SEASON {
		totalTrailler, _ := GetListIDRelatedVideoByZrange(contentId, 0, 1, true)
		if len(totalTrailler) > 0 {
			VodObjOutput.Have_trailer = 1
		}

		VodObjOutput.Related_season = GetListRelatedSeasonBySeason(VodObjOutput.Group_id, VodObjOutput.Id, cacheActive)
	}

	// Write Redis
	dataByte, _ := json.Marshal(VodObjOutput)
	mRedisKV.SetString(keyCache, string(dataByte), TTL_KVCACHE)

	return VodObjOutput, nil
}

/**
Lay cac thong tin content co lien quan den user
- rating
- permission
- watchlater
- progress
- .....
*/

func GetVodInfoPersonal(contentId string, epsId string, userId string, platform string, modelPlatform string, ipUser string, cacheActive bool) (VodDetailOptimizePersonalStruct, error) {
	var ContentDetailPersonalObj VodDetailOptimizePersonalStruct
	ContentDetailPersonalObj.Id = contentId
	ContentDetailPersonalObj.Group_id = contentId
	if epsId != "" {
		ContentDetailPersonalObj.Id = epsId
		ContentDetailPersonalObj.Group_id = contentId
	}

	// Check permission
	GetPermissionVodPersonal(&ContentDetailPersonalObj, userId, platform, modelPlatform)

	//if empty data episode field `ContentDetailPersonalObj.Episode_info` return ""
	if ContentDetailPersonalObj.Episode_info == nil {
		ContentDetailPersonalObj.Episode_info = ""
	}

	// Check link play (VinadataBuildTokenUrl)
	ContentDetailPersonalObj.Link_play.Dash_link_play = BuildTokenUrl(ContentDetailPersonalObj.Link_play.Dash_link_play, "", ipUser)
	ContentDetailPersonalObj.Link_play.Hls_link_play = BuildTokenUrl(ContentDetailPersonalObj.Link_play.Hls_link_play, "", ipUser)

	if len(ContentDetailPersonalObj.Subtitles) == 0 {
		ContentDetailPersonalObj.Subtitles = make([]SubtitlesOuputStruct, 0)
	}
	if len(ContentDetailPersonalObj.Audios) == 0 {
		ContentDetailPersonalObj.Audios = make([]AudiosOutputStruct, 0)
	}
	// if len(ContentDetailPersonalObj.PackageGroup) == 0 {
	// 	ContentDetailPersonalObj.PackageGroup = make([]PackageGroupObjectStruct, 0)
	// }

	if userId != "" {
		// Check Anonymous
		if strings.HasPrefix(userId, "anonymous_") == true {
			return ContentDetailPersonalObj, nil
		}

		//get status watch later
		var Episode_info EpisodeOptimizeOutputStruct
		dataByte, _ := json.Marshal(ContentDetailPersonalObj.Episode_info)
		json.Unmarshal(dataByte, &Episode_info)
		if Episode_info.Id != "" {
			Episode_info.Is_watchlater = watchlater.CheckContentIsWatchLater(Episode_info.Id, userId)
			ContentDetailPersonalObj.Episode_info = Episode_info
			ContentDetailPersonalObj.Is_watchlater = watchlater.CheckContentIsWatchLater(ContentDetailPersonalObj.Group_id, userId)
		} else {
			ContentDetailPersonalObj.Is_watchlater = watchlater.CheckContentIsWatchLater(ContentDetailPersonalObj.Id, userId)
		}

		// Check watch later (User - Detail)

		// Check progress (User - Detail)
		listProgress, err := tracking.GetProgressByListID([]string{ContentDetailPersonalObj.Id}, userId)
		if len(listProgress) > 0 && err == nil {
			ContentDetailPersonalObj.Progress = listProgress[ContentDetailPersonalObj.Id]
		}

		// Check rating (User - Season)
		ContentDetailPersonalObj.User_rating = rating.GetRatingContentByUser(userId, ContentDetailPersonalObj.Group_id)
	}

	return ContentDetailPersonalObj, nil
}

/**
0: 		Khong duoc xem (khong tra link play)
206: 	Duoc phep xem (co tra link play)
207: 	Chỉ user premium dược phép xem (co tra link play)
*/
func GetPermissionVodPersonal(ContentDetailPersonalObj *VodDetailOptimizePersonalStruct, userId string, platform string, modelPlatform string) {
	ContentDetailPersonalObj.Permission = permissionDefaultVod

	// Set Default Map_profile Premiun
	ContentDetailPersonalObj.Link_play.Map_profile.Full_hd.Is_premium = 1

	// Get content detail by id
	dataContent, err := GetFullInfoVodDetail(ContentDetailPersonalObj.Id, platform, true)
	if err != nil {
		Sentry_log(err)
		return
	}

	var ContentDetailPersonalObjTemp VodDetailOptimizePersonalStruct
	dataByte, _ := json.Marshal(dataContent)
	err = json.Unmarshal(dataByte, &ContentDetailPersonalObjTemp)
	if err != nil {
		Sentry_log(err)
		return
	}

	// Khong login quyen xem dua theo quyen Default cua Setting
	if strings.HasPrefix(userId, "anonymous_") == false && userId != "" {
		// Co login duoc phep xem
		ContentDetailPersonalObj.Permission = PERMISSION_VALID
	} else {
		//Process ads for non user
		ContentDetailPersonalObj.Ads = ProcessAdsForUser(userId, false, ContentDetailPersonalObjTemp.Ads, "")
	}

	// Check permission valid
	if ContentDetailPersonalObj.Permission == PERMISSION_REQUIRE_LOGIN {
		ContentDetailPersonalObj.Link_play.Dash_link_play = ""
		ContentDetailPersonalObj.Link_play.Hls_link_play = ""
		return
	}

	// Kiem tra user đã mua package nào thuộc list package trên chưa
	// Neu co permission valid
	var IsUserPremium = false
	var Sub Subscription.SubcriptionObjectStruct
	SubcriptionsOfUser, err := Sub.GetListByUserId(userId)
	if err == nil && len(SubcriptionsOfUser) > 0 {
		IsUserPremium = true
	}

	// Check content in package
	// Xem content co thuoc package nao khong
	// Neu khong permission theo permission default
	// Neu co permission theo permission mua goi / code
	var Pack Package.PackagesObjectStruct
	PackageGroupDenyContents, err := Pack.GetListPackageByDenyContent(ContentDetailPersonalObj.Group_id, true)
	if err == nil && len(PackageGroupDenyContents) > 0 {
		// Content co thuoc 1/n package nao do
		// Set permission theo permission mua goi / code
		if userId == "" {
			ContentDetailPersonalObj.Permission = PERMISSION_REQUIRE_LOGIN // 0
		} else {
			//Check content thuoc goi vip or free
			for _, val := range PackageGroupDenyContents {
				if val.Billing_package_group_id == 10 {
					ContentDetailPersonalObj.Permission = PERMISSION_REQUIRE_PREMIUM
					break
				}

				ContentDetailPersonalObj.Permission = PERMISSION_REQUIRE_PACKAGE
				break
			}
		}
	}
	for _, vSubObj := range SubcriptionsOfUser {
		for _, vPGDC := range PackageGroupDenyContents {
			if vSubObj.Current_package_id == vPGDC.Pk_id {
				ContentDetailPersonalObj.Permission = PERMISSION_VALID
				break
			}
		}
	}

	// Check permission valid
	if ContentDetailPersonalObj.Permission != PERMISSION_VALID {
		ContentDetailPersonalObj.Link_play.Dash_link_play = ""
		ContentDetailPersonalObj.Link_play.Hls_link_play = ""
		return
	}

	// Sync data
	ContentDetailPersonalObj.Link_play = ContentDetailPersonalObjTemp.Link_play
	ContentDetailPersonalObj.Subtitles = ContentDetailPersonalObjTemp.Subtitles
	ContentDetailPersonalObj.Audios = ContentDetailPersonalObjTemp.Audios
	ContentDetailPersonalObj.Drm_service_name = ContentDetailPersonalObjTemp.Drm_service_name
	// ContentDetailPersonalObj.PackageGroup = ContentDetailPersonalObjTemp.PackageGroup

	//Process ads for user
	if len(ContentDetailPersonalObj.Ads) == 0 {
		ContentDetailPersonalObj.Ads = ProcessAdsForUser(userId, IsUserPremium, ContentDetailPersonalObjTemp.Ads, "")
	}

	//copy intro outtro
	ContentDetailPersonalObj.Intro.Start = dataContent.Intro_start
	ContentDetailPersonalObj.Intro.End = dataContent.Intro_end
	ContentDetailPersonalObj.Outtro.Start = dataContent.Outtro_start
	ContentDetailPersonalObj.Outtro.End = dataContent.Outtro_end

	//Khanh DT-11014
	//Get link play from default episode if data request is SEASON

	var dataEpisode EpisodeOptimizeOutputStruct
	if dataContent.Type == VOD_TYPE_SEASON {
		//get episode data

		dataByte, _ := json.Marshal(dataContent.Default_episode)
		err := json.Unmarshal(dataByte, &dataEpisode)
		if err != nil {
			Sentry_log(err)
			return
		}

		//copy ads to episode
		dataEpisode.Ads = ContentDetailPersonalObj.Ads

		ContentDetailPersonalObj.Episode_info = dataEpisode

		//sync link play from episode
		dataByte, _ = json.Marshal(dataContent.Default_episode.Link_play)
		err = json.Unmarshal(dataByte, &ContentDetailPersonalObj.Link_play)
		if err != nil {
			Sentry_log(err)
			return
		}

	} else if dataContent.Type == VOD_TYPE_EPISODE || dataContent.Type == VOD_TYPE_TRAILER {
		dataByte, _ := json.Marshal(dataContent)
		err := json.Unmarshal(dataByte, &dataEpisode)
		if err != nil {
			Sentry_log(err)
			return
		}
		//copy ads to episode
		dataEpisode.Ads = ContentDetailPersonalObj.Ads

		ContentDetailPersonalObj.Episode_info = dataEpisode
	}

	//nếu là smarttv của samsung or lg thì lấy linkplay h265
	if CheckPlatformUsingLinkplayH265(platform, modelPlatform) {
		if ContentDetailPersonalObj.Link_play.H265_dash_link_play != "" {
			ContentDetailPersonalObj.Link_play.Dash_link_play = ContentDetailPersonalObj.Link_play.H265_dash_link_play
		}
		if ContentDetailPersonalObj.Link_play.H265_hls_link_play != "" {
			ContentDetailPersonalObj.Link_play.Hls_link_play = ContentDetailPersonalObj.Link_play.H265_hls_link_play
		}
	}
	ContentDetailPersonalObj.Link_play.H265_dash_link_play = ""
	ContentDetailPersonalObj.Link_play.H265_hls_link_play = ""

	// Handle link play
	ContentDetailPersonalObj.Link_play.Dash_link_play = HandleLinkPlayVieON(ContentDetailPersonalObj.Link_play.Dash_link_play)
	ContentDetailPersonalObj.Link_play.Hls_link_play = HandleLinkPlayVieON(ContentDetailPersonalObj.Link_play.Hls_link_play)

	ContentDetailPersonalObj.Link_play.Dash_link_play = HandleLinkPlayVodByPremium(ContentDetailPersonalObj.Link_play.Dash_link_play, IsUserPremium)
	ContentDetailPersonalObj.Link_play.Hls_link_play = HandleLinkPlayVodByPremium(ContentDetailPersonalObj.Link_play.Hls_link_play, IsUserPremium)

	// Set Permision Profile (Can handle them - config tu setting)
	ContentDetailPersonalObj.Link_play.Map_profile.Full_hd.Is_premium = 1
	ContentDetailPersonalObj.Link_play.Map_profile.Full_hd.Permission = ContentDetailPersonalObj.Permission
	ContentDetailPersonalObj.Link_play.Map_profile.Hd.Permission = ContentDetailPersonalObj.Permission
	ContentDetailPersonalObj.Link_play.Map_profile.Sd.Permission = ContentDetailPersonalObj.Permission

	// Set Permision Audio
	for k, val := range ContentDetailPersonalObj.Audios {
		ContentDetailPersonalObj.Audios[k].Is_premium = 1
		ContentDetailPersonalObj.Audios[k].Permission = ContentDetailPersonalObj.Permission
		if val.Is_default == 1 {
			ContentDetailPersonalObj.Audios[k].Is_premium = 0
			ContentDetailPersonalObj.Audios[k].Permission = PERMISSION_VALID
		}

	}

	// Set Permision Subtitle
	for k, val := range ContentDetailPersonalObj.Subtitles {
		ContentDetailPersonalObj.Subtitles[k].Is_premium = 1
		ContentDetailPersonalObj.Subtitles[k].Permission = ContentDetailPersonalObj.Permission
		if val.Is_default == 1 {
			ContentDetailPersonalObj.Subtitles[k].Is_premium = 0
			ContentDetailPersonalObj.Subtitles[k].Permission = PERMISSION_VALID
		}

	}

	return
}

func GetFullInfoVodDetail(contentId string, platform string, cacheActive bool) (ContentObjOutputStruct, error) {
	platformInfo := Platform(platform)
	var ContentObjOutput ContentObjOutputStruct
	var ImagesMapping ImagesOutputObjectStruct

	vodObj, err := vod.GetVodDetail(contentId, platformInfo.Id, cacheActive)
	if err != nil {
		return ContentObjOutput, err
	}

	dataByte, _ := json.Marshal(vodObj)
	err = json.Unmarshal(dataByte, &ContentObjOutput)
	if err != nil {
		return ContentObjOutput, err
	}

	//get ads
	ContentObjOutput.Ads = GetAds(vodObj.Content_provider_id, vodObj.Type, vodObj.Group_id, cacheActive, platform)

	//get list packages
	ContentObjOutput = GetPermissionDefaultVOD(ContentObjOutput)

	if ContentObjOutput.Type == VOD_TYPE_SEASON {
		// Call func GetVodByGroup get group VOD
		Episode_info := vod.GetDefaultEpisodeByGroupId(ContentObjOutput.Id, platformInfo.Id)

		dataByte, _ = json.Marshal(Episode_info)
		err = json.Unmarshal(dataByte, &ContentObjOutput.Default_episode)
		if err != nil {
			return ContentObjOutput, err
		}
		var ImagesMappingEpisode ImagesOutputObjectStruct
		// switch images follow platform
		switch platformInfo.Type {
		case "web":
			ImagesMappingEpisode.Vod_thumb_big = BuildImage(Episode_info.Images.Web.Vod_thumb_big)
			ImagesMappingEpisode.Thumbnail = BuildImage(Episode_info.Images.Web.Vod_thumb)
		case "smarttv":
			ImagesMappingEpisode.Vod_thumb_big = BuildImage(Episode_info.Images.Smarttv.Vod_thumb_big)
			ImagesMappingEpisode.Thumbnail = BuildImage(Episode_info.Images.Smarttv.Thumbnail)
		case "app":
			ImagesMappingEpisode.Thumbnail = BuildImage(Episode_info.Images.App.Thumbnail)
			ImagesMappingEpisode.Vod_thumb_big = BuildImage(Episode_info.Images.Web.Vod_thumb_big)
		}

		ImagesMappingEpisode = MappingImagesV4(platformInfo.Type, ImagesMappingEpisode, Episode_info.Images, false)
		ContentObjOutput.Default_episode.Images = ImagesMappingEpisode

		// Check is end epi
		ContentObjOutput.Default_episode.Is_end = CheckVodIsEndEps(ContentObjOutput.Default_episode.Group_id, ContentObjOutput.Default_episode.Id, platform)
	}

	ContentObjOutput.Is_end = true
	if ContentObjOutput.Type == VOD_TYPE_EPISODE {
		// Check is end epi
		ContentObjOutput.Is_end = CheckVodIsEndEps(ContentObjOutput.Group_id, ContentObjOutput.Id, platform)
	}

	//Check trailer exists
	if ContentObjOutput.Type == VOD_TYPE_SEASON {
		totalTrailler, _ := GetListIDRelatedVideoByZrange(contentId, 0, 1, true)
		if len(totalTrailler) > 0 {
			ContentObjOutput.Have_trailer = 1
		}
	}

	// switch images follow platform
	switch platformInfo.Type {
	case "web":
		ImagesMapping.Vod_thumb_big = BuildImage(vodObj.Images.Web.Vod_thumb_big)
		ImagesMapping.Thumbnail = BuildImage(vodObj.Images.Web.Vod_thumb)
	case "smarttv":
		ImagesMapping.Vod_thumb_big = BuildImage(vodObj.Images.Smarttv.Banner)
		ImagesMapping.Thumbnail = BuildImage(vodObj.Images.Smarttv.Thumbnail)
	case "app":
		ImagesMapping.Thumbnail = BuildImage(vodObj.Images.App.Thumbnail)
		ImagesMapping.Vod_thumb_big = BuildImage(vodObj.Images.App.Vod_thumb_big)
	}
	ImagesMapping = MappingImagesV4(platformInfo.Type, ImagesMapping, vodObj.Images, false)

	ContentObjOutput.Images = ImagesMapping

	return ContentObjOutput, nil
}

func GetListPackageGroupByVOD(vodObj VodOptimizeOutputStruct) []PackageGroupObjectStruct {
	var ListPackageGroup = make([]PackageGroupObjectStruct, 0)

	// get package
	contentID := vodObj.Id
	if vodObj.Type == VOD_TYPE_EPISODE {
		contentID = vodObj.Group_id
	}

	var keyCache = KV_VOD_LIST_PACKAGE_GROUP + "_" + contentID
	valC, err := mRedisKV.GetString(keyCache)
	if err == nil {
		dataByte, _ := json.Marshal(valC)
		err = json.Unmarshal([]byte(dataByte), &ListPackageGroup)
		if err == nil {
			return ListPackageGroup
		}
	} else {
		//Connect mysql
		db_mysql, _ := ConnectMySQL()
		defer db_mysql.Close()

		// sqlRaw := fmt.Sprintf(`
		// 	SELECT bpg.id, bpg.name
		// 	FROM billing_packages_deny_contents as bpdc
		// 	LEFT JOIN billing_package_group as bpg ON  bpdc.billing_package_group_id = bpg.id
		// 	WHERE bpdc.content_id = "%s" AND bpg.is_active = 1`, contentID)
		dataPackage, err := db_mysql.Query(`
		SELECT bpg.id, bpg.name 
		FROM billing_packages_deny_contents as bpdc
		LEFT JOIN billing_package_group as bpg ON  bpdc.billing_package_group_id = bpg.id
		WHERE bpdc.content_id = ? AND bpg.is_active = 1`, contentID)
		if err == nil {
			for dataPackage.Next() {
				var PackageObject PackageGroupObjectStruct
				err := dataPackage.Scan(&PackageObject.Id, &PackageObject.Name)
				if err != nil {
					continue
				}
				PackageObject.Price = 30000
				PackageObject.Period = "month"
				PackageObject.Period_value = 1
				ListPackageGroup = append(ListPackageGroup, PackageObject)
			}

			// Write cache
			dataByte, _ := json.Marshal(ListPackageGroup)
			mRedisKV.SetString(keyCache, string(dataByte), TTL_KVCACHE*2)
		}
	}

	return ListPackageGroup
}
