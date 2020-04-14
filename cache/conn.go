package cache

import (
	"github.com/gomodule/redigo/redis"
	"time"
	"fmt"
)

var(
	pool *redis.Pool
	redisHost = "113.31.118.136:6379"
	redisPass = "myredis"
)

func newRedisPool() *redis.Pool{
	return &redis.Pool{
		MaxIdle:50,
		MaxActive:30,
		IdleTimeout:time.Second * 300,
		Dial: func() (redis.Conn, error) {
			//1. 打开链接
			c,err := redis.Dial("tcp",redisHost)
			if err != nil {
				fmt.Printf("dail redis error: %s \n",err.Error())
				return nil,err
			}
			//2. 访问认证
			if _,err := c.Do("AUTH",redisPass);err != nil {
				c.Close()
				fmt.Printf("DO AUTH error: %s \n",err.Error())
				return nil,err
			}
			return c,nil
		},
		TestOnBorrow: func(c redis.Conn, t time.Time) error {
			if time.Since(t) < time.Minute {
				return nil
			}
			_,err := c.Do("PING")
			return err
		},
	}
}

func init() {
	pool = newRedisPool()
}

func RedisPool() *redis.Pool{
	return pool
}