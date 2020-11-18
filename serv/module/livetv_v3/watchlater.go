package livetv_v3

import (
	"errors"
	"time"

	. "cm-v5/serv/module"
	. "cm-v5/schema"
)

func GetLivetvWatched(userId string, platform string, cacheActive bool) ([]LiveTVObjectStruct, error) {
	var LiveTVObjectOutput = make([]LiveTVObjectStruct, 0)

	listId := GetListIDLivetvWatchedByUserID(userId)
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

func AddLivetvWatched(userId string, LivetvIdsStr string) error {
	if LivetvIdsStr == "" || userId == "" {
		return errors.New("Data input not empty")
	}

	var LivetvIds []string
	json.Unmarshal([]byte(LivetvIdsStr), &LivetvIds)
	if len(LivetvIds) <= 0 {
		return errors.New("Data input not empty")
	}

	for _, LivetvId := range LivetvIds {
		// Update cache
		UpdateLivetvWatchedCache(userId, LivetvId)
	}

	return nil
}

func UpdateLivetvWatchedCache(userId string, LivetvId string) error {
	var curTime = time.Now().Unix()
	err := mRedisUSC.ZAdd(PREFIX_REDIS_ZRANGE_USC_WATCHED_V4+userId, float64(curTime), LivetvId)
	return err
}

func GetListIDLivetvWatchedByUserID(userId string) []string {
	livetvId := mRedisUSC.ZRevRange(PREFIX_REDIS_ZRANGE_USC_WATCHED_V4+userId, 0, 200)
	return livetvId
}
