package content

import (
	"fmt"
	"net"
	"strings"

	. "cm-v5/serv/module"
	. "cm-v5/schema"
	Package "cm-v5/serv/module/packages"
	rating "cm-v5/serv/module/rating"
	Subscription "cm-v5/serv/module/subscription"
	tracking "cm-v5/serv/module/tracking"
	watchlater "cm-v5/serv/module/watchlater"
)

const (
	SUPER_PREMIUM           = "99"
	PREFIX_USER_FIRST_LOGIN = "user_first_login_"
)

func GetContentTip(contentId string, epsId string, userId string, platform string, cacheActive bool) (VodTipDataOutputObjectStruct, error) {
	var VodTipDataOutput VodTipDataOutputObjectStruct

	// Get content detail by id
	dataContent, err := GetContent(contentId, platform, cacheActive)
	if err != nil {
		return VodTipDataOutput, nil
	}

	dataByte, _ := json.Marshal(dataContent)
	err = json.Unmarshal(dataByte, &VodTipDataOutput)
	if err != nil {
		return VodTipDataOutput, err
	}

	if userId != "" {
		VodTipDataOutput.Is_watchlater = watchlater.CheckContentIsWatchLater(contentId, userId)
	}

	return VodTipDataOutput, nil
}

/**
Lay cac thong tin content co lien quan den user
- rating
- permission
- watchlater
- progress
- .....
*/

