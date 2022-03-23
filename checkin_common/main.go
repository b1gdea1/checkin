package main

import (
	"flag"
	"fmt"
	"github.com/b1gdea1/checkin_common/common_service"
	"github.com/b1gdea1/checkin_common/mysqlDao"
	redisdriver "github.com/b1gdea1/checkin_common/redis"
	"github.com/jinzhu/configor"
	"google.golang.org/grpc"
	"log"
	"net"
	"os"
	"os/signal"
	"syscall"
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
	keyPath := flag.String("p", "C:\\source\\checkin\\config\\CheckinConfig.yml", "配置文件路径")
	flag.Parse()
	log.Printf("配置文件路径：%v", *keyPath)
	
	err := configor.Load(&Config, *keyPath)
	if err != nil {
		return
	}
	fmt.Printf("config: %#v", Config)
	mysqlDao.InitDao(Config.Mysql.Database, Config.Mysql.User, Config.Mysql.Password, Config.Mysql.Port, Config.Mysql.Address)
	redisdriver.PoolInitRedis(Config.Redis.Address, Config.Redis.Auth)
	// s := grpc.NewServer(grpc.Creds(keys.GetServerCreds(Config.KeyPath)))
	s := grpc.NewServer()
	common_service.RegisterStudentServiceServer(s, &common_service.Server{})
	listener, err := net.Listen("tcp", ":8081")
	if err != nil {
		log.Fatal(err)
		return
	}
	errChan := make(chan error)
	go func() {
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

// var pool *redigo.Pool
//
// func init() {
// 	pool = &redigo.Pool{
// 		MaxIdle:     5, // 空闲数
// 		IdleTimeout: 240 * time.Second,
// 		MaxActive:   20, // 最大数
// 		Dial: func() (redigo.Conn, error) {
// 			c, err := redigo.Dial("tcp", "39.99.39.103:6379")
// 			if err != nil {
// 				return nil, err
// 			}
// 			if _, err := c.Do("AUTH", "TGYtgy1234"); err != nil {
// 				c.Close()
// 				return nil, err
// 			}
// 			return c, err
// 		},
// 		TestOnBorrow: func(c redigo.Conn, t time.Time) error {
// 			_, err := c.Do("PING")
// 			return err
// 		},
// 	}
// }

// func main() {
// 	// 从 pool 中取出一个连接
// 	conn := pool.Get()
// 	// 向Redis写入一个数据
// 	_, err := conn.Do("Set", "name", "Rose")
// 	if err != nil {
// 		fmt.Println("conn.Do err is", err)
// 		return
// 	}
// 	// 取出
// 	_, err = redigo.String(conn.Do("Get", "name"))
// 	if err != nil {
// 		fmt.Println("conn.DO Get err is", err)
// 		return
// 	}
// 	fmt.Println(pool.ActiveCount())
// 	fmt.Println(pool.IdleCount())
// 	fmt.Println()
// 	conn.Close()
// 	fmt.Println(pool.ActiveCount())
// 	fmt.Println(pool.IdleCount())
// 	fmt.Println()
// 	conn1 := pool.Get()
// 	// 向Redis写入一个数据
// 	_, err = conn1.Do("Set", "name", "Rose")
// 	if err != nil {
// 		fmt.Println("conn.Do err is", err)
// 		return
// 	}
// 	fmt.Println(pool.ActiveCount())
// 	fmt.Println(pool.IdleCount())
// 	fmt.Println()
// 	conn2 := pool.Get()
// 	// 向Redis写入一个数据
// 	_, err = conn2.Do("Set", "name", "Rose")
// 	if err != nil {
// 		fmt.Println("conn.Do err is", err)
// 		return
// 	}
// 	fmt.Println(pool.ActiveCount())
// 	fmt.Println(pool.IdleCount())
// 	fmt.Println()
// }
