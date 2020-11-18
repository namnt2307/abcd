package middleware

import (
	"strings"
	"time"

	. "cm-v5/serv/module"
	Subscription "cm-v5/serv/module/subscription"
	"github.com/gin-gonic/gin"
	jsoniter "github.com/json-iterator/go"
)

var (
	json                                             = jsoniter.ConfigCompatibleWithStandardLibrary
	topicKafkaRequestLog, serverKafka string
)

func init() {
	serverKafka, _ = CommonConfig.GetString("KAFKA", "uri")
	topicKafkaRequestLog, _ = CommonConfig.GetString("KAFKA", "topic_log_request")
}

func ProcessResponse(c *gin.Context, latency float64, dataBodyRequest string, dataBodyResponse string) {
	//Send msg api slow query
	if latency >= 2 {
		Sentry_log_slow_request(c.Request.URL.String())
	}

	//Config off push data kafka
	if ENABLE_KAFKA == "disable" {
		return
	}

	if c.Writer.Status() != 200 {
		return
	}

	url := c.Request.URL.String()
	if strings.Contains(url, "/menu") == true {
		return
	}
	if strings.Contains(url, "/update-cache") == true {
		return
	}
	if strings.Contains(url, "/tips") == true {
		return
	}
	if strings.Contains(url, "/geo-check") == true {
		return
	}
	if strings.Contains(url, "/qnet") == true {
		return
	}
	if strings.Contains(url, "/episode") == true {
		return
	}
	if strings.Contains(url, "/related_videos") == true {
		return
	}
	if strings.Contains(url, "/related") == true {
		return
	}
	if strings.Contains(url, "/ribbon_v3") == true {
		return
	}
	if strings.Contains(url, "/page_banners_v3") == true {
		return
	}
	if strings.Contains(url, "/page_ribbons_v3") == true {
		return
	}
	if strings.Contains(url, "/health-check") == true {
		return
	}

	if strings.Contains(url, "/shp") == true {
		return
	}

	userId := c.GetString("user_id")
	platform := c.DefaultQuery("platform", "web")
	device_name := c.DefaultQuery("device_name", "")

	dataBodyResponse = FilterDataResponseToLog(c.Request.URL.String(), dataBodyResponse)

	if strings.Contains(url, "/feedback") == true {
		if dataBodyResponse == "{}" {
			return
		}
		dataBodyRequest = DecodeDataformToString(dataBodyRequest)
	}

	var dataInputProcessRequest DataInputLogStruct
	dataInputProcessRequest.Action = "response"
	dataInputProcessRequest.Uri = url
	dataInputProcessRequest.Method = c.Request.Method
	dataInputProcessRequest.User_agent = c.GetHeader("User-Agent")
	dataInputProcessRequest.Status = c.Writer.Status()
	dataInputProcessRequest.Duration = latency
	paramsList := c.Request.URL.Query()
	dataByte, _ := json.Marshal(paramsList)
	dataInputProcessRequest.Params = string(dataByte)
	dataInputProcessRequest.Data = dataBodyRequest
	dataInputProcessRequest.Data_response = dataBodyResponse
	dataInputProcessRequest.Ip, _ = GetClientIPHelper(c.Request, c)
	dataInputProcessRequest.User_id = userId
	dataInputProcessRequest.User_is_premium = 0
	dataInputProcessRequest.Token_access = c.Request.Header.Get("Authorization")

	if device_name != "" {
		platform = platform + ": " + device_name
	}
	dataInputProcessRequest.Platform = platform
	dataInputProcessRequest.Time = time.Now().Format("2006-01-02 15:04:05")

	if userId != "" {
		// Kiem tra user premium
		var Sub Subscription.SubcriptionObjectStruct
		SubcriptionsOfUser, err := Sub.GetListByUserId(userId)
		if err == nil && len(SubcriptionsOfUser) > 0 {
			dataInputProcessRequest.User_is_premium = 1
		}
	}




	var KafkaProducer KafkaProducerModelStruct
	KafkaProducer.SendMessage(topicKafkaRequestLog, serverKafka, dataInputProcessRequest)

	//Save log history payment
	if strings.Contains(url, "/ipn") ||
		strings.Contains(url, "/transaction-momo") ||
		strings.Contains(url, "/transaction-momo-app") ||
		strings.Contains(url, "/transaction-vnpay") ||
		strings.Contains(url, "/transaction-vnpay-qr") ||
		strings.Contains(url, "/transaction-ios") ||
		strings.Contains(url, "/payment-callback-ios") ||
		strings.Contains(url, "/transaction-google") ||
		strings.Contains(url, "/payment-callback-google") {

		go SaveLogHistoryPayment(userId, url, c.Request.Method, dataBodyRequest, dataBodyResponse)
	}

	// if strings.Contains(url, "/tracking/watch") == true {
	// 	go TrackContentViewInDay(dataInputProcessRequest, platform)
	// }

	// user_id := c.GetString("user_id")
	// if "" != user_id {
	// 	TrackUserOnline(user_id)
	// }
}

func FilterDataResponseToLog(url, dataBodyResponse string) string {
	if strings.Contains(url, "/content/") == true || strings.Contains(url, "/slug/content") == true {
		var DataInputLog_VOD DataInputLog_VOD_Struct
		err := json.Unmarshal([]byte(dataBodyResponse), &DataInputLog_VOD)
		if err == nil {
			dataByte, _ := json.Marshal(DataInputLog_VOD)
			return string(dataByte)
		}
	}

	if strings.Contains(url, "/livetv/detail") == true {
		var DataInputLog_LIVETV DataInputLog_LIVETV_Struct
		err := json.Unmarshal([]byte(dataBodyResponse), &DataInputLog_LIVETV.Technical)
		if err == nil {
			dataByte, _ := json.Marshal(DataInputLog_LIVETV)
			return string(dataByte)
		}
	}

	if strings.Contains(url, "/tracking/") == true ||
		strings.Contains(url, "/giftcode") == true ||
		strings.Contains(url, "/slug/events") == true ||
		strings.Contains(url, "/feedback") == true ||
		strings.Contains(url, "/events/") == true {
		return dataBodyResponse
	}

	if strings.Contains(url, "/ipn") ||
		strings.Contains(url, "/transaction-momo") ||
		strings.Contains(url, "/transaction-momo-app") ||
		strings.Contains(url, "/transaction-vnpay") ||
		strings.Contains(url, "/transaction-vnpay-qr") ||
		strings.Contains(url, "/transaction-ios") ||
		strings.Contains(url, "/payment-callback-ios") ||
		strings.Contains(url, "/transaction-google") ||
		strings.Contains(url, "/payment-callback-google") {
		return dataBodyResponse
	}

	return ""
}
