package geo

import (
	"net/http"
	"strings"

	. "cm-v5/serv/module"
	config "cm-v5/serv/module/config"
	"github.com/gin-gonic/gin"
)

var STATUS_GEOCHECK bool
var mRedis RedisModelStruct

const KEY_REDIS_LIST_SUBDOMAIN = "list_subdomain_cdn"

func init() {
	STATUS_GEOCHECK, _ = CommonConfig.GetBool("GEOCHECK", "status_check")
}

func CheckGeoIp(c *gin.Context) {
	var GeoCheckStruct struct {
		Geo_valid      bool     `json:"geo_valid"`
		Service_status int      `json:"service_status"`
		Sub_domains    []string `json:"sub_domains"`
	}
	ipUser, _ := GetClientIPHelper(c.Request, c)
	// log.Println("CheckGeoIp IP:", string(ipUser))
	if STATUS_GEOCHECK == true { //STATUS_GEOCHECK == true then check ip VN
		GeoCheckStruct.Geo_valid = CheckIpIsVN(ipUser)
	} else { //else pass ip
		GeoCheckStruct.Geo_valid = true
	}

	// Get list sub domain
	valC, err := LocalCache.GetValue("LIST_SUBDOMAIN_CDN")
	if err != nil {
		var listData, err = config.GetConfigByKey(KEY_REDIS_LIST_SUBDOMAIN, "web", true)
		if err == nil {
			GeoCheckStruct.Sub_domains = strings.Split(listData.Data.Value, ",")
			LocalCache.SetValue("LIST_SUBDOMAIN_CDN", listData.Data.Value, TTL_LOCALCACHE)
		}
	} else {
		if tmp, ok := valC.(string); ok {
			GeoCheckStruct.Sub_domains = strings.Split(tmp, ",")
		}
	}

	// Check Service Status
	valC, err = LocalCache.GetValue("SERVICE_STATUS")
	if err == nil {
		if tmp, ok := valC.(int); ok {
			GeoCheckStruct.Service_status = tmp
		}
		c.JSON(http.StatusOK, FormatResultAPI(http.StatusOK, "", GeoCheckStruct))
		return
	}

	GeoCheckStruct.Service_status, _ = mRedis.GetInt("SERVICE_STATUS")
	LocalCache.SetValue("SERVICE_STATUS", GeoCheckStruct.Service_status, TTL_LOCALCACHE)

	c.JSON(http.StatusOK, FormatResultAPI(http.StatusOK, "", GeoCheckStruct))
}
