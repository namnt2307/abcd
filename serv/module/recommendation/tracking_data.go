package recommendation

import (
	"fmt"
	"net/http"

	. "cm-v5/serv/module"
	. "cm-v5/schema"
	"github.com/gin-gonic/gin"
)

type TrackingUserDataStruct struct {
	Tracking_id   string `json:"tracking_id" `
	Tracking_type string `json:"tracking_type" `
	Uri           string `json:"uri" `
	Item_id       string `json:"item_id" `
	User_agent    string `json:"user_agent" `
	Ip            string `json:"ip" `
	Token         string `json:"token" `
	Time          string `json:"time" `
	Percentage    int64  `json:"percentage" `
	Duration      int64  `json:"duration" `
	Query         string `json:"query" `
	Count         int64  `json:"count" `
	Action        string `json:"action" `
}

type ResultTrackingUserDataStruct struct {
	Success bool `json:"success" `
}

var YuspConstantsEvents = []string{"VIEW", "FREE_VIEW", "WATCH", "SHARE", "REC_CLICK", "ADD_TO_PLAYLIST", "REMOVE_FROM_PLAYLIST", "LIKE", "DISLIKE", "COMMENT", "LOGIN", "SEARCH", "SUBSCRIBE", "RATING"}

func TrackingEvent(c *gin.Context) {
	trackingID := c.PostForm("tracking_id")
	trackingType := c.PostForm("tracking_type")
	eventType := c.PostForm("event_type")
	itemID := c.PostForm("item_id")
	duration := c.PostForm("duration")
	percentage := c.PostForm("percentage")

	var TrackingUserData TrackingUserDataStruct
	var ResultTrackingUserData ResultTrackingUserDataStruct

	var TrackingData TrackingDataStruct
	TrackingData.Recommendation_id = trackingID
	TrackingData.Type = trackingType
	TrackingUserData, err := GetTracking(c, itemID, TrackingData)
	if err != nil {
		Sentry_log(err)
		ResultTrackingUserData.Success = false
		c.JSON(http.StatusOK, FormatResultAPI(http.StatusOK, "Invalid", ResultTrackingUserData))
	}

	TrackingUserData.Action = eventType

	var newDuration, newPercentage int
	newDuration, _ = StringToInt(duration)
	newPercentage, _ = StringToInt(percentage)
	if newDuration > 0 {
		TrackingUserData.Duration = int64(newDuration)
		TrackingUserData.Percentage = int64(newPercentage)
	}

	//Check event type not empty

	checkEventType, _ := In_array(eventType, YuspConstantsEvents)
	if eventType != "" && checkEventType && TrackingUserData.Item_id != "" {
		//Push data Kafka
		PushDataToKafka(eventType, TrackingUserData)
		ResultTrackingUserData.Success = true
		c.JSON(http.StatusOK, FormatResultAPI(http.StatusOK, "", ResultTrackingUserData))
		return
	}

	ResultTrackingUserData.Success = false
	c.JSON(http.StatusOK, FormatResultAPI(http.StatusOK, "Invalid", ResultTrackingUserData))
}

func GetTracking(c *gin.Context, itemID string, trackingData TrackingDataStruct) (TrackingUserDataStruct, error) {
	var TrackingUserData TrackingUserDataStruct
	TrackingUserData.Tracking_id = trackingData.Recommendation_id
	TrackingUserData.Tracking_type = trackingData.Type
	TrackingUserData.Item_id = itemID
	TrackingUserData.User_agent = c.GetHeader("User-Agent")
	TrackingUserData.Ip = c.ClientIP()
	TrackingUserData.Token = c.GetHeader("Authorization")
	TrackingUserData.Uri = c.Request.Proto + c.Request.Host + fmt.Sprint(c.Request.URL) //c.GetHeader("Content-Location")
	time, _ := GetCurrentTimeStamp()
	TrackingUserData.Time = time

	return TrackingUserData, nil
}
