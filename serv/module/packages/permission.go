package packages

import (
	. "cm-v5/serv/module"
	. "cm-v5/schema"
	Kplus "cm-v5/serv/module/kplus"
	Subscription "cm-v5/serv/module/subscription"
)

// Không bat login
var permissionDefault = PERMISSION_VALID
var listIDLivetvKplus []string

func init() {
	permissionDefault, _ = CommonConfig.GetInt("GLOBAL_CONFIG", "permission_default")
}

func (this *PackagesObjectStruct) GetPackage(content_id string) ([]PackageGroupObjectStruct, int) {
	// Set default
	var Permission = permissionDefault
	var PackageGroup = make([]PackageGroupObjectStruct, 0)

	var keyCache = LOCAL_CONTENT_PERMISSION + "_" + content_id
	//Get package
	valC, err := LocalCache.GetValue(keyCache)
	if err == nil {
		dataByte, _ := json.Marshal(valC)
		json.Unmarshal([]byte(dataByte), &PackageGroup)
	} else {
		//Connect mysql
		db_mysql, _ := ConnectMySQL()
		defer db_mysql.Close()

		dataPackage, err := db_mysql.Query(`
		SELECT bpg.id, bpg.name, bp.price, bp.duration
		FROM billing_packages_deny_contents as bpdc
		LEFT JOIN billing_package_group as bpg ON  bpdc.billing_package_group_id = bpg.id
		LEFT JOIN billing_packages as bp ON  bp.billing_package_group_id = bpg.id
		WHERE bpdc.content_id = ? AND bpg.is_active = 1
		ORDER BY bp.price ASC LIMIT 1`, content_id)

		if err == nil {
			for dataPackage.Next() {
				var PackageObject PackageGroupObjectStruct
				err := dataPackage.Scan(&PackageObject.Id, &PackageObject.Name, &PackageObject.Price, &PackageObject.Period_value)

				if err != nil {
					continue
				}

				PackageObject.Period = "hours"
				PackageGroup = append(PackageGroup, PackageObject)
			}

			// Write local data x2
			LocalCache.SetValue(keyCache, PackageGroup, TTL_LOCALCACHE*2)
		}
	}

	//Check permision
	if len(PackageGroup) > 0 {
		// Content can mua goi gan Permission 208
		Permission = PERMISSION_REQUIRE_PACKAGE
	}

	return PackageGroup, Permission
}

//CheckPermissionUser return valid and len of package SVOD
func CheckPermissionUser(contentId, userId string, statusUserIsPremium int) (bool, int) {

	var Pack PackagesObjectStruct
	PackageGroupDenyContents, _ := Pack.GetListPackageByDenyContent(contentId, true)
	if len(PackageGroupDenyContents) <= 0 {
		return true, len(PackageGroupDenyContents)
	}

	//số lượng gói chứa content này
	numPackages := len(PackageGroupDenyContents)

	//user super premium pass permission
	if statusUserIsPremium == USER_SUPER_PREMIUM {
		return true, numPackages
	}

	// Kiem tra user đã mua package nào thuộc list package trên chưa
	// Neu co permission valid
	var Sub Subscription.SubcriptionObjectStruct
	SubcriptionsOfUser, _ := Sub.GetListByUserId(userId)
	if len(SubcriptionsOfUser) <= 0 {
		// User chưa mua goi
		return false, numPackages
	}

	isLivetvKplus := Kplus.CheckIdLivetvIsKplus(contentId)

	for _, vSubObj := range SubcriptionsOfUser {
		//Gói LG200K và Gói VIP + K+ thì được xem content VIP.
		//Content đang check là livetv của K+ thì bắt buộc mua gói có add các kênh này
		if !isLivetvKplus && IsPkgIncludeContentVip(vSubObj.Current_package_id) {
			return true, numPackages
		}
		for _, vPGDC := range PackageGroupDenyContents {
			if vSubObj.Current_package_id == vPGDC.Pk_id {
				return true, numPackages
			}
		}
	}

	return false, numPackages
}

//IsPkgIncludeContentVip check xem gói được xem content VIP hay không
func IsPkgIncludeContentVip(package_id int) bool {

	//gói LG200K được phép xem content VIP (package_id = 30)
	//gói VIP + K+ được phép xem content VIP
	if package_id == 30 || Kplus.CheckPackageIsVipKplus(package_id) {
		return true
	}

	return false
}
