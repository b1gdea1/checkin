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

func TestAdminExist(t *testing.T) {
	err := configor.Load(&Config, "C:\\source\\config\\CheckinAdminServiceConfig.yml")
	if err != nil {
		return
	}
	fmt.Printf("config: %#v", Config)
	InitDao(Config.Mysql.Database, Config.Mysql.User, Config.Mysql.Password, Config.Mysql.Port, Config.Mysql.Address)
	defer MysqlDb.Close()
	fmt.Println(AdminExist("201930507008"))
}
func TestGetAdministratorID(t *testing.T) {
	err := configor.Load(&Config, "C:\\source\\config\\CheckinAdminServiceConfig.yml")
	if err != nil {
		return
	}
	fmt.Printf("config: %#v", Config)
	InitDao(Config.Mysql.Database, Config.Mysql.User, Config.Mysql.Password, Config.Mysql.Port, Config.Mysql.Address)
	defer MysqlDb.Close()
	fmt.Println(GetAdministratorID("201930507008", "123456"))
}
func TestUpdateQuestions(t *testing.T) {
	err := configor.Load(&Config, "C:\\source\\config\\CheckinAdminServiceConfig.yml")
	if err != nil {
		return
	}
	fmt.Printf("config: %#v", Config)
	InitDao(Config.Mysql.Database, Config.Mysql.User, Config.Mysql.Password, Config.Mysql.Port, Config.Mysql.Address)
	defer MysqlDb.Close()
	fmt.Println(UpdateQuestions([]*Question{
		&Question{
			QuestionID: 1,
			Type:       int32(0),
			Content:    "what is your name?",
			Answers:    "tell me",
			Depends:    "",
		},
		&Question{
			QuestionID: 2,
			Type:       int32(1),
			Content:    "what is your gender?",
			Answers:    "male,female,lgbt",
			Depends:    "",
		},
		&Question{
			QuestionID: 3,
			Type:       int32(0),
			Content:    "since you are male, how long is your hair?",
			Answers:    "tell me",
			Depends:    "2=male",
		},
	}))
}
func TestGetQuestions(t *testing.T) {
	err := configor.Load(&Config, "C:\\source\\config\\CheckinAdminServiceConfig.yml")
	if err != nil {
		return
	}
	fmt.Printf("config: %#v", Config)
	InitDao(Config.Mysql.Database, Config.Mysql.User, Config.Mysql.Password, Config.Mysql.Port, Config.Mysql.Address)
	defer MysqlDb.Close()
	questions, err := GetQuestions()
	if err != nil {
		return
	}
	for _, question := range questions {
		fmt.Println(question)
	}
	
}
