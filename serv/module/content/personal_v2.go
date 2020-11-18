package content

import (
	"errors"
	"fmt"
	"log"
	"strings"
	"time"

	. "cm-v5/serv/module"

	. "cm-v5/schema"
	Package "cm-v5/serv/module/packages"
	rating "cm-v5/serv/module/rating"
	Subscription "cm-v5/serv/module/subscription"
	tracking "cm-v5/serv/module/tracking"
	watchlater "cm-v5/serv/module/watchlater"
)

const (
	ERR_MSG_INVALID_PARAM = "Invalid param(s)"
	KEY_RULE_VOD          = "rule_vod"
)

type RulesObjStruct struct {
	Platform  string `json:"platform" `
	Model     string `json:"model" `
	User_type string `json:"user_type" `
}
type RuleStruct struct {
	User_type int
	Ips       []string
}

func MapContentInfoPersonal(contentId, epsId, userId, platform, modelPlatform, ipUser, typeUser string, cacheActive bool) (ContentDetailPersonalObjStruct, error) {
	var ContentDetailPersonalObj ContentDetailPersonalObjStruct
	var err error

	valueCache, err := mRedis.GetString(KEY_RULE_VOD)
	if err == nil && valueCache != "" {
		var Rule []RuleStruct
		err = json.Unmarshal([]byte(valueCache), &Rule)
		if err == nil {
			var userType int = 0
			var Sub Subscription.SubcriptionObjectStruct
			SubcriptionsOfUser, err := Sub.GetListByUserId(userId)
			if err == nil && len(SubcriptionsOfUser) > 0 {
				userType = 1
			}

			if typeUser != "99" {
				typeUser = fmt.Sprint(userType)
			}

			for _, val := range Rule {
				if val.User_type == userType && CheckIpRule(ipUser, val.Ips) {
					ContentDetailPersonalObj, err = GetContentInfoPersonal_V2(contentId, epsId, userId, platform, modelPlatform, ipUser, typeUser, cacheActive)
					return ContentDetailPersonalObj, err
				}
			}
		}
	}

	ContentDetailPersonalObj, err = GetContentInfoPersonal(contentId, epsId, userId, platform, modelPlatform, ipUser, cacheActive, typeUser)
	return ContentDetailPersonalObj, err
}

