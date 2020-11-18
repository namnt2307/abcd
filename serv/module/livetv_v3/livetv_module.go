package livetv_v3

import (
	"errors"
	"fmt"
	"strings"

	. "cm-v5/schema"
	. "cm-v5/serv/module"
	"cm-v5/serv/module/seo"

	"gopkg.in/mgo.v2/bson"
)

func GetLiveTVByListID(listId []string, platform int, cache bool) ([]LiveTVMongoObjectStruct, error) {
	var LiveTVMongoObjects []LiveTVMongoObjectStruct
	if len(listId) <= 0 {
		return LiveTVMongoObjects, errors.New("GetVODByListID: Empty data")
	}

	if cache {
		dataRedis, err := mRedis.HMGet(PREFIX_REDIS_HASH_LIVE_TV, listId)
		if err == nil && len(dataRedis) > 0 {
			for _, val := range dataRedis {
				if str, ok := val.(string); ok {
					var LiveTVMongoObject LiveTVMongoObjectStruct
					err = json.Unmarshal([]byte(str), &LiveTVMongoObject)
					if err != nil {
						continue
					}
					LiveTVMongoObjects = append(LiveTVMongoObjects, LiveTVMongoObject)
				}
			}
		}
	}

	if len(LiveTVMongoObjects) != len(listId) {
		// Connect DB
		session, db, err := GetCollection()
		if err != nil {
			Sentry_log(err)
			return LiveTVMongoObjects, err
		}
		defer session.Close()

		var where = bson.M{
			"id":   bson.M{"$in": listId},
			"type": 5,
		}

		err = db.C(COLLECTION_LIVE_TV).Find(where).All(&LiveTVMongoObjects)
		if err != nil {
			Sentry_log(err)
			return LiveTVMongoObjects, err
		}

		if len(LiveTVMongoObjects) <= 0 {
			// Remove cache
			for _, id := range listId {
				mRedis.HDel(PREFIX_REDIS_HASH_LIVE_TV, id)
			}
			fmt.Println("where", where)
			return LiveTVMongoObjects, errors.New("GetLiveTVByListID: Empty data")
		}
		// Write cache
		for k, val := range LiveTVMongoObjects {
			//handle image
			LiveTVMongoObjects[k].Image_link = BuildImage(val.Image_link)

			//handle seo object
			LiveTVMongoObjects[k].Seo = seo.FormatSeoLiveTV(val.Id, val.Slug, val.Title, cache)

			dataByte, _ := json.Marshal(LiveTVMongoObjects[k])
			// Write cache follow ID
			mRedis.HSet(PREFIX_REDIS_HASH_LIVE_TV, val.Id, string(dataByte))
			mRedis.HSet(PREFIX_REDIS_HASH_LIVE_TV_SLUG, LiveTVMongoObjects[k].Seo.Url, string(val.Id))

			//set drm key link to livetv have drm k+ or castlab
			drmServiceName := strings.ToLower(val.Drm_service_name)
			if drmServiceName == "k+" || drmServiceName == "castlab" {
				keyCache := KV_REDIS_LIVETV_DRM_KEY + "_" + val.Drm_key
				mRedisKV.SetString(keyCache, val.Id, -1)
			}

		}
	}

	// Check platform
	var LiveTVMongoObjecsTemp []LiveTVMongoObjectStruct
	for _, valVod := range LiveTVMongoObjects {
		exists, _ := In_array(platform, valVod.Platforms)
		if exists || platform == 0 {
			LiveTVMongoObjecsTemp = append(LiveTVMongoObjecsTemp, valVod)
		}
	}

	// Sap xep mang theo thu tu id truyen vao
	var LiveTVMongoObjecsSort []LiveTVMongoObjectStruct
	for _, valId := range listId {
		for _, vodDataSort := range LiveTVMongoObjecsTemp {
			if valId == vodDataSort.Id {
				// Crypt
				// encrypted, _ := AESEncrypt(vodDataSort.Link_play)
				// vodDataSort.Link_play = encrypted

				LiveTVMongoObjecsSort = append(LiveTVMongoObjecsSort, vodDataSort)
				continue
			}
		}
	}

	return LiveTVMongoObjecsSort, nil
}

func GetIdLiveTVByAssetID(assetID string) string {
	keyCache := KV_REDIS_LIVETV_DRM_KEY + "_" + assetID
	data, err := mRedisKV.GetString(keyCache)
	if err != nil {
		return ""
	}
	return data
}
