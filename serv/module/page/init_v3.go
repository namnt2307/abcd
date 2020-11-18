package page

import (
	"errors"
	"fmt"
	"net/http"
	"strings"

	. "cm-v5/schema"
	. "cm-v5/serv/module"
	Packages_premium "cm-v5/serv/module/packages_premium"
	recommendation "cm-v5/serv/module/recommendation"
	seo "cm-v5/serv/module/seo"
	Subscription "cm-v5/serv/module/subscription"
	jsoniter "github.com/json-iterator/go"
	"github.com/gin-gonic/gin"
	"github.com/gosimple/slug"
)

var LIST_SCENARIO_BYWATCHED = map[string]bool{
	"BYWATCHED":          true,
	"BYWATCHED_TVSERIES": true,
	"BYWATCHED_TVSHOW":   true,
	"BYWATCHED_MOVIE":    true,
}

var mRedis RedisModelStruct
var mRedisKV RedisKVModelStruct
var json = jsoniter.ConfigCompatibleWithStandardLibrary

func init() {

}

func GetHomeBannersV3(c *gin.Context) {
	var page_id string = PAGE_HOME_ID
	c.Set("page_id", page_id)
	GetPageBannersV3(c)
}

func GetHomeRibbonsV3(c *gin.Context) {
	var page_id string = PAGE_HOME_ID
	c.Set("page_id", page_id)
	GetPageRibbonsV3(c)
}

func GetPageBannersV3(c *gin.Context) {
	var page_id string = c.Param("page_id")
	if page_id == "" {
		page_id = c.GetString("page_id")
	}

	platform := c.DefaultQuery("platform", "web")
	cacheActive := StringToBool(c.DefaultQuery("cache", "true"))
	limit, err := StringToInt(c.DefaultQuery("limit", "8"))
	if err != nil {
		limit = 8
	}
	limit = 8

	// Hit local data
	var keyCache = LOCAL_DETAIL_RIBBONS_BANNERS_V3 + "_" + platform + "_" + page_id
	if cacheActive {
		valC, err := LocalCache.GetValue(keyCache)
		if err == nil {
			var BannerRibbonsOutputV3 = make([]BannerRibbonsOutputObjectStruct, 0)
			json.Unmarshal([]byte(valC.(string)), &BannerRibbonsOutputV3)
			BannerRibbonsOutputV3 = CheckLocationBannersOutput(BannerRibbonsOutputV3, c)
			c.JSON(http.StatusOK, FormatResultAPI(http.StatusOK, "", BannerRibbonsOutputV3))
			return
		}
	}

	BannerRibbonsOutputV3, err := GetInfoPageBannersV3(page_id, platform, limit, cacheActive)
	if err != nil {
		Sentry_log(err)
		c.JSON(http.StatusBadRequest, FormatResultAPI(http.StatusBadRequest, err.Error(), ""))
		return
	}

	// Write local data
	dataByte, _ := json.Marshal(BannerRibbonsOutputV3)
	LocalCache.SetValue(keyCache, string(dataByte), TTL_LOCALCACHE)

	BannerRibbonsOutputV3 = CheckLocationBannersOutput(BannerRibbonsOutputV3, c)
	c.JSON(http.StatusOK, FormatResultAPI(http.StatusOK, "", BannerRibbonsOutputV3))
}

