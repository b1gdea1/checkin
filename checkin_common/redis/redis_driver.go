package redis

import (
	"fmt"
	redigo "github.com/gomodule/redigo/redis"
	"log"
	"time"
)

var p *redigo.Pool = nil

// PoolInitRedis redis pool
func PoolInitRedis(server string, password string) {
	p = &redigo.Pool{
		MaxIdle:     5, // 空闲数
		IdleTimeout: 240 * time.Second,
		MaxActive:   20, // 最大数
		Dial: func() (redigo.Conn, error) {
			c, err := redigo.Dial("tcp", server)
			if err != nil {
				return nil, err
			}
			if password != "" {
				if _, err := c.Do("AUTH", password); err != nil {
					c.Close()
					return nil, err
				}
			}
			return c, err
		},
		TestOnBorrow: func(c redigo.Conn, t time.Time) error {
			_, err := c.Do("PING")
			return err
		},
	}
}

func UploadUserToken(userID string, token string, ttlSecond int) bool {
	var ret bool
	conn := p.Get()
	defer conn.Close()                    // close不是直接关闭连接，而是将连接放回池中
	result, err := conn.Do("select", "0") // 进入数据库0
	if err != nil {
		return false
	}
	log.Printf("select DB 1:%v", result)
	defer func(conn redigo.Conn) {
		err := conn.Close()
		if err != nil {
			fmt.Println(err)
		}
	}(conn)
	result, err = conn.Do("SET", userID, token)
	if err != nil {
		return false
	}
	ret = result.(string) == "OK"
	result, err = conn.Do("EXPIRE", userID, ttlSecond)
	if err != nil {
		return false
	}
	ret = result.(int64) == 1
	return ret
}
func VerifyToken(userID string, token string) bool {
	conn := p.Get()
	defer conn.Close()
	result, err := conn.Do("select", "0") // 进入一号数据库
	if err != nil {
		return false
	}
	log.Printf("select DB 1:%v", result)
	defer func(conn redigo.Conn) {
		err := conn.Close()
		if err != nil {
			fmt.Println(err)
		}
	}(conn)
	result, err = conn.Do("GET", userID)
	if err != nil {
		return false
	}
	return (string)(result.([]uint8)) == token
}
