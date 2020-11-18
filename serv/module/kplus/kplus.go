package kplus

import (
	"fmt"
	"strings"
	"time"

	. "cm-v5/schema"
	. "cm-v5/serv/module"

	Config "cm-v5/serv/module/config"
	Subcription "cm-v5/serv/module/subscription"
)

var (
	listIdPackagesOfKplus []string
	listIDLivetvKplus     []string
	groupPackageKplus     int
	livetvGroupID         string
)

// Type banner
// 0 => nonVIP
// 1 => VIP
// 2 => VIP + KPLUS

type TypeBanner struct {
	Type                    int    `json:"type"`
	Package_group_id        int    `json:"package_group_id"`
	Livetv_group_id         string `json:"livetv_group_id"`
	Livetv_id               string `json:"livetv_id"`
	Hide_button_buy_package bool   `json:"hide_button_buy_package"`
}

func init() {
	listIDString, _ := CommonConfig.GetString("KPLUS", "list_packages")
	if listIDString != "" {
		listIdPackagesOfKplus = strings.Split(listIDString, ",")
	}

	groupPackageStr, _ := CommonConfig.GetString("KPLUS", "group_package")
	groupPackageKplus, _ = StringToInt(groupPackageStr)
	livetvGroupID, _ = CommonConfig.GetString("KPLUS", "livetv_group_kplus")

	//get id livetv K+
	LIVETV_KPLUS_1_HD_ID, _ := CommonConfig.GetString("KPLUS", "livetv_kplus_1_hd_id")
	LIVETV_KPLUS_NS_HD_ID, _ := CommonConfig.GetString("KPLUS", "livetv_kplus_ns_hd_id")
	LIVETV_KPLUS_PC_HD_ID, _ := CommonConfig.GetString("KPLUS", "livetv_kplus_pc_hd_id")
	LIVETV_KPLUS_PM_HD_ID, _ := CommonConfig.GetString("KPLUS", "livetv_kplus_pm_hd_id")
	listIDLivetvKplus = []string{LIVETV_KPLUS_1_HD_ID, LIVETV_KPLUS_NS_HD_ID, LIVETV_KPLUS_PC_HD_ID, LIVETV_KPLUS_PM_HD_ID}

}

func GetTypeBannerKplus(user_id, platform string) TypeBanner {
	var bannerTypeOutput TypeBanner

	//package group of kplus
	bannerTypeOutput.Package_group_id = groupPackageKplus
	bannerTypeOutput.Livetv_group_id = livetvGroupID
	bannerTypeOutput.Livetv_id = listIDLivetvKplus[0]

	var Sub Subcription.SubcriptionObjectStruct
	listSubcription, err := Sub.GetListByUserId(user_id)

	if err != nil {
		return bannerTypeOutput
	}
	if len(listSubcription) > 0 {
		bannerTypeOutput.Type = 1

		for _, val := range listSubcription {
			if CheckPackageIsVipKplus(val.Current_package_id) {
				bannerTypeOutput.Type = 2

				//get config
				configExpireDays, _ := Config.GetConfigByKey(CONFIG_KEY_PACKAGE_EXPIRED_DAYS, platform, true)
				package_expired_days := 0
				if configExpireDays.Data.Value != "" {
					package_expired_days, _ = StringToInt(configExpireDays.Data.Value)
				}
				expiryDate := date(Sub.Expiry_date)
				toDate := time.Now()
				if GetNumDaysOfTwoDate(toDate, expiryDate) > package_expired_days {
					bannerTypeOutput.Hide_button_buy_package = true
				}

				break
			}
		}
	}

	return bannerTypeOutput
}

func CheckPackageIsVipKplus(package_id int) bool {
	if ok, _ := In_array(fmt.Sprint(package_id), listIdPackagesOfKplus); ok {
		return true
	}
	return false
}

func CheckIdLivetvIsKplus(livetv_id string) bool {
	ok, _ := In_array(livetv_id, listIDLivetvKplus)
	return ok
}

func GetNumDaysOfTwoDate(start, end time.Time) int {

	diff := end.Sub(start)
	days := int(diff.Hours() / 24)

	return days
}

func date(s string) time.Time {
	layout1 := "2006-01-02 15:04:05"
	layout2 := "2006-01-02"
	d, err := time.Parse(layout1, s)
	if err != nil {
		d, _ = time.Parse(layout2, s)
	}
	return d
}
