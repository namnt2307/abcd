package livetv_v3

import (
	"fmt"
	"log"
	"strings"
	"time"

	. "cm-v5/serv/module"

	. "cm-v5/schema"

	Subscription "cm-v5/serv/module/subscription"

	"cm-v5/serv/module/packages"

	Fingering "cm-v5/serv/module/fingering"

	Kplus "cm-v5/serv/module/kplus"
)

const (
	KEY_RULE_LIVETV = "rule_livetv"
)

type RuleStruct struct {
	User_type int
	Ips       []string
}

func RouterGetPermission(DetailLiveTV *DetailLiveTVObjectOutputStruct, userId, ipUser, tokenUser string, statusUserIsPremium int) {
	valueCache, err := mRedis.GetString(KEY_RULE_LIVETV)
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

			if statusUserIsPremium != 99 {
				statusUserIsPremium = userType
			}

			for _, val := range Rule {
				if val.User_type == userType && CheckIpRule(ipUser, val.Ips) {
					GetPermissionListLiveTVDetailV2(DetailLiveTV, userId, ipUser, tokenUser, statusUserIsPremium)
					return
				}
			}
		}
	}
	GetPermissionListLiveTVDetail(DetailLiveTV, userId, ipUser, tokenUser, statusUserIsPremium)
}

func GetPermissionListLiveTVDetailV2(DetailLiveTV *DetailLiveTVObjectOutputStruct, userId, ipUser, tokenUser string, statusUserIsPremium int) {
	// Set link default
	DetailLiveTV.Hls_link_play = ""
	DetailLiveTV.Dash_link_play = ""
	DetailLiveTV.Programme.Hls_link_play = ""
	DetailLiveTV.Programme.Dash_link_play = ""

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
		return
	}

	//get current epg and check permission
	dataCurrentEpg, _ := GetEPGCurrent(DetailLiveTV.Id)
	if dataCurrentEpg.Id != "" && dataCurrentEpg.Ott_disabled == 1 {
		DetailLiveTV.Permission = PERMISSION_NOT_ALLOW
	}

	// VBE-17
	livetvId := DetailLiveTV.Id
	for _, ip := range LIVETV_LIST_IP_LOCK {
		if ip == ipUser {
			livetvId = "dfacafdc-1a67-11eb-adc0-c79cce82064c"
			break
		}
	}
	// VBE-17

	// Handle link play POTF
	t := time.Now()
	infoLiveTv, err := GetLiveTVInfo(livetvId, ipUser, DetailLiveTV.Programme.Time_start, DetailLiveTV.Programme.Duration)
	log.Println(fmt.Sprintf("Total time request api POFT: %f ms", time.Since(t).Seconds()))
	if err != nil {
		log.Println(err)
		return
	}

	DetailLiveTV.Hls_link_play = infoLiveTv.Result.Live_links.Hls
	DetailLiveTV.Dash_link_play = infoLiveTv.Result.Live_links.Dash
	DetailLiveTV.Programme.Hls_link_play = infoLiveTv.Result.Epg_links.Hls
	DetailLiveTV.Programme.Dash_link_play = infoLiveTv.Result.Epg_links.Dash

	//check drm set is vieon
	//Khanh DT-11862
	if DetailLiveTV.Drm_service_name == "vieon" {
		DetailLiveTV.Hls_link_play = VieONBuildDrmTokenUrl(DetailLiveTV.Hls_link_play, DetailLiveTV.Id, "livetv", ipUser, tokenUser)
		DetailLiveTV.Dash_link_play = VieONBuildDrmTokenUrl(DetailLiveTV.Dash_link_play, DetailLiveTV.Id, "livetv", ipUser, tokenUser)

		DetailLiveTV.Programme.Hls_link_play = VieONBuildDrmTokenUrl(DetailLiveTV.Programme.Hls_link_play, DetailLiveTV.Id, "livetv", ipUser, tokenUser)
		DetailLiveTV.Programme.Dash_link_play = VieONBuildDrmTokenUrl(DetailLiveTV.Programme.Dash_link_play, DetailLiveTV.Id, "livetv", ipUser, tokenUser)

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
