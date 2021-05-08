package redisdb

import (
	"log"
	"strconv"
	"time"

	"github.com/gomodule/redigo/redis"
)

const (
	apiUserPrefix     = "api_user:"
	maxRequestsPerMin = 100
	fixedWindow       = 60 // start of minute to end of minute
)

// incrementKey increases the count by 1 if key is present.
// Creates a key, sets the value to 1 with an expiry if the key is missing.
func incrementKey(rdb redis.Conn, key string, seconds int) {
	_, err := rdb.Do("INCR", key)
	if err != nil {
		log.Fatal("redis: unable to set count for key: ", key)
	}

	if seconds > 0 {
		_, err = rdb.Do("EXPIRE", key, seconds)
		if err != nil {
			log.Fatal("redis: unable to set expiry for key: ", key)
		}
	}
}

// getKeyCount retrieves the key count
func getKeyCount(rdb redis.Conn, key string) int {
	exists, err := redis.Int(rdb.Do("EXISTS", key))
	if err != nil {
		log.Fatal("redis: check key exists failed. key: ", key)
	}

	if exists == 0 {
		return 0
	}

	count, err := redis.Int(rdb.Do("GET", key))
	if err != nil {
		log.Fatal("redis: unable to get count for key: ", key)
	}

	return count
}

// getKeyCount creates the key with prefix and suffix
func createKey(ip string, t time.Time) string {
	return apiUserPrefix + ip + ":" + strconv.Itoa(t.Minute())
}

// AllowAPIRequest determines whether an API request should be allowed or not
func AllowAPIRequest(ip string) bool {
	rdb := Rdbpool.Get()
	defer rdb.Close()

	t := time.Now()
	key := createKey(ip, t)

	n := getKeyCount(rdb, key)
	if n >= maxRequestsPerMin {
		return false
	}

	seconds := 0
	if n == 0 {
		seconds = fixedWindow - t.Second()
	}
	incrementKey(rdb, key, seconds)
	return true
}
