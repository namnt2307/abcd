package content

import (
	. "cm-v5/serv/module"
	. "cm-v5/schema"
	vod "cm-v5/serv/module/vod"
	"gopkg.in/mgo.v2/bson"
)

func GetListRelatedSeasonBySeason(group_id, ss_id string, cache bool) []RelatedSeasonObjStruct {
	var lstRelatedSS []RelatedSeasonObjStruct

	listID, _ := GetListIDSeasonByGroupID(group_id, cache)

	if len(listID) > 0 {
		ResultVodData, _ := vod.GetVODByListID(listID, 0, 1, true)
		for _, val := range ResultVodData {
			var relatedSeason RelatedSeasonObjStruct
			relatedSeason.Id = val.Id
			relatedSeason.Seo_url = val.Seo.Url

			// A tien noi season name khong co thi lay title
			relatedSeason.Title = val.Season_name
			if val.Season_name == "" {
				relatedSeason.Title = val.Title
			}
			lstRelatedSS = append(lstRelatedSS, relatedSeason)
		}
	}

	return lstRelatedSS

}

func GetListIDSeasonByGroupID(groupId string, cache bool) ([]string, error) {
	var listID []string
	var keyCache = KV_LIST_ID_SEASON_BY_SHOW + "_" + groupId

	if cache {
		dataCache, err := mRedis.GetString(keyCache)
		if err == nil && dataCache != "" {
			err = json.Unmarshal([]byte(dataCache), &listID)
			if err == nil {
				return listID, nil
			}
		}
	}

	// Connect DB
	session, db, err := GetCollection()
	if err != nil {
		Sentry_log(err)
		return listID, err
	}
	defer session.Close()

	var where = bson.M{
		"group_id": groupId,
		"type":     3, //season
	}

	// Lây list id by group id
	var listIDStruct []struct {
		Id string
	}
	err = db.C(COLLECTION_VOD).Find(where).Sort("odr").All(&listIDStruct)
	if err != nil && err.Error() != "not found" {
		Sentry_log(err)
		return listID, err
	}

	listID = make([]string, 0)

	for _, val := range listIDStruct {
		//check exists id
		if ok, _ := In_array(val.Id, listID); !ok {
			listID = append(listID, val.Id)
		}

	}
	// Write cache
	dataByte, _ := json.Marshal(listID)
	mRedis.SetString(keyCache, string(dataByte), 3600)

	return listID, nil
}
