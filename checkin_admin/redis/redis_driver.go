package redis

import (
	"encoding/json"
	"errors"
	"fmt"
	"github.com/b1gdea1/checkin_admin/base_service"
	"github.com/b1gdea1/checkin_admin/mysqlDao"
	redigo "github.com/gomodule/redigo/redis"
	"log"
	"strings"
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

// 用户token-->0
// 问题--> 1
// 通知-->2

func UploadUserToken(userID string, token string, ttlSecond int) bool {
	var ret bool
	conn := p.Get()
	defer conn.Close()
	result, err := conn.Do("select", "0") // 进入数据库0
	if err != nil {
		return false
	}
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
	result, err := conn.Do("select", "0") // 进入数据库0
	if err != nil {
		return false
	}
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

func GetTodayQuestion(date time.Time, backup time.Time) (questions []*mysqlDao.Question, err error) {
	questions = []*mysqlDao.Question{}
	var mid []*base_service.Question
	conn := p.Get()
	defer conn.Close()
	result, err := conn.Do("select", "1") // 进入数据库1
	if err != nil {
		questions = nil
		return nil, err
	}
	defer func(conn redigo.Conn) {
		err := conn.Close()
		if err != nil {
			log.Println(err)
		}
	}(conn)
	result, err = conn.Do("GET", date.Format("2006-01-02"))
	if result == nil {
		s, err := conn.Do("GET", backup.Format("2006-01-02"))
		if s == nil {
			
			s, err = conn.Do("GET", "2022-01-17")
			log.Printf("0117 is %v", s.([]uint8))
		}
		_, err = conn.Do("SET", date.Format("2006-01-02"), string(s.([]uint8)))
		if err != nil {
			return nil, err
		}
		return nil, err
	}
	log.Printf("GetTodayQuestion result is %v", string(result.([]uint8)))
	err = json.Unmarshal([]byte(result.([]uint8)), &mid)
	if err != nil {
		log.Printf("json.unmarshal for %v failed by %v", result, err)
		return nil, err
	}
	for i, question := range mid {
		questions = append(questions, &mysqlDao.Question{
			QuestionID: int32(i),
			Type:       int32(question.Type),
			Content:    question.Content,
			Answers:    strings.Join(question.Answers, "\t"),
			Depends:    strings.Join(question.Depends, "\t"),
		})
	}
	return questions, nil
}
func PrepareTomorrowQuestion(questions []*base_service.Question) error {
	conn := p.Get()
	defer conn.Close()
	result, err := conn.Do("select", "1") // 进入数据库1
	if err != nil {
		return errors.New("cannot connect to redis")
	}
	defer func(conn redigo.Conn) {
		err := conn.Close()
		if err != nil {
			fmt.Println(err)
		}
	}(conn)
	bytes, err := json.Marshal(questions)
	if err != nil {
		return err
	}
	now := time.Now().AddDate(0, 0, 1)
	timeFormat := now.Format("2006-01-02")
	result, err = conn.Do("SET", timeFormat, string(bytes))
	if err != nil {
		return err
	}
	log.Printf("PrepareTomorrowQuestion result is %v", result)
	return nil
}

func PublishNotice(userID string, msg string, ttl int) bool {
	conn := p.Get()
	defer conn.Close()
	result, err := conn.Do("select", "2") // 进入数据库2
	if err != nil {
		return false
	}
	result, _ = conn.Do("SET", userID, msg)
	r := result.(string) == "OK"
	if !r {
		return false
	}
	if ttl == -1 {
		return true
	}
	result, _ = conn.Do("EXPIRE", userID, ttl)
	r = result.(int64) == 1
	if !r {
		return false
	}
	return true
}
