package packages

import (
	"errors"
	"fmt"
	"time"

	. "cm-v5/serv/module"
	// . "cm-v5/schema"
)

func (this *PackagesObjectStruct) FetchDetailById(idPackage int) (err error) {
	keyCache := "KV_PACKAGE_DETAIL_" + fmt.Sprint(idPackage)

	// Read cache
	valueCache, err := mRedisKV.GetString(keyCache)
	if err == nil && valueCache != "" {
		err = json.Unmarshal([]byte(valueCache), &this)
		if err == nil {
			return nil
		}
	}

	var curDateStr = time.Now().Format("2006-01-02 15:04:05.000000")

	//Connect mysql
	DbMysql, err := ConnectMySQL()
	if err != nil {
		return err
	}
	defer DbMysql.Close()

	// SqlRow := fmt.Sprintf(`
	// 	SELECT id , name , price , duration , expired_date
	// 	FROM billing_packages
	// 	WHERE is_active = 1 and id = %d and expired_date >= "%s"`, idPackage, curDateStr)
	dataPackagesObject, err := DbMysql.Query(`
	SELECT id , name , price , duration, duration_type , expired_date
	FROM billing_packages
	WHERE is_active = 1 and id = ? and expired_date >= ?`, idPackage, curDateStr)
	if err != nil {
		return err
	}

	//fomart result in db
	for dataPackagesObject.Next() {
		err = dataPackagesObject.Scan(&this.Id, &this.Name, &this.Price, &this.Duration, &this.Duration_type, &this.Expired_date)
		if err != nil {
			return err
		}
	}

	if this.Id == 0 {
		return errors.New("Package not valid")
	}

	// Write cache
	dataByte, _ := json.Marshal(this)
	mRedisKV.SetString(keyCache, string(dataByte), TTL_KVCACHE)

	return err
}
