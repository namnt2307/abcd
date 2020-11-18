package search

import (
	"context"
	"fmt"
	. "cm-v5/schema"
	. "cm-v5/serv/module"
	artist "cm-v5/serv/module/artist"
	recommendation "cm-v5/serv/module/recommendation"
	seo "cm-v5/serv/module/seo"
	vod "cm-v5/serv/module/vod"
	"strconv"
	"github.com/gin-gonic/gin"
)

type DataSearchResultStruct struct {
	Title string
	Id    string
	Type  int
}

type DataSearchSuggestStruct struct {
	Type  int
	Name  string
	Title string
}

type DataOutputApiSearch struct {
	Meta struct {
		Page struct {
			Size          int   `json:"size" `
			Current       int   `json:"current" `
			Total_pages   int   `json:"total_pages"`
			Total_results int64 `json:"total_results"`
		} `json:"page" `
		Request_id string `json:"request_id"`
	} `json:"meta" `
	Results []struct {
		Meta struct {
			Id string `json:"id"`
		} `json:"_meta"`
		Type struct {
			Raw string `json:"raw"`
		} `json:"type"`
		Names struct {
			Raw []string `json:"raw"`
		} `json:"names"`
	} `json:"results" `
}

type DataPosApiSearch struct {
	Query string `json:"query" `
	Page  struct {
		Size    int `json:"size" `
		Current int `json:"current" `
	} `json:"page" `
	Filters struct {
		All []interface{} `json:"all"`
	} `json:"filters"`
}

type DataPostSearchFilterPlatform struct {
	Platforms []string `json:"platforms"`
}

type DataPostSearchFilterTag struct {
	Tag_ids []string `json:"tag_ids",omitempty`
}

func init() {
	var err error
	API_ES_HOSTNAME, err = CommonConfig.GetString("APIES", "host_name")
	API_ES_TOKEN, err = CommonConfig.GetString("APIES", "token")
	API_ES_PREFIX, err = CommonConfig.GetString("APIES", "prefix")

	if err != nil {
		API_ES_HOSTNAME = "https://1c65fb04810047b5814f539d8895366b.ent-search.ap-southeast-1.aws.cloud.es.io"
		API_ES_TOKEN = "search-a18ci3rowcq1fkwuope48y4k"
		API_ES_PREFIX = "dev-"

	}

}

func GetSearchResult(keyword, tags, platform, version, token string, page, limit, entityType int, c *gin.Context, cacheActive bool) (SearchResultStruct, error) {
	var SearchResult SearchResultStruct

	var PushDataRecommandFunc = func(c *gin.Context, TrackingData TrackingDataStruct) {
		TrackingUserData, err := recommendation.GetTracking(c, "", TrackingData)
		if err == nil {
			TrackingUserData.Query = keyword
			TrackingUserData.Count = SearchResult.Metadata.Total
			recommendation.PushDataToKafka(SEARCH, TrackingUserData)
		}
	}

	//Lay data tu Elasticsearch
	SearchResult, err := GetSearchByElasticsearch(keyword, tags, platform, version, page, limit, entityType, c.Request.Context(), cacheActive)
	if err != nil {
		Sentry_log(err)
		return SearchResult, err
	}

	//Push data Kafka recommand search
	go PushDataRecommandFunc(c, SearchResult.Tracking_data)

	return SearchResult, nil
}

