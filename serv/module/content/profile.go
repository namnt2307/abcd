package content

import (
	. "cm-v5/serv/module"
	. "cm-v5/schema"
)

func ContentProfileMySQl(content_id string, cacheActive bool) ([]ContentProfileObjStruct, error) {
	var DownloadProfiles = make([]ContentProfileObjStruct, 0)

	// Cache
	keyCache := PREFIX_REDIS_DOWNLOAD_PROFILE_VOD + content_id
	if cacheActive {
		// Get data in cache
		dataCache, err := mRedis.GetString(keyCache)
		// fmt.Println("TIENNM GetListMenuInfoMySQL Read Redis:" , keyCache, err)
		// fmt.Println("TIENNM GetListMenuInfoMySQL Data Redis:" , dataCache)
		if err == nil && dataCache != "" {
			err = json.Unmarshal([]byte(dataCache), &DownloadProfiles)
			if err == nil {
				return DownloadProfiles, nil
			}
		}
	}

	//Connect mysql
	db_mysql, err := ConnectMySQL()
	if err != nil {
		return DownloadProfiles, err
	}
	defer db_mysql.Close()

	dataObj, err := db_mysql.Query(`select name, resolution, size from entity_vod_profiles where entity_id = ? and status = 1 order by size desc`, content_id)
	if err != nil {
		return DownloadProfiles, err
	}

	//fomart result in db
	for dataObj.Next() {
		var DownloadProfile ContentProfileObjStruct
		err = dataObj.Scan(&DownloadProfile.Name, &DownloadProfile.Resolution, &DownloadProfile.Size)
		if err != nil {
			return DownloadProfiles, err
		}
		DownloadProfiles = append(DownloadProfiles, DownloadProfile)
	}

	// Write cache
	dataByte, _ := json.Marshal(DownloadProfiles)
	mRedis.SetString(keyCache, string(dataByte), TTL_REDIS_LV1)

	return DownloadProfiles, nil
}

func GetProfileDownload(ContentDetailPersonalObj *ContentDetailPersonalObjStruct) {
	ContentDetailPersonalObj.DownloadProfile, _ = ContentProfileMySQl(ContentDetailPersonalObj.Id, true)
	return
}