func GetPageRibbonsV3(c *gin.Context) {
	var page_id string = c.Param("page_id")
	if page_id == "" {
		page_id = c.GetString("page_id")
	}

	platform := c.DefaultQuery("platform", "web")
	cacheActive := StringToBool(c.DefaultQuery("cache", "true"))
	limit, err := StringToInt(c.DefaultQuery("limit", "10"))
	if err != nil {
		limit = 10
	}

	// Hit local data
	var keyCache = LOCAL_PAGE_RIBBONS_V3 + "_" + platform + "_" + page_id

	if cacheActive {
		valC, err := LocalCache.GetValue(keyCache)
		if err == nil {
			var PageRibbonsOutputV3 = make([]PageRibbonsOutputObjectStruct, 0)
			json.Unmarshal([]byte(valC.(string)), &PageRibbonsOutputV3)
			PageRibbonsOutputV3 = GetTrackingDataPageRibbonsV3(platform, PageRibbonsOutputV3)
			PageRibbonsOutputV3 = CheckLocationPageRibbonsOutput(PageRibbonsOutputV3, c)
			c.JSON(http.StatusOK, FormatResultAPI(http.StatusOK, "", PageRibbonsOutputV3))
			return
		}
	}

	PageRibbonsOutputV3, err := GetInfoPageRibbonsV3(page_id, platform, limit, cacheActive)
	if err != nil {
		Sentry_log(err)
		c.JSON(http.StatusBadRequest, FormatResultAPI(http.StatusBadRequest, err.Error(), ""))
		return
	}

	// Write local data
	dataByte, _ := json.Marshal(PageRibbonsOutputV3)
	LocalCache.SetValue(keyCache, string(dataByte), TTL_LOCALCACHE)
	PageRibbonsOutputV3 = GetTrackingDataPageRibbonsV3(platform, PageRibbonsOutputV3)
	PageRibbonsOutputV3 = CheckLocationPageRibbonsOutput(PageRibbonsOutputV3, c)
	c.JSON(http.StatusOK, FormatResultAPI(http.StatusOK, "", PageRibbonsOutputV3))
}

func GetRibbonV3(c *gin.Context) {
	platform := c.DefaultQuery("platform", "web")
	cacheActive := StringToBool(c.DefaultQuery("cache", "true"))
	ribbon_id := c.Param("ribbon_id")
	if ribbon_id == "" {
		ribbon_id = c.GetString("ribbon_id")
	}
	limit, err := StringToInt(c.DefaultQuery("limit", "30"))
	if err != nil {
		limit = 30
	}
	page, err := StringToInt(c.DefaultQuery("page", "0"))
	if err != nil {
		page = 0
	}
	order, err := StringToInt(c.DefaultQuery("order", "3"))
	if err != nil {
		order = 3
	}

	//Lay data tu YUSP
	provider := recommendation.GetProviderRecommend()
	//Test co data sua dieu kien thanh khac (copy database tu mongoDB)
	if provider == RECOMMENDATION_NAME_YUSP {
		user_id := c.GetString("user_id")
		device_id := c.GetString("device_id")
		scenario := recommendation.GetScenariosByRibbonID(ribbon_id)
		if scenario != "" {
			//user chưa login không có ribbbon because you watched
			if (user_id == "" || strings.HasPrefix(user_id, "anonymous_") == true) && LIST_SCENARIO_BYWATCHED[scenario] {
				c.JSON(http.StatusOK, FormatResultAPI(http.StatusOK, "", RibbonDetailOutputObjectStruct{}))
				return
			}

			RibbonsOutput, err := GetRibbonByYUSP(user_id, device_id, scenario, ribbon_id, platform, page, limit, cacheActive)
			if err == nil {
				if len(RibbonsOutput.Ribbon_items) <= 0 {
					c.JSON(http.StatusOK, FormatResultAPI(http.StatusOK, "", RibbonsOutput))
					return
				}

				RibbonOutputResultForCache, err := CheckLocationRibbonsOutput(RibbonsOutput, c)
				if err != nil {
					c.JSON(http.StatusBadRequest, FormatResultAPI(http.StatusBadRequest, err.Error(), ""))
					return
				}

				RibbonOutputResultForCache, err = CheckTagVipOutput(RibbonsOutput, c)
				if err != nil {
					c.JSON(http.StatusBadRequest, FormatResultAPI(http.StatusBadRequest, err.Error(), ""))
					return
				}

				c.JSON(http.StatusOK, FormatResultAPI(http.StatusOK, "", RibbonOutputResultForCache))
				return
			}
		}
	}

	// Hit local data
	var keyCache = LOCAL_DETAIL_RIBBONS_V3 + "_" + platform + "_" + ribbon_id + "_" + fmt.Sprint(page) + "_" + fmt.Sprint(limit) + "_" + fmt.Sprint(order)

	if cacheActive {
		valC, err := LocalCache.GetValue(keyCache)
		if err == nil {
			var RibbonsOutput RibbonDetailOutputObjectStruct
			dataByte, _ := json.Marshal(valC)
			json.Unmarshal([]byte(dataByte), &RibbonsOutput)
			GetTrackingDataRibbonV3(platform, &RibbonsOutput)
			RibbonOutputResultForCache, err := CheckLocationRibbonsOutput(RibbonsOutput, c)
			if err != nil {
				c.JSON(http.StatusBadRequest, FormatResultAPI(http.StatusBadRequest, err.Error(), ""))
				return
			}

			RibbonOutputResultForCache, err = CheckTagVipOutput(RibbonsOutput, c)
			if err != nil {
				c.JSON(http.StatusBadRequest, FormatResultAPI(http.StatusBadRequest, err.Error(), ""))
				return
			}

			c.JSON(http.StatusOK, FormatResultAPI(http.StatusOK, "", RibbonOutputResultForCache))
			return
		}
	}

	RibbonOutput, err := GetRibbonInfoV3(ribbon_id, platform, page, limit, order, cacheActive)
	if err != nil {
		Sentry_log(err)
		c.JSON(http.StatusBadRequest, FormatResultAPI(http.StatusBadRequest, err.Error(), ""))
		return
	}

	// Write local data
	LocalCache.SetValue(keyCache, RibbonOutput, TTL_LOCALCACHE)
	GetTrackingDataRibbonV3(platform, &RibbonOutput)
	RibbonOutputResult, err := CheckLocationRibbonsOutput(RibbonOutput, c)
	if err != nil {
		c.JSON(http.StatusBadRequest, FormatResultAPI(http.StatusBadRequest, err.Error(), ""))
		return
	}

	RibbonOutputResult, err = CheckTagVipOutput(RibbonOutput, c)
	if err != nil {
		c.JSON(http.StatusBadRequest, FormatResultAPI(http.StatusBadRequest, err.Error(), ""))
		return
	}

	c.JSON(http.StatusOK, FormatResultAPI(http.StatusOK, "", RibbonOutputResult))
}