func GetSearchByElasticsearch(keyword, tags, platform, version string, page, limit, entityType int, ctx context.Context, cacheActive bool) (SearchResultStruct, error) {
	var SearchResult SearchResultStruct
	var DataPost DataPosApiSearch

	// Thông tin tags and platform
	var dataFilterPlafrom DataPostSearchFilterPlatform
	platformInfo := Platform(platform)
	dataFilterPlafrom.Platforms = []string{fmt.Sprintf(`%d`, platformInfo.Id)}
	DataPost.Filters.All = append(DataPost.Filters.All, dataFilterPlafrom)

	if tags != "" {
		var dataFilterTag DataPostSearchFilterTag
		dataFilterTag.Tag_ids = []string{tags}
		// Data platform
		DataPost.Filters.All = append(DataPost.Filters.All, dataFilterTag)
	}

	DataPost.Query = keyword
	DataPost.Page.Size = limit
	DataPost.Page.Current = page + 1

	Respon, Err := PostAPIJsonTonken(API_ES_HOSTNAME+"/api/as/v1/engines/"+API_ES_PREFIX+"vieon-v3/search", DataPost, API_ES_TOKEN)
	if Err != nil {
		Sentry_log(Err)
		return SearchResult, Err

	}
	//Dich nguoc lai data tra KQ

	var DataOutputEs DataOutputApiSearch
	Err = json.Unmarshal(Respon, &DataOutputEs)
	if Err != nil {
		Sentry_log(Err)
		return SearchResult, Err
	}
	// log.Println(DataOutputEs)
	// End

	SearchResult.Items = make([]SearchItemResultStruct, 0)

	if len(DataOutputEs.Results) > 0 {

		var DataSearchResultStructs []DataSearchResultStruct
		// get list ID & Type
		for _, hit := range DataOutputEs.Results {
			// log.Println(hit)
			var tempS DataSearchResultStruct
			tempS.Title = hit.Names.Raw[0]
			tempS.Id = hit.Meta.Id
			tempS.Type, _ = strconv.Atoi(hit.Type.Raw)

			DataSearchResultStructs = append(DataSearchResultStructs, tempS)
		}
		// Parse data to SearchResult
		if len(DataSearchResultStructs) > 0 {
			var listPeopleID []string
			var listVodID []string

			//Lay danh sach id VOD, Artist
			for _, val := range DataSearchResultStructs {
				switch val.Type {
				case TYPE_SEARCH_ITEM_SEASON:
					listVodID = append(listVodID, val.Id)
				case TYPE_SEARCH_ITEM_ARTIST:
					listPeopleID = append(listPeopleID, val.Id)
				default:
					listVodID = append(listVodID, val.Id)
				}
			}

			//Lay data VOD, Artist trong DB theo ID
			vodDataObjects, _ := vod.GetVODByListID(listVodID, platformInfo.Id, 1, true)
			artistObjects, _ := artist.GetArtistByListID(listPeopleID, true)

			// Thu fix bug search DT-10579
			// TTTTTTTTTTTTTTT tạm comment
			// if len(vodDataObjects) < len(listVodID) {
			// 	//async
			// 	go func(vodDataObjects []VODDataObjectStruct, listVodID []string) {
			// 		ContentReindexElastic(vodDataObjects, listVodID)
			// 	}(vodDataObjects, listVodID)

			// }

			// if len(artistObjects) < len(listPeopleID) {
			// 	go func(artistObjects []ArtistObjectStruct, listPeopleID []string) {
			// 		ArticleReindexElastic(artistObjects, listPeopleID)
			// 	}(artistObjects, listPeopleID)

			// }

			var vodDataObjectsForSearch = make(map[string]VODDataObjectStruct, 0)
			var artistObjectsForSearch = make(map[string]ArtistObjectStruct, 0)

			for _, val := range vodDataObjects {
				vodDataObjectsForSearch[val.Id] = val
			}
			for _, val := range artistObjects {
				artistObjectsForSearch[val.Id] = val
			}

			for position, val := range DataSearchResultStructs {
				var SearchItemResult SearchItemResultStruct
				switch val.Type {
				case TYPE_SEARCH_ITEM_ARTIST:
					if valA, ok := artistObjectsForSearch[val.Id]; ok {
						SearchItemResult.Is_artist = true
						SearchItemResult.Name = valA.Name
						SearchItemResult.Seo = valA.Seo
						SearchItemResult.Job = valA.Job
						SearchItemResult.Type = TYPE_SEARCH_ITEM_ARTIST
						SearchItemResult.Images = valA.Images
						SearchItemResult.Id = valA.Id
					}
				default:
					if valV, ok := vodDataObjectsForSearch[val.Id]; ok {
						var eps = valV.Current_episode
						if eps == "" {
							eps = "0"
						}

						SearchItemResult.Is_artist = false
						SearchItemResult.Episode = valV.Episode
						SearchItemResult.Title = val.Title
						SearchItemResult.Resolution = valV.Resolution
						SearchItemResult.Total_rate = valV.Total_rate
						SearchItemResult.Current_episode = eps
						// SearchItemResult.Slug = valV.Slug
						SearchItemResult.Seo = valV.Seo
						SearchItemResult.Group_id = valV.Group_id
						SearchItemResult.Type = val.Type
						SearchItemResult.Id = valV.Id
						// SearchItemResult.Slug_seo = valV.Slug_seo
						SearchItemResult.Is_premium = valV.Is_premium

						//Khanh DT-15025
						SearchItemResult.Avg_rate = ReCalculatorRating(valV.Min_rate, valV.Avg_rate)
						SearchItemResult.Total_rate = valV.Total_rate

						//Switch images platform
						var ImagesMapping ImagesOutputObjectStruct
						// switch platformInfo.Type {
						// case "web":
						// 	// ImagesMapping.Vod_thumb = BuildImage(valV.Images.Web.Vod_thumb)
						// 	ImagesMapping.Thumbnail = BuildImage(valV.Images.Web.Vod_thumb)
						// case "smarttv":
						// 	ImagesMapping.Thumbnail = BuildImage(valV.Images.Smarttv.Thumbnail)
						// case "app":
						// 	ImagesMapping.Thumbnail = BuildImage(valV.Images.App.Thumbnail)
						// }
						switch platformInfo.Type {
						case "web":
							ImagesMapping.Vod_thumb = BuildImage(valV.Images.Web.Vod_thumb)
							ImagesMapping.Thumbnail = BuildImage(valV.Images.Web.Vod_thumb)
						case "smarttv":
							ImagesMapping.Thumbnail = BuildImage(valV.Images.Smarttv.Thumbnail)
						case "app":
							ImagesMapping.Thumbnail = BuildImage(valV.Images.App.Thumbnail)
						}
						ImagesMapping = MappingImagesV4(platformInfo.Type, ImagesMapping, valV.Images, true)
						SearchItemResult.Images = ImagesMapping
					}

				}
				if SearchItemResult.Id != "" {
					SearchItemResult.Request_id = DataOutputEs.Meta.Request_id
					SearchItemResult.Position = fmt.Sprintf(`%d`, (page*limit)+position)
					SearchResult.Items = append(SearchResult.Items, SearchItemResult)
				}
			}
		}
	}

	SearchResult.Seo = seo.FormatSeoForSearch(keyword, SearchResult)
	SearchResult.Metadata.Limit = limit
	SearchResult.Metadata.Page = page
	SearchResult.Metadata.Keyword = keyword
	SearchResult.Tracking_data = GetTrackingDataSearch(platform, SEARCH_RESULT)
	SearchResult.Metadata.Total = DataOutputEs.Meta.Page.Total_results // searchResult.Hits.TotalHits

	// Write Redis
	// dataByte, _ := json.Marshal(SearchResult)
	// mRedis.SetString(keyCache, string(dataByte), TTL_KVCACHE)
	return SearchResult, nil
}

