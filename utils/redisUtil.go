package utils

import (
	"fmt"
	"github.com/garyburd/redigo/redis"
)

func RedisUtilset(phoneNumber string, verifyCode string) {
	fmt.Println("start redis")
	conn, err := redis.Dial("tcp", "127.0.0.1:6379")
	if err != nil {
		fmt.Println("connect redis error :", err)
		return
	}
	defer conn.Close()
	_, err = conn.Do("SET", phoneNumber, verifyCode)
	_, err = conn.Do("expire", phoneNumber, 300)
	if err != nil {
		fmt.Println("redis set error:", err)
	}
}

func RedisUtilget(phoneNumber string) string {
	conn, err := redis.Dial("tcp", "127.0.0.1:6379")
	if err != nil {
		fmt.Println("connect redis error :", err)
	}
	defer conn.Close()
	name, err := redis.String(conn.Do("GET", phoneNumber))

	if err != nil {
		fmt.Println("redis get error:", err)
	} else {
		fmt.Printf("Got name:%s", name)
	}
	return name
}

var Pool redis.Pool
func init() { //init 用于初始化一些参数，先于main执行
	Pool = redis.Pool{
		MaxIdle:     16,  //最大的空闲连接数
		MaxActive:   32,  //最大的激活连接数
		IdleTimeout: 120, //空闲连接等待时间
		Dial: func() (redis.Conn, error) {
			return redis.Dial("tcp", "localhost:6379")
		},
	}
}

func GetClientMes(key, arg string) string {
	conn := Pool.Get()
	data, _ := conn.Do("HGET", key, arg)
	defer conn.Close()
	return fmt.Sprintf("%s", data)
}
