package recommendation

import (
	"fmt"
	"net/url"
	"strings"

	. "cm-v5/serv/module"
	. "cm-v5/schema"
	jsoniter "github.com/json-iterator/go"
)

var (
	json                                                        = jsoniter.ConfigCompatibleWithStandardLibrary
	urlYUSP, domainYUSP, prefixTopicKafka, topicKafka, serverKafka string
	PREFIX_DATA_YUSP                                            = "data_yusp"

	scenariosYUSP = make(map[string]string , 0)
	providerRecommend = make(map[string]int , 0)
)

type YUSPObjStruct struct {
	ScenarioId       string
	NumberLimit      int
	CookieId         string
	UserId           string
	CurrentItemId    string
	ResultNameValue  interface{}
	ResultNameValues interface{}
	PagingOffset     int
	Query            string
	Date             string
	Type             interface{}
	Sort             string
	NameValue        interface{}
	nameValues       interface{}
	CallObj          interface{}
}

type DataCallYuspObjStruct struct {
	RecommendationId string
	Items            []ItemsYuspStruct
	ItemIds          []string
	PredictionValues []float64
	CurrentItemId    string
	OutputNameValues interface{}
	TotalResults     int `json:"totalResults" `
}

type ItemsYuspStruct struct {
	ItemId     string
	Title      string
	ItemType   string
	Hidden     bool
	FromDate   int64
	ToDate     int64
	NameValues []struct {
		Name  string
		Value string
	}
}

type InfoUserRecommandationStruct struct {
	Device_id string
	User_id   string
}

func init() {
	//Get config Kafka
	serverKafka, _ = CommonConfig.GetString("KAFKA", "uri")
	prefixTopicKafka, _ = CommonConfig.GetString("KAFKA", "prefix")
	topicKafka, _ = CommonConfig.GetString("KAFKA", "topic_rec_click")


	//Get config YUSP
	domainYUSP, _  = CommonConfig.GetString("YUSP", "domain")
	urlYUSP, _  = CommonConfig.GetString("YUSP", "url")
	

	scenarios, _ := CommonConfig.GetString("YUSP", "scenarios")
	if scenarios != "" {
		json.Unmarshal([]byte(scenarios), &scenariosYUSP)
	}
	
	providerStr, _ := CommonConfig.GetString("RECOMMENDATION", "provider")
	if providerStr != "" {
		json.Unmarshal([]byte(providerStr), &providerRecommend)
	}

}

func GetRandomDefaultTrackingData(recommendationType string) TrackingDataStruct {
	recommendationID := GetRandomDefaultRecommendationID(recommendationType)
	trackingData := GetTrackingData(recommendationID, RECOMMENDATION_VIEON)
	return trackingData
}

func GetRandomDefaultRecommendationID(recommendationType string) string {

	// func RandStringAny anh huong den chat luong code => off di => return rong
	return ""
	var charset = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"
	var preFix = RandStringAny(charset, RECOMMENDATION_MIN_RANDOM_STRING, RECOMMENDATION_MAX_RANDOM_STRING)
	var sufFix = RandStringAny(charset, RECOMMENDATION_MIN_RANDOM_STRING, RECOMMENDATION_MAX_RANDOM_STRING)

	if len(recommendationType) > 0 {
		return preFix + "-" + recommendationType + "-" + sufFix
	}
	return ""
}

func GetTrackingData(recommendationID, recommendationSource string) TrackingDataStruct {
	var trackingData TrackingDataStruct
	if len(recommendationID) > 0 {
		trackingData.Type = recommendationSource
		trackingData.Recommendation_id = recommendationID
		return trackingData
	}

	return trackingData
}

func GetProviderRecommend() string {
	var scenario = RECOMMENDATION_NAME_VIEPLAY
	if len(providerRecommend) <= 0 {
		return scenario 
	}
	r := RandNumber(0, 100)
	for k, val := range providerRecommend {
		if val > r {
			scenario = k
		}
	}
	return scenario
}

func GetScenariosByRibbonID(ribbonID string) string {
	if len(scenariosYUSP) <= 0 {
		return ""
	}

	for k, val := range scenariosYUSP {
		if ribbonID == k {
			return val
		}
	}
	return ""
}

func GetScenariosContentDetail() string {
	return "VIDEO_DETAIL"
}


func GetContentDetailYuspRecommendationAPI(user_id, device_id, scenario, platform , content_id string, page, limit int) (DataCallYuspObjStruct, error) {
	resultValue := []string{CONTENT_TYPE, AVAILABLE_EPISODE, TOTAL_EPISODE, SEASON, SERIES}

	//build query get data YUSP
	var buildQuery = url.Values{}
	buildQuery.Add("currentItemId", content_id)
	
	query := BuildQueryQueryYusp(buildQuery, page, limit, user_id, device_id, scenario, platform, resultValue)
	dataYusp, err := CallYuspRecommendationAPI(query)

	return dataYusp, err
}


