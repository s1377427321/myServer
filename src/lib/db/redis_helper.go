package db

import (
	"github.com/gomodule/redigo/redis"
	"time"
	"github.com/astaxie/beego"
)

var pool *redis.Pool

func InitRedis(server,password string)  {
	pool = &redis.Pool{
		MaxIdle:     500,
		IdleTimeout: 3600 * time.Second,
		Dial: func() (redis.Conn, error) {
			c, err := redis.Dial("tcp", server)
			if err != nil {
				beego.Error(err)
				return nil, err
			}

			if password != "" {
				if _, err := c.Do("AUTH", password); err != nil {
					beego.Error(err)
					c.Close()
					return nil, err
				}
			}

			return c, err
		},
	}
}

func GetRedis() redis.Conn {
	return pool.Get()
}

func Send(cmd string, args ...interface{}) error {
	red := GetRedis()

	if red.Err() != nil {
		beego.Error("Get connection from redis-pool failed:", red.Err())
		return red.Err()
	}
	defer red.Close()

	err := red.Send(cmd, args...)
	if err != nil {
		beego.Error("Execute Send request to redis failed:", err)
		return err
	}

	err = red.Flush()
	if err != nil {
		beego.Error("Execute Flush to redis failed:", err)
		return err
	}

	return nil
}

func Do(cmd string, args ...interface{}) (interface{}, error) {
	red := GetRedis()

	if red.Err() != nil {
		beego.Error("Get connection from redis-pool failed:", red.Err())
		return nil, red.Err()
	}
	defer red.Close()

	reply, err := red.Do(cmd, args...)
	if err != nil {
		beego.Error("Execute Do request to redis failed:", err)
		return nil, err
	}

	return reply, nil
}

