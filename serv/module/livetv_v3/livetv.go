package livetv_v3

import (
	"errors"
	"fmt"
	"strings"

	. "cm-v5/schema"
	. "cm-v5/serv/module"
	"cm-v5/serv/module/packages"
	seo "cm-v5/serv/module/seo"

	"gopkg.in/mgo.v2/bson"
)

func GetIdBySlug(livetv_slug string) (string, error) {

	dataRedisBySlugId, _ := mRedis.HMGet(PREFIX_REDIS_HASH_LIVE_TV_SLUG, []string{livetv_slug})
	for _, val := range dataRedisBySlugId {
		if str, ok := val.(string); ok {
			return str, nil
		}
	}

	return "", errors.New("GetIdBySlug: Empty data")
}

func GetDetailLiveTV(livetv_id, epg_slug, platform string, cache bool) (DetailLiveTVObjectOutputStruct, error) {
	var DetailLiveTVObjectOutput DetailLiveTVObjectOutputStruct
	var listID = []string{livetv_id}
	var keyCache = KV_REDIS_DETAIL_LIVE_TV + "_" + livetv_id + "_" + epg_slug + "_" + platform

	if cache {
		valueCache, err := mRedisKV.GetString(keyCache)
		if err == nil && valueCache != "" {
			err = json.Unmarshal([]byte(valueCache), &DetailLiveTVObjectOutput)
			if err == nil {
				return DetailLiveTVObjectOutput, nil
			}
		}
	}

	platformInfo := Platform(platform)
	dataLiveTVTemp, err := GetLiveTVByListID(listID, platformInfo.Id, cache)
	if err != nil || len(dataLiveTVTemp) <= 0 {
		Sentry_log(err)
		return DetailLiveTVObjectOutput, err
	}

	dataByte, _ := json.Marshal(dataLiveTVTemp[0])
	err = json.Unmarshal(dataByte, &DetailLiveTVObjectOutput)
	if err != nil {
		Sentry_log(err)
		return DetailLiveTVObjectOutput, err
	}

	//Lay EPG hiện tại
	if epg_slug == "" {
		DetailLiveTVObjectOutput.Programme, _ = GetEPGCurrent(DetailLiveTVObjectOutput.Id)
		DetailLiveTVObjectOutput.Seo = seo.FormatSeoLiveTV(DetailLiveTVObjectOutput.Id, DetailLiveTVObjectOutput.Slug, DetailLiveTVObjectOutput.Title, cache)
	} else {
		DetailLiveTVObjectOutput.Programme, _ = GetDetailEpgBySlug(epg_slug, cache)
		DetailLiveTVObjectOutput.Programme.Seo = seo.FormatSeoEPG(DetailLiveTVObjectOutput.Slug, DetailLiveTVObjectOutput.Programme.Slug, DetailLiveTVObjectOutput.Title, DetailLiveTVObjectOutput.Programme.Title)
	}

	DetailLiveTVObjectOutput.Programme.Hls_link_play = GetLinkPlayEpg(DetailLiveTVObjectOutput.Hls_link_play, DetailLiveTVObjectOutput.Programme.Time_start, DetailLiveTVObjectOutput.Programme.Duration)

	for _, val := range dataLiveTVTemp[0].Livetv_group {
		DetailLiveTVObjectOutput.Livetv_group = append(DetailLiveTVObjectOutput.Livetv_group, val.Id)
	}

	//handle link player_logo
	if DetailLiveTVObjectOutput.Player_logo != "" {
		DetailLiveTVObjectOutput.Player_logo = BuildImage(DetailLiveTVObjectOutput.Player_logo)
	}

	//handle asset_id nếu dùng drm castlab hoặc K+
	drmName := strings.ToLower(DetailLiveTVObjectOutput.Drm_service_name)
	if drmName == "castlab" || drmName == "k+" {
		DetailLiveTVObjectOutput.Asset_id = dataLiveTVTemp[0].Drm_key
	}

	//Lay thong tin goi + permission
	var Pack packages.PackagesObjectStruct
	DetailLiveTVObjectOutput.PackageGroup, DetailLiveTVObjectOutput.Permission = Pack.GetPackage(DetailLiveTVObjectOutput.Id)

	// Write cache
	dataByte, _ = json.Marshal(DetailLiveTVObjectOutput)
	mRedisKV.SetString(keyCache, string(dataByte), TTL_KVCACHE)

	return DetailLiveTVObjectOutput, nil
}

