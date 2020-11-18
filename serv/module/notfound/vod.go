package notfound

import (
	"fmt"
	. "cm-v5/schema"
	. "cm-v5/serv/module"
	pageRibbon "cm-v5/serv/module/page"
	
)

func GetDataNotFoundPreview(col_id , platform string, page ,limit int, cacheActive bool) (RibbonDetailOutputObjectStruct , error) {
	var keyCache = "NOTFOUND_PREVIEW_" + col_id + "_" + fmt.Sprint(page) + "_" + fmt.Sprint(limit)
	var RibbonOutput RibbonDetailOutputObjectStruct

	if cacheActive {
		valueCache, err := mRedisKV.GetString(keyCache)
		if err == nil && valueCache != "" {
			err = json.Unmarshal([]byte(valueCache), &RibbonOutput)
			return RibbonOutput , nil
		}
	}

	RibbonOutput, err := pageRibbon.GetRibbonInfoV3(col_id, platform, page, limit, 3, cacheActive)
	if err != nil {
		return RibbonOutput , err
	}

	for _, val := range RibbonOutput.Ribbon_items {
		RibbonOutputFinal, err := pageRibbon.GetRibbonInfoV3(val.Id, platform, page, limit, 3, cacheActive)
		if err == nil {
			RibbonOutput = RibbonOutputFinal
			break
		}

	}

	dataByte, _ := json.Marshal(RibbonOutput)
	mRedisKV.SetString(keyCache, string(dataByte), TTL_KVCACHE*5)

	return RibbonOutput , nil

}
