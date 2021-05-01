package redisdb

import (
	"log"
	"math"
	"time"

	"github.com/chynkm/ratesdb/currencystore"
	"github.com/chynkm/ratesdb/datastore"
)

const (
	euro = "EUR"
	days = 31
)

// SaveExchangeRates to Redis
func SaveExchangeRates() {
	rdb := Rdbpool.Get()
	defer rdb.Close()

	dbCurrencies := datastore.GetCurrencies()

	for i := 0; i <= days; i++ {
		currentDate := time.Now().AddDate(0, 0, -i).Format(currencystore.DateLayout)
		exchangeRates := datastore.GetExchangeRates(currentDate)

		dailyExchangeRates := createExchangeRateHash(
			currentDate,
			dbCurrencies,
			exchangeRates,
		)

		for key, dailyExchangeRate := range dailyExchangeRates {
			redisExchangeRates := []interface{}{}
			redisExchangeRates = append(redisExchangeRates, key)
			for code, rate := range dailyExchangeRate {
				redisExchangeRates = append(redisExchangeRates, code, rate)
			}

			if len(redisExchangeRates) > 1 {
				_, err := rdb.Do("HMSET", redisExchangeRates...)
				if err != nil {
					log.Fatal(err)
				}
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
		exchangeRateHash[date+":"+currencyCode] = getExchangeRate(
			currencyCode,
			exchangeRates,
		)
	}

	return exchangeRateHash
}

// getExchangeRate for a base currency
func getExchangeRate(
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
