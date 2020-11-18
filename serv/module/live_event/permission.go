package live_event

import (
	"strings"
	. "cm-v5/serv/module"
	. "cm-v5/schema"
	"cm-v5/serv/module/packages"
)

// Không bat login
var permissionDefault = PERMISSION_VALID

func init() {
	permissionDefault, _ = CommonConfig.GetInt("GLOBAL_CONFIG", "permission_default")
}

func GetPackageLiveEvent(LiveEvent LiveEventOutputObjectStruct) LiveEventOutputObjectStruct {
	var LiveEvenID = LiveEvent.Id

	//Lay thong tin goi + permission.
	var Pack packages.PackagesObjectStruct
	LiveEvent.PackageGroup, LiveEvent.Permission = Pack.GetPackage(LiveEvenID)

	return LiveEvent
}

func GetPermissionLiveEvent(LiveEventOutput *LiveEventOutputObjectStruct, userId, ipUser string, statusUserIsPremium int) {
	// KT Vod co thuoc package nao khong - lay danh sach package id
	if len(LiveEventOutput.PackageGroup) <= 0 {
		if strings.HasPrefix(userId, "anonymous_") == false && userId != "" {
			// Co login duoc phep xem
			LiveEventOutput.Permission = PERMISSION_VALID
		}
	} else {
		//Neu mua thi cho xem phim
		if ok, _ := packages.CheckPermissionUser(LiveEventOutput.Id, userId, statusUserIsPremium); ok {
			LiveEventOutput.Permission = PERMISSION_VALID
		}
	}

	LiveEventOutput.Link_play.Dash_link_play = BuildTokenUrl(LiveEventOutput.Link_play.Dash_link_play, "", ipUser)
	LiveEventOutput.Link_play.Hls_link_play = BuildTokenUrl(LiveEventOutput.Link_play.Hls_link_play, "", ipUser)

	if LiveEventOutput.Permission != PERMISSION_VALID {
		LiveEventOutput.Link_play.Dash_link_play = ""
		LiveEventOutput.Link_play.Hls_link_play = ""
	}

	LiveEventOutput.Link_play.Dash_link_play = HandleLinkPlayVieON(LiveEventOutput.Link_play.Dash_link_play)
	LiveEventOutput.Link_play.Hls_link_play = HandleLinkPlayVieON(LiveEventOutput.Link_play.Hls_link_play)

	return
}

func GetPermissionListLiveEvent(LiveEventOutput []LiveEventOutputObjectStruct, userId, ipUser string, statusUserIsPremium int) []LiveEventOutputObjectStruct {

	var LiveEventOutputPer []LiveEventOutputObjectStruct
	for _, val := range LiveEventOutput {

		// KT Vod co thuoc package nao khong - lay danh sach package id
		if len(val.PackageGroup) <= 0 {
			if strings.HasPrefix(userId, "anonymous_") == false && userId != "" {
				// Co login duoc phep xem
				val.Permission = PERMISSION_VALID
			}
		} else {
			//Neu mua thi cho xem phim
			if ok, _ := packages.CheckPermissionUser(val.Id, userId, statusUserIsPremium); ok {
				val.Permission = PERMISSION_VALID
			}
		}

		//Build token link play
		val.Link_play.Dash_link_play = BuildTokenUrl(val.Link_play.Dash_link_play, "", ipUser)
		val.Link_play.Hls_link_play = BuildTokenUrl(val.Link_play.Hls_link_play, "", ipUser)

		if val.Permission != PERMISSION_VALID {
			val.Link_play.Dash_link_play = ""
			val.Link_play.Hls_link_play = ""
		}

		val.Link_play.Dash_link_play = HandleLinkPlayVieON(val.Link_play.Dash_link_play)
		val.Link_play.Hls_link_play = HandleLinkPlayVieON(val.Link_play.Hls_link_play)
		LiveEventOutputPer = append(LiveEventOutputPer, val)
	}

	return LiveEventOutputPer
}
