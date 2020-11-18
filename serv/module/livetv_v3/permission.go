package livetv_v3

import (
	"fmt"
	"net"
	"strings"

	. "cm-v5/serv/module"
	. "cm-v5/schema"
	"cm-v5/serv/module/packages"

	Fingering "cm-v5/serv/module/fingering"
	Kplus "cm-v5/serv/module/kplus"
	Subscription "cm-v5/serv/module/subscription"
)

// Không bat login
var permissionDefault = 206

func init() {
	permissionDefault, _ = CommonConfig.GetInt("GLOBAL_CONFIG", "permission_default")
}

func GetFavoriteListLiveTV(ListLiveTV []LiveTVObjectStruct, userId string) []LiveTVObjectStruct {
	if strings.HasPrefix(userId, "anonymous_") == true || userId == "" {
		// Co login duoc phep xem
		return ListLiveTV
	}

	// var PackageGroupObject = make([]PackageGroupObjectStruct, 0)
	listId := GetListIDLivetvFavoriteByUserID(userId)
	// fmt.Println("listId: " , listId)
	if len(listId) > 0 {
		for k, val := range ListLiveTV {
			for _, id := range listId {
				if val.Id == id {
					ListLiveTV[k].IsFavorite = true
				}
			}
		}
	}
	return ListLiveTV
}

func GetPermissionListLiveTVDetail(DetailLiveTV *DetailLiveTVObjectOutputStruct, userId, ipUser, tokenUser string, statusUserIsPremium int) {
	DetailLiveTV.Permission = permissionDefault
	if strings.HasPrefix(userId, "anonymous_") == false && userId != "" {
		// Co login duoc phep xem
		DetailLiveTV.Permission = PERMISSION_VALID
	}

	listId := GetListIDLivetvFavoriteByUserID(userId)
	if len(listId) > 0 {
		for _, id := range listId {
			if DetailLiveTV.Id == id {
				DetailLiveTV.IsFavorite = true
			}
		}
	}

	// KT Kenh co thuoc package nao khong - lay danh sach package id
	if len(DetailLiveTV.PackageGroup) > 0 {
		DetailLiveTV.Permission = PERMISSION_REQUIRE_PACKAGE // 208
		if userId == "" {
			DetailLiveTV.Permission = PERMISSION_REQUIRE_LOGIN // 0
		} else {
			if ok, _ := packages.CheckPermissionUser(DetailLiveTV.Id, userId, statusUserIsPremium); ok {
				DetailLiveTV.Permission = PERMISSION_VALID // 206
			}
		}
	}

	if DetailLiveTV.Permission != PERMISSION_VALID {
		DetailLiveTV.Hls_link_play = ""
		DetailLiveTV.Dash_link_play = ""
		DetailLiveTV.Programme.Hls_link_play = ""
		DetailLiveTV.Programme.Dash_link_play = ""
		return
	}

	//get current epg and check permission
	dataCurrentEpg, _ := GetEPGCurrent(DetailLiveTV.Id)
	if dataCurrentEpg.Id != "" && dataCurrentEpg.Ott_disabled == 1 {
		DetailLiveTV.Permission = PERMISSION_NOT_ALLOW
	}

	// Handle task VIEON-250
	var IsUserPremiumMbf bool = false
	var IsUserMbf bool = IsPrivateIP(net.ParseIP(ipUser))

	if statusUserIsPremium == USER_SUPER_PREMIUM {
		IsUserPremiumMbf = true
	}

	if IsUserMbf {
		var Sub Subscription.SubcriptionObjectStruct
		SubcriptionsOfUser, _ := Sub.GetListByUserId(userId)
		for _, vSubObj := range SubcriptionsOfUser {
			for _, val := range LIST_ID_PACKAGE_MBF {
				if fmt.Sprint(vSubObj.Current_package_id) == val {
					IsUserPremiumMbf = true
					break
				}
			}
		}

		//Replace domain mbf
		if IsUserPremiumMbf {
			DetailLiveTV.Hls_link_play = ReplaceDomainMBF(DetailLiveTV.Hls_link_play)
			DetailLiveTV.Dash_link_play = ReplaceDomainMBF(DetailLiveTV.Dash_link_play)

			DetailLiveTV.Programme.Hls_link_play = ReplaceDomainMBF(DetailLiveTV.Programme.Hls_link_play)
			DetailLiveTV.Programme.Dash_link_play = ReplaceDomainMBF(DetailLiveTV.Programme.Dash_link_play)
		}
	}

	//check drm set is vieon
	//Khanh DT-11862
	if DetailLiveTV.Drm_service_name == "vieon" {

		DetailLiveTV.Hls_link_play = BuildTokenUrl(DetailLiveTV.Hls_link_play, "LiveTV", ipUser)
		DetailLiveTV.Dash_link_play = BuildTokenUrl(DetailLiveTV.Dash_link_play, "LiveTV", ipUser)

		DetailLiveTV.Programme.Hls_link_play = BuildTokenUrl(DetailLiveTV.Programme.Hls_link_play, "LiveTV", ipUser)
		DetailLiveTV.Programme.Dash_link_play = BuildTokenUrl(DetailLiveTV.Programme.Dash_link_play, "LiveTV", ipUser)

		DetailLiveTV.Hls_link_play = VieONBuildDrmTokenUrl(DetailLiveTV.Hls_link_play, DetailLiveTV.Id, "livetv", ipUser, tokenUser)
		DetailLiveTV.Dash_link_play = VieONBuildDrmTokenUrl(DetailLiveTV.Dash_link_play, DetailLiveTV.Id, "livetv", ipUser, tokenUser)

		DetailLiveTV.Programme.Hls_link_play = VieONBuildDrmTokenUrl(DetailLiveTV.Programme.Hls_link_play, DetailLiveTV.Id, "livetv", ipUser, tokenUser)
		DetailLiveTV.Programme.Dash_link_play = VieONBuildDrmTokenUrl(DetailLiveTV.Programme.Dash_link_play, DetailLiveTV.Id, "livetv", ipUser, tokenUser)

	} else {
		DetailLiveTV.Hls_link_play = BuildTokenUrl(DetailLiveTV.Hls_link_play, "LiveTV", ipUser)
		DetailLiveTV.Dash_link_play = BuildTokenUrl(DetailLiveTV.Dash_link_play, "LiveTV", ipUser)

		DetailLiveTV.Programme.Hls_link_play = BuildTokenUrl(DetailLiveTV.Programme.Hls_link_play, "LiveTV", ipUser)
		DetailLiveTV.Programme.Dash_link_play = BuildTokenUrl(DetailLiveTV.Programme.Dash_link_play, "LiveTV", ipUser)
	}

	// Handle Link play
	if Kplus.CheckIdLivetvIsKplus(DetailLiveTV.Id) { // DT-14322

		Usi, err := Fingering.GetUniqueStreamingID(userId)
		if Usi == "" || err != nil {
			DetailLiveTV.Hls_link_play = ""
			DetailLiveTV.Dash_link_play = ""
			DetailLiveTV.Link_play = ""
			DetailLiveTV.Programme.Hls_link_play = ""
			DetailLiveTV.Programme.Dash_link_play = ""
			DetailLiveTV.Permission = PERMISSION_REQUEST_LIMIT
		}
		DetailLiveTV.Usi = Usi
	}

	DetailLiveTV.Link_play = DetailLiveTV.Hls_link_play

	return
}
