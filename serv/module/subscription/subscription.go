package subscription

import (
	"time"

	. "cm-v5/serv/module"
	// . "cm-v5/schema"
)

func (this *SubcriptionObjectStruct) GetListByUserId(idUser string) (SubcriptionObjects []SubcriptionObjectStruct, err error) {

	keyCache := "USC_LIST_SUBCRIPTION_TEMP_" + idUser

	// Read cache
	valueCache, err := mRedisUSC.GetString(keyCache)
	if err == nil && valueCache != "" {
		err = json.Unmarshal([]byte(valueCache), &SubcriptionObjects)
		if err == nil {
			return SubcriptionObjects, nil
		}
	}

	//Read DB
	SubcriptionObjects, err = this.GetListByUserIdWithDB(idUser)
	if err != nil {
		return SubcriptionObjects, err
	}

	return SubcriptionObjects, nil

}

func (this *SubcriptionObjectStruct) GetListByUserIdWithDB(idUser string) (SubcriptionObjects []SubcriptionObjectStruct, err error) {
	//Connect mysql
	DbMysql, err := ConnectMySQL()
	if err != nil {
		return SubcriptionObjects, err
	}
	defer DbMysql.Close()

	dataSubcriptionObjects, err := DbMysql.Query(`
	SELECT bs.id, bs.is_trialed, bs.created_at, bs.status, bs.current_package_id, bs.expiry_date, bpg.type, COALESCE(bs.recurring, 0)
	FROM billing_subscription AS bs
	LEFT JOIN billing_packages AS bp ON bp.id = bs.current_package_id
	LEFT JOIN billing_package_group AS bpg ON bpg.id = bp.billing_package_group_id
	WHERE bs.user_id = ? and bs.status = 1`, idUser)
	if err != nil {
		return SubcriptionObjects, err
	}

	//fomart result in db
	for dataSubcriptionObjects.Next() {
		var SubcriptionObject SubcriptionObjectStruct
		err = dataSubcriptionObjects.Scan(&SubcriptionObject.Id, &SubcriptionObject.Is_trialed, &SubcriptionObject.Created_at, &SubcriptionObject.Status, &SubcriptionObject.Current_package_id, &SubcriptionObject.Expiry_date, &SubcriptionObject.Type, &SubcriptionObject.Recurring)
		if err != nil {
			continue
		}

		SubcriptionObjects = append(SubcriptionObjects, SubcriptionObject)
	}

	// Write cache subcription
	keyCache := "USC_LIST_SUBCRIPTION_TEMP_" + idUser
	dataByte, _ := json.Marshal(SubcriptionObjects)
	mRedisUSC.SetString(keyCache, string(dataByte), -1)

	return SubcriptionObjects, err
}

func (this *SubcriptionObjectStruct) Save() (err error) {
	// Get time to day
	t := time.Now()
	CurentTime := t.Format("2006-01-02 15:04:05")

	// Connect MySQL
	DbMysql, err := ConnectMySQL()
	if err != nil {
		return err
	}
	defer DbMysql.Close()

	// Check Insert / Update
	if this.Id == "" {
		// Insert
		this.Id = string(UUIDV4())
		this.Created_at = CurentTime
		this.Updated_at = CurentTime

		StmtIns, err := DbMysql.Prepare(`
			INSERT billing_subscription
			SET id=?, start_date=?, expiry_date=?, is_trialed=?, method=?
			, suspended_date=?, created_at=?, updated_at=?, status=?, current_package_id=?
			, last_success_transaction_id=?, last_transaction_id=?, next_package_id=?
			, paypal_agreement_id=?, user_id=?`)
		if err != nil {
			return err
		}
		_, err = StmtIns.Exec(this.Id, this.Start_date, this.Expiry_date, this.Is_trialed, this.Method,
			this.Suspended_date, this.Created_at, this.Updated_at, this.Status, this.Current_package_id,
			this.Last_success_transaction_id, this.Last_transaction_id, this.Next_package_id,
			this.Paypal_agreement_id, this.User_id)
		if err != nil {
			return err
		}

	} else {
		// update
		this.Updated_at = CurentTime
	}

	//Write cache code
	this.GetListByUserIdWithDB(this.User_id)

	return err
}
