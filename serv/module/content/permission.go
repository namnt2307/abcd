package content

import (
	// "fmt"
	"strings"
	// "time"

	// . "cm-v5/serv/module"
	. "cm-v5/schema"
)

func GetPermissionDefaultVOD(ContentObjOutput ContentObjOutputStruct) ContentObjOutputStruct {
	// Set default
	ContentObjOutput.Default_episode.Permission = permissionDefaultVod
	ContentObjOutput.Permission = permissionDefaultVod
	ContentObjOutput.PackageGroup = make([]PackageGroupObjectStruct, 0)

	return ContentObjOutput
}

// User co dang nhap 		=> permission = 206
// User khong dang nhap 	=> permission = 0
func GetPermissionVOD(ContentObjOutput ContentObjOutputStruct, userId string) ContentObjOutputStruct {

	// KT Vod co thuoc package nao khong - lay danh sach package id
	if len(ContentObjOutput.PackageGroup) <= 0 {

		if strings.HasPrefix(userId, "anonymous_") == false && userId != "" {
			// Co login duoc phep xem
			ContentObjOutput.Permission = PERMISSION_VALID
			ContentObjOutput.Default_episode.Permission = PERMISSION_VALID
		}
	} else {
		//Neu mua thi cho xem phim
		// if CheckPermissionUser(ContentObjOutput, userId) {
		ContentObjOutput.Default_episode.Permission = PERMISSION_VALID
		ContentObjOutput.Permission = PERMISSION_VALID
		// }
	}
	return ContentObjOutput
}

// func CheckPermissionUser(ContentObjOutput ContentObjOutputStruct, userId string) bool {
// 	return true
// 	//Connect mysql
// 	db_mysql, _ := ConnectMySQL()
// 	defer db_mysql.Close()

// 	// Get time to day
// 	t := time.Now()
// 	var (
// 		currentDate            string = t.Format("2006-01-02 15:04:05")
// 		existsCheckUserPackage int
// 		justStringPackage      []string
// 	)
// 	for _, val := range ContentObjOutput.PackageGroup {
// 		justStringPackage = append(justStringPackage, fmt.Sprint(val.Id))
// 	}

// 	//Kiem tra user co mua goi xem content
// 	sqlRaw := fmt.Sprintf(`
// 		SELECT bp.id
// 		FROM billing_transaction as bt
// 		LEFT JOIN billing_packages as bp ON bt.package_id = bp.id
// 		WHERE bt.user_id = "%s" AND bt.status = 1 AND bt.expiry_date >= "%s" AND bp.billing_package_group_id IN ('%s')`, userId, currentDate, strings.Join(justStringPackage, "','"))
// 	err := db_mysql.QueryRow(sqlRaw)
// 	if existsCheckUserPackage <= 0 || err != nil {
// 		return false
// 	}
// 	return true
// }
