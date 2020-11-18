package drm_auth

import (
	// "fmt"
	"log"
	"net/http"
	"strings"
	"encoding/hex"

	. "cm-v5/serv/module"
	. "cm-v5/schema"
	"github.com/gin-gonic/gin"
	jsoniter "github.com/json-iterator/go"
)

var mRedis RedisModelStruct
var mRedisKV RedisKVModelStruct
var json = jsoniter.ConfigCompatibleWithStandardLibrary
var arrStrBlock []string

func init() {
	BLOCK_BY_STRINGs, _ := CommonConfig.GetString("DRM", "block_by_strings")
	if BLOCK_BY_STRINGs != "" {
		BLOCK_BY_STRINGs = strings.Replace(BLOCK_BY_STRINGs, "...", " ", -1)
		arrStrBlock = strings.Split(BLOCK_BY_STRINGs, ",")
	}

}

func CheckAuthDrmCastlab(c *gin.Context) {

	var infoResponseErr CastlabLicenseResponseError
	infoResponseErr.Message = "not granted"
	infoResponseErr.RedirectUrl = "https://vieon.vn/"

	var dataInput interface{}
	err := c.BindJSON(&dataInput)
	if err != nil {
		c.JSON(http.StatusOK, FormatResultAPI(http.StatusOK, "", infoResponseErr))
		return
	}

	dataByte, err := json.Marshal(dataInput)
	if err != nil {
		c.JSON(http.StatusOK, FormatResultAPI(http.StatusOK, "", infoResponseErr))
		return
	}

	var licenseRequest CastlabLicenseRequest
	err = json.Unmarshal(dataByte, &licenseRequest)
	if err != nil {
		c.JSON(http.StatusOK, FormatResultAPI(http.StatusOK, "", infoResponseErr))
		return
	}

	//check asset id là livetv hay vod
	isLiveTV := strings.Index(licenseRequest.Asset, "livetv_") == 0

	//check permission user with content id
	status := CheckValidLicenseCastLab(licenseRequest, isLiveTV)
	if status == true {
		infoResponse := GenerateInfoResponseCastlab(licenseRequest, isLiveTV)

		//write log success
		go func(dataInput interface{}, infoResponse interface{}) {
			InsertRespLog(dataInput, infoResponse, "castlab")
		}(dataInput, infoResponse)
		c.JSON(http.StatusOK, FormatResultAPI(http.StatusOK, "", infoResponse))
	} else {

		//write log error
		go func(dataInput interface{}, infoResponse interface{}) {
			InsertRespLog(dataInput, infoResponse, "castlab")
		}(dataInput, infoResponseErr)

		c.JSON(http.StatusOK, FormatResultAPI(http.StatusOK, "", infoResponseErr))
	}

}

func CheckAuthDrmVieON(c *gin.Context) {

	contentID := c.DefaultQuery("id", "")
	drmType := c.DefaultQuery("type", "")
	userToken := c.DefaultQuery("token", "")

	strAgent := c.Request.UserAgent()

	// check user agent
	if strings.Index(strAgent, "ExoPlayerLib") >= 0 && strings.Index(strAgent, "app-name") < 0 && strings.Index(strAgent, "ReactNativeVideo") < 0 && strings.Index(strAgent, "VieON") < 0 {
		// block
		decoded, _ := hex.DecodeString("4337344437334636364633463942444135305635304246324242364331313930")
		c.Header("Content-Type", "application/octet-stream")
		c.String(http.StatusOK, "%s", decoded)
		return
	}

	for _ , val := range arrStrBlock {
		if strings.Index(strAgent, val) >= 0 {
			log.Println("val:" , val)
			// block
			decoded, _ := hex.DecodeString("4337344437334636364633463942444135305635304246324242364331313930")
			c.Header("Content-Type", "application/octet-stream")
			c.String(http.StatusOK, "%s", decoded)
			return
		}
	}

	if userToken == "" || contentID == "" || (drmType != "livetv" && drmType != "vod") {
		log.Println("CheckAuthDrmVieON param not correct")
		c.AbortWithStatus(400)
		return
	}

	// decode token
	jwt, err := LocalAuthVerify(userToken)
	if err != nil {
		log.Println("decode jwt error", err)
		c.AbortWithStatus(400)
		return

	}

	var keyCache = KV_CHECK_DRM_VIEON + "_" + contentID + "_" + drmType + "_" + jwt.Subject

	valueCache, err := mRedis.GetString(keyCache)
	if err == nil && valueCache != "" {
		c.Header("Content-Type", "application/octet-stream")
		c.String(http.StatusOK, "%s", valueCache)
		return
	}

	var licenseRequest VieONLicenseRequest
	licenseRequest.Id = contentID
	licenseRequest.Token = userToken
	licenseRequest.Type = drmType
	licenseRequest.User_id = jwt.Subject

	statusUserIsPremium := jwt.Ispremium

	//check permission user with content id
	status := CheckValidLicenseVieON(licenseRequest, statusUserIsPremium)
	if status == true {
		infoResponse := GenerateInfoResponseVieON(licenseRequest)

		if infoResponse != "" {
			//write log success 1hour
			mRedis.SetString(keyCache, infoResponse, TTL_KVCACHE)

			//async
			// go InsertRespLog(licenseRequest, infoResponse, "vieon")

			c.Header("Content-Type", "application/octet-stream")
			c.String(http.StatusOK, "%s", infoResponse)
			return
		}

	} else {
		log.Println("CheckValidLicenseVieON fail")
	}

	// //write log error
	// go func(dataInput interface{}, infoResponse interface{}) {
	// 	InsertRespLog(dataInput, infoResponse, "vieon")
	// }(dataInput, infoResponseErr)

	// c.JSON(http.StatusOK, FormatResultAPI(http.StatusOK, "", infoResponseErr))

	c.AbortWithStatus(400)

}
