package content

import (
	"fmt"
	"strings"
	"time"

	. "cm-v5/serv/module"
	. "cm-v5/schema"
)

func GetAdsMySQL(ProviderId string, arrAds []AdsOutputStruct, cache bool, platform string) []AdsOutputStruct {
	var keyCache = KV_PROVIDER_ADS_DETAIL + ProviderId + platform

	if ProviderId == "" {
		return arrAds
	}

	if cache {
		valueCache, err := mRedisKV.GetString(keyCache)
		if err == nil && valueCache != "" {
			err = json.Unmarshal([]byte(valueCache), &arrAds)
			if err == nil {
				return arrAds
			}
		}
	}
	//Connect mysql
	DbMysql, err := ConnectMySQL()
	if err != nil {
		return arrAds
	}
	defer DbMysql.Close()

	// SqlRow := fmt.Sprintf(`
	// SELECT ads.url, ads.repeat, ads.type as t, platform
	// FROM content_provider_ads as ads
	// JOIN content_provider as cp ON cp.id = ads.content_provider_id
	// WHERE (cp.status = 1 and ads.status and ads.content_provider_id="%s" and platform="%s")
	// or (cp.status = 1 and ads.status and ads.content_provider_id="%s" and platform = "general")
	// ORDER BY platform ASC `, ProviderId, platform, ProviderId)

	type AdsOutputMySQLStruct struct {
		Url      string `json:"url" `
		Repeat   int    `json:"repeat" `
		Type     string `json:"type" `
		Platform string `json:"platform" `
		Group    string `json:"group"`
	}
	var arrRawAds []AdsOutputMySQLStruct
	rawQuery := `
	SELECT ads.url, ads.repeat, ads.type as t, platform ,ads.group
	FROM content_provider_ads as ads
	JOIN content_provider as cp ON cp.id = ads.content_provider_id
	WHERE cp.status = 1 and ads.status = 1 and ads.content_provider_id=? AND ( platform=? OR platform = "general")
	ORDER BY ads.group DESC`
	resMySql, err := DbMysql.Query(rawQuery, ProviderId, platform)

	if err == nil {
		for resMySql.Next() {
			var objAds AdsOutputMySQLStruct
			err = resMySql.Scan(&objAds.Url, &objAds.Repeat, &objAds.Type, &objAds.Platform, &objAds.Group)
			if err != nil {
				continue
			} else {
				arrRawAds = append(arrRawAds, objAds)
			}
		}
	}

	// Filter data arrRawAds
	for _, data_raw := range arrRawAds {
		var objAds AdsOutputStruct
		objAds.Repeat = data_raw.Repeat
		objAds.Type = data_raw.Type
		objAds.Group = data_raw.Group
		randomNum := RandNumber(1, 10000)
		randomString := fmt.Sprintf(`%d%d`, time.Now().Unix(), randomNum)
		objAds.Url = strings.Replace(data_raw.Url, "[ADS_RANDOM_INT]", randomString, 1)

		statusExists := false

		for i, data_ad := range arrAds {
			if data_raw.Type == data_ad.Type && objAds.Group == data_ad.Group {
				statusExists = true
				if data_raw.Platform != "general" {
					arrAds[i] = objAds
				}
			}
		}

		// //nếu chưa tồn tại
		if statusExists == false {
			arrAds = append(arrAds, objAds)
		}
	}

	// Write cache
	dataByte, _ := json.Marshal(arrAds)
	mRedisKV.SetString(keyCache, string(dataByte), -1)

	return arrAds
}
