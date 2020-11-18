package message

import (
	"fmt"
	"time"

	. "cm-v5/serv/module"
	. "cm-v5/schema"
	live_event "cm-v5/serv/module/live_event"
	vod "cm-v5/serv/module/vod"
	"gopkg.in/mgo.v2/bson"
)

type DataResponseStruct struct {
	Status string `json:"status" `
}

func GetMessageByUser(userId string, platform string, page int, limit int, cacheActive bool) (MessageOutputObjectStruct, error) {
	var MessageOutputObject MessageOutputObjectStruct
	MessageItemOutputObjects, err := GetMessageListAllByUser(userId, platform, cacheActive)
	if err != nil {
		return MessageOutputObject, err
	}

	var start, end, total int
	total = len(MessageItemOutputObjects)
	start = page * limit
	end = start + limit
	if end > total {
		end = total
	}

	var MessageItemOutputObjectsWithPage = make([]MessageItemOutputObjectStruct, 0)
	if start < end {
		MessageItemOutputObjectsWithPage = MessageItemOutputObjects[start:end]
	}

	MessageOutputObject.Items = MessageItemOutputObjectsWithPage
	MessageOutputObject.Metadata.Limit = limit
	MessageOutputObject.Metadata.Page = page
	MessageOutputObject.Metadata.Total = total

	return MessageOutputObject, nil
}
func GetMessageListAllByUser(userId string, platform string, cacheActive bool) ([]MessageItemOutputObjectStruct, error) {
	MessageItemOutputObjects, err := GetMessageListAll(userId, platform, cacheActive)
	if err != nil {
		return MessageItemOutputObjects, err
	}
	// Get Msg status by user
	var listMsgStatus = make(map[string]int)
	var keyCache = PREFIX_REDIS_USC_MESSAGE_STATUS + userId
	valueCache, err := mRedisUSC.GetString(keyCache)
	// fmt.Println("mRedisUSC:", err)
	if err == nil && valueCache != "" {
		err = json.Unmarshal([]byte(valueCache), &listMsgStatus)
		if err != nil {
			Sentry_log(err)
		}
	}

	var MessageItemOutputObjectNews []MessageItemOutputObjectStruct
	for _, val := range MessageItemOutputObjects {
		if val.Is_push_all == 0 {
			// Push theo user
			// Check user
			isExists, _ := In_array(userId, val.User_ids)
			if isExists == false {
				continue
			}
		}

		if msgStatus, ok := listMsgStatus[val.Id]; ok {
			val.Status = msgStatus
			if msgStatus == 2 {
				continue
			}
		}
		MessageItemOutputObjectNews = append(MessageItemOutputObjectNews, val)
	}
	return MessageItemOutputObjectNews, nil
}

func GetMessageListAll(userId, platform string, cacheActive bool) ([]MessageItemOutputObjectStruct, error) {
	var MessageItemOutputObjects []MessageItemOutputObjectStruct

	var keyCache = MESSAGE_LIST + "_" + platform
	if cacheActive {
		valueCache, err := mRedis.GetString(keyCache)
		if err == nil && valueCache != "" {
			err = json.Unmarshal([]byte(valueCache), &MessageItemOutputObjects)
			if err == nil {
				return MessageItemOutputObjects, nil
			}
		}
	}

	var MessageObjects []MessageObjectStruct
	var platformDetail = Platform(platform)
	var currentTime = time.Now().Unix()
	var expireTime = currentTime + TTL_REDIS_LV1

	// Connect MongoDB
	session, db, err := GetCollection()
	if err != nil {
		return MessageItemOutputObjects, nil
	}
	defer session.Close()

	var where = bson.M{
		"status":    bson.M{"$in": []int{1, 5}}, //status pushed and pushed by schedule
		"platforms": bson.M{"$in": []string{platform}},
	}

	err = db.C(COLLECTION_NOTIFICATION).Find(where).Sort("-created_at").All(&MessageObjects)
	if err != nil {
		return MessageItemOutputObjects, nil
	}

	for _, val := range MessageObjects {

		var MessageItemOutputObject MessageItemOutputObjectStruct
		MessageItemOutputObject.Id = val.Id
		MessageItemOutputObject.Title = val.Title
		MessageItemOutputObject.Content = val.Content
		MessageItemOutputObject.Expire = val.Expire
		MessageItemOutputObject.Entity_id = val.Entity_id
		MessageItemOutputObject.Entity_type = val.Entity_type
		MessageItemOutputObject.Episode_id = val.Episode_id
		MessageItemOutputObject.Created_at = val.Created_at
		MessageItemOutputObject.Is_push_all = val.Is_push_all
		MessageItemOutputObject.User_ids = val.User_ids
		MessageItemOutputObject.Image = val.Image
		MessageItemOutputObject.More_info = val.More_info

		if val.Expire <= currentTime {
			continue
		}

		if expireTime > val.Expire && val.Expire > 0 {
			expireTime = val.Expire
		}

		if val.Entity_type == 4 {
			/// EPS
			VODDataObjects, _ := vod.GetVODByListID([]string{val.Episode_id}, platformDetail.Id, 1, true)
			if len(VODDataObjects) > 0 {
				MessageItemOutputObject.Seo = VODDataObjects[0].Seo
			}
		} else if val.Entity_type == 1 || val.Entity_type == 2 || val.Entity_type == 3 || val.Entity_type == 6 {
			// Vod
			VODDataObjects, _ := vod.GetVODByListID([]string{val.Entity_id}, platformDetail.Id, 1, true)
			if len(VODDataObjects) > 0 {
				MessageItemOutputObject.Seo = VODDataObjects[0].Seo
			}
		} else if val.Entity_type == 0 {
			LiveEventData, _ := live_event.GetDetailByID(val.Entity_id, platformDetail.Type, true)
			MessageItemOutputObject.Seo = LiveEventData.Seo
		} else {
			// LiveTVData, _ := livetv.GetDetailLiveTV_VTVCab_WithUserInfo(userId, val.Entity_id, platformDetail.Type, true)
			// MessageItemOutputObject.Seo = LiveTVData.Seo
		}
		MessageItemOutputObjects = append(MessageItemOutputObjects, MessageItemOutputObject)
	}

	expireTime = expireTime - currentTime

	if expireTime <= 0 {
		expireTime = int64(TTL_KVCACHE)
	}

	// Write Redis
	dataByte, _ := json.Marshal(MessageItemOutputObjects)
	mRedis.SetString(keyCache, string(dataByte), time.Duration(expireTime))

	return MessageItemOutputObjects, nil
}

