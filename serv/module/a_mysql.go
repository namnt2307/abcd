package module

import (
	"database/sql"
	"log"
	_ "github.com/go-sql-driver/mysql"
)

type MySQLCli struct {
	db *sql.DB
}

var (
	// instanceMySQLCli *MySQLCli = nil
	hostMySQL         string
	portMySQL         string
	usernameMySQL     string
	passwordMySQL     string
	dbnameMySQL       string
	dbnameMySQLStream string
)

func init() {
	hostMySQL, _ = CommonConfig.GetString("MYSQL", "host")
	portMySQL, _ = CommonConfig.GetString("MYSQL", "port")
	usernameMySQL, _ = CommonConfig.GetString("MYSQL", "username")
	passwordMySQL, _ = CommonConfig.GetString("MYSQL", "password")
	dbnameMySQL, _ = CommonConfig.GetString("MYSQL", "dbname")
	dbnameMySQLStream, _ = CommonConfig.GetString("MYSQL", "dbname_stream")

}

func ConnectMySQL() (db *sql.DB, err error) {
	var instanceMySQLCli = new(MySQLCli)
	instanceMySQLCli.db, err = sql.Open("mysql", usernameMySQL+":"+passwordMySQL+"@tcp("+hostMySQL+":"+portMySQL+")/"+dbnameMySQL)
	if err != nil {
		log.Fatalf("Connect MySQL: %s\n", err)
		Sentry_log(err)
		return instanceMySQLCli.db, err
	}

	return instanceMySQLCli.db, nil
}
