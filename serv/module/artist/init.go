package artist

import (
	"fmt"
	"net/http"
	"strings"

	. "cm-v5/serv/module"
	. "cm-v5/schema"
	"github.com/gin-gonic/gin"
	jsoniter "github.com/json-iterator/go"

	Subscription "cm-v5/serv/module/subscription"
	Packages_premium "cm-v5/serv/module/packages_premium"
)

var mRedis RedisModelStruct
var mRedisKV RedisKVModelStruct
var json = jsoniter.ConfigCompatibleWithStandardLibrary

func init() {

}

func GetContentArtistBySlug(c *gin.Context) {
	peopleSlug := c.PostForm("people_slug")
	peopleSlug = strings.Replace(peopleSlug, "/nghe-si/", "", -1)
	cacheActive := StringToBool(c.DefaultQuery("cache", "true"))
	//Lay id nghe si
	people_id, err := GetIdArtistBySlug(peopleSlug, cacheActive)
	if err != nil || people_id == "" {
		c.JSON(http.StatusBadRequest, FormatResultAPI(http.StatusBadRequest, err.Error(), "Content not found"))
		return
	}

	c.Set("people_id", people_id)
	GetContentArtistByID(c)
}

func GetContentArtistByID(c *gin.Context) {
	peopleID := c.Param("people_id")
	if peopleID == "" {
		peopleID = c.GetString("people_id")
	}
	page, err := StringToInt(c.DefaultPostForm("page", "0"))
	if err != nil || page > PAGE_MAX {
		page = 0
	}
	limit, err := StringToInt(c.DefaultPostForm("limit", "30"))
	if err != nil || limit > LIMIT_MAX {
		limit = 30
	}
	sort, err := StringToInt(c.DefaultPostForm("sort", "0"))
	if err != nil {
		sort = 0
	}
	platforms := c.DefaultQuery("platform", "web")
	cacheActive := StringToBool(c.DefaultQuery("cache", "true"))

	// Hit local data
	var keyCache = LOCAL_ARTIST_VOD + "_" + fmt.Sprint(page) + "_" + fmt.Sprint(limit) + "_" + fmt.Sprint(sort) + "_" + peopleID + "_" + platforms
	if cacheActive {
		valC, err := LocalCache.GetValue(keyCache)
		if err == nil {
			var dataContentArtist VodArtistOutputObjectStruct
			json.Unmarshal([]byte(valC.(string)), &dataContentArtist)
			dataContentArtist , _ = CheckTagVipOutput(dataContentArtist , c) 
			c.JSON(http.StatusOK, FormatResultAPI(http.StatusOK, "", dataContentArtist))
			return
		}
	}

	dataContentArtist, err := GetVodArtist(peopleID, platforms, page, limit, sort, cacheActive)
	if err != nil {
		c.JSON(http.StatusBadRequest, FormatResultAPI(http.StatusBadRequest, err.Error(), "Content not found"))
		return
	}

	// Write local data
	LocalCache.SetValue(keyCache, dataContentArtist, TTL_LOCALCACHE)

	dataContentArtist , _ = CheckTagVipOutput(dataContentArtist , c) 
	c.JSON(http.StatusOK, FormatResultAPI(http.StatusOK, "", dataContentArtist))
}

func GetArtistBySlug(c *gin.Context) {
	peopleSlug := c.PostForm("people_slug")
	peopleSlug = strings.Replace(peopleSlug, "/nghe-si/", "", -1)
	cacheActive := StringToBool(c.DefaultQuery("cache", "true"))

	//Lay id nghe si
	people_id, err := GetIdArtistBySlug(peopleSlug, cacheActive)
	if err != nil && people_id == "" {
		c.JSON(http.StatusBadRequest, FormatResultAPI(http.StatusBadRequest, err.Error(), "Content not found"))
		return
	}

	c.Set("people_id", people_id)
	GetArtistID(c)
}

func GetArtistID(c *gin.Context) {
	artistID := c.Param("people_id")
	if artistID == "" {
		artistID = c.GetString("people_id")
	}

	cacheActive := StringToBool(c.DefaultQuery("cache", "true"))
	if artistID == "" {
		c.JSON(http.StatusBadRequest, FormatResultAPI(http.StatusBadRequest, "Content not found", "Content not found"))
		return
	}

	// Hit local data
	var keyCache = LOCAL_ARTIST + "_" + artistID
	if cacheActive {
		valC, err := LocalCache.GetValue(keyCache)
		if err == nil {
			c.JSON(http.StatusOK, FormatResultAPI(http.StatusOK, "", valC))
			return
		}
	}

	dataArtist, err := GetInfoArtist(artistID, cacheActive)
	if err != nil {
		c.JSON(http.StatusBadRequest, FormatResultAPI(http.StatusBadRequest, err.Error(), "Content not found"))
		return
	}

	var ArtistOutputObject ArtistOutputObjectStruct
	dataByte, _ := json.Marshal(dataArtist)
	json.Unmarshal(dataByte, &ArtistOutputObject)
	

	// Write local data
	LocalCache.SetValue(keyCache, ArtistOutputObject, TTL_LOCALCACHE)
	c.JSON(http.StatusOK, FormatResultAPI(http.StatusOK, "", ArtistOutputObject))
}