func GetSlugPageRibbonsV3(c *gin.Context) {
	var page_slug string = c.PostForm("page_slug")

	// Check slug exists
	page_id := seo.CheckExistsSEOBySlug(page_slug)
	if page_id != "" {
		c.Set("page_id", page_id)
		GetPageRibbonsV3(c)
		return
	}

	// Get page id by slug
	page_id, err := GetPageIdBySlug(page_slug)
	if err != nil {
		Sentry_log(err)
		c.JSON(http.StatusBadRequest, FormatResultAPI(http.StatusBadRequest, err.Error(), ""))
		return
	}

	c.Set("page_id", page_id)
	GetPageRibbonsV3(c)
}

func GetSlugPageBannersV3(c *gin.Context) {
	var page_slug string = c.PostForm("page_slug")

	// Check slug exists
	page_id := seo.CheckExistsSEOBySlug(page_slug)
	if page_id != "" {
		c.Set("page_id", page_id)
		GetPageBannersV3(c)
		return
	}

	// Get page id by slug
	page_id, err := GetPageIdBySlug(page_slug)
	if err != nil {
		Sentry_log(err)
		c.JSON(http.StatusBadRequest, FormatResultAPI(http.StatusBadRequest, err.Error(), ""))
		return
	}

	c.Set("page_id", page_id)
	GetPageBannersV3(c)
}

func GetSlugRibbonV3(c *gin.Context) {
	var ribbon_slug string = c.PostForm("ribbon_slug")

	// Check slug exists
	ref_id := seo.CheckExistsSEOBySlug(ribbon_slug)
	if ref_id != "" {
		c.Set("ribbon_id", ref_id)
		GetRibbonV3(c)
		return
	}
	// Get page id by slug
	ribbon_id, err := GetRibbonIdBySlugV3(ribbon_slug)
	if err != nil {
		Sentry_log(err)
		c.JSON(http.StatusBadRequest, FormatResultAPI(http.StatusBadRequest, err.Error(), ""))
		return
	}

	c.Set("ribbon_id", ribbon_id)
	GetRibbonV3(c)
}

