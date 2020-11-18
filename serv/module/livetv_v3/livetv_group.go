package livetv_v3

import (
	// . "ott-backend-go/serv/module_v4"

	. "cm-v5/serv/module"
	. "cm-v5/schema"
	"gopkg.in/mgo.v2/bson"
)

func GetLiveTVGroup(platform string, cache bool) ([]LiveTVGroupOutputObject, error) {
	var LiveTVGroupOutputs []LiveTVGroupOutputObject
	var keyCache = PREFIX_REDIS_LIVE_TV_GROUP + "_" + platform

	if cache {
		valueCache, err := mRedis.GetString(keyCache)
		if err == nil && valueCache != "" {
			err = json.Unmarshal([]byte(valueCache), &LiveTVGroupOutputs)
			if err == nil {
				return LiveTVGroupOutputs, nil
			}
		}
	}

	// Connect DB
	session, db, err := GetCollection()
	if err != nil {
		Sentry_log(err)
		return LiveTVGroupOutputs, err
	}
	defer session.Close()

	platformInfo := Platform(platform)
	var where = bson.M{
		"status":    bson.M{"$in": []int{STATUS_LIVETV_GROUP_PUBLIC, STATUS_LIVETV_GROUP_FOR_SUPER_PREMIUM}},
		"platforms": platformInfo.Id,
	}

	err = db.C(COLLECTION_LIVE_TV_GROUP).Find(where).Sort("odr").All(&LiveTVGroupOutputs)
	if err != nil {
		Sentry_log(err)
		return LiveTVGroupOutputs, err
	}
	// Write cache
	dataByte, _ := json.Marshal(LiveTVGroupOutputs)
	mRedis.SetString(keyCache, string(dataByte), TTL_REDIS_LV1)

	return LiveTVGroupOutputs, nil
}

//CheckGroupSuperPremium Khanh DT-10720
func CheckGroupSuperPremium(LiveTVGroupOutputs []LiveTVGroupOutputObject, userId string) []LiveTVGroupOutputObject {

	var NewLiveTVGroupOutputs = make([]LiveTVGroupOutputObject, 0)
	if userId != "" {
		infoUser, err := GetInfoUserById(userId)
		if err == nil && infoUser.Is_premium == USER_SUPER_PREMIUM {
			return LiveTVGroupOutputs
		}
	}
	//Loại trừ các nhóm kênh dành cho user super premium
	for _, val := range LiveTVGroupOutputs {
		if val.Status != STATUS_LIVETV_GROUP_FOR_SUPER_PREMIUM { //status 99
			NewLiveTVGroupOutputs = append(NewLiveTVGroupOutputs, val)
		}
	}
	return NewLiveTVGroupOutputs
}
