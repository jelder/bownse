package main

import (
	"github.com/garyburd/redigo/redis"
	. "github.com/jelder/goenvmap"
	"time"
)

var (
	RedisPool *redis.Pool
)

func init() {
	RedisPool = &redis.Pool{
		MaxIdle:     3,
		IdleTimeout: 240 * time.Second,
		Dial: func() (c redis.Conn, err error) {
			c, err = redis.DialURL(ENV["REDIS_URL"])
			if err != nil {
				panic(err)
			}
			return c, err
		},
		TestOnBorrow: func(c redis.Conn, t time.Time) error {
			_, err := c.Do("PING")
			return err
		},
	}
}
