package main

import (
	"flag"
	"fmt"
	"github.com/b1gdea1/checkin_admin/admin_service"
	"github.com/b1gdea1/checkin_admin/mysqlDao"
	"github.com/b1gdea1/checkin_admin/redis"
	_ "github.com/go-sql-driver/mysql"
	"github.com/jinzhu/configor"
	"google.golang.org/grpc"
	"log"
	"net"
	"os"
	"os/signal"
	"syscall"
	"time"
)

var Config = struct {
	Mysql struct {
		Database string
		Address  string
		User     string `default:"root"`
		Password string `required:"true" env:"DBPassword"`
		Port     int    `default:"3306"`
	}
	Redis struct {
		Address string
		Auth    string `required:"true"`
	}
}{}

func main() {
	keyPath := flag.String("p", "C:\\source\\checkin\\config\\CheckinAdminServiceConfig.yml", "配置文件路径")
	flag.Parse()
	log.Printf("配置文件路径：%v", *keyPath)
	
	err := configor.Load(&Config, *keyPath)
	if err != nil {
		return
	}
	fmt.Printf("config: %#v", Config)
	mysqlDao.InitDao(Config.Mysql.Database, Config.Mysql.User, Config.Mysql.Password, Config.Mysql.Port, Config.Mysql.Address)
	redis.PoolInitRedis(Config.Redis.Address, Config.Redis.Auth)
	s := grpc.NewServer()
	admin_service.RegisterAdminServiceServer(s, &admin_service.CheckinServiceAsAdministrator{})
	
	go func() { // 每日定时任务
		for {
			questions, err := redis.GetTodayQuestion(time.Now().AddDate(0, 0, 1), time.Now())
			if err != nil {
				return
			}
			ret := mysqlDao.UpdateQuestions(questions)
			if !ret {
				log.Fatal("update question failed")
			}
			now := time.Now()
			// 计算下一个零点
			next := now.Add(time.Hour * 24)
			next = time.Date(next.Year(), next.Month(), next.Day(), 0, 0, 0, 0, next.Location())
			t := time.NewTimer(next.Sub(now))
			<-t.C
		}
	}()
	
	errChan := make(chan error)
	go func() {
		listener, err := net.Listen("tcp", ":8082")
		if err != nil {
			log.Fatal(err)
			return
		}
		
		err = s.Serve(listener)
		if err != nil {
			return
		}
		if err != nil {
			log.Println(err)
			errChan <- err
		}
	}()
	go func() {
		c := make(chan os.Signal)
		signal.Notify(c, syscall.SIGINT, syscall.SIGTERM) // 中断 和 中止
		errChan <- fmt.Errorf("%s", <-c)
	}()
	getErr := <-errChan
	fmt.Println(getErr)
}
