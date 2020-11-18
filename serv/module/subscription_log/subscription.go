package subscription_log

import (
	. "cm-v5/serv/module"
	// "fmt"
	"time"
)

func (this *SubcriptionLogObjectStruct) Save() (err error) {
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
	if this.Id == 0 {
		// Insert
		this.Created_date = CurentTime
		this.Updated_date = CurentTime

		StmtIns, err := DbMysql.Prepare(`
			INSERT billing_subscription_logs
			SET  user_id=?,  package_id=?, subscribed_date=?, expired_date=?
			, is_active=?, created_date=?, updated_date=?`)
		if err != nil {
			return err
		}
		_, err = StmtIns.Exec(this.User_id, this.Package_id, this.Subscribed_date, this.Expired_date,
			this.Is_active, this.Created_date, this.Updated_date)
		if err != nil {
			return err
		}

	}

	return err
}
