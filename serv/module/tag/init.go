package tag

import (
	"fmt"
	"net/http"
	"strings"

	. "cm-v5/schema"
	. "cm-v5/serv/module"
	recommendation "cm-v5/serv/module/recommendation"
	Subscription "cm-v5/serv/module/subscription"
	Packages_premium "cm-v5/serv/module/packages_premium"

	"github.com/gin-gonic/gin"
	"github.com/gosimple/slug"
	jsoniter "github.com/json-iterator/go"
)

var mRedis RedisModelStruct
var mRedisKV RedisKVModelStruct
var json = jsoniter.ConfigCompatibleWithStandardLibrary

func init() {

}

func GetTagsBySlug(c *gin.Context) {
	tags := c.PostForm("tags")
	tags = strings.Replace(tags, "/", "", -1)
	if tags == "" {
		c.JSON(http.StatusBadRequest, FormatResultAPI(http.StatusBadRequest, "Tags required", "Tags required"))
		return
	}

	//Parse string to array
	var arrSlugTag []string
	arrSlugTag, err := ParseStringToArray(tags)
	if err != nil {
		c.JSON(http.StatusBadRequest, FormatResultAPI(http.StatusBadRequest, err.Error(), "Fail"))
		return
	}

	//Lay danh sach tag id
	listTagID, err := GetTagIDBySlug(arrSlugTag, false)
	if err != nil || len(listTagID) <= 0 {
		c.JSON(http.StatusBadRequest, FormatResultAPI(http.StatusBadRequest, "Tags not found", "Tags not found"))
		return
	}

	justString := strings.Join(listTagID, `","`)
	c.Set("tags", `["`+justString+`"]`)
	GetTagsById(c)
}

func GetTagsById(c *gin.Context) {
	tags := c.GetString("tags")
	if tags == "" {
		tags = c.PostForm("tags")
	}

	page, err := StringToInt(c.DefaultPostForm("page", "0"))
	if err != nil || page > PAGE_MAX {
		page = 0
	}

	limit, err := StringToInt(c.DefaultPostForm("limit", "30"))
	if err != nil || limit > LIMIT_MAX {
		limit = 43
	}

	sort, err := StringToInt(c.DefaultPostForm("sort", "3"))
	if err != nil {
		sort = 3
	}

	platform := c.DefaultQuery("platform", "web")
	cacheActive := StringToBool(c.DefaultQuery("cache", "true"))

	if tags == "" {
		c.JSON(http.StatusBadRequest, FormatResultAPI(http.StatusBadRequest, "Tags required", "Tags required"))
		return
	}

	//Parse string to array
	var arrTagIds []string
	arrTagIds, err = ParseStringToArray(tags)
	if err != nil {
		c.JSON(http.StatusBadRequest, FormatResultAPI(http.StatusBadRequest, err.Error(), "Fail"))
		return
	}
	if len(arrTagIds) <= 0 {
		c.JSON(http.StatusBadRequest, FormatResultAPI(http.StatusBadRequest, "Content not found", "Content not found"))
		return
	}

	// Hit local data
	var keyCache = LOCAL_TAGS + "_" + strings.Join(arrTagIds, "-") + "_" + fmt.Sprint(sort) + "_" + fmt.Sprint(page) + "_" + fmt.Sprint(limit) + "_" + platform
	if cacheActive {
		valC, err := LocalCache.GetValue(keyCache)
		if err == nil {
			var dataVodTags TagsOutputObjectStruct
			dataByte, _ := json.Marshal(valC)
			json.Unmarshal(dataByte, &dataVodTags)
			dataVodTags = GetTrackingDataTags(platform, dataVodTags)
			dataVodTags = CheckLocationTags(dataVodTags, c)
			dataVodTags , _ = CheckTagVipOutput(dataVodTags, c)

			c.JSON(http.StatusOK, FormatResultAPI(http.StatusOK, "", dataVodTags))
			return
		}
	}

	dataVodTags, err := GetVODByTagIds(arrTagIds, sort, page, limit, platform, cacheActive)
	if err != nil {
		c.JSON(http.StatusBadRequest, FormatResultAPI(http.StatusBadRequest, err.Error(), "Content not found"))
		return
	}

	// Write local data
	LocalCache.SetValue(keyCache, dataVodTags, TTL_LOCALCACHE)
	dataVodTags = GetTrackingDataTags(platform, dataVodTags)
	dataVodTags = CheckLocationTags(dataVodTags, c)
	dataVodTags , _ = CheckTagVipOutput(dataVodTags, c)
	

	c.JSON(http.StatusOK, FormatResultAPI(http.StatusOK, "", dataVodTags))
}

