package cache

import (
	"fmt"
	"net/http"
	"time"

	. "cm-v5/schema"
	. "cm-v5/serv/module"
	page "cm-v5/serv/module/page"

	// "fmt"
	"github.com/gin-gonic/gin"
)


func UpdateCacheRibbonV3(c *gin.Context) {
	var ribbon_id string = c.Param("ribbon_id")
	if ribbon_id == "" {
		ribbon_id = c.GetString("ribbon_id")
	}
	for _, val := range platforms {
		RibbonDetailOutputObject, _ := page.GetRibbonInfoV3(ribbon_id, val, 0, 30, 3, false)

		if len(RibbonDetailOutputObject.Menus) > 0 {
			for _, menu := range RibbonDetailOutputObject.Menus {
				// Check Banner / Ribbon
				if RibbonDetailOutputObject.Type == RIB_TYPE_MASTER_BANNER {
					// Sentry_log_with_msg("UpdateCacheRibbonV3 -> GetListRibbonBannerByPageIDV3 - " + menu.Id + " - " + val)
					// fmt.Println("UpdateCacheRibbonV3 -> GetListRibbonBannerByPageIDV3 - " + menu.Id + " - " + val)
					page.GetRibbonBannerV3(menu.Id, val, false)
				} else {
					// Sentry_log_with_msg("UpdateCacheRibbonV3 -> GetListRibbonByPageIDV3 - " + menu.Id + " - " + val)
					page.GetListRibbonV3(menu.Id, val, false)
				}
			}
		}
	}
	c.JSON(http.StatusOK, FormatResultAPI(http.StatusOK, "", "Update cache done"))
}

func ClearCacheRibbonV3(ribbon_id string) {
	var ribbonSlug string
	for _, val := range platforms {
		RibbonDetailOutputObject, _ := page.GetRibbonInfoV3(ribbon_id, val, 0, 30, 3, false)
		if RibbonDetailOutputObject.Id != "" {
			ribbonSlug = RibbonDetailOutputObject.Seo.Url
		}

		if len(RibbonDetailOutputObject.Menus) > 0 {
			for _, menu := range RibbonDetailOutputObject.Menus {
				// Check Banner / Ribbon
				if RibbonDetailOutputObject.Type == RIB_TYPE_MASTER_BANNER {
					// Sentry_log_with_msg("UpdateCacheRibbonV3 -> GetListRibbonBannerByPageIDV3 - " + menu.Id + " - " + val)
					// fmt.Println("UpdateCacheRibbonV3 -> GetListRibbonBannerByPageIDV3 - " + menu.Id + " - " + val)
					page.GetRibbonBannerV3(menu.Id, val, false)
				} else {
					// Sentry_log_with_msg("UpdateCacheRibbonV3 -> GetListRibbonByPageIDV3 - " + menu.Id + " - " + val)
					page.GetListRibbonV3(menu.Id, val, false)
				}
			}
		}
	}
	if ribbonSlug != "" {
		var keyCache = "DETAIL_RIB_SLUG_V3_" + ribbonSlug
		mRedisKV.Del(keyCache)
		fmt.Println("delete key", keyCache)
	}

	// fmt.Println("Update cache done")
}

func ClearCachePageRibbonV3(menu_id string, typeRib int) {

	for _, val := range platforms {
		if typeRib == RIB_TYPE_MASTER_BANNER {
			// Sentry_log_with_msg("UpdateCacheRibbonV3 -> GetListRibbonBannerByPageIDV3 - " + menu_id + " - " + val)
			page.GetRibbonBannerV3(menu_id, val, false)
		} else {
			// Sentry_log_with_msg("UpdateCacheRibbonV3 -> GetListRibbonByPageIDV3 - " + menu_id + " - " + val)
			page.GetListRibbonV3(menu_id, val, false)
		}
	}
	// fmt.Println("Update cache done")
}

func DelHashKeyCacheRibbon(ribId string) {
	var mRedis RedisModelStruct
	for _, platformStr := range platforms {
		var keyCache = LIST_RIB_ITEM_V3 + "_" + ribId + "_" + platformStr
		mRedis.Del(keyCache)
	}
}

func ClearCacheRibbonItemV3(rib_item_id string, rib_ids []string, status int) {
	var listRibItemID []string
	listRibItemID = append(listRibItemID, rib_item_id)
	shortcutObjectStruct, _ := page.GetListRibItemPushCache(listRibItemID, "web", false)
	for _, platformStr := range platforms {
		platformInfo := Platform(platformStr)

		for _, ribId := range rib_ids {
			var keyCache = LIST_RIB_ITEM_V3 + "_" + ribId + "_" + platformStr

			//xóa cache zrange nếu rib item không publish
			if len(shortcutObjectStruct) == 0 {
				page.RemoveDataToCacheZRange(keyCache, rib_item_id)
				continue
			}
			ribItem := shortcutObjectStruct[0]
			platformExists, _ := In_array(platformInfo.Id, ribItem.Platforms)

			//add cache zrange nếu rib item có status == 1 và có platform tương ứng
			if status == 1 && platformExists {
				page.AddDataToCacheZRange(keyCache, rib_item_id, ribItem.Odr)
			} else {
				page.RemoveDataToCacheZRange(keyCache, rib_item_id)
			}
		}

	}
}

func UpdateCachePageBanner(c *gin.Context) {
	var page_id string = c.Param("page_id")

	go func(page_id string) {
		time.Sleep(60 * time.Second)
		for _, val := range platforms {
			// Sentry_log_with_msg("UpdateCachePageBanner -> GetListRibbonBannerByPageID - " + page_id + " - " + val)
			// page.GetListRibbonBannerByPageID(page_id, val, false)
			page.GetRibbonBannerV3(page_id, val, false)
		}
	}(page_id)

	c.JSON(http.StatusOK, FormatResultAPI(http.StatusOK, "", "Update cache done"))
}

func UpdateCachePageRibbon(c *gin.Context) {
	var page_id string = c.Param("page_id")

	go func(page_id string) {
		time.Sleep(60 * time.Second)
		for _, val := range platforms {
			// Sentry_log_with_msg("UpdateCachePageRibbon -> GetListRibbonByPageID - " + page_id + " - " + val)
			// page.GetListRibbonByPageID(page_id, val, false)
			page.GetListRibbonV3(page_id, val, false)
		}
	}(page_id)

	c.JSON(http.StatusOK, FormatResultAPI(http.StatusOK, "", "Update cache done"))
}
