package rating

import (
	"errors"
	"math"

	. "cm-v5/serv/module"
	. "cm-v5/schema"
	vod "cm-v5/serv/module/vod"
	"gopkg.in/mgo.v2/bson"
)

type DataResponseStruct struct {
	Status string `json:"status" `
}

func RatingByUser(userId string, contentId string, platform string, point int) (DataResponseStruct, error) {
	var dataResponse DataResponseStruct
	if contentId == "" || userId == "" {
		dataResponse.Status = "fail"
		return dataResponse, errors.New("Data input not empty")
	}
	if point <= 0 || point > 5 {
		dataResponse.Status = "fail"
		return dataResponse, errors.New("Point not valid")
	}

	// Get Rating info of content
	RatingDetailContent, err := RatingInfoByContentId(contentId, platform)
	if err != nil {
		dataResponse.Status = "fail"
		return dataResponse, err
	}
	// DT-13431 Get from DB if no cache

	var oldRating int = 0

	// Check oldRating by user id
	oldRating = GetRatingContentByUser(userId, contentId)
	if oldRating <= 0 {
		// Chua tung rating
		RatingDetailContent.Total_rate = RatingDetailContent.Total_rate + 1
		RatingDetailContent.Total_point = RatingDetailContent.Total_point + float64(point)
	} else {
		RatingDetailContent.Total_point = RatingDetailContent.Total_point + float64(point) - float64(oldRating)
	}

	// Recount Avg_rate
	if RatingDetailContent.Total_rate > 0 {
		RatingDetailContent.Avg_rate = float64(RatingDetailContent.Total_point) / float64(RatingDetailContent.Total_rate)
	} else {
		RatingDetailContent.Avg_rate = 5
	}

	// handle rating
	RatingDetailContent.Avg_rate = math.Round(RatingDetailContent.Avg_rate*10) / 10

	// Write cache
	UpdateRatingContentByUser(userId, contentId, point)

	var keyCacheContentInfo = PREFIX_REDIS_CONTENT_RATING + contentId

	dataByte, _ := json.Marshal(RatingDetailContent)
	mRedis.SetString(keyCacheContentInfo, string(dataByte), 0)
	// DT-13431 Check insert DB
	InsertRatingMongoDB(RatingDetailContent, contentId)

	dataResponse.Status = "success"
	return dataResponse, nil
}

func RatingInfoByContentId(contentId string, platform string) (RatingDetailContentStruct, error) {
	var RatingDetailContent RatingDetailContentStruct

	// Read cache
	var keyCache = PREFIX_REDIS_CONTENT_RATING + contentId
	valueCache, err := mRedis.GetString(keyCache)
	if err == nil && valueCache != "" {
		err = json.Unmarshal([]byte(valueCache), &RatingDetailContent)
		if err == nil {
			return RatingDetailContent, nil
		}
	}

	var platformDetail = Platform(platform)
	VODDataObjects, err := vod.GetVODByListID([]string{contentId}, platformDetail.Id, 1, true)
	if len(VODDataObjects) <= 0 || err != nil {
		return RatingDetailContent, err
	}

	RatingDetailContent.Avg_rate = VODDataObjects[0].Avg_rate
	RatingDetailContent.Total_rate = VODDataObjects[0].Total_rate
	total_point := VODDataObjects[0].Avg_rate * float64(VODDataObjects[0].Total_rate)
	RatingDetailContent.Total_point = total_point

	// DT-13431 GET from Mongo nếu tồn tại update lại
	RatingDetailContent = GetDataRatingContentInMongo(contentId)

	// Write cache
	dataByte, _ := json.Marshal(RatingDetailContent)
	mRedis.SetString(keyCache, string(dataByte), 0)
	return RatingDetailContent, nil
}

func GetDataRatingContentInMongo(contentId string) RatingDetailContentStruct {
	var RatingDetailContent RatingDetailContentStruct
	session, db, err := GetCollection()
	if err != nil {
		Sentry_log(err)
	}
	defer session.Close()
	var where = bson.M{
		"contentid": contentId,
	}
	var rateObj RateObjectStruct
	err = db.C(COLLECTION_RATE).Find(where).One(&rateObj)
	// valueFromDb := int64(0)
	if err == nil && rateObj.Total_rate > 0 {
		RatingDetailContent.Avg_rate = rateObj.Avg_rate
		RatingDetailContent.Total_rate = rateObj.Total_rate
		RatingDetailContent.Total_point = rateObj.Total_point
	}

	return RatingDetailContent
}

func UpdateRatingContentByUser(userId string, contentId string, point int) error {
	var keyCache = PREFIX_REDIS_USC_RATING + userId
	mRedisUSC.HSet(keyCache, contentId, point)
	return nil
}

func GetRatingContentByUser(userId string, contentId string) int {
	var keyCache = PREFIX_REDIS_USC_RATING + userId
	dataRedis, err := mRedisUSC.HGet(keyCache, contentId)
	if err == nil && dataRedis != "" {
		pointInCache, err := StringToInt(dataRedis)
		if err == nil {
			return pointInCache
		}
	}
	return 0
}

func InsertRatingMongoDB(RatingDetailContent RatingDetailContentStruct, contentId string) error {
	if RatingDetailContent.Total_rate%2 == 0 {
		var rateObj RateObjectStruct
		rateObj.Avg_rate = RatingDetailContent.Avg_rate
		rateObj.Total_rate = RatingDetailContent.Total_rate
		rateObj.Total_point = RatingDetailContent.Total_point
		rateObj.ContentId = contentId
		go func(rateObj RateObjectStruct) {
			//connect db
			session, db, err := GetCollection()
			if err != nil {
				Sentry_log(err)
				return
			}
			defer session.Close()
			var where = bson.M{
				"contentid": rateObj.ContentId,
			}
			_, err = db.C(COLLECTION_RATE).Upsert(where, rateObj)

			if err != nil {
				return
			}
		}(rateObj)

	}
	//Update cache vod_hash
	dataRedis, err := mRedis.HGet(PREFIX_REDIS_HASH_CONTENT, contentId)
	if err != nil {
		// Sentry_log(err)
		return err
	}
	if str, ok := dataRedis.(string); ok {
		var vodDataObject VODDataObjectStruct
		err = json.Unmarshal([]byte(str), &vodDataObject)
		if err != nil {
			return err
		}
		vodDataObject.Avg_rate = ReCalculatorRating(vodDataObject.Min_rate, RatingDetailContent.Avg_rate)
		vodDataObject.Total_rate = RatingDetailContent.Total_rate
		// Set cache
		dataByte, _ := json.Marshal(vodDataObject)
		mRedis.HSet(PREFIX_REDIS_HASH_CONTENT, contentId, string(dataByte))

	}

	// var keyCache = PREFIX_REDIS_USC_RATING + userId
	// mRedisUSC.HSet(keyCache, contentId, point)
	return nil
}