func GetTrackingDataPageRibbonsV3(platform string, PageRibbonsOutputs []PageRibbonsOutputObjectStruct) []PageRibbonsOutputObjectStruct {
	for k, val := range PageRibbonsOutputs {
		trackingType := platform + "_" + val.Name
		trackingType = slug.Make(trackingType)
		PageRibbonsOutputs[k].Tracking_data = recommendation.GetRandomDefaultTrackingData(strings.ToUpper(trackingType))
	}
	return PageRibbonsOutputs
}

func CheckLocationPageRibbonsOutput(PageRibbonsOutputs []PageRibbonsOutputObjectStruct, c *gin.Context) []PageRibbonsOutputObjectStruct {
	valCheckIpVN := 2
	PageRibbonsOutputsTemps := []PageRibbonsOutputObjectStruct{}
	ipUser, _ := GetClientIPHelper(c.Request, c)
	user_id := c.GetString("user_id")

	for i, s := 0, len(PageRibbonsOutputs); i < s; i++ {
		scenario := recommendation.GetScenariosByRibbonID(PageRibbonsOutputs[i].Id)
		//user chưa login không có ribbbon because you watched
		if scenario != "" && (user_id == "" || strings.HasPrefix(user_id, "anonymous_") == true) && scenario == SCENARIO_BECAUSE_YOU_WATCHED {
			PageRibbonsOutputs[i].Name = ""
		}

		//user chưa login không trả về ribbon because you watched
		// if user_id == "" {
		// 	scenario := recommendation.GetScenariosByRibbonID(PageRibbonsOutputs[i].Id)
		// 	if scenario == SCENARIO_BECAUSE_YOU_WATCHED {
		// 		continue
		// 	}
		// }
		if PageRibbonsOutputs[i].Geo_check == 1 {
			// Content nay  Chi duoc xem o VN
			if valCheckIpVN == 2 {
				valCheckIpVN = 1 // ip vn
				// Check ip user is VN
				if CheckIpIsVN(ipUser) == false {
					valCheckIpVN = 0 // ip nuoc ngoai
				}
			}
			if valCheckIpVN == 1 {
				// IP nuoc ngoai => an item
				PageRibbonsOutputsTemps = append(PageRibbonsOutputsTemps, PageRibbonsOutputs[i])
			}
		} else {
			PageRibbonsOutputsTemps = append(PageRibbonsOutputsTemps, PageRibbonsOutputs[i])
		}
	}

	return PageRibbonsOutputsTemps
}

func GetTrackingDataRibbonV3(platform string, RibbonsOutput *RibbonDetailOutputObjectStruct) {
	trackingType := platform + "_" + RibbonsOutput.Name
	trackingType = slug.Make(trackingType)
	RibbonsOutput.Tracking_data = recommendation.GetRandomDefaultTrackingData(strings.ToUpper(trackingType))
	return
}

func CheckLocationRibbonsOutput(RibbonsOutput RibbonDetailOutputObjectStruct, c *gin.Context) (RibbonsOutputTemps RibbonDetailOutputObjectStruct, err error) {
	RibbonsOutputTemps = RibbonsOutput
	Ribbon_items := []RibbonItemOutputObjectStruct{}
	ipUser, _ := GetClientIPHelper(c.Request, c)

	// Check localtion ribbon
	if RibbonsOutput.Geo_check == 1 {
		if CheckIpIsVN(ipUser) == false {
			// ip nuoc ngoai
			return RibbonsOutputTemps, errors.New("Ribbon not validate your country")
		}
	}

	// Check localtion item
	valCheckIpVN := 2
	total := 0
	for i, s := 0, len(RibbonsOutput.Ribbon_items); i < s; i++ {

		// Check link play (VinadataBuildTokenUrl)
		Hls_link_play := RibbonsOutput.Ribbon_items[i].Link_play.Hls_link_play
		if Hls_link_play != "" {
			RibbonsOutput.Ribbon_items[i].Link_play.Hls_link_play = BuildTokenUrl(Hls_link_play, "", ipUser)
		}

		Dash_link_play := RibbonsOutput.Ribbon_items[i].Link_play.Dash_link_play
		if Dash_link_play != "" {
			RibbonsOutput.Ribbon_items[i].Link_play.Dash_link_play = BuildTokenUrl(Dash_link_play, "", ipUser)
		}

		if RibbonsOutput.Ribbon_items[i].Geo_check == 1 {
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
				Ribbon_items = append(Ribbon_items, RibbonsOutput.Ribbon_items[i])
			}
		} else {
			Ribbon_items = append(Ribbon_items, RibbonsOutput.Ribbon_items[i])
		}
	}

	RibbonsOutputTemps.Ribbon_items = Ribbon_items
	RibbonsOutputTemps.Metadata.Total = RibbonsOutputTemps.Metadata.Total - total

	return RibbonsOutputTemps, nil
}

