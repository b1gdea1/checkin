package redis

import (
	"fmt"
	"github.com/gomodule/redigo/redis"
	"testing"
)

func Test1(t *testing.T) {
	conn, err := redis.Dial("tcp", "39.99.39.103:6379")
	if err != nil {
		fmt.Println(err)
		return
	}
	defer conn.Close()
	// fmt.Println(conn)
	result, err := conn.Do("auth", "TGYtgy1234")
	if err != nil {
		
		return
	}
	fmt.Printf("%T\n", result)
	fmt.Println(result)
	reslut, err := conn.Do("set", "name", "tgy")
	if err != nil {
		return
	}
	fmt.Println(reslut)
	// result, _ = conn.Do("get", "name")
	result, err = redis.String(conn.Do("get", "name"))
	if err != nil {
		return
	}
	fmt.Println(result)
}
