package redisdb

import "github.com/gomodule/redigo/redis"

var (
	Rdbpool    *redis.Pool
	LatestDate string
)

const (
	euro = "EUR"
	Days = 30
)