func ActionMessage(userId string, messageId string, action string) (DataResponseStruct, error) {
	var dataResponse DataResponseStruct
	dataResponse.Status = "success"

	var listMsgStatus = make(map[string]int)
	var keyCache = PREFIX_REDIS_USC_MESSAGE_STATUS + userId
	valueCache, err := mRedisUSC.GetString(keyCache)
	if err == nil && valueCache != "" {
		json.Unmarshal([]byte(valueCache), &listMsgStatus)
	}

	switch action {
	case "mark_all":
		MessageItemOutputObjects, _ := GetMessageListAllByUser(userId, "web", true)
		for _, val := range MessageItemOutputObjects {
			listMsgStatus[val.Id] = 1
		}
	case "mark":
		if _, ok := listMsgStatus[messageId]; !ok {
			listMsgStatus[messageId] = 1
		}
	case "unmark":
		if vStatus, ok := listMsgStatus[messageId]; ok {
			if vStatus == 1 {
				delete(listMsgStatus, messageId)
			}
		}
	case "delete":
		listMsgStatus[messageId] = 2
	}

	if action == "unmark_all" {
		var listMsgStatusTemp = make(map[string]int)
		for key, val := range listMsgStatus {
			if val == 2 {
				listMsgStatusTemp[key] = 2
			}
		}
		listMsgStatus = listMsgStatusTemp
	}

	dataByte, _ := json.Marshal(listMsgStatus)
	mRedisUSC.SetString(keyCache, string(dataByte), 0)
	return dataResponse, nil
}

func CountMessage(userId string, platform string, cacheActive bool) (interface{}, error) {
	type DataResponseStruct struct {
		Unmark int `json:"unmark" `
		Mark   int `json:"mark" `
	}
	var dataResponse DataResponseStruct

	MessageItemOutputObjects, err := GetMessageListAllByUser(userId, platform, true)
	if err != nil {
		return dataResponse, err
	}

	for _, val := range MessageItemOutputObjects {
		switch val.Status {
		case 0:
			dataResponse.Unmark = dataResponse.Unmark + 1
		case 1:
			dataResponse.Mark = dataResponse.Mark + 1
		}
	}
	return dataResponse, nil
}

func ConvertTimeCreated(timeInput int64) int64 {
	t := time.Unix(timeInput, 0)
	y := t.Year()
	mon := t.Month().String()[:3]
	d := t.Day()
	s := t.Hour()
	h := t.Minute()
	i := t.Second()

	layout := "02/Jan/2006:15:04:05"

	timeOutput := fmt.Sprintf("%02d/%s/%d:%02d:%02d:%d", int(d), mon, y, s, h, i)
	uotput, _ := time.Parse(layout, timeOutput)

	//UTC
	utcLocation, _ := time.LoadLocation("UTC")
	uotput = uotput.In(utcLocation)

	return uotput.Unix()
}