func GetContentInfoPersonal(contentId, epsId, userId, platform, modelPlatform, ipUser string, cacheActive bool, statusUserIsPremium string) (ContentDetailPersonalObjStruct, error) {
	var ContentDetailPersonalObj ContentDetailPersonalObjStruct
	ContentDetailPersonalObj.Id = contentId
	ContentDetailPersonalObj.Group_id = contentId
	if epsId != "" {
		ContentDetailPersonalObj.Id = epsId
		ContentDetailPersonalObj.Group_id = contentId
	}

	// Check permission
	GetPermissionContentPersonal(&ContentDetailPersonalObj, userId, ipUser, platform, modelPlatform, statusUserIsPremium)

	// Handle task VIEON-226
	if platform == "mobile_web" {
		ContentDetailPersonalObj.Link_play.Dash_link_play = strings.Replace(ContentDetailPersonalObj.Link_play.Dash_link_play, "/playlist.mpd", "/playlist-mwb.mpd", -1)
		ContentDetailPersonalObj.Link_play.Hls_link_play = strings.Replace(ContentDetailPersonalObj.Link_play.Hls_link_play, "/playlist.m3u8", "/playlist-mwb.m3u8", -1)
	}

	ContentDetailPersonalObj.Link_play.Dash_link_play = BuildTokenUrl(ContentDetailPersonalObj.Link_play.Dash_link_play, "", ipUser)
	ContentDetailPersonalObj.Link_play.Hls_link_play = BuildTokenUrl(ContentDetailPersonalObj.Link_play.Hls_link_play, "", ipUser)

	// Ngày 11/08 anh hải đòi bỏ code này
	// if platform == "lg_tv" {
	// 	ContentDetailPersonalObj.Thumbs.Vtt = ""
	// 	ContentDetailPersonalObj.Thumbs.Image = ""
	// }

	// Handle lite version DT-14042
	ContentDetailPersonalObj.Thumbs.Vtt_lite = strings.Replace(ContentDetailPersonalObj.Thumbs.Vtt, "/thumbs.vtt", "/lite/thumbs.vtt", -1)

	// Check link play (VinadataBuildTokenUrl)
	ContentDetailPersonalObj.Thumbs.Vtt = BuildTokenUrl(ContentDetailPersonalObj.Thumbs.Vtt, "", ipUser)
	ContentDetailPersonalObj.Thumbs.Image = BuildTokenUrl(ContentDetailPersonalObj.Thumbs.Image, "", ipUser)
	ContentDetailPersonalObj.Thumbs.Vtt_lite = BuildTokenUrl(ContentDetailPersonalObj.Thumbs.Vtt_lite, "", ipUser)

	if userId != "" {
		// Check Anonimuos
		if strings.HasPrefix(userId, "anonymous_") == true {
			return ContentDetailPersonalObj, nil
		}

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
			if ContentDetailPersonalObj.Progress >= 2 {
				ContentDetailPersonalObj.Progress = ContentDetailPersonalObj.Progress - 2 // lui 2s
			}
		}

		// Check rating (User - Season)
		ContentDetailPersonalObj.User_rating = rating.GetRatingContentByUser(userId, ContentDetailPersonalObj.Group_id)

		// Get download info
		// GetProfileDownload(&ContentDetailPersonalObj)
	}

	// Start [Request] Thay đổi logic force login với AVOD (https://report.datvietvac.vn/browse/VIEON-381)
	keyCache := PREFIX_USER_FIRST_LOGIN + userId
	if mRedis.Exists(keyCache) > 0 {
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

/**
0: 		Khong duoc xem (khong tra link play)
206: 	Duoc phep xem (co tra link play)
207: 	Chỉ user premium dược phép xem (co tra link play)
*/
func GetPermissionContentPersonal(ContentDetailPersonalObj *ContentDetailPersonalObjStruct, userId, ipUser, platform, modelPlatform, statusUserIsPremium string) {
	ContentDetailPersonalObj.Permission = permissionDefaultVod
	// Set Default Map_profile Premiun
	ContentDetailPersonalObj.Link_play.Map_profile.Full_hd.Is_premium = 1

	// Get content detail by id
	dataContent, err := GetContent(ContentDetailPersonalObj.Id, platform, true)
	if err != nil {
		Sentry_log(err)
		return
	}

	// xe cu ri ti
	if dataContent.Type == VOD_TYPE_EPISODE && dataContent.Group_id != ContentDetailPersonalObj.Group_id {
		return
	}

	//check xem content có sử dụng ads hay không?
	statusEnableAds := dataContent.Enable_ads
	groupAds := dataContent.Custom_ads
	listAds := dataContent.Ads
	var dataContentSs ContentObjOutputStruct
	if dataContent.Type == VOD_TYPE_EPISODE {
		// Get content detail by id
		dataContentSs, err = GetContent(dataContent.Group_id, platform, true)
		if err != nil {
			Sentry_log(err)
			return
		}

		//Đối với episode thì check xem season có cho phép ads hay không
		statusEnableAds = dataContentSs.Enable_ads
		groupAds = dataContentSs.Custom_ads
		listAds = dataContentSs.Ads
	}

	//copy intro outtro
	ContentDetailPersonalObj.Intro.Start = dataContent.Intro_start
	ContentDetailPersonalObj.Intro.End = dataContent.Intro_end
	ContentDetailPersonalObj.Outtro.Start = dataContent.Outtro_start
	ContentDetailPersonalObj.Outtro.End = dataContent.Outtro_end

	var ContentDetailPersonalObjTemp ContentDetailPersonalObjStruct
	dataByte, _ := json.Marshal(dataContent)

	err = json.Unmarshal(dataByte, &ContentDetailPersonalObjTemp)
	if err != nil {
		Sentry_log(err)
		return
	}

	// Khong login quyen xem dua theo quyen Default cua Setting
	if strings.HasPrefix(userId, "anonymous_") == false && userId != "" {
		// Co login duoc phep xem
		ContentDetailPersonalObj.Permission = PERMISSION_VALID

	} else { //user chưa login

		//Process ads for non user
		ContentDetailPersonalObj.Ads = ProcessAdsForUser(userId, false, listAds, groupAds)
	}

	// Kiem tra is_vip
	// EPS co is_vip = 1
	Is_vip := dataContent.Is_vip
	ContentDetailPersonalObj.Is_vip = Is_vip
	if Is_vip != 1 && dataContent.Type == VOD_TYPE_EPISODE {
		Is_vip = dataContentSs.Is_vip
		// Cap nhat theo task DT-13637
		ContentDetailPersonalObj.Is_vip = Is_vip
	}

	// Check permission valid
	if ContentDetailPersonalObj.Permission == PERMISSION_REQUIRE_LOGIN {
		ContentDetailPersonalObj.Link_play.Dash_link_play = ""
		ContentDetailPersonalObj.Link_play.Hls_link_play = ""
		return
	}

	// Kiem tra user đã mua package nào thuộc list package trên chưa
	// Neu co permission valid
	var IsUserPremium = false
	var Sub Subscription.SubcriptionObjectStruct
	SubcriptionsOfUser, err := Sub.GetListByUserId(userId)
	if err == nil && len(SubcriptionsOfUser) > 0 {
		IsUserPremium = true
	} else if err != nil {
		fmt.Println(err)
		Sentry_log(err)
	}

	// Handle task VIEON-250
	var IsUserMbf bool = IsPrivateIP(net.ParseIP(ipUser))
	var IsUserPremiumMbf bool = false

	// User SuperVip mac dinh Premium, khong can mua goi
	if statusUserIsPremium == SUPER_PREMIUM {
		IsUserPremium = true
		IsUserPremiumMbf = true
	}

	// Kiem tra content co can check subs hay k
	if Is_vip == 1 && statusUserIsPremium != SUPER_PREMIUM {
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

			// Lấy thông tin gói group id
			// TienNM - Off vì FE chưa dùng
			// ContentDetailPersonalObj.PackageGroup, _ = Pack.GetPackage(ContentDetailPersonalObj.Group_id)
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
				// Lấy thông tin gói theo id
				// TienNM - Off vì FE chưa dùng
				// ContentDetailPersonalObj.PackageGroup, _ = Pack.GetPackage(ContentDetailPersonalObj.Id)
			}
		}

		for _, vSubObj := range SubcriptionsOfUser {
			// DT-13438 Xem user đăng ký LG VIP được phép xem nội dung VIP
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
		ContentDetailPersonalObj.Link_play.Dash_link_play = ""
		ContentDetailPersonalObj.Link_play.Hls_link_play = ""
		return
	}

	// Handle domain free data MBF
	if IsUserMbf {
		for _, vSubObj := range SubcriptionsOfUser {
			for _, val := range LIST_ID_PACKAGE_MBF {
				if fmt.Sprint(vSubObj.Current_package_id) == val {
					IsUserPremiumMbf = true
					break
				}
			}
		}

		if IsUserPremiumMbf {
			ContentDetailPersonalObjTemp.Link_play.Dash_link_play = ReplaceDomainMBF(ContentDetailPersonalObjTemp.Link_play.Dash_link_play)
			ContentDetailPersonalObjTemp.Link_play.Hls_link_play = ReplaceDomainMBF(ContentDetailPersonalObjTemp.Link_play.Hls_link_play)
		}
	}

	// Sync data
	ContentDetailPersonalObj.Link_play = ContentDetailPersonalObjTemp.Link_play
	ContentDetailPersonalObj.Subtitles = ContentDetailPersonalObjTemp.Subtitles
	ContentDetailPersonalObj.Audios = ContentDetailPersonalObjTemp.Audios
	ContentDetailPersonalObj.Drm_service_name = ContentDetailPersonalObjTemp.Drm_service_name

	//Process ads for user
	if len(ContentDetailPersonalObj.Ads) == 0 && statusEnableAds == 1 {
		ContentDetailPersonalObj.Ads = ProcessAdsForUser(userId, IsUserPremium, listAds, groupAds)
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

	// Handle link play
	ContentDetailPersonalObj.Link_play.Dash_link_play = HandleLinkPlayVieON(ContentDetailPersonalObj.Link_play.Dash_link_play)
	ContentDetailPersonalObj.Link_play.Hls_link_play = HandleLinkPlayVieON(ContentDetailPersonalObj.Link_play.Hls_link_play)

	if strings.Contains(ContentDetailPersonalObj.Link_play.Hls_link_play, "?DVR") == false || platform == "ios" {
		ContentDetailPersonalObj.Thumbs.Vtt = HandleLinkThumbVtt(ContentDetailPersonalObj.Link_play.Hls_link_play)
		ContentDetailPersonalObj.Thumbs.Image = HandleLinkThumbImage(ContentDetailPersonalObj.Link_play.Hls_link_play)
	}

	ContentDetailPersonalObj.Link_play.Dash_link_play = HandleLinkPlayVodByPremium(ContentDetailPersonalObj.Link_play.Dash_link_play, IsUserPremium)
	ContentDetailPersonalObj.Link_play.Hls_link_play = HandleLinkPlayVodByPremium(ContentDetailPersonalObj.Link_play.Hls_link_play, IsUserPremium)

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

	// Set Permision Audio
	for k, val := range ContentDetailPersonalObj.Audios {
		ContentDetailPersonalObj.Audios[k].Is_premium = 1
		ContentDetailPersonalObj.Audios[k].Permission = ContentDetailPersonalObj.Permission
		if val.Is_default == 1 {
			ContentDetailPersonalObj.Audios[k].Is_premium = 0
			ContentDetailPersonalObj.Audios[k].Permission = PERMISSION_VALID
		}

	}

	// Set Permision Subtitle
	for k, val := range ContentDetailPersonalObj.Subtitles {
		ContentDetailPersonalObj.Subtitles[k].Is_premium = 1
		ContentDetailPersonalObj.Subtitles[k].Permission = ContentDetailPersonalObj.Permission
		if val.Is_default == 1 {
			ContentDetailPersonalObj.Subtitles[k].Is_premium = 0
			ContentDetailPersonalObj.Subtitles[k].Permission = PERMISSION_VALID
		}

	}

	return
}

func ProcessAdsForUser(userId string, Is_premium bool, listAds []AdsOutputStruct, groupAds string) []AdsOutputStruct {

	if groupAds == "" {
		groupAds = "default"
	}
	//user VIP empty ads
	adsResult := make([]AdsOutputStruct, 0)

	//nếu không có group nào thì mặc định là group default
	if groupAds == "" {
		groupAds = "default"
	}

	//user chưa login
	if userId == "" || strings.HasPrefix(userId, "anonymous_") == true {
		for _, val := range listAds {
			if strings.Index(val.Type, "nologin-") == 0 && val.Group == groupAds {

				//remove str "nologin-"
				val.Type = strings.Replace(val.Type, "nologin-", "", 1)

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
		}
	} else if !Is_premium { //đã login nhưng ko có PREMIUM
		for _, val := range listAds {
			if strings.Index(val.Type, "nologin-") == -1 && val.Group == groupAds {

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
		}
	}

	return adsResult
}

func CheckUserAgentAndroidAndIOSVersionLock(platform, user_agent string) bool {
	// Check user agent ios
	if platform == "ios" {
		parseUserAgentIos := strings.Split(user_agent, " ")
		if len(parseUserAgentIos) > 0 {
			strParseUserAgentIos := strings.Replace(parseUserAgentIos[0], "VieON/", "", -1)
			versionIos := strParseUserAgentIos[:1]
			if versionIos == "3" || versionIos == "4" {
				return true
			}
		}
		return false
	}

	// Check user agent android
	if platform == "android" {
		parseUserAgentAndroid := strings.Split(user_agent, "_")
		if len(parseUserAgentAndroid) > 1 {
			strParseUserAgentAndroid := parseUserAgentAndroid[1]
			versionAndroid := strParseUserAgentAndroid[:1]
			if versionAndroid == "3" || versionAndroid == "4" {
				return true
			}
		}
	}

	return false
}
