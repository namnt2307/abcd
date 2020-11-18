package search

import (
	"net/http"
	"strings"

	. "cm-v5/schema"
	. "cm-v5/serv/module"
	Packages_premium "cm-v5/serv/module/packages_premium"
	Subscription "cm-v5/serv/module/subscription"

	"github.com/gin-gonic/gin"
	jsoniter "github.com/json-iterator/go"
	// "log"
)

var mRedis RedisModelStruct
var mRedisUSC RedisUSCModelStruct
var json = jsoniter.ConfigCompatibleWithStandardLibrary

func init() {

}

func Search(c *gin.Context) {
	keyword := c.Query("keyword")
	if keyword == "" {
		c.JSON(http.StatusBadRequest, FormatResultAPI(http.StatusBadRequest, "Bad keyword!", ""))
		return
	}
	tags := c.Query("tags")
	platform := c.DefaultQuery("platform", "web")
	cacheActive := StringToBool(c.DefaultQuery("cache", "true"))
	limit, err := StringToInt(c.DefaultQuery("limit", "30"))
	if err != nil || limit > LIMIT_MAX {
		limit = 30
	}
	page, err := StringToInt(c.DefaultQuery("page", "0"))
	if err != nil || page > PAGE_MAX {
		page = 0
	}
	version := c.Query("version")
	entityType, err := StringToInt(c.DefaultQuery("entity_type", "0"))
	if err != nil {
		entityType = 0
	}
	token := c.GetHeader("Authorization")

	dataSearch, err := GetSearchResult(keyword, tags, platform, version, token, page, limit, entityType, c, cacheActive)
	if err != nil {
		Sentry_log(err)
		c.JSON(http.StatusBadRequest, FormatResultAPI(http.StatusBadRequest, err.Error(), ""))
		return
	}

	dataSearch, _ = CheckTagVipOutput(dataSearch, c)

	c.JSON(http.StatusOK, FormatResultAPI(http.StatusOK, "", dataSearch))
}

func SearchSuggest(c *gin.Context) {
	keyword := c.Query("keyword")
	if keyword == "" {
		c.JSON(http.StatusBadRequest, FormatResultAPI(http.StatusBadRequest, "Bad keyword!", ""))
		return
	}
	tags := c.Query("tags")
	platform := c.DefaultQuery("platform", "web")
	cacheActive := StringToBool(c.DefaultQuery("cache", "true"))
	limit, err := StringToInt(c.DefaultQuery("limit", "15"))
	if err != nil || limit > LIMIT_MAX {
		limit = 30
	}
	page, err := StringToInt(c.DefaultQuery("page", "0"))
	if err != nil || page > PAGE_MAX {
		page = 0
	}
	version := c.Query("version")
	entityType, err := StringToInt(c.DefaultQuery("entity_type", "0"))
	if err != nil {
		entityType = 0
	}
	token := c.GetHeader("Authorization")

	dataSearchSuggest, err := GetSearchSuggest(keyword, tags, platform, version, token, page, limit, entityType, c.Request.Context(), cacheActive)
	if err != nil {
		Sentry_log(err)
		c.JSON(http.StatusBadRequest, FormatResultAPI(http.StatusBadRequest, err.Error(), ""))
		return
	}
	c.JSON(http.StatusOK, FormatResultAPI(http.StatusOK, "", dataSearchSuggest))
}

func GetSearchKeyword(c *gin.Context) {
	platform := c.DefaultQuery("platform", "web")
	cacheActive := StringToBool(c.DefaultQuery("cache", "true"))
	limit, err := StringToInt(c.DefaultQuery("limit", "10"))
	if err != nil || limit > LIMIT_MAX {
		limit = 30
	}
	page, err := StringToInt(c.DefaultQuery("page", "0"))
	if err != nil || page > PAGE_MAX {
		page = 0
	}
	token := c.GetHeader("Authorization")

	dataSeachKeyword, err := GetSeachKeywordResult(platform, token, page, limit, cacheActive)
	if err != nil {
		Sentry_log(err)
		c.JSON(http.StatusBadRequest, FormatResultAPI(http.StatusBadRequest, err.Error(), ""))
		return
	}

	c.JSON(http.StatusOK, FormatResultAPI(http.StatusOK, "", dataSeachKeyword))
}

