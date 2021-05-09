package redisdb

import (
	"log"
	"math"
	"time"

	"github.com/chynkm/ratesdb/currencystore"
	"github.com/chynkm/ratesdb/datastore"
	"github.com/gomodule/redigo/redis"
)

const (
	euro        = "EUR"
	rate_prefix = "rate:"
	days        = 30 // maximum days exchange rate to load inside Redis at startup
)

var LatestDate string

// SaveExchangeRates to Redis
// Insert data of previous days exchange rate when a day is missing
// Previous day is always present since we are fetching it from the DB
func SaveExchangeRates() {
	rdb := Rdbpool.Get()
	defer rdb.Close()

	// Flush the Redis data before fresh data is inserted
	_, err := rdb.Do("FLUSHDB")
	if err != nil {
		log.Fatal(err)
	}

	for i := days; i >= 0; i-- {
		date := time.Now().AddDate(0, 0, -i).Format(currencystore.DateLayout)
		if i == 0 {
			LatestDate = date
		}
		saveExchangeRateForDate(rdb, date)
	}
}

// saveExchangeRateForDate save exchange rate to Redis for a single date
func saveExchangeRateForDate(rdb redis.Conn, date string) {
	dbCurrencies := datastore.GetCurrencies()

	exchangeRates := datastore.GetExchangeRates(date)

	dailyExchangeRates := createExchangeRateHash(
		date,
		dbCurrencies,
		exchangeRates,
	)

	for key, dailyExchangeRate := range dailyExchangeRates {
		redisExchangeRates := []interface{}{}
		redisExchangeRates = append(redisExchangeRates, rate_prefix+key)
		for code, rate := range dailyExchangeRate {
			redisExchangeRates = append(redisExchangeRates, code, rate)
		}

		_, err := rdb.Do("HMSET", redisExchangeRates...)
		if err != nil {
			log.Fatal(err)
		}
	}
}

// createExchangeRateHash generates hash for exchange rates of different
// currencies for the given date
func createExchangeRateHash(
	date string,
	currencies map[string]int,
	exchangeRates map[string]float64,
) map[string]map[string]float64 {
	exchangeRateHash := map[string]map[string]float64{}

	for currencyCode, _ := range currencies {
		exchangeRateHash[date+":"+currencyCode] = calculateExchangeRate(
			currencyCode,
			exchangeRates,
		)
	}

	return exchangeRateHash
}

// calculateExchangeRate for a base currency
func calculateExchangeRate(
	baseCurrencyCode string,
	exchangeRates map[string]float64,
) map[string]float64 {
	if baseCurrencyCode == euro {
		return exchangeRates
	}

	baseCurrencyExchangeRate := map[string]float64{}
	for currencyCode, rate := range exchangeRates {
		roundRate := rate / exchangeRates[baseCurrencyCode]
		roundRate = math.Round(roundRate*10000) / 10000
		baseCurrencyExchangeRate[currencyCode] = roundRate
	}

	return baseCurrencyExchangeRate
}

// GetExchangeRate retrieves the rate for the day from Redis
func GetExchangeRate(date, from, to string) (map[string]interface{}, error) {
	rdb := Rdbpool.Get()
	defer rdb.Close()

	key := rate_prefix + date + ":" + from

	exists, err := redis.Int(rdb.Do("EXISTS", key))
	if err != nil {
		log.Println("redis: check key exists failed in GetExchangeRate. key: ", key)
		return map[string]interface{}{}, err
	}

	if exists == 0 {
		saveExchangeRateForDate(rdb, date)
	}

	if to != "" {
		rate, err := redis.Float64(rdb.Do("HGET", key, to))
		if err != nil {
			log.Println("redis: unable to retrieve HGET exchange rate for: ", date, from, to)
			return map[string]interface{}{}, err
		}

		return map[string]interface{}{to: rate}, nil
	}

	redisRates, err := redis.Values(rdb.Do("HGETALL", key))
	if err != nil {
		log.Println("redis: unable to retrieve HGETALL exchange rate for: ", date, from)
		return map[string]interface{}{}, err
	}

	rates := map[string]interface{}{}
	for i := 0; i < len(redisRates); i += 2 {
		currencyCode, _ := redis.String(redisRates[i], nil)
		exchangeRate, _ := redis.Float64(redisRates[i+1], nil)

		if exchangeRate == 0 {
			rates[currencyCode] = nil
		} else {
			rates[currencyCode] = exchangeRate
		}
	}

	return rates, err
}
