package content

import (
	// "fmt"

	"errors"
	"fmt"
	"strconv"
	"strings"

	. "cm-v5/serv/module"
	. "cm-v5/schema"
	seo "cm-v5/serv/module/seo"
	tracking "cm-v5/serv/module/tracking"
	vod "cm-v5/serv/module/vod"
	watchlater "cm-v5/serv/module/watchlater"
	"gopkg.in/mgo.v2/bson"
)

func GetEpisode(contentId string, entitySlug string, page int, limit int, platforms string, cacheActive bool) (EpisodeObjOutputStruct, error) {
	platform := Platform(platforms)
	var EpisodeObjOutput EpisodeObjOutputStruct
	var ItemsEpisodeObjOutput = make([]ItemsEpisodeObjOutputStruct, 0)
	EpisodeObjOutput.Items = ItemsEpisodeObjOutput

	if entitySlug != "" {

		vodId, err := vod.GetVodIdBySlug(entitySlug, platform.Id)
		if err != nil {
			return EpisodeObjOutput, err
		}
		if vodId == "" {
			return EpisodeObjOutput, errors.New("not found")
		}

		contentId = vodId
	}

	var keyCache = KV_REDIS_LIST_EPISODE + "_" + platforms + "_" + contentId + "_" + strconv.Itoa(page) + "_" + strconv.Itoa(limit)
	if cacheActive {
		valueCache, err := mRedisKV.GetString(keyCache)
		if err == nil && valueCache != "" {
			err = json.Unmarshal([]byte(valueCache), &EpisodeObjOutput)
			if err == nil {
				return EpisodeObjOutput, nil
			}
		}
	}

	// Call func GetVodByListGroup get group VOD
	listEpisode, err := vod.GetVodByListGroup(contentId, platform.Id, page, limit, cacheActive)
	if err != nil {
		return EpisodeObjOutput, err
	}

	var total int = 0

	dataContent, _ := GetContent(contentId, platforms, true)
	if len(listEpisode) > 0 {
		//Return format sort
		// sort.Slice(listEpisode, func(i, j int) bool { return listEpisode[i].Odr < listEpisode[j].Odr })
		dataByte, _ := json.Marshal(listEpisode)
		err = json.Unmarshal(dataByte, &EpisodeObjOutput.Items)
		if err != nil {
			return EpisodeObjOutput, err
		}
		//Switch images platform
		for k, vodData := range listEpisode {
			var ImagesMapping ImagesOutputObjectStruct
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
			EpisodeObjOutput.Items[k].Images = ImagesMapping
			EpisodeObjOutput.Items[k].Resolution = dataContent.Resolution
		}
	}
	total, _ = GetTotalEpisodeVOD(contentId, platforms, cacheActive)
	//Pagination
	EpisodeObjOutput.Metadata.Total = total
	EpisodeObjOutput.Metadata.Limit = limit
	EpisodeObjOutput.Metadata.Page = page

	//get data by content id to create seo info
	EpisodeObjOutput.Seo = seo.FormatSeoEpisodeList(dataContent.Title, dataContent.Seo.Url)

	// Write Redis
	dataByte, _ := json.Marshal(EpisodeObjOutput)
	mRedisKV.SetString(keyCache, string(dataByte), TTL_KVCACHE)

	return EpisodeObjOutput, nil
}

func GetGroupIdBySlugSeason(entitySlugSeason, platforms string, cacheActive bool) (string, error) {
	platform := Platform(platforms)
	var keyCache = EPISODE_GROUP_ID_BY_SLUG + "_" + entitySlugSeason + "_" + strconv.Itoa(platform.Id)

	if cacheActive {
		value, err := mRedis.GetString(keyCache)
		if err == nil {
			return value, nil
		}
	}
	// Connect DB
	session, db, err := GetCollection()
	if err != nil {
		return "", err
	}
	defer session.Close()

	var where = bson.M{
		"slug_seo": entitySlugSeason,
	}

	var VODDataObject VODDataObjectStruct
	err = db.C(COLLECTION_VOD).Find(where).One(&VODDataObject)
	if err != nil && err.Error() != "not found" {
		return "", err
	}

	// Write cache
	mRedis.SetString(keyCache, VODDataObject.Id, TTL_REDIS_LV1)

	return VODDataObject.Id, nil
}

