package packages

import (
	"fmt"

	. "cm-v5/serv/module"
	. "cm-v5/schema"
)

func (this *PackagesObjectStruct) GetListPackageByDenyContent(idContent string, cacheActive bool) (PackageGroupDenyContents []PackageGroupDenyContentStruct, err error) {
	keyCache := "KV_PACKAGEs_BY_CONTENT_" + idContent

	if cacheActive {
		// Read cache
		valueCache, err := mRedisKV.GetString(keyCache)
		if err == nil && valueCache != "" {
			err = json.Unmarshal([]byte(valueCache), &PackageGroupDenyContents)
			if err == nil {
				return PackageGroupDenyContents, nil
			}
		}
	}

	//Connect mysql
	DbMysql, err := ConnectMySQL()
	if err != nil {
		return PackageGroupDenyContents, err
	}
	defer DbMysql.Close()

	// Get list group package id + package id
	SqlRow := fmt.Sprintf(`
		SELECT bpdc.id , bpdc.billing_package_group_id , bpdc.content_id , 
				bpdc.name ,  bp.id as pk_id
		FROM billing_packages_deny_contents as bpdc
		JOIN billing_packages as bp ON bp.billing_package_group_id = bpdc.billing_package_group_id
		WHERE bpdc.content_id = "%s" AND bp.is_active IN (1, 2)`, idContent)

	dataPackages, err := DbMysql.Query(SqlRow)
	if err != nil {
		return PackageGroupDenyContents, err
	}

	//fomart result in db
	for dataPackages.Next() {
		var PackageGroupDenyContent PackageGroupDenyContentStruct
		err = dataPackages.Scan(&PackageGroupDenyContent.Id, &PackageGroupDenyContent.Billing_package_group_id,
			&PackageGroupDenyContent.Content_id, &PackageGroupDenyContent.Name,
			&PackageGroupDenyContent.Pk_id)
		if err != nil {
			continue
		}
		PackageGroupDenyContents = append(PackageGroupDenyContents, PackageGroupDenyContent)
	}

	// Write cache
	dataByte, _ := json.Marshal(PackageGroupDenyContents)
	mRedisKV.SetString(keyCache, string(dataByte), TTL_REDIS_LV1)

	return PackageGroupDenyContents, err
}