func AddHistorySearch(c *gin.Context) {
	user_id := c.GetString("user_id")
	content_ids := c.PostForm("content_ids")
	artist_ids := c.PostForm("artist_ids")
	keyword := c.PostForm("keyword")
	data, err := AddHistorySearchByUser(user_id, keyword, content_ids, artist_ids)
	if err != nil {
		Sentry_log(err)
		c.JSON(http.StatusBadRequest, FormatResultAPI(http.StatusBadRequest, err.Error(), ""))
		return
	}
	// VBE-3-Send-Click-EsAws
	platform := c.DefaultQuery("platform", "web")
	position := c.DefaultQuery("position", "")
	request_id := c.DefaultQuery("request_id", "")
	errEs := SendClickEsAwsByUser(user_id, keyword, content_ids, artist_ids, platform, request_id, position)
	if errEs != nil {
		Sentry_log(errEs)
		c.JSON(http.StatusBadRequest, FormatResultAPI(http.StatusBadRequest, errEs.Error(), ""))
	}
	c.JSON(http.StatusOK, FormatResultAPI(http.StatusOK, "", data))
}

func RemoveHistorySearch(c *gin.Context) {
	user_id := c.GetString("user_id")
	data, err := RemoveHistorySearchByUser(user_id)
	if err != nil {
		Sentry_log(err)
		c.JSON(http.StatusBadRequest, FormatResultAPI(http.StatusBadRequest, err.Error(), ""))
		return
	}
	c.JSON(http.StatusOK, FormatResultAPI(http.StatusOK, "", data))
}

func GetHistorySearch(c *gin.Context) {
	platform := c.DefaultQuery("platform", "web")
	user_id := c.GetString("user_id")
	page, err := StringToInt(c.DefaultQuery("page", "0"))
	if err != nil || page > PAGE_MAX {
		page = 0
	}
	limit, err := StringToInt(c.DefaultQuery("limit", "15"))
	if err != nil || limit > LIMIT_MAX {
		limit = 15
	}

	data, err := GetHistorySearchByUserID(user_id, platform, page, limit)
	if err != nil {
		Sentry_log(err)
		c.JSON(http.StatusBadRequest, FormatResultAPI(http.StatusBadRequest, err.Error(), ""))
		return
	}
	c.JSON(http.StatusOK, FormatResultAPI(http.StatusOK, "", data))
}

func CheckLocationSearch(SearchResult SearchResultStruct, c *gin.Context) SearchResultStruct {
	ipUser, _ := GetClientIPHelper(c.Request, c)
	// Check localtion item
	valCheckIpVN := 2
	var total int64

	var SearchItemResultTemps []SearchItemResultStruct
	for i, s := 0, len(SearchResult.Items); i < s; i++ {
		if SearchResult.Items[i].Geo_check == 1 {
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
				SearchItemResultTemps = append(SearchItemResultTemps, SearchResult.Items[i])
			}
		} else {
			SearchItemResultTemps = append(SearchItemResultTemps, SearchResult.Items[i])
		}
	}

	SearchResult.Items = SearchItemResultTemps
	SearchResult.Metadata.Total = SearchResult.Metadata.Total - total

	return SearchResult
}

func CheckTagVipOutput(SearchResult SearchResultStruct, c *gin.Context) (SearchResultTemp SearchResultStruct, err error) {
	SearchResultTemp = SearchResult

	userId := c.GetString("user_id")

	// Kiểm tra User chưa login => Return
	if strings.HasPrefix(userId, "anonymous_") == true || userId == "" {
		return SearchResultTemp, nil
	}

	var userType bool = false // chua mua goi
	var Sub Subscription.SubcriptionObjectStruct
	SubcriptionsOfUser, err := Sub.GetListByUserId(userId)
	if err == nil && len(SubcriptionsOfUser) > 0 {
		userType = true // Co gói đang active
	}

	// Kiểm tra User chưa mua gói => Return
	if userType == false {
		return SearchResultTemp, nil
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
		var Items []SearchItemResultStruct
		for i := 0; i < len(SearchResult.Items); i++ {
			var Item = SearchResult.Items[i]
			if Item.Is_premium == 1 {
				// La content có tag VIP
				Item.Is_premium = 0
			}
			Items = append(Items, Item)
		}
		SearchResultTemp.Items = Items
		return SearchResultTemp, nil
	}

	// User Premium
	// Kiểm tra từng item is_premium = 1, kiểm tra gói của content và gói của user
	listContentPremium, _ := Packages_premium.GetListContentIdInPackagePremium(true)

	Items := []SearchItemResultStruct{}
	for i := 0; i < len(SearchResult.Items); i++ {
		var Item = SearchResult.Items[i]

		if Item.Is_premium == 1 {
			if ok, _ := In_array(Item.Id, listContentPremium); ok {
				Item.Is_premium = 0
			}
		}
		Items = append(Items, Item)
	}

	SearchResultTemp.Items = Items
	return SearchResultTemp, nil
}
