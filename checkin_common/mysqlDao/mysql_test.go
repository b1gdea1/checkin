package mysqlDao

import (
	"fmt"
	"github.com/jinzhu/configor"
	"testing"
)

var Config = struct {
	APPName string `default:"app name"`
	Port    int
	Mysql   struct {
		Database string
		Address  string
		User     string `default:"root"`
		Password string `required:"true" env:"DBPassword"`
		Port     int    `default:"3306"`
	}
	KeyPath string
}{}

func TestGetUserID(t *testing.T) {
	err := configor.Load(&Config, "C:\\source\\config\\CheckinUserServiceConfig.yml")
	if err != nil {
		return
	}
	fmt.Printf("config: %#v", Config)
	InitDao(Config.Mysql.Database, Config.Mysql.User, Config.Mysql.Password, Config.Mysql.Port, Config.Mysql.Address)
	defer MysqlDb.Close()
	fmt.Println(GetStudentID("201930507008", "123456"))
}
func TestInsertUserInfo(t *testing.T) {
	err := configor.Load(&Config, "C:\\source\\config\\CheckinUserServiceConfig.yml")
	if err != nil {
		return
	}
	fmt.Printf("config: %#v", Config)
	InitDao(Config.Mysql.Database, Config.Mysql.User, Config.Mysql.Password, Config.Mysql.Port, Config.Mysql.Address)
	defer MysqlDb.Close()
	fmt.Println(InsertStudentFull("201930507009", "222322", "1223"))
}
func TestInsertStudentSimple(t *testing.T) {
	err := configor.Load(&Config, "C:\\source\\config\\CheckinUserServiceConfig.yml")
	if err != nil {
		return
	}
	fmt.Printf("config: %#v", Config)
	InitDao(Config.Mysql.Database, Config.Mysql.User, Config.Mysql.Password, Config.Mysql.Port, Config.Mysql.Address)
	defer MysqlDb.Close()
	fmt.Println(InsertStudentSimple("2019305070339", "222322"))
	
}
