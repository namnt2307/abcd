package tag

import (
	"fmt"
	"strings"

	. "cm-v5/schema"
	. "cm-v5/serv/module"
	seo "cm-v5/serv/module/seo"

	"gopkg.in/mgo.v2/bson"
)

func GetTagIDBySlug(listSlug []string, cacheActive bool) ([]string, error) {
	var listTagIds []string
	var keyCache = TAGS_VOD + "_" + strings.Join(listSlug, "-")

	if cacheActive {
		valueCache, err := mRedis.GetString(keyCache)
		if err == nil && valueCache != "" {
			err = json.Unmarshal([]byte(valueCache), &listTagIds)
			if err == nil {
				return listTagIds, nil
			}
		}
	}

	// slugSplits := strings.Split(listSlug[0], "-col#tag-")
	// if len(slugSplits) == 2 {
		// Check slug exists
		ref_id := seo.CheckExistsSEOBySlug(listSlug[0])
		if ref_id != "" {
			listTagIds = []string{ref_id}
		}
	// }""

	if len(listTagIds) <= 0 {
		//Connect mysql
		db_mysql, err := ConnectMySQL()
		if err != nil {
			return listTagIds, err
		}
		defer db_mysql.Close()

		queryTagSlug := strings.Join(listSlug, "','")
		sqlRaw := fmt.Sprintf(`SELECT id FROM tag WHERE slug IN ('%s')`, queryTagSlug)
		tagIds, err := db_mysql.Query(sqlRaw)
		if err != nil {
			return listTagIds, err
		}

		for tagIds.Next() {
			var id string
			err = tagIds.Scan(&id)
			if err != nil {
				return listTagIds, err
			}
			listTagIds = append(listTagIds, id)
		}
	}
	

	// Write Redis
	dataByte, _ := json.Marshal(listTagIds)
	mRedis.SetString(keyCache, string(dataByte), TTL_REDIS_LV1)

	return listTagIds, nil
}

func GetListIDTagVOD(listTagID []string, sort, page, limit int, platform string, cacheActive bool) ([]string, error) {
	var listVodID []string
	var keyCache = TAGS_VOD_ID + "_" + strings.Join(listTagID, "-") + "_" + fmt.Sprint(sort) + "_" + fmt.Sprint(page) + "_" + fmt.Sprint(limit) + "_" + platform
	if cacheActive {
		dataCache, err := mRedis.GetString(keyCache)
		if err == nil && dataCache != "" {
			err = json.Unmarshal([]byte(dataCache), &listVodID)
			if err == nil {
				return listVodID, nil
			}
		}
	}

	// Connect DB
	session, db, err := GetCollection()
	if err != nil {
		return listVodID, err
	}
	defer session.Close()

	var field_sort string
	switch sort {
	case 3:
		field_sort = "-view"
	case 2:
		field_sort = "-release_date"
	case 1:
		field_sort = "release_date"
	}

	platformInfo := Platform(platform)
	var where = bson.M{
		"tags.id":   bson.M{"$in": listTagID},
		"platforms": platformInfo.Id,
	}

	var vodDataObjects []VODDataObjectStruct
	err = db.C(COLLECTION_VOD).Find(where).Sort(field_sort).Skip(page * limit).Limit(limit).All(&vodDataObjects)
	if err != nil {
		return listVodID, err
	}

	if len(vodDataObjects) > 0 {
		for _, val := range vodDataObjects {
			listVodID = append(listVodID, val.Id)
		}
	}

	// Write cache
	dataByte, _ := json.Marshal(listVodID)
	mRedis.SetString(keyCache, string(dataByte), TTL_REDIS_LV1)

	return listVodID, nil
}
