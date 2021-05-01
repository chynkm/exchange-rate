package redisdb

import "github.com/gomodule/redigo/redis"

var Rdbpool *redis.Pool

const (
	euro = "EUR"
	days = 31
)
