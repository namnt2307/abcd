package config

import (
	"errors"

	. "cm-v5/serv/module"
	. "cm-v5/schema"
)

func GetConfigByKey(key, platform string, cacheActive bool) (ConfigOutputObjectStruct, error) {
	var keyCache = "config_text_" + "_" + key + "_" + platform
	var configObjOutput ConfigOutputObjectStruct

	if key == "" {
		return configObjOutput, errors.New("Empty key")
	}

	if cacheActive {
		valueCache, err := mRedisKV.GetString(keyCache)
		if err == nil {
			err = json.Unmarshal([]byte(valueCache), &configObjOutput.Data)
			return configObjOutput, err
		}
		
	}
	//Connect mysql
	DbMysql, err := ConnectMySQL()
	if err != nil {
		return configObjOutput, err
	}
	defer DbMysql.Close()

	// SqlRow := fmt.Sprintf(`
	// SELECT ads.url, ads.repeat, ads.type as t, platform
	// FROM content_provider_ads as ads
	// JOIN content_provider as cp ON cp.id = ads.content_provider_id
	// WHERE (cp.status = 1 and ads.status and ads.content_provider_id="%s" and platform="%s")
	// or (cp.status = 1 and ads.status and ads.content_provider_id="%s" and platform = "general")
	// ORDER BY platform ASC `, ProviderId, platform, ProviderId)

	type ConfigObjectStruct struct {
		Key      string `json:"key"`
		Value    string `json:"value"`
		Type     string `json:"type"`
		Platform string `json:"platform"`
	}
	var configObj ConfigObjectStruct
	rawQuery := `SELECT config.key, value, type, platform FROM config WHERE config.key=? AND (platform=? OR platform = "general")`
	resMySql, err := DbMysql.Query(rawQuery, key, platform)
	if err == nil {
		for resMySql.Next() {
			var confRaw ConfigObjectStruct
			err = resMySql.Scan(&confRaw.Key, &confRaw.Value, &confRaw.Type, &confRaw.Platform)
			if err != nil {
				continue
			}
			if configObj.Key == "" || confRaw.Platform != "general" {
				configObj = confRaw
			}
		}
	}

	dataByte, _ := json.Marshal(configObj)
	err = json.Unmarshal(dataByte, &configObjOutput.Data)
	if err != nil {
		return configObjOutput, err
	}

	// Write Redis
	dataByte, _ = json.Marshal(configObjOutput.Data)
	mRedisKV.SetString(keyCache, string(dataByte), -1)
	return configObjOutput, nil
}
