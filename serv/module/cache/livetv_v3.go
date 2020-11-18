package cache

import (
	"fmt"
	"net/http"
	"strings"

	livetv "cm-v5/serv/module/livetv_v3"

	. "cm-v5/serv/module"
	// . "cm-v5/schema"
	"github.com/gin-gonic/gin"
)

var mRedisKV RedisKVModelStruct

func UpdateLiveTVGroup(c *gin.Context) {
	for _, val := range platforms {
		livetv.GetLiveTVGroup(val, false)
	}
	c.JSON(http.StatusOK, FormatResultAPI(http.StatusOK, "", "Update cache done"))
}

func UpdateCacheLiveTV(c *gin.Context) {
	livetv_id := c.DefaultQuery("livetv_id", "")
	livetv_group_ids := c.DefaultQuery("livetv_group_ids", "")
	// for _, val := range platforms {
	// 	platformInfo := Platform(val)
	// 	livetv.GetLiveTVByListID([]string{livetv_id}, platformInfo.Id, false)
	// }
	livetv.GetLiveTVByListID([]string{livetv_id}, 0, false)
	//clear cache livetv in group
	if livetv_group_ids != "" {
		arrGroupId := strings.Split(livetv_group_ids, ",")
		fmt.Println("arrGroupId ,", arrGroupId)
		for _, val := range arrGroupId {
			for _, platform := range platforms {
				livetv.GetLiveTVByGroupId(val, platform, 0, 60, false)
			}
		}
	}
	c.JSON(http.StatusOK, FormatResultAPI(http.StatusOK, "", "Update cache done"))
}

func UpdateCacheLiveTVEpg(c *gin.Context) {
	// livetv_id := c.DefaultQuery("livetv_id", "")
	group_id := c.DefaultQuery("group_id", "")
	start, _ := StringToInt(c.DefaultQuery("start", ""))
	// end, _ := StringToInt(c.DefaultQuery("end", ""))

	// strDate := GetDateFromTimestamp(int64(start))
	startTime, endTime := GetStartEndTimeFromTimestamp(int64(start))
	listEpg, _ := livetv.GetEpgByLivetvID(group_id, startTime, endTime, false)

	for _, val := range listEpg {
		livetv.GetDetailEpgBySlug(val.Seo.Url, false)
	}

	// keyCache1 := KV_REDIS_EPG_BY_LIVE_TV_ID + "_" + group_id + "_" + strDate
	// keyCache := KV_REDIS_EPG_BY_LIVE_TV_ID + "_" + group_id + "_" + fmt.Sprint(startTime) + "_" + fmt.Sprint(endTime)

	// livetv.GetEpgByLivetvID(group_id, start, end, false)

	// mRedisKV.Del(keyCache1)
	// mRedisKV.Del(keyCache)
	// fmt.Println("del key", keyCache)

	c.JSON(http.StatusOK, FormatResultAPI(http.StatusOK, "", "Update cache done"))
}