//get_common_yusp_recommendation_api
func GetCommonYuspRecommendationAPI(user_id, device_id, scenario, platform string, page, limit int) (DataCallYuspObjStruct, error) {
	resultValue := []string{CONTENT_TYPE, AVAILABLE_EPISODE, TOTAL_EPISODE, SEASON, SERIES}

	//build query get data YUSP
	var buildQuery = url.Values{}
	query := BuildQueryQueryYusp(buildQuery, page, limit, user_id, device_id, scenario, platform, resultValue)
	dataYusp, err := CallYuspRecommendationAPI(query)

	return dataYusp, err
}

//get_search_key_yusp_recommendation_api
func GetSearchKeyYuspRecommendationAPI(keyword, user_id, device_id, scenario, platform, sort string, page, limit int) (DataCallYuspObjStruct, error) {
	resultValue := []string{VALUE}
	var buildQuery = url.Values{}
	if keyword != "" {
		buildQuery.Add("query", keyword)
	}
	query := BuildQueryQueryYusp(buildQuery, page, limit, user_id, device_id, scenario, platform, resultValue)
	dataYusp, err := CallYuspRecommendationAPI(query)
	return dataYusp, err
}

func GetSearchYuspRecommendationAPI(keyword, user_id, device_id, scenario, platform, sort string, page, limit int) (DataCallYuspObjStruct, error) {
	resultValue := []string{CONTENT_TYPE, AVAILABLE_EPISODE, TOTAL_EPISODE, SEASON, SERIES}
	//build query get data YUSP
	var buildQuery = url.Values{}
	buildQuery.Add("query", keyword)
	query := BuildQueryQueryYusp(buildQuery, page, limit, user_id, device_id, scenario, platform, resultValue)
	dataYusp, err := CallYuspRecommendationAPI(query)
	return dataYusp, err
}

// call_yusp_recommendation_api
func CallYuspRecommendationAPI(query string) (DataCallYuspObjStruct, error) {
	var DatCallYuspObj DataCallYuspObjStruct

	//Call API
	url := domainYUSP + urlYUSP + "?" + query
	responseData, _ := GetAPI(url)
	//Format lai data tra ve
	err := json.Unmarshal(responseData, &DatCallYuspObj)
	if err != nil {
		Sentry_log(err)
		return DatCallYuspObj, err
	}
	return DatCallYuspObj, nil
}

func GetScenarioByPlatform(scenario, platform string) string {
	var platformScenario string

	switch platform {
	case "smarttv":
		platformScenario = DEVICEKEY_SMARTTV + scenario
	case "web":
		platformScenario = DEVICEKEY_WEB + scenario
	case "ios":
		platformScenario = DEVICEKEY_IOS + scenario
	case "android":
		platformScenario = DEVICEKEY_ANDROID + scenario
	case "samsung_tv":
		platformScenario = DEVICEKEY_SMARTTV + scenario
	case "sony_androidtv":
		platformScenario = DEVICEKEY_SMARTTV + scenario
	case "lg_tv":
		platformScenario = DEVICEKEY_SMARTTV + scenario
	case "mobile_web":
		platformScenario = DEVICEKEY_WEB + scenario
	default:
		platformScenario = DEVICEKEY_WEB + scenario
	}
	return platformScenario
}

func BuildQueryQueryYusp(buildQuery url.Values, page, limit int, user_id, device_id, scenario, platform string, resultValue []string) string {
	pagingOffset := page * limit
	numberLimit := limit
	platformScenario := GetScenarioByPlatform(scenario, platform)
	scenarioId := platformScenario

	buildQuery.Add("scenarioId", scenarioId)
	buildQuery.Add("pagingOffset", fmt.Sprint(pagingOffset))
	buildQuery.Add("numberLimit", fmt.Sprint(numberLimit))
	buildQuery.Add("cookieId", device_id)
	buildQuery.Add("userId", user_id)
	buildQuery.Add("sort", "relevance")

	if len(resultValue) > 0 {
		for _, val := range resultValue {
			buildQuery.Add("resultNameValue", val)
		}
	}

	return HttpBuildQuery(buildQuery)
}

func GetUserInfo(token string) InfoUserRecommandationStruct {
	var InfoUserRecommandation InfoUserRecommandationStruct
	if token != "" {
		jwt, err := LocalAuthVerify(token)
		if err != nil {
			Sentry_log(err)
			return InfoUserRecommandation
		}

		InfoUserRecommandation.Device_id = jwt.DeviceId
		InfoUserRecommandation.User_id = jwt.Subject
		if strings.HasPrefix(InfoUserRecommandation.User_id, "anonymous_") {
			InfoUserRecommandation.User_id = ""
		}
	}

	return InfoUserRecommandation
}

func GenerateCookieID(InfoUserRecommandation InfoUserRecommandationStruct) string {
	// hashMD5 := md5.New()
	// dataByte, _ := json.Marshal(InfoUserRecommandation)
	// io.WriteString(hashMD5, string(dataByte))
	// return fmt.Sprintf("%x", hashMD5.Sum(nil))
	return InfoUserRecommandation.Device_id
}

func PushDataToKafka(action string, trackingData TrackingUserDataStruct) {
	if trackingData.Item_id != "" || action == SEARCH {
		trackingData.Action = action
		var KafkaProducer KafkaProducerModelStruct
		KafkaProducer.SendMessage(prefixTopicKafka+topicKafka, serverKafka, trackingData)
	}
}
