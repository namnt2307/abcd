package version_data

import (
	"fmt"
	"net/http"
	"time"

	. "cm-v5/serv/module"
	. "cm-v5/schema"
	"github.com/gin-gonic/gin"
	jsoniter "github.com/json-iterator/go"
)

var json = jsoniter.ConfigCompatibleWithStandardLibrary

func GetVersionData(c *gin.Context) {
	cacheActive := StringToBool(c.DefaultQuery("cache", "true"))
	currentTime := time.Now().Unix()
	
	var mRedis RedisModelStruct
	var keyCache string = "VERSION_DATA"
	var VersionDataObj VersionDataObjStruct

	if cacheActive {
		value, err := mRedis.GetString(keyCache)
		if err == nil && value != "" {
			err = json.Unmarshal([]byte(value), &VersionDataObj)
			if err == nil {
				c.JSON(http.StatusOK, FormatResultAPI(http.StatusOK, "", VersionDataObj))
				return
			}
		}
	}

	VersionDataObj.Version_data_smarttv = fmt.Sprint(currentTime)
	dataByte, _ := json.Marshal(VersionDataObj)

	mRedis.SetString(keyCache, string(dataByte), 60*30)
	c.JSON(http.StatusOK, FormatResultAPI(http.StatusOK, "", VersionDataObj))
	return
}


func GetVersionV4(c *gin.Context) {
	var json = jsoniter.ConfigCompatibleWithStandardLibrary
	var mRedis RedisModelStruct
	var keyCache string = "VERSION_DATA_V4"

	var VersionObj = make(map[string]string , 0)
	value, err := mRedis.GetString(keyCache)
	if err == nil && value != "" {
		err = json.Unmarshal([]byte(value), &VersionObj)
		if err == nil {
			fmt.Println("GetVersionV4 Error" , err)
			c.JSON(http.StatusOK, FormatResultAPI(http.StatusOK, "", VersionObj))
			return
		}
	}
	c.JSON(http.StatusOK, FormatResultAPI(http.StatusOK, "", VersionObj))
}
func SetVersionV4(c *gin.Context) {
	var formDate map[string]string
	err := c.ShouldBind(&formDate)
	if err == nil {
		var mRedis RedisModelStruct
		dataByte, _ := json.Marshal(formDate)
		mRedis.SetString("VERSION_DATA_V4" , string(dataByte) , -1)
		c.JSON(http.StatusOK, FormatResultAPI(http.StatusOK, "", formDate))
		return
	}

	c.JSON(http.StatusBadRequest, FormatResultAPI(http.StatusBadRequest, err.Error(), ""))
}

func CheckVersionV4(c *gin.Context) {
	var mRedis RedisModelStruct
	var keyCache string = "VERSION_DATA_V4"
	var VersionObj = make(map[string]string , 0)

	version := c.DefaultQuery("version", "")

	value, err := mRedis.GetString(keyCache)
	if err == nil && value != "" {
		err = json.Unmarshal([]byte(value), &VersionObj)
		if err == nil {
			if val, ok := VersionObj[version]; ok {
				c.JSON(http.StatusOK, FormatResultAPI(http.StatusOK, "", val))
				return
			}
			c.JSON(http.StatusOK, FormatResultAPI(http.StatusOK, "", ""))
			return
		}
		fmt.Println("CheckVersionV4 Error" , err)
	}
	c.JSON(http.StatusOK, FormatResultAPI(http.StatusOK, "", ""))
}