func GetTotalEpisodeVOD(groupId, platforms string, cacheActive bool) (int, error) {

	var keyCache = CONTENT_LIST_EPS + groupId
	total := mRedis.ZCountAll(keyCache)
	return int(total), nil

	// platform := Platform(platforms)
	// var keyCache = EPISODE_VOD_TOTAL + "_" + groupId + "_" + strconv.Itoa(platform.Id)
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
	// 	"type":      4,
	// 	"platforms": platform.Id,
	// }

	// //Get total episode
	// total, _ := db.C(COLLECTION_VOD).Find(where).Count()
	// mRedis.SetInt(keyCache, total, TTL_REDIS_LV1)

	// return total, nil
}

func GetUserInfoDataEpisode(userId string, dataEpisode EpisodeObjOutputStruct) EpisodeObjOutputStruct {
	if STATUS_LISTEPISODE_GETUSERINFO == "true" {
		return dataEpisode
	}
	// Get Watchlater
	var listIdContent []string
	for k, val := range dataEpisode.Items {
		dataEpisode.Items[k].Is_watchlater = watchlater.CheckContentIsWatchLater(val.Id, userId)
		listIdContent = append(listIdContent, val.Id)
	}

	// get progress eps
	if len(listIdContent) > 0 {
		listProgress, err := tracking.GetProgressByListID(listIdContent, userId)
		if len(listProgress) > 0 && err == nil {
			for k, val := range dataEpisode.Items {
				if progress, ok := listProgress[val.Id]; ok {
					dataEpisode.Items[k].Progress = progress
				}
			}
		}
	}
	return dataEpisode
}

func GetRangeEpisode(groupId, platform string, cacheActive bool) []string {
	var rangeEpisode = make([]string, 0)
	var platformInfo = Platform(platform)

	total, _ := GetTotalEpisodeVOD(groupId, platform, cacheActive)
	if total > 0 {
		var totalPage float64 = float64(total) / float64(LIMIT_EPS_PER_RANGE)
		for page := 0.0; page < totalPage; page++ {
			listEpisode, err := vod.GetVodByListGroup(groupId, platformInfo.Id, int(page), LIMIT_EPS_PER_RANGE, cacheActive)
			if err != nil {
				return rangeEpisode
			}

			// Get Title Start - End
			if len(listEpisode) > 0 {
				var startTitle = listEpisode[0].Title
				var endTitle = listEpisode[len(listEpisode)-1].Episode
				var rangeTitle string = startTitle
				if startTitle != "Tập "+fmt.Sprint(endTitle) {
					rangeTitle = startTitle + " - " + fmt.Sprint(endTitle)
				}
				rangeEpisode = append(rangeEpisode, rangeTitle)
			}
		}
	}
	return rangeEpisode
}

func CheckVodIsEndEps(groupId, contentId, platform string) bool {
	var platformInfo = Platform(platform)
	total, _ := GetTotalEpisodeVOD(groupId, platform, true)
	if total > 0 {
		var totalPage float64 = float64(total) / float64(LIMIT_EPS_PER_RANGE)
		listEpisode, err := vod.GetVodByListGroup(groupId, platformInfo.Id, int(totalPage), LIMIT_EPS_PER_RANGE, true)
		if err != nil {
			return false
		}
		if len(listEpisode) > 0 {
			endEpisodeID := listEpisode[0].Id
			endEpisodeNum := listEpisode[0].Episode
			endEpisodeSlugSeo := listEpisode[0].Slug_seo
			for i := 1; i < len(listEpisode); i++ {
				//nếu tập lớn hơn hoặc nếu số tập bằng và slug_seo lớn hơn item ghi nhận trước đó thì ghi nhận item mới này
				//strings.Compare(a,b)) = 1 if a > b
				if listEpisode[i].Episode > endEpisodeNum ||
					(listEpisode[i].Episode == endEpisodeNum && strings.Compare(listEpisode[i].Slug_seo, endEpisodeSlugSeo) == 1) {
					endEpisodeID = listEpisode[i].Id
					endEpisodeNum = listEpisode[i].Episode
					endEpisodeSlugSeo = listEpisode[i].Slug_seo
				}
			}
			if endEpisodeID == contentId {
				return true
			}
		}
	}
	return false
}