func GetLiveTVByGroupId(livetv_group_id, platform string, page, limit int, cache bool) (LiveTVObjectOutputStruct, error) {
	var LiveTVObjectOutput LiveTVObjectOutputStruct
	LiveTVObjectOutput.Items = make([]LiveTVObjectStruct, 0)
	var keyCache = KV_REDIS_LIVE_TV_BY_GROUP + "_" + livetv_group_id + "_" + platform + "_" + fmt.Sprint(page) + "_" + fmt.Sprint(limit)

	if cache {
		valueCache, err := mRedisKV.GetString(keyCache)
		if err == nil && valueCache != "" {
			err = json.Unmarshal([]byte(valueCache), &LiveTVObjectOutput)
			if err == nil {
				return LiveTVObjectOutput, nil
			}
		}
	}

	listLiveTVId, err := GetListLiveTVIdByGroupId(livetv_group_id, platform, page, limit, cache)
	if err != nil || len(listLiveTVId) <= 0 {
		Sentry_log(err)
		return LiveTVObjectOutput, err
	}

	platformInfo := Platform(platform)
	dataLiveTVTemp, err := GetLiveTVByListID(listLiveTVId, platformInfo.Id, cache)
	if err != nil || len(dataLiveTVTemp) <= 0 {
		Sentry_log(err)
		return LiveTVObjectOutput, err
	}

	dataByte, _ := json.Marshal(dataLiveTVTemp)
	err = json.Unmarshal(dataByte, &LiveTVObjectOutput.Items)
	if err != nil {
		Sentry_log(err)
		return LiveTVObjectOutput, err
	}

	// for k, val := range dataLiveTVTemp {
	// 	for _, valG := range val.Livetv_group {
	// 		LiveTVObjectOutput.Items[k].Livetv_group = append(LiveTVObjectOutput.Items[k].Livetv_group, valG.Id)
	// 	}

	// }

	LiveTVObjectOutput.Metadata.Total, _ = GetTotalLiveTvByGroup(livetv_group_id, platform, cache)
	LiveTVObjectOutput.Metadata.Page = page
	LiveTVObjectOutput.Metadata.Limit = limit

	// Write cache
	dataByte, _ = json.Marshal(LiveTVObjectOutput)
	mRedisKV.SetString(keyCache, string(dataByte), TTL_KVCACHE)
	return LiveTVObjectOutput, nil
}

func GetListLiveTVIdByGroupId(livetv_group_id, platform string, page, limit int, cache bool) ([]string, error) {

	var keyCache = PREFIX_REDIS_LIST_ID_LIVE_TV_ID_BY_GROUP + "_" + livetv_group_id + "_" + platform
	var start = page * limit
	var stop = (start + limit) - 1

	if cache {
		vals, err := mRedis.ZRange(keyCache, int64(start), int64(stop))
		if err == nil && len(vals) > 0 {
			return vals, nil
		}
	}

	var listId []string
	// Connect DB
	session, db, err := GetCollection()
	if err != nil {
		Sentry_log(err)
		return listId, err
	}
	defer session.Close()

	platformInfo := Platform(platform)
	var where = bson.M{
		"livetv_group.id": livetv_group_id,
		"platforms":       platformInfo.Id,
		"type":            5,
	}

	var LiveTVMongoObjects []LiveTVMongoObjectStruct
	err = db.C(COLLECTION_LIVE_TV).Find(where).Sort("livetv_group.odr").Skip(page * limit).Limit(limit).All(&LiveTVMongoObjects)
	if err != nil && err.Error() != "not found" {
		Sentry_log(err)
		return listId, err
	}

	var listLiveTVIdAll []string

	//Write cache
	mRedis.Del(keyCache)
	fmt.Println("clear cache", keyCache)
	for _, val := range LiveTVMongoObjects {
		for _, valG := range val.Livetv_group {
			if livetv_group_id == valG.Id {
				listLiveTVIdAll = append(listLiveTVIdAll, val.Id)
				mRedis.ZAdd(keyCache, float64(valG.Odr), val.Id)
				continue
			}
		}

	}

	for i := start; i <= stop; i++ {
		if i < len(listLiveTVIdAll) {
			listId = append(listId, listLiveTVIdAll[i])
		}
	}

	return listId, nil
}

func GetTotalLiveTvByGroup(livetv_group_id string, platform string, cache bool) (int, error) {
	platformInfo := Platform(platform)
	var keyCache = TOTAL_LIVE_TV_BY_GROUP + "_" + livetv_group_id + "_" + platform
	if cache {
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
		"livetv_group.id": livetv_group_id,
		"platforms":       platformInfo.Id,
		"type":            5,
	}

	total, _ := db.C(COLLECTION_LIVE_TV).Find(where).Count()
	mRedis.SetInt(keyCache, total, TTL_REDIS_LV1)

	return total, nil
}