func GetSeachByYusp(keyword, user_id, device_id, scenario, platform, sort string, page, limit int) (SearchResultStruct, error) {
	var SearchResult SearchResultStruct
	var platformInfo = Platform(platform)

	SearchResult.Items = make([]SearchItemResultStruct, 0)
	dataSearchForYusp, err := recommendation.GetSearchYuspRecommendationAPI(keyword, user_id, device_id, scenario, platform, sort, page, limit)
	SearchResult.Seo = seo.FormatSeoForSearch(keyword, SearchResult)
	SearchResult.Metadata.Limit = limit
	SearchResult.Metadata.Page = page
	SearchResult.Metadata.Keyword = keyword
	if err != nil || len(dataSearchForYusp.ItemIds) <= 0 {
		Sentry_log(err)
		return SearchResult, err
	}

	vodDataObjects, _ := vod.GetVODByListID(dataSearchForYusp.ItemIds, platformInfo.Id, 1, true)
	for _, valV := range vodDataObjects {
		var eps = valV.Current_episode
		if eps == "" {
			eps = "0"
		}
		var SearchItemResult SearchItemResultStruct
		SearchItemResult.Is_artist = false
		SearchItemResult.Episode = valV.Episode
		SearchItemResult.Title = valV.Title
		SearchItemResult.Resolution = valV.Resolution
		SearchItemResult.Total_rate = valV.Total_rate
		SearchItemResult.Current_episode = eps
		SearchItemResult.Slug = valV.Slug
		SearchItemResult.Seo = valV.Seo
		SearchItemResult.Group_id = valV.Group_id
		SearchItemResult.Type = valV.Type
		SearchItemResult.Id = valV.Id
		SearchItemResult.Slug_seo = valV.Slug_seo
		SearchItemResult.Is_premium = valV.Is_premium

		//Khanh DT-15025
		SearchItemResult.Avg_rate = ReCalculatorRating(valV.Min_rate, valV.Avg_rate)
		SearchItemResult.Total_rate = valV.Total_rate

		//Switch images platform
		var ImagesMapping ImagesOutputObjectStruct
		switch platformInfo.Type {
		case "web":
			ImagesMapping.Vod_thumb = BuildImage(valV.Images.Web.Vod_thumb)
			ImagesMapping.Thumbnail = BuildImage(valV.Images.Web.Vod_thumb)
		case "smarttv":
			ImagesMapping.Thumbnail = BuildImage(valV.Images.Smarttv.Thumbnail)
		case "app":
			ImagesMapping.Thumbnail = BuildImage(valV.Images.App.Thumbnail)
		}
		ImagesMapping = MappingImagesV4(platformInfo.Type, ImagesMapping, valV.Images, true)
		SearchItemResult.Images = ImagesMapping
		SearchResult.Items = append(SearchResult.Items, SearchItemResult)
	}

	SearchResult.Metadata.Total = int64(dataSearchForYusp.TotalResults)
	SearchResult.Tracking_data = recommendation.GetTrackingData(dataSearchForYusp.RecommendationId, RECOMMENDATION_NAME_YUSP)

	return SearchResult, nil
}

