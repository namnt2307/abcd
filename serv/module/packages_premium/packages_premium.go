package packages_premium

import (
	"fmt"

	. "cm-v5/schema"

	. "cm-v5/serv/module"
)

func GetListContentIdInPackagePremium(cacheActive bool) (listContentPremium []string, err error) {
	keyCache := "KV_LIST_CONTENT_ID_PREMIUM"

	if cacheActive {
		// Read cache
		valueCacheLocal, err := LocalCache.GetValue(keyCache)
		valueCache := fmt.Sprintf("%v", valueCacheLocal)
		if err == nil && valueCache != "" {
			err = json.Unmarshal([]byte(valueCache), &listContentPremium)
			if err == nil {
				return listContentPremium, nil
			}
		}

		valueCache, err = mRedisKV.GetString(keyCache)
		if err == nil && valueCache != "" {
			err = json.Unmarshal([]byte(valueCache), &listContentPremium)
			if err == nil {
				return listContentPremium, nil
			}
		}
	}

	//Connect mysql
	DbMysql, err := ConnectMySQL()
	if err != nil {
		return listContentPremium, err
	}
	defer DbMysql.Close()

	SqlRow := fmt.Sprintf(`
		SELECT DISTINCT bpdc.content_id
		FROM billing_packages_deny_contents as bpdc
		JOIN billing_package_group as bpg ON bpg.id = bpdc.billing_package_group_id
		WHERE bpg.type = 2`)

	dataPackages, err := DbMysql.Query(SqlRow)
	if err != nil {
		return listContentPremium, err
	}

	//fomart result in db
	for dataPackages.Next() {
		var contentId string
		err = dataPackages.Scan(&contentId)
		if err != nil {
			continue
		}
		listContentPremium = append(listContentPremium, contentId)
	}

	// Write cache
	dataByte, _ := json.Marshal(listContentPremium)
	mRedisKV.SetString(keyCache, string(dataByte), TTL_REDIS_5_MINUTE)

	LocalCache.SetValue(keyCache, string(dataByte), TTL_LOCALCACHE)

	return listContentPremium, err
}
