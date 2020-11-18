package livetv_v3

import (
	// "encoding/json"
	"errors"
	"sync"
	"time"
	// "fmt"
	. "cm-v5/serv/module"
	. "cm-v5/schema"
	"gopkg.in/mgo.v2/bson"
)

func GetLiveFavorite(userId string, platform string, cache bool) ([]LiveTVObjectStruct, error) {
	var LiveTVObjectOutput = make([]LiveTVObjectStruct, 0)
	listId := GetListIDLivetvFavoriteByUserID(userId)
	if len(listId) <= 0 {
		return LiveTVObjectOutput, errors.New("GetLiveFavorite: Empty data ")
	}

	platformInfo := Platform(platform)
	dataLiveTVTemp, err := GetLiveTVByListID(listId, platformInfo.Id, true)
	if err != nil || len(dataLiveTVTemp) <= 0 {
		Sentry_log(err)
		return LiveTVObjectOutput, err
	}

	dataByte, _ := json.Marshal(dataLiveTVTemp)
	err = json.Unmarshal(dataByte, &LiveTVObjectOutput)
	if err != nil {
		Sentry_log(err)
		return LiveTVObjectOutput, err
	}

	return LiveTVObjectOutput, nil
}

func GetListIDLivetvFavoriteByUserID(userId string) []string {
	// Check Exists Data in cache
	if mRedisUSC.Exists(PREFIX_REDIS_ZRANGE_USC_FAVORITE_V4+userId) == 0 {
		// No have Data in cache
		// Push data from DB
		PushLivetvFavoriteToCache(userId)
		// fmt.Println("Read DB GetListIDLivetvFavoriteByUserID")
	}

	livetvId := mRedisUSC.ZRevRange(PREFIX_REDIS_ZRANGE_USC_FAVORITE_V4+userId, 0, 200)
	return livetvId
}

func AddLivetvFavorite(userId string, livetvIdStr string) error {
	if livetvIdStr == "" || userId == "" {
		return errors.New("Data input not empty")
	}

	var livetvIds []string
	json.Unmarshal([]byte(livetvIdStr), &livetvIds)
	if len(livetvIds) <= 0 {
		return errors.New("Data input not empty")
	}

	var wg sync.WaitGroup
	for k, livetvId := range livetvIds {
		wg.Add(1)

		// Update cache
		status, _ := UpdateLivetvFavoriteCache(userId, livetvId)

		// Update Mongo
		go func(userId string, livetvId string, status int64) {
			UpdateLivetvFavoriteDB(userId, livetvId, status)
			wg.Done()
		}(userId, livetvId, status)

		if (k+1)%10 == 0 {
			wg.Wait()
		}
	}

	return nil
}

func UpdateLivetvFavoriteCache(userId string, livetvId string) (int64, error) {
	var curTime = time.Now().Unix()
	val := mRedisUSC.Incr(PREFIX_REDIS_USC_FAVORITE_V4 + userId + "_" + livetvId)
	if val%2 == 0 {
		// UnFavorite -> Remove cache
		mRedisUSC.ZRem(PREFIX_REDIS_ZRANGE_USC_FAVORITE_V4+userId, livetvId)
	} else {
		// Favorite -> add cache
		mRedisUSC.ZAdd(PREFIX_REDIS_ZRANGE_USC_FAVORITE_V4+userId, float64(curTime), livetvId)
	}
	return val, nil
}

func UpdateLivetvFavoriteDB(userId string, livetvId string, status int64) error {
	// Get data in DB
	session, db, err := GetCollection()
	if err != nil {
		return err
	}
	defer session.Close()

	status = status % 2

	var where = bson.M{
		"user_id":    userId,
		"content_id": livetvId,
	}
	var update = bson.M{
		"user_id":      userId,
		"content_id":   livetvId,
		"entity_type":  5,
		"updated_date": time.Now().Unix(),
		"status":       status,
	}
	_, err = db.C(COLLECTION_USC_USER_FAVORITE).Upsert(where, update)
	if err != nil {
		return err
	}
	return nil
}

func PushLivetvFavoriteToCache(userId string) error {
	/*
	   GetAll se bi timeout khi data qua lon
	   Can optimize them
	*/
	// Get data in DB
	session, db, err := GetCollection()
	if err != nil {
		return err
	}
	defer session.Close()

	var UscUserLivetvFavoriteContentObjs []UscUserLivetvFavoriteContentObjStruct
	var where = bson.M{
		"user_id": userId,
		"status":  1,
	}

	err = db.C(COLLECTION_USC_USER_FAVORITE).Find(where).Sort("-updated_date").All(&UscUserLivetvFavoriteContentObjs)
	if err != nil {
		Sentry_log(err)
		return err
	}

	for _, val := range UscUserLivetvFavoriteContentObjs {
		mRedisUSC.ZAdd(PREFIX_REDIS_ZRANGE_USC_FAVORITE_V4+userId, float64(val.Updated_date), val.Content_id)
	}
	return nil
}
