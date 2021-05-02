package redisdb

import (
	"log"
	"math"
	"time"

	"github.com/chynkm/ratesdb/currencystore"
	"github.com/chynkm/ratesdb/datastore"
	"github.com/gomodule/redigo/redis"
)

// SaveExchangeRates to Redis
// Insert data of previous days exchange rate when a day is missing
// Previous day is always present since we are fetching it from the DB
func SaveExchangeRates() {
	rdb := Rdbpool.Get()
	defer rdb.Close()

	dbCurrencies := datastore.GetCurrencies()

	for i := Days; i >= 0; i-- {
		date := time.Now().AddDate(0, 0, -i).Format(currencystore.DateLayout)
		if i == 0 {
			LatestDate = date
		}
		exchangeRates := datastore.GetExchangeRates(date)

		dailyExchangeRates := createExchangeRateHash(
			date,
			dbCurrencies,
			exchangeRates,
		)

		for key, dailyExchangeRate := range dailyExchangeRates {
			redisExchangeRates := []interface{}{}
			redisExchangeRates = append(redisExchangeRates, key)
			for code, rate := range dailyExchangeRate {
				redisExchangeRates = append(redisExchangeRates, code, rate)
			}

			_, err := rdb.Do("HMSET", redisExchangeRates...)
			if err != nil {
				log.Fatal(err)
			}
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
func GetExchangeRate(date, from, to string) float64 {
	rdb := Rdbpool.Get()
	defer rdb.Close()

	rate, err := redis.Float64(rdb.Do("HGET", date+":"+from, to))
	if err != nil {
		log.Println("redis: unable to retrieve exchange rate for: ", date, from, to)
		return 0
	}

	return rate
}