func GetArtistRelatedBySlug(c *gin.Context) {
	peopleSlug := c.PostForm("people_slug")
	peopleSlug = strings.Replace(peopleSlug, "/nghe-si/", "", -1)
	cacheActive := StringToBool(c.DefaultQuery("cache", "true"))

	//Lay id nghe si
	people_id, err := GetIdArtistBySlug(peopleSlug, cacheActive)
	if err != nil || people_id == "" {
		c.JSON(http.StatusBadRequest, FormatResultAPI(http.StatusBadRequest, err.Error(), "Content not found"))
		return
	}

	c.Set("people_id", people_id)
	GetArtistRelatedByID(c)

}

func GetArtistRelatedByID(c *gin.Context) {
	peopleID := c.Param("people_id")
	if peopleID == "" {
		peopleID = c.GetString("people_id")
	}
	cacheActive := StringToBool(c.DefaultQuery("cache", "true"))
	page, err := StringToInt(c.DefaultPostForm("page", "0"))
	if err != nil || page > PAGE_MAX {
		page = 0
	}
	limit, err := StringToInt(c.DefaultPostForm("limit", "5"))
	if err != nil || limit > LIMIT_MAX {
		limit = 5
	}
	platform := c.DefaultQuery("platform", "web")

	// Hit local data
	var keyCache = LOCAL_ARTIST_RELEATED + "_" + peopleID + "_" + fmt.Sprint(page) + "_" + fmt.Sprint(limit) + "_" + platform
	if cacheActive {
		valC, err := LocalCache.GetValue(keyCache)
		if err == nil {
			c.JSON(http.StatusOK, FormatResultAPI(http.StatusOK, "", valC))
			return
		}
	}

	dataArtistRelated, err := GetArtistRelated(peopleID, platform, page, limit, cacheActive)
	if err != nil {
		c.JSON(http.StatusBadRequest, FormatResultAPI(http.StatusBadRequest, err.Error(), "not found"))
		return
	}

	// Write local data
	LocalCache.SetValue(keyCache, dataArtistRelated, TTL_LOCALCACHE)
	c.JSON(http.StatusOK, FormatResultAPI(http.StatusOK, "", dataArtistRelated))
}

func CheckLocationContentArtist(VodArtistOutput VodArtistOutputObjectStruct, c *gin.Context) VodArtistOutputObjectStruct {
	ipUser, _ := GetClientIPHelper(c.Request, c)
	// Check localtion item
	valCheckIpVN := 2
	total := 0

	var ItemsVodArtistTemps []ItemsVodArtistOutputStruct
	for i, s := 0, len(VodArtistOutput.Items); i < s; i++ {
		if VodArtistOutput.Items[i].Geo_check == 1 {
			// Content nay  Chi duoc xem o VN
			if valCheckIpVN == 2 {
				valCheckIpVN = 1 // ip vn
				// Check ip user is VN
				if CheckIpIsVN(ipUser) == false {
					valCheckIpVN = 0 // ip nuoc ngoai
				}
			}
			if valCheckIpVN == 0 {
				total = total + 1
			} else {
				ItemsVodArtistTemps = append(ItemsVodArtistTemps, VodArtistOutput.Items[i])
			}
		} else {
			ItemsVodArtistTemps = append(ItemsVodArtistTemps, VodArtistOutput.Items[i])
		}
	}

	VodArtistOutput.Items = ItemsVodArtistTemps
	VodArtistOutput.Metadata.Total = VodArtistOutput.Metadata.Total - total

	return VodArtistOutput
}


func CheckTagVipOutput(VodArtistOutput VodArtistOutputObjectStruct, c *gin.Context) (VodArtistOutputTemp VodArtistOutputObjectStruct, err error) {
	VodArtistOutputTemp = VodArtistOutput

	userId := c.GetString("user_id")

	// Kiểm tra User chưa login => Return
	if strings.HasPrefix(userId, "anonymous_") == true || userId == "" {
		return VodArtistOutputTemp, nil
	}

	var userType bool = false // chua mua goi
	var Sub Subscription.SubcriptionObjectStruct
	SubcriptionsOfUser, err := Sub.GetListByUserId(userId)
	if err == nil && len(SubcriptionsOfUser) > 0 {
		userType = true // Co gói đang active
	}

	// Kiểm tra User chưa mua gói => Return
	if userType == false {
		return VodArtistOutputTemp, nil
	}
	
	var userIsVipOrPremium bool = false // False: Premium , True: VIP
	for _, val := range SubcriptionsOfUser {
		if val.Type == 1 {
			userIsVipOrPremium = true
			break
		}
	}

	// Kiểm tra User có gói VIP
	// chuyển is_premium = 1 => 0
	if userIsVipOrPremium {
		var Items []ItemsVodArtistOutputStruct
		for i := 0; i < len(VodArtistOutput.Items); i++ {
			var Item = VodArtistOutput.Items[i]
			if Item.Is_premium == 1 {
				// La content có tag VIP
				Item.Is_premium = 0 
			} 
			Items = append(Items, Item)
		}
		VodArtistOutputTemp.Items = Items
		return VodArtistOutputTemp, nil
	}

	// User Premium
	// Kiểm tra từng item is_premium = 1, kiểm tra gói của content và gói của user
	listContentPremium , _ := Packages_premium.GetListContentIdInPackagePremium(true)

	Items := []ItemsVodArtistOutputStruct{}
	for i := 0; i < len(VodArtistOutput.Items); i++ {
		var Item = VodArtistOutput.Items[i]

		if Item.Is_premium == 1 {
			if ok, _ := In_array(Item.Id, listContentPremium); ok {
				Item.Is_premium = 0
			}
		} 
		Items = append(Items, Item)
	}

	VodArtistOutputTemp.Items = Items
	return VodArtistOutputTemp, nil
}
