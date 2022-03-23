package mysqlDao

import (
	"database/sql"
	"fmt"
	_ "github.com/go-sql-driver/mysql"
	"log"
)

var MysqlDb *sql.DB
var MysqlDbErr error

func InitDao(DATABASE string, UserName string, PassWord string, PORT int, HOST string) {
	dbDSN := fmt.Sprintf("%s:%s@tcp(%s:%d)/%s", UserName, PassWord, HOST, PORT, DATABASE)
	MysqlDb, MysqlDbErr = sql.Open("mysql", dbDSN)
	if MysqlDbErr != nil {
		log.Println("创建数据库连接失败")
		panic(MysqlDbErr)
		return
	}
	fmt.Println("创建数据库连接成功")
}
