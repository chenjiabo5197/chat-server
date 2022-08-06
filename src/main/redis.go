package main

import (
	"github.com/garyburd/redigo/redis"
	"time"
)

//定义一个全局的pool,传给UserDao实例对象
var pool *redis.Pool

func initPool(address string, maxIdle int, maxActive int, idleTimeout time.Duration) {
	pool = &redis.Pool {
		MaxIdle : maxIdle,  //最大空闲的连接数
		MaxActive : maxActive, //表示和redis的最大连接数，0表示无限制
		IdleTimeout : idleTimeout, //最大空闲时间
		Dial : func() (redis.Conn, error) { //初始化连接的代码，连接哪个redis
			return redis.Dial("tcp", address)
		},
	}
}


