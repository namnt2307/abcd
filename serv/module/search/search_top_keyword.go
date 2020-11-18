package search

import (
	"fmt"

	. "cm-v5/serv/module"
	. "cm-v5/schema"
	recommendation "cm-v5/serv/module/recommendation"
)

func GetSeachKeywordResult(platform, token string, page, limit int, cacheActive bool) (SearchTopKeywordOutputStruct, error) {
	var SearchTopKeywordOutput SearchTopKeywordOutputStruct
	//Lay data tu Mysql
	SearchTopKeywordOutput, err := GetSearchKeywordByMysql(platform, page, limit, cacheActive)
	if err != nil {
		Sentry_log(err)
		return SearchTopKeywordOutput, err
	}
	return SearchTopKeywordOutput, nil
}

func GetSearchKeywordByYusp(keyword, user_id, device_id, scenario, platform, sort string, page, limit int) (SearchTopKeywordOutputStruct, error) {
	var SearchTopKeywordOutput SearchTopKeywordOutputStruct
	SearchTopKeywordOutput.Items = make([]SearchTopKeywordResultStruct, 0)

	dataSearchKeywordForYusp, err := recommendation.GetSearchKeyYuspRecommendationAPI(keyword, user_id, device_id, scenario, platform, sort, page, limit)
	if err != nil || len(dataSearchKeywordForYusp.ItemIds) <= 0 {
		Sentry_log(err)
		return SearchTopKeywordOutput, err
	}

	for k, val := range dataSearchKeywordForYusp.Items {
		var SearchTopKeywordResult SearchTopKeywordResultStruct
		SearchTopKeywordResult.Id = val.ItemId
		SearchTopKeywordResult.Keyword = val.Title
		SearchTopKeywordResult.Search_count = int(dataSearchKeywordForYusp.PredictionValues[k])
		SearchTopKeywordOutput.Items = append(SearchTopKeywordOutput.Items, SearchTopKeywordResult)
	}

	SearchTopKeywordOutput.Tracking_data = recommendation.GetTrackingData(dataSearchKeywordForYusp.RecommendationId, RECOMMENDATION_NAME_YUSP)
	return SearchTopKeywordOutput, nil
}

func GetSearchKeywordByMysql(platform string, page, limit int, cacheActive bool) (SearchTopKeywordOutputStruct, error) {
	var SearchTopKeywordOutput SearchTopKeywordOutputStruct
	var keyCache = PREFIX_REDIS_SEARCH_TOP_KEYWORD + "_" + fmt.Sprint(page) + "_" + fmt.Sprint(limit) + "_" + platform
	if cacheActive {
		valueCache, err := mRedis.GetString(keyCache)
		if err == nil && valueCache != "" {
			err = json.Unmarshal([]byte(valueCache), &SearchTopKeywordOutput)
			if err == nil {
				return SearchTopKeywordOutput, nil
			}
		}
	}
	//Connect mysql
	db_mysql, err := ConnectMySQL()
	if err != nil {
		Sentry_log(err)
		return SearchTopKeywordOutput, err
	}
	defer db_mysql.Close()

	topKeywordObj, err := db_mysql.Query("SELECT id , keyword , search_count FROM top_keyword WHERE status = 1 ORDER BY `order` ASC LIMIT ? OFFSET ?", limit, page)
	if err != nil {
		Sentry_log(err)
		return SearchTopKeywordOutput, err
	}

	SearchTopKeywordOutput.Items = make([]SearchTopKeywordResultStruct, 0)
	//fomart result in db
	for topKeywordObj.Next() {
		var SearchTopKeywordResult SearchTopKeywordResultStruct
		err = topKeywordObj.Scan(&SearchTopKeywordResult.Id, &SearchTopKeywordResult.Keyword, &SearchTopKeywordResult.Search_count)
		if err != nil {
			return SearchTopKeywordOutput, err
		}
		SearchTopKeywordOutput.Items = append(SearchTopKeywordOutput.Items, SearchTopKeywordResult)
	}
	// Write Redis
	dataByte, _ := json.Marshal(SearchTopKeywordOutput)
	mRedis.SetString(keyCache, string(dataByte), TTL_REDIS_2_HOURS)

	SearchTopKeywordOutput.Tracking_data = GetTrackingDataSearch(platform, SEARCH_KEYWORD)
	return SearchTopKeywordOutput, nil
}