func GetTagsCategory(c *gin.Context) {
	cacheActive := StringToBool(c.DefaultQuery("cache", "true"))

	// Hit local data
	var keyCache = LOCAL_TAGS_CATEGORY
	if cacheActive {
		valC, err := LocalCache.GetValue(keyCache)
		if err == nil {
			c.JSON(http.StatusOK, FormatResultAPI(http.StatusOK, "", valC))
			return
		}
	}

	dateTags, err := GetTagsCategorys(cacheActive)
	if err != nil {
		c.JSON(http.StatusBadRequest, FormatResultAPI(http.StatusBadRequest, err.Error(), "Content not found"))
		return
	}

	// Write local data
	LocalCache.SetValue(keyCache, dateTags, TTL_LOCALCACHE)
	c.JSON(http.StatusOK, FormatResultAPI(http.StatusOK, "", dateTags))
}

func GetTrackingDataTags(platform string, TagsOutput TagsOutputObjectStruct) TagsOutputObjectStruct {
	var tagName string
	if len(TagsOutput.Metadata.Tags) > 0 {
		for _, val := range TagsOutput.Metadata.Tags {
			tagName += val.Name + "_"
		}
	}
	trackingType := platform + "_" + tagName + "ITEMS"
	trackingType = slug.Make(trackingType)
	TagsOutput.Tracking_data = recommendation.GetRandomDefaultTrackingData(strings.ToUpper(trackingType))
	return TagsOutput
}

func CheckLocationTags(TagsOutput TagsOutputObjectStruct, c *gin.Context) TagsOutputObjectStruct {
	ipUser, _ := GetClientIPHelper(c.Request, c)
	// Check localtion item
	valCheckIpVN := 2
	total := 0

	var ItemTagsOutputTemps []TagsItemsObjectStruct
	for i, s := 0, len(TagsOutput.Items); i < s; i++ {
		if TagsOutput.Items[i].Geo_check == 1 {
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
				ItemTagsOutputTemps = append(ItemTagsOutputTemps, TagsOutput.Items[i])
			}
		} else {
			ItemTagsOutputTemps = append(ItemTagsOutputTemps, TagsOutput.Items[i])
		}
	}

	TagsOutput.Items = ItemTagsOutputTemps
	TagsOutput.Metadata.Total = TagsOutput.Metadata.Total - total

	return TagsOutput
}

func CheckTagVipOutput(TagsOutput TagsOutputObjectStruct, c *gin.Context) (TagsOutputTemp TagsOutputObjectStruct, err error) {
	TagsOutputTemp = TagsOutput

	userId := c.GetString("user_id")

	// Kiểm tra User chưa login => Return
	if strings.HasPrefix(userId, "anonymous_") == true || userId == "" {
		return TagsOutputTemp, nil
	}

	var userType bool = false // chua mua goi
	var Sub Subscription.SubcriptionObjectStruct
	SubcriptionsOfUser, err := Sub.GetListByUserId(userId)
	if err == nil && len(SubcriptionsOfUser) > 0 {
		userType = true // Co gói đang active
	}

	// Kiểm tra User chưa mua gói => Return
	if userType == false {
		return TagsOutputTemp, nil
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
		var Items []TagsItemsObjectStruct
		for i := 0; i < len(TagsOutput.Items); i++ {
			var Item = TagsOutput.Items[i]
			if Item.Is_premium == 1 {
				// La content có tag VIP
				Item.Is_premium = 0 
			} 
			Items = append(Items, Item)
		}
		TagsOutputTemp.Items = Items
		return TagsOutputTemp, nil
	}

	// User Premium
	// Kiểm tra từng item is_premium = 1, kiểm tra gói của content và gói của user
	listContentPremium , _ := Packages_premium.GetListContentIdInPackagePremium(true)

	Items := []TagsItemsObjectStruct{}
	for i := 0; i < len(TagsOutput.Items); i++ {
		var Item = TagsOutput.Items[i]

		if Item.Is_premium == 1 {
			if ok, _ := In_array(Item.Id, listContentPremium); ok {
				Item.Is_premium = 0
			}
		} 
		Items = append(Items, Item)
	}

	TagsOutputTemp.Items = Items
	return TagsOutputTemp, nil
}
