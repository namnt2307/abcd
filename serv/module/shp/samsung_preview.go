package shp

import (
	"fmt"
	"log"
	"time"

	. "cm-v5/serv/module"
	. "cm-v5/schema"
	pageRibbon "cm-v5/serv/module/page"
)

var (
	SAMSUNG_SHP_MENUID     string
	SAMSUNG_SHP_LIMIT_ITEM int
	// mRedisKV           RedisKVModelStruct
)

func init() {
	var err error
	SAMSUNG_SHP_MENUID, _ = CommonConfig.GetString("SMARTHUB_PREVIEW", "samsung_menuid")
	SAMSUNG_SHP_LIMIT_ITEM, err = CommonConfig.GetInt("SMARTHUB_PREVIEW", "samsung_limit_item")
	if err != nil || SAMSUNG_SHP_LIMIT_ITEM <= 0 {
		SAMSUNG_SHP_LIMIT_ITEM = 10
	} else {
		SAMSUNG_SHP_LIMIT_ITEM = int(SAMSUNG_SHP_LIMIT_ITEM)
	}

}

func GetDataSamSungPreview(userID string, cache bool) SamsungSmartHubPreview {

	var keyCache = "SAMSUNG_SMARTHUB_PREVIEW_" + SAMSUNG_SHP_MENUID + "_" + fmt.Sprint(SAMSUNG_SHP_LIMIT_ITEM) + "_"
	var dataSHP SamsungSmartHubPreview

	timeNow := time.Now()

	if SAMSUNG_SHP_MENUID == "" {
		log.Println("err read config SAMSUNG_SHP_MENUID")
		return dataSHP
	}

	if cache {
		valueCache, err := mRedisKV.GetString(keyCache)
		if err == nil && valueCache != "" {
			err = json.Unmarshal([]byte(valueCache), &dataSHP)
			if err == nil && dataSHP.Expires > timeNow.Unix() {
				return dataSHP
			}
		}
	}

	//get list ribbon by menu id
	ListRibbon, _ := pageRibbon.GetListRibbonV3(SAMSUNG_SHP_MENUID, "samsung_tv", true)

	for _, val := range ListRibbon {
		var section SamsungSHPSection
		section.Title = val.Name
		//get list ribbon item
		dataRibbon, _ := pageRibbon.GetRibbonInfoV3(val.Id, "samsung_tv", 0, SAMSUNG_SHP_LIMIT_ITEM, 3, true)
		if len(dataRibbon.Ribbon_items) > 0 {
			for i, rib_item := range dataRibbon.Ribbon_items {
				var item SamsungSHPItem
				item.Title = rib_item.Title
				item.Subtitle = rib_item.Subtitle
				item.Image_url = rib_item.Images.Thumbnail
				item.Image_ratio = "16by9"
				item.Is_playable = true
				item.Action_data = fmt.Sprintf(`{"id":"%s","slug":"%s"}`, rib_item.Id, rib_item.Slug)
				item.Position = i + 1
				section.Tiles = append(section.Tiles, item)
			}
		}
		dataSHP.Sections = append(dataSHP.Sections, section)
	}

	dataSHP.Expires = timeNow.Add(time.Hour * 10).Unix()

	dataByte, _ := json.Marshal(dataSHP)
	mRedisKV.SetString(keyCache, string(dataByte), TTL_KVCACHE*5)

	return dataSHP

}