func GetContentInfoPersonal_V2(contentId, epsId, userId, platform, modelPlatform, ipUser, typeUser string, cacheActive bool) (ContentDetailPersonalObjStruct, error) {
	var ContentDetailPersonalObj ContentDetailPersonalObjStruct
	ContentDetailPersonalObj.Id = contentId
	ContentDetailPersonalObj.Group_id = contentId
	if epsId != "" {
		ContentDetailPersonalObj.Id = epsId
		ContentDetailPersonalObj.Group_id = contentId
	}

	// Set data default
	ContentDetailPersonalObj.Permission = permissionDefaultVod
	ContentDetailPersonalObj.Link_play.Map_profile.Full_hd.Is_premium = 1

	// input content_id = eps_id return error
	if contentId == epsId {
		return ContentDetailPersonalObj, errors.New(ERR_MSG_INVALID_PARAM)
	}

	// Get content detail by id
	dataContent, err := GetContent(ContentDetailPersonalObj.Id, platform, true)
	if err != nil {
		return ContentDetailPersonalObj, err
	}

	// if info video is eps => data id = input eps_id and data group_id = contentId else return error
	if dataContent.Type == VOD_TYPE_EPISODE && dataContent.Id != epsId && dataContent.Group_id != contentId {
		return ContentDetailPersonalObj, errors.New(ERR_MSG_INVALID_PARAM)
	}

	// Kiem tra user đã mua package nào thuộc list package trên chưa
	// Neu co permission valid
	var IsUserPremium = false
	var Sub Subscription.SubcriptionObjectStruct
	SubcriptionsOfUser, err := Sub.GetListByUserId(userId)
	if err == nil && len(SubcriptionsOfUser) > 0 {
		IsUserPremium = true
	} else if err != nil {
		Sentry_log(err)
	}

	// User SuperVip mac dinh Premium, khong can mua goi
	if !IsUserPremium && typeUser == SUPER_PREMIUM {
		IsUserPremium = true
	}

	// Khong login quyen xem dua theo quyen Default cua Setting
	if strings.HasPrefix(userId, "anonymous_") == false && userId != "" {
		// Co login duoc phep xem
		ContentDetailPersonalObj.Permission = PERMISSION_VALID
	}

	//check xem content có sử dụng ads hay không?
	statusEnableAds := dataContent.Enable_ads
	groupAds := dataContent.Custom_ads
	listAds := dataContent.Ads
	Is_vip := dataContent.Is_vip
	if dataContent.Type == VOD_TYPE_EPISODE {
		// Get content detail by id
		dataContentSs, err := GetContent(dataContent.Group_id, platform, true)
		if err != nil {
			Sentry_log(err)
			return ContentDetailPersonalObj, errors.New("Content not exits")
		}
		statusEnableAds = dataContentSs.Enable_ads
		groupAds = dataContentSs.Custom_ads
		listAds = dataContentSs.Ads
		if Is_vip != 1 {
			Is_vip = dataContentSs.Is_vip
		}
	}

	//copy data
	ContentDetailPersonalObj.Intro.Start = dataContent.Intro_start
	ContentDetailPersonalObj.Intro.End = dataContent.Intro_end
	ContentDetailPersonalObj.Outtro.Start = dataContent.Outtro_start
	ContentDetailPersonalObj.Outtro.End = dataContent.Outtro_end
	ContentDetailPersonalObj.Drm_service_name = dataContent.Drm_service_name

	//Process ads for user
	if statusEnableAds == 1 {
		ContentDetailPersonalObj.Ads = ProcessAdsForUser(userId, IsUserPremium, listAds, groupAds)
	}

	// Kiem tra is_vip
	// EPS co is_vip = 1
	ContentDetailPersonalObj.Is_vip = Is_vip

	if userId != "" && !strings.HasPrefix(userId, "anonymous_") {
		// Check watch later (User - Detail)
		if epsId != "" {
			ContentDetailPersonalObj.Is_watchlater = watchlater.CheckContentIsWatchLater(ContentDetailPersonalObj.Group_id, userId)
		} else {
			ContentDetailPersonalObj.Is_watchlater = watchlater.CheckContentIsWatchLater(ContentDetailPersonalObj.Id, userId)
		}

		// Check progress (User - Detail)
		listProgress, err := tracking.GetProgressByListID([]string{ContentDetailPersonalObj.Id}, userId)
		if len(listProgress) > 0 && err == nil {
			ContentDetailPersonalObj.Progress = listProgress[ContentDetailPersonalObj.Id]
			if ContentDetailPersonalObj.Progress >= 5 {
				ContentDetailPersonalObj.Progress = ContentDetailPersonalObj.Progress - 5 // lui 5s
			}
		}

		// Check rating (User - Season)
		ContentDetailPersonalObj.User_rating = rating.GetRatingContentByUser(userId, ContentDetailPersonalObj.Group_id)
	}

	// Check permission valid
	if ContentDetailPersonalObj.Permission == PERMISSION_REQUIRE_LOGIN {
		return ContentDetailPersonalObj, nil
	}

	// Kiem tra content co can check subs hay k
	if Is_vip == 1 && typeUser != SUPER_PREMIUM {
		// Check content in package
		// Xem content co thuoc package nao khong
		// Neu khong permission theo permission default
		// Neu co permission theo permission mua goi / code
		var Pack Package.PackagesObjectStruct
		PackageGroupDenyContents, err := Pack.GetListPackageByDenyContent(ContentDetailPersonalObj.Group_id, true)
		if err == nil && len(PackageGroupDenyContents) > 0 {
			// Content co thuoc 1/n package nao do
			// Set permission theo permission mua goi / code
			if userId == "" {
				ContentDetailPersonalObj.Permission = PERMISSION_REQUIRE_LOGIN // 0
			} else {
				//Check content thuoc goi vip or free
				for _, val := range PackageGroupDenyContents {
					if val.Billing_package_group_id == 10 {
						ContentDetailPersonalObj.Permission = PERMISSION_REQUIRE_PREMIUM
						break
					}
					ContentDetailPersonalObj.Permission = PERMISSION_REQUIRE_PACKAGE
					break
				}
			}
		}

		// Nếu seasion không thuộc gói nào, kiểm tra tập có thuộc gói nào không ??
		// Không phải phim lẻ ( Group_id != Id )
		if len(PackageGroupDenyContents) <= 0 && ContentDetailPersonalObj.Group_id != ContentDetailPersonalObj.Id {
			PackageGroupDenyContents, err = Pack.GetListPackageByDenyContent(ContentDetailPersonalObj.Id, true)
			if err == nil && len(PackageGroupDenyContents) > 0 {
				// Content co thuoc 1/n package nao do
				// Set permission theo permission mua goi / code
				if userId == "" {
					ContentDetailPersonalObj.Permission = PERMISSION_REQUIRE_LOGIN // 0
				} else {
					//Check content thuoc goi vip or free
					for _, val := range PackageGroupDenyContents {
						if val.Billing_package_group_id == 10 {
							ContentDetailPersonalObj.Permission = PERMISSION_REQUIRE_PREMIUM
							break
						}
						ContentDetailPersonalObj.Permission = PERMISSION_REQUIRE_PACKAGE
						break
					}
				}
			}
		}

		for _, vSubObj := range SubcriptionsOfUser {
			// DT-13438 Xem user đăng ký LG VIP hoặc VIP + K+ được phép xem nội dung VIP
			if Package.IsPkgIncludeContentVip(vSubObj.Current_package_id) {
				ContentDetailPersonalObj.Permission = PERMISSION_VALID
				break
			}

			for _, vPGDC := range PackageGroupDenyContents {
				if vSubObj.Current_package_id == vPGDC.Pk_id {
					ContentDetailPersonalObj.Permission = PERMISSION_VALID
					break
				}
			}
		}
	}

	// Check permission valid
	if ContentDetailPersonalObj.Permission != PERMISSION_VALID {
		return ContentDetailPersonalObj, nil
	}

	// Set Permision Profile (Can handle them - config tu setting)
	ContentDetailPersonalObj.Link_play.Map_profile.Four_k.Is_premium = 1
	ContentDetailPersonalObj.Link_play.Map_profile.Four_k.Permission = ContentDetailPersonalObj.Permission
	ContentDetailPersonalObj.Link_play.Map_profile.Full_hd.Is_premium = 1
	ContentDetailPersonalObj.Link_play.Map_profile.Full_hd.Permission = ContentDetailPersonalObj.Permission
	ContentDetailPersonalObj.Link_play.Map_profile.Hd.Permission = ContentDetailPersonalObj.Permission
	ContentDetailPersonalObj.Link_play.Map_profile.Sd.Permission = ContentDetailPersonalObj.Permission

	if IsUserPremium == false {
		ContentDetailPersonalObj.Link_play.Map_profile.Four_k.Permission = PERMISSION_REQUIRE_PACKAGE
		ContentDetailPersonalObj.Link_play.Map_profile.Full_hd.Permission = PERMISSION_REQUIRE_PACKAGE
	}

	// Can optimail them cho nay. Update tam thoi de khong loi POTF
	var videoInfo VideoInfoObjectStruct
	if dataContent.Type != VOD_TYPE_SEASON {
		// Build rules get info video
		var RulesObj RulesObjStruct
		RulesObj.Model = modelPlatform
		RulesObj.Platform = platform
		RulesObj.User_type = typeUser

		t := time.Now()
		videoInfo, err = GetVideoInfo(ContentDetailPersonalObj.Id, ipUser, RulesObj)
		log.Println(fmt.Sprintf("Total time request api POFT: %f ms", time.Since(t).Seconds()))
		if err != nil {
			Sentry_log(err)
			return ContentDetailPersonalObj, nil
		}
	}

	ContentDetailPersonalObj.Link_play.Dash_link_play = videoInfo.Result.Link_play.Dash_link_play
	ContentDetailPersonalObj.Link_play.Hls_link_play = videoInfo.Result.Link_play.Hls_link_play

	// Set humbs
	ContentDetailPersonalObj.Thumbs.Vtt = videoInfo.Result.Thumbs.Vtt
	ContentDetailPersonalObj.Thumbs.Image = videoInfo.Result.Thumbs.Image
	ContentDetailPersonalObj.Thumbs.Vtt_lite = videoInfo.Result.Thumbs.Vtt_lite

	for _, val := range videoInfo.Result.Map_profile {
		if val.Width >= 3840 {
			ContentDetailPersonalObj.Link_play.Map_profile.Four_k.Height = val.Height
			ContentDetailPersonalObj.Link_play.Map_profile.Four_k.Width = val.Width
			ContentDetailPersonalObj.Link_play.Map_profile.Four_k.Height = val.Height
			ContentDetailPersonalObj.Link_play.Map_profile.Four_k.Bandwidth = val.Bandwidth
			ContentDetailPersonalObj.Link_play.Map_profile.Four_k.Name = val.Name
			continue
		} else if val.Width >= 1920 && val.Width < 3840 {
			ContentDetailPersonalObj.Link_play.Map_profile.Full_hd.Height = val.Height
			ContentDetailPersonalObj.Link_play.Map_profile.Full_hd.Width = val.Width
			ContentDetailPersonalObj.Link_play.Map_profile.Full_hd.Height = val.Height
			ContentDetailPersonalObj.Link_play.Map_profile.Full_hd.Bandwidth = val.Bandwidth
			ContentDetailPersonalObj.Link_play.Map_profile.Full_hd.Name = "Full HD"
			continue
		} else if val.Width >= 1280 && val.Width < 1920 {
			ContentDetailPersonalObj.Link_play.Map_profile.Hd.Height = val.Height
			ContentDetailPersonalObj.Link_play.Map_profile.Hd.Width = val.Width
			ContentDetailPersonalObj.Link_play.Map_profile.Hd.Height = val.Height
			ContentDetailPersonalObj.Link_play.Map_profile.Hd.Bandwidth = val.Bandwidth
			ContentDetailPersonalObj.Link_play.Map_profile.Hd.Name = val.Name
			continue
		} else {
			if ContentDetailPersonalObj.Link_play.Map_profile.Sd.Width == 0 || ContentDetailPersonalObj.Link_play.Map_profile.Sd.Width < val.Width {
				ContentDetailPersonalObj.Link_play.Map_profile.Sd.Height = val.Height
				ContentDetailPersonalObj.Link_play.Map_profile.Sd.Width = val.Width
				ContentDetailPersonalObj.Link_play.Map_profile.Sd.Height = val.Height
				ContentDetailPersonalObj.Link_play.Map_profile.Sd.Bandwidth = val.Bandwidth
				ContentDetailPersonalObj.Link_play.Map_profile.Sd.Name = "SD"
			}

		}
	}

	//nếu là smarttv của samsung or lg thì lấy linkplay h265
	if CheckPlatformUsingLinkplayH265(platform, modelPlatform) {
		if ContentDetailPersonalObj.Link_play.H265_dash_link_play != "" {
			ContentDetailPersonalObj.Link_play.Dash_link_play = ContentDetailPersonalObj.Link_play.H265_dash_link_play
			ContentDetailPersonalObj.Link_play.Map_profile.Four_k.Width = 3840
			ContentDetailPersonalObj.Link_play.Map_profile.Four_k.Name = "4K"
		}
		if ContentDetailPersonalObj.Link_play.H265_hls_link_play != "" {
			ContentDetailPersonalObj.Link_play.Hls_link_play = ContentDetailPersonalObj.Link_play.H265_hls_link_play
			ContentDetailPersonalObj.Link_play.Map_profile.Four_k.Width = 3840
			ContentDetailPersonalObj.Link_play.Map_profile.Four_k.Name = "4K"
		}
	}
	ContentDetailPersonalObj.Link_play.H265_dash_link_play = ""
	ContentDetailPersonalObj.Link_play.H265_hls_link_play = ""

	// Set Permision Audio
	ContentDetailPersonalObj.Audios = dataContent.Audios
	for k, val := range ContentDetailPersonalObj.Audios {
		ContentDetailPersonalObj.Audios[k].Is_premium = 1
		ContentDetailPersonalObj.Audios[k].Permission = ContentDetailPersonalObj.Permission
		// if IsUserPremium == false {
		// 	ContentDetailPersonalObj.Audios[k].Permission = PERMISSION_REQUIRE_PACKAGE
		// 	ContentDetailPersonalObj.Audios[k].Permission = PERMISSION_REQUIRE_PACKAGE
		// }
		if val.Is_default == 1 {
			ContentDetailPersonalObj.Audios[k].Is_premium = 0
			ContentDetailPersonalObj.Audios[k].Permission = PERMISSION_VALID
		}

	}

	// Set Permision Subtitle
	ContentDetailPersonalObj.Subtitles = dataContent.Subtitles
	for k, val := range ContentDetailPersonalObj.Subtitles {
		ContentDetailPersonalObj.Subtitles[k].Is_premium = 1
		ContentDetailPersonalObj.Subtitles[k].Permission = ContentDetailPersonalObj.Permission
		// if IsUserPremium == false {
		// 	ContentDetailPersonalObj.Subtitles[k].Permission = PERMISSION_REQUIRE_PACKAGE
		// 	ContentDetailPersonalObj.Subtitles[k].Permission = PERMISSION_REQUIRE_PACKAGE
		// }
		if val.Is_default == 1 {
			ContentDetailPersonalObj.Subtitles[k].Is_premium = 0
			ContentDetailPersonalObj.Subtitles[k].Permission = PERMISSION_VALID
		}

	}

	// Start [Request] Thay đổi logic force login với AVOD (https://report.datvietvac.vn/browse/VIEON-381)
	keyCache := PREFIX_USER_FIRST_LOGIN + userId
	if mRedis.Exists(keyCache) > 0 && len(ContentDetailPersonalObj.Ads) > 0 {
		var AdsOutputObjs []AdsOutputStruct
		for _, val := range ContentDetailPersonalObj.Ads {
			if val.Type == "pre" {
				continue
			}
			AdsOutputObjs = append(AdsOutputObjs, val)
		}

		ContentDetailPersonalObj.Ads = AdsOutputObjs
		mRedis.Del(keyCache)
	}
	// End [Request] Thay đổi logic force login với AVOD (https://report.datvietvac.vn/browse/VIEON-381)
	return ContentDetailPersonalObj, nil
}

func ProcessAdsForUserV2(userId, groupAds string, listAds []AdsOutputStruct) []AdsOutputStruct {

	//user VIP empty ads
	adsResult := make([]AdsOutputStruct, 0)

	//nếu không có group nào thì mặc định là group default
	if groupAds == "" {
		groupAds = "default"
	}

	for _, val := range listAds {
		if strings.Index(val.Type, "nologin-") == 0 && val.Group == groupAds {

			//remove str "nologin-"
			val.Type = strings.Replace(val.Type, "nologin-", "", 1)
		}

		//check tồn tại type
		statusExists := false
		for _, item := range adsResult {
			if item.Type == val.Type {
				statusExists = true
				break
			}
		}
		if statusExists == false {
			adsResult = append(adsResult, val)
		}
	}

	return adsResult
}
