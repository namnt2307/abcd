package topviews

import (
	. "cm-v5/schema"
	. "cm-v5/serv/module"
	"fmt"
)

func GetTopViewsRankingByContentID(content_id string, cacheActive bool) int {
	// Hit local data
	var keyCache = LOCAL_TOPVIEWS_VOD + "_" + content_id
	if cacheActive {
		valC, err := LocalCache.GetValue(keyCache)
		if err == nil {
			num, _ := StringToInt(fmt.Sprint(valC))
			return num
		}
	}

	var rankingNumber = 0

	dataTopviews, err := GetDataTopViewsMongo(cacheActive)

	if err == nil {
		for _, val := range dataTopviews {
			if val.Id_content == content_id {
				rankingNumber = val.Ranking
				break
			}
		}
	}

	// Write local data
	LocalCache.SetValue(keyCache, rankingNumber, TTL_LOCALCACHE)

	return rankingNumber
}

func GetDataTopViewsMongo(cacheActive bool) ([]ViewsContentStruct, error) {
	var keyCache = KV_DATA_TOPVIEWS_CONTENT

	var topViewsData = make([]ViewsContentStruct, 0)

	if cacheActive {
		valueCache, err := mRedisKV.GetString(keyCache)
		if err == nil && valueCache != "" {
			err = json.Unmarshal([]byte(valueCache), &topViewsData)
			if err == nil {
				return topViewsData, nil
			}
		}
	}

	// Connect MongoDB
	session, db, err := GetCollection()
	if err != nil {
		return topViewsData, err
	}
	defer session.Close()

	err = db.C(COLLECTION_TOPVIEWS_CONTENT_FINAL).Find(nil).Sort("ranking").Limit(10).All(&topViewsData)
	if err != nil {
		return topViewsData, err
	}

	// Write Redis
	dataByte, _ := json.Marshal(topViewsData)
	mRedisKV.SetString(keyCache, string(dataByte), TTL_REDIS_5_MINUTE)

	return topViewsData, nil

}