func CheckTagVipOutput(RibbonsOutput RibbonDetailOutputObjectStruct, c *gin.Context) (RibbonsOutputTemps RibbonDetailOutputObjectStruct, err error) {
	RibbonsOutputTemps = RibbonsOutput

	userId := c.GetString("user_id")

	// Kiểm tra User chưa login => Return
	if strings.HasPrefix(userId, "anonymous_") == true || userId == "" {
		return RibbonsOutputTemps, nil
	}

	var userType bool = false // chua mua goi
	var Sub Subscription.SubcriptionObjectStruct
	SubcriptionsOfUser, err := Sub.GetListByUserId(userId)
	if err == nil && len(SubcriptionsOfUser) > 0 {
		userType = true // Co gói đang active
	}

	// Kiểm tra User chưa mua gói => Return
	if userType == false {
		return RibbonsOutputTemps, nil
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
		Ribbon_items := []RibbonItemOutputObjectStruct{}
		for i := 0; i < len(RibbonsOutput.Ribbon_items); i++ {
			var Ribbon_item = RibbonsOutput.Ribbon_items[i]
			if Ribbon_item.Is_premium == 1 {
				// La content có tag VIP
				Ribbon_item.Is_premium = 0
			}
			Ribbon_items = append(Ribbon_items, Ribbon_item)
		}
		RibbonsOutputTemps.Ribbon_items = Ribbon_items
		return RibbonsOutputTemps, nil
	}

	// User Premium
	// Kiểm tra từng item is_premium = 1, kiểm tra gói của content và gói của user
	listContentPremium, _ := Packages_premium.GetListContentIdInPackagePremium(true)

	Ribbon_items := []RibbonItemOutputObjectStruct{}
	for i := 0; i < len(RibbonsOutput.Ribbon_items); i++ {
		var Ribbon_item = RibbonsOutput.Ribbon_items[i]

		if Ribbon_item.Is_premium == 1 {
			if ok, _ := In_array(Ribbon_item.Id, listContentPremium); ok {
				Ribbon_item.Is_premium = 0
			}
		}
		Ribbon_items = append(Ribbon_items, Ribbon_item)
	}

	RibbonsOutputTemps.Ribbon_items = Ribbon_items
	return RibbonsOutputTemps, nil
}

func CheckLocationBannersOutput(BannersOutput []BannerRibbonsOutputObjectStruct, c *gin.Context) []BannerRibbonsOutputObjectStruct {
	BannersOutputTemps := []BannerRibbonsOutputObjectStruct{}
	ipUser, _ := GetClientIPHelper(c.Request, c)

	// Check localtion item
	valCheckIpVN := 2
	total := 0

	for i, s := 0, len(BannersOutput); i < s; i++ {

		// Check link play (VinadataBuildTokenUrl)
		Hls_link_play := BannersOutput[i].Link_play.Hls_link_play
		if Hls_link_play != "" {
			BannersOutput[i].Link_play.Hls_link_play = BuildTokenUrl(Hls_link_play, "", ipUser)
		}

		Dash_link_play := BannersOutput[i].Link_play.Dash_link_play
		if Dash_link_play != "" {
			BannersOutput[i].Link_play.Dash_link_play = BuildTokenUrl(Dash_link_play, "", ipUser)
		}

		if BannersOutput[i].Geo_check == 1 {
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
				BannersOutputTemps = append(BannersOutputTemps, BannersOutput[i])
			}
		} else {
			BannersOutputTemps = append(BannersOutputTemps, BannersOutput[i])
		}
	}

	return BannersOutputTemps
}