func GetSearchSuggest(keyword, tags, platform, version, token string, page, limit, entityType int, ctx context.Context, cacheActive bool) (SearchSuggestOutputStruct, error) {
	var SearchSuggest SearchSuggestOutputStruct
	//Lay data tu Elasticsearch
	SearchSuggest, err := GetSearchSuggestNewByElasticsearch(keyword, tags, platform, version, page, limit, entityType, ctx, cacheActive)
	if err != nil {
		Sentry_log(err)
		return SearchSuggest, err
	}
	return SearchSuggest, nil
}

func GetSearchSuggestByYusp(keyword, user_id, device_id, scenario, platform, sort string, page, limit int) (SearchSuggestOutputStruct, error) {
	var SearchSuggest SearchSuggestOutputStruct
	SearchSuggest.Items = []string{}

	dataSearchSuggestForYusp, err := recommendation.GetSearchKeyYuspRecommendationAPI(keyword, user_id, device_id, scenario, platform, sort, page, limit)
	if err != nil || len(dataSearchSuggestForYusp.Items) <= 0 {
		Sentry_log(err)
		return SearchSuggest, err
	}

	for _, val := range dataSearchSuggestForYusp.Items {
		SearchSuggest.Items = append(SearchSuggest.Items, val.Title)
	}

	SearchSuggest.Tracking_data = recommendation.GetTrackingData(dataSearchSuggestForYusp.RecommendationId, RECOMMENDATION_NAME_YUSP)
	return SearchSuggest, nil
}

func GetSearchSuggestNewByElasticsearch(keyword, tags, platform, version string, page, limit, entityType int, ctx context.Context, cacheActive bool) (SearchSuggestOutputStruct, error) {
	var SearchSuggest SearchSuggestOutputStruct
	SearchSuggest.Items = []string{}

	// Call API POST

	var DataPost DataPosApiSearch

	// Thông tin tags and platform
	var dataFilterPlafrom DataPostSearchFilterPlatform
	platformInfo := Platform(platform)
	dataFilterPlafrom.Platforms = []string{fmt.Sprintf(`%d`, platformInfo.Id)}
	DataPost.Filters.All = append(DataPost.Filters.All, dataFilterPlafrom)

	if tags != "" {
		var dataFilterTag DataPostSearchFilterTag
		dataFilterTag.Tag_ids = []string{tags}
		// Data platform
		DataPost.Filters.All = append(DataPost.Filters.All, dataFilterTag)
	}

	DataPost.Query = keyword
	DataPost.Page.Size = limit
	DataPost.Page.Current = page + 1

	Respon, Err := PostAPIJsonTonken(API_ES_HOSTNAME+"/api/as/v1/engines/"+API_ES_PREFIX+"vieon-v3/search", DataPost, API_ES_TOKEN)
	if Err != nil {
		Sentry_log(Err)
		return SearchSuggest, Err

	}
	//Dich nguoc lai data tra KQ

	var DataOutputEs DataOutputApiSearch
	Err = json.Unmarshal(Respon, &DataOutputEs)
	if Err != nil {
		Sentry_log(Err)
		return SearchSuggest, Err
	}
	// log.Println(DataOutputEs)

	if len(DataOutputEs.Results) > 0 {
		for _, val := range DataOutputEs.Results {
			SearchSuggest.Items = append(SearchSuggest.Items, val.Names.Raw[0])
		}
	}
	// End

	SearchSuggest.Tracking_data = GetTrackingDataSearch(platform, ZERO_RESULT)
	// log.Println(SearchSuggest)
	return SearchSuggest, nil
}
