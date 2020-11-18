package search

import (
	"strings"
	"time"

	. "cm-v5/schema"
	. "cm-v5/serv/module"
	artist "cm-v5/serv/module/artist"
	vod "cm-v5/serv/module/vod"
)

type DataResponseStruct struct {
	Success bool `json:"success" `
}

type DetailHistorySearchStruct struct {
	Id        string          `json:"id" `
	Is_artist bool            `json:"is_artist" `
	Name      string          `json:"name,omitempty" `
	Slug      string          `json:"slug,omitempty" `
	Title     string          `json:"title,omitempty" `
	Seo       SeoObjectStruct `json:"seo" `
}

type DataHistorySearchStruct struct {
	Items    []DetailHistorySearchStruct `json:"items" `
	Metadata struct {
		Limit  int    `json:"limit" `
		Cursor string `json:"cursor" `
	} `json:"metadata" `
}

var (
	API_ES_HOSTNAME string
	API_ES_TOKEN    string
	API_ES_PREFIX   string
)

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

func AddHistorySearchByUser(userId string, keyword string, contentIds string, artistIds string) (DataResponseStruct, error) {
	var DataResponse DataResponseStruct
	DataResponse.Success = true

	//Parse string to array
	var arrID []string
	if contentIds != "" {
		arrIDContent, _ := ParseStringToArray(contentIds)
		for _, val := range arrIDContent {
			arrID = append(arrID, "content/"+val)
		}
	}
	if artistIds != "" {
		arrIDArtist, _ := ParseStringToArray(artistIds)
		for _, val := range arrIDArtist {
			arrID = append(arrID, "artist/"+val)
		}
	}
	if len(arrID) <= 0 {
		DataResponse.Success = false
		return DataResponse, nil
	}

	var curTime = time.Now().Unix()
	for _, val := range arrID {
		mRedisUSC.ZAdd(PREFIX_REDIS_ZRANGE_USC_HISTORY_SEARCH+userId, float64(curTime), val)
	}
	return DataResponse, nil
}

func RemoveHistorySearchByUser(userId string) (DataResponseStruct, error) {
	var DataResponse DataResponseStruct
	DataResponse.Success = true
	mRedisUSC.Del(PREFIX_REDIS_ZRANGE_USC_HISTORY_SEARCH + userId)
	return DataResponse, nil
}

func GetHistorySearchByUserID(userId string, platform string, page int, limit int) (DataHistorySearchStruct, error) {
	var platformDetail = Platform(platform)
	var DataHistorySearch DataHistorySearchStruct
	DataHistorySearch.Items = make([]DetailHistorySearchStruct, 0)
	DataHistorySearch.Metadata.Limit = limit

	var start = page * limit
	var stop = (start + limit) - 1
	listIds := mRedisUSC.ZRevRange(PREFIX_REDIS_ZRANGE_USC_HISTORY_SEARCH+userId, int64(start), int64(stop))

	for _, val := range listIds {
		valId := strings.Split(val, "/")
		var DetailHistorySearch DetailHistorySearchStruct
		if len(valId) > 1 {
			switch valId[0] {
			case "content":
				valContent, err := vod.GetVodDetail(valId[1], platformDetail.Id, true)
				if err == nil {
					DetailHistorySearch.Id = valContent.Id
					DetailHistorySearch.Is_artist = false
					if valContent.Movie.Title != "" {
						DetailHistorySearch.Title = valContent.Movie.Title + " " + valContent.Title
					} else {
						DetailHistorySearch.Title = valContent.Title
					}
					DetailHistorySearch.Slug = valContent.Slug
					DetailHistorySearch.Seo = valContent.Seo
					DataHistorySearch.Items = append(DataHistorySearch.Items, DetailHistorySearch)
				}
			case "artist":
				valArtist, err := artist.GetInfoArtist(valId[1], true)
				if err == nil {
					DetailHistorySearch.Id = valArtist.Id
					DetailHistorySearch.Is_artist = true
					DetailHistorySearch.Name = valArtist.Name
					DetailHistorySearch.Seo = valArtist.Seo
					DataHistorySearch.Items = append(DataHistorySearch.Items, DetailHistorySearch)
				}
			}
		}
	}

	return DataHistorySearch, nil
}

type DataPosApiClickSearch struct {
	Query       string   `json:"query" `
	Request_id  string   `json:"request_id" `
	Document_id string   `json:"document_id" `
	Tags        []string `json:"tags" `
}

func SendClickEsAwsByUser(userId string, keyword string, contentIds string, artistIds string, platform string, request_id string, position string) error {
	// Call API POST
	//Parse string to array
	if contentIds != "" {
		arrIDContent, _ := ParseStringToArray(contentIds)
		for _, val := range arrIDContent {
			var DataPost DataPosApiClickSearch
			DataPost.Query = keyword
			DataPost.Document_id = val

			DataPost.Tags = []string{"user:" + userId, platform}
			if position != "" {
				DataPost.Tags = append(DataPost.Tags, "item:"+position)
			}
			if request_id != "" {
				DataPost.Request_id = request_id
			}
			_, Err := PostAPIJsonTonken(API_ES_HOSTNAME+"/api/as/v1/engines/"+API_ES_PREFIX+"vieon-v3/click", DataPost, API_ES_TOKEN)
			if Err != nil {
				Sentry_log(Err)
				return Err
			}

		}
	}
	return nil
}
