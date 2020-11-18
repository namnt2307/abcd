package transaction

import (
	"fmt"
	"time"

	. "cm-v5/serv/module"
	// . "cm-v5/schema"
)

func (this *TransactionObjectStruct) GetByUserAndPackage() (err error) {
	keyCache := "KV_TRANSACTION_DETAIL_" + this.User_id + "_" + this.Active_code_id

	// Read cache
	valueCache, err := mRedisKV.GetString(keyCache)
	if err == nil && valueCache != "" {
		err = json.Unmarshal([]byte(valueCache), &this)
		if err == nil {
			return nil
		}
	}

	//Connect mysql
	DbMysql, err := ConnectMySQL()
	if err != nil {
		return err
	}
	defer DbMysql.Close()

	SqlRow := fmt.Sprintf(`
		SELECT id , user_id , package_id , status , expiry_date
		FROM billing_transaction
		WHERE user_id = "%s" and active_code_id = "%s" and status = 1`, this.User_id, this.Active_code_id)

	dataTransactionObject, err := DbMysql.Query(SqlRow)
	if err != nil {
		return err
	}

	//fomart result in db
	for dataTransactionObject.Next() {
		err = dataTransactionObject.Scan(&this.Id, &this.User_id, &this.Package_id, &this.Status, &this.Expiry_date)
		if err != nil {
			continue
		}
	}

	if err == nil && this.Id != "" {
		// Write cache
		dataByte, _ := json.Marshal(this)
		mRedisKV.SetString(keyCache, string(dataByte), TTL_KVCACHE)
	}

	return err
}

func (this *TransactionObjectStruct) Save() (err error) {
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
			INSERT billing_transaction
			SET id=?, start_date=?, expiry_date=?, is_trialed=?, type=?, platform=?, currency=?, amount=?, setup_fee=?, payment_method=?,
			readable_id=?, created_at=?, updated_at=?, status=?, active_code_id=?, package_id=?, user_id=?`)
		if err != nil {
			return err
		}

		_, err = StmtIns.Exec(this.Id, this.Start_date, this.Expiry_date, this.Is_trialed, this.Type, this.Platform,
			this.Currency, this.Amount, this.Setup_fee, this.Payment_method,
			this.Readable_id, this.Created_at, this.Updated_at, this.Status,
			this.Active_code_id, this.Package_id, this.User_id)
		if err != nil {
			return err
		}
	} else {
		// update
		this.Updated_at = CurentTime
	}

	return nil
}
